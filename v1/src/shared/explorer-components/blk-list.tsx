import React from "react";
import { Blk, BlkBlockItem } from "@/shared/explorer-components/blk-block-item";

type Props = {
  blks?: Array<Blk | null | undefined> | null | undefined;
};

export const BlkList: React.FC<Props> = ({ blks }) => {
  if (!blks) {
    return <></>;
  }

  return (
    <div>
      <div data-items="true">
        {blks.map((blk) => blk && <BlkBlockItem key={blk.hash} blk={blk} />)}
      </div>
    </div>
  );
};
