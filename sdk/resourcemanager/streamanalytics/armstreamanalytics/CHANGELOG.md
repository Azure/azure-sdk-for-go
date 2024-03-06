# Release History

## 2.0.0-beta.1 (2024-01-26)
### Breaking Changes

- Type of `OutputProperties.SizeWindow` has been changed from `*float32` to `*int32`
- Function `*AzureMachineLearningWebServiceFunctionBinding.GetFunctionBinding` has been removed
- Function `*AzureMachineLearningWebServiceFunctionRetrieveDefaultDefinitionParameters.GetFunctionRetrieveDefaultDefinitionParameters` has been removed
- Struct `AzureMachineLearningWebServiceFunctionBinding` has been removed
- Struct `AzureMachineLearningWebServiceFunctionBindingProperties` has been removed
- Struct `AzureMachineLearningWebServiceFunctionBindingRetrievalProperties` has been removed
- Struct `AzureMachineLearningWebServiceFunctionRetrieveDefaultDefinitionParameters` has been removed
- Struct `AzureMachineLearningWebServiceInputColumn` has been removed
- Struct `AzureMachineLearningWebServiceInputs` has been removed
- Struct `AzureMachineLearningWebServiceOutputColumn` has been removed
- Field `Table` of struct `AzureSQLReferenceInputDataSourceProperties` has been removed

### Features Added

- New value `EventSerializationTypeCustomClr`, `EventSerializationTypeDelta` added to enum type `EventSerializationType`
- New enum type `BlobWriteMode` with values `BlobWriteModeAppend`, `BlobWriteModeOnce`
- New enum type `EventGridEventSchemaType` with values `EventGridEventSchemaTypeCloudEventSchema`, `EventGridEventSchemaTypeEventGridEventSchema`
- New enum type `InputWatermarkMode` with values `InputWatermarkModeNone`, `InputWatermarkModeReadWatermark`
- New enum type `OutputWatermarkMode` with values `OutputWatermarkModeNone`, `OutputWatermarkModeSendCurrentPartitionWatermark`, `OutputWatermarkModeSendLowestWatermarkAcrossPartitions`
- New enum type `QueryTestingResultStatus` with values `QueryTestingResultStatusCompilerError`, `QueryTestingResultStatusRuntimeError`, `QueryTestingResultStatusStarted`, `QueryTestingResultStatusSuccess`, `QueryTestingResultStatusTimeout`, `QueryTestingResultStatusUnknownError`
- New enum type `ResourceType` with values `ResourceTypeMicrosoftStreamAnalyticsStreamingjobs`
- New enum type `SKUCapacityScaleType` with values `SKUCapacityScaleTypeAutomatic`, `SKUCapacityScaleTypeManual`, `SKUCapacityScaleTypeNone`
- New enum type `SampleInputResultStatus` with values `SampleInputResultStatusErrorConnectingToInput`, `SampleInputResultStatusNoEventsFoundInRange`, `SampleInputResultStatusReadAllEventsInRange`
- New enum type `TestDatasourceResultStatus` with values `TestDatasourceResultStatusTestFailed`, `TestDatasourceResultStatusTestSucceeded`
- New enum type `UpdatableUdfRefreshType` with values `UpdatableUdfRefreshTypeBlocking`, `UpdatableUdfRefreshTypeNonblocking`
- New enum type `UpdateMode` with values `UpdateModeRefreshable`, `UpdateModeStatic`
- New function `*AzureDataExplorerOutputDataSource.GetOutputDataSource() *OutputDataSource`
- New function `*AzureMachineLearningServiceFunctionBinding.GetFunctionBinding() *FunctionBinding`
- New function `*AzureMachineLearningServiceFunctionRetrieveDefaultDefinitionParameters.GetFunctionRetrieveDefaultDefinitionParameters() *FunctionRetrieveDefaultDefinitionParameters`
- New function `*AzureMachineLearningStudioFunctionBinding.GetFunctionBinding() *FunctionBinding`
- New function `*AzureMachineLearningStudioFunctionRetrieveDefaultDefinitionParameters.GetFunctionRetrieveDefaultDefinitionParameters() *FunctionRetrieveDefaultDefinitionParameters`
- New function `*CSharpFunctionBinding.GetFunctionBinding() *FunctionBinding`
- New function `*CSharpFunctionRetrieveDefaultDefinitionParameters.GetFunctionRetrieveDefaultDefinitionParameters() *FunctionRetrieveDefaultDefinitionParameters`
- New function `*ClientFactory.NewSKUClient() *SKUClient`
- New function `*CustomClrSerialization.GetSerialization() *Serialization`
- New function `*DeltaSerialization.GetSerialization() *Serialization`
- New function `*EventGridStreamInputDataSource.GetStreamInputDataSource() *StreamInputDataSource`
- New function `*FileReferenceInputDataSource.GetReferenceInputDataSource() *ReferenceInputDataSource`
- New function `*GatewayMessageBusOutputDataSource.GetOutputDataSource() *OutputDataSource`
- New function `*GatewayMessageBusStreamInputDataSource.GetStreamInputDataSource() *StreamInputDataSource`
- New function `*PostgreSQLOutputDataSource.GetOutputDataSource() *OutputDataSource`
- New function `*RawOutputDatasource.GetOutputDataSource() *OutputDataSource`
- New function `*RawReferenceInputDataSource.GetReferenceInputDataSource() *ReferenceInputDataSource`
- New function `*RawStreamInputDataSource.GetStreamInputDataSource() *StreamInputDataSource`
- New function `NewSKUClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SKUClient, error)`
- New function `*SKUClient.NewListPager(string, string, *SKUClientListOptions) *runtime.Pager[SKUClientListResponse]`
- New function `*SubscriptionsClient.CompileQuery(context.Context, string, CompileQuery, *SubscriptionsClientCompileQueryOptions) (SubscriptionsClientCompileQueryResponse, error)`
- New function `*SubscriptionsClient.BeginSampleInput(context.Context, string, SampleInput, *SubscriptionsClientBeginSampleInputOptions) (*runtime.Poller[SubscriptionsClientSampleInputResponse], error)`
- New function `*SubscriptionsClient.BeginTestInput(context.Context, string, TestInput, *SubscriptionsClientBeginTestInputOptions) (*runtime.Poller[SubscriptionsClientTestInputResponse], error)`
- New function `*SubscriptionsClient.BeginTestOutput(context.Context, string, TestOutput, *SubscriptionsClientBeginTestOutputOptions) (*runtime.Poller[SubscriptionsClientTestOutputResponse], error)`
- New function `*SubscriptionsClient.BeginTestQuery(context.Context, string, TestQuery, *SubscriptionsClientBeginTestQueryOptions) (*runtime.Poller[SubscriptionsClientTestQueryResponse], error)`
- New struct `AzureDataExplorerOutputDataSource`
- New struct `AzureDataExplorerOutputDataSourceProperties`
- New struct `AzureMachineLearningServiceFunctionBinding`
- New struct `AzureMachineLearningServiceFunctionBindingProperties`
- New struct `AzureMachineLearningServiceFunctionBindingRetrievalProperties`
- New struct `AzureMachineLearningServiceFunctionRetrieveDefaultDefinitionParameters`
- New struct `AzureMachineLearningServiceInputColumn`
- New struct `AzureMachineLearningServiceInputs`
- New struct `AzureMachineLearningServiceOutputColumn`
- New struct `AzureMachineLearningStudioFunctionBinding`
- New struct `AzureMachineLearningStudioFunctionBindingProperties`
- New struct `AzureMachineLearningStudioFunctionBindingRetrievalProperties`
- New struct `AzureMachineLearningStudioFunctionRetrieveDefaultDefinitionParameters`
- New struct `AzureMachineLearningStudioInputColumn`
- New struct `AzureMachineLearningStudioInputs`
- New struct `AzureMachineLearningStudioOutputColumn`
- New struct `CSharpFunctionBinding`
- New struct `CSharpFunctionBindingProperties`
- New struct `CSharpFunctionBindingRetrievalProperties`
- New struct `CSharpFunctionRetrieveDefaultDefinitionParameters`
- New struct `CompileQuery`
- New struct `CustomClrSerialization`
- New struct `CustomClrSerializationProperties`
- New struct `DeltaSerialization`
- New struct `DeltaSerializationProperties`
- New struct `EventGridStreamInputDataSource`
- New struct `EventGridStreamInputDataSourceProperties`
- New struct `External`
- New struct `FileReferenceInputDataSource`
- New struct `FileReferenceInputDataSourceProperties`
- New struct `GatewayMessageBusOutputDataSource`
- New struct `GatewayMessageBusOutputDataSourceProperties`
- New struct `GatewayMessageBusSourceProperties`
- New struct `GatewayMessageBusStreamInputDataSource`
- New struct `GatewayMessageBusStreamInputDataSourceProperties`
- New struct `GetStreamingJobSKUResult`
- New struct `GetStreamingJobSKUResultSKU`
- New struct `GetStreamingJobSKUResults`
- New struct `InputWatermarkProperties`
- New struct `LastOutputEventTimestamp`
- New struct `OutputWatermarkProperties`
- New struct `PostgreSQLDataSourceProperties`
- New struct `PostgreSQLOutputDataSource`
- New struct `PostgreSQLOutputDataSourceProperties`
- New struct `QueryCompilationError`
- New struct `QueryCompilationResult`
- New struct `QueryFunction`
- New struct `QueryInput`
- New struct `QueryTestingResult`
- New struct `RawInputDatasourceProperties`
- New struct `RawOutputDatasource`
- New struct `RawOutputDatasourceProperties`
- New struct `RawReferenceInputDataSource`
- New struct `RawStreamInputDataSource`
- New struct `RefreshConfiguration`
- New struct `SKUCapacity`
- New struct `SampleInput`
- New struct `SampleInputResult`
- New struct `TestDatasourceResult`
- New struct `TestInput`
- New struct `TestOutput`
- New struct `TestQuery`
- New struct `TestQueryDiagnostics`
- New field `AuthenticationMode` in struct `AzureSQLReferenceInputDataSourceProperties`
- New field `AuthenticationMode` in struct `AzureSynapseDataSourceProperties`
- New field `AuthenticationMode` in struct `AzureSynapseOutputDataSourceProperties`
- New field `BlobPathPrefix`, `BlobWriteMode` in struct `BlobOutputDataSourceProperties`
- New field `BlobName`, `DeltaPathPattern`, `DeltaSnapshotRefreshRate`, `FullSnapshotRefreshRate`, `SourcePartitionCount` in struct `BlobReferenceInputDataSourceProperties`
- New field `AuthenticationMode` in struct `DocumentDbOutputDataSourceProperties`
- New field `PartitionCount` in struct `EventHubDataSourceProperties`
- New field `PartitionCount` in struct `EventHubOutputDataSourceProperties`
- New field `PartitionCount`, `PrefetchCount` in struct `EventHubStreamInputDataSourceProperties`
- New field `UserAssignedIdentities` in struct `Identity`
- New field `WatermarkSettings` in struct `InputProperties`
- New field `LastOutputEventTimestamps`, `WatermarkSettings` in struct `OutputProperties`
- New field `WatermarkSettings` in struct `ReferenceInputProperties`
- New field `Capacity` in struct `SKU`
- New field `AuthenticationMode` in struct `StorageAccount`
- New field `WatermarkSettings` in struct `StreamInputProperties`
- New field `SKU` in struct `StreamingJob`
- New field `Externals` in struct `StreamingJobProperties`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/streamanalytics/armstreamanalytics` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).