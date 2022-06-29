module.exports = async function addChainToMM (mmCfg) {
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
          rpcUrls: mmCfg.rpcUrls,
          blockExplorerUrls: [blockscoutURL]
        }]
      })
    }
    return true;
  } catch (error) {
    console.error(`failed to add chain to metamask: ${error}`)
    return false;
  }
}
