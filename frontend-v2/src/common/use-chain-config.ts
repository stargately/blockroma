export type ChainConfig = {
  chainId: number;
  chainName: string;
  symbol: string;
  rpcUrls: string[];
  decimals: number;
  networkPath: string;
};

export const chainConfig = {
  chainId: 7778,
  chainName: "Blockroma",
  symbol: "BLO",
  rpcUrls: ["https://example.com"],
  decimals: 18,
  networkPath: "",
};

export function useChainConfig(): ChainConfig {
  return chainConfig;
}
