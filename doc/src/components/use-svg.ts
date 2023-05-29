import { useEffect } from "react";

export const useSvg = () => {
  useEffect(() => {
    const SVGInject = require("@iconfu/svg-inject");
    SVGInject(document.querySelectorAll("[data-inject-svg]"));
  }, []);
};
