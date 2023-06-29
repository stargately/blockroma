import { MyServer } from "@/server/start-server";
import { Block } from "@/model/block";
import { Transaction } from "@/model/transaction";
import { Address } from "@/model/address";
import * as Relay from "graphql-relay";
import { Token } from "@/model/token";
import { TokenTransfer } from "@/model/token-transfer";

export class IndexedChainService {
  server: MyServer;

  constructor(server: MyServer) {
    this.server = server;
  }

  async numberToBlock(number: number): Promise<Block | null> {
    return this.server.gateways.dbCon
      .getRepository(Block)
      .findOneBy({ number });
  }

  async chainMeta(): Promise<{
    blockHeight: number;
    totalTransactions: number;
    totalAddresses: number;
  }> {
    const [blockHeight, totalTransactions, totalAddresses] = await Promise.all([
      this.server.gateways.dbCon
        .createQueryBuilder("block", "b")
        .select("MAX(b.number)", "max")
        .getRawOne(),

      this.server.gateways.dbCon.getRepository(Transaction).count(),

      this.server.gateways.dbCon.getRepository(Address).count(),
    ]);

    return {
      blockHeight: blockHeight.max,
      totalTransactions,
      totalAddresses,
    };
  }

  async insertBlocks(blocks: Array<Block>): Promise<void> {
    await this.server.gateways.dbCon
      .createQueryBuilder()
      .insert()
      .into(Block)
      .values(blocks)
      .orIgnore(true)
      .execute();
  }

  async insertTransactions(transactions: Array<Transaction>): Promise<void> {
    if (!transactions.length) {
      return;
    }
    await this.server.gateways.dbCon
      .createQueryBuilder()
      .insert()
      .into(Transaction)
      .values(transactions)
      .orIgnore(true)
      .execute();
  }

  async insertAddresses(addresses: Array<Address>): Promise<void> {
    await this.server.gateways.dbCon
      .createQueryBuilder()
      .insert()
      .into(Address)
      .values(addresses)
      .orIgnore(true)
      .execute();
  }

  async updateAddressForHigherBlock(address: Address): Promise<void> {
    await this.server.gateways.dbCon
      .createQueryBuilder()
      .update(Address)
      .set(address)
      .where("fetchedCoinBalanceBlockNumber < :bn and hash = :hash", {
        bn: address.fetchedCoinBalanceBlockNumber,
        hash: address.hash,
      })
      .execute();
  }

  async updateTxGasUsed(txHash: Buffer, gasUsed: string): Promise<void> {
    await this.server.gateways.dbCon
      .getRepository(Transaction)
      .update({ hash: txHash }, { gasUsed });
  }

  async getTransactionByHash(hash: Buffer): Promise<Transaction | null> {
    const txs = await this.getTransactions(
      { txHash: hash },
      { first: 1, after: "0" }
    );
    return txs.data[0];
  }

  async getAddressByHash(hash: Buffer): Promise<Address | null> {
    return this.server.gateways.dbCon
      .getRepository(Address)
      .findOneBy({ hash });
  }

  // async getBlocks(
  //   args: ConnectionArgs
  // ): Promise<WithPageInfo<Block>> {
  //   const query = this.server.gateways.dbCon
  //     .createQueryBuilder()
  //     .select()
  //     .from(Block, "b");
  //
  //   return {
  //     data: [], b
  //     pageInfo: emptyPage,
  //   };
  // }

  async getBlocks(
    orFilters: { minerHash?: Buffer },
    args: ConnectionArgs
  ): Promise<WithPageInfo<Block & { numTxs: number }>> {
    const query = this.server.gateways.dbCon
      .createQueryBuilder()
      .select([
        `*`,
        `(SELECT COUNT(*) FROM "transaction" WHERE "transaction"."blockNumber" = "b"."number") as "numTxs"`,
      ])
      .from(Block, "b");

    let count;
    if (orFilters.minerHash) {
      query.where("b.minerHash = :minerHash", {
        fromAddress: orFilters.minerHash,
      });
      count = await query.getCount();
    } else {
      // avoid full block table query
      const result = await this.server.gateways.dbCon.query(
        `SELECT n_live_tup AS estimate FROM pg_stat_all_tables WHERE relname = 'block';`
      );
      count = Number(result[0]?.estimate);
    }

    // Forward pagination
    if (args.first !== undefined) {
      const offset = Number(args.after ?? 10);
      const data = await query
        .orderBy({
          "b.number": "DESC",
        })
        .limit(args.first)
        .offset(offset)
        .execute();

      const next = offset + args.first;
      const prev = count - offset;

      return {
        data,
        pageInfo: {
          hasNextPage: next < count,
          startCursor: String(offset + args.first),
          hasPreviousPage: prev < count,
          endCursor: String(prev),
        },
      };
    }
    // Backward pagination
    if (args.last !== undefined) {
      const offset = Number(args.before ?? 0);
      const data = await query
        .orderBy({ "b.number": "ASC" })
        .limit(args.last)
        .offset(offset)
        .execute();
      const next = offset + args.last;
      const prev = count - offset;

      return {
        data,
        pageInfo: {
          hasNextPage: next < count,
          startCursor: String(offset + args.last),
          hasPreviousPage: prev < count,
          endCursor: String(prev),
        },
      };
    }

    return {
      data: [],
      pageInfo: emptyPage,
    };
  }

  async getNumTxs({
    addressHash,
    blockNumber,
  }: {
    addressHash?: Buffer;
    blockNumber?: number;
  }): Promise<number> {
    const query = this.server.gateways.dbCon
      .createQueryBuilder()
      .select()
      .from(Transaction, "tx");

    if (addressHash) {
      query
        .where("tx.fromAddressHash = :fromAddress", {
          fromAddress: addressHash,
        })
        .orWhere("tx.toAddressHash = :toAddress", {
          toAddress: addressHash,
        });
    } else {
      query.where("tx.blockNumber = :blockNumber", { blockNumber });
    }

    return query.getCount();
  }

  // TODO(dora): not sure how to achieve https://github.com/benjamin658/typeorm-cursor-pagination, it seems broken
  // use numbered pagination for now
  async getTransactions(
    orFilters: { addressHash?: Buffer; blockNumber?: number; txHash?: Buffer },
    args: ConnectionArgs
  ): Promise<WithPageInfo<Transaction>> {
    const query = this.server.gateways.dbCon
      .createQueryBuilder()
      .select()
      .from(Transaction, "tx");

    let count;
    if (orFilters.addressHash) {
      query
        .where("tx.fromAddressHash = :fromAddress", {
          fromAddress: orFilters.addressHash,
        })
        .orWhere("tx.toAddressHash = :toAddress", {
          toAddress: orFilters.addressHash,
        });
      count = await query.getCount();
    } else if (orFilters.blockNumber) {
      query.where("tx.blockNumber = :blockNumber", {
        blockNumber: orFilters.blockNumber,
      });
      count = await query.getCount();
    } else if (orFilters.txHash) {
      query.where("tx.hash = :hash", {
        hash: orFilters.txHash,
      });
      count = await query.getCount();
    } else {
      // avoid full transaction table query
      const result = await this.server.gateways.dbCon.query(
        `SELECT n_live_tup AS estimate FROM pg_stat_all_tables WHERE relname = 'transaction';`
      );
      count = Number(result[0]?.estimate);
    }

    // Forward pagination
    if (args.first !== undefined) {
      const offset = Number(args.after ?? 10);
      const data = await query
        .orderBy({
          "tx.blockNumber": "DESC",
        })
        .limit(args.first)
        .offset(offset)
        .execute();

      const next = offset + args.first;
      const prev = count - offset;

      return {
        data,
        pageInfo: {
          hasNextPage: next < count,
          startCursor: String(offset + args.first),
          hasPreviousPage: prev < count,
          endCursor: String(prev),
        },
      };
    }
    // Backward pagination
    if (args.last !== undefined) {
      const offset = Number(args.after ?? 0);
      const data = await query
        .orderBy({ "tx.blockNumber": "ASC" })
        .limit(args.last)
        .offset(offset)
        .execute();
      const next = offset + args.last;
      const prev = count - offset;
      return {
        data,
        pageInfo: {
          hasNextPage: next < count,
          startCursor: String(offset + args.last),
          hasPreviousPage: prev < count,
          endCursor: String(prev),
        },
      };
    }

    return {
      data: [],
      pageInfo: emptyPage,
    };
  }

  async getToken(
    tokenContractAddressHash: Buffer | undefined
  ): Promise<Token | null> {
    return this.server.gateways.dbCon
      .getRepository(Token)
      .findOneBy({ contractAddress: tokenContractAddressHash });
  }

  // TODO(dora): not sure how to achieve https://github.com/benjamin658/typeorm-cursor-pagination, it seems broken
  // use numbered pagination for now
  async getTokens(
    orFilters: { type?: string; symbol?: string },
    args: ConnectionArgs
  ): Promise<WithPageInfo<Token>> {
    const query = this.server.gateways.dbCon
      .createQueryBuilder()
      .select()
      .from(Token, "tk");

    let count;
    if (orFilters.type) {
      query.where("tk.type = :type", {
        type: orFilters.type,
      });
      count = await query.getCount();
    } else if (orFilters.symbol) {
      query.where("tk.symbol = :symbol", {
        symbol: orFilters.symbol,
      });
      count = await query.getCount();
    } else {
      // avoid full transaction table query
      const result = await this.server.gateways.dbCon.query(
        `SELECT n_live_tup AS estimate FROM pg_stat_all_tables WHERE relname = 'token';`
      );
      count = Number(result[0]?.estimate);
    }

    // Forward pagination
    if (args.first !== undefined) {
      const after = Number(args.after ?? 10);
      const data = await query
        .orderBy({
          "tk.updatedAt": "DESC",
        })
        .limit(args.first)
        .offset(after)
        .execute();

      const next = after + args.first;
      const prev = count - after;

      return {
        data,
        pageInfo: {
          hasNextPage: next < count,
          startCursor: String(after),
          hasPreviousPage: prev < count,
          endCursor: String(next),
        },
      };
    }
    // Backward pagination
    if (args.last !== undefined) {
      const before = Number(args.before ?? 0);
      const data = await query
        .orderBy({ "tk.updatedAt": "ASC" })
        .limit(args.last)
        .offset(before)
        .execute();
      const next = before + args.last;
      const prev = count - before;
      return {
        data,
        pageInfo: {
          hasNextPage: next < count,
          startCursor: String(before),
          hasPreviousPage: prev < count,
          endCursor: String(next),
        },
      };
    }

    return {
      data: [],
      pageInfo: emptyPage,
    };
  }

  async getTransferByTxHash(
    transactionHash: Buffer
  ): Promise<WithPageInfo<TokenTransfer>> {
    const query = await this.server.gateways.dbCon
      .createQueryBuilder()
      .select()
      .from(TokenTransfer, "txTransfer")
      .where("txTransfer.transactionHash = :transactionHash", {
        transactionHash,
      })
      .execute();

    return {
      data: query,
      pageInfo: emptyPage,
    };
  }
}

export const emptyPage = {
  hasNextPage: false,
  startCursor: "0",
  hasPreviousPage: false,
  endCursor: "0",
};

interface ConnectionArgs {
  before?: Relay.ConnectionCursor;
  after?: Relay.ConnectionCursor;
  first?: number;
  last?: number;
}

type WithPageInfo<T> = {
  data: Array<T>;
  pageInfo: PageInfo;
};

type PageInfo = {
  hasNextPage: boolean;
  startCursor: Relay.ConnectionCursor | null;

  hasPreviousPage: boolean;
  endCursor: Relay.ConnectionCursor | null;
};
