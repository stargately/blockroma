const routePrefix = "/ethw/iceberg";

module.exports = {
  server: {
    siteOrigin: process.env.SITE_ORIGIN,
    routePrefix,
    noSecurityHeadersRoutes: {
      [`${routePrefix}/api-gateway/`]: true,
      [`${routePrefix}/api/`]: true,
    },
    noCsrfRoutes: {
      [`${routePrefix}/api-gateway/`]: true,
      [`${routePrefix}/api/`]: true,
    },
  },
  chain: {
    chainId: 10002,
    chainName: "ethw-iceberg-testnet",
    symbol: "ETHW",
    rpcUrls: ["https://iceberg.ethereumpow.org"],
    decimals: 18,
    networkPath: "",
  },
  indexer: {
    catchup: {
      enabled: false,
    },
    realtime: {
      enabled: true,
    },
  },
};
