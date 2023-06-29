export function divDecimals(
  num?: string | number | null,
  decimals?: number | string | null
): string {
  if (!num) {
    return "0";
  }
  if (!decimals) {
    return String(num);
  }
  return (Number(num) / 10 ** Number(decimals))
    .toFixed(20)
    .replace(/\.?0*$/, "");
}
