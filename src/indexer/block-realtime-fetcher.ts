import { MyServer } from "@/server/start-server";
import { logger } from "onefx/lib/integrated-gateways/logger";

export function mountBlockRealtimeFetcher(server: MyServer) {
  if (!server.config.indexer.realtime.enabled) {
    return;
  }

  server.gateways.chainProvider.on("block", async (blockNumber: number) => {
    logger.info(`listening on block ${blockNumber}`);
    await server.service.importerService.importRange([
      blockNumber,
      blockNumber,
    ]);
  });
}
