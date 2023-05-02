import { Status } from "@/model/transaction";
import { Field, ID, Int, ObjectType, registerEnumType } from "type-graphql";
import { BufferScalar } from "@/api-gateway/resolvers/types/buffer-scalar";
import {
  ConnectionType,
  EdgeType,
} from "@/api-gateway/resolvers/types/connection-type";
import { Decimal } from "./decimal-scalar";

registerEnumType(Status, {
  name: "Status",
});

@ObjectType({ description: "Models a Web3 transaction." })
export class Transaction {
  @Field(() => Int)
  blockNumber?: number;

  @Field(() => String)
  createdContractAddressHash?: string;

  @Field(() => String)
  cumulativeGasUsed?: Decimal;

  @Field(() => String)
  error?: string;

  @Field(() => BufferScalar)
  fromAddressHash?: Buffer;

  @Field(() => String)
  gas?: Decimal;

  @Field(() => String)
  gasPrice?: Decimal;

  @Field(() => String)
  gasUsed?: Decimal;

  @Field(() => BufferScalar)
  hash?: Buffer;

  @Field(() => ID, { description: "The ID of an object", nullable: false })
  id: string;

  @Field(() => Int)
  index?: number;

  @Field(() => BufferScalar)
  input?: Buffer;

  @Field()
  timestamp: Date;

  // internalTransactions(
  //   after: String
  // before: String;
  // count: Int;
  // first: Int;
  // last: Int): InternalTransactionConnection;
  //

  @Field(() => Int)
  nonce?: number;

  @Field(() => Int)
  type?: number;

  @Field(() => String)
  r?: string;

  @Field(() => String)
  s?: string;

  @Field(() => Status)
  status?: Status;

  @Field(() => BufferScalar)
  toAddressHash?: Buffer;

  @Field(() => String)
  v?: string;

  @Field(() => String)
  value?: string;

  @Field(() => String)
  revertReason?: string;

  @Field(() => String)
  maxPriorityFeePerGas?: Decimal;

  @Field(() => String)
  maxFeePerGas?: Decimal;
}

@ObjectType()
export class TransactionEdge extends EdgeType("transaction", Transaction) {}

@ObjectType()
export class TransactionConnection extends ConnectionType<TransactionEdge>(
  "transaction",
  TransactionEdge
) {}
