import React, { useEffect } from "react";
import useDocusaurusContext from "@docusaurus/useDocusaurusContext";
import Layout from "@theme/Layout";
import { useAos } from "@site/src/components/use-aos";
import { SeoHead } from "@site/src/components/seo-head";
import { Birdview } from "@site/src/components/birdview";
import { CompetitiveAnalysis } from "@site/src/components/competitive-analysis";
import { SubscribeFormLarge } from "@site/src/components/subscribe-form/subscribe-form-large";
import { SiteFooter } from "@site/src/site-footer";
import { useSvg } from "@site/src/components/use-svg";
const { beConfig } = require("../config");

const opts = {
  title: "Blockroma",

  heroTitle: "Open-source EVM-compatible blockchain explorer",
  heroDesc:
    "Built with TypeScript, React, and PostgreSQL. Easy to customize and deploy in minutes.",

  heroCtaText: "Launch App",
  heroCtaTextLink: "https://blockroma.com/ethw/mainnet/",
  heroCtaTextSecondary: "Fork me on Github",
  heroCtaTextSecondaryLink: "https://github.com/stargately/blockroma",

  heroCtaSecondaryText: "Interested in a hassle-free managed solution?",
  heroCtaSecondaryTextButton: "Contact Sales",
  heroCtaSecondaryTextButtonLink: "https://t.me/mikethrift",

  featTitle: "Solution tailored to the fast-moving web3 industry",

  feat1Title: "Open Source",
  feat1Desc:
    "A blockchain explorer visualizing your blockchain transactions, blocks, ERC20 and ERC721 tokens in realtime, with a community-driven development.",
  feat1CtaText: "Learn More",
  feat1CtaLink: "https://stargately.com/docs/blockroma",

  feat2Title: "Modern Stack",
  feat2Desc:
    "Extensibility and maintainability achieved by TypeScript, React, GraphQL, and PostgreSQL with the largest developer community. Deploy in <strike>hours</strike> minutes.",
  feat2CtaText: "Learn More",
  feat2CtaLink: "https://github.com/stargately/blockroma",

  feat3Title: "Managed Hosts",
  feat3Desc:
    "We provide custom themes, premium supports, and priority updates to managed hosting on blockroma.com. It's worry-free with zero operation costs for you.",
  feat3CtaText: "Contact Sales",
  feat3CtaLink: "https://t.me/mikethrift",

  onboardingTitle: "Fits into your needs",
  onboardingDesc: "Integrations for agility and extensibility",

  integrations: [
    {
      icon: "assets/img/logos/product/eth.svg",
      title: "Ethereum",
      desc: "Live Demo",
      link: "https://blockroma.com/eth/mainnet/",
    },
    {
      icon: "assets/img/logos/product/graphql.svg",
      title: "GraphQL",
      desc: "With playground",
      link: "https://blockroma.com/eth/mainnet/api-gateway/",
    },
    {
      icon: "assets/img/logos/product/postgresql.svg",
      title: "PostgreSQL",
      desc: "For Indexed Data",
      link: "https://stargately.com/docs/blockroma/#how-does-it-work",
    },
  ],

  integrationsCta: "",
  integrationsCtaLink: "#",

  onboardingSteps: [
    [
      "1. Download source code",
      "git clone https://github.com/stargately/blockroma",
    ],
    [
      "2. Customize to your chain",
      "follow instructions",
      "https://stargately.com/docs/blockroma#how-does-it-work",
    ],
    ["3. Deploy the explorer", "use any NodeJS cloud vendor"],
    [
      "Or you need a managed solution?",
      "contact us",
      "https://t.me/mikethrift",
    ],
  ],
};

export default function Home(): JSX.Element {
  const { siteConfig } = useDocusaurusContext();
  useEffect(() => {
    document.querySelector("body")?.classList.add("loaded");
  }, []);

  useAos();
  useSvg();

  return (
    <Layout title={beConfig.title} description={siteConfig.tagline}>
      <SeoHead />

      <main>
        <section className="has-divider">
          <div className="container">
            <div className="row align-items-center justify-content-between o-hidden">
              <div
                className="col-md-6 order-sm-2 mb-5 mb-sm-0"
                data-aos="fade-left"
              >
                <img src="assets/img/saas-3.svg" alt="Image" />
              </div>
              <div className="col-md-6 pr-xl-5 order-sm-1">
                <h1 className="display-4">{opts.heroTitle}</h1>
                <p className="lead">{opts.heroDesc}</p>
                <form className="d-sm-flex mb-2 mt-4">
                  <a
                    className="btn btn-lg btn-primary mr-1 mb-1"
                    target="_blank"
                    href={opts.heroCtaTextLink}
                    rel="nofollow noreferrer noopener"
                  >
                    {opts.heroCtaText}
                  </a>
                  <a
                    className="btn btn-lg btn-secondary mr-1 mb-1"
                    target="_blank"
                    href={opts.heroCtaTextSecondaryLink}
                    rel="nofollow noreferrer noopener"
                  >
                    {opts.heroCtaTextSecondary}
                  </a>
                </form>
                <span className="text-small text-muted">
                  {opts.heroCtaSecondaryText}{" "}
                  <a href={opts.heroCtaSecondaryTextButtonLink}>
                    {opts.heroCtaSecondaryTextButton}
                  </a>
                </span>
              </div>
            </div>
          </div>
          <div className="divider">
            <img
              className="bg-primary-alt"
              src="assets/img/dividers/divider-1.svg"
              alt="divider graphic"
              data-inject-svg=""
            />
          </div>
        </section>
        <section className="bg-primary-alt">
          <div className="container">
            <div className="row mb-4">
              <div className="col">
                <h2 className="h1">{opts.featTitle}</h2>
              </div>
            </div>
            <div className="row">
              <div className="col-md-4" data-aos="fade-up" data-aos-delay={100}>
                <div className="card card-body">
                  <div className="icon-round mb-3 mb-md-4 bg-primary">
                    <img
                      className="icon bg-primary"
                      src="assets/img/icons/theme/communication/chat-check.svg"
                      alt="icon"
                      data-inject-svg=""
                    />
                  </div>
                  <h4>{opts.feat1Title}</h4>
                  <p>{opts.feat1Desc}</p>
                  {opts.feat1CtaLink && (
                    <a href={opts.feat1CtaLink}>{opts.feat1CtaText}</a>
                  )}
                </div>
              </div>
              <div className="col-md-4" data-aos="fade-up" data-aos-delay={200}>
                <div className="card card-body">
                  <div className="icon-round mb-3 mb-md-4 bg-primary">
                    <img
                      className="icon bg-primary"
                      src="assets/img/icons/theme/shopping/wallet-3.svg"
                      alt="icon"
                      data-inject-svg=""
                    />
                  </div>
                  <h4>{opts.feat2Title}</h4>
                  <p dangerouslySetInnerHTML={{ __html: opts.feat2Desc }}></p>
                  {opts.feat2CtaLink && (
                    <a href={opts.feat2CtaLink}>{opts.feat2CtaText}</a>
                  )}
                </div>
              </div>
              <div className="col-md-4" data-aos="fade-up" data-aos-delay={300}>
                <div className="card card-body">
                  <div className="icon-round mb-3 mb-md-4 bg-primary">
                    <img
                      className="icon bg-primary"
                      src="assets/img/icons/theme/food/sushi.svg"
                      alt="icon"
                      data-inject-svg=""
                    />
                  </div>
                  <h4>{opts.feat3Title}</h4>
                  <p>{opts.feat3Desc}</p>
                  {opts.feat3CtaLink && (
                    <a href={opts.feat3CtaLink}>{opts.feat3CtaText}</a>
                  )}
                </div>
              </div>
            </div>
          </div>
        </section>

        <Birdview />

        <section className="bg-primary-alt">
          <div className="container">
            <div className="row justify-content-between">
              <div className="col-xl-5 col-lg-6" data-aos="fade-right">
                <div className="row justify-content-center">
                  <div className="col-md-8 col-lg mb-4 mb-lg-0">
                    <img src="assets/img/saas-2.svg" alt="Image" />
                  </div>
                </div>
              </div>
              <div className="col-lg-6">
                <h3 className="h1">{opts.onboardingTitle}</h3>
                <p className="lead">{opts.onboardingDesc}</p>
                <div className="row mt-5 o-hidden">
                  {opts.integrations.map((integration, idx) => (
                    <div
                      key={integration.title}
                      className="col-sm-6"
                      data-aos="fade-left"
                      data-aos-delay={idx * 100}
                    >
                      <a
                        href={integration.link}
                        className="card card-sm card-body flex-row align-items-center hover-shadow-3d"
                      >
                        <img
                          className=""
                          src={integration.icon}
                          alt="icon"
                          data-inject-svg=""
                        />
                        <div className="ml-3">
                          <h5 className="mb-0">{integration.title}</h5>
                          <span>{integration.desc}</span>
                        </div>
                      </a>
                    </div>
                  ))}
                </div>
                {opts.integrationsCta && (
                  <a href={opts.integrationsCtaLink} className="hover-arrow">
                    {opts.integrationsCta}
                  </a>
                )}
              </div>
            </div>
          </div>
        </section>

        <section className="has-divider bg-primary-alt">
          <div className="container pt-5">
            <div className="row justify-content-between align-items-center o-hidden">
              <div className="col-md-6 col-lg-5">
                {opts.onboardingSteps.map((s) => (
                  <div
                    key={s[0]}
                    className="d-flex mb-4"
                    data-aos="fade-up"
                    data-aos-delay="NaN"
                  >
                    <div className="process-circle bg-primary" />
                    <div className="ml-3">
                      <h5>{s[0]}</h5>
                      {s[2] ? (
                        <a href={s[2]} rel="noopener nofollow noreferrer">
                          {s[1]}
                        </a>
                      ) : (
                        <p>{s[1]}</p>
                      )}
                    </div>
                  </div>
                ))}
              </div>
              <div className="col-md-6" data-aos="fade-left">
                <img src="assets/img/saas-4.svg" alt="Image" />
              </div>
            </div>
          </div>
          <div className="divider">
            <img
              src="assets/img/dividers/divider-2.svg"
              alt="graphical divider"
              data-inject-svg=""
            />
          </div>
        </section>

        <CompetitiveAnalysis />

        <section className="p-0">
          <div className="container">
            <div className="row justify-content-center">
              <div className="col-xl-6 col-lg-7 col-md-8 mb-lg-n7 layer-2">
                <img
                  src="assets/img/saas-1.svg"
                  alt="Image"
                  data-aos="fade-up"
                />
              </div>
            </div>
          </div>
        </section>

        <section className="bg-primary text-light has-divider">
          <div className="divider flip-y">
            <img
              src="assets/img/dividers/divider-3.svg"
              alt="graphical divider"
              data-inject-svg=""
            />
          </div>
          <SubscribeFormLarge />
        </section>

        <SiteFooter />

        <a
          href="#"
          className="btn back-to-top btn-primary btn-round"
          data-smooth-scroll
          data-aos="fade-up"
          data-aos-offset="2000"
          data-aos-mirror="true"
          data-aos-once="false"
        >
          <img
            className="icon"
            src="assets/img/icons/theme/navigation/arrow-up.svg"
            alt="arrow-up icon"
            data-inject-svg
          />
        </a>
      </main>
    </Layout>
  );
}
