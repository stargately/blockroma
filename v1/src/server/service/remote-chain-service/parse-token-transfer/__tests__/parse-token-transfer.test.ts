import test from "ava";
import { parseTokenTransfers } from "@/server/service/remote-chain-service/parse-token-transfer/parse-token-transfer";
import { Log } from "@ethersproject/abstract-provider";
import { hexToBuffer } from "@/server/service/remote-chain-service/hex-utils";
import { ParsedTokenTransfers } from "@/server/service/remote-chain-service/parse-token-transfer/token-transfer-parser-utils";

test("parse/1 parses logs for tokens and token transfers", async (t) => {
  const parsed = parseTokenTransfers([
    {
      transactionIndex: 0,
      blockNumber: 3530917,
      transactionHash:
        "0x43dfd761974e8c3351d285ab65bee311454eb45b149a015fe7804a33252f19e5",
      address: "0xf2eec76e45b328df99a34fa696320a262cb92154",
      topics: [
        "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
        "0x000000000000000000000000556813d9cc20acfe8388af029a679d34a63388db",
        "0x00000000000000000000000092148dd870fa1b7c4700f2bd7f44238821c26f73",
      ],
      data: "0x000000000000000000000000000000000000000000000000ebec21ee1da40000",
      logIndex: 8,
      blockHash:
        "0x79594150677f083756a37eee7b97ed99ab071f502104332cb3835bac345711ca",
      removed: false,
    },
    {
      transactionIndex: 0,
      blockNumber: 3586935,
      transactionHash:
        "0x8425a9b81a9bd1c64861110c1a453b84719cb0361d6fa0db68abf7611b9a890e",
      address: "0x6ea5ec9cb832e60b6b1654f5826e9be638f276a5",
      topics: [
        "0x55e10366a5f552746106978b694d7ef3bbddec06bd5f9b9d15ad46f475c653ef",
        "0x00000000000000000000000063b0595bb7a0b7edd0549c9557a0c8aee6da667b",
        "0x000000000000000000000000f3089e15d0c23c181d7f98b0878b560bfe193a1d",
        "0xc89efdaa54c0f20c7adf612882df0950f5a951637e0307cdcb4c672f298b8bc6",
      ],
      data: "0x",
      logIndex: 0,
      blockHash:
        "0x79594150677f083756a37eee7b97ed99ab071f502104332cb3835bac345711ca",
      removed: false,
    },
    {
      transactionIndex: 0,
      blockNumber: 3664064,
      transactionHash:
        "0x4011d9a930a3da620321589a54dc0ca3b88216b4886c7a7c3aaad1fb17702d35",
      address: "0x91932e8c6776fb2b04abb71874a7988747728bb2",
      topics: [
        "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
        "0x0000000000000000000000009851ba177554eb07271ac230a137551e6dd0aa84",
        "0x000000000000000000000000dccb72afee70e60b0c1226288fe86c01b953e8ac",
        "0x00000000000000000000000000000000000000000000000000000000000000b7",
      ],
      data: "0x000000000000000000000000000000000000000000000000ebec21ee1da40000",
      logIndex: 1,
      blockHash:
        "0x79594150677f083756a37eee7b97ed99ab071f502104332cb3835bac345711ca",
      removed: false,
    },
  ]);
  t.deepEqual(parsed.tokenTransfers.length, 2);
  t.truthy(parsed.tokenTransfers.filter((it) => it.type === "ERC-20").length);
  t.truthy(parsed.tokenTransfers.filter((it) => it.type === "ERC-721").length);
});

test("parses ERC-721 transfer with addresses in data field", async (t) => {
  const log: Log = {
    address: "0x58Ab73CB79c8275628E0213742a85B163fE0A9Fb",
    blockNumber: 8683457,
    data: "0x00000000000000000000000058ab73cb79c8275628e0213742a85b163fe0a9fb000000000000000000000000be8cdfc13ffda20c844ac3da2b53a23ac5787f1e0000000000000000000000000000000000000000000000000000000000003a5b",
    topics: [
      "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
    ],
    logIndex: 2,
    transactionHash:
      "0x6d2dd62c178e55a13b65601f227c4ffdd8aa4e3bcb1f24731363b4f7619e92c8",
    blockHash:
      "0x79594150677f083756a37eee7b97ed99ab071f502104332cb3835bac345711ca",
    removed: false,
    transactionIndex: 1,
  };

  const expected: ParsedTokenTransfers = {
    tokenTransfers: [
      {
        blockNumber: log.blockNumber,
        logIndex: log.logIndex,
        fromAddress: Buffer.from(
          "58ab73cb79c8275628e0213742a85b163fe0a9fb",
          "hex"
        ),
        toAddress: Buffer.from(
          "be8cdfc13ffda20c844ac3da2b53a23ac5787f1e",
          "hex"
        ),
        block: Buffer.from(
          "79594150677f083756a37eee7b97ed99ab071f502104332cb3835bac345711ca",
          "hex"
        ),
        tokenContractAddress: hexToBuffer(log.address),
        tokenId: "14939",
        transaction: hexToBuffer(log.transactionHash),
        type: "ERC-721",
      },
    ],
  };

  t.deepEqual(parseTokenTransfers([log]), expected);
});

test("parses erc1155 token transfer", (t) => {
  const log: Log = {
    address: "0x58Ab73CB79c8275628E0213742a85B163fE0A9Fb",
    blockNumber: 8683457,
    data: "0x1000000000000c520000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001",
    topics: [
      "0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62",
      "0x0000000000000000000000009c978f4cfa1fe13406bcc05baf26a35716f881dd",
      "0x0000000000000000000000009c978f4cfa1fe13406bcc05baf26a35716f881dd",
      "0x0000000000000000000000009c978f4cfa1fe13406bcc05baf26a35716f881dd",
    ],
    logIndex: 2,
    transactionIndex: 333,
    transactionHash:
      "0x6d2dd62c178e55a13b65601f227c4ffdd8aa4e3bcb1f24731363b4f7619e92c8",
    blockHash:
      "0x79594150677f083756a37eee7b97ed99ab071f502104332cb3835bac345711ca",
    removed: false,
  };
  const expected: ParsedTokenTransfers = {
    tokenTransfers: [
      {
        amount: "1",
        block: hexToBuffer(
          "0x79594150677f083756a37eee7b97ed99ab071f502104332cb3835bac345711ca"
        ),
        blockNumber: 8683457,
        fromAddress: hexToBuffer("0x9c978f4cfa1fe13406bcc05baf26a35716f881dd"),
        logIndex: 2,
        toAddress: hexToBuffer("0x9c978f4cfa1fe13406bcc05baf26a35716f881dd"),
        tokenContractAddress: hexToBuffer(
          "0x58Ab73CB79c8275628E0213742a85B163fE0A9Fb"
        ),
        tokenId:
          "7237005577332282011952059972634123378909214838582411639295170840059424276480",
        type: "ERC-1155",
        transaction: hexToBuffer(
          "0x6d2dd62c178e55a13b65601f227c4ffdd8aa4e3bcb1f24731363b4f7619e92c8"
        ),
      },
    ],
  };
  t.deepEqual(parseTokenTransfers([log]), expected);
});

test("parses erc1155 batch token transfer", async (t) => {
  const log: Log = {
    address: "0x58Ab73CB79c8275628E0213742a85B163fE0A9Fb",
    blockNumber: 8683457,
    data: "0x000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000001388",
    topics: [
      "0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb",
      "0x0000000000000000000000006c943470780461b00783ad530a53913bd2c104d3",
      "0x0000000000000000000000006c943470780461b00783ad530a53913bd2c104d3",
      "0x0000000000000000000000006c943470780461b00783ad530a53913bd2c104d3",
    ],
    logIndex: 2,
    transactionHash:
      "0x6d2dd62c178e55a13b65601f227c4ffdd8aa4e3bcb1f24731363b4f7619e92c8",
    blockHash:
      "0x79594150677f083756a37eee7b97ed99ab071f502104332cb3835bac345711ca",
    removed: false,
    transactionIndex: 999,
  };
  const expected: ParsedTokenTransfers = {
    tokenTransfers: [
      {
        block: hexToBuffer(
          "0x79594150677f083756a37eee7b97ed99ab071f502104332cb3835bac345711ca"
        ),
        blockNumber: 8683457,
        fromAddress: hexToBuffer("0x6c943470780461b00783ad530a53913bd2c104d3"),
        logIndex: 2,
        toAddress: hexToBuffer("0x6c943470780461b00783ad530a53913bd2c104d3"),
        tokenContractAddress: hexToBuffer(
          "0x58Ab73CB79c8275628E0213742a85B163fE0A9Fb"
        ),
        tokenId: undefined,
        tokenIds: ["680564733841876926926749214863536422912"],
        type: "ERC-1155",
        transaction: hexToBuffer(
          "0x6d2dd62c178e55a13b65601f227c4ffdd8aa4e3bcb1f24731363b4f7619e92c8"
        ),
        amounts: ["5000"],
      },
    ],
  };
  t.deepEqual(parseTokenTransfers([log]), expected);
});

test("logs error with unrecognized token transfer format", async (t) => {
  const log: Log = {
    address: "0x58Ab73CB79c8275628E0213742a85B163fE0A9Fb",
    blockNumber: 8_683_457,
    blockHash:
      "0x79594150677f083756a37eee7b97ed99ab071f502104332cb3835bac345711ca",
    data: "0x",
    topics: [
      "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
    ],
    logIndex: 2,
    transactionHash:
      "0x6d2dd62c178e55a13b65601f227c4ffdd8aa4e3bcb1f24731363b4f7619e92c8",
    removed: false,
    transactionIndex: 999,
  };
  t.deepEqual(parseTokenTransfers([log]), {
    tokenTransfers: [],
  });
});
