import React, { useEffect, useState } from "react";
import { styled } from "./styled";
import { contentPadding, maxContentWidth } from "./styles/style-padding";
import { media, PALM_WIDTH } from "./styles/style-media";
import { transition } from "./styles/style-animation";
import OutsideClickHandler from "react-outside-click-handler";
import { CommonMargin } from "./common-margin";
import { Hamburger } from "./icons/hamburger.svg";
import { Cross } from "./icons/cross.svg";
import { Icon } from "./icon";
import Link from "next/link";

export const TOP_BAR_HEIGHT = 52;

const Flex = styled("div", () => ({
  flexDirection: "row",
  display: "flex",
  boxSizing: "border-box",
}));

const Menu = styled("div", {
  display: "flex!important",
  [media.palm]: {
    display: "none!important",
  },
});

const BarPlaceholder = styled("div", ({ $theme }) => {
  const height = TOP_BAR_HEIGHT / 2;
  return {
    display: "block",
    padding: `${height}px ${height}px ${height}px ${height}px`,
    backgroundColor: $theme.colors.nav01,
  };
});

const LogoWrapper = styled("a", {
  height: `${TOP_BAR_HEIGHT}px`,
  display: "flex",
  alignItems: "center",
});

function Logo(): JSX.Element {
  return (
    <LogoWrapper href="/">
      <Icon url={"/favicon.svg"} />
    </LogoWrapper>
  );
}

const MaxWidth = styled("div", () => ({
  display: "flex",
  flexDirection: "row",
  ...maxContentWidth,
  justifyContent: "space-between",
  alignItems: "center",
}));

const A = styled(Link, ({ $theme }) => ({
  marginLeft: "14px",
  textDecoration: "none",
  ":hover": {
    color: $theme.colors.primary,
  },
  transition,
  color: $theme.colors.text01,
  [media.palm]: {
    boxSizing: "border-box",
    width: "100%",
    padding: "16px 0 16px 0",
    borderBottom: "1px #EDEDED solid",
  },
}));

// @ts-ignore
const Dropdown = styled("div", ({ $theme }) => ({
  backgroundColor: $theme.colors.nav01,
  display: "flex",
  flexDirection: "column",
  ...contentPadding,
  position: "fixed",
  top: TOP_BAR_HEIGHT,
  "z-index": "1",
  width: "100vw",
  height: "100vh",
  alignItems: "flex-end!important",
  boxSizing: "border-box",
}));

function HamburgerBtn({
  displayMobileMenu,
  children,
  onClick,
}: {
  displayMobileMenu: boolean;
  children: Array<JSX.Element> | JSX.Element;
  onClick: () => void;
}): JSX.Element {
  const Styled = styled("div", ({ $theme }) => ({
    ":hover": {
      color: $theme.colors.primary,
    },
    display: "none!important",
    [media.palm]: {
      display: "flex!important",
      ...(displayMobileMenu ? { display: "none!important" } : {}),
    },
    cursor: "pointer",
    justifyContent: "center",
  }));
  return <Styled onClick={onClick}>{children}</Styled>;
}

function CrossBtn({
  displayMobileMenu,
  children,
}: {
  displayMobileMenu: boolean;
  children: Array<JSX.Element> | JSX.Element;
}): JSX.Element {
  const Styled = styled("div", ({ $theme }) => ({
    ":hover": {
      color: $theme.colors.primary,
    },
    display: "none!important",
    [media.palm]: {
      display: "none!important",
      ...(displayMobileMenu ? { display: "flex!important" } : {}),
    },
    cursor: "pointer",
    justifyContent: "center",
    padding: "5px",
  }));
  return <Styled>{children}</Styled>;
}

const BrandText = styled("a", ({ $theme }) => ({
  textDecoration: "none",
  ":hover": {
    color: $theme.colors.primary,
  },
  fontWeight: 700,
  transition,
  color: $theme.colors.text01,
  marginLeft: 0,
}));

const Bar = styled("nav", ({ $theme }) => ({
  display: "flex",
  flexDirection: "row",
  justifyContent: "center",
  alignItems: "center",
  lineHeight: `${TOP_BAR_HEIGHT}px`,
  height: `${TOP_BAR_HEIGHT}px`,
  backgroundColor: $theme.colors.nav01,
  color: $theme.colors.text01,
  borderBottom: "1px solid rgb(238, 236, 236)",
  position: "fixed",
  top: "0px",
  left: "0px",
  "z-index": "70",
  ...contentPadding,
  boxSizing: "border-box",
}));

export const TopBar = (): JSX.Element => {
  const routePrefix = "/";

  const [displayMobileMenu, setDisplayMobileMenu] = useState(false);
  useEffect(() => {
    window.addEventListener("resize", () => {
      if (
        document.documentElement &&
        document.documentElement.clientWidth > PALM_WIDTH
      ) {
        setDisplayMobileMenu(false);
      }
    });
  }, []);

  const hideMobileMenu = (): void => {
    setDisplayMobileMenu(false);
  };

  const renderMenu = (): JSX.Element => (
    <>
      <A href={`${routePrefix}`} onClick={hideMobileMenu}>
        Home
      </A>
      <A href={`${routePrefix}about-us`} onClick={hideMobileMenu}>
        About Us
      </A>
    </>
  );

  const renderMobileMenu = (): JSX.Element | null => {
    if (!displayMobileMenu) {
      return null;
    }

    return (
      <OutsideClickHandler onOutsideClick={hideMobileMenu}>
        <Dropdown>{renderMenu()}</Dropdown>
      </OutsideClickHandler>
    );
  };

  return (
    <div>
      <Bar>
        <MaxWidth>
          <Flex>
            <Logo />
            <CommonMargin />
            <BrandText href="/">OneFx</BrandText>
          </Flex>
          <Flex>
            <Menu>{renderMenu()}</Menu>
          </Flex>
          <HamburgerBtn
            displayMobileMenu={displayMobileMenu}
            onClick={() => {
              setDisplayMobileMenu(true);
            }}
          >
            <Hamburger />
          </HamburgerBtn>
          <CrossBtn displayMobileMenu={displayMobileMenu}>
            <Cross />
          </CrossBtn>
        </MaxWidth>
      </Bar>
      <BarPlaceholder />
      {renderMobileMenu()}
    </div>
  );
};
