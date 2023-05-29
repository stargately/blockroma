import React from "react";
import classnames from "classnames";

export const CompetitiveAnalysis: React.FC = () => {
  const meta = {
    competitiveAnalysisTitle: "Why choose Blockroma?",
    competitiveAnalysisDesc:
      "It's built for blockchains that are moving fast with modern stack.",
    competitiveAnalysisDimensions: [
      "Multi-chain support",
      "EVM-compatible",
      "Deploy in minutes",
      "Easy to extend",
      "Integrated wallet",
    ],
    competitors: [
      {
        head: "/img/blockroma-favicon.png",
        primary: true,
        items: [true, true, true, true, true],
      },
      {
        head: "https://tp-misc.b-cdn.net/blockroma-blockscout-logo.jpg",
        items: [true, true, false, false, false],
      },
      {
        head: "https://tp-misc.b-cdn.net/blockroma-mintscan-logo.jpg",
        items: [true, false, false, true, false],
      },
    ],
    competitiveAnalysisCtaText: "",
    competitiveAnalysisCtaLink: "#",

    recommendations: [
      {
        icon: "assets/img/avatars/jensen-w.jpg",
        words: "We develop upon it for our blockchain network.",
        name: "Jensen W.",
        title: "Engineering Manager, MBA",
      },
      {
        icon: "assets/img/avatars/yuan-l.jpg",
        words:
          "It's amazing to set up Blockroma in 10 minutes instead of hours.",
        name: "Yuan L.",
        title: "Software Engineer",
      },
      {
        icon: "assets/img/avatars/nico-m.jpg",
        words:
          "Blockroma simplifies our tech stack but I am hoping to see more features.",
        name: "Nico M.",
        title: "Blockchain Engineer",
      },
    ],
  };
  return (
    <section>
      <div className="container">
        <div className="row justify-content-center text-center mb-6">
          <div className="col-xl-8 col-lg-9">
            <h2 className="display-4 mx-xl-6">
              {meta.competitiveAnalysisTitle}
            </h2>
            <p className="lead">{meta.competitiveAnalysisDesc}</p>
          </div>
        </div>
        <div className="row justify-content-center">
          <div className="" data-aos="fade-up">
            <table className="table1 pricing-table pricing-table-competitors">
              <thead>
                <tr>
                  <th scope="col" />
                  {meta.competitors.map((competitor) => (
                    <th
                      scope="col"
                      className={classnames({
                        "bg-primary-alt": competitor.primary,
                      })}
                      key={competitor.head}
                    >
                      <img
                        src={competitor.head}
                        alt="Image"
                        width={48}
                        height={48}
                      />
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {meta.competitiveAnalysisDimensions.map((dimension, row) => (
                  <tr key={dimension}>
                    <th scope="row">
                      <span className="h6 mb-0 d-block">{dimension}</span>
                      <span className="text-small text-muted" />
                    </th>
                    {meta.competitors.map((competitor) =>
                      competitor.items[row] ? (
                        <td
                          className={classnames({
                            "bg-primary-alt": competitor.primary,
                          })}
                          key={competitor.head}
                        >
                          <div
                            className={classnames(
                              "icon-round icon-round-xs",
                              competitor.primary
                                ? "bg-success"
                                : "bg-primary-3-alt"
                            )}
                          >
                            <img
                              className={classnames(
                                "icon",
                                competitor.primary ? "bg-light" : "bg-primary-3"
                              )}
                              src="assets/img/icons/interface/check.svg"
                              alt="check interface icon"
                              data-inject-svg=""
                            />
                          </div>
                        </td>
                      ) : (
                        <td key={competitor.head} />
                      )
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
            {meta.competitiveAnalysisCtaText && (
              <div className="text-right">
                <a
                  href={meta.competitiveAnalysisCtaLink}
                  className="hover-arrow"
                >
                  {meta.competitiveAnalysisCtaText}
                </a>
              </div>
            )}
          </div>
        </div>

        <div className="row justify-content-center mt-6">
          <div className="col-xl-7 col-lg-8 col-md-10">
            <div>
              {meta.recommendations.map((r) => (
                <div className="carousel-cell mx-md-4" key={r.name}>
                  <div className="card card-body flex-row py-4">
                    <img
                      src={r.icon}
                      alt="Avatar"
                      className="avatar1 avatar-lg"
                    />
                    <div className="ml-3">
                      <h4>“{r.words}”</h4>
                      <div className="avatar-author d-block">
                        <h6>{r.name}</h6>
                        <span>{r.title}</span>
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};
