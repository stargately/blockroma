import { useQuery } from "@apollo/client";
import { queryTokens } from "@/shared/token-container/data/query";
import { Tokens } from "@/shared/token-container/data/__generated__/Tokens";

const pageSize = 50;

export const useQueryTokens = (querySymbol: string) => {
  const { data, loading, fetchMore } = useQuery<Tokens>(queryTokens, {
    ssr: false,
    variables: {
      after: "0",
      first: pageSize,
      symbol: querySymbol,
    },
  });
  return {
    tokensFetchMore: async (cursor?: string | null, symbol?: string) => {
      if (!cursor) {
        return;
      }
      await fetchMore({
        variables: {
          after: cursor,
          first: pageSize,
          symbol,
        },
      });
    },
    tokensData: data,
    tokensLoading: loading,
  };
};
