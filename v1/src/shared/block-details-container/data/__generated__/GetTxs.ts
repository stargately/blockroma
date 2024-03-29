/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

import { Status } from "./../../../../../__generated__/globalTypes";

// ====================================================
// GraphQL query operation: GetTxs
// ====================================================

export interface GetTxs_transactions_pageInfo {
  hasNextPage: boolean | null;
  endCursor: string | null;
  startCursor: string | null;
  hasPreviousPage: boolean | null;
}

export interface GetTxs_transactions_edges_node {
  id: string;
  timestamp: any | null;
  hash: any | null;
  blockNumber: number | null;
  value: string | null;
  valueWithDecimal: string | null;
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

export interface GetTxs_transactions_edges {
  node: GetTxs_transactions_edges_node | null;
}

export interface GetTxs_transactions {
  pageInfo: GetTxs_transactions_pageInfo | null;
  edges: (GetTxs_transactions_edges | null)[] | null;
}

export interface GetTxs {
  transactions: GetTxs_transactions | null;
}

export interface GetTxsVariables {
  blockNumber?: number | null;
  first?: number | null;
  after?: string | null;
}
