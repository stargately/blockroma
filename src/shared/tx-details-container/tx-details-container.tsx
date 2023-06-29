import * as React from "react";
import { useQueryTx } from "@/shared/tx-details-container/hooks/use-query-tx";
import { useParams } from "react-router";
import { TxStatus } from "@/shared/explorer-components/tx-status";
import { normalizeTokenValue } from "@/shared/common/normalize-token-value";
import { getGasUsedPercent } from "@/shared/common/get-gas-used-percent";
import { CopyToClipboard } from "@/shared/explorer-components/copy-to-clipboard";
import { TickingTs } from "@/shared/explorer-components/ticking-ts";
import format from "date-fns/format";
import { useChainConfig } from "@/shared/common/use-chain-config";
import { t } from "onefx/lib/iso-i18n";
import { assetURL } from "onefx/lib/asset-url";
import { DataInput } from "../explorer-components/data-input";
import { TokenTransferContainer } from "../token-transfer-container/token-transfer-container";
import { divDecimals } from "../common/div-decimals";

export function TxDetailsContainer(): JSX.Element {
  const params = useParams<{ txHash: string }>();
  const chainConfig = useChainConfig();

  const { data, loading, error, refetch } = useQueryTx({ hash: params.txHash });
  if (loading) {
    // TODO(dora):
    return <></>;
  }
  if (error && error.graphQLErrors[0]?.extensions?.code === "NOT_FOUND") {
    // TODO(dora):
    return (
      <>
        Transaction not found
        <button onClick={refetch}>Refetch</button>
      </>
    );
  }

  const tx = data?.transaction;
  let txFee = "0";
  try {
    txFee = (Number(tx?.gasUsed ?? 0) * Number(tx?.gasPrice ?? 0)).toString();
  } catch (err) {
    console.error(`failed to calc txFee: ${err}`);
  }

  return (
    <main className="pt-4">
      <p className="alert alert-info" role="alert" />
      <p className="alert alert-danger" role="alert" />
      <section className="container">
        <section
          className="fs-14"
          data-page="transaction-details"
          data-page-transaction-hash={tx?.hash}
        >
          <div className="row">
            <div className="col-md-12">
              <div className="card mb-3">
                <div className="card-body">
                  <h1 className="card-title margin-bottom-1">
                    <div
                      style={{
                        display: "inline-block",
                        verticalAlign: "bottom",
                        lineHeight: "25px",
                      }}
                    >
                      {t("tx.tx_details")}
                    </div>
                  </h1>

                  <dl className="row">
                    <dt className="col-sm-3 col-lg-2 text-muted">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title={t("tx.tx_hash.tip")}
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      {t("tx.tx_hash")}
                    </dt>
                    <dd
                      className="col-sm-9 col-lg-10"
                      style={{ wordBreak: "break-all" }}
                    >
                      <span
                        className="transaction-details-address"
                        data-test="transaction_detail_hash"
                      >
                        {tx?.hash}{" "}
                      </span>

                      <CopyToClipboard
                        value={tx?.hash}
                        reason={t("tx.copy_tx_hash")}
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
                        title={t("tx.result.tip")}
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      {t("tx.result")}
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      <TxStatus status={tx?.status} />
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
                        title={t("tx.status.tip")}
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      {t("tx.status")}
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      <span className="mr-4">
                        <span data-transaction-status="Confirmed">
                          <div className="bs-label success large">
                            Confirmed
                          </div>
                        </span>

                        {/*

                         //TODO(dora) confirmed by how many blocks?
                         <span className="bs-label large ml-2 confirmations-label">
                         Confirmed by{" "}
                         <span data-selector="block-confirmations">594</span>{" "}
                         blocks
                         </span>
                         */}
                      </span>
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
                        title={t("tx.block.tip")}
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      {t("tx.block")}
                    </dt>
                    <dd
                      className="col-sm-9 col-lg-10"
                      data-selector="block-number"
                    >
                      <a
                        className="transaction__link"
                        href={assetURL(`block/${tx?.blockNumber}`)}
                      >
                        {tx?.blockNumber}
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
                        title={t("tx.timestamp.tip")}
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      {t("tx.timestamp")}
                    </dt>
                    <dd
                      className="col-sm-9 col-lg-10"
                      data-selector="block-timestamp"
                    >
                      <i className="fa-regular fa-clock" />{" "}
                      <span>
                        <TickingTs timestamp={tx?.timestamp} />
                        {" | "}
                        {format(
                          new Date(tx?.timestamp),
                          "MMM-d-y hh:mm:ss a x"
                        )}{" "}
                        UTC
                      </span>
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
                        title={t("tx.from.tip")}
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      {t("tx.from")}
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      <a href={assetURL(`address/${tx?.fromAddressHash}`)}>
                        {tx?.fromAddressHash}
                      </a>

                      <CopyToClipboard
                        reason="Copy From Address"
                        value={tx?.fromAddressHash}
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
                        title={t("tx.to.tip")}
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      {t("tx.to")}
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      <a href={assetURL(`address/${tx?.toAddressHash}`)}>
                        {tx?.toAddressHash}
                      </a>

                      <CopyToClipboard
                        reason="Copy To Address"
                        value={tx?.toAddressHash}
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
                        title={t("tx.value.tip")}
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      {t("tx.value")}
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      {" "}
                      {normalizeTokenValue(tx?.value)} {chainConfig.symbol}
                      {/*
                      TODO(dora): coin balance price

                      (
                      <span
                        data-wei-value={252981808720000000}
                        data-usd-exchange-rate={1.0}
                      >
                        $0.252982 USD
                      </span>
                      )
                      */}
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
                        title={t("tx.fee.tip")}
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      {t("tx.fee")}
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      {txFee} Gwei
                      {/*
                      TODO(dora):
                      (
                      <span
                        data-wei-value={tx?.gas}
                        data-usd-exchange-rate={1.0}
                      >
                        $0.000042 USD
                      </span>
                      )

                      */}
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
                        title={t("tx.gas_price.tip")}
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      {t("tx.gas_price")}
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      {" "}
                      {divDecimals(tx?.gasPrice, chainConfig.decimals)}{" "}
                      {chainConfig.symbol}{" "}
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
                        title={t("tx.tx_type.tip")}
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      {t("tx.tx_type")}
                    </dt>
                    <dd className="col-sm-9 col-lg-10"> 0</dd>
                  </dl>
                  <hr />
                  <dl className="row">
                    <dt className="col-sm-3 col-lg-2 text-muted">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title={t("tx.gas_limit.tip")}
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      {t("tx.gas_limit")}
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      {" "}
                      {Number(tx?.gas).toLocaleString()}{" "}
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
                        title={t("tx.max_fee_per_gas.tip")}
                      >
                        <i className="fa-solid fa-info-circle"></i>{" "}
                      </span>
                      {t("tx.max_fee_per_gas")}
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      {" "}
                      {divDecimals(tx?.maxFeePerGas)} {chainConfig.symbol}
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
                        title="User defined maximum fee (tip) per unit of gas paid to validator for transaction prioritization."
                      >
                        <i className="fa-solid fa-info-circle"></i>{" "}
                      </span>
                      Max Priority Fee per Gas
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      {" "}
                      {divDecimals(tx?.maxPriorityFeePerGas)}{" "}
                      {chainConfig.symbol}
                    </dd>
                  </dl>

                  <dl className="row">
                    <dt className="col-sm-3 col-lg-2 text-muted transaction-gas-used">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title="tx.gas_used.tip"
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      {t("tx.gas_used")}
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      {" "}
                      {tx?.gasUsed} | {getGasUsedPercent(tx?.gasUsed, tx?.gas)}%
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
                        title={t("tx.nonce.tip")}
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      {t("tx.nonce")}
                      <span
                        className="index-label ml-2"
                        data-toggle="tooltip"
                        title="Index position of Transaction in the block."
                      >
                        {t("tx.pos")}
                      </span>
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      {" "}
                      {tx?.nonce}
                      <span className="index-label ml-2">{tx?.index}</span>{" "}
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
                        title=""
                        data-original-title={t("tx.raw_input.tip")}
                      >
                        <i className="fa-solid fa-info-circle"></i>{" "}
                      </span>
                      {t("tx.raw_input")}
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      <DataInput input={tx?.input} />
                    </dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* TODO(dora): more tx info */}

        <div className="card">
          <div id="transaction-tabs" className="card-tabs js-card-tabs">
            <a
              className="card-tab active"
              href="/xdai/mainnet/tx/0x5a2ee4a79e55a0bc0c7f5c187fcc83ba8deaaa69614d9da2df3c6a42ada7d87c/token-transfers"
            >
              Token Transfers
            </a>
          </div>

          <div className="card-body">
            <TokenTransferContainer transactionHash={tx?.hash} />
          </div>
        </div>
      </section>
    </main>
  );
}
