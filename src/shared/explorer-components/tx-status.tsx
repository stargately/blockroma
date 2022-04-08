import * as React from "react";

type Prosp = {
  status: "ERROR" | "OK"  | null | undefined;
};

export function TxStatus({ status }: Prosp): JSX.Element {
  if (status === "OK") {
    return (
      <span className="mr-4" data-transaction-status="Success">
        <i
          style={{ color: "#20b760" }}
          className="fa-regular fa-check-circle"
        />
        {" "}Success
      </span>
    );
  }

  return (
    <span className="mr-4" data-transaction-status="Error: Reverted">
      <i
        style={{ color: "#dc3545" }}
        className="fa-solid fa-exclamation-circle"
      />
      {" "}Error
    </span>
  );
}
