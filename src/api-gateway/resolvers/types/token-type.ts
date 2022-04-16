import { Field, ObjectType } from "type-graphql";
import { BufferScalar } from "@/api-gateway/resolvers/types/buffer-scalar";

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
