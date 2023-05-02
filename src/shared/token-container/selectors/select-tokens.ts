import { Tokens } from "@/shared/token-container/data/__generated__/Tokens";

export const selectTokens = (tks?: Tokens) => {
  return tks?.tokens?.edges?.map((ed) => ed?.node);
};
