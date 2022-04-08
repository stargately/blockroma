import { Ctx, Field, ObjectType, Query, Resolver } from "type-graphql";
import { ResolverContext } from "@/api-gateway/resolver-context";

@ObjectType()
class ChainMeta {
  @Field()
  avgBlockTimeSec: number;

  @Field()
  totalTransactions: number;

  @Field()
  blockHeight: number;

  @Field()
  totalAddresses: number;
}

@Resolver()
export class ChainMetaResolver {
  @Query(() => ChainMeta, { description: "Gets the blockchain's metadata" })
  async chainMeta(@Ctx() ctx: ResolverContext): Promise<ChainMeta> {
    const meta = await ctx.service.indexedChainService.chainMeta();

    return {
      ...meta,
      // TODO(dora): read from genesis
      avgBlockTimeSec: 5,
    };
  }
}
