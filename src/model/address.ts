import {
  Column,
  CreateDateColumn,
  Entity,
  PrimaryColumn,
  UpdateDateColumn,
} from "typeorm";

@Entity()
export class Address {
  @PrimaryColumn("bytea")
  hash: Buffer;

  @Column("numeric", { precision: 100 })
  fetchedCoinBalance: string;

  @Column("integer")
  fetchedCoinBalanceBlockNumber: number;

  @Column("bytea", { nullable: true })
  contractCode?: Buffer;

  @Column("integer", { nullable: true })
  nonce?: number;

  @Column("boolean", { default: false })
  decompiled: boolean;

  @Column("boolean", { default: false })
  verified: boolean;

  @Column("integer", { nullable: true })
  gasUsed?: number;

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
