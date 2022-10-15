import type { NextPage } from 'next';
import Head from 'next/head';

const Home: NextPage = () => {
  return (
    <div>
      <Head>
        <title>Doodemo Cook</title>
        <meta
          name="description"
          content="Doodemo Cook"
        />
        <link
          rel="icon"
          href="/favicon.ico"
        />
      </Head>
      <div>Doodemo Cook</div>
    </div>
  );
};

export default Home;
