import { TokenTransfer } from "@/model/token-transfer";
import { Log } from "@ethersproject/abstract-provider";

export type ParsedTokenTransfers = {
  tokenTransfers: TokenTransfer[];
};

export type ParsedTokenTransfer = {
  tokenTransfer: TokenTransfer & { type: string };
};

export const erc20AndErc721TokenTransferFirstTopic =
  "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef";
export const erc1155SingleTransferSignature =
  "0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62";
export const erc1155BatchTransferSignature =
  "0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb";

export function trancate(addr: string): string {
  return addr.replace(/^0x000000000000000000000000/, "");
}

export interface Parser {
  matchInput(input: string): boolean;
  matchLog(log: Log): boolean;
  parse(log: Log): ParsedTokenTransfer | null;
}
