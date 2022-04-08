import { MyServer } from "@/server/start-server";
import { getRanges } from "@/server/service/utils/get-ranges";
import { BigNumber, ethers } from "ethers";
import { Block } from "@/model/block";
import { Transaction } from "@/model/transaction";
import { Address } from "@/model/address";
import { logger } from "onefx/lib/integrated-gateways/logger";

function hexToBuffer(hex: string): Buffer {
  return Buffer.from(hex.slice(2), "hex");
}

function hexToNumber(hex: string): number {
  return BigNumber.from(hex).toNumber();
}

function hexToDecimal(hex: string): string {
  return BigNumber.from(hex).toString();
}

type Hex = string;

type RawTransaction = {
  blockHash: Hex;
  blockNumber: Hex;
  from: Hex;
  gas: Hex;
  gasPrice: Hex;
  hash: Hex;
  input: Hex;
  nonce: Hex;
  to: Hex;
  transactionIndex: Hex;
  value: Hex;
  type: Hex;
  v: Hex;
  r: Hex;
  s: Hex;
};

type RawBlock = {
  difficulty: Hex;
  extraData: Hex;
  gasLimit: Hex;
  gasUsed: Hex;
  hash: Hex;
  logsBloom: Hex;
  miner: Hex;
  mixHash: Hex;
  nonce: Hex;
  number: Hex;
  parentHash: Hex;
  receiptsRoot: Hex;
  sha3Uncles: Hex;
  size: Hex;
  stateRoot: Hex;
  timestamp: Hex;
  totalDifficulty: Hex;
  transactions: RawTransaction[];
  transactionsRoot: Hex;
  uncles: Hex[];
};

export class RemoteChainService {
  server: MyServer;

  constructor(server: MyServer) {
    this.server = server;
  }

  async missingBlockNumberRanges(
    first: number,
    last: number
  ): Promise<Array<[number, number]>> {
    const rangeMin = Math.min(first, last);
    const rangeMax = Math.max(first, last);

    const missingBlockNumbers = await this.server.gateways.dbCon.manager.query(
      `(SELECT distinct b1.number FROM generate_series(($1)::integer, ($2)::integer) AS b1(number) WHERE NOT EXISTS (SELECT 1 FROM block b2 WHERE b2.number=b1.number AND b2.consensus) ORDER BY b1.number ASC)`,
      [rangeMin, rangeMax]
    );
    const numberList = missingBlockNumbers.map(
      (it: { number: number }) => it.number
    );
    return getRanges(numberList);
  }

  async getRawBlock(i: number): Promise<RawBlock> {
    const rawBlock = await this.server.gateways.chainProvider.send(
      "eth_getBlockByNumber",
      [ethers.utils.hexValue(i), true]
    );

    return rawBlock;
  }

  async fetchBlockByRange(
    end: number,
    start: number
  ): Promise<{
    blocks: Array<Block>;
    transactions: Array<Transaction>;
    addresses: Array<Address>;
  }> {
    const min = Math.min(start, end);
    const max = Math.max(start, end);

    const transactions: Array<Transaction> = [];
    const addressMap: Record<string, number> = {};
    // TODO(dora): should optimize
    const blockBatch: Array<Block> = [];

    for (let i = max; i >= min; i -= 1) {
      try {
        // eslint-disable-next-line no-await-in-loop
        const block = await this.getRawBlock(i);
        const blockNumber = hexToNumber(block.number);
        const timestamp = new Date(
          BigNumber.from(block.timestamp).toNumber() * 1000
        );
        if (!addressMap[block.miner] || addressMap[block.miner] < blockNumber) {
          addressMap[block.miner] = blockNumber;
        }

        const processingTime = new Date();
        let cumulativeGasUsed = BigNumber.from(0);

        for (const tx of block.transactions) {
          const gas = BigNumber.from(tx.gas);
          cumulativeGasUsed = cumulativeGasUsed.add(gas);

          if (!addressMap[tx.from] || addressMap[tx.from] < blockNumber) {
            addressMap[tx.from] = blockNumber;
          }
          if (!addressMap[tx.to] || addressMap[tx.to] < blockNumber) {
            addressMap[tx.to] = blockNumber;
          }

          transactions.push({
            block: hexToBuffer(block.hash),
            hash: hexToBuffer(tx.hash),
            blockNumber,
            cumulativeGasUsed: cumulativeGasUsed.toString(),
            earliestProcessingStart: processingTime,
            error: "", // TODO(dora): how to?
            gas: gas.toString(),
            gasPrice: hexToDecimal(tx.gasPrice),
            gasUsed: "0", // TODO(dora): how to?
            index: hexToNumber(tx.transactionIndex),
            // createdContractCodeIndexedAt: "", // TODO(dora)
            input: hexToBuffer(tx.input),
            nonce: hexToNumber(tx.nonce),
            r: hexToDecimal(tx.r),
            s: hexToDecimal(tx.s),
            status: 1, // TODO(dora)
            v: hexToDecimal(tx.v),
            value: hexToDecimal(tx.value),
            // revertReason: "", // TODO(dora)
            // maxPriorityFeePerGas: "", // TODO(dora)
            // maxFeePerGas: "", // TODO(dora)
            // type: "", // TODO(dora)
            fromAddress: hexToBuffer(tx.from),
            fromAddressHash: hexToBuffer(tx.from),
            toAddress: hexToBuffer(tx.to),
            toAddressHash: hexToBuffer(tx.to),
            timestamp,
          });
        }

        const bBlock: Block = {
          hash: hexToBuffer(block.hash),
          consensus: true, // TODO(dora): ?
          difficulty: BigNumber.from(block.difficulty).toString(),
          gasLimit: BigNumber.from(block.gasLimit).toString(),
          gasUsed: BigNumber.from(block.gasUsed).toString(),
          nonce: hexToBuffer(block.nonce),
          number: hexToNumber(block.number),
          size: BigNumber.from(block.size).toNumber(),
          timestamp,
          totalDifficulty: BigNumber.from(block.totalDifficulty).toString(),
          isEmpty: !!block.transactions.length,
          // TODO(dora): need to get it from ethers.js
          baseFeePerGas: "0",
          parentHash: hexToBuffer(block.parentHash),
          miner: hexToBuffer(block.miner),
        };

        blockBatch.push(bBlock);
      } catch (err) {
        logger.error(`failed to fetch block ${i}`);
      }
    }

    const addresses = Object.keys(addressMap).map((addressHash: string) => ({
      hash: hexToBuffer(addressHash),
      fetchedCoinBalanceBlockNumber: addressMap[addressHash],
      coinBalancePromise: this.server.gateways.chainProvider.getBalance(
        addressHash,
        addressMap[addressHash]
      ),
    }));

    const balances = await Promise.all(
      addresses.map((address) => address.coinBalancePromise)
    );

    const resolvedAddresses: Array<Address> = addresses.map((address, i) => ({
      hash: address.hash,
      fetchedCoinBalanceBlockNumber: address.fetchedCoinBalanceBlockNumber,
      fetchedCoinBalance: balances[i].toString(),
      decompiled: false,
      verified: false,
    }));

    return {
      blocks: blockBatch,
      transactions,
      addresses: resolvedAddresses,
    };
  }
}
