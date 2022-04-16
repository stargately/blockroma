import {
  Column,
  CreateDateColumn,
  Entity,
  Index,
  ManyToOne,
  PrimaryColumn,
  UpdateDateColumn,
} from "typeorm";
import { Address } from "@/model/address";

/**
 * Represents a token.
 *
 *   ## Token Indexing
 *
 *   The following types of tokens are indexed:
 *
 *   * ERC-20
 *   * ERC-721
 *   * ERC-1155
 *
 *   ## Token Specifications
 *
 *   * [ERC-20](https://github.com/ethereum/EIPs/blob/master/EIPS/eip-20.md)
 *   * [ERC-721](https://github.com/ethereum/EIPs/blob/master/EIPS/eip-721.md)
 *   * [ERC-777](https://github.com/ethereum/EIPs/blob/master/EIPS/eip-777.md)
 *   * [ERC-1155](https://github.com/ethereum/EIPs/blob/master/EIPS/eip-1155.md)
 */
@Entity()
export class Token {
  @Column("varchar", { length: 255 })
  name: string;

  @Index("symbolIndex")
  @Column("varchar", { length: 255 })
  symbol: string;

  @Column("numeric")
  totalSupply: string;

  @Column("numeric")
  decimals: string;

  @Index()
  @Column("varchar", { length: 255 })
  type: string;

  @Index({ unique: true })
  @ManyToOne(() => Address)
  @PrimaryColumn("bytea", { unique: true, name: "contractAddressHash" })
  contractAddress: Buffer;

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

  // @Column("integer")
  // holderCount: number;
  //
  // @Column({nullable: true})
  // bridged?: boolean;
  //
  // @Column()
  // cataloged: boolean;

  @Column({ nullable: true })
  skipMetadata?: boolean;
}
