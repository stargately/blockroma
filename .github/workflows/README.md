# CI/CD Workflows

This repository has separate CI workflows for each subproject to ensure efficient and targeted testing.

## Workflow Overview

### 1. V1 CI (`v1-ci.yml`)
**Path:** `.github/workflows/v1-ci.yml`  
**Triggers:** 
- Push to `main` branch (when `v1/**` files change)
- Pull requests to `main` branch (when `v1/**` files change)

**Steps:**
- Install Node.js 16
- Install dependencies with Yarn
- Run linter
- Build the project
- Run tests with coverage

### 2. Frontend V2 CI (`frontend-v2-ci.yml`)
**Path:** `.github/workflows/frontend-v2-ci.yml`  
**Triggers:**
- Push to `main` branch (when `frontend-v2/**` files change)
- Pull requests to `main` branch (when `frontend-v2/**` files change)

**Steps:**
- Install Node.js 18
- Install dependencies with Yarn
- Run linter
- Build the Next.js project
- Run tests

### 3. Doc CI (`doc-ci.yml`)
**Path:** `.github/workflows/doc-ci.yml`  
**Triggers:**
- Push to `main` branch (when `doc/**` files change)
- Pull requests to `main` branch (when `doc/**` files change)

**Steps:**
- Install Node.js 16
- Install dependencies with Yarn
- Run TypeScript type checking
- Build the Docusaurus documentation

### 4. Soroban Indexer CI (`soroban-indexer-ci.yml`)
**Path:** `.github/workflows/soroban-indexer-ci.yml`  
**Triggers:**
- Push to `main` branch (when `soroban/indexer/**` files change)
- Pull requests to `main` branch (when `soroban/indexer/**` files change)

**Steps:**
- Install Go 1.23
- Download and verify dependencies
- Run format check (non-blocking)
- Run go vet (non-blocking)
- Build the binary
- Run tests with coverage
- Upload coverage to Codecov

## Path Filters

Each workflow includes path filters to optimize CI runs. The workflow will only trigger when:
1. Files in the specific subproject directory change, OR
2. The workflow file itself changes

This prevents unnecessary CI runs when changes are made to unrelated subprojects.

## Notes

- **V1** and **Doc** use Node.js 16 as specified in their package.json
- **Frontend V2** uses Node.js 18 as it's a Next.js project with newer dependencies
- **Soroban Indexer** uses Go 1.23 as specified in go.mod
- Format and vet checks in the Soroban Indexer workflow are set to `continue-on-error: true` to allow the build to succeed despite pre-existing code style issues
- All Node.js projects use Yarn for package management
