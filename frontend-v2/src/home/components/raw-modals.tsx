import * as React from "react";

export function RawModals(): JSX.Element {
  return (
    <>
      <div
        className="modal fade"
        id="errorStatusModal"
        tabIndex={-1}
        role="dialog"
        aria-hidden="true"
      >
        <div
          className="modal-dialog modal-dialog-centered modal-status"
          role="document"
        >
          <div className="modal-content">
            <div className="modal-status-graph modal-status-graph-error">
              <svg xmlns="http://www.w3.org/2000/svg" width={85} height={86}>
                <defs>
                  <filter
                    id="errora"
                    width={85}
                    height={86}
                    x={0}
                    y={0}
                    filterUnits="userSpaceOnUse"
                  >
                    <feOffset dy={6} in="SourceAlpha" />
                    <feGaussianBlur result="blurOut" stdDeviation="3.464" />
                    <feFlood floodColor="#C80A40" result="floodOut" />
                    <feComposite in="floodOut" in2="blurOut" operator="atop" />
                    <feComponentTransfer>
                      <feFuncA slope=".6" type="linear" />
                    </feComponentTransfer>
                    <feMerge>
                      <feMergeNode />
                      <feMergeNode in="SourceGraphic" />
                    </feMerge>
                  </filter>
                  <filter id="errorb">
                    <feOffset dy={-4} in="SourceAlpha" />
                    <feGaussianBlur result="blurOut" stdDeviation="2.828" />
                    <feFlood floodColor="#FF0D51" result="floodOut" />
                    <feComposite
                      in="floodOut"
                      in2="blurOut"
                      operator="out"
                      result="compOut"
                    />
                    <feComposite in="compOut" in2="SourceAlpha" operator="in" />
                    <feComponentTransfer>
                      <feFuncA slope=".5" type="linear" />
                    </feComponentTransfer>
                    <feBlend in2="SourceGraphic" />
                  </filter>
                </defs>
                <path
                  fill="#FFF"
                  fillRule="evenodd"
                  d="M54.738 36.969L70.342 52.58c3.521 3.524 3.521 9.237 0 12.76a9.015 9.015 0 0 1-12.754 0L41.984 49.729 26.381 65.34a9.015 9.015 0 0 1-12.754 0c-3.522-3.523-3.522-9.236 0-12.76l15.604-15.611-15.572-15.58c-3.522-3.524-3.522-9.237 0-12.76a9.013 9.013 0 0 1 12.753 0l15.572 15.58 15.572-15.58a9.015 9.015 0 0 1 12.754 0c3.522 3.523 3.522 9.236 0 12.76l-15.572 15.58z"
                  filter="url(#errorb)"
                />
              </svg>
            </div>
            <button
              type="button"
              className="close close-modal"
              data-dismiss="modal"
              aria-label="Close"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width={18} height={18}>
                <path
                  fillRule="evenodd"
                  d="M10.435 8.983l7.261 7.261a1.027 1.027 0 1 1-1.452 1.452l-7.261-7.261-7.262 7.262a1.025 1.025 0 1 1-1.449-1.45l7.262-7.261L.273 1.725A1.027 1.027 0 1 1 1.725.273l7.261 7.261 7.23-7.231a1.025 1.025 0 1 1 1.449 1.45l-7.23 7.23z"
                />
              </svg>
            </button>
            <div className="modal-body modal-status-body">
              <h2 className="modal-status-title" />
              <p
                className="modal-status-text"
                style={{ wordBreak: "break-all" }}
              />
              <div className="modal-status-button-wrapper">
                <button className="btn-line" type="button" data-dismiss="modal">
                  <span className="btn-line-text">Ok</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div
        className="modal fade"
        id="warningStatusModal"
        tabIndex={-1}
        role="dialog"
        aria-hidden="true"
      >
        <div
          className="modal-dialog modal-dialog-centered modal-status"
          role="document"
        >
          <div className="modal-content">
            <div className="modal-status-graph modal-status-graph-warning">
              <svg xmlns="http://www.w3.org/2000/svg" width={43} height={85}>
                <defs>
                  <filter
                    id="warninga"
                    width={43}
                    height={85}
                    x={0}
                    y={0}
                    filterUnits="userSpaceOnUse"
                  >
                    <feOffset dy={6} in="SourceAlpha" />
                    <feGaussianBlur result="blurOut" stdDeviation="3.464" />
                    <feFlood floodColor="#C16502" result="floodOut" />
                    <feComposite in="floodOut" in2="blurOut" operator="atop" />
                    <feComponentTransfer>
                      <feFuncA slope=".6" type="linear" />
                    </feComponentTransfer>
                    <feMerge>
                      <feMergeNode />
                      <feMergeNode in="SourceGraphic" />
                    </feMerge>
                  </filter>
                  <filter id="warningb">
                    <feOffset dy={-4} in="SourceAlpha" />
                    <feGaussianBlur result="blurOut" stdDeviation="2.828" />
                    <feFlood floodColor="#FF8502" result="floodOut" />
                    <feComposite
                      in="floodOut"
                      in2="blurOut"
                      operator="out"
                      result="compOut"
                    />
                    <feComposite in="compOut" in2="SourceAlpha" operator="in" />
                    <feComponentTransfer>
                      <feFuncA slope=".5" type="linear" />
                    </feComponentTransfer>
                    <feBlend in2="SourceGraphic" />
                  </filter>
                </defs>
                <path
                  fill="#FFF"
                  fillRule="evenodd"
                  d="M30.23 18.848L26 40h-.1a5.003 5.003 0 0 1-9.8 0H16l-4.23-21.152A9.958 9.958 0 0 1 11 15c0-5.523 4.477-10 10-10s10 4.477 10 10a9.958 9.958 0 0 1-.77 3.848zM21 49a9 9 0 0 1 9 9 9 9 0 0 1-9 9 9 9 0 0 1-9-9 9 9 0 0 1 9-9z"
                  filter="url(#warningb)"
                />
              </svg>
            </div>
            <button
              type="button"
              className="close close-modal"
              data-dismiss="modal"
              aria-label="Close"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width={18} height={18}>
                <path
                  fillRule="evenodd"
                  d="M10.435 8.983l7.261 7.261a1.027 1.027 0 1 1-1.452 1.452l-7.261-7.261-7.262 7.262a1.025 1.025 0 1 1-1.449-1.45l7.262-7.261L.273 1.725A1.027 1.027 0 1 1 1.725.273l7.261 7.261 7.23-7.231a1.025 1.025 0 1 1 1.449 1.45l-7.23 7.23z"
                />
              </svg>
            </button>
            <div className="modal-body modal-status-body">
              <h2 className="modal-status-title" />
              <p className="modal-status-text" style={{}} />
              <div className="modal-status-button-wrapper">
                <button className="btn-line" type="button" data-dismiss="modal">
                  <span className="btn-line-text">Ok</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div
        className="modal fade"
        id="successStatusModal"
        tabIndex={-1}
        role="dialog"
        aria-hidden="true"
      >
        <div
          className="modal-dialog modal-dialog-centered modal-status"
          role="document"
        >
          <div className="modal-content">
            <div className="modal-status-graph modal-status-graph-success">
              <svg xmlns="http://www.w3.org/2000/svg" width={99} height={80}>
                <defs>
                  <filter
                    id="successa"
                    width={99}
                    height={80}
                    x={0}
                    y={0}
                    filterUnits="userSpaceOnUse"
                  >
                    <feOffset dy={6} in="SourceAlpha" />
                    <feGaussianBlur result="blurOut" stdDeviation="3.464" />
                    <feFlood floodColor="#0E796B" result="floodOut" />
                    <feComposite in="floodOut" in2="blurOut" operator="atop" />
                    <feComponentTransfer>
                      <feFuncA slope=".6" type="linear" />
                    </feComponentTransfer>
                    <feMerge>
                      <feMergeNode />
                      <feMergeNode in="SourceGraphic" />
                    </feMerge>
                  </filter>
                  <filter id="successb">
                    <feOffset dy={-4} in="SourceAlpha" />
                    <feGaussianBlur result="blurOut" stdDeviation="2.828" />
                    <feFlood floodColor="#15BBA6" result="floodOut" />
                    <feComposite
                      in="floodOut"
                      in2="blurOut"
                      operator="out"
                      result="compOut"
                    />
                    <feComposite in="compOut" in2="SourceAlpha" operator="in" />
                    <feComponentTransfer>
                      <feFuncA slope=".5" type="linear" />
                    </feComponentTransfer>
                    <feBlend in2="SourceGraphic" />
                  </filter>
                </defs>
                <path
                  fill="#FFF"
                  fillRule="evenodd"
                  d="M84.31 21.389L46.381 59.34a8.993 8.993 0 0 1-6.399 2.642 8.998 8.998 0 0 1-6.402-2.642L13.629 39.389a9.022 9.022 0 0 1 0-12.76 9.022 9.022 0 0 1 12.76 0l13.593 13.593L71.557 8.629a9.013 9.013 0 0 1 12.753 0c3.522 3.523 3.522 9.236 0 12.76z"
                  filter="url(#successb)"
                />
              </svg>
            </div>
            <button
              type="button"
              className="close close-modal"
              data-dismiss="modal"
              aria-label="Close"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width={18} height={18}>
                <path
                  fillRule="evenodd"
                  d="M10.435 8.983l7.261 7.261a1.027 1.027 0 1 1-1.452 1.452l-7.261-7.261-7.262 7.262a1.025 1.025 0 1 1-1.449-1.45l7.262-7.261L.273 1.725A1.027 1.027 0 1 1 1.725.273l7.261 7.261 7.23-7.231a1.025 1.025 0 1 1 1.449 1.45l-7.23 7.23z"
                />
              </svg>
            </button>
            <div className="modal-body modal-status-body">
              <h2 className="modal-status-title" />
              <p className="modal-status-text" style={{}} />
              <div className="modal-status-button-wrapper">
                <button className="btn-line" type="button" data-dismiss="modal">
                  <span className="btn-line-text">Ok</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div
        className="modal fade"
        id="questionStatusModal"
        tabIndex={-1}
        role="dialog"
        aria-hidden="true"
      >
        <div
          className="modal-dialog modal-dialog-centered modal-status"
          role="document"
        >
          <div className="modal-content">
            <div className="modal-status-graph modal-status-graph-question">
              <svg xmlns="http://www.w3.org/2000/svg" width={62} height={85}>
                <defs>
                  <filter
                    id="questiona"
                    width={62}
                    height={85}
                    x={0}
                    y={0}
                    filterUnits="userSpaceOnUse"
                  >
                    <feOffset dy={6} in="SourceAlpha" />
                    <feGaussianBlur result="blurOut" stdDeviation="3.464" />
                    <feFlood floodColor="#2674AF" result="floodOut" />
                    <feComposite in="floodOut" in2="blurOut" operator="atop" />
                    <feComponentTransfer>
                      <feFuncA slope=".6" type="linear" />
                    </feComponentTransfer>
                    <feMerge>
                      <feMergeNode />
                      <feMergeNode in="SourceGraphic" />
                    </feMerge>
                  </filter>
                  <filter id="questionb">
                    <feOffset dy={-4} in="SourceAlpha" />
                    <feGaussianBlur result="blurOut" stdDeviation="2.828" />
                    <feFlood floodColor="#329AE9" result="floodOut" />
                    <feComposite
                      in="floodOut"
                      in2="blurOut"
                      operator="out"
                      result="compOut"
                    />
                    <feComposite in="compOut" in2="SourceAlpha" operator="in" />
                    <feComponentTransfer>
                      <feFuncA slope=".5" type="linear" />
                    </feComponentTransfer>
                    <feBlend in2="SourceGraphic" />
                  </filter>
                </defs>
                <path
                  fill="#FFF"
                  fillRule="evenodd"
                  d="M30.385 46.419c1.149 0 2.146-.344 2.995-1.033.847-.689 1.381-1.664 1.6-2.928.273-1.205.889-2.368 1.846-3.487.957-1.119 2.339-2.454 4.144-4.004 1.915-1.779 3.474-3.3 4.678-4.563 1.203-1.262 2.228-2.784 3.077-4.564.848-1.779 1.272-3.731 1.272-5.855 0-2.87-.794-5.438-2.38-7.707-1.587-2.267-3.774-4.046-6.565-5.338-2.79-1.292-5.963-1.937-9.518-1.937-3.174 0-6.333.531-9.478 1.593-3.147 1.062-5.95 2.54-8.412 4.434-.931.747-1.6 1.522-2.01 2.325-.41.804-.616 1.808-.616 3.013 0 1.78.479 3.287 1.436 4.521.957 1.235 2.12 1.851 3.488 1.851 1.149 0 2.653-.545 4.513-1.636l1.97-1.033c1.531-.918 2.913-1.621 4.144-2.11a10.075 10.075 0 0 1 3.734-.732c1.53 0 2.734.374 3.61 1.12.875.747 1.313 1.752 1.313 3.014 0 1.263-.315 2.41-.943 3.444-.63 1.033-1.574 2.268-2.832 3.702-1.751 1.895-3.118 3.746-4.103 5.554-.984 1.808-1.477 4.062-1.477 6.759 0 1.78.397 3.158 1.19 4.133.793.977 1.9 1.464 3.324 1.464zm.164 20.579c2.297 0 4.198-.803 5.703-2.411 1.504-1.606 2.257-3.587 2.257-5.941 0-2.353-.753-4.334-2.257-5.941-1.505-1.607-3.406-2.411-5.703-2.411-2.244 0-4.117.804-5.621 2.411-1.506 1.607-2.257 3.588-2.257 5.941 0 2.354.751 4.335 2.257 5.941 1.504 1.608 3.377 2.411 5.621 2.411z"
                  filter="url(#questionb)"
                />
              </svg>
            </div>
            <button
              type="button"
              className="close close-modal"
              data-dismiss="modal"
              aria-label="Close"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width={18} height={18}>
                <path
                  fillRule="evenodd"
                  d="M10.435 8.983l7.261 7.261a1.027 1.027 0 1 1-1.452 1.452l-7.261-7.261-7.262 7.262a1.025 1.025 0 1 1-1.449-1.45l7.262-7.261L.273 1.725A1.027 1.027 0 1 1 1.725.273l7.261 7.261 7.23-7.231a1.025 1.025 0 1 1 1.449 1.45l-7.23 7.23z"
                />
              </svg>
            </button>
            <div className="modal-body modal-status-body">
              <h2 className="modal-status-title" />
              <p className="modal-status-text" style={{}} />
              <div className="modal-status-button-wrapper">
                <button className="btn-line except" type="button">
                  <span className="btn-line-text">No</span>
                </button>
                <button className="btn-line accept" type="button">
                  <span className="btn-line-text">Yes</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
