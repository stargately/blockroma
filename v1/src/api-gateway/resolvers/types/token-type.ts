import { ArgsType, Field, ObjectType } from "type-graphql";
import { BufferScalar } from "@/api-gateway/resolvers/types/buffer-scalar";
import {
  ConnectionType,
  EdgeType,
} from "@/api-gateway/resolvers/types/connection-type";
import * as Relay from "graphql-relay";

@ObjectType()
export class Token {
  @Field(() => String)
  name: string;

  @Field(() => String)
  symbol?: string;

  @Field()
  totalSupply: string;

  @Field()
  decimals: string;

  @Field()
  type: string;

  @Field(() => BufferScalar)
  contractAddress: Buffer;

  @Field(() => Boolean, { nullable: true })
  skipMetadata?: boolean;
}

@ObjectType()
export class TokenEdge extends EdgeType("token", Token) {}

@ObjectType()
export class TokenConnection extends ConnectionType<TokenEdge>(
  "token",
  TokenEdge
) {}

@ArgsType()
export class TokensArgs implements Relay.ConnectionArguments {
  @Field(() => String, {
    nullable: true,
    description: "token symbol",
  })
  symbol?: string;

  @Field(() => String, {
    nullable: true,
    description: "Paginate before opaque cursor",
  })
  before?: Relay.ConnectionCursor;

  @Field(() => String, {
    nullable: true,
    description: "Paginate after opaque cursor",
  })
  after?: Relay.ConnectionCursor;

  @Field(() => Number, { nullable: true, description: "Paginate first" })
  first?: number;

  @Field(() => Number, { nullable: true, description: "Paginate last" })
  last?: number;
}
