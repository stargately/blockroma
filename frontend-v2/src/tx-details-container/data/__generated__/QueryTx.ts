/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.


// ====================================================
// GraphQL query operation: QueryTx
// ====================================================

import {Status} from "@/shared/__generated__/globalTypes";

export interface QueryTx_transaction {
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
  revertReason: string | null;
  maxPriorityFeePerGas: string | null;
  maxFeePerGas: string | null;
}

export interface QueryTx {
  transaction: QueryTx_transaction | null;
}

export interface QueryTxVariables {
  hash: any;
}
