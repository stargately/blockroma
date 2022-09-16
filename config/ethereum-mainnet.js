const routePrefix = "/eth/mainnet";

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
    chainId: 1,
    chainName: "Ethereum",
    symbol: "ETH",
    rpcUrls: JSON.parse(process.env.ETHEREUM_MAINNET_RPC_URLS),
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
