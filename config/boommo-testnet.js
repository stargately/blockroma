const { config } = require("dotenv");
config();

const routePrefix = "/bmo/testnet";

module.exports = {
  project: "blockroma",
  chain: {
    chainId: 7778,
    chainName: "BoomMo Chain",
    symbol: "BMO",
    rpcUrls: ["https://api-testnet.boommo.com"],
    decimals: 18,
    networkPath: routePrefix,
  },
  server: {
    siteOrigin: process.env.SITE_ORIGIN,
    routePrefix,
    port: process.env.PORT || 4134,
    proxy: false,
    staticDir: "./dist",
    delayInitMiddleware: false,
    cookie: {
      secrets: ["insecure plain text", "insecure secret here"],
    },
    noSecurityHeadersRoutes: {
      [`${routePrefix}/api-gateway/`]: true,
      [`${routePrefix}/api/`]: true,
    },
    noCsrfRoutes: {
      [`${routePrefix}/api-gateway/`]: true,
      [`${routePrefix}/api/`]: true,
    },
  },
  indexer: {
    catchup: {
      enabled: false,
    },
    realtime: {
      enabled: false,
    },
  },
  gateways: {
    logger: {
      enabled: true,
      level: "debug",
    },
    postgresql: {
      uri: process.env.DATABASE_URL,
      ssl: String(process.env.DATABASE_URL).includes(".com"),
    },
  },
  analytics: {
    gaMeasurementId: "G-VKPSYZ2K22",
  },
};
