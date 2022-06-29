import { MyServer } from "@/server/start-server";
import { createConnection, DataSource } from "typeorm";
import { ChainProviderPool } from "@/server/gateway/chain-provider-pool";

export type Gateways = {
  chainProvider: ChainProviderPool;
  dbCon: DataSource;
};

export async function setGateways(server: MyServer): Promise<void> {
  server.gateways = server.gateways || {};
  if (server.config.gateways.postgresql.uri) {
    server.gateways.dbCon = await createConnection({
      type: "postgres",
      url: server.config.gateways.postgresql.uri,
      entities: [`${__dirname}/../../model/**/*{.js,.ts}`],
      synchronize: true,
      ssl: server.config.gateways.postgresql.ssl,
    });
  }

  const chainConfig = server.config.chain;
  server.gateways.chainProvider = new ChainProviderPool(chainConfig.rpcUrls, {
    name: chainConfig.chainName,
    chainId: chainConfig.chainId,
  });
}
