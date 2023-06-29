import React, { useState } from "react";
import { useParams } from "react-router";
import { normalizeTokenValue } from "@/shared/common/normalize-token-value";
import { AddressTransactionsContainer } from "@/shared/address-details-container/address-transactions-container";
import { useChainConfig } from "@/shared/common/use-chain-config";
import { t } from "onefx/lib/iso-i18n";
import { assetURL } from "onefx/lib/asset-url";
import { useGetAddrDetails } from "@/shared/address-details-container/hooks/use-get-addr-details";
import { QrModal } from "./components/qr-modal";
import { CopyAddress } from "../explorer-components/copy-address";

export interface Addr {
  fetchedCoinBalance?: string | null;
  fetchedCoinBalanceBlockNumber?: number | null;
  nonce?: number | null;
  hash?: any;
  hashQr?: string | null;
  numTxs?: number | null;
  gasUsed?: number | null;
}

export const AddressDetailsContainer: React.FC = () => {
  const chainConfig = useChainConfig();
  const params = useParams<{ addressHash: string }>();
  const { addressHash: rawAddressHash } = params;
  const addressHash = rawAddressHash.toLowerCase();
  const [qrModalOpen, setQrModalOpen] = useState(false);

  const { data, loading, error, refetch } = useGetAddrDetails({
    hash: addressHash,
  });
  if (loading) {
    return <></>;
  }
  if (error && error.graphQLErrors[0]?.extensions?.code !== "NOT_FOUND") {
    // TODO(dora):
    return <button onClick={refetch}>error</button>;
  }
  const addr = data?.address;

  if (!addr) {
    return <></>;
  }

  return (
    <main className="js-ad-dependant-pt pt-5">
      <p className="alert alert-info" role="alert" />
      <p className="alert alert-danger" role="alert" />
      <section className="container">
        <section
          className="address-overview"
          data-page="address-details"
          data-page-address-hash={addressHash}
          data-async-counters={`/address-counters?id=${addressHash}`}
        >
          <div className="row js-ad-dependant-mb-2 js-ad-dependant-mb-5-reverse mb-2">
            <div className="col-md-12 js-ad-dependant-mb-2 mb-2">
              <div className="card js-ad-dependant-mb-2 mb-2">
                <div className="card-body fs-14" style={{ lineHeight: "31px" }}>
                  <h1 className="card-title lg-card-title mb-2-desktop">
                    <div className="title-with-label">Address Details</div>
                    <span className="overview-title-buttons float-right">
                      <CopyAddress addressHash={addr.hash} />
                      <span
                        className="overview-title-item"
                        data-target="#qrModal"
                        data-toggle="modal"
                        onClick={() => setQrModalOpen(true)}
                      >
                        <span
                          className="btn-qr-icon i-tooltip-2"
                          data-toggle="tooltip"
                          data-placement="top"
                          title="QR Code"
                          aria-label="Show QR Code"
                        >
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            viewBox="0 0 32.5 32.5"
                            width={32}
                            height={32}
                          >
                            <path
                              fillRule="evenodd"
                              d="M22.5 24.5v-2h2v2h-2zm-1-4v-1h1v1h-1zm1-3h2v2h-2v-2zm1-2h-5a1 1 0 0 1-1-1v-5a1 1 0 0 1 1-1h5a1 1 0 0 1 1 1v5a1 1 0 0 1-1 1zm-1-5h-3v3h3v-3zm-8 14h-5a1 1 0 0 1-1-1v-5a1 1 0 0 1 1-1h5a1 1 0 0 1 1 1v5a1 1 0 0 1-1 1zm-1-5h-3v3h3v-3zm1-4h-5a1 1 0 0 1-1-1v-5a1 1 0 0 1 1-1h5a1 1 0 0 1 1 1v5a1 1 0 0 1-1 1zm-1-5h-3v3h3v-3zm6 9h-2v-2h2v2zm1 1h-1v-1h1v1zm0 1v1h-1v-1h1zm-1 3h-2v-2h2v2z"
                            />
                          </svg>
                        </span>
                      </span>
                    </span>
                  </h1>
                  <h3
                    className="address-detail-hash-title mb-4 "
                    data-test="address_detail_hash"
                  >
                    {addressHash}
                  </h3>
                  <dl className="row">
                    <dt className="col-sm-4 col-md-4 col-lg-3 text-muted">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title={`Address balance in ${chainConfig.symbol} (doesn't include ERC20, ERC721, ERC1155 tokens).`}
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Balance
                    </dt>
                    <dd
                      className="col-sm-8 col-md-8 col-lg-9"
                      data-test="address_balance"
                    >
                      {normalizeTokenValue(addr.fetchedCoinBalance)}{" "}
                      {chainConfig.symbol}
                      {/*

                      // TODO(dora) coin balance


                                            <span className="address-current-balance">
                        (
                        <span
                          data-wei-value={45253703500000000000}
                          data-usd-exchange-rate="0.0188876"
                          data-placement="top"
                          data-toggle="tooltip"
                          data-html="true"
                          data-original-title="@ 0.0188876/BMO"
                        >
                          $0.854734 USD
                        </span>
                        )
                      </span>
                      */}
                    </dd>
                  </dl>
                  <dl className="row" data-test="outside_of_dropdown">
                    <dt className="col-sm-4 col-md-4 col-lg-3 text-muted">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title="All tokens in the account and total value."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Tokens
                    </dt>
                    <dd
                      className="col-sm-8 col-md-8 col-lg-9"
                      data-test="address_tokens"
                      data-selector="balance-card"
                    >
                      <div
                        className="address-current-balance"
                        data-token-balance-dropdown
                        data-api_path={`/address/${addressHash}/token-balances`}
                      >
                        <div className="d-flex">
                          <span
                            data-tokens-count
                            style={{ lineHeight: "31px" }}
                          >
                            0 tokens
                          </span>
                          <div
                            className="dropdown-menu dropdown-menu-right token-balance-dropdown p-0"
                            aria-labelledby="dropdown-tokens"
                          >
                            <div data-dropdown-items className="dropdown-items">
                              <div className="position-relative">
                                <svg
                                  className="position-absolute dropdown-search-icon"
                                  viewBox="0 0 16 17"
                                  xmlns="http://www.w3.org/2000/svg"
                                  width={16}
                                  height={17}
                                >
                                  <path
                                    fill="#7DD79F"
                                    fillRule="evenodd"
                                    d="M15.713 15.727a.982.982 0 0 1-1.388 0l-2.289-2.29C10.773 14.403 9.213 15 7.5 15A7.5 7.5 0 1 1 15 7.5c0 1.719-.602 3.284-1.575 4.55l2.288 2.288a.983.983 0 0 1 0 1.389zM7.5 2a5.5 5.5 0 1 0 0 11 5.5 5.5 0 1 0 0-11z"
                                  />
                                </svg>
                                <input
                                  className="w-100 dropdown-search-field"
                                  id="token_search_name"
                                  name="token_search[name]"
                                  placeholder="Search tokens"
                                  type="text"
                                  data-filter-dropdown-tokens
                                />
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                    </dd>
                  </dl>
                  <dl className="row address-transactions-count-item">
                    <dt className="col-sm-4 col-md-4 col-lg-3 text-muted">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title="Number of transactions related to this address."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      {t("nav.txs")}
                    </dt>
                    <dd
                      className="col-sm-8 col-md-8 col-lg-9"
                      data-test="address_transaction_count"
                    >
                      <a
                        href="#txs"
                        className="page-link bs-label large btn-no-border-link-to-tems"
                        data-selector="transaction-count"
                      >
                        {addr.numTxs ?? 0} {t("nav.txs")}
                      </a>
                    </dd>
                  </dl>
                  <dl className="row address-transfers-count-item">
                    <dt className="col-sm-4 col-md-4 col-lg-3 text-muted">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title="Number of transfers to/from this address."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Transfers
                    </dt>
                    <dd
                      className="col-sm-8 col-md-8 col-lg-9"
                      data-test="address_transfer_count"
                    >
                      <a
                        href={assetURL(
                          `address/${addressHash}/token-transfers#transfers`
                        )}
                        className="page-link bs-label large btn-no-border-link-to-tems"
                        data-selector="transfer-count"
                      >
                        0 Transfers
                      </a>
                    </dd>
                  </dl>
                  <dl className="row address-nonce-item">
                    <dt className="col-sm-4 col-md-4 col-lg-3 text-muted">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title="Nonce is transaction count of the account."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Nonce
                    </dt>

                    <dd
                      className="col-sm-8 col-md-8 col-lg-9"
                      data-test="address_transaction_count"
                    >
                      <span data-selector="Nonce">{addr.nonce ?? 0}</span>
                    </dd>
                  </dl>
                  <dl className="row address-gas-used-item">
                    <dt className="col-sm-4 col-md-4 col-lg-3 text-muted">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title="Gas used by the address."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Gas Used
                    </dt>
                    <dd
                      className="col-sm-8 col-md-8 col-lg-9"
                      data-test="address_gas_used"
                    >
                      <span data-selector="gas-usage-count">
                        {addr.gasUsed ?? 0}
                      </span>
                    </dd>
                  </dl>
                  <dl className="row">
                    <dt className="col-sm-4 col-md-4 col-lg-3 text-muted">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title="Block number in which the address was updated."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Last Balance Update
                    </dt>
                    <dd
                      className="col-sm-8 col-md-8 col-lg-9"
                      data-test="address_last_balance_update"
                    >
                      <a
                        className="tile-title-lg"
                        href={assetURL(
                          `block/${addr.fetchedCoinBalanceBlockNumber}`
                        )}
                      >
                        {addr.fetchedCoinBalanceBlockNumber ?? 0}
                      </a>
                    </dd>
                  </dl>
                  <dl
                    className="row address-validation-count-item"
                    style={{ display: "none" }}
                  >
                    <dt className="col-sm-4 col-md-4 col-lg-3 text-muted">
                      <span
                        className="i-tooltip-2 "
                        data-boundary="window"
                        data-container="body"
                        data-html="true"
                        data-placement="top"
                        data-toggle="tooltip"
                        title="Number of blocks validated by this validator."
                      >
                        <i className="fa-solid fa-info-circle" />{" "}
                      </span>
                      Blocks Validated
                    </dt>
                    <dd
                      className="col-sm-8 col-md-8 col-lg-9"
                      data-test="address_blocks_validated"
                    >
                      <span data-selector="validation-count"></span>
                    </dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>
        </section>

        <QrModal
          addressHash={addressHash}
          open={qrModalOpen}
          hashQr={addr.hashQr}
          onClose={() => setQrModalOpen(false)}
        />

        <div className="ad-container mb-2" style={{}}>
          <div className="coinzilla" data-zone="C-26660bf627543e46851" />
        </div>

        <AddressTransactionsContainer addressHash={addressHash} />
      </section>
    </main>
  );
};
