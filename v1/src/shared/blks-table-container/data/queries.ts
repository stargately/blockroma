import { gql } from "@apollo/client";

export const queryBlocks = gql`
  query QueryBlocks($first: Float, $after: String) {
    blocks(first: $first, after: $after) {
      edges {
        node {
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
      pageInfo {
        hasNextPage
        endCursor
        startCursor
        hasPreviousPage
      }
    }
  }
`;
