const routePrefix = "/ethw/mainnet";

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
    chainId: 10001,
    chainName: "ETHW-mainnet",
    symbol: "ETHW",
    rpcUrls: ["https://api.blockeden.xyz/ethw/QyLuprp5WC5y2XsVB65o"],
    decimals: 18,
    networkPath: routePrefix,
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
