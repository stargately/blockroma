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
      enabled: true,
      blocksBatchSize: 200,

      blockNumberRanges: [
        [1039, 1039],
        [225791, 225791],
        [1044, 1044],
        [255055, 255055],
        [276263, 276263],
        // [0, "latest"],
      ],
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
  // csp: {
  //   enabled:false,
  //   "default-src": ["none"],
  //   "manifest-src": ["self"],
  //   "style-src": ["self", "unsafe-inline", "https://fonts.googleapis.com/css"],
  //   "frame-src": [],
  //   "connect-src": [
  //     "self",
  //     "https://www.google-analytics.com/",
  //     ...(process.env.API_GATEWAY_URL ? [process.env.API_GATEWAY_URL] : []),
  //   ],
  //   "child-src": ["self"],
  //   "font-src": ["self", "data:", "https://fonts.gstatic.com/"],
  //   "img-src": ["*", "data:"],
  //   "media-src": ["self"],
  //   "object-src": ["self"],
  //   "script-src": ["self", "https://www.googletagmanager.com/"],
  // },
};
