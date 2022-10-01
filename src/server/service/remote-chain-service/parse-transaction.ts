import { RawTransaction } from "@/server/service/remote-chain-service/parse-block";
import { Transaction } from "@/model/transaction";
import {
  hexToBuffer,
  hexToDecimal,
  hexToNumber,
} from "@/server/service/remote-chain-service/hex-utils";
import { logger } from "onefx/lib/integrated-gateways/logger";

export function parseTransaction(
  tx: RawTransaction,
  externalProps: {
    earliestProcessingStart: Date;
    blockNumber: number;
    timestamp: Date;
    gas: string;
    block: Buffer;
    cumulativeGasUsed: string;
  }
): Transaction | null {
  try {
    return {
      hash: hexToBuffer(tx.hash),
      error: "", // TODO(dora): how to?
      gasPrice: hexToDecimal(tx.gasPrice),
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
      maxPriorityFeePerGas:
        tx.maxPriorityFeePerGas && hexToDecimal(tx.maxPriorityFeePerGas),
      maxFeePerGas: tx.maxFeePerGas && hexToDecimal(tx.maxFeePerGas),
      type: tx.type ? hexToNumber(tx.type) : undefined,
      fromAddress: hexToBuffer(tx.from),
      fromAddressHash: hexToBuffer(tx.from),
      // @ts-ignore
      toAddress: tx.to && hexToBuffer(tx.to), // TODO(tian) is null processed well?
      // @ts-ignore
      toAddressHash: tx.to && hexToBuffer(tx.to), // is null processed well?
      ...externalProps,
    };
  } catch (err) {
    logger.error(
      `failed to parse the transaction: ${JSON.stringify(tx)}: ${err} ${
        err instanceof Error && err.stack
      }`
    );
    return null;
  }
}
