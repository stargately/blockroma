import React, { FormEventHandler } from "react";

type Props = {
  position: "top" | "bottom";
  numPages: number;
  curPage: number;
  setCurPage: (curPage: number) => void;
};

export const Pagination: React.FC<Props> = ({
  curPage,
  numPages,
  setCurPage,
  position,
}) => {
  const prev = curPage - 1;
  const next = curPage + 1;

  const onSubmit: FormEventHandler<HTMLFormElement> = (event) => {
    event.preventDefault();
    // @ts-ignore
    const val = String(event?.target[0].value).trim();
    setCurPage(parseInt(val, 10) || 1);
  };

  return (
    <div className={`pagination-container mlm17 mrm18 position-${position}`}>
      <ul className="pagination align-end">
        {prev > 1 && (
          <li className={`page-item${curPage === 1 ? " active" : ""}`}>
            <a
              className="page-link page-link-light-hover"
              data-page-number="1"
              onClick={() => setCurPage(1)}
            >
              1
            </a>
          </li>
        )}

        {prev > 3 && (
          <li className="page-item disabled">
            <a
              className="page-link page-link-light-hover"
              data-page-number="..."
            >
              ...
            </a>
          </li>
        )}

        {prev >= 1 && (
          <li className="page-item">
            <a
              className="page-link page-link-light-hover"
              data-page-number={prev}
              onClick={() => setCurPage(prev)}
            >
              {prev}
            </a>
          </li>
        )}

        <li className="page-item active">
          <a
            className="page-link page-link-light-hover"
            data-page-number={curPage}
            onClick={() => setCurPage(curPage)}
          >
            {curPage}
          </a>
        </li>

        {next <= numPages && (
          <li className="page-item">
            <a
              className="page-link page-link-light-hover"
              data-page-number={next}
              onClick={() => setCurPage(next)}
            >
              {next}
            </a>
          </li>
        )}

        {next < numPages - 1 && (
          <li className="page-item disabled">
            <a
              className="page-link page-link-light-hover"
              data-page-number="..."
            >
              ...
            </a>
          </li>
        )}

        {next < numPages && (
          <li className="page-item">
            <a
              className="page-link page-link-light-hover"
              data-page-number="200"
              onClick={() => setCurPage(numPages)}
            >
              {numPages}
            </a>
          </li>
        )}
      </ul>
      <ul className="pagination fml5 go-to">
        <li className="page-link no-hover tb ml10">Go to</li>
        <li className="page-item">
          <form onSubmit={onSubmit}>
            <input
              className="page-number"
              id="page-number"
              type="text"
              size={3}
            />
            <input className="d-none" type="submit" />
          </form>
        </li>
      </ul>
    </div>
  );
};
