import { formatDistanceToNowStrict } from "date-fns";
import React, { useState } from "react";

export function TickingTs({
  timestamp,
  className,
  inTile,
}: {
  timestamp: string;
  className?: string;
  inTile?: boolean;
}): JSX.Element {
  const fromNow = formatDistanceToNowStrict(new Date(timestamp), {
    addSuffix: true,
  });
  const [val, setVal] = useState(fromNow);
  setInterval(() => {
    const newNow = formatDistanceToNowStrict(new Date(timestamp), {
      addSuffix: true,
    });
    setVal(newNow);
  }, 1000);
  return (
    <span
      className={className}
      data-from-now={timestamp}
      in-tile={String(Boolean(inTile))}
    >
      {val}
    </span>
  );
}
