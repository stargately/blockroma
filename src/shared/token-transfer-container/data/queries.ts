import { gql } from "@apollo/client";

export const queryTokenTransfer = gql`
  query TokenTransfer($transactionHash: Buffer) {
    tokenTransfer(transactionHash: $transactionHash) {
      edges {
        node {
          id
          transactionHash
          logIndex
          fromAddress
          toAddress
          amount
          tokenId
          tokenContractAddress
          block
          blockNumber
          amounts
          tokenIds
          createdAt
          updatedAt
          type
        }
      }
    }
  }
`;

export const queryTokenDetails = gql`
  query TokenDetails($tokenContractAddressHash: Buffer!) {
    token(tokenContractAddressHash: $tokenContractAddressHash) {
      name
      symbol
      totalSupply
      decimals
      type
      contractAddress
      skipMetadata
    }
  }
`;
