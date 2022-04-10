import { useSelector } from "react-redux";

export type ChainConfig = {
  chainId: number;
  chainName: string;
  symbol: string;
  rpcUrl: string;
  decimals: number;
  networkPath: string;
};

export function useChainConfig(): ChainConfig {
  return useSelector(
    (state: { base: { chain: ChainConfig } }) => state.base.chain
  );
}
