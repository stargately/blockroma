import React, { useState } from "react";
import { useQueryTokens } from "@/shared/token-container/hooks/use-query-tokens";
import { selectTokens } from "@/shared/token-container/selectors/select-tokens";
import { shortenHash } from "@/shared/common/shorten-hash";
import { useRouter } from "next/router";
import { assetURL } from "@/shared/common/asset-url";
import Link from "next/link";

export const TokensContainer = () => {
  const router = useRouter();
  const searchParams = router.query;
  const querySymbol = searchParams["symbol"];

  const [symbol, setSymbol] = useState(String(querySymbol ?? "").toUpperCase());

  const updateSymbol = (s: string) => {
    // Create a new URLSearchParams object
    const p = new URLSearchParams(location.search);
    // Set the new query value
    p.set("symbol", s);
    // Update the URL with the new query string
    router.push({ search: p.toString() });
    setSymbol(s);
  };

  const { tokensData, tokensFetchMore } = useQueryTokens(symbol);
  const { tokens, currentCursor, hasNextPage } = selectTokens(tokensData);

  const handleKeyUp = async (event: any) => {
    if (event.key === "Enter") {
      const uppercase = String(event.target.value).toUpperCase();
      updateSymbol(uppercase);
      await tokensFetchMore("0", uppercase);
    }
  };

  return (
    <main className="pt-4">
      <p className="alert alert-info" role="alert"></p>
      <p className="alert alert-danger" role="alert"></p>
      <section
        className="container"
        data-page="tokens"
        data-chain-id={10}
        data-display-token-icons="false"
      >
        <div className="ad mb-3" style={{ display: "none" }}>
          <span className="ad-prefix" />:{" "}
          <img className="ad-img-url" width={20} height={20} />{" "}
          <b>
            <span className="ad-name" />
          </b>{" "}
          - <span className="ad-short-description" />{" "}
          <a className="ad-url">
            <b>
              <span className="ad-cta-button" />
            </b>
          </a>
        </div>
        <div className="card">
          <div
            className="card-body"
            data-async-load=""
            data-async-listing="/optimism/mainnet/tokens"
            data-no-self-calls=""
          >
            <h1 className="card-title list-title-description">Tokens</h1>
            <div
              className="list-top-pagination-container-wrapper tokens-list-search-input-outer-container d-flex"
              style={{ float: "right" }}
            >
              <label className="tokens-list-search-input-container tokens mr-3">
                <input
                  data-search-field=""
                  className="form-control tokens-list-search-input search-input"
                  type="text"
                  name="filter"
                  defaultValue={symbol}
                  placeholder="Token symbol"
                  id="search-text-input"
                  onKeyUp={handleKeyUp}
                />
              </label>
              <div
                className="pagination-container position-top "
                data-pagination-container=""
              >
                <ul className="pagination">
                  <li className="page-item">
                    <Link
                      className="page-link no-hover"
                      href=""
                      data-page-number=""
                    >
                      Page {currentCursor}
                    </Link>
                  </li>

                  {hasNextPage && (
                    <li className="page-item">
                      <a
                        className="page-link"
                        data-next-page-button=""
                        onClick={(e) => {
                          e.preventDefault();
                          return tokensFetchMore(
                            tokensData?.tokens?.pageInfo?.endCursor,
                            symbol,
                          );
                        }}
                      >
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          width={6}
                          height={10}
                        >
                          <path
                            fillRule="evenodd"
                            d="M5.715 5.715c-.064.064-.141.102-.217.144L1.642 9.715A.959.959 0 1 1 .285 8.358L3.642 5 .285 1.642A.959.959 0 1 1 1.642.285L5.498 4.14c.075.043.153.081.217.145A.949.949 0 0 1 5.989 5a.949.949 0 0 1-.274.715z"
                          />
                        </svg>
                      </a>
                    </li>
                  )}
                </ul>
              </div>
            </div>
            <div className="addresses-table-container">
              <div className="stakes-table-container">
                <table>
                  <thead>
                    <tr>
                      <th className="stakes-table-th">
                        <div className="stakes-table-th-content">&nbsp;</div>
                      </th>
                      <th className="stakes-table-th">
                        <div>&nbsp;</div>
                      </th>
                      <th className="stakes-table-th">
                        <div className="stakes-table-th-content">Token</div>
                      </th>
                      <th className="stakes-table-th">
                        <div className="stakes-table-th-content">Address</div>
                      </th>
                      {/* <th className="stakes-table-th"> */}
                      {/*  <div className="stakes-table-th-content"> */}
                      {/*    Circulating Market Cap */}
                      {/*  </div> */}
                      {/* </th> */}
                      <th className="stakes-table-th">
                        <div className="stakes-table-th-content">
                          Total Supply
                        </div>
                      </th>
                      {/* <th className="stakes-table-th"> */}
                      {/*  <div className="stakes-table-th-content"> */}
                      {/*    Holders Count */}
                      {/*  </div> */}
                      {/* </th> */}
                    </tr>
                  </thead>
                  <tbody data-items="" data-selector="top-tokens-list">
                    {tokens?.map((t) => (
                      <tr key={t?.contractAddress}>
                        <td className="stakes-td">
                          <span className="color-lighten"></span>
                        </td>

                        <td className="token-icon"></td>

                        <td className="stakes-td">
                          <Link
                            className="text-truncate"
                            data-test="token_link"
                            href={assetURL(`address/${t?.contractAddress}`)}
                          >
                            {t?.name} ({t?.symbol})
                          </Link>
                        </td>

                        <td className="stakes-td">
                          <Link
                            data-test="address_hash_link"
                            href={assetURL(`address/${t?.contractAddress}`)}
                          >
                            <span className="contract-address">
                              <span
                                data-toggle="tooltip"
                                data-placement="top"
                                title=""
                                data-custom-class=""
                              >
                                <span className="d-none d-md-none d-lg-inline d-xl-inline">
                                  {t?.name}
                                </span>
                                <span className="d-inline d-md-inline d-lg-none d-xl-none">
                                  {t?.name}
                                </span>
                                <span>
                                  {" "}
                                  ({shortenHash(t?.contractAddress)})
                                </span>
                              </span>
                            </span>
                          </Link>
                        </td>

                        {/* <td className="stakes-td"> */}
                        {/*  <span data-selector="circulating-market-cap-usd"> */}
                        {/*    N/A */}
                        {/*  </span> */}
                        {/* </td> */}
                        <td className="stakes-td">
                          <span data-test="token_supply">
                            {Number(t?.totalSupply).toLocaleString()}
                          </span>{" "}
                          {t?.symbol}
                        </td>
                        {/* <td className="stakes-td"> */}
                        {/*  <span className="mr-4"> */}
                        {/*    <span data-test="transaction_count">??</span> */}
                        {/*  </span> */}
                        {/* </td> */}
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>

            <div
              className="pagination-container  position-bottom"
              data-pagination-container=""
            >
              <ul className="pagination">
                <li className="page-item">
                  <Link
                    className="page-link no-hover"
                    href=""
                    data-page-number=""
                  >
                    Page {currentCursor}
                  </Link>
                </li>

                {hasNextPage && (
                  <li className="page-item">
                    <Link
                      className="page-link"
                      href="#"
                      onClick={() => {
                        return tokensFetchMore(
                          tokensData?.tokens?.pageInfo?.endCursor,
                          symbol,
                        );
                      }}
                      data-next-page-button=""
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width={6}
                        height={10}
                      >
                        <path
                          fillRule="evenodd"
                          d="M5.715 5.715c-.064.064-.141.102-.217.144L1.642 9.715A.959.959 0 1 1 .285 8.358L3.642 5 .285 1.642A.959.959 0 1 1 1.642.285L5.498 4.14c.075.043.153.081.217.145A.949.949 0 0 1 5.989 5a.949.949 0 0 1-.274.715z"
                        />
                      </svg>
                    </Link>
                  </li>
                )}
              </ul>
            </div>
          </div>
        </div>
      </section>
    </main>
  );
};
