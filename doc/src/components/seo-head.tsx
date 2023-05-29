import React from "react";
import Head from "@docusaurus/Head";
import { beConfig } from "@site/src/config";

export const SeoHead: React.FC = () => {
  return (
    <Head>
      <meta name="twitter:site" content="@BlockEdenHQ" />
      <meta name="twitter:image" content={beConfig.previewImageUrl} />
      <meta name="twitter:title" content={beConfig.title} />
      <meta name="twitter:description" content={beConfig.description} />
      <meta property="og:description" content={beConfig.description} />
      <meta
        property="og:image"
        name="og:image"
        content={beConfig.previewImageUrl}
      />

      <link
        href="assets/css/loaders/loader-typing.css"
        rel="stylesheet"
        type="text/css"
        media="all"
      />
      <link
        rel="preload"
        as="font"
        href="assets/fonts/Inter-UI-upright.var.woff2"
        type="font/woff2"
        crossOrigin="anonymous"
      />
      <link
        rel="preload"
        as="font"
        href="assets/fonts/Inter-UI.var.woff2"
        type="font/woff2"
        crossOrigin="anonymous"
      />
    </Head>
  );
};
