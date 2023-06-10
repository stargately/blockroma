// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require("prism-react-renderer/themes/github");
const { beConfig } = require("./src/config");
const darkCodeTheme = require("prism-react-renderer/themes/dracula");

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: "blockroma.com",
  titleDelimiter: "-",
  tagline: beConfig.tagline,
  url: "https://blockroma.com",
  baseUrl: "/",
  onBrokenLinks: "throw",
  onBrokenMarkdownLinks: "warn",
  favicon: "img/blockroma-favicon.png",
  trailingSlash: true,

  headTags: [
    {
      tagName: "link",
      attributes: {
        rel: "shortcut icon",
        href: "https://blockroma.com/img/favicon.svg",
      },
    },
    {
      tagName: "link",
      attributes: {
        rel: "apple-touch-icon",
        href: "https://blockroma.com/img/favicon.svg",
      },
    },
  ],

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: "facebook", // Usually your GitHub org/user name.
  projectName: "docusaurus", // Usually your repo name.

  // Even if you don't use internalization, you can use this field to set useful
  // metadata like html lang. For example, if your site is Chinese, you may want
  // to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: "en",
    locales: ["en"],
  },

  presets: [
    [
      "classic",
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          editUrl: "https://github.com/stargately/blockroma/edit/main",
          sidebarCollapsed: false,
          sidebarPath: require.resolve("./sidebars.js"),
        },
        blog: {
          showReadingTime: true,
        },
        theme: {
          customCss: [
            require.resolve("./src/css/custom.css"),
            require.resolve("./static/assets/css/theme-saas.css"),
          ],
        },
        gtag: {
          trackingID: "G-VKPSYZ2K22",
        },
        sitemap: {
          changefreq: "daily",
          priority: 0.5,
          filename: "sitemap.xml",
        },
      }),
    ],
  ],

  stylesheets: ["https://fonts.googleapis.com/css?family=Inter"],

  themeConfig:
    /** @type {import('docusaurus-preset-openapi').ThemeConfig} */
    ({
      metadata: [
        {
          name: "keywords",
          content:
            "Blockroma, Open-source, EVM-compatible, blockchain explorer",
        },
      ],

      navbar: {
        title: "Blockroma",
        logo: {
          alt: "blockroma.com",
          src: "img/blockroma-favicon.png",
        },
        items: [
          {
            to: "docs/intro",
            activeBasePath: "docs",
            label: "Docs",
            position: "left",
          },
          {
            to: "/blog/",
            label: "Blog",
            position: "left",
          },
          {
            to: "https://blockroma.com/ethw/mainnet/",
            label: "Launch App",
            position: "right",
          },
        ],
      },
      footer: {
        style: "dark",
        links: [],
        copyright: `Copyright © ${new Date().getFullYear()} Blockroma`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },
      colorMode: {
        defaultMode: "light",
        disableSwitch: true,
        respectPrefersColorScheme: false,
      },
    }),
};

module.exports = config;
