import React from "react";
import {
  Tx,
  TxTransactionItem,
} from "@/shared/explorer-components/tx-transaction-item";

type Props = {
  txs?: Array<Tx | null | undefined> | null | undefined;
};

export function TxsList({ txs }: Props): JSX.Element {
  if (!txs) {
    return <></>;
  }
  return (
    <>
      {/*
      // TODO(dora): pagination
      <div
        className="pagination-container position-top "
        data-pagination-container
      >
        <ul className="pagination">
          <li className="page-item">
            <a
              className="page-link"
              href="#"
              data-first-page-button
              style={{ display: "none" }}
            >
              First
            </a>
          </li>
          <li className="page-item">
            <a className="page-link" href="#" data-prev-page-button>
              <svg xmlns="http://www.w3.org/2000/svg" width={6} height={10}>
                <path
                  fillRule="evenodd"
                  d="M2.358 5l3.357 3.358a.959.959 0 1 1-1.357 1.357L.502 5.859c-.076-.042-.153-.08-.217-.144A.949.949 0 0 1 .011 5a.949.949 0 0 1 .274-.715c.064-.064.142-.102.217-.145L4.358.285a.959.959 0 1 1 1.357 1.357L2.358 5z"
                />
              </svg>
            </a>
          </li>
          <li className="page-item">
            <a className="page-link no-hover" href="#" data-page-number>
              Page 1
            </a>
          </li>
          <li className="page-item">
            <a
              className="page-link"
              href="/block/26559654/transactions"
              data-next-page-button
            >
              <svg xmlns="http://www.w3.org/2000/svg" width={6} height={10}>
                <path
                  fillRule="evenodd"
                  d="M5.715 5.715c-.064.064-.141.102-.217.144L1.642 9.715A.959.959 0 1 1 .285 8.358L3.642 5 .285 1.642A.959.959 0 1 1 1.642.285L5.498 4.14c.075.043.153.081.217.145A.949.949 0 0 1 5.989 5a.949.949 0 0 1-.274.715z"
                />
              </svg>
            </a>
          </li>
        </ul>
      </div>

      */}

      <div data-items="true">
        {txs.map((tx) => tx && <TxTransactionItem key={tx.hash} tx={tx} />)}
      </div>
    </>
  );
}
