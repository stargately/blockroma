/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

import { Status } from "./../../../../../__generated__/globalTypes";

// ====================================================
// GraphQL query operation: QueryTx
// ====================================================

export interface QueryTx_transaction {
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
  gasPriceWithDecimal: string | null;
  index: number | null;
  input: any | null;
  nonce: number | null;
  r: string | null;
  s: string | null;
  v: string | null;
  revertReason: string | null;
  maxPriorityFeePerGas: string | null;
  maxPriorityFeePerGasWithDecimal: string | null;
  maxFeePerGas: string | null;
  maxFeePerGasWithDecimal: string | null;
}

export interface QueryTx {
  transaction: QueryTx_transaction | null;
}

export interface QueryTxVariables {
  hash: any;
}
