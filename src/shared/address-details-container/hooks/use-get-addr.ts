import { useQuery } from "@apollo/client";
import { queryAddressByHash } from "@/shared/address-details-container/data/queries";
import { QueryAddressByHash } from "@/shared/address-details-container/data/__generated__/QueryAddressByHash";

export const useGetAddr = ({
  hash,
  first,
  after,
}: {
  hash?: string;
  first: number;
  after: number;
}) => {
  const { loading, data, error, refetch } = useQuery<QueryAddressByHash>(
    queryAddressByHash,
    {
      ssr: false,
      variables: {
        hash,
        first,
        after: String(after),
      },
    }
  );
  return { loading, data, error, refetch };
};
