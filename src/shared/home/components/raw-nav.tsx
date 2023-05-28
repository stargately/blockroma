import React, { FormEventHandler, useEffect, useRef, useState } from "react";
import { assetURL } from "onefx/lib/asset-url";
import OutsideClickHandler from "react-outside-click-handler";
import { useChainConfig } from "@/shared/common/use-chain-config";
import { CommonMargin } from "@/shared/common/common-margin";
import { t } from "onefx/lib/iso-i18n";
import { MultiChainDropdown } from "@/shared/home/components/multi-chain-dropdown";
import { useParams } from "react-router";
import { DarkModeToggle } from "./dark-mode-toggle";

function DesktopSearch({ searchVal }: { searchVal: string }): JSX.Element {
  const desktopTextInput = useRef<HTMLInputElement>(null);
  const [focusedField, setFocusedField] = useState(false);
  const [dangerField, setDangerField] = useState(false);
  useEffect(() => {
    desktopTextInput?.current?.addEventListener("focus", () => {
      setFocusedField(true);
      setDangerField(false);
    });
    desktopTextInput?.current?.addEventListener("blur", () => {
      setFocusedField(false);
      setDangerField(false);
    });
    document.addEventListener("keyup", (event) => {
      if (event.key === "/") {
        desktopTextInput?.current?.focus();
      }
    });
    if (desktopTextInput.current) {
      desktopTextInput.current.value = searchVal;
    }
  }, [desktopTextInput, searchVal]);

  const onSubmit: FormEventHandler<HTMLFormElement> = (event) => {
    setDangerField(false);
    event.preventDefault();
    // @ts-ignore
    const val = String(event?.target[0].value).trim();
    if (val.length === 66) {
      window.location.replace(assetURL(`tx/${val}`));
    } else if (val.length === 42) {
      window.location.replace(assetURL(`address/${val}`));
    } else if (parseInt(val, 10) > 0) {
      window.location.replace(assetURL(`block/${val}`));
    } else {
      window.location.replace(assetURL(`tokens/?symbol=${val}`));
    }
  };

  let fieldLabel = focusedField ? " focused-field" : "";
  if (dangerField) {
    fieldLabel = " danger-field";
  }

  return (
    <form
      onSubmit={onSubmit}
      className="search-form d-lg-flex d-inline-block mobile-search-hide"
    >
      <div
        className="input-group"
        style={{ width: "100%" }}
        title={t("nav.search")}
      >
        <div
          className={`form-control search-control me auto mobile-search-hide${fieldLabel}`}
        >
          <div
            className="autoComplete_wrapper"
            role="combobox"
            aria-owns="autoComplete_list_1"
            aria-haspopup="true"
            aria-expanded="false"
          >
            <input
              ref={desktopTextInput}
              id="main-search-autocomplete"
              className="main-search-autocomplete"
              data-test="search_input"
              data-chain-id={99}
              data-display-token-icons
              type="text"
              tabIndex={1}
              aria-controls="autoComplete_list_1"
              aria-autocomplete="both"
              placeholder={t("nav.search")}
            />
            <ul id="autoComplete_list_1" role="listbox" hidden />
          </div>
        </div>
        <div className="input-group-append left">
          <button className="input-group-text" id="search-icon">
            <svg
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
          </button>
        </div>
        <div className="input-group-append right desktop-only">
          <div
            id="slash-icon"
            className="input-group-text border"
            data-placement="bottom"
            data-toggle="tooltip"
            data-original-title="Press / and focus will be moved to the search field"
            data-template="<div class='tooltip tooltip-pale-color' role='tooltip'><div class='arrow'></div><div class='tooltip-inner'></div></div>"
            style={focusedField ? { display: "none" } : {}}
          >
            /
          </div>
        </div>
      </div>
      <button
        className="btn btn-outline-success my-2 my-sm-0 sr-only hidden"
        type="submit"
      >
        Search
      </button>
    </form>
  );
}

function MobileSearch({ searchVal }: { searchVal: string }): JSX.Element {
  const textInput = useRef<HTMLInputElement>(null);
  const [focusedField, setFocusedField] = useState(false);
  const [dangerField, setDangerField] = useState(false);
  useEffect(() => {
    textInput?.current?.addEventListener("focus", () => {
      setFocusedField(true);
      setDangerField(false);
    });
    textInput?.current?.addEventListener("blur", () => {
      setFocusedField(false);
      setDangerField(false);
    });
    if (textInput.current) {
      textInput.current.value = searchVal;
    }
  }, [textInput, searchVal]);

  const onSubmit: FormEventHandler<HTMLFormElement> = (event) => {
    setDangerField(false);
    event.preventDefault();
    // @ts-ignore
    const val = event?.target[0].value;
    if (val.length === 66) {
      window.location.replace(assetURL(`tx/${val}`));
    } else if (val.length === 42) {
      window.location.replace(assetURL(`address/${val}`));
    } else if (parseInt(val, 10) > 0) {
      window.location.replace(assetURL(`block/${val}`));
    } else {
      setDangerField(true);
    }
  };

  let fieldLabel = focusedField ? " focused-field" : "";
  if (dangerField) {
    fieldLabel = " danger-field";
  }

  return (
    <form
      onSubmit={onSubmit}
      className="search-form d-lg-flex d-inline-block mobile-search-show"
    >
      <div
        className="input-group"
        style={{ width: "100%" }}
        title={t("nav.search")}
      >
        <div
          className={`form-control search-control me auto mobile-search-show${fieldLabel}`}
        >
          <div
            className="autoComplete_wrapper"
            role="combobox"
            aria-owns="autoComplete_list_2"
            aria-haspopup="true"
            aria-expanded="false"
          >
            <input
              ref={textInput}
              id="main-search-autocomplete-mobile"
              className="main-search-autocomplete"
              data-test="search_input"
              data-chain-id={99}
              data-display-token-icons
              type="text"
              tabIndex={1}
              aria-controls="autoComplete_list_2"
              aria-autocomplete="both"
              placeholder={t("nav.search")}
            />
            <ul id="autoComplete_list_2" role="listbox" hidden />
          </div>
        </div>
        <div className="input-group-append left">
          <button className="input-group-text" id="search-icon">
            <svg
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
          </button>
        </div>
        <div className="input-group-append right desktop-only">
          <div
            id="slash-icon"
            className="input-group-text border"
            data-placement="bottom"
            data-toggle="tooltip"
            data-original-title="Press / and focus will be moved to the search field"
            data-template="<div class='tooltip tooltip-pale-color' role='tooltip'><div class='arrow'></div><div class='tooltip-inner'></div></div>"
          >
            /
          </div>
        </div>
      </div>
      <button
        className="btn btn-outline-success my-2 my-sm-0 sr-only hidden"
        type="submit"
      >
        Search
      </button>
    </form>
  );
}

export const RawNav: React.FC = () => {
  const [barClapsed, setBarClapsed] = useState(true);
  const chainConfig = useChainConfig();
  const { blockNumber, addressHash, txHash } = useParams<{
    blockNumber?: string;
    addressHash?: string;
    txHash?: string;
  }>();
  const searchVal = blockNumber ?? addressHash ?? txHash ?? "";

  return (
    <OutsideClickHandler onOutsideClick={() => setBarClapsed(true)}>
      <nav
        className="navbar navbar-expand-lg navbar-primary"
        data-selector="navbar"
        id="top-navbar"
      >
        <div className="container-fluid navbar-container">
          <a
            className="navbar-brand"
            data-test="header_logo"
            href={assetURL("")}
          >
            <img
              className="navbar-logo"
              id="navbar-logo"
              src={assetURL("favicon.svg")}
              alt={chainConfig.symbol}
            />
            <CommonMargin />
            {chainConfig.chainName}
          </a>
          <button
            onClick={() => setBarClapsed(!barClapsed)}
            className="navbar-toggler"
            type="button"
            data-toggle="collapse"
            data-target="#navbarSupportedContent"
            aria-controls="navbarSupportedContent"
            aria-expanded="false"
            aria-label="Toggle navigation"
          >
            <span className="navbar-toggler-icon" />
          </button>
          <div
            className={`collapse navbar-collapse${barClapsed ? "" : " show"}`}
            id="navbarSupportedContent"
          >
            <ul className="navbar-nav">
              <li className="nav-item dropdown">
                <a
                  className="nav-link topnav-nav-link"
                  href={assetURL("blocks")}
                  id="navbarBlocksDropdown"
                  role="button"
                  data-toggle="dropdown"
                  aria-haspopup="true"
                  aria-expanded="false"
                >
                  <span className="nav-link-icon">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      width={14}
                      height={16}
                    >
                      <path
                        fill="#D2DAE9"
                        fillRule="evenodd"
                        d="M7 0L-.007 3.498 7 6.996l7.007-3.498L7 0z"
                      />
                      <path
                        fill="#C1C7D4"
                        fillRule="evenodd"
                        d="M0 5.005v7.674L6.004 16V8.326L0 5.005z"
                      />
                      <path
                        fill="#A3A9B5"
                        fillRule="evenodd"
                        d="M7.996 8.326V16L14 12.679V5.005L7.996 8.326z"
                      />
                    </svg>
                  </span>
                  {t("nav.blocks")}
                </a>
              </li>
              <li className="nav-item dropdown" id="activeTransactions">
                <a
                  href={assetURL("txs")}
                  role="button"
                  id="navbarTransactionsDropdown"
                  className="nav-link topnav-nav-link"
                  data-toggle="dropdown"
                  aria-haspopup="true"
                  aria-expanded="false"
                >
                  <span className="nav-link-icon">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      width={16}
                      height={14}
                    >
                      <path
                        fill="#BDC4D2"
                        fillRule="evenodd"
                        d="M15 5h-4.519l1.224 1.224a1.037 1.037 0 0 1-1.465 1.465L7.311 4.76a1.017 1.017 0 0 1-.24-.409 1.024 1.024 0 0 1-.066-.327C7.005 4.016 7 4.009 7 4c-.012-.28.082-.562.295-.776L10.224.295a1.037 1.037 0 0 1 1.465 1.465L10.448 3H15a1 1 0 0 1 0 2z"
                      />
                      <path
                        fill="#A3A9B5"
                        fillRule="evenodd"
                        d="M8.788 10.587c-.038.058-.063.122-.114.173l-2.925 2.929a1.033 1.033 0 0 1-1.462 0 1.037 1.037 0 0 1 0-1.465L5.509 11H1a1 1 0 0 1 0-2h4.541L4.303 7.76a1.037 1.037 0 0 1 0-1.465 1.033 1.033 0 0 1 1.462 0L8.69 9.224c.203.203.303.47.302.736.001.014.008.026.008.04 0 .224-.087.42-.212.587z"
                      />
                    </svg>
                  </span>
                  {t("nav.txs")}
                </a>
              </li>
              <li className="nav-item dropdown">
                <a
                  href={assetURL("tokens/")}
                  role="button"
                  id="navbarAPIsDropdown"
                  className="nav-link topnav-nav-link"
                  data-toggle="dropdown"
                  aria-haspopup="true"
                  aria-expanded="false"
                >
                  <span className="nav-link-icon">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      viewBox="0 0 24 24"
                      enableBackground="new 0 0 24 24"
                      width={16}
                      height={16}
                    >
                      <path
                        fill="#A3A9B5"
                        fillRule="evenodd"
                        d="M 16.5 0 C 12.357864 0 9 1.790861 9 4 C 9 6.209139 12.357864 8 16.5 8 C 20.642136 8 24 6.209139 24 4 C 24 1.790861 20.642136 0 16.5 0 z M 9.0625 5.5 C 9.0235 5.664 9 5.83 9 6 L 9 7 C 9 7.849 9.49775 8.63325 10.34375 9.28125 C 11.80575 9.57825 13.06825 10.10225 14.03125 10.78125 C 14.80225 10.92425 15.638 11 16.5 11 C 20.642 11 24 9.209 24 7 L 24 6 C 24 5.83 23.9765 5.664 23.9375 5.5 C 23.4755 7.473 20.324 9 16.5 9 C 12.676 9 9.5245 7.473 9.0625 5.5 z M 23.9375 8.5 C 23.4755 10.473 20.324 12 16.5 12 C 16.079 12 15.68325 11.9735 15.28125 11.9375 C 15.74625 12.5585 15.997 13.266 16 14 C 16.166 14.005 16.331 14 16.5 14 C 20.642 14 24 12.209 24 10 L 24 9 C 24 8.83 23.9765 8.664 23.9375 8.5 z M 7.5 10 C 3.358 10 0 11.791 0 14 C 0 16.209 3.358 18 7.5 18 C 11.642 18 15 16.209 15 14 C 15 11.791 11.642 10 7.5 10 z M 23.9375 11.5 C 23.4755 13.473 20.324 15 16.5 15 C 16.331 15 16.166 15.005 16 15 L 16 17 C 16.166 17.005 16.331 17 16.5 17 C 20.642 17 24 15.209 24 13 L 24 12 C 24 11.83 23.9765 11.664 23.9375 11.5 z M 23.9375 14.5 C 23.4755 16.473 20.324 18 16.5 18 C 16.331 18 16.166 18.005 16 18 L 16 20 C 16.166 20.005 16.331 20 16.5 20 C 20.642 20 24 18.209 24 16 L 24 15 C 24 14.83 23.9765 14.664 23.9375 14.5 z M 0.0625 15.5 C 0.0235 15.664 0 15.83 0 16 L 0 17 C 0 19.209 3.358 21 7.5 21 C 11.642 21 15 19.209 15 17 L 15 16 C 15 15.83 14.9765 15.664 14.9375 15.5 C 14.4755 17.473 11.324 19 7.5 19 C 3.676 19 0.5245 17.473 0.0625 15.5 z M 0.0625 18.5 C 0.0235 18.664 0 18.83 0 19 L 0 20 C 0 22.209 3.358 24 7.5 24 C 11.642 24 15 22.209 15 20 L 15 19 C 15 18.83 14.9765 18.664 14.9375 18.5 C 14.4755 20.473 11.324 22 7.5 22 C 3.676 22 0.5245 20.473 0.0625 18.5 z"
                      />
                    </svg>
                  </span>
                  Tokens
                </a>
              </li>
              <li className="nav-item dropdown">
                <a
                  href={assetURL("api-gateway/")}
                  role="button"
                  id="navbarAPIsDropdown"
                  className="nav-link topnav-nav-link"
                  data-toggle="dropdown"
                  aria-haspopup="true"
                  aria-expanded="false"
                >
                  <span className="nav-link-icon">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      width={16}
                      height={16}
                    >
                      <path
                        fill="#D2DAE9"
                        fillRule="evenodd"
                        d="M15 3H5a1 1 0 0 1-2 0H1a1 1 0 0 1 0-2h2a1 1 0 0 1 2 0h10a1 1 0 0 1 0 2z"
                      />
                      <path
                        fill="#BDC4D2"
                        fillRule="evenodd"
                        d="M15 9h-2a1 1 0 0 1-2 0H1a1 1 0 0 1 0-2h10a1 1 0 0 1 2 0h2a1 1 0 0 1 0 2z"
                      />
                      <path
                        fill="#A3A9B5"
                        fillRule="evenodd"
                        d="M15 15H5a1 1 0 0 1-2 0H1a1 1 0 0 1 0-2h2a1 1 0 0 1 2 0h10a1 1 0 0 1 0 2z"
                      />
                    </svg>
                  </span>
                  API
                </a>
              </li>

              <MultiChainDropdown chainName={chainConfig.chainName} />
            </ul>
            <DarkModeToggle />
            <DesktopSearch searchVal={searchVal} />
          </div>
        </div>
        <MobileSearch searchVal={searchVal} />
      </nav>
    </OutsideClickHandler>
  );
};
