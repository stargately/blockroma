# CI Workflows

This directory contains GitHub Actions workflows for continuous integration (CI) of all subprojects in the Blockroma repository.

## Overview

Each subproject has its own dedicated CI workflow that runs on:
- **Pull requests** to the `main` branch
- **Pushes** to the `main` branch

Workflows are triggered only when changes are made to the specific subproject directory or the workflow file itself (using path filters).

## Workflows

### 1. v1 CI (`v1-ci.yml`)

**Subproject:** `v1/` - Main blockchain explorer application

**Technology:** Node.js 16 + TypeScript

**Steps:**
- Install dependencies with `npm install --legacy-peer-deps`
- Run i18n format check
- Run ESLint
- Build the project
- Run tests with coverage

**Path filters:** `v1/**`, `.github/workflows/v1-ci.yml`

### 2. frontend-v2 CI (`frontend-v2-ci.yml`)

**Subproject:** `frontend-v2/` - Next.js frontend application

**Technology:** Node.js 18 + Next.js + TypeScript + Yarn

**Steps:**
- Install dependencies with `yarn install --frozen-lockfile`
- Run Next.js linter
- Build the project
- Run tests (lint + format check)

**Path filters:** `frontend-v2/**`, `.github/workflows/frontend-v2-ci.yml`

### 3. doc CI (`doc-ci.yml`)

**Subproject:** `doc/` - Docusaurus documentation site

**Technology:** Node.js 16 + Docusaurus

**Steps:**
- Install dependencies with `npm install`
- Run TypeScript type checking (non-blocking due to pre-existing issues)
- Build the documentation site

**Path filters:** `doc/**`, `.github/workflows/doc-ci.yml`

### 4. soroban-indexer CI (`soroban-indexer-ci.yml`)

**Subproject:** `soroban/indexer/` - Stellar Soroban blockchain indexer

**Technology:** Go 1.23

**Steps:**
- Download and verify dependencies
- Check code formatting
- Run `go vet` (non-blocking due to pre-existing issues)
- Build the indexer
- Run tests with coverage
- Upload coverage to Codecov

**Path filters:** `soroban/indexer/**`, `.github/workflows/soroban-indexer-ci.yml`

## Notes

### Pre-existing Issues

Some workflows have steps marked as `continue-on-error: true` to handle pre-existing issues in the codebase:
- **doc CI**: TypeScript type checking errors in `docusaurus.config.js`
- **soroban-indexer CI**: Go vet warnings about bit shifts

These are documented and should be addressed separately but don't block CI.

### Dependency Management

- **v1**: Uses `--legacy-peer-deps` flag due to peer dependency conflicts
- **frontend-v2**: Uses Yarn with frozen lockfile for reproducible builds
- **doc**: Uses npm with standard installation
- **soroban/indexer**: Uses Go modules

## Adding New Workflows

When adding a new subproject or workflow:

1. Create a new workflow file in `.github/workflows/`
2. Name it descriptively (e.g., `<subproject>-ci.yml`)
3. Configure path filters to only trigger on relevant changes
4. Add appropriate build/test steps for the technology stack
5. Document the workflow in this README

## Local Testing

You can test the build steps locally before pushing:

```bash
# v1
cd v1 && npm install --legacy-peer-deps && npm run build && npm test

# frontend-v2
cd frontend-v2 && yarn install --frozen-lockfile && yarn build && yarn test

# doc
cd doc && npm install && npm run build

# soroban/indexer
cd soroban/indexer && make deps && make build && make test
```

## Troubleshooting

### Workflow not triggering
- Check if your changes match the path filters
- Ensure the workflow file is valid YAML
- Check GitHub Actions tab for errors

### Build failures
- Run the steps locally to reproduce the issue
- Check if dependencies need updating
- Review the build logs in GitHub Actions

### Cache issues
- Workflows use caching to speed up builds
- If you suspect cache corruption, re-run the workflow or clear caches

## Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Workflow Syntax Reference](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions)
