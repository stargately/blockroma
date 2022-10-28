# Blockroma

![Discord](https://img.shields.io/discord/817917152467222529?label=Discord&logo=discord&style=flat~~~~)

![Blockroma](https://tp-misc.b-cdn.net/blockroma-v0.1.png)

## Getting Started

Blockroma is a blockchain explorer for Ethereum web3 compatible blockchains. It is built with the modern web stack - TypeScript, KOA, React, SASS, Apollo GraphQL, TypeORM, and PostgreSQL. It is open-sourced under [GPL license](#license).

![blockroma-v0.1-demo](https://tp-misc.b-cdn.net/blockroma-v0.1-2.gif)

> The project is still in development with an unstable version 0.1.0. As a result, there might be breaking changes before 1.0.0.

Feature List

- [x] realtime + catchup indexer, based on web3 JSONRPC and WSS APIs, for blocks, txs, and addresses.
- [x] GraphQL APIs and PostgreSQL Data models for blocks, txs, and addresses.
- [x] basic web UI for home page, search bar, blocks, txs, and addresses.
- [x] dark mode and light mode.
- [x] i18n / internationalization / multi-language / English, Japanese, Chinese
- [x] customizable for other blockchains.
- [ ] developer guide.
- [ ] tx contract details and logs.
- [ ] ERC20. [Design](./src/server/service/remote-chain-service/parse-token-transfer/README.md)
- [ ] ERC721.
- [ ] more accurate gas fee calculation.
- [ ] better error handling and loading state.
- [ ] onboard more blockchains.

### Download the project

```bash
git clone git@github.com:stargately/blockroma.git
```

### Run the project

This is intended for \*nix users. If you use Windows, go to [Run on Windows](#run-on-windows). Let's first prepare the environment.

```bash
cd blockroma

npm install

# prepare environment variable
cp ./.env.tmpl ./.env
```

#### Development mode

To run your project in development mode, run:

```bash
npm run watch
```

The development site will be available at [http://localhost:4134](http://localhost:4134).

#### Production Mode

It's sometimes useful to run a project in production mode, for example, to check bundle size or to debug a production-only issue. To run your project in production mode locally, run:

```bash
npm run build-production
NODE_ENV=production npm run start
```

#### NPM scripts

- `npm run test`: test the whole project and generate a test coverage
- `npm run ava ./path/to/test-file.js`: run a specific test file
- `npm run build`: build source code from `src` to `dist`
- `npm run lint`: run the linter
- `npm run kill`: kill the node server occupying the port 5000.

## License

Special thanks to the blockscout project. We are not fans of Elixir but we used its JS and style files, and thus comply with its GPL license.

## Star and Fork the Project

If you like the project or want to support its future development, please give it a star ⭐️ and fork it!
