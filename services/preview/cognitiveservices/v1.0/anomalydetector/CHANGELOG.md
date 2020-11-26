
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Const `Weekly` type has been changed from `Granularity` to `TimeGranularity`
- Const `Yearly` type has been changed from `Granularity` to `TimeGranularity`
- Const `Daily` type has been changed from `Granularity` to `TimeGranularity`
- Const `Hourly` type has been changed from `Granularity` to `TimeGranularity`
- Const `Monthly` type has been changed from `Granularity` to `TimeGranularity`
- Type of `ChangePointDetectRequest.Series` has been changed from `*[]Point` to `*[]TimeSeriesPoint`
- Type of `ChangePointDetectRequest.Granularity` has been changed from `Granularity` to `TimeGranularity`
- Const `Secondly` has been removed
- Const `Minutely` has been removed
- Function `BaseClient.LastDetectSender` has been removed
- Function `BaseClient.EntireDetectResponder` has been removed
- Function `BaseClient.ChangePointDetectSender` has been removed
- Function `BaseClient.ChangePointDetectResponder` has been removed
- Function `BaseClient.LastDetect` has been removed
- Function `BaseClient.LastDetectResponder` has been removed
- Function `BaseClient.ChangePointDetectPreparer` has been removed
- Function `PossibleGranularityValues` has been removed
- Function `BaseClient.EntireDetectSender` has been removed
- Function `BaseClient.LastDetectPreparer` has been removed
- Function `BaseClient.EntireDetectPreparer` has been removed
- Function `BaseClient.ChangePointDetect` has been removed
- Function `BaseClient.EntireDetect` has been removed
- Struct `APIError` has been removed
- Struct `Point` has been removed
- Struct `Request` has been removed

## New Content

- Const `PerMinute` is added
- Const `PerSecond` is added
- Function `BaseClient.DetectEntireSeriesSender(*http.Request) (*http.Response,error)` is added
- Function `BaseClient.DetectLastPointResponder(*http.Response) (LastDetectResponse,error)` is added
- Function `PossibleTimeGranularityValues() []TimeGranularity` is added
- Function `BaseClient.DetectEntireSeries(context.Context,DetectRequest) (EntireDetectResponse,error)` is added
- Function `BaseClient.DetectLastPoint(context.Context,DetectRequest) (LastDetectResponse,error)` is added
- Function `BaseClient.DetectLastPointSender(*http.Request) (*http.Response,error)` is added
- Function `BaseClient.DetectChangePointSender(*http.Request) (*http.Response,error)` is added
- Function `BaseClient.DetectLastPointPreparer(context.Context,DetectRequest) (*http.Request,error)` is added
- Function `BaseClient.DetectChangePoint(context.Context,ChangePointDetectRequest) (ChangePointDetectResponse,error)` is added
- Function `BaseClient.DetectEntireSeriesPreparer(context.Context,DetectRequest) (*http.Request,error)` is added
- Function `BaseClient.DetectEntireSeriesResponder(*http.Response) (EntireDetectResponse,error)` is added
- Function `BaseClient.DetectChangePointResponder(*http.Response) (ChangePointDetectResponse,error)` is added
- Function `BaseClient.DetectChangePointPreparer(context.Context,ChangePointDetectRequest) (*http.Request,error)` is added
- Struct `DetectRequest` is added
- Struct `Error` is added
- Struct `TimeSeriesPoint` is added

