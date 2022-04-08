import {
  defaultThemeCode,
  ThemeCode,
} from "onefx/lib/styletron-react/theme-provider";

const storeTheme = (newTheme: ThemeCode) => {
  try {
    localStorage.setItem("theme", newTheme);
  } catch (err) {
    // eslint-disable-next-line no-console
    console.error(err);
  }
};

export function baseReducer(
  initialState: { themeCode?: ThemeCode } = { themeCode: defaultThemeCode },
  action: { type: string; payload: ThemeCode }
): { themeCode?: ThemeCode } {
  if (action.type === "SET_THEME") {
    const themeCode = action.payload === "light" ? "light" : "dark";
    window.document &&
      window.document.documentElement.setAttribute("data-theme", themeCode);
    if (themeCode === "dark") {
      window.document.documentElement.classList.add("dark-theme-applied");
    } else {
      window.document.documentElement.classList.remove("dark-theme-applied");
    }

    storeTheme(themeCode);
    return {
      ...initialState,
      themeCode,
    };
  }
  if (!initialState.themeCode) {
    initialState.themeCode = defaultThemeCode;
  }
  return initialState;
}

export function actionSetTheme(themeCode: ThemeCode): {
  type: string;
  payload: ThemeCode;
} {
  return {
    type: "SET_THEME",
    payload: themeCode,
  };
}
