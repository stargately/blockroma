import { Log } from "@ethersproject/abstract-provider";
import { hexToBuffer } from "@/server/service/remote-chain-service/hex-utils";
import { BigNumber, ethers } from "ethers";
import {
  erc1155BatchTransferSignature,
  ParsedTokenTransfer,
  Parser,
  trancate,
} from "@/server/service/remote-chain-service/parse-token-transfer/token-transfer-parser-utils";

const decoder = ethers.utils.defaultAbiCoder;

const signatures = [
  // safeTransferFrom(address,address,uint256,uint256,bytes)
  "0xf242432a",
  // safeBatchTransferFrom(address,address,uint256[],uint256[],bytes)
  "0x2eb2c2d6",
];

export const erc1155BatchTransfer: Parser = {
  matchInput(input: string): boolean {
    return signatures.some((sig) => input.startsWith(sig));
  },

  matchLog(log: Log) {
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
    };
  },
};
