# Release History

## 2.0.0 (2023-11-16)

### Breaking Changes
* Removed `testutil.GetEnv`, use `github.com/Azure/azure-sdk-for-go/sdk/internal/recording.GetEnvVariable` instead.
* Remove `testutil.GenerateAlphaNumericID`, use `github.com/Azure/azure-sdk-for-go/sdk/internal/recording.GenerateAlphaNumericID` instead.

### Other Changes
* Update dependencies: `github.com/Azure/azure-sdk-for-go/sdk/internal@v1.5.0`

## 1.1.2 (2023-03-03)

### Other Changes
* Deprecate `testutil.GetEnv`, use `github.com/Azure/azure-sdk-for-go/sdk/internal/recording.GetEnvVariable` instead
* Deprecate `testutil.GenerateAlphaNumericID`, use `github.com/Azure/azure-sdk-for-go/sdk/internal/recording.GenerateAlphaNumericID` instead
* Migrating all test recording files to assets repo.=

## 1.1.1 (2022-08-30)

### Bugs Fixed
* Fix seed not stable with `GenerateAlphaNumericID` when playback

## 1.1.0 (2022-08-24)

### Features Added
* Add `GenerateAlphaNumericID` to testutil

## 1.0.1 (2022-06-23)

### Other Changes
* Upgrade `azcore` version and change test `poller` method

## 1.0.0 (2022-05-16)

### Features Added
* Export FakeCredential

### Other Changes
* Upgrade dependencies of azcore, azidentity and armresources to the stable version

## 0.3.0 (2022-04-08)

### Breaking Changes
* Upgrade to generic version for test helper

## 0.2.0 (2022-03-16)

### Features Added
* Add helper method for ARM template deployment
* Add delegate stop method return for start recording

## 0.1.0 (2022-03-10)

### Features Added
* Add test util for resource manager

