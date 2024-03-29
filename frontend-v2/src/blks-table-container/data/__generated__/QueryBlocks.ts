/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL query operation: QueryBlocks
// ====================================================

export interface QueryBlocks_blocks_edges_node {
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

export interface QueryBlocks_blocks_edges {
  node: QueryBlocks_blocks_edges_node | null;
}

export interface QueryBlocks_blocks_pageInfo {
  hasNextPage: boolean | null;
  endCursor: string | null;
  startCursor: string | null;
  hasPreviousPage: boolean | null;
}

export interface QueryBlocks_blocks {
  edges: (QueryBlocks_blocks_edges | null)[] | null;
  pageInfo: QueryBlocks_blocks_pageInfo | null;
}

export interface QueryBlocks {
  blocks: QueryBlocks_blocks | null;
}

export interface QueryBlocksVariables {
  first?: number | null;
  after?: string | null;
}
