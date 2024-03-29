import * as React from "react";
import { createRef, useEffect, useState } from "react";

type Props = {
  value: string;
  reason: string;
};

export function CopyToClipboard({ value, reason }: Props): JSX.Element {
  const [copied, setCopied] = useState(false);
  const text = copied ? "Copied!" : reason;
  const ref = createRef<HTMLSpanElement>();
  useEffect(() => {
    ref.current?.addEventListener("mouseleave", () => {
      setCopied(false);
    });
    ref.current?.addEventListener("touchcancel", () => {
      setCopied(false);
    });
  }, [ref]);

  return (
    <span
      aria-label={text}
      data-clipboard-text={value}
      className="btn-copy-icon btn-copy-icon-small btn-copy-icon-custom btn-copy-icon-no-borders i-tooltip-2"
      data-placement="top"
      data-toggle="tooltip"
      title={text}
      style={{}}
      ref={ref}
      onClick={() => {
        setCopied(true);

        navigator.clipboard.writeText(value);
      }}
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 32.5 32.5"
        width={32}
        height={32}
      >
        <path
          fillRule="evenodd"
          d="M23.5 20.5a1 1 0 0 1-1-1v-9h-9a1 1 0 0 1 0-2h10a1 1 0 0 1 1 1v10a1 1 0 0 1-1 1zm-3-7v10a1 1 0 0 1-1 1h-10a1 1 0 0 1-1-1v-10a1 1 0 0 1 1-1h10a1 1 0 0 1 1 1zm-2 1h-8v8h8v-8z"
        />
      </svg>
    </span>
  );
}
