/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL query operation: QueryAddressDetailsByHash
// ====================================================

export interface QueryAddressDetailsByHash_address {
  fetchedCoinBalance: string | null;
  fetchedCoinBalanceWithDecimal: string | null;
  fetchedCoinBalanceBlockNumber: number | null;
  nonce: number | null;
  hash: any;
  gasUsed: number | null;
  hashQr: string | null;
  numTxs: number | null;
}

export interface QueryAddressDetailsByHash {
  address: QueryAddressDetailsByHash_address | null;
}

export interface QueryAddressDetailsByHashVariables {
  hash: any;
}
