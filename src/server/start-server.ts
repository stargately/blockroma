import config from "config";
import { Config, Server } from "onefx/lib/server";
import { setModel } from "@/model";
import { mountBlockRealtimeFetcher } from "@/indexer/block-realtime-fetcher";
import { mountBlockCatchupFetcher } from "@/indexer/block-catchup-fetcher";
import { logger } from "onefx/lib/integrated-gateways/logger";
import { ChainConfig } from "@/shared/common/use-chain-config";
import { Gateways, setGateways } from "./gateway/gateway";
import { setMiddleware } from "./middleware";
import { setServerRoutes } from "./server-routes";
import { Service, setService } from "./service";

export type MyServer = Server & {
  gateways: Gateways;
  config: Config & {
    chain: ChainConfig;
    indexer: {
      catchup: {
        enabled?: boolean;
        blocksBatchSize: number;
        blockNumberRanges: Array<[number | "latest", number | "latest"]>;
      };
      realtime: {
        enabled?: boolean;
      };
    };
    gateways: {
      postgresql: {
        uri: string;
        ssl: boolean;
      };
    };
    server: {
      siteOrigin: string;
    };
  };
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  resolvers: any;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  model: any;
  service: Service;
};

export async function startServer(): Promise<Server> {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const server: MyServer = new Server(config as any as Config) as MyServer;
  server.app.proxy = Boolean(config.get("server.proxy"));
  logger.info(`chain config: ${JSON.stringify(server.config.chain, null, 2)}`);
  await setGateways(server);
  setMiddleware(server);
  setModel(server);
  setServerRoutes(server);
  setService(server);
  mountBlockRealtimeFetcher(server);

  mountBlockCatchupFetcher({ server }).catch((err) => {
    logger.error(
      `failed to mountBlockCatchupFetcher: ${err}\n${
        err instanceof Error && err.stack
      }`
    );
  });

  const port = Number(process.env.PORT || config.get("server.port"));
  server.listen(port);
  return server;
}
