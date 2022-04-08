import { GraphQLScalarType, Kind } from "graphql";

export const BufferScalar = new GraphQLScalarType({
  name: "Buffer",
  serialize(value: unknown): string {
    // check the type of received value
    if (!(value instanceof Buffer)) {
      throw new Error("BufferScalar can only serialize Buffer values");
    }
    return `0x${value.toString("hex")}`; // value sent to the client
  },

  parseValue(value: unknown): Buffer {
    // check the type of received value
    if (typeof value !== "string") {
      throw new Error("BufferScalar can only parse string values");
    }
    return Buffer.from(value.replace(/^0x/, ""), "hex"); // value from the client input variables
  },

  parseLiteral(ast): Buffer {
    // check the type of received value
    if (ast.kind !== Kind.STRING) {
      throw new Error("BufferScalar can only parse string values");
    }
    return Buffer.from(ast.value.replace(/^0x/, ""), "hex"); // value from the client query
  },
});
