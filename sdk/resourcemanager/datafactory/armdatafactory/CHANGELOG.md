# Release History

## 3.0.0 (2022-10-27)
### Breaking Changes

- Type of `SynapseSparkJobReference.ReferenceName` has been changed from `*string` to `interface{}`

### Features Added

- New field `WorkspaceResourceID` in struct `AzureSynapseArtifactsLinkedServiceTypeProperties`
- New field `DisablePublish` in struct `FactoryRepoConfiguration`
- New field `DisablePublish` in struct `FactoryGitHubConfiguration`
- New field `DisablePublish` in struct `FactoryVSTSConfiguration`
- New field `ScriptBlockExecutionTimeout` in struct `ScriptActivityTypeProperties`
- New field `PythonCodeReference` in struct `SynapseSparkJobActivityTypeProperties`
- New field `FilesV2` in struct `SynapseSparkJobActivityTypeProperties`


## 2.0.0 (2022-10-10)
### Breaking Changes

- Type of `SQLMISource.StoredProcedureParameters` has been changed from `map[string]*StoredProcedureParameter` to `interface{}`
- Type of `SQLMISink.StoredProcedureParameters` has been changed from `map[string]*StoredProcedureParameter` to `interface{}`
- Type of `AzureSQLSink.StoredProcedureParameters` has been changed from `map[string]*StoredProcedureParameter` to `interface{}`
- Type of `SQLServerSource.StoredProcedureParameters` has been changed from `map[string]*StoredProcedureParameter` to `interface{}`
- Type of `SQLServerSink.StoredProcedureParameters` has been changed from `map[string]*StoredProcedureParameter` to `interface{}`
- Type of `AzureSQLSource.StoredProcedureParameters` has been changed from `map[string]*StoredProcedureParameter` to `interface{}`
- Type of `SQLSink.StoredProcedureParameters` has been changed from `map[string]*StoredProcedureParameter` to `interface{}`
- Type of `SQLSource.StoredProcedureParameters` has been changed from `map[string]*StoredProcedureParameter` to `interface{}`
- Type of `AmazonRdsForSQLServerSource.StoredProcedureParameters` has been changed from `map[string]*StoredProcedureParameter` to `interface{}`


## 1.3.0 (2022-09-07)
### Features Added

- New const `NotebookParameterTypeBool`
- New const `NotebookReferenceTypeNotebookReference`
- New const `NotebookParameterTypeString`
- New const `SparkJobReferenceTypeSparkJobDefinitionReference`
- New const `NotebookParameterTypeInt`
- New const `BigDataPoolReferenceTypeBigDataPoolReference`
- New const `NotebookParameterTypeFloat`
- New type alias `NotebookParameterType`
- New type alias `SparkJobReferenceType`
- New type alias `NotebookReferenceType`
- New type alias `BigDataPoolReferenceType`
- New function `*AzureSynapseArtifactsLinkedService.GetLinkedService() *LinkedService`
- New function `PossibleBigDataPoolReferenceTypeValues() []BigDataPoolReferenceType`
- New function `PossibleNotebookParameterTypeValues() []NotebookParameterType`
- New function `*SynapseSparkJobDefinitionActivity.GetExecutionActivity() *ExecutionActivity`
- New function `*GoogleSheetsLinkedService.GetLinkedService() *LinkedService`
- New function `*SynapseNotebookActivity.GetExecutionActivity() *ExecutionActivity`
- New function `PossibleNotebookReferenceTypeValues() []NotebookReferenceType`
- New function `PossibleSparkJobReferenceTypeValues() []SparkJobReferenceType`
- New function `*SynapseNotebookActivity.GetActivity() *Activity`
- New function `*SynapseSparkJobDefinitionActivity.GetActivity() *Activity`
- New struct `AzureSynapseArtifactsLinkedService`
- New struct `AzureSynapseArtifactsLinkedServiceTypeProperties`
- New struct `BigDataPoolParametrizationReference`
- New struct `GoogleSheetsLinkedService`
- New struct `GoogleSheetsLinkedServiceTypeProperties`
- New struct `NotebookParameter`
- New struct `SynapseNotebookActivity`
- New struct `SynapseNotebookActivityTypeProperties`
- New struct `SynapseNotebookReference`
- New struct `SynapseSparkJobActivityTypeProperties`
- New struct `SynapseSparkJobDefinitionActivity`
- New struct `SynapseSparkJobReference`


## 1.2.0 (2022-06-15)
### Features Added

- New field `ClientSecret` in struct `RestServiceLinkedServiceTypeProperties`
- New field `Resource` in struct `RestServiceLinkedServiceTypeProperties`
- New field `Scope` in struct `RestServiceLinkedServiceTypeProperties`
- New field `TokenEndpoint` in struct `RestServiceLinkedServiceTypeProperties`
- New field `ClientID` in struct `RestServiceLinkedServiceTypeProperties`


## 1.1.0 (2022-05-30)
### Features Added

- New function `GlobalParameterResource.MarshalJSON() ([]byte, error)`
- New struct `GlobalParameterListResponse`
- New struct `GlobalParameterResource`
- New struct `GlobalParametersClientCreateOrUpdateOptions`
- New struct `GlobalParametersClientCreateOrUpdateResponse`
- New struct `GlobalParametersClientDeleteOptions`
- New struct `GlobalParametersClientDeleteResponse`
- New struct `GlobalParametersClientGetOptions`
- New struct `GlobalParametersClientGetResponse`
- New struct `GlobalParametersClientListByFactoryOptions`
- New struct `GlobalParametersClientListByFactoryResponse`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datafactory/armdatafactory` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).## 1.2.0 (2022-06-15)
### Features Added

- New field `ClientSecret` in struct `RestServiceLinkedServiceTypeProperties`
- New field `Resource` in struct `RestServiceLinkedServiceTypeProperties`
- New field `Scope` in struct `RestServiceLinkedServiceTypeProperties`
- New field `TokenEndpoint` in struct `RestServiceLinkedServiceTypeProperties`
- New field `ClientID` in struct `RestServiceLinkedServiceTypeProperties`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datafactory/armdatafactory` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).