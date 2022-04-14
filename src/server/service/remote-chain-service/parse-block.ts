import { BigNumber } from "ethers";
import { Block } from "@/model/block";
import {
  hexToBuffer,
  hexToNumber,
} from "@/server/service/remote-chain-service/hex-utils";
import { logger } from "onefx/lib/integrated-gateways/logger";

type Hex = string;

export type RawTransaction = {
  blockHash: Hex;
  blockNumber: Hex;
  from: Hex;
  gas: Hex;
  gasPrice: Hex;
  hash: Hex;
  input: Hex;
  nonce: Hex;
  to?: Hex;
  transactionIndex: Hex;
  value: Hex;
  type: Hex;
  v: Hex;
  r: Hex;
  s: Hex;
  maxFeePerGas?: Hex;
  maxPriorityFeePerGas?: Hex;
};

export type RawBlock = {
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

export function parseBlock(block: RawBlock): Block | null {
  try {
    return {
      hash: hexToBuffer(block.hash),
      consensus: true, // TODO(dora): ?
      difficulty: BigNumber.from(block.difficulty).toString(),
      gasLimit: BigNumber.from(block.gasLimit).toString(),
      gasUsed: BigNumber.from(block.gasUsed).toString(),
      nonce: hexToBuffer(block.nonce),
      number: hexToNumber(block.number),
      size: BigNumber.from(block.size).toNumber(),
      timestamp: new Date(BigNumber.from(block.timestamp).toNumber() * 1000),
      totalDifficulty: BigNumber.from(block.totalDifficulty).toString(),
      isEmpty: !!block.transactions.length,
      // TODO(dora): need to get it from ethers.js
      baseFeePerGas: "0",
      parentHash: hexToBuffer(block.parentHash),
      miner: hexToBuffer(block.miner),
    };
  } catch (err) {
    logger.error(
      `failed to parse the block: ${JSON.stringify(block)}: ${err} ${
        err instanceof Error && err.stack
      }`
    );
    return null;
  }
}
