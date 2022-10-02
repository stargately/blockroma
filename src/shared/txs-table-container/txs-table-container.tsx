import React, { useState } from "react";
import { useGetTxs } from "@/shared/block-details-container/hooks/use-get-txs";
import { TxsList } from "@/shared/explorer-components/txs-list";
import { useLocation } from "onefx/lib/react-router";
import { useHistory } from "react-router-dom";
import { Pagination } from "@/shared/explorer-components/pagination";
import { t } from "onefx/lib/iso-i18n";
import { paginationProcessTotalNumPage } from "@/shared/common/functions/paginations";

export const TxsTableContainer: React.FC = () => {
  return (
    <main className="js-ad-dependant-pt pt-5">
      <p className="alert alert-info" role="alert" />
      <p className="alert alert-danger" role="alert" />
      <section className="container" data-page="transaction-list">
        <div className="card">
          <div className="card-body">
            <h1 className="card-title list-title-description">
              Validated Transactions
            </h1>
            <div
              data-selector="channel-batching-message"
              className="d-none"
              style={{ display: "none" }}
            >
              <div data-selector="reload-button" className="alert alert-info">
                <a href="#" className="alert-link">
                  <span data-selector="channel-batching-count" /> More
                  transactions have come in
                </a>
              </div>
            </div>
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

            <TableWithPagination />
          </div>
        </div>
      </section>
    </main>
  );
};

const TableWithPagination = () => {
  const location = useLocation();
  const history = useHistory();
  const search = new URLSearchParams(location.search);
  const initialPage = Number(search.get("page")) || 1;
  const [curPage, setCurPage] = useState(initialPage);
  const pageSize = 20;

  const setCurPageWithSideEffect = (p: number) => {
    setCurPage(p);
    history.push({ search: `?page=${p}` });
  };
  const { data, loading, error, refetch } = useGetTxs(
    { first: pageSize, after: (curPage - 1) * pageSize },
    {}
  );
  if (loading) {
    // TODO(dora)
    return <></>;
  }

  const txs = data?.transactions?.edges?.map((e) => e?.node);

  const numPage = paginationProcessTotalNumPage(data?.transactions);

  return (
    <>
      {error && (
        <button
          onClick={refetch}
          data-error-message
          className="alert alert-danger col-12 text-left"
        >
          <span className="alert-link">{t("info.err")}</span>
        </button>
      )}

      {!txs?.length && (
        <div data-empty-response-message style={{ display: "none" }}>
          <div className="tile tile-muted text-center">
            <span data-selector="empty-internal-transactions-list">
              There are no transactions.
            </span>
          </div>
        </div>
      )}

      {txs && txs?.length > 0 ? (
        <>
          <Pagination
            setCurPage={setCurPageWithSideEffect}
            curPage={curPage}
            numPages={numPage}
            position="top"
          />

          <div data-selector="transactions-list">
            <TxsList txs={txs} />
          </div>

          <Pagination
            setCurPage={setCurPageWithSideEffect}
            position="bottom"
            curPage={curPage}
            numPages={numPage}
          />
        </>
      ) : (
        <>
          <div data-empty-response-message style={{ display: "none" }}>
            <div className="tile tile-muted text-center">
              <span data-selector="empty-internal-transactions-list">
                There are no transactions.
              </span>
            </div>
          </div>
        </>
      )}
    </>
  );
};
