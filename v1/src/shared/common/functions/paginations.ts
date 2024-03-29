export interface PageInfo {
  hasNextPage: boolean | null;
  endCursor: string | null;
  startCursor: string | null;
  hasPreviousPage: boolean | null;
}

export interface WithPageInfo {
  pageInfo: PageInfo | null;
}

export function paginationProcessTotalNumPage(
  withPageInfo: WithPageInfo | undefined | null
): number {
  if (!withPageInfo) {
    return 1;
  }
  const { pageInfo } = withPageInfo;
  const totalRecord =
    parseInt(pageInfo?.endCursor || "20", 10) +
    parseInt(pageInfo?.startCursor || "20", 10);
  return Math.round(totalRecord / 20) - 1;
}
