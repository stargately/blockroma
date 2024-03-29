import {
  Column,
  CreateDateColumn,
  Entity,
  Index,
  ManyToOne,
  PrimaryColumn,
  UpdateDateColumn,
} from "typeorm";
import { Transaction } from "@/model/transaction";
import { Address } from "@/model/address";
import { Block } from "@/model/block";

/**
 * TokenTransfer represents a token transfer between addresses for a given token.
 *
 *   ## Overview
 *
 *   Token transfers are special cases from a `t:Explorer.Chain.Log.t/0`. A token
 *   transfer is always signified by the value from the `first_topic` in a log. That value
 *   is always `0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef`.
 *
 *   ## Data Mapping From a Log
 *
 *   Here's how a log's data maps to a token transfer:
 *
 *   | Log                | Token Transfer               | Description                     |
 *   |--------------------|------------------------------|---------------------------------|
 *   | `:secondTopic`     | `:fromAddressHash`           | Address sending tokens          |
 *   | `:thirdTopic`      | `:toAddressHash`             | Address receiving tokens        |
 *   | `:data`            | `:amount`                    | Amount of tokens transferred    |
 *   | `:transactionHash` | `:transactionHash`           | Transaction of the transfer     |
 *   | `:addressHash`     | `:tokenContractAddressHash`  | Address of token's contract     |
 *   | `:index`           | `:logIndex`                  | Index of log in transaction     |
 */
@Entity("tokenTransfer")
export class TokenTransfer {
  @Index()
  @PrimaryColumn("bytea", { name: "transactionHash" })
  @ManyToOne(() => Transaction)
  transaction: Buffer;

  transactionHash?: Buffer;

  @Column("integer")
  logIndex: number;

  @Index()
  @ManyToOne(() => Address, {})
  fromAddress: Buffer;

  fromAddressHash?: Buffer;

  @Index()
  @ManyToOne(() => Address)
  toAddress: Buffer;

  toAddressHash?: Buffer;

  @Column("numeric", { nullable: true })
  amount?: string;

  @Column("numeric", { precision: 78, nullable: true })
  tokenId?: string;

  @Index()
  @ManyToOne(() => Address)
  tokenContractAddress: Buffer;

  tokenContractAddressHash?: Buffer;

  @PrimaryColumn("bytea", { name: "blockHash" })
  @ManyToOne(() => Block)
  block: Buffer;

  blockHash?: Buffer;

  @Index()
  @Column("integer")
  blockNumber: number;

  @Column("numeric", { array: true, nullable: true })
  amounts?: string[];

  @Column("numeric", { precision: 78, array: true, nullable: true })
  tokenIds?: string[];

  @CreateDateColumn({
    type: "timestamp",
    default: () => "CURRENT_TIMESTAMP(6)",
  })
  createdAt?: Date;

  @UpdateDateColumn({
    type: "timestamp",
    default: () => "CURRENT_TIMESTAMP(6)",
    onUpdate: "CURRENT_TIMESTAMP(6)",
  })
  updatedAt?: Date;

  type?: string;
}
