import { Gateways } from "@/server/gateway/gateway";
import { Service } from "@/server/service";

export type ResolverContext = {
  gateways: Gateways;
  service: Service;
};
