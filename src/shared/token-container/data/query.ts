import { gql } from "@apollo/client";

export const queryTokens = gql`
  query Tokens($after: String, $first: Float) {
    tokens(after: $after, first: $first) {
      edges {
        cursor
        node {
          name
          symbol
          totalSupply
          decimals
          type
          contractAddress
          skipMetadata
        }
      }
      pageInfo {
        endCursor
        hasNextPage
      }
    }
  }
`;
