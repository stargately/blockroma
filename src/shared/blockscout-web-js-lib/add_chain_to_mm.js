const mmCfg = {
  chainId: 7778,
  chainName: "BoomMo Chain",
  symbol: "BMO",
  rpcUrl: "https://api-testnet.boommo.com",
  decimals: 18,
  networkPath: "",
}

module.exports = async function addChainToMM () {
  try {
    const chainID = await window.ethereum.request({ method: 'eth_chainId' })
    const chainIDFromEnvVar = mmCfg.chainId
    const chainIDHex = chainIDFromEnvVar && `0x${chainIDFromEnvVar.toString(16)}`
    const blockscoutURL = `${window.location.protocol  }//${  window.location.host  }${mmCfg.networkPath}`
    if (chainID !== chainIDHex) {
      await window.ethereum.request({
        method: 'wallet_addEthereumChain',
        params: [{
          chainId: chainIDHex,
          chainName: mmCfg.chainName,
          nativeCurrency: {
            name: mmCfg.symbol,
            symbol: mmCfg.symbol,
            decimals: mmCfg.decimals,
          },
          rpcUrls: [mmCfg.rpcUrl],
          blockExplorerUrls: [blockscoutURL]
        }]
      })
    }
    return true;
  } catch (error) {
    console.error(error)
    return false;
  }
}
