import { useTokenTransfer } from "@/shared/token-transfer-container/hooks/use-token-transfer";
import * as React from "react";
import { selectTokenTransfer } from "@/shared/token-transfer-container/selectors/select-token-transfer";
import { assetURL } from "onefx/lib/asset-url";
import { useTokenDetails } from "@/shared/token-transfer-container/hooks/use-token-details";
import { selectTokenDetails } from "@/shared/token-transfer-container/selectors/select-token-details";
import { divDecimals } from "@/shared/common/div-decimals";
import { t } from "onefx/lib/iso-i18n";
import { useCallback } from "react";

type Props = {
  transactionHash?: string;
};

export const TokenTransferContainer: React.FC<Props> = ({
  transactionHash,
}) => {
  const { transferData, transferLoading, transferError, transferRefetch } =
    useTokenTransfer(transactionHash);
  const transfers = selectTokenTransfer(transferData);

  const {
    tokenDetailsData,
    tokenDetailsLoading,
    tokenDetailsError,
    tokenDetailsRefetch,
  } = useTokenDetails(transfers?.at(0)?.tokenContractAddress);
  const tokenDetails = selectTokenDetails(tokenDetailsData);

  const refetch = useCallback(() => {
    (async () => {
      await transferRefetch();
      await tokenDetailsRefetch();
    })();
  }, []);

  return (
    <>
      <h2 className="card-title list-title-description">Token Transfers</h2>
      <div className="list-top-pagination-container-wrapper">
        <div
          className="pagination-container position-top "
          data-pagination-container=""
        >
          <ul className="pagination"></ul>
        </div>
      </div>

      {(transferError || tokenDetailsError) && (
        <button
          onClick={refetch}
          className="alert alert-danger col-12 text-left"
        >
          <span className="alert-link">{t("info.err")}</span>
        </button>
      )}

      <div data-empty-response-message="" style={{ display: "none" }}>
        <div
          className="tile tile-muted text-center"
          data-selector="empty-logs-list"
        >
          There are no token transfers for this transaction
        </div>
      </div>
      <div data-items="">
        {!(transferLoading || tokenDetailsLoading) && !transfers?.length && (
          <div data-empty-response-message>
            <div
              className="tile tile-muted text-center"
              data-selector="empty-transactions-list"
            >
              No token transfers in this transaction.
            </div>
          </div>
        )}

        {!(transferLoading || tokenDetailsLoading) &&
          transfers?.map((tf) => (
            <div
              className="tile tile-type-token-transfer fade-in"
              key={tf.transactionHash}
            >
              <div className="row justify-content-end">
                <div className="col-12 col-md-4 col-lg-2 d-flex align-items-center justify-content-start justify-content-lg-center tile-label">
                  Token Transfer
                </div>

                <div className="col-12 col-md-8 col-lg-10 d-flex flex-column text-nowrap">
                  <a
                    className="text-truncate"
                    data-test="transaction_hash_link"
                    href={assetURL(`tx/${tf.transactionHash}`)}
                  >
                    {tf.transactionHash}
                  </a>
                  <span className="text-nowrap">
                    <a
                      data-test="address_hash_link"
                      href={assetURL(`address/${tf.transactionHash}`)}
                    >
                      <span>
                        <span className="d-none d-md-none d-xl-inline">
                          {tf.fromAddress}
                        </span>
                      </span>
                    </a>{" "}
                    →{" "}
                    <a
                      data-test="address_hash_link"
                      href={assetURL(`address/${tf.toAddress}`)}
                    >
                      <span>{tf.toAddress}</span>
                    </a>
                  </span>

                  <span className="tile-title">
                    {divDecimals(tf.amount, tokenDetails?.decimals)}
                    <a
                      data-test="token_link"
                      href={assetURL(`address/${tf.tokenContractAddress}`)}
                    >
                      {" "}
                      {tokenDetails?.symbol}
                    </a>
                  </span>
                </div>
              </div>
            </div>
          ))}
      </div>
      <div
        className="pagination-container  position-bottom"
        data-pagination-container=""
      >
        <ul className="pagination"></ul>
      </div>
    </>
  );
};
