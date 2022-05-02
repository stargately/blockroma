import React from "react";
import { assetURL } from "onefx/lib/asset-url";

export const TxTokenTransferItem: React.FC = () => {
  const blockNumber = 123;
  return (
    <div
      className="tile tile-type-token-transfer fade-in tile-status--success"
      data-test="token-transfer"
      data-identifier-hash="0xa26665a77410800c96ee764adc2f0aa9cfa0c797c2c332ee73422777aa91a89b"
    >
      <div
        className="row tile-body"
        data-selector="token-transfers-toggle"
        data-test="chain_transaction"
      >
        {/* Color Block */}
        <div className="tile-transaction-type-block col-md-2 d-flex flex-row flex-md-column">
          <div>
            <span className="tile-label" data-test="transaction_type">
              Token Transfer
            </span>
          </div>
          <span
            className="tile-status-label ml-2 ml-md-0"
            data-test="transaction_status"
          >
            Success
          </span>
        </div>
        {/* Content */}
        <div className="col-md-7 col-lg-8 d-flex flex-column pr-2 pr-sm-2 pr-md-0">
          <span>
            <div className="text-truncate d-flex">
              <a
                className="text-truncate"
                data-test="transaction_hash_link"
                href="/poa/core/tx/0xa26665a77410800c96ee764adc2f0aa9cfa0c797c2c332ee73422777aa91a89b"
              >
                0xa26665a77410800c96ee764adc2f0aa9cfa0c797c2c332ee73422777aa91a89b
              </a>
              <div className="bs-label method ml-1">CreateEvent</div>
            </div>
          </span>
          <span>
            <a
              data-test="address_hash_link"
              href="/poa/core/address/0x621c2a125ec4a6d8a7c7a655a18a2868d35eb43c"
            >
              <span data-address-hash="0x621C2a125ec4A6D8A7C7A655A18a2868d35eb43C">
                <span className="d-none d-md-none d-xl-inline">
                  0x621C2a125ec4A6D8A7C7A655A18a2868d35eb43C
                </span>
                <span className="d-md-inline-block d-xl-none">
                  0x621c2a–5eb43c
                </span>
              </span>
            </a>
            →
            <a
              data-test="address_hash_link"
              href="/poa/core/address/0xd789a607ceac2f0e14867de4eb15b15c9ffb5859"
            >
              <span
                className="contract-address"
                data-address-hash="0xd789a607CEac2f0E14867de4EB15b15C9FFB5859"
              >
                <span
                  data-toggle="tooltip"
                  data-placement="top"
                  title="0xd789a607CEac2f0E14867de4EB15b15C9FFB5859"
                  data-custom-class
                >
                  <span className="d-none d-md-none d-lg-inline d-xl-inline">
                    ArianeeStore
                  </span>
                  <span className="d-inline d-md-inline d-lg-none d-xl-none">
                    ArianeeS..re
                  </span>
                  <span> (0xd789a6–fb5859)</span>
                </span>
              </span>
            </a>
          </span>
          <span className="d-flex flex-md-row flex-column mt-3 mt-md-0">
            <span className="tile-title">0 BMO</span>
            <span className="ml-0 ml-md-1 text-nowrap">0.003637884 TX Fee</span>
          </span>
          {/* Transfer */}
          <div className="d-flex flex-column mt-2">
            <div
              className="text-nowrap-small-screen row mt-3 mt-sm-0"
              data-test="token_transfer"
            >
              <span
                className="col-xs-12 col-lg-5"
                style={{ display: "inline-table" }}
              >
                <span className="d-inline-block tile-type-token-transfer-short-name">
                  <a
                    data-test="address_hash_link"
                    href="/poa/core/address/0xd789a607ceac2f0e14867de4eb15b15c9ffb5859"
                  >
                    <span
                      className="contract-address"
                      data-address-hash="0xd789a607CEac2f0E14867de4EB15b15C9FFB5859"
                    >
                      0xd789a6–fb5859
                    </span>
                  </a>
                </span>
                →
                <span className="d-inline-block tile-type-token-transfer-short-name">
                  <a
                    data-test="address_hash_link"
                    href="/poa/core/address/0x7d20a8d54f955b4483a66ab335635ab66e151c51"
                  >
                    <span
                      className="contract-address"
                      data-address-hash="0x7d20a8D54F955b4483A66aB335635ab66e151c51"
                    >
                      0x7d20a8–151c51
                    </span>
                  </a>
                </span>
              </span>
              <span className="col-xs-12 col-lg-4 ml-3 ml-sm-0">
                0.02105596608669953
                <a
                  data-test="token_link"
                  href="/poa/core/token/0x55d536e4d6c1993d8ef2e2a4ef77f02088419420"
                >
                  ARIA
                </a>
              </span>
            </div>
            <div
              className="collapse token-transfer-toggle"
              id="transaction-0xa26665a77410800c96ee764adc2f0aa9cfa0c797c2c332ee73422777aa91a89b"
            >
              <div
                className="text-nowrap-small-screen row mt-3 mt-sm-0"
                data-test="token_transfer"
              >
                <span
                  className="col-xs-12 col-lg-5"
                  style={{ display: "inline-table" }}
                >
                  <span className="d-inline-block tile-type-token-transfer-short-name">
                    <a
                      data-test="address_hash_link"
                      href="/poa/core/address/0xd789a607ceac2f0e14867de4eb15b15c9ffb5859"
                    >
                      <span
                        className="contract-address"
                        data-address-hash="0xd789a607CEac2f0E14867de4EB15b15C9FFB5859"
                      >
                        0xd789a6–fb5859
                      </span>
                    </a>
                  </span>
                  →
                  <span className="d-inline-block tile-type-token-transfer-short-name">
                    <a
                      data-test="address_hash_link"
                      href="/poa/core/address/0x3ae23695016bda778c24d62fb607b3d285eed4ae"
                    >
                      <span
                        className="contract-address"
                        data-address-hash="0x3ae23695016bDA778c24d62fB607B3D285eEd4Ae"
                      >
                        0x3ae236–eed4ae
                      </span>
                    </a>
                  </span>
                </span>
                <span className="col-xs-12 col-lg-4 ml-3 ml-sm-0">
                  0.08422386434679812
                  <a
                    data-test="token_link"
                    href="/poa/core/token/0x55d536e4d6c1993d8ef2e2a4ef77f02088419420"
                  >
                    ARIA
                  </a>
                </span>
              </div>
              <div
                className="text-nowrap-small-screen row mt-3 mt-sm-0"
                data-test="token_transfer"
              >
                <span
                  className="col-xs-12 col-lg-5"
                  style={{ display: "inline-table" }}
                >
                  <span className="d-inline-block tile-type-token-transfer-short-name">
                    <a
                      data-test="address_hash_link"
                      href="/poa/core/address/0xd789a607ceac2f0e14867de4eb15b15c9ffb5859"
                    >
                      <span
                        className="contract-address"
                        data-address-hash="0xd789a607CEac2f0E14867de4EB15b15C9FFB5859"
                      >
                        0xd789a6–fb5859
                      </span>
                    </a>
                  </span>
                  →
                  <span className="d-inline-block tile-type-token-transfer-short-name">
                    <a
                      data-test="address_hash_link"
                      href="/poa/core/address/0xa79b29ad7e0196c95b87f4663ded82fbf2e3add8"
                    >
                      <span
                        className="contract-address"
                        data-address-hash="0xA79B29AD7e0196C95B87f4663ded82Fbf2E3ADD8"
                      >
                        0xa79b29–e3add8
                      </span>
                    </a>
                  </span>
                </span>
                <span className="col-xs-12 col-lg-4 ml-3 ml-sm-0">
                  0.04211193217339906
                  <a
                    data-test="token_link"
                    href="/poa/core/token/0x55d536e4d6c1993d8ef2e2a4ef77f02088419420"
                  >
                    ARIA
                  </a>
                </span>
              </div>
            </div>
          </div>
          <div className="token-tile-view-more">
            <a
              data-selector="token-transfer-open"
              data-test="token_transfers_expansion"
              data-toggle="collapse"
              href="#transaction-0xa26665a77410800c96ee764adc2f0aa9cfa0c797c2c332ee73422777aa91a89b"
            >
              View More Transfers
            </a>
            <a
              className="d-none"
              data-selector="token-transfer-close"
              data-toggle="collapse"
              href="#transaction-0xa26665a77410800c96ee764adc2f0aa9cfa0c797c2c332ee73422777aa91a89b"
            >
              View Less Transfers
            </a>
          </div>
        </div>
        {/* Block info */}
        <div className="col-md-3 col-lg-2 d-flex flex-row flex-md-column flex-nowrap justify-content-center text-md-right mt-3 mt-md-0 tile-bottom">
          <span className="mr-2 mr-md-0 order-1">
            <a href={assetURL(`block/${blockNumber}`)}>Block #{blockNumber}</a>
          </span>
          <span
            className="mr-2 mr-md-0 order-2"
            in-tile="true"
            data-from-now="2022-03-27 07:15:05.000000Z"
          >
            38 seconds ago
          </span>
        </div>
      </div>
    </div>
  );
};
