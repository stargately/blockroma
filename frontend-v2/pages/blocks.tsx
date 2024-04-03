import React from "react";
import type { GetStaticProps, InferGetStaticPropsType } from "next";
import { serverSideTranslations } from "next-i18next/serverSideTranslations";
import { RawNav } from "@/shared/home/components/raw-nav";
import { BlksTableContainer } from "@/shared/blks-table-container/blks-table-container";
import { SeoHead } from "@/shared/common/seo-head";

export default function BlocksPage() {
  return (
    <div>
      <SeoHead config={{ title: "Blocks" }} />
      <RawNav />
      <BlksTableContainer />
    </div>
  );
}

// or getServerSideProps: GetServerSideProps<Props> = async ({ locale })
export const getStaticProps: GetStaticProps<{}> = async ({ locale }) => ({
  props: {
    ...(await serverSideTranslations(locale ?? "en", ["common"])),
  },
});
