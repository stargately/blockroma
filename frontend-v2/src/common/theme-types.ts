type Font = {
  fontSize: string | number;
  lineHeight: string | number;
  fontFamily?: string;
  letterSpacing?: string | number;
};

export type Theme = {
  colors: {
    primary: string;
    secondary: string;

    black: string;
    black10: string;
    black20: string;
    black40: string;
    black60: string;
    black80: string;
    black95: string;

    text01: string;
    textReverse: string;

    white: string;

    error: string; //	Error
    success: string; //	Success
    warning: string; //	Warning
    information: string; //	Information

    nav01: string; //	Global top bar
    nav02: string; //	CTA footer
    nav03: string; //	Global footer
  };
  sizing: Array<string>;
  fonts: Array<Font>;
};
