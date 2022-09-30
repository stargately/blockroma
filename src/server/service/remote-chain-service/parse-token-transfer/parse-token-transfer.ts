import { Log } from "@ethersproject/abstract-provider";
import { logger } from "onefx/lib/integrated-gateways/logger";
import { erc20TokenTransfer } from "@/server/service/remote-chain-service/parse-token-transfer/erc20-token-transfer";
import { erc721TopicAsAddressesParser } from "@/server/service/remote-chain-service/parse-token-transfer/erc721-topic-as-addresses-parser";
import { erc721InfoInDataParser } from "@/server/service/remote-chain-service/parse-token-transfer/erc721-info-in-data-parser";
import { erc1155BatchTransfer } from "@/server/service/remote-chain-service/parse-token-transfer/erc1155-batch-transfer";
import { erc1155SingleTransfer } from "@/server/service/remote-chain-service/parse-token-transfer/erc1155-single-transfer";
import {
  ParsedTokenTransfers,
  Parser,
} from "@/server/service/remote-chain-service/parse-token-transfer/token-transfer-parser-utils";

export type RawToken = {
  contractAddress: Buffer;
  type: "ERC-721" | "ERC-20" | "ERC-1155";
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
    dedupTokens[t.contractAddress.toString("hex")] = t;
  }
  return Object.values(dedupTokens);
}

export function matchTokenTransferInput(input: string): boolean {
  return parsers.some((parser) => parser.matchInput(input));
}

// https://github.com/blockscout/blockscout/blob/master/apps/indexer/lib/indexer/transform/token_transfers.ex
export function parseTokenTransfers(logs: Log[]): ParsedTokenTransfers {
  const tokens = [];
  const tokenTransfers = [];

  for (const log of logs) {
    for (const parser of parsers) {
      if (parser.matchLog(log)) {
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
