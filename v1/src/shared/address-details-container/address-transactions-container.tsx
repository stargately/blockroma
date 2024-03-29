import React from "react";
import { AdrTransactionItem } from "@/shared/explorer-components/adr-transaction.item";
import { useGetAddr } from "@/shared/address-details-container/hooks/use-get-addr";
import { t } from "onefx/lib/iso-i18n";
import { assetURL } from "onefx/lib/asset-url";

type Props = {
  addressHash: string;
};

export const AddressTransactionsContainer: React.FC<Props> = ({
  addressHash,
}) => {
  const { data, loading, error, refetch } = useGetAddr({
    hash: addressHash,
    first: 20,
    after: 0,
  });
  if (loading) {
    // TODO(dora):
    return <></>;
  }

  const txs = data?.address?.transactions?.edges?.map((e) => e?.node);

  return (
    <section data-page="address-transactions" id="txs">
      <div className="card">
        <div className="card-tabs js-card-tabs">
          <a
            className="card-tab active"
            href={assetURL(`address/${addressHash}/transactions`)}
          >
            Transactions
          </a>
        </div>
        <div
          className="card-body"
          data-async-listing={assetURL(`address/${addressHash}/transactions`)}
        >
          <div
            data-selector="channel-disconnected-message"
            style={{ display: "none" }}
          >
            <div
              data-selector="reload-button"
              className="alert alert-danger"
              style={{ padding: "0.75rem 0rem", cursor: "pointer" }}
            >
              <span className="alert alert-danger">
                Connection Lost, click to load newer transactions
              </span>
            </div>
          </div>
          <div className="clearfix">
            <h2 className="card-title float-left">{t("nav.txs")}</h2>
          </div>
          <div
            data-selector="channel-batching-message"
            style={{ display: "none" }}
          >
            <div
              data-selector="reload-transactions-button"
              className="alert alert-info"
            >
              <a href="#" className="alert-link">
                <span data-selector="channel-batching-count" /> More
                transactions have come in
              </a>
            </div>
          </div>

          {error && error.graphQLErrors[0]?.extensions?.code !== "NOT_FOUND" && (
            <button
              data-error-message
              className="alert alert-danger col-12 text-left"
              onClick={refetch}
            >
              <span className="alert-link">{t("info.err")}</span>
            </button>
          )}

          {(!txs || !txs.length) && (
            <div data-empty-response-message>
              <div
                className="tile tile-muted text-center"
                data-selector="empty-transactions-list"
              >
                There are no transactions for this address.
              </div>
            </div>
          )}

          {!!txs && (
            <div data-items>
              {txs.map((tx) => (
                <AdrTransactionItem
                  key={tx?.hash}
                  tx={tx}
                  selfAddressHash={addressHash}
                />
              ))}
            </div>
          )}

          {/*
          TODO(dora): download csv

                    <div className="transaction-bottom-panel">
            <div csv-download className="download-all-transactions">
              Download{" "}
              <a
                className="download-all-transactions-link"
                href="/csv-export?address=0x11F3fb5677c84131377BD9762Ee2ef451eEF47DB&type=transactions"
              >
                CSV
                <svg xmlns="http://www.w3.org/2000/svg" width={14} height={16}>
                  <path
                    fill="#333333"
                    fillRule="evenodd"
                    d="M13 16H1c-.999 0-1-1-1-1V1s-.004-1 1-1h6l7 7v8s-.032 1-1 1zm-1-8c0-.99-1-1-1-1H8s-1 .001-1-1V3c0-.999-1-1-1-1H2v12h10V8z"
                  />
                </svg>
              </a>
            </div>
          </div>

          */}
        </div>
      </div>
    </section>
  );
};
