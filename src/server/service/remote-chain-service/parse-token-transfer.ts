import { TokenTransfer } from "@/model/token-transfer";
import { Log } from "@ethersproject/abstract-provider";
import {
  hexToBuffer,
  hexToDecimal,
} from "@/server/service/remote-chain-service/hex-utils";
import { BigNumber, ethers } from "ethers";
import { logger } from "onefx/lib/integrated-gateways/logger";

const decoder = ethers.utils.defaultAbiCoder;

const erc20AndErc721TokenTransferFirstTopic =
  "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef";
const erc1155SingleTransferSignature =
  "0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62";
const erc1155BatchTransferSignature =
  "0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb";

export const transferFunctionSignature = "0xa9059cbb";

export type RawToken = {
  contractAddress: Buffer;
  type: "ERC-721" | "ERC-20" | "ERC-1155";
};

export type ParsedTokenTransfers = {
  tokenTransfers: TokenTransfer[];
  tokens: RawToken[];
};

interface Parser {
  condition(log: Log): boolean;
  parse(log: Log): ParsedTokenTransfer | null;
}

const erc721TopicAsAddressesParser = {
  condition(log: Log) {
    return Boolean(
      log.topics[0] === erc20AndErc721TokenTransferFirstTopic &&
        log.topics[1] &&
        log.topics[2] &&
        log.topics[3]
    );
  },
  parse: function parseErc721TopicAsAddresses(log: Log): ParsedTokenTransfer {
    const [tokenId] = decoder.decode(["uint256"], log.topics[3]);
    return {
      tokenTransfer: {
        blockNumber: log.blockNumber,
        logIndex: log.logIndex,
        block: hexToBuffer(log.blockHash),
        fromAddress: Buffer.from(trancate(log.topics[1]), "hex"),
        toAddress: Buffer.from(trancate(log.topics[2]), "hex"),
        tokenContractAddress: hexToBuffer(log.address),
        tokenId: tokenId.toString() ?? "0",
        transaction: hexToBuffer(log.transactionHash),
        type: "ERC-721",
      },
      token: {
        contractAddress: hexToBuffer(log.address),
        type: "ERC-721",
      },
    };
  },
};

const erc721InfoInDataParser = {
  condition(log: Log) {
    return Boolean(
      log.topics[0] === erc20AndErc721TokenTransferFirstTopic &&
        !log.topics[1] &&
        !log.topics[2] &&
        !log.topics[3]
    );
  },
  parse: function parseErc721InfoInData(log: Log): ParsedTokenTransfer {
    const [fromAddressHash, toAddressHash, tokenId] = decoder.decode(
      ["address", "address", "uint256"],
      log.data
    );
    return {
      tokenTransfer: {
        blockNumber: log.blockNumber,
        logIndex: log.logIndex,
        block: hexToBuffer(log.blockHash),
        fromAddress: hexToBuffer(fromAddressHash),
        toAddress: hexToBuffer(toAddressHash),
        tokenContractAddress: hexToBuffer(log.address),
        tokenId: tokenId.toString() ?? "0",
        transaction: hexToBuffer(log.transactionHash),
        type: "ERC-721",
      },
      token: {
        contractAddress: hexToBuffer(log.address),
        type: "ERC-721",
      },
    };
  },
};

type ParsedTokenTransfer = {
  tokenTransfer: TokenTransfer & { type: string };
  token: RawToken;
};

function trancate(addr: string): string {
  return addr.replace(/^0x000000000000000000000000/, "");
}

const erc20TokenTransfer = {
  condition(log: Log) {
    return Boolean(
      log.topics[0] === erc20AndErc721TokenTransferFirstTopic &&
        log.topics[1] &&
        log.topics[2] &&
        !log.topics[3]
    );
  },
  parse: function parseErc20TokenTransfer(
    log: Log
  ): ParsedTokenTransfer | null {
    return {
      tokenTransfer: {
        amount: hexToDecimal(log.data),
        blockNumber: log.blockNumber,
        block: hexToBuffer(log.blockHash),
        logIndex: log.logIndex,
        fromAddress: Buffer.from(trancate(log.topics[1]), "hex"),
        toAddress: Buffer.from(trancate(log.topics[2]), "hex"),
        tokenContractAddress: hexToBuffer(log.address),
        transaction: hexToBuffer(log.transactionHash),
        tokenId: undefined,
        type: "ERC-20",
      },
      token: {
        contractAddress: hexToBuffer(log.address),
        type: "ERC-20",
      },
    };
  },
};

const erc1155BatchTransfer = {
  condition(log: Log) {
    return log.topics[0] === erc1155BatchTransferSignature;
  },
  parse: function parseErc1155BatchTransfer(
    log: Log
  ): ParsedTokenTransfer | null {
    const [tokenIds, values] = decoder.decode(
      ["uint256[]", "uint256[]"],
      log.data
    );
    return {
      tokenTransfer: {
        blockNumber: log.blockNumber,
        block: hexToBuffer(log.blockHash),
        logIndex: log.logIndex,
        fromAddress: Buffer.from(trancate(log.topics[1]), "hex"),
        toAddress: Buffer.from(trancate(log.topics[2]), "hex"),
        tokenContractAddress: hexToBuffer(log.address),
        transaction: hexToBuffer(log.transactionHash),
        tokenId: undefined,
        tokenIds: tokenIds?.map((id: BigNumber) => id.toString()),
        amounts: values?.map((val: BigNumber) => val.toString()),
        type: "ERC-1155",
      },
      token: {
        contractAddress: hexToBuffer(log.address),
        type: "ERC-1155",
      },
    };
  },
};

const erc1155SingleTransfer = {
  condition(log: Log) {
    return log.topics[0] === erc1155SingleTransferSignature;
  },
  parse: function parseErc1155SingleTransfer(
    log: Log
  ): ParsedTokenTransfer | null {
    const [tokenId, value] = decoder.decode(["uint256", "uint256"], log.data);
    return {
      tokenTransfer: {
        amount: hexToDecimal(value),
        blockNumber: log.blockNumber,
        block: hexToBuffer(log.blockHash),
        logIndex: log.logIndex,
        fromAddress: Buffer.from(trancate(log.topics[2]), "hex"),
        toAddress: Buffer.from(trancate(log.topics[3]), "hex"),
        tokenContractAddress: hexToBuffer(log.address),
        transaction: hexToBuffer(log.transactionHash),
        type: "ERC-1155",
        tokenId: tokenId.toString() ?? "0",
      },
      token: {
        contractAddress: hexToBuffer(log.address),
        type: "ERC-1155",
      },
    };
  },
};

const parsers: Parser[] = [
  erc20TokenTransfer,
  erc721TopicAsAddressesParser,
  erc721InfoInDataParser,
  erc1155BatchTransfer,
  erc1155SingleTransfer,
];

function getUniqueTokens(tokens: RawToken[]): RawToken[] {
  const dedupTokens: Record<string, RawToken> = {};
  for (const t of tokens) {
    dedupTokens[t.contractAddress.toString()] = t;
  }
  return Object.values(dedupTokens);
}

export function parseTokenTransfers(logs: Log[]): ParsedTokenTransfers {
  const tokens = [];
  const tokenTransfers = [];

  for (const log of logs) {
    for (const parser of parsers) {
      if (parser.condition(log)) {
        try {
          const ret = parser.parse(log);
          if (ret) {
            tokens.push(ret.token);
            tokenTransfers.push(ret.tokenTransfer);
          }
        } catch (err) {
          logger.error(`unknown token transfer, failed to parse logs: ${logs}`);
        }
      }
    }
  }

  const uniqueTokens = getUniqueTokens(tokens);

  return {
    tokens: uniqueTokens,
    tokenTransfers,
  };
}
