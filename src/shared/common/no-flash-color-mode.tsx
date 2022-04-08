// Need to be inlined to prevent dark mode FOUC
import { ThemeCode } from "onefx/lib/styletron-react";

const storageKey = "theme";
export const noFlashColorMode = ({
  defaultMode,
  respectPrefersColorScheme = true,
}: {
  defaultMode: ThemeCode;
  respectPrefersColorScheme?: boolean;
}): string => {
  return `(function() {
  var defaultMode = '${defaultMode}';
  var respectPrefersColorScheme = ${respectPrefersColorScheme};

  function setDataThemeAttribute(theme) {
    document.documentElement.setAttribute('data-theme', theme);
    if (theme === "dark") {
      document.documentElement.classList.add("dark-theme-applied");
    } else {
      document.documentElement.classList.remove("dark-theme-applied");
    }
  }

  function getStoredTheme() {
    var theme = null;
    try {
      theme = localStorage.getItem('${storageKey}');
    } catch (err) {}
    return theme;
  }

  var storedTheme = getStoredTheme();
  if (storedTheme !== null) {
    setDataThemeAttribute(storedTheme);
  } else {
    if (
      respectPrefersColorScheme &&
      window.matchMedia('(prefers-color-scheme: dark)').matches
    ) {
      setDataThemeAttribute('dark');
    } else if (
      respectPrefersColorScheme &&
      window.matchMedia('(prefers-color-scheme: light)').matches
    ) {
      setDataThemeAttribute('light');
    } else {
      setDataThemeAttribute(defaultMode === 'dark' ? 'dark' : 'light');
    }
  }
})();`;
};
