import { queryBlocks } from "@/shared/blks-table-container/data/queries";
import { useQuery } from "@apollo/client";
import { QueryBlocks } from "@/shared/blks-table-container/data/__generated__/QueryBlocks";

export const useQueryBlocks = (
  {
    first,
    after,
  }: {
    first: number;
    after: number;
  },
  { pollInterval = undefined }: { pollInterval?: number }
) => {
  const { loading, data, error, refetch } = useQuery<QueryBlocks>(queryBlocks, {
    ssr: false,
    variables: {
      first,
      after: String(after),
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
