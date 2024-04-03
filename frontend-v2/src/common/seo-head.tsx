import Head from "next/head";
import React from "react";
import { chainConfig } from "@/shared/common/use-chain-config";

export const defaultSeoConfig = {
  title: chainConfig.chainName,
  tagline: `${chainConfig.chainName} Blockchain Explorer`,
  previewImageUrl: "/favicon.svg",
  description:
    "Blockroma is a tool for inspecting and analyzing EVM based blockchains. Blockchain explorer for Ethereum Networks.",
  twitter: "@stargately",
};

export const SeoHead: React.FC<{
  config?: Partial<typeof defaultSeoConfig>;
}> = (config = { config: defaultSeoConfig }) => {
  const seoConfig = {
    ...defaultSeoConfig,
    ...config.config,
  };
  return (
    <Head>
      <title>{seoConfig.title}</title>
      <meta name="twitter:site" content={seoConfig.twitter} />
      <meta name="twitter:image" content={seoConfig.previewImageUrl} />
      <meta name="twitter:title" content={seoConfig.title} />
      <meta name="twitter:description" content={seoConfig.description} />
      <meta property="og:description" content={seoConfig.description} />
      <meta
        property="og:image"
        name="og:image"
        content={seoConfig.previewImageUrl}
      />
    </Head>
  );
};
