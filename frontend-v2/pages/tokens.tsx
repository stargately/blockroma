import React from "react";
import type { GetStaticProps, InferGetStaticPropsType } from "next";
import { serverSideTranslations } from "next-i18next/serverSideTranslations";
import { RawNav } from "@/shared/home/components/raw-nav";
import { TokensContainer } from "@/shared/token-container/tokens-container";
import { SeoHead } from "@/shared/common/seo-head";

export default function TokensPage() {
  return (
    <div>
      <SeoHead config={{ title: "Tokens" }} />
      <RawNav />
      <TokensContainer />
    </div>
  );
}

// or getServerSideProps: GetServerSideProps<Props> = async ({ locale })
export const getStaticProps: GetStaticProps<{}> = async ({ locale }) => ({
  props: {
    ...(await serverSideTranslations(locale ?? "en", ["common"])),
  },
});
