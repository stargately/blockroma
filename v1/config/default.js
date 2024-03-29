const { config } = require("dotenv");
config();

const routePrefix = "";

module.exports = {
  project: "blockroma",
  chain: {
    chainId: 10001,
    chainName: "ETHW-mainnet",
    symbol: "ETHW",
    rpcUrls: ["https://api.blockeden.xyz/ethw/QyLuprp5WC5y2XsVB65o"],
    decimals: 18,
    networkPath: routePrefix,
  },
  server: {
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
      enabled: true,
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
