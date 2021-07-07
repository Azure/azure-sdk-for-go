# Unreleased

## Additive Changes

### New Constants

1. QueryTestingResultStatus.CompilerError
1. QueryTestingResultStatus.RuntimeError
1. QueryTestingResultStatus.Started
1. QueryTestingResultStatus.Success
1. QueryTestingResultStatus.Timeout
1. QueryTestingResultStatus.UnknownError
1. SampleInputResultStatus.ErrorConnectingToInput
1. SampleInputResultStatus.NoEventsFoundInRange
1. SampleInputResultStatus.ReadAllEventsInRange
1. TestDatasourceResultStatus.TestFailed
1. TestDatasourceResultStatus.TestSucceeded
1. TypeBasicOutputDataSource.TypeRaw
1. TypeBasicReferenceInputDataSource.TypeBasicReferenceInputDataSourceTypeRaw
1. TypeBasicStreamInputDataSource.TypeBasicStreamInputDataSourceTypeRaw

### New Funcs

1. *RawOutputDatasource.UnmarshalJSON([]byte) error
1. *RawReferenceInputDataSource.UnmarshalJSON([]byte) error
1. *RawStreamInputDataSource.UnmarshalJSON([]byte) error
1. *SubscriptionsSampleInputMethodFuture.UnmarshalJSON([]byte) error
1. *SubscriptionsTestInputMethodFuture.UnmarshalJSON([]byte) error
1. *SubscriptionsTestOutputMethodFuture.UnmarshalJSON([]byte) error
1. *SubscriptionsTestQueryMethodFuture.UnmarshalJSON([]byte) error
1. AzureDataLakeStoreOutputDataSource.AsRawOutputDatasource() (*RawOutputDatasource, bool)
1. AzureFunctionOutputDataSource.AsRawOutputDatasource() (*RawOutputDatasource, bool)
1. AzureSQLDatabaseOutputDataSource.AsRawOutputDatasource() (*RawOutputDatasource, bool)
1. AzureSQLReferenceInputDataSource.AsRawReferenceInputDataSource() (*RawReferenceInputDataSource, bool)
1. AzureSynapseOutputDataSource.AsRawOutputDatasource() (*RawOutputDatasource, bool)
1. AzureTableOutputDataSource.AsRawOutputDatasource() (*RawOutputDatasource, bool)
1. BlobOutputDataSource.AsRawOutputDatasource() (*RawOutputDatasource, bool)
1. BlobReferenceInputDataSource.AsRawReferenceInputDataSource() (*RawReferenceInputDataSource, bool)
1. BlobStreamInputDataSource.AsRawStreamInputDataSource() (*RawStreamInputDataSource, bool)
1. DocumentDbOutputDataSource.AsRawOutputDatasource() (*RawOutputDatasource, bool)
1. EventHubOutputDataSource.AsRawOutputDatasource() (*RawOutputDatasource, bool)
1. EventHubStreamInputDataSource.AsRawStreamInputDataSource() (*RawStreamInputDataSource, bool)
1. EventHubV2OutputDataSource.AsRawOutputDatasource() (*RawOutputDatasource, bool)
1. EventHubV2StreamInputDataSource.AsRawStreamInputDataSource() (*RawStreamInputDataSource, bool)
1. IoTHubStreamInputDataSource.AsRawStreamInputDataSource() (*RawStreamInputDataSource, bool)
1. OutputDataSource.AsRawOutputDatasource() (*RawOutputDatasource, bool)
1. PossibleQueryTestingResultStatusValues() []QueryTestingResultStatus
1. PossibleSampleInputResultStatusValues() []SampleInputResultStatus
1. PossibleTestDatasourceResultStatusValues() []TestDatasourceResultStatus
1. PowerBIOutputDataSource.AsRawOutputDatasource() (*RawOutputDatasource, bool)
1. QueryCompilationError.MarshalJSON() ([]byte, error)
1. QueryCompilationResult.MarshalJSON() ([]byte, error)
1. QueryTestingResult.MarshalJSON() ([]byte, error)
1. RawOutputDatasource.AsAzureDataLakeStoreOutputDataSource() (*AzureDataLakeStoreOutputDataSource, bool)
1. RawOutputDatasource.AsAzureFunctionOutputDataSource() (*AzureFunctionOutputDataSource, bool)
1. RawOutputDatasource.AsAzureSQLDatabaseOutputDataSource() (*AzureSQLDatabaseOutputDataSource, bool)
1. RawOutputDatasource.AsAzureSynapseOutputDataSource() (*AzureSynapseOutputDataSource, bool)
1. RawOutputDatasource.AsAzureTableOutputDataSource() (*AzureTableOutputDataSource, bool)
1. RawOutputDatasource.AsBasicOutputDataSource() (BasicOutputDataSource, bool)
1. RawOutputDatasource.AsBlobOutputDataSource() (*BlobOutputDataSource, bool)
1. RawOutputDatasource.AsDocumentDbOutputDataSource() (*DocumentDbOutputDataSource, bool)
1. RawOutputDatasource.AsEventHubOutputDataSource() (*EventHubOutputDataSource, bool)
1. RawOutputDatasource.AsEventHubV2OutputDataSource() (*EventHubV2OutputDataSource, bool)
1. RawOutputDatasource.AsOutputDataSource() (*OutputDataSource, bool)
1. RawOutputDatasource.AsPowerBIOutputDataSource() (*PowerBIOutputDataSource, bool)
1. RawOutputDatasource.AsRawOutputDatasource() (*RawOutputDatasource, bool)
1. RawOutputDatasource.AsServiceBusQueueOutputDataSource() (*ServiceBusQueueOutputDataSource, bool)
1. RawOutputDatasource.AsServiceBusTopicOutputDataSource() (*ServiceBusTopicOutputDataSource, bool)
1. RawOutputDatasource.MarshalJSON() ([]byte, error)
1. RawReferenceInputDataSource.AsAzureSQLReferenceInputDataSource() (*AzureSQLReferenceInputDataSource, bool)
1. RawReferenceInputDataSource.AsBasicReferenceInputDataSource() (BasicReferenceInputDataSource, bool)
1. RawReferenceInputDataSource.AsBlobReferenceInputDataSource() (*BlobReferenceInputDataSource, bool)
1. RawReferenceInputDataSource.AsRawReferenceInputDataSource() (*RawReferenceInputDataSource, bool)
1. RawReferenceInputDataSource.AsReferenceInputDataSource() (*ReferenceInputDataSource, bool)
1. RawReferenceInputDataSource.MarshalJSON() ([]byte, error)
1. RawStreamInputDataSource.AsBasicStreamInputDataSource() (BasicStreamInputDataSource, bool)
1. RawStreamInputDataSource.AsBlobStreamInputDataSource() (*BlobStreamInputDataSource, bool)
1. RawStreamInputDataSource.AsEventHubStreamInputDataSource() (*EventHubStreamInputDataSource, bool)
1. RawStreamInputDataSource.AsEventHubV2StreamInputDataSource() (*EventHubV2StreamInputDataSource, bool)
1. RawStreamInputDataSource.AsIoTHubStreamInputDataSource() (*IoTHubStreamInputDataSource, bool)
1. RawStreamInputDataSource.AsRawStreamInputDataSource() (*RawStreamInputDataSource, bool)
1. RawStreamInputDataSource.AsStreamInputDataSource() (*StreamInputDataSource, bool)
1. RawStreamInputDataSource.MarshalJSON() ([]byte, error)
1. ReferenceInputDataSource.AsRawReferenceInputDataSource() (*RawReferenceInputDataSource, bool)
1. SampleInputResult.MarshalJSON() ([]byte, error)
1. ServiceBusQueueOutputDataSource.AsRawOutputDatasource() (*RawOutputDatasource, bool)
1. ServiceBusTopicOutputDataSource.AsRawOutputDatasource() (*RawOutputDatasource, bool)
1. StreamInputDataSource.AsRawStreamInputDataSource() (*RawStreamInputDataSource, bool)
1. SubscriptionsClient.CompileQueryMethod(context.Context, CompileQuery, string) (QueryCompilationResult, error)
1. SubscriptionsClient.CompileQueryMethodPreparer(context.Context, CompileQuery, string) (*http.Request, error)
1. SubscriptionsClient.CompileQueryMethodResponder(*http.Response) (QueryCompilationResult, error)
1. SubscriptionsClient.CompileQueryMethodSender(*http.Request) (*http.Response, error)
1. SubscriptionsClient.SampleInputMethod(context.Context, SampleInput, string) (SubscriptionsSampleInputMethodFuture, error)
1. SubscriptionsClient.SampleInputMethodPreparer(context.Context, SampleInput, string) (*http.Request, error)
1. SubscriptionsClient.SampleInputMethodResponder(*http.Response) (SampleInputResult, error)
1. SubscriptionsClient.SampleInputMethodSender(*http.Request) (SubscriptionsSampleInputMethodFuture, error)
1. SubscriptionsClient.TestInputMethod(context.Context, TestInput, string) (SubscriptionsTestInputMethodFuture, error)
1. SubscriptionsClient.TestInputMethodPreparer(context.Context, TestInput, string) (*http.Request, error)
1. SubscriptionsClient.TestInputMethodResponder(*http.Response) (TestDatasourceResult, error)
1. SubscriptionsClient.TestInputMethodSender(*http.Request) (SubscriptionsTestInputMethodFuture, error)
1. SubscriptionsClient.TestOutputMethod(context.Context, TestOutput, string) (SubscriptionsTestOutputMethodFuture, error)
1. SubscriptionsClient.TestOutputMethodPreparer(context.Context, TestOutput, string) (*http.Request, error)
1. SubscriptionsClient.TestOutputMethodResponder(*http.Response) (TestDatasourceResult, error)
1. SubscriptionsClient.TestOutputMethodSender(*http.Request) (SubscriptionsTestOutputMethodFuture, error)
1. SubscriptionsClient.TestQueryMethod(context.Context, TestQuery, string) (SubscriptionsTestQueryMethodFuture, error)
1. SubscriptionsClient.TestQueryMethodPreparer(context.Context, TestQuery, string) (*http.Request, error)
1. SubscriptionsClient.TestQueryMethodResponder(*http.Response) (QueryTestingResult, error)
1. SubscriptionsClient.TestQueryMethodSender(*http.Request) (SubscriptionsTestQueryMethodFuture, error)

### Struct Changes

#### New Structs

1. CompileQuery
1. QueryCompilationError
1. QueryCompilationResult
1. QueryFunction
1. QueryInput
1. QueryTestingResult
1. RawInputDatasourceProperties
1. RawOutputDatasource
1. RawOutputDatasourceProperties
1. RawReferenceInputDataSource
1. RawStreamInputDataSource
1. SampleInput
1. SampleInputResult
1. SubscriptionsSampleInputMethodFuture
1. SubscriptionsTestInputMethodFuture
1. SubscriptionsTestOutputMethodFuture
1. SubscriptionsTestQueryMethodFuture
1. TestDatasourceResult
1. TestInput
1. TestOutput
1. TestQuery
1. TestQueryDiagnostics
