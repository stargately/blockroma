export function shortenHash(addr: string): string {
  try {
    return `${addr.slice(0, 8)}-${addr.slice(addr.length - 6, addr.length)}`;
  } catch (err) {
    return "";
  }
}
