import { Link } from "react-router-dom";
import * as React from "react";
import { BlockTransactionsListContainer } from "@/shared/block-details-container/block-transactions-list-container";
import { t } from "onefx/lib/iso-i18n";

export function BlockTransactions({
  blockNumber,
}: {
  blockNumber: number;
}): JSX.Element {
  return (
    <section>
      <div className="card mb-3">
        <div className="card-tabs js-card-tabs">
          <Link
            className="card-tab active noCaret"
            to={`/block/${blockNumber}/transactions`}
          >
            {t("nav.txs")}
          </Link>
        </div>

        <BlockTransactionsListContainer blockNumber={blockNumber} />
      </div>
    </section>
  );
}
