import { Helmet } from "onefx/lib/react-helmet";
import { assetURL } from "onefx/lib/asset-url";
import React from "react";

export function PageStyles({ stylePath }: { stylePath: string }): JSX.Element {
  return (
    <Helmet
      link={[
        // styles
        {
          rel: "stylesheet",
          type: "text/css",
          href: assetURL(stylePath),
        },
      ]}
    ></Helmet>
  );
}
