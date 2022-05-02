const routePrefix = "/eth/mainnet";

module.exports = {
  server: {
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
    rpcUrl: process.env.ETHEREUM_MAINNET_RPC_URL,
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
