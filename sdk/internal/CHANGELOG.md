# Release History

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
