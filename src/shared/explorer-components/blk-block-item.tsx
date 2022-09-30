import React from "react";
import { shortenHash } from "@/shared/common/shorten-hash";
import { assetURL } from "onefx/lib/asset-url";
import { TickingTs } from "./ticking-ts";

export type Blk = {
  consensus: boolean | null;
  difficulty: string | null;
  gasLimit: string | null;
  gasUsed: string | null;
  hash: any | null;
  miner: any | null;
  nonce: any | null;
  number: number | null;
  parentHash: any | null;
  size: number | null;
  numTxs: number | null;
  timestamp: any | null;
  totalDifficulty: string | null;
};

type Props = {
  blk: Blk;
};

export const BlkBlockItem: React.FC<Props> = ({ blk }) => {
  const gasUsedPercent = (
    (parseFloat(blk.gasUsed ?? "0") / parseFloat(blk.gasLimit ?? "0")) *
    100
  ).toFixed(1);

  return (
    <div
      className="tile tile-type-block fade-up"
      data-selector="block-tile"
      data-block-number={blk.number}
      data-block-hash={blk.hash}
      style={{}}
    >
      <div className="row">
        <div className="tile-transaction-type-block col-md-2 d-flex flex-row flex-md-column">
          <a
            className="tile-label"
            data-selector="block-number"
            href={assetURL(`block/${blk.number}`)}
          >
            #{blk.number}
          </a>
          <span
            className="tile-status-label font-weight-400"
            data-test="transaction_type"
          >
            Block
          </span>
        </div>
        <div className="col-md-6 col-lg-7">
          <div>
            {/* transactions */}
            <span className="mr-2">{blk.numTxs} transactions</span>
            {/* size */}
            <span className="mr-2">
              {" "}
              {Number(blk.size).toLocaleString()} bytes{" "}
            </span>
            {/* age */}
            <TickingTs timestamp={blk.timestamp} />
          </div>
          <div className="text-nowrap text-truncate mt-3 mt-md-0">
            {/* validator */}
            Validator{" "}
            <a
              data-test="address_hash_link"
              href={assetURL(`address/${blk.miner}`)}
            >
              <span data-address-hash={blk.miner}>
                <span
                  data-toggle="tooltip"
                  data-placement="top"
                  title={blk.miner}
                  data-custom-class
                >
                  <span>{shortenHash(blk.miner)}</span>
                </span>
              </span>
            </a>
          </div>

          {/*
          TODO(dora): block reward

          <div className="text-nowrap text-truncate mt-3 mt-md-0">
            Reward
            <span className="ml-2">2 BMO</span>
          </div>
          */}
        </div>
        <div className="col-md-4 col-lg-3 text-md-right d-flex flex-column align-items-md-end justify-content-md-end mt-3 mt-md-0">
          {/* Priority Fee TODO(dora) */}
          {/* <span> 0 BMO Priority Fees </span> */}
          {/* Burnt Fees */}
          {/* <span> 0 BMO Burnt Fees </span> */}
          {/* Gas Limit */}
          <span>
            {" "}
            {parseFloat(blk.gasLimit || "0").toLocaleString()} Gas Limit{" "}
          </span>
          {/* Gas Used */}
          <div className="mr-3 mr-md-0">
            {blk.gasUsed} ({gasUsedPercent}%) Gas Used
          </div>
          {/* Progress bar */}
          <div className="progress">
            <div
              className="progress-bar"
              role="progressbar"
              style={{ width: "0%" }}
              aria-valuenow={parseFloat(gasUsedPercent) * 100}
              aria-valuemin={0}
              aria-valuemax={100}
            ></div>
          </div>
        </div>
      </div>
    </div>
  );
};
