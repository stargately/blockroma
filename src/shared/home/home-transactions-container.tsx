import * as React from "react";
import { useGetTxs } from "@/shared/block-details-container/hooks/use-get-txs";
import { TxTransactionItem } from "@/shared/explorer-components/tx-transaction-item";

//         TODO(dora): more tx coming in
// eslint-disable-next-line @typescript-eslint/no-unused-vars
// @ts-ignore
function MoreTxComingIn(): JSX.Element {
  return (
    <div data-selector="channel-batching-message" style={{}}>
      <div
        data-selector="reload-transactions-button"
        className="alert alert-info"
      >
        <a href="#" className="alert-link">
          <span data-selector="channel-batching-count">516</span> More
          transactions have come in
        </a>
      </div>
    </div>
  );
}

export function HomeTransactionsContainer(): JSX.Element {
  const { data, refetch, loading, error } = useGetTxs(
    { first: 20, after: 0 },
    { pollInterval: 5000 }
  );

  if (loading) {
    // TODO(dora)
    return <></>;
  }

  const txs = data?.transactions?.edges?.map((e) => e?.node);

  return (
    <div className="card card-chain-transactions">
      <div className="card-body">
        <a className="btn-line float-right" href="/txs">
          View All Transactions
        </a>
        <h2 className="card-title lg-card-title">Transactions</h2>
        {/*
        TODO(dora): more tx coming in
        <MoreTxComingIn/>
        */}

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
              There are no transactions for this chain.
            </div>
          </div>
        )}

        <span
          data-selector="transactions-list"
          data-transactions-path="/recent-transactions"
        >
          {!!txs?.length &&
            txs.map((tx) => <TxTransactionItem key={tx?.hash} tx={tx} />)}
        </span>
      </div>
    </div>
  );
}
