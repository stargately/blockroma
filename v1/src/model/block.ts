import {
  Column,
  CreateDateColumn,
  Entity,
  Index,
  ManyToOne,
  OneToMany,
  PrimaryColumn,
  UpdateDateColumn,
} from "typeorm";
import { Transaction } from "@/model/transaction";
import { Address } from "@/model/address";

@Entity()
export class Block {
  @PrimaryColumn("bytea")
  hash: Buffer;

  // `consensus`
  //   * `true` - this is a block on the longest consensus agreed upon chain.
  //   * `false` - this is an uncle block from a fork.
  @Column()
  consensus: boolean;

  @Column("numeric", { precision: 100 })
  difficulty: string;

  @Column("numeric", { precision: 100 })
  gasLimit: string;

  @Column("numeric", { precision: 100 })
  gasUsed: string;

  @Column("bytea")
  nonce: Buffer;

  @Index("blockNumberIndex", { unique: true })
  @Column("integer")
  number: number;

  // `size` - The size of the block in bytes.
  @Column("integer")
  size: number;

  @Column("timestamp")
  timestamp: Date;

  @Column("numeric", { precision: 100 })
  totalDifficulty: string;

  @Column("numeric", { precision: 100 })
  baseFeePerGas: string;

  @Column()
  isEmpty: boolean;

  @OneToMany(() => Transaction, (transaction) => transaction.block)
  transactions?: Transaction[];

  @Column("bytea")
  parentHash: Buffer;

  @ManyToOne(() => Address)
  @Column("bytea")
  miner: Buffer;

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
