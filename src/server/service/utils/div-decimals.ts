import config from "config";
import { BigNumber } from "ethers";

export function divDecimals(
  num: string | number | null | undefined,
  decimals: number | string | null = config.get("chain.decimals")
): string {
  if (!num) {
    return "0";
  }
  if (!decimals) {
    return String(num);
  }
  return BigNumber.from(num).div(BigNumber.from(10).pow(decimals)).toString();
}
