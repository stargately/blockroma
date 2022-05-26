import React from "react";
import { shortenHash } from "@/shared/common/shorten-hash";
import { normalizeTokenValue } from "@/shared/common/normalize-token-value";
import { TickingTs } from "@/shared/explorer-components/ticking-ts";
import { useChainConfig } from "@/shared/common/use-chain-config";
import { assetURL } from "onefx/lib/asset-url";

export type Tx = {
  id: string;
  timestamp: any | null;
  hash: any | null;
  blockNumber: number | null;
  value: string | null;
  gasUsed: string | null;
  cumulativeGasUsed: string | null;
  error: string | null;
  fromAddressHash: any | null;
  toAddressHash: any | null;
  status: any;
  gas: string | null;
  gasPrice: string | null;
  index: number | null;
  input: any | null;
  nonce: number | null;
  r: string | null;
  s: string | null;
  v: string | null;
};

type Props = {
  tx?: Tx | undefined | null;
};

export const TxTransactionItem: React.FC<Props> = ({ tx }) => {
  const chainConfig = useChainConfig();
  if (!tx) {
    return <></>;
  }
  return (
    <div
      className="tile tile-type-transaction fade-in tile-status--success"
      data-test="transaction"
      data-identifier-hash={tx.hash}
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
              Transaction
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
                href={assetURL(`tx/${tx.hash}`)}
              >
                {tx.hash}
              </a>
              <div className="bs-label method ml-1">Transfer</div>
            </div>
          </span>
          <span>
            <a
              data-test="address_hash_link"
              href={assetURL(`address/${tx.fromAddressHash}`)}
            >
              <span data-address-hash={tx.fromAddressHash}>
                <span className="d-none d-md-none d-xl-inline">
                  {tx.fromAddressHash}
                </span>
                <span className="d-md-inline-block d-xl-none">
                  {shortenHash(tx.fromAddressHash)}
                </span>
              </span>
            </a>
            →
            <a
              data-test="address_hash_link"
              href={assetURL(`address/${tx.toAddressHash}`)}
            >
              <span data-address-hash={tx.toAddressHash}>
                <span className="d-none d-md-none d-xl-inline">
                  {tx.toAddressHash}
                </span>
                <span className="d-md-inline-block d-xl-none">
                  {shortenHash(tx.toAddressHash)}
                </span>
              </span>
            </a>
          </span>
          <span className="d-flex flex-md-row flex-column mt-3 mt-md-0">
            <span className="tile-title">
              {normalizeTokenValue(tx.value)} {chainConfig.symbol}
            </span>
            <span className="ml-0 ml-md-1 text-nowrap">
              {tx.gasPrice} TX Fee
            </span>
          </span>
          {/* Transfer */}
        </div>
        {/* Block info */}
        <div className="col-md-3 col-lg-2 d-flex flex-row flex-md-column flex-nowrap justify-content-center text-md-right mt-3 mt-md-0 tile-bottom">
          <span className="mr-2 mr-md-0 order-1">
            <a href={assetURL(`block/${tx.blockNumber}`)}>
              Block #{tx.blockNumber}
            </a>
          </span>
          <TickingTs
            className="mr-2 mr-md-0 order-2"
            inTile={true}
            timestamp={tx.timestamp}
          />
        </div>
      </div>
    </div>
  );
};
