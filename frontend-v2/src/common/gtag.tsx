import Script from "next/script";
import React from "react";

export const Gtag: React.FC = () => {
  return (
    <>
      <Script
        src="https://www.googletagmanager.com/gtag/js?id=G-VKPSYZ2K22"
        strategy="afterInteractive"
      />
      <Script id="google-analytics" strategy="afterInteractive">
        {`
          window.dataLayer = window.dataLayer || [];
          function gtag(){window.dataLayer.push(arguments);}
          gtag('js', new Date());

          gtag('config', 'GA_MEASUREMENT_ID');
        `}
      </Script>
    </>
  );
};
