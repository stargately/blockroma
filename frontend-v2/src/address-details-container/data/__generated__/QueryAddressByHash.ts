/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL query operation: QueryAddressByHash
// ====================================================

import {Status} from "@/shared/__generated__/globalTypes";

export interface QueryAddressByHash_address_transactions_pageInfo {
  hasNextPage: boolean | null;
  startCursor: string | null;
  hasPreviousPage: boolean | null;
  endCursor: string | null;
}

export interface QueryAddressByHash_address_transactions_edges_node {
  id: string;
  timestamp: any | null;
  hash: any | null;
  blockNumber: number | null;
  value: string | null;
  gasUsed: string | null;
  cumulativeGasUsed: string | null;
  error: string | null;
  fromAddressHash: any | null;
  toAddressHash: any | null;
  status: Status | null;
  gas: string | null;
  gasPrice: string | null;
  index: number | null;
  input: any | null;
  nonce: number | null;
  r: string | null;
  s: string | null;
  v: string | null;
}

export interface QueryAddressByHash_address_transactions_edges {
  node: QueryAddressByHash_address_transactions_edges_node | null;
  cursor: string | null;
}

export interface QueryAddressByHash_address_transactions {
  pageInfo: QueryAddressByHash_address_transactions_pageInfo | null;
  edges: (QueryAddressByHash_address_transactions_edges | null)[] | null;
}

export interface QueryAddressByHash_address {
  fetchedCoinBalance: string | null;
  fetchedCoinBalanceBlockNumber: number | null;
  nonce: number | null;
  hash: any;
  gasUsed: number | null;
  hashQr: string | null;
  numTxs: number | null;
  transactions: QueryAddressByHash_address_transactions | null;
}

export interface QueryAddressByHash {
  address: QueryAddressByHash_address | null;
}

export interface QueryAddressByHashVariables {
  hash: any;
  first?: number | null;
  after?: string | null;
}
