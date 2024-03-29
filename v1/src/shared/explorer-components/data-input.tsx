import React, { useEffect, useState } from "react";
import OutsideClickHandler from "react-outside-click-handler";
import { CopyToClipboard } from "./copy-to-clipboard";

export const DataInput: React.FC<{ input: string }> = ({ input }) => {
  const [showMenu, setShowMenu] = useState(false);
  const [showUtf8, setShowUtf8] = useState(false);
  useEffect(() => {
    setShowMenu(false);
  }, [showUtf8]);
  useEffect(() => {
    if (showUtf8) {
      setCopy(hexToUtf8(input));
    } else {
      setCopy(input);
    }
  }, [showUtf8]);

  const [copy, setCopy] = useState(input);

  return (
    <>
      <div className="d-flex mb-3 justify-content-between">
        <div className={`dropdown${showMenu ? " show" : ""}`}>
          <button
            className="btn-dropdown-line large dropdown-toggle"
            type="button"
            id="tx-input-decoding-button"
            data-toggle="dropdown"
            aria-haspopup="true"
            aria-expanded="false"
            onClick={() => setShowMenu(!showMenu)}
          >
            {showUtf8 ? "UTF-8" : "Hex (Default)"}
          </button>
          <OutsideClickHandler onOutsideClick={() => setShowMenu(false)}>
            <div
              className={`dropdown-menu${showMenu ? " show" : ""}`}
              aria-labelledby="transaction-input-decoding-button"
            >
              <a
                href="#"
                onClick={(e) => {
                  e.preventDefault();
                  setShowUtf8(false);
                }}
                className="dropdown-item large tx-input-dropdown"
                data-target=".tx-raw-input"
                data-target-to-hide=".tx-utf8-input"
                id="tx-dropdown-raw"
              >
                Hex (Default)
              </a>
              <a
                href="#"
                onClick={(e) => {
                  e.preventDefault();
                  setShowUtf8(true);
                }}
                className="dropdown-item large tx-input-dropdown"
                data-target=".tx-utf8-input"
                data-target-to-hide=".tx-raw-input"
                id="tx-dropdown-utf8"
              >
                UTF-8
              </a>
            </div>
          </OutsideClickHandler>
        </div>
        <div className="btn-copy-tx-raw-input-container">
          <CopyToClipboard value={copy} reason="Copy Txn Input" />
        </div>
      </div>

      <div className="tx-raw-input" style={{ display: showUtf8 ? "none" : "" }}>
        <div className="tile tile-muted">
          <pre className="pre-scrollable pre-scrollable-shorty pre-wrap mb-0">
            <code>{input}</code>
          </pre>
        </div>
      </div>

      <div
        className="tx-utf8-input"
        style={{ display: showUtf8 ? "" : "none" }}
      >
        <div className="tile tile-muted">
          <pre className="pre-scrollable pre-scrollable-shorty pre-wrap mb-0">
            <code>{hexToUtf8(input.replace(/^0x/, ""))}</code>
          </pre>
        </div>
      </div>
    </>
  );
};

function hexToUtf8(hex: string): string {
  const result = [];
  for (let i = 0; i < hex.length; i += 2) {
    result.push(String.fromCharCode(parseInt(hex.substr(i, 2), 16)));
  }
  return result.join("");
}
