import { getTxs } from "@/shared/block-details-container/data/quries";
import { useQuery } from "@apollo/client";
import { GetTxs } from "@/shared/block-details-container/data/__generated__/GetTxs";

export const useGetTxs = (
  filters: {
    blockNumber?: number;
    first: number;
    after: number;
  },
  { pollInterval }: { pollInterval?: number }
) => {
  const { loading, data, error, refetch } = useQuery<GetTxs>(getTxs, {
    ssr: false,
    variables: {
      ...filters,
      after: String(filters.after),
    },
    pollInterval,
    ...(pollInterval
      ? {
          fetchPolicy: "no-cache",
        }
      : {}),
  });
  return { loading, data, error, refetch };
};
