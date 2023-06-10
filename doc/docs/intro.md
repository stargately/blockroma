---
title: Blockroma
---

People use blockchains to make online transactions or payments. However, blockchain APIs and data storage are not built for data analysis purposes. This makes it hard for data scientists or analysts to step into those TB-level volumes of data and find out insights using SQL queries.

[Blockroma](https://blockroma.com) is a blockchain indexer and explorer that extracts, transforms, and saves EVM-based blockchain data into PostgreSQL database, GraphQL APIs, and web-browsable UIs.

## How does it work?

Ethererum blockchain serves all the block and transaction data via its JSON RPC API. The blockchain **realtime** indexer will listen to the latest block height and fetch the block accordingly. And meanwhile, the **catchup** indexer will scan a range of blocks. Then we apply parsers to those raw data and then insert results into the database.

![Blockroma and Ethereum](https://tp-misc.b-cdn.net/blockchain-explorer-architecture@2x.png)

Here is how to run Blockroma against Ethereum:

```
git clone git@github.com:stargately/blockroma.git
cd blockroma

nvm use lts/gallium
npm install

# prepare environment variable
cp ./.env.tmpl ./.env
```

In the `.env` file, fill in your database URL and mainnet endpoint, for example,

```shell
DATABASE_URL=postgres://localhost:5432/blockroma_eth
ETHEREUM_MAINNET_RPC_URL=https://mainnet.infura.io/v3/TODO_TOKEN
```

And then, you could customize the indexer in `config/ethereum-mainnet.js`

```js
module.exports = {
  chain: {
    chainId: 1,
    chainName: "Ethereum",
    symbol: "ETH",
    rpcUrl: process.env.ETHEREUM_MAINNET_RPC_URL,
    decimals: 18,
    networkPath: "",
  },
  indexer: {
    catchup: {
      enabled: false, // enable if you want
      blocksBatchSize: 200,
      blockNumberRanges: [
        [0, "latest"], // range to scan
      ],
    },
    realtime: {
      enabled: true,
    },
  },
};
```

Once you setup the indexer and explorer, you could run it with

```
NODE_ENV=production npm run build
NODE_CONFIG_ENV=ethereum-mainnet npm run start
```

## How to browse Blockchain data?

Then, open your browser with http://localhost:4134/, and you can visit

- http://localhost:4134/blocks for blocks
- http://localhost:4134/txs for transactions
- http://localhost:4134/api-gateway/ for APIs

Or simply tap your `/` key and input your hash, address or block height, and finally hit `Enter` to search.

![Search with blockroma](https://tp-misc.b-cdn.net/search-bar-blockroma.png)

## How to run SQL queries on Ethereum?

Connect to the database that is serving Blockroma service and then you will see `address`, `block`, `transaction` for ETH and then ERC20, ERC721, ERC1155 `token` and `tokenTransfer` tables.
