import { gql } from "@apollo/client";

export const queryTokens = gql`
  query Tokens(
    $after: String
    $first: Float
    $before: String
    $last: Float
    $symbol: String
  ) {
    tokens(
      after: $after
      first: $first
      before: $before
      last: $last
      symbol: $symbol
    ) {
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
        startCursor
        endCursor
        hasNextPage
        hasPreviousPage
      }
    }
  }
`;
