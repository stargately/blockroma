import React from "react";
import { shortenHash } from "@/shared/common/shorten-hash";
import { normalizeTokenValue } from "@/shared/common/normalize-token-value";
import { TickingTs } from "@/shared/explorer-components/ticking-ts";
import { useChainConfig } from "@/shared/common/use-chain-config";
import { assetURL } from "onefx/lib/asset-url";
import { Status } from "../../../__generated__/globalTypes";

type AdrTx = {
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
  status: Status | null;
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
  tx?: AdrTx | null;
  selfAddressHash: string;
};

export const AdrTransactionItem: React.FC<Props> = ({
  tx,
  selfAddressHash,
}) => {
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
          <span>
            <MaybeClickableAddress
              hash={tx.fromAddressHash}
              self={selfAddressHash}
            />
            →
            <MaybeClickableAddress
              hash={tx.toAddressHash}
              self={selfAddressHash}
            />
          </span>
          <span className="d-flex flex-md-row flex-column mt-3 mt-md-0">
            <span className="tile-title">
              {normalizeTokenValue(tx.value)} {chainConfig.symbol}
            </span>
            <span className="ml-0 ml-md-1 text-nowrap">
              {tx.gasUsed} TX Fee
            </span>
          </span>
          {/* Transfer */}
        </div>
        {/* Block info */}
        <div className="col-md-3 col-lg-2 d-flex flex-row flex-md-column flex-nowrap justify-content-center text-md-right mt-3 mt-md-0 tile-bottom">
          <span className="mr-2 mr-md-0 order-1">
            <a href={assetURL(`block/${tx.blockNumber}`)}>Block #{tx.blockNumber}</a>
          </span>
          <TickingTs
            className="mr-2 mr-md-0 order-2"
            timestamp={tx.timestamp}
            inTile={true}
          />
          <span className="mr-2 mr-md-0 order-0 order-md-3">
            {tx.toAddressHash === selfAddressHash ? (
              <span
                data-test="transaction_type"
                className="badge badge-success tile-badge"
              >
                IN
              </span>
            ) : (
              <span
                data-test="transaction_type"
                className="badge badge-danger tile-badge"
              >
                OUT
              </span>
            )}
          </span>
        </div>
      </div>
    </div>
  );
};

function MaybeClickableAddress({
  hash,
  self,
}: {
  hash: string;
  self: string;
}): JSX.Element {
  const addr = (
    <span data-address-hash={hash}>
      <span className="d-none d-md-none d-xl-inline">{hash}</span>
      <span className="d-md-inline-block d-xl-none">{shortenHash(hash)}</span>
    </span>
  );
  if (hash !== self) {
    return (
      <a data-test="address_hash_link" href={assetURL(`address/${hash}`)}>
        {addr}
      </a>
    );
  }

  return addr;
}
