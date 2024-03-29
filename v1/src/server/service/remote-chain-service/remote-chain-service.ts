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
import {
  matchTokenTransferInput,
  parseTokenTransfers,
} from "@/server/service/remote-chain-service/parse-token-transfer/parse-token-transfer";
import { TokenTransfer } from "@/model/token-transfer";
import { logger } from "onefx/lib/integrated-gateways/logger";

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

  async getRawBlock(i: number): Promise<RawBlock | null> {
    let rawBlock: RawBlock | null = null;
    for (let t = 0; t < 3; t += 1) {
      // eslint-disable-next-line no-await-in-loop
      rawBlock = await this.server.gateways.chainProvider
        .get()
        .send("eth_getBlockByNumber", [ethers.utils.hexValue(i), true]);
      if (rawBlock) {
        break;
      }
      logger.warn(`failed to fetch raw block ${i} for sequence ${t}`);
    }

    return rawBlock;
  }

  async fetchBlockByRange(
    end: number,
    start: number
  ): Promise<{
    blocks: Array<Block>;
    transactions: Array<Transaction>;
    addresses: Array<Address>;
    tokenTransfers: TokenTransfer[];
  }> {
    const min = Math.min(start, end);
    const max = Math.max(start, end);

    const transactions: Array<Transaction> = [];
    const addressMap: Record<string, number> = {};
    // TODO(dora): should optimize
    const blockBatch: Array<Block> = [];
    let tokenTransfers: TokenTransfer[] = [];

    for (let i = max; i >= min; i -= 1) {
      // eslint-disable-next-line no-await-in-loop
      const rawBlock = await this.getRawBlock(i);
      if (!rawBlock) {
        logger.warn(`failed to fetch raw block ${i}; skipping...`);
        // eslint-disable-next-line no-continue
        continue;
      }

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

          // reduce remote calls by checking input signature. only potential transfer call is worth fetching receipt.
          if (matchTokenTransferInput(tx.input)) {
            try {
              const receipt =
                // eslint-disable-next-line no-await-in-loop
                await this.server.gateways.chainProvider
                  .get()
                  .getTransactionReceipt(tx.hash);
              const { tokenTransfers: tt } = parseTokenTransfers(receipt.logs);
              tokenTransfers = [...tokenTransfers, ...tt];
              for (const t of tokenTransfers) {
                t.toAddress &&
                  maybeAddToAddress(
                    addressMap,
                    blockNumber,
                    `0x${t.toAddress.toString("hex")}`
                  );
                t.fromAddress &&
                  maybeAddToAddress(
                    addressMap,
                    blockNumber,
                    `0x${t.fromAddress.toString("hex")}`
                  );
                t.tokenContractAddress &&
                  maybeAddToAddress(
                    addressMap,
                    blockNumber,
                    `0x${t.tokenContractAddress.toString("hex")}`
                  );
              }
            } catch (err) {
              logger.error(`failed to parse transfer: ${err}`);
            }
          }
        }
      }
    }

    const addresses = await this.fetchAddressBalances(addressMap);

    return {
      blocks: blockBatch,
      transactions,
      addresses,
      tokenTransfers,
    };
  }

  private async fetchAddressBalances(
    addressToBlockNumber: Record<string, number>
  ): Promise<Array<Address>> {
    const addresses = Object.keys(addressToBlockNumber).map(
      (addressHash: string) => ({
        hash: hexToBuffer(addressHash),
        fetchedCoinBalanceBlockNumber: addressToBlockNumber[addressHash],
        coinBalancePromise: this.server.gateways.chainProvider
          .get()
          .getBalance(addressHash, addressToBlockNumber[addressHash]),
      })
    );

    // TODO(tian) should handle error
    const toSettle = await Promise.allSettled(
      addresses.map((address) => address.coinBalancePromise)
    );
    const balances = toSettle.map((t) => {
      if (t.status !== "rejected") {
        return t.value;
      }
      return undefined;
    });

    return addresses.map((address, i) => ({
      hash: address.hash,
      fetchedCoinBalanceBlockNumber: address.fetchedCoinBalanceBlockNumber,
      fetchedCoinBalance: (balances[i] ?? 0).toString(),
      decompiled: false,
      verified: false,
    }));
  }
}
