import { gql } from "@apollo/client";

export const queryBlock = gql`
  query QueryBlock($blockNumber: Int!) {
    block(number: $blockNumber) {
      consensus
      difficulty
      gasLimit
      gasUsed
      hash
      miner
      nonce
      number
      parentHash
      size
      timestamp
      totalDifficulty
      numTxs
    }
  }
`;

export const getTxs = gql`
  query GetTxs($blockNumber: Int, $first: Float, $after: String) {
    transactions(blockNumber: $blockNumber, first: $first, after: $after) {
      pageInfo {
        hasNextPage
        endCursor
        startCursor
        hasPreviousPage
      }
      edges {
        node {
          id
          timestamp
          hash
          blockNumber
          value
          gasUsed
          cumulativeGasUsed
          error
          fromAddressHash
          toAddressHash
          status
          gas
          gasPrice
          index
          input
          nonce
          r
          s
          v
        }
      }
    }
  }
`;
