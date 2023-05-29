import React, { useRef } from "react";
import classnames from "classnames";
import { SubscribeForm } from "@site/src/components/subscribe-form/subscribe-form";

export const SiteFooter: React.FC = () => {
  const opts = {
    ctaText: "Start visualizing and indexing your blockchain",
    ctaButtonText: "Contact Sales",
    ctaButtonLink: "#",

    navigateTitle: "Navigate",
    navigateItems: [
      ["Stargately Home", "https://stargately.com/"],
      ["Beancount.io", "https://beancount.io/"],
      ["10x.pub", "https://10x.pub/"],
    ],

    contactTitle: "Contact",
    contactAddress: ["601 Van Ness Avenue", "San Francisco, CA"],
    contactPhone: [
      "Mike Thrift",
      "https://t.me/mikethrift",
      "Singapore Time (SGT) Mon - Fri 7am - 3pm",
    ],
    contactEmail: "dev@stargately.com",

    social: [
      // {
      //   name: "instagram",
      //   url: "#",
      // },
      // {
      //   name: "twitter",
      //   url: "#",
      // },
      // {
      //   name: "youtube",
      //   url: "#",
      // },
      // {
      //   name: "medium",
      //   url: "#",
      // },
      // {
      //   name: "facebook",
      //   url: "#",
      // },
    ],
    tradeMark: "@ " + new Date().getFullYear() + " Stargately",
  };

  return (
    <footer className="pb-4 bg-primary-3 text-light" id="footer">
      <div className="container">
        <div className="row mb-5">
          <div className="col">
            <div className="card card-body border-0 o-hidden mb-0 bg-primary text-light">
              <div className="position-relative d-flex flex-column flex-md-row justify-content-between align-items-center">
                <div className="h3 text-center mb-md-0">{opts.ctaText}</div>
                <a href={opts.ctaButtonText} className="btn btn-lg btn-white">
                  {opts.ctaButtonText}
                </a>
              </div>
              <div className="decoration layer-0">
                <img
                  className="bg-primary-2"
                  src="assets/img/decorations/deco-blob-1.svg"
                  alt="deco-blob-1 decoration"
                  data-inject-svg=""
                />
              </div>
            </div>
          </div>
        </div>
        <div className="row mb-5">
          <div className="col-6 col-lg-3 col-xl-2">
            <h5>{opts.navigateTitle}</h5>
            <ul className="nav flex-column">
              {opts.navigateItems.map((it) => {
                return (
                  <li className="nav-item" key={it[0]}>
                    <a href={it[1]} className="nav-link">
                      {it[0]}
                    </a>
                  </li>
                );
              })}
            </ul>
          </div>
          <div className="col-6 col-lg">
            <h5>{opts.contactTitle}</h5>
            <ul className="list-unstyled">
              <li className="mb-3 d-flex">
                <img
                  className="icon"
                  src="assets/img/icons/theme/map/marker-1.svg"
                  alt="marker-1 icon"
                  data-inject-svg=""
                />
                <div className="ml-3">
                  <span>
                    {opts.contactAddress.map((it) => (
                      <span key={it}>
                        {it}
                        <br />
                      </span>
                    ))}
                  </span>
                </div>
              </li>
              <li className="mb-3 d-flex">
                <img
                  className="icon"
                  src="assets/img/icons/theme/communication/call-1.svg"
                  alt="call-1 icon"
                  data-inject-svg=""
                />
                <div className="ml-3">
                  <a
                    href={opts.contactPhone[1]}
                    rel="noreferrer nofollow noopener"
                  >
                    {opts.contactPhone[0]}
                  </a>
                  <span className="d-block text-muted text-small">
                    {opts.contactPhone[2]}
                  </span>
                </div>
              </li>
              <li className="mb-3 d-flex">
                <img
                  className="icon"
                  src="assets/img/icons/theme/communication/mail.svg"
                  alt="mail icon"
                  data-inject-svg=""
                />
                <div className="ml-3">
                  <a href="#">{opts.contactEmail}</a>
                </div>
              </li>
            </ul>
          </div>
          <SubscribeForm />
        </div>
        {/*<div className="row justify-content-center mb-2">*/}
        {/*  <div className="col-auto">*/}
        {/*    <ul className="nav">*/}
        {/*      {opts.social.map(it => (*/}
        {/*        <li key={it.name} className="nav-item">*/}
        {/*          <a href={it.url} className="nav-link">*/}
        {/*            <img*/}
        {/*              className="icon "*/}
        {/*              src={`assets/img/icons/social/${it.name}.svg`}*/}
        {/*              alt={`${it.name} social icon`}*/}
        {/*              data-inject-svg=""*/}
        {/*            />*/}
        {/*          </a>*/}
        {/*        </li>*/}
        {/*      ))}*/}
        {/*    </ul>*/}
        {/*  </div>*/}
        {/*</div>*/}
      </div>
    </footer>
  );
};
