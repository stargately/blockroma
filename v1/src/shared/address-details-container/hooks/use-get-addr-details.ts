import { useQuery } from "@apollo/client";
import { queryAddressDetailsByHash } from "@/shared/address-details-container/data/queries";
import { QueryAddressDetailsByHash } from "@/shared/address-details-container/data/__generated__/QueryAddressDetailsByHash";

export const useGetAddrDetails = ({ hash }: { hash?: string }) => {
  const { loading, data, error, refetch } = useQuery<QueryAddressDetailsByHash>(
    queryAddressDetailsByHash,
    {
      ssr: false,
      variables: {
        hash,
      },
    }
  );
  return { loading, data, error, refetch };
};
