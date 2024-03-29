export const assetURL = (path: string): string => {
  if (path[0] !== "/") {
    return `/${path}`;
  }
  return path;
};
