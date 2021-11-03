# Release History

## 0.8.2 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed

### Other Changes

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
