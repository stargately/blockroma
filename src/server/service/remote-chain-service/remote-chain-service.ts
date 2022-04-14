import { MyServer } from "@/server/start-server";
import { getRanges } from "@/server/service/utils/get-ranges";
import { BigNumber, ethers } from "ethers";
import { Block } from "@/model/block";
import { Transaction } from "@/model/transaction";
import { Address } from "@/model/address";
import {
  parseBlock,
  RawBlock,
} from "@/server/service/remote-chain-service/parse-block";
import {
  hexToBuffer,
  hexToNumber,
} from "@/server/service/remote-chain-service/hex-utils";
import { parseTransaction } from "@/server/service/remote-chain-service/parse-transaction";

function maybeAddToAddress(
  addressMap: Record<string, number>,
  curBlockNumber: number,
  address?: string
) {
  if (!address) {
    return;
  }
  if (!addressMap[address] || addressMap[address] < curBlockNumber) {
    addressMap[address] = curBlockNumber;
  }
}

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
      // eslint-disable-next-line no-await-in-loop
      const rawBlock = await this.getRawBlock(i);
      const blockNumber = hexToNumber(rawBlock.number);

      maybeAddToAddress(addressMap, blockNumber, rawBlock.miner);

      const processingTime = new Date();
      let cumulativeGasUsed = BigNumber.from(0);

      const block = parseBlock(rawBlock);
      if (block) {
        blockBatch.push(block);

        for (const tx of rawBlock.transactions) {
          const gas = BigNumber.from(tx.gas);
          cumulativeGasUsed = cumulativeGasUsed.add(gas);

          maybeAddToAddress(addressMap, blockNumber, tx.from);
          maybeAddToAddress(addressMap, blockNumber, tx.to);

          const transaction = parseTransaction(tx, {
            block: hexToBuffer(rawBlock.hash),
            gas: gas.toString(),
            blockNumber: block.number,
            earliestProcessingStart: processingTime,
            timestamp: new Date(
              BigNumber.from(rawBlock.timestamp).toNumber() * 1000
            ),
            cumulativeGasUsed: cumulativeGasUsed.toString(),
          });
          if (transaction) {
            transactions.push(transaction);
          }
        }
      }
    }

    const addresses = await this.fetchAddressBalances(addressMap);

    return {
      blocks: blockBatch,
      transactions,
      addresses,
    };
  }

  private async fetchAddressBalances(
    addressToBlockNumber: Record<string, number>
  ): Promise<Array<Address>> {
    const addresses = Object.keys(addressToBlockNumber).map(
      (addressHash: string) => ({
        hash: hexToBuffer(addressHash),
        fetchedCoinBalanceBlockNumber: addressToBlockNumber[addressHash],
        coinBalancePromise: this.server.gateways.chainProvider.getBalance(
          addressHash,
          addressToBlockNumber[addressHash]
        ),
      })
    );

    const balances = await Promise.all(
      addresses.map((address) => address.coinBalancePromise)
    );

    return addresses.map((address, i) => ({
      hash: address.hash,
      fetchedCoinBalanceBlockNumber: address.fetchedCoinBalanceBlockNumber,
      fetchedCoinBalance: balances[i].toString(),
      decompiled: false,
      verified: false,
    }));
  }
}
