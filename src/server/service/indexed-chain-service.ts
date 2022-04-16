import { MyServer } from "@/server/start-server";
import { Block } from "@/model/block";
import { Transaction } from "@/model/transaction";
import { Address } from "@/model/address";
import * as Relay from "graphql-relay";
import { Token } from "@/model/token";

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

    if (orFilters.minerHash) {
      query.where("b.minerHash = :minerHash", {
        fromAddress: orFilters.minerHash,
      });
    }

    const count = await query.getCount();

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

    if (orFilters.addressHash) {
      query
        .where("tx.fromAddressHash = :fromAddress", {
          fromAddress: orFilters.addressHash,
        })
        .orWhere("tx.toAddressHash = :toAddress", {
          toAddress: orFilters.addressHash,
        });
    } else if (orFilters.blockNumber) {
      query.where("tx.blockNumber = :blockNumber", {
        blockNumber: orFilters.blockNumber,
      });
    } else if (orFilters.txHash) {
      query.where("tx.hash = :hash", {
        hash: orFilters.txHash,
      });
    }

    const count = await query.getCount();

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
