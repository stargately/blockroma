export function getGasUsedPercent(
  gasUsed: string | null | undefined,
  gasLimit: string | null | undefined
): string {
  const result =
    (parseFloat(gasUsed ?? "0") / parseFloat(gasLimit ?? "0")) * 100;

  // Check if result is NaN, return "0.0" if true, otherwise return the toFixed(1) value.
  return Number.isNaN(result) ? "0.0" : result.toFixed(1);
}
