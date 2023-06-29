import { TokenTransfer } from "@/shared/token-transfer-container/data/__generated__/TokenTransfer";

export const selectTokenTransfer = (resp: TokenTransfer | null | undefined) => {
  return resp?.tokenTransfer?.edges
    ?.map((ed) => ed?.node)
    ?.map((node) => ({
      transactionHash: node?.transactionHash,
      fromAddress: node?.fromAddress,
      toAddress: node?.toAddress,
      amount: node?.amount,
      tokenContractAddress: node?.tokenContractAddress,
    }));
};
