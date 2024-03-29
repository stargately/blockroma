/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL query operation: TokenDetails
// ====================================================

export interface TokenDetails_token {
  name: string | null;
  symbol: string | null;
  totalSupply: string | null;
  decimals: string | null;
  type: string | null;
  contractAddress: any | null;
  skipMetadata: boolean | null;
}

export interface TokenDetails {
  token: TokenDetails_token | null;
}

export interface TokenDetailsVariables {
  tokenContractAddressHash: any;
}
