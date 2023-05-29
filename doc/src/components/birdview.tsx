import React, { useState } from "react";
import classNames from "classnames";

export const Birdview: React.FC = () => {
  const opts = {
    birdViewTitle: "Modern and stylish",
    birdViewDesc:
      "Welcome to our Blockchain Explorer, a seamless blend of modern design and effortless functionality. Echoing our ethos of simplicity and efficiency, our tool is engineered to deliver a user-centric experience in your blockchain interactions",

    birdViewTile1: "Realtime Explorer",
    birdViewTile1ImgSrc: "assets/img/icons/theme/shopping/box-2.svg",
    birdViewTile1Popovers: [],
    birdViewTile1Screenshot:
      "https://camo.githubusercontent.com/32af9e9099ba6746129b30c137ca6af438b26838d6c61eb8ff57dd678c8891ef/68747470733a2f2f74702d6d6973632e622d63646e2e6e65742f626c6f636b726f6d612d76302e312d322e676966",

    birdViewTile2: "Transactions",
    birdViewTile2ImgSrc: "assets/img/icons/theme/shopping/chart-pie.svg",
    birdViewTile2Screenshot:
      "https://tp-misc.b-cdn.net/blockroma-transactions.png",

    birdViewTile3: "Dark Mode",
    birdViewTile3ImgSrc: "assets/img/icons/theme/design/saturation.svg",
    birdViewTile3Screenshot:
      "https://tp-misc.b-cdn.net/blockroma-screenshot-dark-mode.png",
  };
  const [activeIdx, setActiveIdx] = useState(0);
  return (
    <section>
      <div className="container">
        <div className="row justify-content-center text-center mb-6">
          <div className="col-xl-8 col-lg-9">
            <h2 className="display-4 mx-xl-6">{opts.birdViewTitle}</h2>
            <p className="lead">{opts.birdViewDesc}</p>
          </div>
        </div>
        <div className="row justify-content-center mb-5">
          <div className="col-xl-11">
            <ul className="nav justify-content-center" role="tablist">
              <li className="nav-item mx-1">
                <a
                  className={classNames("nav-link", {
                    active: activeIdx === 0,
                  })}
                  href="#saas-tab-1"
                  data-toggle="tab"
                  role="tab"
                  aria-controls="saas-tab-1"
                  aria-selected="true"
                  onClick={(e) => {
                    e.preventDefault();
                    setActiveIdx(0);
                  }}
                >
                  <div className="icon-round icon-round-sm bg-primary">
                    <img
                      className="icon bg-primary"
                      src={opts.birdViewTile1ImgSrc}
                      alt="box-2 icon"
                      data-inject-svg=""
                    />
                  </div>
                  {opts.birdViewTile1}
                </a>
              </li>
              <li className="nav-item mx-1">
                <a
                  className={classNames("nav-link", {
                    active: activeIdx === 1,
                  })}
                  href="#saas-tab-2"
                  data-toggle="tab"
                  role="tab"
                  aria-controls="saas-tab-2"
                  aria-selected="false"
                  onClick={(e) => {
                    e.preventDefault();
                    setActiveIdx(1);
                  }}
                >
                  <div className="icon-round icon-round-sm bg-primary">
                    <img
                      className="icon bg-primary"
                      src={opts.birdViewTile2ImgSrc}
                      alt="chart-pie icon"
                      data-inject-svg=""
                    />
                  </div>
                  {opts.birdViewTile2}
                </a>
              </li>
              <li className="nav-item mx-1">
                <a
                  className={classNames("nav-link", {
                    active: activeIdx === 2,
                  })}
                  href="#saas-tab-3"
                  data-toggle="tab"
                  role="tab"
                  aria-controls="saas-tab-3"
                  aria-selected="false"
                  onClick={(e) => {
                    e.preventDefault();
                    setActiveIdx(2);
                  }}
                >
                  <div className="icon-round icon-round-sm bg-primary">
                    <img
                      className="icon bg-primary"
                      src={opts.birdViewTile3ImgSrc}
                      alt="saturation icon"
                      data-inject-svg=""
                    />
                  </div>
                  {opts.birdViewTile3}
                </a>
              </li>
            </ul>
          </div>
        </div>

        <div className="row justify-content-center">
          <div className="col-xl-11" data-aos="fade-up">
            <div className="tab-content">
              <div
                className={classNames("tab-pane fade", {
                  "show active": activeIdx === 0,
                })}
                id="saas-tab-1"
                role="tabpanel"
                aria-labelledby="saas-tab-1"
              >
                <div className="popover-image">
                  {/*<div*/}
                  {/*  className="popover-hotspot bg-primary-2"*/}
                  {/*  style={{top: "25%", left: "10%"}}*/}
                  {/*  data-toggle="popover"*/}
                  {/*  title="Hotspot title"*/}
                  {/*  data-content="And here's some amazing content. It's very engaging. Right?"*/}
                  {/*/>*/}
                  {/*<div*/}
                  {/*  className="popover-hotspot bg-primary-2"*/}
                  {/*  style={{top: "10%", left: "90%"}}*/}
                  {/*  data-toggle="popover"*/}
                  {/*  title="Hotspot title"*/}
                  {/*  data-content="And here's some amazing content. It's very engaging. Right?"*/}
                  {/*/>*/}
                  {/*<div*/}
                  {/*  className="popover-hotspot bg-primary-2"*/}
                  {/*  style={{top: "90%", left: "20%"}}*/}
                  {/*  data-toggle="popover"*/}
                  {/*  title="Hotspot title"*/}
                  {/*  data-content="And here's some more amazing content. It's very engaging. Right?"*/}
                  {/*/>*/}
                  <img
                    src={opts.birdViewTile1Screenshot}
                    alt="Image"
                    className="rounded border shadow-lg"
                  />
                </div>
              </div>
              <div
                className={classNames("tab-pane fade", {
                  "show active": activeIdx === 1,
                })}
                id="saas-tab-2"
                role="tabpanel"
                aria-labelledby="saas-tab-2"
              >
                <div className="popover-image">
                  {/*<div*/}
                  {/*  className="popover-hotspot bg-primary-2"*/}
                  {/*  style={{top: "25%", left: "90%"}}*/}
                  {/*  data-toggle="popover"*/}
                  {/*  title="Hotspot title"*/}
                  {/*  data-content="And here's some amazing content. It's very engaging. Right?"*/}
                  {/*/>*/}
                  {/*<div*/}
                  {/*  className="popover-hotspot bg-primary-2"*/}
                  {/*  style={{top: "90%", left: "10%"}}*/}
                  {/*  data-toggle="popover"*/}
                  {/*  title="Hotspot title"*/}
                  {/*  data-content="And here's some amazing content. It's very engaging. Right?"*/}
                  {/*/>*/}
                  {/*<div*/}
                  {/*  className="popover-hotspot bg-primary-2"*/}
                  {/*  style={{top: "30%", left: "60%"}}*/}
                  {/*  data-toggle="popover"*/}
                  {/*  title="Hotspot title"*/}
                  {/*  data-content="And here's some more amazing content. It's very engaging. Right?"*/}
                  {/*/>*/}
                  <img
                    src={opts.birdViewTile2Screenshot}
                    alt="Image"
                    className="rounded border shadow-lg"
                  />
                </div>
              </div>
              <div
                className={classNames("tab-pane fade", {
                  "show active": activeIdx === 2,
                })}
                id="saas-tab-3"
                role="tabpanel"
                aria-labelledby="saas-tab-3"
              >
                <div className="popover-image">
                  {/*<div*/}
                  {/*  className="popover-hotspot bg-primary-2"*/}
                  {/*  style={{top: "2%", left: "2%"}}*/}
                  {/*  data-toggle="popover"*/}
                  {/*  title="Hotspot title"*/}
                  {/*  data-content="And here's some amazing content. It's very engaging. Right?"*/}
                  {/*/>*/}
                  {/*<div*/}
                  {/*  className="popover-hotspot bg-primary-2"*/}
                  {/*  style={{top: "35%", left: "95%"}}*/}
                  {/*  data-toggle="popover"*/}
                  {/*  title="Hotspot title"*/}
                  {/*  data-content="And here's some amazing content. It's very engaging. Right?"*/}
                  {/*/>*/}
                  {/*<div*/}
                  {/*  className="popover-hotspot bg-primary-2"*/}
                  {/*  style={{top: "90%", left: "50%"}}*/}
                  {/*  data-toggle="popover"*/}
                  {/*  title="Hotspot title"*/}
                  {/*  data-content="And here's some more amazing content. It's very engaging. Right?"*/}
                  {/*/>*/}
                  <img
                    src={opts.birdViewTile3Screenshot}
                    alt="Image"
                    className="rounded border shadow-lg"
                  />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};
