import Document, { Html, Head, Main, NextScript } from "next/document";
import { Provider as StyletronProvider } from "styletron-react";
import { styletron } from "../src/styletron";
import { defaultSeoConfig } from "../src/common/seo-head";

const beConfig = defaultSeoConfig;

class MyDocument extends Document {
  static async getInitialProps(context) {
    const renderPage = () =>
      context.renderPage({
        enhanceApp: (App) => (props) =>
          (
            <StyletronProvider value={styletron}>
              <App {...props} />
            </StyletronProvider>
          ),
      });

    const initialProps = await Document.getInitialProps({
      ...context,
      renderPage,
    });
    const stylesheets = styletron.getStylesheets() || [];
    return { ...initialProps, stylesheets };
  }

  render() {
    return (
      <Html className="dark-theme-applied">
        <Head>
          <meta name="twitter:site" content={beConfig.twitter} />
          <meta name="twitter:image" content={beConfig.previewImageUrl} />
          <meta name="twitter:title" content={beConfig.title} />
          <meta name="twitter:description" content={beConfig.description} />
          <meta property="twitter:card" content="summary_large_image" />
          <meta property="og:title" content={beConfig.title} />
          <meta property="og:description" content={beConfig.description} />
          <meta
            property="og:image"
            name="og:image"
            content={beConfig.previewImageUrl}
          />
          <meta charSet="utf-8" />

          <link
            rel="stylesheet"
            href="https://fonts.googleapis.com/css2?family=Inter&display=swap"
          />

          {this.props.stylesheets.map((sheet, i) => (
            <style
              className="_styletron_hydrate_"
              dangerouslySetInnerHTML={{ __html: sheet.css }}
              media={sheet.attrs.media}
              data-hydrate={sheet.attrs["data-hydrate"]}
              key={i}
            />
          ))}
        </Head>
        <body>
          <Main />
          <NextScript />
        </body>
      </Html>
    );
  }
}

export default MyDocument;
