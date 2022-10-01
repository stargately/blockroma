import { Log } from "@ethersproject/abstract-provider";
import {
  hexToBuffer,
  hexToDecimal,
} from "@/server/service/remote-chain-service/hex-utils";
import {
  erc1155SingleTransferSignature,
  ParsedTokenTransfer,
  Parser,
  trancate,
} from "@/server/service/remote-chain-service/parse-token-transfer/token-transfer-parser-utils";
import { ethers } from "ethers";
import { erc1155BatchTransfer } from "@/server/service/remote-chain-service/parse-token-transfer/erc1155-batch-transfer";

const decoder = ethers.utils.defaultAbiCoder;

export const erc1155SingleTransfer: Parser = {
  matchInput(input: string): boolean {
    return erc1155BatchTransfer.matchInput(input);
  },

  matchLog(log: Log) {
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
    };
  },
};
