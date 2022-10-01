import { Log } from "@ethersproject/abstract-provider";
import {
  hexToBuffer,
  hexToDecimal,
} from "@/server/service/remote-chain-service/hex-utils";
import {
  erc20AndErc721TokenTransferFirstTopic,
  ParsedTokenTransfer,
  Parser,
  trancate,
} from "@/server/service/remote-chain-service/parse-token-transfer/token-transfer-parser-utils";

// ERC-20 or for old ERC-721, ERC-1155 token versions
// transfer(address, uint256)
const transferFunctionSignature = "0xa9059cbb";
const unknownFunctionSignature = "0xf907fc5b";

export const erc20TokenTransfer: Parser = {
  matchInput(input: string): boolean {
    return (
      input.startsWith(transferFunctionSignature) ||
      input.startsWith(unknownFunctionSignature)
    );
  },

  matchLog(log: Log) {
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
    };
  },
};
