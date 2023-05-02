import { Field, Int, ObjectType } from "type-graphql";
import { TransactionConnection } from "@/api-gateway/resolvers/types/transaction-type";
import { BufferScalar } from "@/api-gateway/resolvers/types/buffer-scalar";
import { Decimal } from "./decimal-scalar";

@ObjectType({ description: "A stored representation of a Web3 address." })
export class Address {
  @Field(() => String)
  fetchedCoinBalance?: Decimal;

  @Field(() => Int)
  fetchedCoinBalanceBlockNumber?: number;

  @Field()
  nonce?: number;

  @Field()
  gasUsed?: number;

  @Field()
  hashQr?: string;

  @Field()
  numTxs?: number;

  @Field(() => BufferScalar, { nullable: false })
  hash: Buffer;

  @Field(() => TransactionConnection)
  transactions?: TransactionConnection;
}
