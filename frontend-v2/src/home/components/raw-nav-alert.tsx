import * as React from "react";

export function RawNavAlert(): JSX.Element {
  return (
    <div
      className="alert alert-warning text-center mb-0 p-3"
      data-selector="indexed-status"
    >
      BMO is joining the BMO ecosystem, and token holders can now swap BMO for
      STAKE on the BMO chain! More info and instructions{" "}
      <a href="https://www.poa.network">here</a>
    </div>
  );
}
