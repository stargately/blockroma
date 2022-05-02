module.exports = {
  server: {
    routePrefix: "/eth/mainnet",
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
