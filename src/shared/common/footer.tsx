import React, { useState } from "react";
import { assetURL } from "onefx/lib/asset-url";
import { useChainConfig } from "@/shared/common/use-chain-config";
import { t } from "onefx/lib/iso-i18n";
import { TOP_BAR_HEIGHT } from "./top-bar";

const addChainToMM = require("../blockscout-web-js-lib/add_chain_to_mm");

export const FOOTER_HEIGHT = 0;

export const FOOTER_ABOVE = {
  minHeight: `calc(100vh - ${FOOTER_HEIGHT + TOP_BAR_HEIGHT}px)`,
};

export function Footer(): JSX.Element {
  const [mmAdded, setMmAdded] = useState(false);
  const chainConfig = useChainConfig();
  return (
    <footer className="footer">
      <div className="footer-body container">
        <div className="row footer-logo-row">
          <div className="col-md-12">
            <a className="footer-brand" href={assetURL("")}>
              <img
                className="footer-logo"
                src={assetURL("favicon.png")}
                alt="BMO"
              />
            </a>
          </div>
        </div>
        <div className="row">
          <div className="col-xs-12 col-lg-3">
            <p className="footer-info-text">{t("meta.description")}</p>
            <div className="footer-social-icons">
              <a
                href="https://github.com/stargately/blockroma"
                rel="noreferrer"
                target="_blank"
                className="footer-social-icon"
                title="Github"
              >
                <div className="footer-social-icon-container fontawesome-icon github" />
              </a>
              <a
                href="https://www.twitter.com/puncsky/"
                rel="noreferrer"
                target="_blank"
                className="footer-social-icon"
                title="Twitter"
              >
                <div className="footer-social-icon-container fontawesome-icon twitter" />
              </a>
              <a
                href="https://t.me/system_design_and_architecture"
                rel="noreferrer"
                target="_blank"
                className="footer-social-icon"
                title="Telegram"
              >
                <div className="footer-social-icon-container fontawesome-icon telegram" />
              </a>
            </div>
          </div>
          <div className="col-xs-12 col-md-4 col-lg-3 footer-list">
            <h3>Blockroma</h3>
            <ul>
              <li>
                <a
                  href="https://github.com/stargately/blockroma/issues/new?labels=Blockroma&body=%2ADescribe+your+issue+here.%2A%0A%0A%23%23%23+Environment%0A%2A+Elixir+Version%3A+1.12.3%0A%2A+Erlang+Version%3A+24%0A%2A+Blockroma+Version%3A+v4.1.2-beta%0A%0A%2A+User+Agent%3A+%60Mozilla%2F5.0+%28Macintosh%3B+Intel+Mac+OS+X+10_15_7%29+AppleWebKit%2F537.36+%28KHTML%2C+like+Gecko%29+Chrome%2F99.0.4844.84+Safari%2F537.36%60%0A%0A%23%23%23+Steps+to+reproduce%0A%0A%2ATell+us+how+to+reproduce+this+issue.+If+possible%2C+push+up+a+branch+to+your+fork+with+a+regression+test+we+can+run+to+reproduce+locally.%2A%0A%0A%23%23%23+Expected+Behaviour%0A%0A%2ATell+us+what+should+happen.%2A%0A%0A%23%23%23+Actual+Behaviour%0A%0A%2ATell+us+what+happens+instead.%2A%0A&title=BMO%3A+%3CIssue+Title%3E"
                  rel="noreferrer"
                  className="footer-link"
                  target="_blank"
                >
                  Submit an Issue
                </a>
              </li>
              <li>
                <a
                  href="https://github.com/stargately/blockroma"
                  rel="noreferrer"
                  className="footer-link"
                >
                  Contribute
                </a>
              </li>
              <li>
                <a
                  href="https://discord.gg/Pb5YbK3ykN"
                  rel="noreferrer"
                  className="footer-link"
                >
                  Chat (#blockroma)
                </a>
              </li>
              <li>
                <a
                  onClick={async () => {
                    await setMmAdded(await addChainToMM(chainConfig));
                  }}
                  className="footer-link js-btn-add-chain-to-mm btn-add-chain-to-mm in-footer"
                  style={{ cursor: "pointer" }}
                >
                  {mmAdded ? "BMO Added" : "Add BMO"}
                </a>
              </li>
            </ul>
          </div>
          <div className="col-xs-12 col-md-4 col-lg-3 footer-list">
            <h3>Main Networks</h3>
            <ul>
              <li>
                <a
                  href="https://blockroma.com/eth/mainnet/"
                  rel="norefferer"
                  className="footer-link"
                >
                  {" "}
                  Ethereum{" "}
                </a>
              </li>
              <li>
                <a
                  href="https://blockroma.com/etc/mainnet"
                  rel="norefferer"
                  className="footer-link"
                >
                  {" "}
                  Ethereum Classic{" "}
                </a>
              </li>
            </ul>
          </div>
          <div className="col-xs-12 col-md-4 col-lg-3 footer-list">
            <h3>Test Networks</h3>
            <ul>
              <li>
                <a
                  href="https://blockroma.com/poa/sokol"
                  rel="noreferrer"
                  className="footer-link"
                >
                  {" "}
                  Sokol{" "}
                </a>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </footer>
  );
}
