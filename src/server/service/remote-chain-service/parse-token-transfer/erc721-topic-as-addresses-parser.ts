import { Log } from "@ethersproject/abstract-provider";
import { hexToBuffer } from "@/server/service/remote-chain-service/hex-utils";
import {
  erc20AndErc721TokenTransferFirstTopic,
  ParsedTokenTransfer,
  Parser,
  trancate,
} from "@/server/service/remote-chain-service/parse-token-transfer/token-transfer-parser-utils";
import { ethers } from "ethers";
import { erc721InfoInDataParser } from "@/server/service/remote-chain-service/parse-token-transfer/erc721-info-in-data-parser";

const decoder = ethers.utils.defaultAbiCoder;

export const erc721TopicAsAddressesParser: Parser = {
  matchInput(input: string): boolean {
    return erc721InfoInDataParser.matchInput(input);
  },

  matchLog(log: Log) {
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
    };
  },
};
