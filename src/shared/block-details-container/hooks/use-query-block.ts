import { useQuery } from "@apollo/client";
import { queryBlock } from "@/shared/block-details-container/data/quries";
import { QueryBlock } from "../data/__generated__/QueryBlock";

export const useQueryBlock = (blockNumber: number) => {
  const { loading, data, error, refetch } = useQuery<QueryBlock>(queryBlock, {
    ssr: false,
    variables: { blockNumber },
  });
  return { loading, data, error, refetch };
};
