# Release History

## 0.8.1 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed

### Other Changes

## 0.8.0 (2021-10-20)
* Renamed log constant type and values to conform to guidelines.
* Exports `RecordMode`, `PlaybackMode`, and `LiveMode` for determining test mode
* When running in `LiveMode` no traffic will be routed to the proxy and the `StartRecording`/`StopRecording` methods are no-ops.
* Added support for running tests in parallel

## 0.7.1 (2021-09-28)
* add `mock.NewTrackedCloser` to help test when `Close` is called
