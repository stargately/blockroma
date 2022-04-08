import React from "react";
import { useGetTxs } from "@/shared/block-details-container/hooks/use-get-txs";
import { TxsList } from "@/shared/explorer-components/txs-list";

type Props = {
  blockNumber: number;
};

export const BlockTransactionsListContainer: React.FC<Props> = ({
  blockNumber,
}) => {
  const { loading, error, data, refetch } = useGetTxs(
    {
      blockNumber,
      first: 1000,
      after: 0,
    },
    {}
  );
  if (loading) {
    // TODO(dora)
    return <></>;
  }

  const txs = data?.transactions?.edges?.map((e) => e?.node);

  return (
    <div
      className="card-body"
      data-async-load
      data-async-listing="/block/26559654/transactions"
      id="txs"
    >
      {error && (
        <button
          data-error-message
          className="alert alert-danger col-12 text-left"
          onClick={refetch}
        >
          <span className="alert-link">
            Something went wrong, click to reload.
          </span>
        </button>
      )}

      {!txs?.length && (
        <div data-empty-response-message>
          <div
            className="tile tile-muted text-center"
            data-selector="empty-transactions-list"
          >
            There are no transactions for this block.
          </div>
        </div>
      )}

      {!!txs?.length && <TxsList txs={txs} />}
    </div>
  );
};
