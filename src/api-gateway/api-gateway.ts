import { ApolloServer } from "apollo-server-koa";
import path from "path";
import "reflect-metadata";
import { buildSchema } from "type-graphql";
import { MyServer } from "@/server/start-server";
import { NonEmptyArray } from "type-graphql/dist/interfaces/NonEmptyArray";
import { ExplorerResolver } from "@/api-gateway/resolvers/explorer-resolver";
import { AddressResolver } from "@/api-gateway/resolvers/address-resolver";
import { ResolverContext } from "@/api-gateway/resolver-context";
import { ChainMetaResolver } from "@/api-gateway/resolvers/chain-meta.resolver";
import { ApolloServerPluginLandingPageGraphQLPlayground } from "apollo-server-core";
import { MetaResolver } from "./resolvers/meta-resolver";

export async function setApiGateway(server: MyServer): Promise<void> {
  const resolvers: NonEmptyArray<Function> = [
    MetaResolver,
    ExplorerResolver,
    AddressResolver,
    ChainMetaResolver,
  ];
  server.resolvers = resolvers;

  const sdlPath = path.resolve(__dirname, "api-gateway.graphql");
  const schema = await buildSchema({
    resolvers,
    emitSchemaFile: {
      path: sdlPath,
      commentDescriptions: true,
    },
    validate: false,
    nullableByDefault: true,
  });

  const apollo = new ApolloServer({
    schema,
    introspection: true,
    plugins: [ApolloServerPluginLandingPageGraphQLPlayground()],
    context: async (): Promise<ResolverContext> => ({
      gateways: server.gateways,
      service: server.service,
    }),
  });
  const gPath = `${server.config.server.routePrefix || ""}/api-gateway/`;
  await apollo.start();
  apollo.applyMiddleware({ app: server.app, path: gPath });
}
