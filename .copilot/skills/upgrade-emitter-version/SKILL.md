---
name: upgrade-emitter-version
description: Use when upgrading TypeSpec emitter packages in eng/emitter-package.json to their latest versions, or when asked to update TypeSpec compiler and Azure tools versions.
---

# Upgrade Emitter and TypeSpec Versions

Upgrade packages in `eng/emitter-package.json` to their latest versions.

## Prerequisites

- Node.js LTS installed
- [npm-check-updates](https://www.npmjs.com/package/npm-check-updates) installed globally (`npm install -g npm-check-updates`)
- [tsp-client](https://www.npmjs.com/package/@azure-tools/typespec-client-generator-cli) installed globally (`npm install -g @azure-tools/typespec-client-generator-cli`)

## Steps

### Step 1: Change to the eng directory

Change the working directory to the `eng/` folder at the repository root.

### Step 2: Rename emitter-package.json to package.json

Rename `emitter-package.json` to `package.json` so that `ncu` can process it.

### Step 3: Run ncu to upgrade packages

Run `ncu -u` to upgrade all packages in `package.json` to their latest versions.

### Step 4: Ensure versions are absolute

Verify that all versions in `package.json` are absolute (no `~` or `^` prefixes). If any version has a `~` or `^` prefix, remove the prefix so only the version number remains.

### Step 5: Rename package.json back to emitter-package.json

Rename `package.json` back to `emitter-package.json`.

### Step 6: Generate the lock file

Run `tsp-client generate-lock-file` to update `emitter-package-lock.json`.

### Step 7: Commit the changes

Commit both `emitter-package.json` and `emitter-package-lock.json`. Do **not** commit `node_modules` or `package.json`.
