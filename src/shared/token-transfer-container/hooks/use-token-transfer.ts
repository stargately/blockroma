import { queryTokenTransfer } from "@/shared/token-transfer-container/data/queries";
import { TokenTransfer } from "@/shared/token-transfer-container/data/__generated__/TokenTransfer";
import { useQuery } from "@apollo/client";

export const useTokenTransfer = (transactionHash?: string) => {
  const { data, loading, error, refetch } = useQuery<TokenTransfer>(
    queryTokenTransfer,
    {
      ssr: false,
      skip: !transactionHash,
      variables: { transactionHash },
    }
  );
  return {
    transferData: data,
    transferLoading: loading,
    transferError: error,
    transferRefetch: refetch,
  };
};
