# Release History

## 1.8.0 (2024-05-09)

### Breaking Changes

* Unexported `Module` and `Version` constants

### Other Changes

* Removed default sanitizers added in v1.6.0 (the test proxy itself now includes these)

## 1.7.0 (2024-05-01)

### Features Added
* Support for local repo override (via presence of eng/target_proxy_version.txt) of invoked test-proxy version.

* `RemoveRegisteredSanitizers` selectively disables sanitizers the test proxy enables by
  default since version 1.0.0-dev.20240422.1

### Breaking Changes

* Deprecated the `go-vcr` based test recording API. Its methods now return errors or panic.
* Changed value of `recording.SanitizedValue` from "sanitized" to "Sanitized" to match the
  test proxy

## 1.6.0 (2024-04-16)

### Features Added

* Options types for `SetBodilessMatcher` and `SetDefaultMatcher` now embed `RecordingOptions`
* Added a collection of default sanitizers for test recordings

## 1.5.2 (2024-02-06)

### Bugs Fixed

* Prevent `exported.Payload` from panicking in the rare event `*http.Response.Body` is `nil`.

### Other Changes

* Update dependencies.

## 1.5.1 (2023-12-06)

### Bugs Fixed

* Recording will restore the original scheme/host after making a successful HTTP(s) call.

## 1.5.0 (2023-11-02)

### Features Added

* Added a new `NonRetriableError` func to the `errorinfo` package. New func serves as an error wrapper for non-retriable errors in the `azure-sdk-for-go/sdk` folder.

## 1.4.0 (2023-10-17)

### Features Added

* Add support for auto-installing the test proxy standalone tooling in the test recording package

### Other Changes

* Updated dependencies.

## 1.3.0 (2023-04-04)

### Features Added
* Added package `poller` which exports various LRO helpers to aid in the creation of custom `PollerHandler[T]`.
* Added package `exported` which contains payload helpers needed by the `poller` package and exported in `azcore`.

## 1.2.0 (2023-03-02)

### Features Added

* Add random alphanumeric string generation support for test-proxy recording framework.

### Bugs Fixed

* Store RNG seed in recordings.

## 1.1.2 (2022-12-12)

### Features Added

- Export user agent formatting code that used to be in azcore's policy_telemetry.go so it can be shared with non-HTTP clients (ie: azservicebus/azeventhubs). ([#19681](https://github.com/Azure/azure-sdk-for-go/pull/19681))

### Other Changes
* Prevented data races in `recording` ([#18763](https://github.com/Azure/azure-sdk-for-go/issues/18763))

## 1.1.1 (2022-11-09)

### Bugs Fixed
* Fixed a race condition in `temporal.Resource[TResource, TState].Get`.

## 1.1.0 (2022-10-20)

### Features Added

* Support test recording assets external to repository

## 1.0.1 (2022-08-22)

### Bugs Fixed
* Don't modify the original *http.Request during recording/perf as it causes failures during retries.

## 1.0.0 (2022-05-12)

### Features Added
* Added `temporal` package for handling of temporal resources.

### Breaking Changes
* Removed logging constants from the `log` package.
* Removed `atomic` package as it wasn't being used.

### Other Changes
* Updated build constraints and Go version to 1.18.

## 0.9.2 (2022-04-04)

### Features Added
* Set sanitizers at the test level by specifying `RecordingOptions.TestInstance`
* Added the `perf` library for performance testing client SDKs

## 0.9.1 (2022-02-01)

### Features Added
* Adds a `CustomDefaultMatcher` that adds headers `:path`, `:authority`, `:method`, and `:scheme` to the default matcher.

## 0.9.0 (2022-01-24)

### Breaking Changes
* The `x-recording-file` is now encoded in the body of a `Start` request, previously was included in a header [#16876](https://github.com/Azure/azure-sdk-for-go/pull/16876).

## 0.8.3 (2021-12-07)

### Features Added
* If `AZURE_RECORD_MODE` is not set, default to `playback`
* If `PROXY_CERT` is not set, try to find it based on the `GOPATH` environment variable and the path to `eng/common/testproxy/dotnet-devcert.crt`
* Adds `NewRecordingHTTPClient()` method which returns an `azcore.Transporter` interface that routes requests to the test proxy [#16221](https://github.com/Azure/azure-sdk-for-go/pull/16221).
* Adds the `SetBodilessMatcher` method [#16256](https://github.com/Azure/azure-sdk-for-go/pull/16256)
* Added variables storage to the `Stop` function. Pass in a `map[string]interface{}` to the `Stop` method options and the values can be retrieved with the `GetVariables(t *testing.T)` function [#16375](https://github.com/Azure/azure-sdk-for-go/pull/16375).

### Breaking Changes
* Renames `ResetSanitizers` to `ResetProxy` [#16256](https://github.com/Azure/azure-sdk-for-go/pull/16256)

## 0.8.2 (2021-11-11)

### Features Added
* Adding `RecordingOptions.RouteURL` to handle routing requests to the test proxy

### Bugs Fixed
* Adding recording sanitizers has no effect when running in `LiveMode`

## 0.8.1 (2021-10-21)

### Features Added
* Exports `RecordMode`, `PlaybackMode`, and `LiveMode` for determining test mode
* When running in `LiveMode` no traffic will be routed to the proxy and the `StartRecording`/`StopRecording` methods are no-ops.
* Adds the following sanitizers: `BodyKeySanitizer`, `BodyRegexSanitizer`, `ContinuationSanitizer`, `GeneralRegexSanitizer`, `HeaderRegexSanitizer`, `OAuthResponseSanitizer`, `RemoveHeaderSanitizer`, `URISanitizer`, `URISubscriptionIDSanitizer`
* Adds `ResetSanitizer` for removing all sanitizers in an active session.
* Add `GroupForReplace` to the `recording.RecordingOptions` option for use in sanitizers.

### Breaking Changes
* Removes `testing.T` parameter from `GetEnvVariable`
* Renames `IdHeader` to `IDHeader` and `UpstreamUriHeader` to `UpstreamURIHeader`
* Renames `StartRecording` to `Start` and `StopRecording` to `Stop`

## 0.8.0 (2021-10-20)
* Renamed log constant type and values to conform to guidelines.
* Added support for running tests in parallel
* Tests marked as LiveOnly will bypass the proxy

## 0.7.1 (2021-09-28)
* add `mock.NewTrackedCloser` to help test when `Close` is called
