# Release History

## 0.4.4 (2026-01-28)

### Bugs Fixed

- Fixed `version` command to skip CHANGELOG.md update when the file doesn't exist and `--sdkversion` is specified.
- Updated `version` command to only update the version in the latest CHANGELOG.md entry (not regenerate changelog content) when `--sdkversion` is specified.
- Use `format.Node` instead of `printer.Fprint` to keep the format aligned with gofmt.

## 0.4.3 (2025-12-17)

### Bugs Fixed

- Add back `go.mod.tpl` for Swagger onboard case since it is a must for `go generate`.

## 0.4.2 (2025-12-08)

### Bugs Fixed

- Move logic of determining preview version after code generation to avoid incorrect version calculation.
- Fix wrong override logic for stable/beta version determination.

## 0.4.1 (2025-12-02)

### Bugs Fixed

- Fix the wrong constant value for beta release type.

## 0.4.0 (2025-12-01)

### Features Added

- Add `changelog` command to generate and update changelog content for SDK packages based on code changes.
- Add `version` command to calculate and update version numbers across all version-related files.
- Add comprehensive helper utilities for package status determination, path resolution, and version management.
- Support detecting parameter renaming as breaking change in changelog generation.

### Breaking Changes

- Split version calculation and changelog generation logic into separate packages (`changelog` and `version`).
- Move constants and enums from `cmd/v2/common` to shared `utils` package for better reusability.
- Refactor `GenerateForSingleRPNamespace` and `GenerateForSingleTypeSpec` to use status-based generator selection.
- Remove template files for CHANGELOG.md and go.mod (now generated programmatically).

### Other Changes

- Upgrade Go version to 1.24.0 for dependency management.

## 0.3.1 (2025-11-13)

### Bugs Fixed

- Fix wrong parsing logic for module name config for swagger with major version suffix.

## 0.3.0 (2025-11-11)

### Features Added

- Add `build` command to build the SDK package.
- Add back `force-stable-version` flag to `release-v2` command to support generating stable version even if input-files contains preview version.

## 0.2.2 (2025-10-14)

### Other Changes

- Fixed Go version and prevented toolchain upgrade in `go get` commands by adding `toolchain@none go@1.23.0` parameters.

## 0.2.1 (2025-09-25)

### Other Changes

- Secure tsp-client usage by using pinned tsp-client versions from `eng/common/tsp-client`.

### Bugs Fixed

- Fixed `FuncFilter` to be compatible with unchanged parameter signatures.

## 0.2.0 (2025-09-18)

### Features Added

- Add `environment` command to check and validate environment prerequisites for Azure Go SDK generation.
- Add `generate` command to generate Azure Go SDK packages from TypeSpec specifications.
- Support major version suffix in the `module` flag to specify the major version of the generated module.

### Bugs Fixed

- Use module as the root of API view generation path.

### Breaking Changes

- Remove `go-version` flag from all commands. It is useless since the code generator could handle it.
- Remove `version-number` flag from all commands. It is no longer supported since the version is now configured in the `module` flag.

### Bugs Fixed

- Refined dependency upgrade logic to explicitly upgrade `azcore` and `azidentity` dependencies instead of using generic `go get -u ./... toolchain@none`

## 0.1.0 (2025-07-21)
- Publish versioning package for generator.
