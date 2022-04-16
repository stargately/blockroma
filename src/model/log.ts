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

@Index(["address"])
@Index(["address", "transaction"])
@Index(["transaction", "index"])
@Index(["transaction", "block", "index"], { unique: true })
@Entity()
export class Log {
  @Column("bytea")
  data: Buffer;

  @PrimaryColumn("integer")
  index: number;

  @Index()
  @Column("varchar", { length: 255 })
  type: string;

  @Index()
  @Column("varchar", { length: 255 })
  firstTopic: string;

  @Index()
  @Column("varchar", { length: 255 })
  secondTopic: string;

  @Index()
  @Column("varchar", { length: 255 })
  thirdTopic: string;

  @Index()
  @Column("varchar", { length: 255 })
  fourthTopic: string;

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

  @ManyToOne(() => Address)
  @Column("bytea", { name: "addressHash" })
  address: Buffer;

  addressHash: Buffer;

  @Index()
  @ManyToOne(() => Transaction)
  @PrimaryColumn("bytea", { name: "transactionHash" })
  transaction: Buffer;

  transactionHash: Buffer;

  @PrimaryColumn("bytea", { name: "blockHash" })
  @ManyToOne(() => Block)
  block: Buffer;

  blockHash: Buffer;

  @Index()
  @Column("integer")
  blockNumber: number;
}
