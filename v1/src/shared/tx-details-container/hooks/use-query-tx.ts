import { useQuery } from "@apollo/client";
import { queryTx } from "@/shared/tx-details-container/data/query-tx";
import { QueryTx } from "../data/__generated__/QueryTx";

export const useQueryTx = ({ hash }: { hash: string }) => {
  const { data, loading, error, refetch } = useQuery<QueryTx>(queryTx, {
    ssr: false,
    variables: {
      hash,
    },
  });
  return {
    data,
    loading,
    error,
    refetch,
  };
};
