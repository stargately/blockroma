import { Field, ObjectType, Int, ID } from "type-graphql";
import {
  ConnectionType,
  EdgeType,
} from "@/api-gateway/resolvers/types/connection-type";
import { BufferScalar } from "@/api-gateway/resolvers/types/buffer-scalar";

@ObjectType({ description: "Models a token transfer." })
export class TokenTransfer {
  @Field(() => ID)
  id: string;

  @Field(() => BufferScalar)
  transactionHash?: Buffer;

  @Field(() => Int)
  logIndex?: number;

  @Field(() => BufferScalar)
  fromAddress?: Buffer;

  @Field(() => BufferScalar)
  toAddress?: Buffer;

  @Field(() => String)
  amount?: string;

  @Field(() => String)
  tokenId?: string;

  @Field(() => BufferScalar)
  tokenContractAddress?: Buffer;

  @Field(() => BufferScalar)
  block?: Buffer;

  @Field(() => Int)
  blockNumber?: number;

  @Field(() => [String])
  amounts?: string[];

  @Field(() => [String])
  tokenIds?: string[];

  @Field()
  createdAt?: Date;

  @Field()
  updatedAt?: Date;

  @Field(() => String)
  type?: string;
}

@ObjectType()
export class TokenTransferEdge extends EdgeType(
  "tokenTransfer",
  TokenTransfer
) {}

@ObjectType()
export class TokenTransferConnection extends ConnectionType<TokenTransferEdge>(
  "tokenTransfer",
  TokenTransferEdge
) {}
