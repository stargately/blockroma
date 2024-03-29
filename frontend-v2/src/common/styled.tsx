/* tslint:disable:no-any */
import React, { PropsWithChildren } from "react";
import {
  createStyled,
  StyletronComponent,
  StyletronWrapper,
} from "styletron-react";
import { driver, getInitialStyle, StyleObject } from "styletron-standard";
import { Theme } from "./theme-types";

const wrapper: StyletronWrapper = (StyledComponent) =>
  function withThemeHOC(props: PropsWithChildren<any>): JSX.Element {
    return (
      <>{($theme: any) => <StyledComponent {...props} $theme={$theme} />}</>
    );
  };

export interface StyledFn {
  <
    C extends keyof JSX.IntrinsicElements | React.ComponentType<any>,
    P extends Record<string, unknown>,
  >(
    component: C,
    style: (props: P & { $theme: Theme }) => StyleObject,
  ): StyletronComponent<
    Pick<
      React.ComponentProps<C>,
      Exclude<keyof React.ComponentProps<C>, { className: string }>
    > &
      P
  >;

  <C extends keyof JSX.IntrinsicElements | React.ComponentType<any>>(
    component: C,
    style: StyleObject,
  ): StyletronComponent<
    Pick<
      React.ComponentProps<C>,
      Exclude<keyof React.ComponentProps<C>, { className: string }>
    >
  >;
}

export const styled: StyledFn = createStyled({
  wrapper,
  getInitialStyle,
  driver,
});
