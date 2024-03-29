/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL query operation: Tokens
// ====================================================

export interface Tokens_tokens_edges_node {
  name: string | null;
  symbol: string | null;
  totalSupply: string | null;
  decimals: string | null;
  type: string | null;
  contractAddress: any | null;
  skipMetadata: boolean | null;
}

export interface Tokens_tokens_edges {
  cursor: string | null;
  node: Tokens_tokens_edges_node | null;
}

export interface Tokens_tokens_pageInfo {
  startCursor: string | null;
  endCursor: string | null;
  hasNextPage: boolean | null;
  hasPreviousPage: boolean | null;
}

export interface Tokens_tokens {
  edges: (Tokens_tokens_edges | null)[] | null;
  pageInfo: Tokens_tokens_pageInfo | null;
}

export interface Tokens {
  tokens: Tokens_tokens | null;
}

export interface TokensVariables {
  after?: string | null;
  first?: number | null;
  before?: string | null;
  last?: number | null;
  symbol?: string | null;
}
