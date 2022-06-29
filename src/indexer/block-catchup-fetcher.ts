/* eslint-disable no-await-in-loop */
import { MyServer } from "@/server/start-server";
import { logger } from "onefx/lib/integrated-gateways/logger";
import { chunkRanges } from "@/server/service/utils/get-ranges";

type TaskOpts = {
  server: MyServer;
};

export async function mountBlockCatchupFetcher(opts: TaskOpts): Promise<void> {
  const config = opts.server.config.indexer.catchup;
  logger.info(`catchup indexer is enabled: ${config.enabled}`);
  if (!config.enabled) {
    return;
  }

  const { blockNumberRanges } = config;
  logger.debug(
    `block-catchup-fetcher:   catching up: ${JSON.stringify(blockNumberRanges)}`
  );

  for (const rootRange of blockNumberRanges) {
    const lastBlockNumber = (
      rootRange[1] === "latest"
        ? await opts.server.gateways.chainProvider.get().getBlockNumber()
        : rootRange[1]
    ) as number;

    const firstBlockNumber = (
      rootRange[0] === "latest"
        ? await opts.server.gateways.chainProvider.get().getBlockNumber()
        : rootRange[0]
    ) as number;

    const missingRanges =
      await opts.server.service.remoteChainService.missingBlockNumberRanges(
        firstBlockNumber,
        lastBlockNumber
      );

    logger.debug(
      `block-catchup-fetcher: missing range: ${JSON.stringify(missingRanges)}`
    );

    const chunkedRanges = chunkRanges(missingRanges, config.blocksBatchSize);

    for (const range of chunkedRanges) {
      logger.info(`catching up block [${range[0]}, ${range[1]}]`);
      await opts.server.service.importerService.importRange(range);
    }
  }

  logger.info("all blocks caught up");
}
