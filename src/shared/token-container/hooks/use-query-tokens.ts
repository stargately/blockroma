import { useQuery } from "@apollo/client";
import { queryTokens } from "@/shared/token-container/data/query";
import { Tokens } from "@/shared/token-container/data/__generated__/Tokens";

export const useQueryTokens = () => {
  const { data, loading } = useQuery<Tokens>(queryTokens, {
    ssr: false,
    variables: {
      after: "0",
      first: 50,
    },
  });
  return {
    tokensData: data,
    tokensloading: loading,
  };
};
