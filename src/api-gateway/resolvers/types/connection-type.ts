import * as Relay from "graphql-relay";
import { ObjectType, Field, ArgsType, ClassType } from "type-graphql";

@ObjectType()
class PageInfo implements Relay.PageInfo {
  @Field(() => Boolean)
  hasNextPage: boolean;

  @Field(() => Boolean)
  hasPreviousPage: boolean;

  @Field(() => String, { nullable: true })
  startCursor: Relay.ConnectionCursor | null;

  @Field(() => String, { nullable: true })
  endCursor: Relay.ConnectionCursor | null;
}

@ArgsType()
export class ConnectionArgs implements Relay.ConnectionArguments {
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

export function EdgeType<NodeType>(
  nodeName: string,
  nodeType: ClassType<NodeType>
) {
  @ObjectType(`${nodeName}Edge`, { isAbstract: true })
  abstract class Edge implements Relay.Edge<NodeType> {
    @Field(() => nodeType)
    node: NodeType;

    @Field(() => String, {
      description: "Used in `before` and `after` args",
    })
    cursor: Relay.ConnectionCursor;
  }

  return Edge;
}

type ExtractNodeType<TEdgeType> = TEdgeType extends Relay.Edge<infer NodeType>
  ? NodeType
  : never;

export function ConnectionType<
  TEdgeType extends Relay.Edge<NodeType>,
  NodeType = ExtractNodeType<TEdgeType>
>(nodeName: string, edgeClass: ClassType<TEdgeType>) {
  @ObjectType(`${nodeName}Connection`, { isAbstract: true })
  abstract class Connection implements Relay.Connection<NodeType> {
    @Field(() => PageInfo)
    pageInfo: PageInfo;

    @Field(() => [edgeClass])
    edges: TEdgeType[];
  }

  return Connection;
}
