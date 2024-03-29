import { BlockDetailsContainer } from "@/shared/block-details-container/block-details-container";
import { GetStaticProps } from "next";
import { serverSideTranslations } from "next-i18next/serverSideTranslations";
import { RawNav } from "@/shared/home/components/raw-nav";
import React from "react";
import nextI18NextConfig from "../../next-i18next.config.js";

export default function BlockDetailsPage() {
  return (
    <>
      <RawNav />
      <BlockDetailsContainer />
    </>
  );
}

export const getStaticProps: GetStaticProps<{}> = async ({ locale }) => ({
  props: {
    ...(await serverSideTranslations(
      locale ?? "en",
      ["common"],
      nextI18NextConfig,
    )),
  },
});

export async function getStaticPaths() {
  return { paths: [], fallback: true };
}
