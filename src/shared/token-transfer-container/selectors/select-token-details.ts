import { TokenDetails } from "@/shared/token-transfer-container/data/__generated__/TokenDetails";

export const selectTokenDetails = (resp: TokenDetails | null | undefined) => {
  return resp?.token;
};
