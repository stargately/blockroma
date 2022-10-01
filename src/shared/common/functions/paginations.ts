export function paginationProcessTotalNumPage(data: object | undefined, dataKey: string): number {
  if (!data) {return 1}
  // @ts-ignore
  const pageInfo = data[dataKey]?.pageInfo;
  const totalRecord = parseInt(pageInfo?.endCursor || '20', 10) + parseInt(pageInfo?.startCursor || '20', 10)
  // @ts-ignore
  // eslint-disable-next-line no-bitwise
  return Math.round(totalRecord / 20) - 1
}
