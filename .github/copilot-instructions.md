# GitHub Copilot Instructions

## File Modification Rules

### Lock Files

**NEVER modify, update, or regenerate yarn.lock files under any circumstances.**

- `yarn.lock` files are managed by the package manager and should only be updated through explicit `yarn install` or `yarn add/remove` commands
- Do not suggest changes to yarn.lock files
- Do not include yarn.lock files in any code modifications
- If dependencies need to be updated, only modify package.json and let the user run yarn commands manually

### Files to Never Modify

The following files should never be modified by Copilot:
- `**/yarn.lock` - Yarn lock files
- `.git/**` - Git internal files
