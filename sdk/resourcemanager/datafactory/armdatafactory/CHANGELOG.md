# Release History

## 2.0.0 (2022-09-07)
### Breaking Changes

- Type of `LinkedServiceReference.Type` has been changed from `*LinkedServiceReferenceType` to `*Type`
- Const `LinkedServiceReferenceTypeLinkedServiceReference` has been removed
- Type alias `LinkedServiceReferenceType` has been removed
- Function `PossibleLinkedServiceReferenceTypeValues` has been removed

### Features Added

- New const `TypeLinkedServiceReference`
- New const `NotebookReferenceTypeNotebookReference`
- New const `NotebookParameterTypeInt`
- New const `SparkJobReferenceTypeSparkJobDefinitionReference`
- New const `NotebookParameterTypeBool`
- New const `NotebookParameterTypeString`
- New const `NotebookParameterTypeFloat`
- New const `BigDataPoolReferenceTypeBigDataPoolReference`
- New type alias `Type`
- New type alias `SparkJobReferenceType`
- New type alias `BigDataPoolReferenceType`
- New type alias `NotebookParameterType`
- New type alias `NotebookReferenceType`
- New function `*SynapseSparkJobDefinitionActivity.GetActivity() *Activity`
- New function `*SynapseNotebookActivity.GetExecutionActivity() *ExecutionActivity`
- New function `PossibleNotebookParameterTypeValues() []NotebookParameterType`
- New function `*GoogleSheetsLinkedService.GetLinkedService() *LinkedService`
- New function `PossibleBigDataPoolReferenceTypeValues() []BigDataPoolReferenceType`
- New function `PossibleSparkJobReferenceTypeValues() []SparkJobReferenceType`
- New function `*SynapseNotebookActivity.GetActivity() *Activity`
- New function `*AzureSynapseArtifactsLinkedService.GetLinkedService() *LinkedService`
- New function `PossibleTypeValues() []Type`
- New function `PossibleNotebookReferenceTypeValues() []NotebookReferenceType`
- New function `*SynapseSparkJobDefinitionActivity.GetExecutionActivity() *ExecutionActivity`
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