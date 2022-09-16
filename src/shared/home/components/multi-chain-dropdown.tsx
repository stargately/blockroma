import React, { useState } from "react";
import OutsideClickHandler from "react-outside-click-handler";

type Props = {
  chainName: string;
};

export const chainSwitchOpts = {
  mainnets: [
    ["/eth/mainnet/", "Ethereum"],
    ["/ethw/mainnet/", "ETHW-mainnet"],
  ],
  testnets: [
    ["/", "BoomMo Chain"],
    ["/ethw/iceberg/", "ETHW-iceberg-testnet"],
  ],
};

export function MultiChainDropdown({ chainName }: Props) {
  const [showDropdown, setShowDropdown] = useState(false);

  return (
    <li className="nav-item dropdown">
      <a
        className="nav-link topnav-nav-link active-icon js-show-network-selector"
        href="#"
        id="navbarDropdown"
        role="button"
        data-toggle="dropdown"
        aria-haspopup="true"
        aria-expanded="false"
        onClick={() => setShowDropdown(true)}
      >
        <span className="nav-link-icon">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width={6}
            height={6}
            className="active-dot-icon"
          >
            <circle cx={3} cy={3} r={3} fill="#80d6a1" />
          </svg>
        </span>
        {chainName}
      </a>

      <OutsideClickHandler onOutsideClick={() => setShowDropdown(false)}>
        <div
          className={`dropdown-menu dropdown-menu-right${
            showDropdown ? " show" : ""
          }`}
          aria-labelledby="navbarDropdown"
        >
          <a className="dropdown-item header division">Mainnets</a>

          {chainSwitchOpts.mainnets.map((it) => (
            <a className="dropdown-item" key={it[0]} href={it[0]}>
              {it[1]}
            </a>
          ))}

          <a className="dropdown-item header division">Testnets</a>
          {chainSwitchOpts.testnets.map((it) => (
            <a className="dropdown-item" key={it[0]} href={it[0]}>
              {it[1]}
            </a>
          ))}
        </div>
      </OutsideClickHandler>
    </li>
  );
}
