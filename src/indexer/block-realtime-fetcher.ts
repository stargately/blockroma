import { MyServer } from "@/server/start-server";
import { logger } from "onefx/lib/integrated-gateways/logger";

export function mountBlockRealtimeFetcher(server: MyServer) {
  logger.info(
    `realtime indexer is enabled: ${server.config.indexer.realtime.enabled}`
  );
  if (!server.config.indexer.realtime.enabled) {
    return;
  }

  server.gateways.chainProvider
    .get()
    .on("block", async (blockNumber: number) => {
      logger.info(`listening on block ${blockNumber}`);
      await server.service.importerService.importRange([
        blockNumber,
        blockNumber,
      ]);
    });
}
