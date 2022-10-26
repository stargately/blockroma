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
    ["/bmo/testnet/", "BoomMo Chain"],
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
        {chainName}{" "}
        <svg
          width={12}
          fill="#BDC4D2"
          height={12}
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 320 512"
        >
          <path d="M137.4 374.6c12.5 12.5 32.8 12.5 45.3 0l128-128c9.2-9.2 11.9-22.9 6.9-34.9s-16.6-19.8-29.6-19.8L32 192c-12.9 0-24.6 7.8-29.6 19.8s-2.2 25.7 6.9 34.9l128 128z" />
        </svg>
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
