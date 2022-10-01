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

const parsers: Parser[] = [
  erc20TokenTransfer,
  erc721TopicAsAddressesParser,
  erc721InfoInDataParser,
  erc1155BatchTransfer,
  erc1155SingleTransfer,
];

export function matchTokenTransferInput(input: string): boolean {
  return parsers.some((parser) => parser.matchInput(input));
}

// https://github.com/blockscout/blockscout/blob/master/apps/indexer/lib/indexer/transform/token_transfers.ex
export function parseTokenTransfers(logs: Log[]): ParsedTokenTransfers {
  const tokenTransfers = [];

  for (const log of logs) {
    for (const parser of parsers) {
      if (parser.matchLog(log)) {
        try {
          const ret = parser.parse(log);
          if (ret) {
            tokenTransfers.push(ret.tokenTransfer);
          }
        } catch (err) {
          logger.error(`unknown token transfer, failed to parse logs: ${logs}`);
        }
      }
    }
  }

  return {
    tokenTransfers,
  };
}
