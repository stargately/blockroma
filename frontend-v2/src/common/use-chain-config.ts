import { useSelector } from "react-redux";

export type ChainConfig = {
  chainId: number;
  chainName: string;
  symbol: string;
  rpcUrls: string[];
  decimals: number;
  networkPath: string;
};

export function useChainConfig(): ChainConfig {
  return {
    chainId: 7778,
    chainName: "Blockroma",
    symbol: "BLO",
    rpcUrls: ["https://example.com"],
    decimals: 18,
    networkPath: "",
  };
}
