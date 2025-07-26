# Release History

## 0.2.0 (UNRELEASED)

### Features Added

- Add `environment` command to check and validate environment prerequisites for Azure Go SDK generation.
- Add `generate` command to generate Azure Go SDK packages from TypeSpec specifications.

### Breaking Changes

- Remove `go-version` flag from all commands. It is useless since the code generator could handle it.

### Bugs Fixed

- Refined dependency upgrade logic to explicitly upgrade `azcore` and `azidentity` dependencies instead of using generic `go get -u ./... toolchain@none`

## 0.1.0 (2025-07-21)
- Publish versioning package for generator.
