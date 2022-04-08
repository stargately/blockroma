import { Args, ArgsType, Ctx, Field, Int, Query, Resolver } from "type-graphql";
import {
  Block,
  BlockConnection,
} from "@/api-gateway/resolvers/types/block-type";
import { Transaction } from "@/api-gateway/resolvers/types/transaction-type";
import {
  Address,
  TransactionConnection,
} from "@/api-gateway/resolvers/types/address-type";
import { ResolverContext } from "@/api-gateway/resolver-context";
import { ApolloError, UserInputError } from "apollo-server-koa";
import { BufferScalar } from "@/api-gateway/resolvers/types/buffer-scalar";
import { logger } from "onefx/lib/integrated-gateways/logger";
import { emptyPage } from "@/server/service/indexed-chain-service";
import * as Relay from "graphql-relay";
import QRCode from "qrcode";

@ArgsType()
class BlockRequest {
  @Field(() => Int, { nullable: false })
  number: number;
}

@ArgsType()
class AddressRequest {
  @Field(() => BufferScalar, { nullable: false })
  hash: Buffer;
}

@ArgsType()
class TransactionRequest {
  @Field(() => BufferScalar, { nullable: false })
  hash: Buffer;
}

function txHashId(category: string, hash: Buffer) {
  return Buffer.from(`${category}:0x${hash.toString("hex")}`, "utf8").toString(
    "base64"
  );
}

@ArgsType()
export class TransactionsArgs implements Relay.ConnectionArguments {
  @Field(() => Int)
  blockNumber: number;

  @Field(() => String, {
    nullable: true,
    description: "Paginate before opaque cursor",
  })
  before?: Relay.ConnectionCursor;

  @Field(() => String, {
    nullable: true,
    description: "Paginate after opaque cursor",
  })
  after?: Relay.ConnectionCursor;

  @Field(() => Number, { nullable: true, description: "Paginate first" })
  first?: number;

  @Field(() => Number, { nullable: true, description: "Paginate last" })
  last?: number;
}

@ArgsType()
export class BlocksArgs implements Relay.ConnectionArguments {
  @Field(() => String, {
    nullable: true,
    description: "Paginate before opaque cursor",
  })
  before?: Relay.ConnectionCursor;

  @Field(() => String, {
    nullable: true,
    description: "Paginate after opaque cursor",
  })
  after?: Relay.ConnectionCursor;

  @Field(() => Number, { nullable: true, description: "Paginate first" })
  first?: number;

  @Field(() => Number, { nullable: true, description: "Paginate last" })
  last?: number;
}

@Resolver()
export class ExplorerResolver {
  @Query(() => Block, { description: "Gets a block by number" })
  async block(
    @Args() request: BlockRequest,
    @Ctx() ctx: ResolverContext
  ): Promise<Block> {
    const block = await ctx.service.indexedChainService.numberToBlock(
      request.number
    );
    if (!block) {
      throw new ApolloError(
        `Block number ${request.number} was not found.`,
        "NOT_FOUND"
      );
    }
    const numTxs = await ctx.service.indexedChainService.getNumTxs({
      blockNumber: request.number,
    });
    return {
      ...block,
      numTxs,
    };
  }

  @Query(() => BlockConnection, {
    description: "Gets blocks by a range of number",
  })
  async blocks(
    @Args() args: BlocksArgs,
    @Ctx() ctx: ResolverContext
  ): Promise<BlockConnection> {
    const blksWithPageInfo = await ctx.service.indexedChainService.getBlocks(
      {},
      args
    );

    const blks = blksWithPageInfo.data.map(
      // @ts-ignore
      (blk) => ({
        node: { ...blk, id: txHashId("Block", blk.hash) },
        cursor: (args.first ? args.after : args.before) ?? "",
      })
    );

    return {
      edges: blks,
      pageInfo: blksWithPageInfo.pageInfo,
    };
  }

  @Query(() => TransactionConnection)
  async transactions(
    @Args() args: TransactionsArgs,
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
          { blockNumber: args.blockNumber },
          args
        );

      const txs = txsWithPageInfo.data.map(
        // @ts-ignore
        (tx) => ({
          node: { ...tx, id: txHashId("Transaction", tx.hash) },
          cursor: (args.first ? args.after : args.before) ?? "",
        })
      );

      return {
        edges: txs,
        pageInfo: txsWithPageInfo.pageInfo,
      };
    } catch (err) {
      logger.error(`failed to get transactions: ${err}: ${err.stack}`);
      return {
        edges: [],
        pageInfo: emptyPage,
      };
    }
  }

  @Query(() => Transaction)
  async transaction(
    @Args() request: TransactionRequest,
    @Ctx() ctx: ResolverContext
  ): Promise<Transaction> {
    const tx = await ctx.service.indexedChainService.getTransactionByHash(
      request.hash
    );
    if (!tx) {
      throw new ApolloError(
        `Transaction hash 0x${request.hash} was not found.`,
        "NOT_FOUND"
      );
    }
    // @ts-ignore
    tx.id = txHashId("Transaction", tx.hash);
    // @ts-ignore
    return tx;
  }

  @Query(() => Address)
  async address(
    @Args() request: AddressRequest,
    @Ctx() ctx: ResolverContext
  ): Promise<Address> {
    const address = await ctx.service.indexedChainService.getAddressByHash(
      request.hash
    );
    if (!address) {
      throw new ApolloError(
        `Address hash 0x${request.hash} was not found.`,
        "NOT_FOUND"
      );
    }

    const numTxs = await ctx.service.indexedChainService.getNumTxs({
      addressHash: request.hash,
    });

    return {
      ...address,
      hashQr: await QRCode.toDataURL(`0x${address.hash.toString("hex")}`),
      numTxs,
    };
  }
}
