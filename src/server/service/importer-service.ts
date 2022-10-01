import { logger } from "onefx/lib/integrated-gateways/logger";
import { MyServer } from "@/server/start-server";
import { Address } from "@/model/address";
import { TokenTransfer } from "@/model/token-transfer";
import { Token } from "@/model/token";
import { ethers } from "ethers";
import { tokenMetadataAbi } from "@/server/service/remote-chain-service/token-metadata-abi";

export class ImporterService {
  server: MyServer;

  constructor(server: MyServer) {
    this.server = server;
  }

  async importRange(range: [number, number]): Promise<void> {
    try {
      const { transactions, blocks, addresses, tokenTransfers } =
        await this.server.service.remoteChainService.fetchBlockByRange(
          range[0],
          range[1]
        );

      await this.upsertAddresses(addresses);
      await this.server.service.indexedChainService.insertBlocks(blocks);
      await this.server.service.indexedChainService.insertTransactions(
        transactions
      );
      await this.trackErc20ByTokenTransfers(tokenTransfers);
      await this.server.gateways.dbCon
        .getRepository(TokenTransfer)
        .insert(tokenTransfers);
    } catch (err) {
      // TODO(dora): should handle error
      logger.error(
        `failed to process range ${JSON.stringify(range)}, ${err} ${
          err instanceof Error && err.stack
        }`
      );
    }
  }

  async upsertAddresses(addresses: Address[]): Promise<void> {
    try {
      await this.server.service.indexedChainService.insertAddresses(addresses);
    } catch (err) {
      logger.warn(`duplicate address inserted: ${err}`);
    }

    await Promise.all(
      addresses.map((addr) =>
        this.server.service.indexedChainService.updateAddressForHigherBlock(
          addr
        )
      )
    );
  }

  private async trackErc20ByTokenTransfers(
    tokenTransfers: TokenTransfer[]
  ): Promise<void> {
    // deduplicate token contract addresses
    const dedupErc20Address = new Set<Buffer>();
    tokenTransfers.filter(tt => tt.type === "ERC-20").forEach(tt => {
      dedupErc20Address.add(tt.tokenContractAddress);
    })

    // update token table with contracts by their addresses
    for (const address of dedupErc20Address) {
      // eslint-disable-next-line no-await-in-loop
      const storedToken = await this.server.gateways.dbCon
        .getRepository(Token)
        .findOneBy({ contractAddress: address });
      if (!storedToken || !storedToken.skipMetadata) {
        const contract = new ethers.Contract(
          `0x${address.toString("hex")}`,
          tokenMetadataAbi,
          this.server.gateways.chainProvider.get()
        );
        // eslint-disable-next-line no-await-in-loop
        const [name, decimals, totalSupply, symbol] = await Promise.all([
          contract.name(),
          contract.decimals(),
          contract.totalSupply(),
          contract.symbol(),
        ]);
        const token: Token = {
          contractAddress: address,
          name,
          decimals,
          totalSupply: totalSupply.toString(),
          symbol,
          type: "ERC-20",
        };
        // eslint-disable-next-line no-await-in-loop
        await this.server.gateways.dbCon
          .getRepository(Token)
          .upsert(token, ["contractAddress"]);
      }
    }
  }
}
