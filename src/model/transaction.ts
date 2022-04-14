import {
  Column,
  CreateDateColumn,
  Entity,
  Index,
  ManyToOne,
  PrimaryColumn,
  UpdateDateColumn,
} from "typeorm";
import { Block } from "@/model/block";
import { Address } from "@/model/address";

// eslint-disable-next-line no-shadow
export enum Status {
  ERROR = 0,
  OK = 1,
}

@Entity()
export class Transaction {
  @PrimaryColumn("bytea")
  hash: Buffer;

  @ManyToOne(() => Block, (block) => block.transactions)
  block: Buffer;

  // TODO(dora) create desc index
  // https://stackoverflow.com/questions/69850518/typeorm-index-creation
  @Index()
  @ManyToOne(() => Address)
  fromAddress: Buffer;

  fromAddressHash: Buffer;

  // TODO(dora) create desc index
  // https://stackoverflow.com/questions/69850518/typeorm-index-creation
  @Index()
  @ManyToOne(() => Address)
  toAddress: Buffer;

  toAddressHash: Buffer;

  @Index()
  @Column("integer")
  blockNumber: number;

  @Column("numeric", { precision: 100 })
  cumulativeGasUsed: string;

  @Column("timestamp")
  earliestProcessingStart: Date;

  @Column("varchar", { length: 255 })
  error: string;

  @Column("numeric", { precision: 100 })
  gas: string;

  @Column("numeric", { precision: 100 })
  gasPrice: string;

  @Column("numeric", { precision: 100, nullable: true })
  gasUsed?: string;

  @Column("integer")
  index: number;

  @Column("timestamp", { nullable: true })
  createdContractCodeIndexedAt?: Date;

  @Column("bytea")
  input: Buffer;

  @Column("integer")
  nonce: number;

  @Column("numeric", { precision: 100 })
  r: string;

  @Column("numeric", { precision: 100 })
  s: string;

  @Column("integer")
  status: Status;

  @Column("numeric", { precision: 100 })
  v: string;

  @Column("numeric", { precision: 100 })
  value: string;

  @Column("text", { nullable: true })
  revertReason?: string;

  @Column("numeric", { precision: 100, nullable: true })
  maxPriorityFeePerGas?: string;

  @Column("numeric", { precision: 100, nullable: true })
  maxFeePerGas?: string;

  // `type` - New transaction type identifier introduced in EIP 2718 (Berlin HF)
  @Column("integer", { nullable: true })
  type?: number;

  @Column("timestamp")
  timestamp: Date;

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
}
