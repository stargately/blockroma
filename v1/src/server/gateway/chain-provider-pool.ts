import { ethers } from "ethers";
import { Network } from "@ethersproject/networks/src.ts/types";

export class ChainProviderPool {
  providers: ethers.providers.JsonRpcProvider[] = [];

  cur = 0;

  constructor(urls: string[], network: Network) {
    for (const conn of urls) {
      this.providers.push(
        new ethers.providers.JsonRpcProvider(
          {
            url: conn,
            timeout: 5000,
          },
          {
            name: network.name,
            chainId: network.chainId,
          }
        )
      );
    }
  }

  get() {
    this.cur = (this.cur + 1) % this.providers.length;
    return this.providers[this.cur];
  }
}
