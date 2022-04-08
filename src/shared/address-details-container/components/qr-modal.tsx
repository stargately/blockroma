import * as React from "react";
import { useEffect } from "react";
import OutsideClickHandler from "react-outside-click-handler";
import { assetURL } from "onefx/lib/asset-url";

type Props = {
  addressHash: string;
  open: boolean;
  onClose: () => void;
  hashQr?: string | null;
};

export function QrModal({
  addressHash,
  open,
  onClose,
  hashQr,
}: Props): JSX.Element {
  useEffect(() => {
    const cl = window.document.getElementsByTagName("body")[0].classList;
    if (open) {
      cl.add("modal-open");
    } else {
      cl.remove("modal-open");
    }

    window.addEventListener("keyup", (event) => {
      if (event.key === "Escape") {
        onClose();
      }
    });
  }, [open]);
  return (
    <div
      className={`modal fade${open ? " show" : ""}`}
      id="qrModal"
      tabIndex={-1}
      role="dialog"
      style={{ display: open ? "block" : "none" }}
      aria-labelledby="qrModalLabel"
      aria-hidden={open ? "false" : "true"}
      aria-modal={open ? "true" : "false"}
    >
      <div className="modal-dialog modal-sm" role="document">
        <OutsideClickHandler onOutsideClick={onClose}>
          <div className="modal-content">
            <div className="modal-header">
              <h2 className="modal-title" id="qrModalLabel">
                QR Code
              </h2>
              <button
                type="button"
                className="close"
                data-dismiss="modal"
                aria-label="Close"
                onClick={onClose}
              >
                <span aria-hidden="true">×</span>
              </button>
            </div>
            <div className="modal-body">
              <img
                src={
                  hashQr || assetURL("images/errors-img/etc-tx-not-found.png")
                }
                className="qr-code"
                alt="qr_code"
                title={addressHash}
              />
            </div>
            <div className="modal-footer">
              <button
                onClick={onClose}
                type="button"
                className="btn btn-primary"
                data-dismiss="modal"
              >
                Close
              </button>
            </div>
          </div>
        </OutsideClickHandler>
      </div>
    </div>
  );
}
