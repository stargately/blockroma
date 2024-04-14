import React, { useState } from "react";
import { AdrTransactionItem } from "@/shared/explorer-components/adr-transaction.item";
import { useGetAddr } from "@/shared/address-details-container/hooks/use-get-addr";
import { assetURL } from "@/shared/common/asset-url";
import { useTranslation } from "next-i18next";
import {Pagination} from "@/shared/explorer-components/pagination";
import {useRouter} from "next/router";
import {paginationProcessTotalNumPage} from "@/shared/common/functions/paginations";

type Props = {
  addressHash: string;
};

export const AddressTransactionsContainer: React.FC<Props> = ({
  addressHash,
}) => {
  const { t } = useTranslation("common");
  const router = useRouter();
  const query = router.query;
  const initialPage = Number(query.page) || 1;
  const [currentPage, setCurrentPage] = useState(initialPage); // Manage current page, starts from 0


  // Pagination variables
  const pageSize = 10;

  const { data, loading, error, refetch } = useGetAddr({
    hash: addressHash,
    first: pageSize,
    after: (currentPage - 1) * pageSize,
  });
  const setCurPageWithSideEffect = async (page: number) => {
    setCurrentPage(page);
    await router.push({ pathname: router.pathname, query: { ...query, page: page.toString() } }, undefined, { shallow: true });
  };

  // Logic to determine the number of pages
  const totalPages = paginationProcessTotalNumPage(data?.address?.transactions);

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
          {/* Render transaction items or messages accordingly */}
          {error && error.graphQLErrors[0]?.extensions?.code !== "NOT_FOUND" && (
            <button
              className="alert alert-danger col-12 text-left"
              onClick={() => refetch()}
            >
              <span className="alert-link">{t("info.err")}</span>
            </button>
          )}
          {!txs?.length ? (
            <div
              className="tile tile-muted text-center"
              data-selector="empty-transactions-list"
            >
              There are no transactions for this address.
            </div>
          ) : (
            <>
              <Pagination setCurPage={setCurPageWithSideEffect} curPage={currentPage} numPages={totalPages} position="bottom" />
              <div className={"mt-4"}>
              {txs.map((tx) => (
                <AdrTransactionItem
                  key={tx?.hash}
                  tx={tx}
                  selfAddressHash={addressHash}
                />
              ))}
            </div>
              <Pagination setCurPage={setCurPageWithSideEffect} curPage={currentPage} numPages={totalPages} position="bottom" />

            </>
          )}


        </div>
      </div>
    </section>
  );
};
