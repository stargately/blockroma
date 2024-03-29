/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL query operation: TokenTransfer
// ====================================================

export interface TokenTransfer_tokenTransfer_edges_node {
  id: string | null;
  transactionHash: any | null;
  logIndex: number | null;
  fromAddress: any | null;
  toAddress: any | null;
  amount: string | null;
  amountWithDecimals: string | null;
  tokenId: string | null;
  tokenContractAddress: any | null;
  block: any | null;
  blockNumber: number | null;
  amounts: (string | null)[] | null;
  tokenIds: (string | null)[] | null;
  createdAt: any | null;
  updatedAt: any | null;
  type: string | null;
}

export interface TokenTransfer_tokenTransfer_edges {
  node: TokenTransfer_tokenTransfer_edges_node | null;
}

export interface TokenTransfer_tokenTransfer {
  edges: (TokenTransfer_tokenTransfer_edges | null)[] | null;
}

export interface TokenTransfer {
  tokenTransfer: TokenTransfer_tokenTransfer | null;
}

export interface TokenTransferVariables {
  transactionHash?: any | null;
}
