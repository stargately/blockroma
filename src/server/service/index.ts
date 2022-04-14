import { RemoteChainService } from "@/server/service/remote-chain-service/remote-chain-service";
import { MyServer } from "@/server/start-server";
import { IndexedChainService } from "@/server/service/indexed-chain-service";
import { ImporterService } from "@/server/service/importer-service";

export type Service = {
  remoteChainService: RemoteChainService;
  indexedChainService: IndexedChainService;
  importerService: ImporterService;
};

export function setService(server: MyServer): void {
  server.service = server.service ?? {};
  server.service.remoteChainService = new RemoteChainService(server);
  server.service.indexedChainService = new IndexedChainService(server);
  server.service.importerService = new ImporterService(server);
}
