import "../styles/globals.css";
import type { AppProps } from "next/app";
import React from "react";
import Head from "next/head";
import "react-toastify/dist/ReactToastify.css";
import { ToastContainer } from "react-toastify";
import { ProgressContextProvider } from "../hooks/useProgress";

function App({ Component, pageProps }: AppProps) {
  return (
    <>
      <Head>
        <title>Filter It!</title>
        <meta
          name="description"
          content="Apply filters to your favourite images and share them with others!"
        />
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link
          href="https://fonts.googleapis.com/css2?family=Libre+Baskerville:ital,wght@0,400;0,700;1,400&display=swap"
          rel="stylesheet"
        />
        <link
          href="https://fonts.googleapis.com/css2?family=Rouge+Script&display=swap"
          rel="stylesheet"
        />
      </Head>

      <ToastContainer />
      <ProgressContextProvider>
        <Component {...pageProps} />
      </ProgressContextProvider>
    </>
  );
}

export default App;
