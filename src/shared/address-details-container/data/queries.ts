import { gql } from "@apollo/client";

export const queryAddressByHash = gql`
  query QueryAddressByHash($hash: Buffer!, $first: Float, $after: String) {
    address(hash: $hash) {
      transactions(first: $first, after: $after) {
        pageInfo {
          hasNextPage
          startCursor

          hasPreviousPage
          endCursor
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
          cursor
        }
      }
    }
  }
`;

export const queryAddressDetailsByHash = gql`
  query QueryAddressDetailsByHash($hash: Buffer!) {
    address(hash: $hash) {
      fetchedCoinBalance
      fetchedCoinBalanceBlockNumber
      nonce
      hash
      gasUsed
      hashQr
      numTxs
    }
  }
`;
