import { Provider as StyletronProvider } from "styletron-react";
import "../styles/globals.css";
import type { AppProps } from "next/app";
import { styletron } from "../src/styletron";
import "../src/stylesheets/app.scss";
import { Provider as ReduxProvider } from "react-redux";
import { configureStore } from "../src/common/configure-store";
import { appWithTranslation } from "next-i18next";
import { ApolloProvider } from "@apollo/client";
import { apolloClient } from "@/shared/common/apollo-client";
import nextI18NextConfig from "../next-i18next.config.js";

export default appWithTranslation(function App({
  Component,
  pageProps,
}: AppProps) {
  return (
    <ApolloProvider client={apolloClient}>
      <ReduxProvider store={configureStore({ base: { theme: "dark" } })}>
        <StyletronProvider value={styletron}>
          <Component {...pageProps} />
        </StyletronProvider>
      </ReduxProvider>
    </ApolloProvider>
  );
}, nextI18NextConfig);
