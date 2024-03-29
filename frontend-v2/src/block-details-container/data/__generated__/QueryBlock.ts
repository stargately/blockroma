/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL query operation: QueryBlock
// ====================================================

export interface QueryBlock_block {
  consensus: boolean | null;
  difficulty: string | null;
  gasLimit: string | null;
  gasUsed: string | null;
  hash: any | null;
  miner: any | null;
  nonce: any | null;
  number: number | null;
  parentHash: any | null;
  size: number | null;
  timestamp: any | null;
  totalDifficulty: string | null;
  numTxs: number | null;
}

export interface QueryBlock {
  block: QueryBlock_block | null;
}

export interface QueryBlockVariables {
  blockNumber: number;
}
