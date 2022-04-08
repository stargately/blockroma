import { Switch } from "onefx/lib/react-router";
import { Route } from "onefx/lib/react-router-dom";
import { styled } from "onefx/lib/styletron-react";
import React from "react";
import { Footer, FOOTER_ABOVE } from "@/shared/common/footer";
import { Head } from "@/shared/common/head";
import { NotFound } from "@/shared/common/not-found";
import { ScrollToTop } from "@/shared/common/scroll-top";
import { fonts } from "@/shared/common/styles/style-font";
import { Home } from "@/shared/home/home";
import { useGtag } from "@/shared/common/use-gtag";
import { RawNav } from "@/shared/home/components/raw-nav";
import { BlockDetailsContainer } from "@/shared/block-details-container/block-details-container";
import { PageStyles } from "@/shared/common/page-styles";
import { BlksTableContainer } from "@/shared/blks-table-container/blks-table-container";
import { TxsTableContainer } from "./txs-table-container/txs-table-container";
import { AddressDetailsContainer } from "./address-details-container/address-details-container";
import { TxDetailsContainer } from "./tx-details-container/tx-details-container";

export function App(): JSX.Element {
  useGtag();
  return (
    <RootStyle>
      <Head />
      <RawNav />
      <div style={FOOTER_ABOVE}>
        <ScrollToTop>
          <Switch>
            <Route exact path="/">
              <Home />
              <PageStyles stylePath="stylesheets/main-page.css" />
            </Route>
            <Route exact path="/block/:blockNumber">
              <BlockDetailsContainer />
              <PageStyles stylePath="stylesheets/app.css" />
            </Route>
            <Route exact path="/txs">
              <TxsTableContainer />
              <PageStyles stylePath="stylesheets/app.css" />
            </Route>
            <Route exact path="/blocks">
              <BlksTableContainer />
              <PageStyles stylePath="stylesheets/app.css" />
            </Route>
            <Route exact path="/address/:addressHash*">
              <AddressDetailsContainer />
              <PageStyles stylePath="stylesheets/app.css" />
            </Route>
            <Route exact path="/tx/:txHash*">
              <TxDetailsContainer />
              <PageStyles stylePath="stylesheets/app.css" />
            </Route>
            <Route path="*">
              <NotFound />
              <PageStyles stylePath="stylesheets/app.css" />
            </Route>
          </Switch>
        </ScrollToTop>
      </div>
      <Footer />
    </RootStyle>
  );
}

const RootStyle = styled("div", ({ $theme }) => ({
  ...fonts.body,
  backgroundColor: $theme?.colors.black10,
  color: $theme?.colors.text01,
  textRendering: "optimizeLegibility",
}));
