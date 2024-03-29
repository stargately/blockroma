# Blockroma Frontend

![Blockroma](https://tp-misc.b-cdn.net/blockroma-v0.1.png)

## Getting Started

Blockroma is a blockchain explorer. It is built with the modern web stack - TypeScript, KOA, React, SASS, Apollo GraphQL, TypeORM, and PostgreSQL. It is open-sourced under [GPL license](#license).

![blockroma-v0.1-demo](https://tp-misc.b-cdn.net/blockroma-v0.1-2.gif)

> The project is still in development with an unstable version 0.1.0. As a result, there might be breaking changes before 1.0.0.

### Download the project

```bash
git clone git@github.com:blockedenhq/blockroma-frontend.git
```

### Run the project

```bash
cd blockroma-frontend

yarn install
```

#### Development mode

To run your project in development mode, run:

```bash
yarn dev
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
