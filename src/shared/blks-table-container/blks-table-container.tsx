import React, { useState } from "react";
import { useQueryBlocks } from "@/shared/blks-table-container/hooks/use-query-blocks";
import { BlkList } from "@/shared/explorer-components/blk-list";
import { Pagination } from "@/shared/explorer-components/pagination";
import { useHistory } from "react-router-dom";
import { useLocation } from "onefx/lib/react-router";
import { t } from "onefx/lib/iso-i18n";
import { paginationProcessTotalNumPage } from "@/shared/common/functions/paginations";

export const BlksTableContainer: React.FC = () => {
  return (
    <main className="js-ad-dependant-pt pt-5">
      <p className="alert alert-info" role="alert" />
      <p className="alert alert-danger" role="alert" />
      <section className="container" data-page="block-list">
        <div className="ad mb-3" style={{ display: "none" }}>
          <span className="ad-prefix" />:{" "}
          <img className="ad-img-url" width={20} height={20} />{" "}
          <b>
            <span className="ad-name" />
          </b>{" "}
          - <span className="ad-short-description" />{" "}
          <a className="ad-url">
            <b>
              <span className="ad-cta-button" />
            </b>
          </a>
        </div>
        <div className="card">
          <div className="card-body" data-async-listing="/poa/core/blocks">
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
                  Connection Lost, click to load newer blocks
                </span>
              </div>
            </div>
            <h1 className="card-title list-title-description">
              {t("nav.blocks")}
            </h1>

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
  const { loading, data, error, refetch } = useQueryBlocks(
    { first: pageSize, after: (curPage - 1) * pageSize },
    {}
  );
  if (loading) {
    // TODO(dora)
    return <></>;
  }

  const blks = data?.blocks?.edges?.map((e) => e?.node);

  const numPage = paginationProcessTotalNumPage(data?.blocks);

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

      {blks && blks?.length > 0 ? (
        <>
          <div className="list-top-pagination-container-wrapper">
            <Pagination
              setCurPage={setCurPageWithSideEffect}
              curPage={curPage}
              numPages={numPage}
              position="top"
            />
          </div>
          <BlkList blks={blks} />
          <div className="list-bottom-pagination-container-wrapper">
            <Pagination
              position="bottom"
              setCurPage={setCurPageWithSideEffect}
              curPage={curPage}
              numPages={numPage}
            />
          </div>
        </>
      ) : (
        <>
          <div data-empty-response-message>
            <span>There are no blocks.</span>
          </div>
        </>
      )}
    </>
  );
};
