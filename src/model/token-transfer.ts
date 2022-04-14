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

@Entity("tokenTransfer")
export class TokenTransfer {
  @Index()
  @PrimaryColumn("bytea")
  @ManyToOne(() => Transaction)
  transaction: Buffer;

  transactionHash: Buffer;

  @Column("integer")
  logIndex: number;

  @Index()
  @ManyToOne(() => Address)
  fromAddress: Buffer;

  fromAddressHash: Buffer;

  @Index()
  @ManyToOne(() => Address)
  toAddress: Buffer;

  toAddressHash: Buffer;

  @Column("numeric")
  amount: string;

  @Column("numeric", { precision: 78 })
  tokenId: string;

  @Index()
  @ManyToOne(() => Address)
  tokenContractAddress: Buffer;

  tokenContractAddressHash: Buffer;

  blockHash: Buffer;

  @PrimaryColumn("bytea")
  @ManyToOne(() => Block)
  block: Buffer;

  @Index()
  @Column("integer")
  blockNumber: number;

  @Column("numeric", { array: true })
  amounts: string[];

  @Column("numeric", { precision: 78, array: true })
  tokenIds: string[];

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
