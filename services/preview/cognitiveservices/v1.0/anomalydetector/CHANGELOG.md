Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Const `Yearly` type has been changed from `Granularity` to `TimeGranularity`
- Const `Daily` type has been changed from `Granularity` to `TimeGranularity`
- Const `Hourly` type has been changed from `Granularity` to `TimeGranularity`
- Const `Monthly` type has been changed from `Granularity` to `TimeGranularity`
- Const `Weekly` type has been changed from `Granularity` to `TimeGranularity`
- Type of `ChangePointDetectRequest.Granularity` has been changed from `Granularity` to `TimeGranularity`
- Type of `ChangePointDetectRequest.Series` has been changed from `*[]Point` to `*[]TimeSeriesPoint`
- Const `Secondly` has been removed
- Const `Minutely` has been removed
- Function `BaseClient.LastDetect` has been removed
- Function `BaseClient.LastDetectPreparer` has been removed
- Function `BaseClient.ChangePointDetect` has been removed
- Function `BaseClient.ChangePointDetectPreparer` has been removed
- Function `BaseClient.EntireDetect` has been removed
- Function `BaseClient.LastDetectResponder` has been removed
- Function `PossibleGranularityValues` has been removed
- Function `BaseClient.ChangePointDetectSender` has been removed
- Function `BaseClient.ChangePointDetectResponder` has been removed
- Function `BaseClient.LastDetectSender` has been removed
- Function `BaseClient.EntireDetectPreparer` has been removed
- Function `BaseClient.EntireDetectSender` has been removed
- Function `BaseClient.EntireDetectResponder` has been removed
- Struct `APIError` has been removed
- Struct `Point` has been removed
- Struct `Request` has been removed

## New Content

- New const `PerMinute`
- New const `PerSecond`
- New function `BaseClient.DetectEntireSeriesPreparer(context.Context, DetectRequest) (*http.Request, error)`
- New function `BaseClient.DetectLastPointPreparer(context.Context, DetectRequest) (*http.Request, error)`
- New function `BaseClient.DetectChangePointSender(*http.Request) (*http.Response, error)`
- New function `BaseClient.DetectLastPointSender(*http.Request) (*http.Response, error)`
- New function `BaseClient.DetectChangePoint(context.Context, ChangePointDetectRequest) (ChangePointDetectResponse, error)`
- New function `BaseClient.DetectEntireSeriesSender(*http.Request) (*http.Response, error)`
- New function `PossibleTimeGranularityValues() []TimeGranularity`
- New function `BaseClient.DetectChangePointResponder(*http.Response) (ChangePointDetectResponse, error)`
- New function `BaseClient.DetectEntireSeriesResponder(*http.Response) (EntireDetectResponse, error)`
- New function `BaseClient.DetectChangePointPreparer(context.Context, ChangePointDetectRequest) (*http.Request, error)`
- New function `BaseClient.DetectLastPoint(context.Context, DetectRequest) (LastDetectResponse, error)`
- New function `BaseClient.DetectLastPointResponder(*http.Response) (LastDetectResponse, error)`
- New function `BaseClient.DetectEntireSeries(context.Context, DetectRequest) (EntireDetectResponse, error)`
- New struct `DetectRequest`
- New struct `Error`
- New struct `TimeSeriesPoint`
