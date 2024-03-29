import { BigNumber } from "ethers";

export function hexToBuffer(hex: string): Buffer {
  return Buffer.from(hex.slice(2), "hex");
}

export function hexToNumber(hex: string): number {
  return BigNumber.from(hex).toNumber();
}

export function hexToDecimal(hex: string): string {
  return BigNumber.from(hex).toString();
}
