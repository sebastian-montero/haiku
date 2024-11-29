import { Html, Head, Main, NextScript } from "next/document";

export default function Document() {
  return (
    <Html lang="en">
      <Head />
      <head>
        <title>haiku</title>
        <link
          rel="icon"
          href="data:image/svg+xml,<svg xmlns=%22http://www.w3.org/2000/svg%22 viewBox=%220 0 32 32%22><text y=%2224%22 font-size=%2224%22>⿻</text></svg>"
        />
        <meta
          name="description"
          content="haiku is a creative app where you can write anything in real time, share your process, and let others watch your creative flow."
        />
      </head>
      <body className="antialiased">
        <Main />
        <NextScript />
      </body>
    </Html>
  );
}
