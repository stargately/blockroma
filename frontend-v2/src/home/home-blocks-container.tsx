import React from "react";
import { useQueryBlocks } from "@/shared/blks-table-container/hooks/use-query-blocks";
import { BlkTile } from "@/shared/home/components/blk-tile";
import { assetURL } from "@/shared/common/asset-url";
import { useTranslation } from "next-i18next";
import Link from "next/link";

export const HomeBlocksContainer: React.FC = () => {
  const { t } = useTranslation("common");
  const { data, loading, error, refetch } = useQueryBlocks(
    { first: 4, after: 0 },
    { pollInterval: 5000 },
  );

  if (loading) {
    // TODO(dora)
    return <></>;
  }

  const blks = data?.blocks?.edges?.map((e) => e?.node);

  return (
    <div className="card card-chain-blocks js-ad-dependant-mb-3 mb-3">
      <div className="card-body">
        <Link className="btn-line float-right" href={assetURL("blocks")}>
          View All Blocks
        </Link>
        <h2 className="card-title">{t("nav.blocks")}</h2>

        {error && (
          <button
            data-error-message
            className="alert alert-danger col-12 text-left"
            onClick={() => refetch()}
          >
            <span className="alert-link">{t("info.err")}</span>
          </button>
        )}

        <div
          className="row"
          data-selector="chain-block-list"
          data-url="/chain-blocks"
        >
          {!!blks?.length &&
            blks.map((blk) => <BlkTile key={blk?.hash} blk={blk} />)}
        </div>
      </div>
    </div>
  );
};
