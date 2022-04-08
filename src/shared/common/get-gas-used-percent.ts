export function getGasUsedPercent(
  gasUsed: string | null | undefined,
  gasLimit: string | null | undefined
): string {
  return (
    (parseFloat(gasUsed ?? "0") / parseFloat(gasLimit ?? "0")) *
    100
  ).toFixed(1);
}
