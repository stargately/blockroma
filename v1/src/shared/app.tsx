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
import { TokensContainer } from "./token-container/tokens-container";

const routes = [
  {
    path: "/",
    style: "stylesheets/app.css",
    component: <Home />,
    exact: true,
  },
  {
    path: "/block/:blockNumber",
    style: "stylesheets/app.css",
    component: <BlockDetailsContainer />,
    exact: true,
  },
  {
    path: "/txs",
    style: "stylesheets/app.css",
    component: <TxsTableContainer />,
    exact: true,
  },
  {
    path: "/blocks",
    style: "stylesheets/app.css",
    component: <BlksTableContainer />,
    exact: true,
  },
  {
    path: "/address/:addressHash*",
    style: "stylesheets/app.css",
    component: <AddressDetailsContainer />,
    exact: true,
  },
  {
    path: "/tx/:txHash*",
    style: "stylesheets/app.css",
    component: <TxDetailsContainer />,
    exact: true,
  },
  {
    path: "/tokens",
    style: "stylesheets/app.css",
    component: <TokensContainer />,
    exact: true,
  },
  {
    path: "*",
    style: "stylesheets/app.css",
    component: <NotFound />,
  },
];

export function App(): JSX.Element {
  useGtag();
  return (
    <RootStyle>
      <Head />
      <ScrollToTop>
        <Switch>
          {routes.map((r) => (
            <Route key={r.path} exact={r.exact} path={r.path}>
              <RawNav />
              <div style={FOOTER_ABOVE}>
                {r.component}
                <PageStyles stylePath={r.style} />
              </div>
            </Route>
          ))}
        </Switch>
      </ScrollToTop>
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
