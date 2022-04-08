import { logger } from "onefx/lib/integrated-gateways/logger";
import { MyServer } from "@/server/start-server";

export class ImporterService {
  server: MyServer;

  constructor(server: MyServer) {
    this.server = server;
  }

  async importRange(range: [number, number]): Promise<void> {
    const { transactions, blocks, addresses } =
      await this.server.service.remoteChainService.fetchBlockByRange(
        range[0],
        range[1]
      );

    try {
      try {
        await this.server.service.indexedChainService.insertAddresses(
          addresses
        );
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

      await this.server.service.indexedChainService.insertBlocks(blocks);

      if (transactions.length) {
        await this.server.service.indexedChainService.insertTransactions(
          transactions
        );
      }
    } catch (err) {
      // TODO(dora): should handle error
      logger.error(
        `failed to process range ${JSON.stringify(range)}, ${err} ${err.stack}`
      );
    }
  }
}
