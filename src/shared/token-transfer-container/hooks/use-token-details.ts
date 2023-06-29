import { useQuery } from "@apollo/client";
import { queryTokenDetails } from "@/shared/token-transfer-container/data/queries";
import { TokenDetails } from "@/shared/token-transfer-container/data/__generated__/TokenDetails";

export const useTokenDetails = (tokenContractAddressHash?: string) => {
  const { data, loading, error, refetch } = useQuery<TokenDetails>(
    queryTokenDetails,
    {
      ssr: false,
      skip: !tokenContractAddressHash,
      variables: {
        tokenContractAddressHash,
      },
    }
  );
  return {
    tokenDetailsData: data,
    tokenDetailsLoading: loading,
    tokenDetailsError: error,
    tokenDetailsRefetch: refetch,
  };
};
