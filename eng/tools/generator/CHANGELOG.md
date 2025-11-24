# Release History

## 0.3.2 (unreleased)

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
