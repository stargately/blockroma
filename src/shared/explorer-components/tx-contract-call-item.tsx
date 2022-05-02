import React from "react";
import { assetURL } from "onefx/lib/asset-url";

export const TxContractCallItem: React.FC = () => {
  const blockNumber = 123;
  return (
    <div
      className="tile tile-type-contract-call fade-in tile-status--success"
      data-test="contract-call"
      data-identifier-hash="0x9fdc11e3189b0903c93efe9c0afaa01c2e1fcacb03e47db1de544ba409e0226a"
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
              Contract Call
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
                href="/poa/core/tx/0x9fdc11e3189b0903c93efe9c0afaa01c2e1fcacb03e47db1de544ba409e0226a"
              >
                0x9fdc11e3189b0903c93efe9c0afaa01c2e1fcacb03e47db1de544ba409e0226a
              </a>
              <div className="bs-label method ml-1">0x344639b4</div>
            </div>
          </span>
          <span>
            <a
              data-test="address_hash_link"
              href="/poa/core/address/0x0cdde70d140dd0928b21520162c3031caa0ceb5a"
            >
              <span data-address-hash="0x0cdDE70d140DD0928b21520162C3031Caa0cEb5a">
                <span className="d-none d-md-none d-xl-inline">
                  0x0cdDE70d140DD0928b21520162C3031Caa0cEb5a
                </span>
                <span className="d-md-inline-block d-xl-none">
                  0x0cdde7–0ceb5a
                </span>
              </span>
            </a>
            →
            <a
              data-test="address_hash_link"
              href="/poa/core/address/0x981c44040cb6150a2b8a7f63fb182760505bf666"
            >
              <span
                className="contract-address"
                data-address-hash="0x981C44040CB6150a2b8A7F63FB182760505bf666"
              >
                <span className="d-none d-md-none d-xl-inline">
                  0x981C44040CB6150a2b8A7F63FB182760505bf666
                </span>
                <span className="d-md-inline-block d-xl-none">
                  0x981c44–5bf666
                </span>
              </span>
            </a>
          </span>
          <span className="d-flex flex-md-row flex-column mt-3 mt-md-0">
            <span className="tile-title">0 BMO</span>
            <span className="ml-0 ml-md-1 text-nowrap">
              0.003077109375 TX Fee
            </span>
          </span>
          {/* Transfer */}
        </div>
        {/* Block info */}
        <div className="col-md-3 col-lg-2 d-flex flex-row flex-md-column flex-nowrap justify-content-center text-md-right mt-3 mt-md-0 tile-bottom">
          <span className="mr-2 mr-md-0 order-1">
            <a href={assetURL(`block/${blockNumber}`)}>Block #26559654</a>
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
