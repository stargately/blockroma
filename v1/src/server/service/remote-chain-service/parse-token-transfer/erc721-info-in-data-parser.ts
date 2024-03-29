import { Log } from "@ethersproject/abstract-provider";
import { hexToBuffer } from "@/server/service/remote-chain-service/hex-utils";
import {
  erc20AndErc721TokenTransferFirstTopic,
  ParsedTokenTransfer,
  Parser,
} from "@/server/service/remote-chain-service/parse-token-transfer/token-transfer-parser-utils";
import { ethers } from "ethers";

const decoder = ethers.utils.defaultAbiCoder;

const signatures = [
  // transferFrom(address,address,uint256)
  "0x23b872dd",
  // safeTransferFrom(address,address,uint256)
  "0x42842e0e",
  // safeTransferFrom(address,address,uint256,bytes)
  "0xb88d4fde",
];

export const erc721InfoInDataParser: Parser = {
  matchInput(input: string): boolean {
    return signatures.some((sig) => input.startsWith(sig));
  },

  matchLog(log: Log) {
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
    };
  },
};
