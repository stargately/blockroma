import { BlockDetailsContainer } from "@/shared/block-details-container/block-details-container";
import { GetStaticProps } from "next";
import { serverSideTranslations } from "next-i18next/serverSideTranslations";
import { RawNav } from "@/shared/home/components/raw-nav";
import React from "react";
import nextI18NextConfig from "../../next-i18next.config.js";
import { AddressDetailsContainer } from "@/shared/address-details-container/address-details-container";
import { SeoHead } from "@/shared/common/seo-head";

export default function AddressHashPage() {
  return (
    <>
      <SeoHead config={{ title: "Address Details" }} />
      <RawNav />
      <AddressDetailsContainer />
    </>
  );
}

export const getStaticProps: GetStaticProps<{}> = async ({ locale }) => ({
  props: {
    ...(await serverSideTranslations(
      locale ?? "en",
      ["common"],
      nextI18NextConfig
    )),
  },
});

export async function getStaticPaths() {
  return { paths: [], fallback: true };
}
