import React from "react";
import { actionSetTheme } from "@/shared/common/base-reducer";
import { connect } from "react-redux";
import { AdContainer } from "@/shared/home/components/ad-container";
import { HomeTransactionsContainer } from "@/shared/home/home-transactions-container";
import { HomeBlocksContainer } from "@/shared/home/home-blocks-container";
import { RawModals } from "./components/raw-modals";

export const Home = connect(
  (state: { base: { themeCode: "dark" | "light" } }) => ({
    themeCode: state.base.themeCode,
  }),
  (dispatch) => ({
    actionSetTheme: (themeCode: "dark" | "light") => {
      dispatch(actionSetTheme(themeCode));
    },
  })
)((): JSX.Element => {
  return (
    <div>
      <div className="layout-container">
        {/*
                  TODO(dora): skip for now
        <RawNavAlert />
        */}
        <main className="js-ad-dependant-pt pt-5">
          <p className="alert alert-info" role="alert" />
          <p className="alert alert-danger" role="alert" />
          {/*
          TODO(dora): skip for now
          <DashboardBanner />
          */}
          <section className="container">
            <HomeBlocksContainer />

            <AdContainer />

            <HomeTransactionsContainer />
          </section>
        </main>
      </div>

      <RawModals />

      <div className="offline-ui offline-ui-up">
        <div className="offline-ui-content" />
        <a className="offline-ui-retry" />
      </div>
      <div className="offline-ui offline-ui-up">
        <div className="offline-ui-content" />
        <a className="offline-ui-retry" />
      </div>
    </div>
  );
});
