# Release History

## 0.2.1 (Unreleased)

### Other Changes

- Secure tsp-client usage by using pinned tsp-client versions from `eng/common/tsp-client`.

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
