export function getRanges(array: Array<number>): Array<[number, number]> {
  const ranges: Array<[number, number]> = [];
  let rstart: number;
  let rend: number;
  for (let i = 0; i < array.length; i += 1) {
    rstart = array[i];
    rend = rstart;
    while (array[i + 1] - array[i] === 1) {
      rend = array[i + 1]; // increment the index if the numbers sequential
      i += 1;
    }
    ranges.push([rstart, rend]);
  }
  return ranges;
}

export function chunkRanges(
  missingRanges: Array<[number, number]>,
  blocksBatchSize: number
): Array<[number, number]> {
  const chunkedRanges: Array<[number, number]> = [];
  for (const range of missingRanges) {
    if (range[1] - range[0] < blocksBatchSize) {
      chunkedRanges.push(range);
    } else {
      let start = range[0];
      let end = start + blocksBatchSize - 1;
      while (end <= range[1] || start <= range[1]) {
        chunkedRanges.push([start, Math.min(end, range[1])]);
        start = end + 1;
        end = start + blocksBatchSize - 1;
      }
    }
  }
  return chunkedRanges;
}
