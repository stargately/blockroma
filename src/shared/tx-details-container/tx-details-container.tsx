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
import { DataInput } from "../explorer-components/data-input";

export function TxDetailsContainer(): JSX.Element {
  const params = useParams<{ txHash: string }>();
  const chainConfig = useChainConfig();
  function divDecimals(num?: string | null): string {
    if (!num) {
      return "0";
    }
    return (Number(num) / 10 ** chainConfig.decimals)
      .toFixed(20)
      .replace(/\.?0*$/, "");
  }
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
                      Transaction Details
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
                        title="Unique character string (TxID) assigned to every verified transaction."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Transaction Hash
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
                        reason="Copy Transaction Hash"
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
                        title="Current transaction state: Success, Failed (Error), or Pending (In Process)"
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Result
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
                        title="The status of the transaction: Confirmed or Unconfirmed."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Status
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
                        title="Block number containing the transaction."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Block
                    </dt>
                    <dd
                      className="col-sm-9 col-lg-10"
                      data-selector="block-number"
                    >
                      <a
                        className="transaction__link"
                        href={`/block/${tx?.blockNumber}`}
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
                        title="Date & time of transaction inclusion, including length of time for confirmation."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Timestamp
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
                        title="Address (external or contract) sending the transaction."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      From
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      <a href={`/address/${tx?.fromAddressHash}`}>
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
                        title="Address (external or contract) receiving the transaction."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      To
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      <a href={`/address/${tx?.toAddressHash}`}>
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
                        title="Value sent in the native token (and USD) if applicable."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Value
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
                        title="Total transaction fee."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Transaction Fee
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
                        title="Price per unit of gas specified by the sender. Higher gas prices can prioritize transaction inclusion during times of high usage."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Gas Price
                    </dt>
                    <dd className="col-sm-9 col-lg-10">
                      {" "}
                      {divDecimals(tx?.gasPrice)} {chainConfig.symbol}{" "}
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
                        title="Transaction type, introduced in EIP-2718."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Transaction Type
                    </dt>
                    <dd className="col-sm-9 col-lg-10"> 0 </dd>
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
                        title="Maximum gas amount approved for the transaction."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Gas Limit
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
                        title="Maximum total amount per unit of gas a user is willing to pay for a transaction, including base fee and priority fee."
                      >
                        <i className="fa-solid fa-info-circle"></i>{" "}
                      </span>
                      Max Fee per Gas
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
                        title="Actual gas amount used by the transaction."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Gas Used by Transaction
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
                        title="Transaction number from the sending address. Each transaction sent from an address increments the nonce by 1."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Nonce
                      <span
                        className="index-label ml-2"
                        data-toggle="tooltip"
                        title="Index position of Transaction in the block."
                      >
                        Position
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
                        data-original-title="Binary data included with the transaction. See input / logs below for additional info."
                      >
                        <i className="fa-solid fa-info-circle"></i>{" "}
                      </span>
                      Raw Input
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
      </section>
    </main>
  );
}
