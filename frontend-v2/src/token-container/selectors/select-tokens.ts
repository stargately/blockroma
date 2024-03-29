import { Tokens } from "@/shared/token-container/data/__generated__/Tokens";

export const selectTokens = (tks?: Tokens) => {
  return {
    tokens: tks?.tokens?.edges?.map((ed) => ed?.node),
    currentCursor: tks?.tokens?.edges?.at(0)?.cursor,
    hasNextPage: tks?.tokens?.pageInfo?.hasNextPage,
  };
};
