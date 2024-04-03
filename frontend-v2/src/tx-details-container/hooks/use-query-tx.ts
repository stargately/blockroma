import { useQuery } from "@apollo/client";
import { queryTx } from "@/shared/tx-details-container/data/query-tx";
import { QueryTx } from "../data/__generated__/QueryTx";

function isEthHash(hash: string): boolean {
  const ethHashRegex = /^0x[a-fA-F0-9]{64}$/;
  return ethHashRegex.test(hash);
}

export const useQueryTx = ({ hash }: { hash: string }) => {
  const { data, loading, error, refetch } = useQuery<QueryTx>(queryTx, {
    ssr: false,
    variables: {
      hash,
    },
    pollInterval: isEthHash(hash) ? 5000 : undefined,
  });
  return {
    data,
    loading,
    error,
    refetch: () => refetch(),
  };
};
