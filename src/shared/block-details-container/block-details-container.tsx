import * as React from "react";
import { useQueryBlock } from "@/shared/block-details-container/hooks/use-query-block";
import { useParams } from "react-router";
import { Link } from "react-router-dom";
import { getGasUsedPercent } from "@/shared/common/get-gas-used-percent";
import { CopyToClipboard } from "@/shared/explorer-components/copy-to-clipboard";
import { TickingTs } from "@/shared/explorer-components/ticking-ts";
import format from "date-fns/format";
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
                      Block Details
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
                          Block {blockNumber}
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
                        title="The block height of a particular block is defined as the number of blocks preceding it in the blockchain."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Block Height
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
                        title="Date & time at which block was produced."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Timestamp
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
                        title="The number of transactions in the block."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Transactions
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      <a
                        href="#txs"
                        className="page-link bs-label large btn-no-border-link-to-tems"
                      >
                        {numTxs} Transactions
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
                        title="A block producer who successfully included the block onto the blockchain."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Validator{" "}
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      <span>
                        <a
                          data-test="address_hash_link"
                          href={`/address/${miner}`}
                        >
                          {miner}
                        </a>
                      </span>
                      <CopyToClipboard value={miner} reason="Copy Address" />
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
                        title="Size of the block in bytes."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Size
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
                        title="The SHA256 hash of the block."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Hash
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      {hash}

                      <CopyToClipboard value={hash} reason="Copy Hash" />
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
                        title="The hash of the block from which this block was generated."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Parent Hash
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      {blockNumber > 0 ? (
                        <a
                          className="transaction__link"
                          href={`/block/${blockNumber - 1}`}
                        >
                          {parentHash}
                        </a>
                      ) : (
                        parentHash
                      )}

                      <CopyToClipboard
                        value={parentHash}
                        reason="Copy Parent Hash"
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
                        title="Block difficulty for miner, used to calibrate block generation time (Note: constant in BMO based networks)."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Difficulty
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
                        title="Total difficulty of the chain until this block."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Total Difficulty
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
                        title="The total gas amount used in the block and its percentage of gas filled in the block."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Gas Used
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
                        title="Total gas limit provided by all transactions in the block."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Gas Limit
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
                      Nonce
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
        <div className="ad-container mb-2" style={{}}>
          <div className="coinzilla" data-zone="C-26660bf627543e46851" />
        </div>
        <BlockTransactions blockNumber={blockNumber} />
      </section>
    </main>
  );
}
