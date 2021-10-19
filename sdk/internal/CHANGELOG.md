# Release History

## 0.7.2 (Unreleased)
* Exports `RecordMode`, `PlaybackMode`, and `LiveMode` for determining test mode
* Renames `StartRecording` to `Start` and `StopRecording` to `Stop`
* When running in `LiveMode` no traffic will be routed to the proxy and the `StartRecording`/`StopRecording` methods are no-ops.
* Adds the following sanitizers: `BodyKeySanitizer`, `BodyRegexSanitizer`, `ContinuationSanitizer`, `GeneralRegexSanitizer`, `HeaderRegexSanitizer`, `OAuthResponseSanitizer`, `RemoveHeaderSanitizer`, `URISanitizer`, `URISubscriptionIDSanitizer`
* Adds `ResetSanitizer` for removing all sanitizers in an active session.
* Renames `IdHeader` to `IDHeader` and `UpstreamUriHeader` to `UpstreamURIHeader`
* Add `GroupForReplace` to the `recording.RecordingOptions` option for use in sanitizers.
* Removes `testing.T` parameter from `GetEnvVariable`

## 0.7.1 (2021-09-28)
* add `mock.NewTrackedCloser` to help test when `Close` is called