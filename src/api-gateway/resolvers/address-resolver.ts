import { Args, Ctx, FieldResolver, Resolver, Root } from "type-graphql";
import { Address } from "@/api-gateway/resolvers/types/address-type";
import { ConnectionArgs } from "@/api-gateway/resolvers/types/connection-type";
import { ResolverContext } from "@/api-gateway/resolver-context";
import { logger } from "onefx/lib/integrated-gateways/logger";
import { UserInputError } from "apollo-server-koa";
import { emptyPage } from "@/server/service/indexed-chain-service";
import { TransactionConnection } from "@/api-gateway/resolvers/types/transaction-type";

@Resolver(() => Address)
export class AddressResolver {
  @FieldResolver(() => TransactionConnection)
  async transactions(
    @Root() address: Address,
    @Args() args: ConnectionArgs,
    @Ctx() ctx: ResolverContext
  ): Promise<TransactionConnection> {
    // https://relay.dev/graphql/connections.htm check params
    if (args.first && args.last) {
      throw new UserInputError(
        "Including a value for both first and last is strongly discouraged"
      );
    }

    try {
      const txsWithPageInfo =
        await ctx.service.indexedChainService.getTransactions(
          { addressHash: address.hash },
          args
        );

      const txs = txsWithPageInfo.data.map(
        // @ts-ignore
        (tx) => ({
          node: { ...tx, id: txHashId(tx.hash) },
          cursor: (args.first ? args.after : args.before) ?? "",
        })
      );

      return {
        edges: txs,
        pageInfo: txsWithPageInfo.pageInfo,
      };
    } catch (err) {
      logger.error(
        `failed to get transactions: ${err}: ${
          err instanceof Error && err.stack
        }`
      );
      return {
        edges: [],
        pageInfo: emptyPage,
      };
    }
  }
}

function txHashId(hash: Buffer) {
  return Buffer.from(`Transaction:0x${hash.toString("hex")}`, "utf8").toString(
    "base64"
  );
}
