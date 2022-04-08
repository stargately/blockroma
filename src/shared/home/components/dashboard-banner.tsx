import React from "react";

export const DashboardBanner: React.FC = () => {
  return (
    <div className="dashboard-banner-container" data-page="chain-details">
      <div className="container">
        <div className="dashboard-banner">
          <div className="dashboard-banner-network-graph">
            <div className="dashboard-banner-chart">
              <button
                data-chart-error-message
                className="alert alert-danger col-12 text-left mt-5"
                style={{ display: "none" }}
              >
                <span>There was a problem loading the chart.</span>
              </button>
              <canvas
                data-chart="historyChart"
                data-history_chart_paths='{
                        "market":"/market-history-chart"

                        ,


                        "transaction":"/transaction-history-chart"

}'
                data-history_chart_config='{"market":["price","market_cap"],"transactions":["transactions_per_day"]}'
                width={1034}
                height={296}
                style={{
                  display: "block",
                  boxSizing: "border-box",
                  height: "148px",
                  width: "517px",
                }}
              ></canvas>
            </div>
            <div className="dashboard-banner-chart-legend">
              <div className="dashboard-banner-chart-legend-item">
                <span className="dashboard-banner-chart-legend-label">
                  Gas tracker
                </span>
                <div className="dashboard-banner-chart-legend-value-container">
                  <span className="dashboard-banner-chart-legend-value inline">
                    <div>
                      <div
                        className="d-flex flex-row"
                        style={{ height: "20px", lineHeight: "20px" }}
                      >
                        <div className="flex-column">39.0 Gwei</div>
                        <span
                          data-toggle="tooltip"
                          data-placement="top"
                          data-html="true"
                          data-template="<div class='tooltip tooltip-inversed-color tooltip-gas-tracker' role='tooltip'><div class='arrow'></div><div class='tooltip-inner'></div></div>"
                          title="
                        <div class='custom-tooltip header'>Gas tracker</div>
                        <div>
                          <div class='custom-tooltip description left d-flex'><span>Slow</span><span class='custom-tooltip description right'>39.0 Gwei</span></div>
                          <div class='custom-tooltip description left d-flex'><span>Average</span><span class='custom-tooltip description right'>39.0 Gwei</span></div>
                          <div class='custom-tooltip description left d-flex'><span>Fast</span><span class='custom-tooltip description right'>39.0 Gwei</span></div>
                        </div>
                        "
                        >
                          <span
                            style={{
                              display: "inline-block",
                              height: "20px",
                              width: "12px",
                            }}
                            className="fontawesome-icon info-circle ml-1"
                          />
                        </span>
                      </div>
                    </div>
                  </span>
                </div>
              </div>
              <div className="dashboard-banner-chart-legend-item">
                <span className="dashboard-banner-chart-legend-label">
                  Daily Transactions
                </span>
                <span
                  className="dashboard-banner-chart-legend-value"
                  data-selector="tx_per_day"
                >
                  87,580
                  <span
                    data-toggle="tooltip"
                    data-placement="top"
                    data-html="true"
                    data-template="<div class='tooltip tooltip-inversed-color tooltip-gas-usage' role='tooltip'><div class='arrow'></div><div class='tooltip-inner'></div></div>"
                    title="<div class='custom-tooltip-header'>Gas used</div><div class='custom-tooltip-description'><b>9,693,816,212<b></div>"
                  >
                    <i
                      style={{ color: "#ffffff" }}
                      className="fa fa-info-circle ml-1"
                    />
                  </span>
                </span>
              </div>
            </div>
          </div>
          <div className="dashboard-banner-network-plain-container">
            <div className="dashboard-banner-network-stats">
              <div className="dashboard-banner-network-stats-item dashboard-banner-network-stats-item-1">
                <span className="dashboard-banner-network-stats-label">
                  Average block time
                </span>
                <span
                  className="dashboard-banner-network-stats-value"
                  data-selector="average-block-time"
                >
                  5 seconds
                </span>
              </div>
              <div className="dashboard-banner-network-stats-item dashboard-banner-network-stats-item-2">
                <span className="dashboard-banner-network-stats-label">
                  Total transactions
                </span>
                <div className="d-flex">
                  <span
                    className="dashboard-banner-network-stats-value"
                    data-selector="transaction-count"
                  >
                    53,185,348
                  </span>
                </div>
              </div>
              <div className="dashboard-banner-network-stats-item dashboard-banner-network-stats-item-3">
                <span className="dashboard-banner-network-stats-label">
                  Total blocks
                </span>
                <span
                  className="dashboard-banner-network-stats-value"
                  data-selector="block-count"
                >
                  26,556,683
                </span>
              </div>
              <div className="dashboard-banner-network-stats-item dashboard-banner-network-stats-item-4">
                <span className="dashboard-banner-network-stats-label">
                  Wallet addresses
                </span>
                <span
                  className="dashboard-banner-network-stats-value"
                  data-selector="address-count"
                >
                  19,438,885
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
