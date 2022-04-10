import { MyServer } from "@/server/start-server";
import { createConnection, DataSource } from "typeorm";
import { ethers, providers } from "ethers";

export type Gateways = {
  chainProvider: providers.JsonRpcProvider;
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
  server.gateways.chainProvider = new ethers.providers.JsonRpcProvider(
    {
      url: chainConfig.rpcUrl,
      timeout: 5000,
    },
    {
      name: chainConfig.chainName,
      chainId: chainConfig.chainId,
    }
  );
}
