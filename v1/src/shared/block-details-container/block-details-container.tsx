import * as React from "react";
import { useQueryBlock } from "@/shared/block-details-container/hooks/use-query-block";
import { useParams } from "react-router";
import { Link } from "react-router-dom";
import { getGasUsedPercent } from "@/shared/common/get-gas-used-percent";
import { CopyToClipboard } from "@/shared/explorer-components/copy-to-clipboard";
import { TickingTs } from "@/shared/explorer-components/ticking-ts";
import format from "date-fns/format";
import { t } from "onefx/lib/iso-i18n";
import { assetURL } from "onefx/lib/asset-url";
import { BlockTransactions } from "./block-transactions";

export function BlockDetailsContainer(): JSX.Element {
  const params = useParams<{ blockNumber: string }>();
  const blockNumber = parseInt(params.blockNumber, 10);
  const { data, loading, error } = useQueryBlock(blockNumber);
  if (loading) {
    // TODO(dora)
    return <></>;
  }
  if (error) {
    // TODO(dora)
    return <></>;
  }
  if (!data || !data.block) {
    // TODO(dora)
    return <></>;
  }
  const {
    difficulty,
    gasLimit,
    gasUsed,
    hash,
    miner,
    nonce,
    parentHash,
    size,
    timestamp,
    totalDifficulty,
    numTxs,
  } = data.block;
  const gasUsedPercent = getGasUsedPercent(gasUsed, gasLimit);
  return (
    <main className="js-ad-dependant-pt pt-5">
      <p className="alert alert-info" role="alert" />
      <p className="alert alert-danger" role="alert" />
      <section className="container">
        <section>
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
          <div className="row">
            <div className="col-md-12 js-ad-dependant-mb-2 mb-2">
              <div className="card js-ad-dependant-mb-2 mb-2">
                <div className="card-body fs-14" style={{ lineHeight: "32px" }}>
                  <dl className="pagination-container">
                    <h1 className="card-title" data-test="detail_type">
                      {t("bk.block_details")}
                    </h1>
                    <ul className="pagination">
                      <li className="page-item">
                        <Link
                          className="page-link"
                          to={`/block/${blockNumber - 1}`}
                          data-prev-page-button
                          data-placement="top"
                          data-toggle="tooltip"
                          title="View previous block"
                        >
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            width={6}
                            height={10}
                          >
                            <path
                              fillRule="evenodd"
                              d="M2.358 5l3.357 3.358a.959.959 0 1 1-1.357 1.357L.502 5.859c-.076-.042-.153-.08-.217-.144A.949.949 0 0 1 .011 5a.949.949 0 0 1 .274-.715c.064-.064.142-.102.217-.145L4.358.285a.959.959 0 1 1 1.357 1.357L2.358 5z"
                            />
                          </svg>
                        </Link>
                      </li>
                      <li className="page-item">
                        <a
                          className="page-link no-hover"
                          href="#"
                          data-page-number
                        >
                          {t("bk.block")} {blockNumber}
                        </a>
                      </li>
                      <li className="page-item">
                        <Link
                          className="page-link"
                          to={`/block/${blockNumber + 1}`}
                          data-next-page-button
                          data-placement="top"
                          data-toggle="tooltip"
                          title="View next block"
                        >
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            width={6}
                            height={10}
                          >
                            <path
                              fillRule="evenodd"
                              d="M5.715 5.715c-.064.064-.141.102-.217.144L1.642 9.715A.959.959 0 1 1 .285 8.358L3.642 5 .285 1.642A.959.959 0 1 1 1.642.285L5.498 4.14c.075.043.153.081.217.145A.949.949 0 0 1 5.989 5a.949.949 0 0 1-.274.715z"
                            />
                          </svg>
                        </Link>
                      </li>
                    </ul>
                  </dl>
                  <dl className="row">
                    <dt className="col-sm-3 col-lg-2 text-muted">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title={t("bk.block_height.tip")}
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      {t("bk.block_height")}
                    </dt>
                    <dd
                      className="col-sm-9 col-lg-10"
                      data-test="block_detail_number"
                    >
                      {blockNumber}
                    </dd>
                  </dl>
                  <dl className="row">
                    <dt className="col-sm-3 col-lg-2 text-muted">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title={t("bk.timestamp.tip")}
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      {t("bk.timestamp")}
                    </dt>
                    <dd
                      className="col-sm-9 col-lg-10"
                      data-from-now={timestamp}
                    >
                      <TickingTs timestamp={timestamp} /> |{" "}
                      {format(new Date(timestamp), "MMM-d-y hh:mm:ss a x")} UTC
                    </dd>
                  </dl>
                  <dl className="row">
                    <dt className="col-sm-3 col-lg-2 text-muted">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title={t("bk.txs.tip")}
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      {t("bk.txs")}
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      <a
                        href="#txs"
                        className="page-link bs-label large btn-no-border-link-to-tems"
                      >
                        {numTxs} {t("nav.txs")}
                      </a>
                    </dd>
                  </dl>
                  <dl className="row">
                    <dt className="col-sm-3 col-lg-2 text-muted">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title={t("bk.validator.tip")}
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      {t("bk.validator")}{" "}
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      <span>
                        <a
                          data-test="address_hash_link"
                          href={assetURL(`address/${miner}`)}
                        >
                          {miner}
                        </a>
                      </span>
                      <CopyToClipboard
                        value={miner}
                        reason={t("bk.copy_address")}
                      />
                    </dd>
                  </dl>
                  <dl className="row">
                    <dt className="col-sm-3 col-lg-2 text-muted">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title={t("bk.size.tip")}
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      {t("bk.size")}
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      {Number(size).toLocaleString()} bytes
                    </dd>
                  </dl>
                  <dl className="row">
                    <dt className="col-sm-3 col-lg-2 text-muted">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title={t("bk.hash.tip")}
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      {t("bk.hash")}
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      {hash}

                      <CopyToClipboard
                        value={hash}
                        reason={t("bk.copy_hash")}
                      />
                    </dd>
                  </dl>

                  <dl className="row">
                    <dt className="col-sm-3 col-lg-2 text-muted">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title={t("bk.parent_hash.tip")}
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      {t("bk.parent_hash")}
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      {blockNumber > 0 ? (
                        <a
                          className="transaction__link"
                          href={assetURL(`block/${blockNumber - 1}`)}
                        >
                          {parentHash}
                        </a>
                      ) : (
                        parentHash
                      )}

                      <CopyToClipboard
                        value={parentHash}
                        reason={t("bk.copy_parent_hash")}
                      />
                    </dd>
                  </dl>

                  <dl className="row">
                    <dt className="col-sm-3 col-lg-2 text-muted">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title={t("bk.difficulty.tip")}
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      {t("bk.difficulty")}
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      {Number(difficulty).toLocaleString()}
                    </dd>
                  </dl>
                  <dl className="row">
                    <dt className="col-sm-3 col-lg-2 text-muted">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title={t("bk.total_difficulty.tip")}
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      {t("bk.total_difficulty")}
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      {Number(totalDifficulty).toLocaleString()}
                    </dd>
                  </dl>
                  <dl className="row">
                    <dt className="col-sm-3 col-lg-2 text-muted">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title={t("bk.gas_used.tip")}
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      {t("bk.gas_used")}
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      {Number(gasUsed).toLocaleString()} | {gasUsedPercent}%
                    </dd>
                  </dl>
                  <dl className="row">
                    <dt className="col-sm-3 col-lg-2 text-muted">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title={t("bk.gas_limit.tip")}
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      {t("bk.gas_limit")}
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      {Number(gasLimit).toLocaleString()}
                    </dd>
                  </dl>
                  <dl className="row">
                    <dt className="col-sm-3 col-lg-2 text-muted">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title="64-bit hash of value verifying proof-of-work (note: null for BMO chains)."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      {t("bk.once")}
                    </dt>
                    <dd className="col-sm-9 col-lg-10">{nonce}</dd>
                  </dl>
                  {/*
                  TODO(dora): baseFeePerGas

                  <dl className="row">
                    <dt className="col-sm-3 col-lg-2 text-muted">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title="Minimum fee required per unit of gas. Fee adjusts based on network congestion."
                      >
                        <i className="fa-solid fa-info-circle"/>
                      </span>
                      Base Fee per Gas
                    </dt>
                    <dd className="col-sm-9 col-lg-10">{0} Gwei</dd>
                  </dl>

                  TODO(dora): Burnt Fees
<dl className="row">
                    <dt className="col-sm-3 col-lg-2 text-muted">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title="BMO burned from transactions included in the block (Base fee (per unit of gas) * Gas Used)."
                      >
                        <i className="fa-solid fa-info-circle"/>
                      </span>
                      Burnt Fees
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      <i className="fa-solid fa-fire i-tooltip-2"/> 0 BMO
                    </dd>
                  </dl>

                  // TODO(dora): Priority Fee / Tip
                  <dl className="row">
                    <dt className="col-sm-3 col-lg-2 text-muted">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title="User-defined tips sent to validator for transaction priority/inclusion."
                      >
                        <i className="fa-solid fa-info-circle"/>
                      </span>
                      Priority Fee / Tip
                    </dt>
                    <dd className="col-sm-9 col-lg-10">0 BMO</dd>
                  </dl>

                  // TODO(dora): block reward
                   <hr/>
                  <dl className="row">
                    <dt className="col-sm-3 col-lg-2 text-muted">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title="Amount of distributed reward. Miners receive a static block reward + Tx fees + uncle fees."
                      >
                        <i className="fa-solid fa-info-circle"/>
                      </span>
                      BMO Mania Reward
                    </dt>
                    <dd className="col-sm-9 col-lg-10">0.5 BMO</dd>
                  </dl>
                  <dl className="row">
                    <dt className="col-sm-3 col-lg-2 text-muted">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title="Amount of distributed reward. Miners receive a static block reward + Tx fees + uncle fees."
                      >
                        <i className="fa-solid fa-info-circle"/>
                      </span>
                      Validator Reward
                    </dt>
                    <dd className="col-sm-9 col-lg-10">1.0175886759375 BMO</dd>
                  </dl>
                  <dl className="row">
                    <dt className="col-sm-3 col-lg-2 text-muted">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title="Amount of distributed reward. Miners receive a static block reward + Tx fees + uncle fees."
                      >
                        <i className="fa-solid fa-info-circle"/>
                      </span>
                      Emission Reward
                    </dt>
                    <dd className="col-sm-9 col-lg-10">0.5 BMO</dd>
                  </dl>
                  */}
                </div>
              </div>
            </div>
          </div>
        </section>
        <BlockTransactions blockNumber={blockNumber} />
      </section>
    </main>
  );
}
