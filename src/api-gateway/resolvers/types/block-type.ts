import { Field, Int, ObjectType } from "type-graphql";
import { Decimal } from "@/api-gateway/resolvers/types/decimal-scalar";
import { BufferScalar } from "@/api-gateway/resolvers/types/buffer-scalar";
import {
  ConnectionType,
  EdgeType,
} from "@/api-gateway/resolvers/types/connection-type";

@ObjectType({
  description:
    'A package of data that contains zero or more transactions, the hash of the previous block ("parent"), and optionally other data. Because each block (except for the initial "genesis block") points to the previous block, the data structure that they form is called a "blockchain".\n',
})
export class Block {
  @Field(() => Boolean)
  consensus: boolean;

  @Field(() => String)
  difficulty: Decimal;

  @Field(() => String)
  gasLimit: Decimal;

  @Field(() => String)
  gasUsed: Decimal;

  @Field(() => BufferScalar)
  hash: Buffer;

  @Field(() => BufferScalar)
  miner: Buffer;

  @Field(() => BufferScalar)
  nonce: Buffer;

  @Field(() => Int)
  number: number;

  @Field(() => BufferScalar)
  parentHash: Buffer;

  @Field(() => Int)
  size: number;

  @Field()
  timestamp: Date;

  @Field(() => String)
  totalDifficulty: Decimal;

  @Field(() => Int)
  numTxs: number;
}

@ObjectType()
export class BlockEdge extends EdgeType("block", Block) {}

@ObjectType()
export class BlockConnection extends ConnectionType<BlockEdge>(
  "block",
  BlockEdge
) {}

@ObjectType({ description: "A stored representation of a Web3 address." })
export class Address {
  @Field(() => String)
  fetchedCoinBalance?: Decimal;

  @Field(() => Int)
  fetchedCoinBalanceBlockNumber?: number;

  @Field(() => BufferScalar, { nullable: false })
  hash: Buffer;

  @Field(() => BlockConnection)
  blocks?: BlockConnection;
}
