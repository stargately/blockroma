export function normalizeTokenValue(val?: string | null): string {
  if (!val) {
    return "0";
  }
  return (parseFloat(val) / 1000000000000000000).toString();
}
