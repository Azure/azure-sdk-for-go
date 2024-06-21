# Release History

## 2.3.0 (2024-06-21)
### Features Added

- New enum type `AutoRunState` with values `AutoRunStateAutoRunDisabled`, `AutoRunStateAutoRunEnabled`
- New struct `ImageTemplateAutoRun`
- New field `AutoRun`, `ManagedResourceTags` in struct `ImageTemplateProperties`
- New field `VMProfile` in struct `ImageTemplateUpdateParametersProperties`
- New field `ContainerInstanceSubnetID` in struct `VirtualNetworkConfig`


## 2.2.0 (2023-12-22)
### Features Added

- New enum type `OnBuildError` with values `OnBuildErrorAbort`, `OnBuildErrorCleanup`
- New struct `ErrorAdditionalInfo`
- New struct `ErrorDetail`
- New struct `ErrorResponse`
- New struct `ImageTemplatePropertiesErrorHandling`
- New struct `ImageTemplateUpdateParametersProperties`
- New field `ErrorHandling` in struct `ImageTemplateProperties`
- New field `Properties` in struct `ImageTemplateUpdateParameters`


## 2.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.0.0 (2023-04-28)
### Breaking Changes

- Type of `ImageTemplateIdentity.UserAssignedIdentities` has been changed from `map[string]*ComponentsVrq145SchemasImagetemplateidentityPropertiesUserassignedidentitiesAdditionalproperties` to `map[string]*UserAssignedIdentity`
- Struct `ComponentsVrq145SchemasImagetemplateidentityPropertiesUserassignedidentitiesAdditionalproperties` has been removed

### Features Added

- New value `ProvisioningStateCanceled` added to enum type `ProvisioningState`
- New value `RunSubStateOptimizing` added to enum type `RunSubState`
- New value `SharedImageStorageAccountTypePremiumLRS` added to enum type `SharedImageStorageAccountType`
- New enum type `VMBootOptimizationState` with values `VMBootOptimizationStateDisabled`, `VMBootOptimizationStateEnabled`
- New function `*ClientFactory.NewTriggersClient() *TriggersClient`
- New function `*DistributeVersioner.GetDistributeVersioner() *DistributeVersioner`
- New function `*DistributeVersionerLatest.GetDistributeVersioner() *DistributeVersioner`
- New function `*DistributeVersionerSource.GetDistributeVersioner() *DistributeVersioner`
- New function `*ImageTemplateFileValidator.GetImageTemplateInVMValidator() *ImageTemplateInVMValidator`
- New function `*SourceImageTriggerProperties.GetTriggerProperties() *TriggerProperties`
- New function `*TriggerProperties.GetTriggerProperties() *TriggerProperties`
- New function `NewTriggersClient(string, azcore.TokenCredential, *arm.ClientOptions) (*TriggersClient, error)`
- New function `*TriggersClient.BeginCreateOrUpdate(context.Context, string, string, string, Trigger, *TriggersClientBeginCreateOrUpdateOptions) (*runtime.Poller[TriggersClientCreateOrUpdateResponse], error)`
- New function `*TriggersClient.BeginDelete(context.Context, string, string, string, *TriggersClientBeginDeleteOptions) (*runtime.Poller[TriggersClientDeleteResponse], error)`
- New function `*TriggersClient.Get(context.Context, string, string, string, *TriggersClientGetOptions) (TriggersClientGetResponse, error)`
- New function `*TriggersClient.NewListByImageTemplatePager(string, string, *TriggersClientListByImageTemplateOptions) *runtime.Pager[TriggersClientListByImageTemplateResponse]`
- New struct `DistributeVersionerLatest`
- New struct `DistributeVersionerSource`
- New struct `ImageTemplateFileValidator`
- New struct `ImageTemplatePropertiesOptimize`
- New struct `ImageTemplatePropertiesOptimizeVMBoot`
- New struct `SourceImageTriggerProperties`
- New struct `TargetRegion`
- New struct `Trigger`
- New struct `TriggerCollection`
- New struct `TriggerStatus`
- New struct `UserAssignedIdentity`
- New field `Optimize` in struct `ImageTemplateProperties`
- New field `TargetRegions` in struct `ImageTemplateSharedImageDistributor`
- New field `Versioning` in struct `ImageTemplateSharedImageDistributor`
- New field `ExactVersion` in struct `ImageTemplateSharedImageVersionSource`
- New field `URI` in struct `ImageTemplateVhdDistributor`


## 1.2.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.2.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.1.0 (2022-05-30)
### Features Added

- New const `ProvisioningErrorCodeBadValidatorType`
- New const `ProvisioningErrorCodeUnsupportedValidatorType`
- New const `RunSubStateValidating`
- New const `ProvisioningErrorCodeNoValidatorScript`
- New const `ProvisioningErrorCodeBadStagingResourceGroup`
- New function `ImageTemplatePowerShellValidator.MarshalJSON() ([]byte, error)`
- New function `*ImageTemplatePropertiesValidate.UnmarshalJSON([]byte) error`
- New function `*ImageTemplateShellValidator.UnmarshalJSON([]byte) error`
- New function `*ImageTemplatePowerShellValidator.GetImageTemplateInVMValidator() *ImageTemplateInVMValidator`
- New function `ImageTemplateShellValidator.MarshalJSON() ([]byte, error)`
- New function `*ImageTemplateShellValidator.GetImageTemplateInVMValidator() *ImageTemplateInVMValidator`
- New function `ImageTemplatePropertiesValidate.MarshalJSON() ([]byte, error)`
- New function `*ImageTemplateInVMValidator.GetImageTemplateInVMValidator() *ImageTemplateInVMValidator`
- New function `*ImageTemplatePowerShellValidator.UnmarshalJSON([]byte) error`
- New struct `ImageTemplateInVMValidator`
- New struct `ImageTemplatePowerShellValidator`
- New struct `ImageTemplatePropertiesValidate`
- New struct `ImageTemplateShellValidator`
- New struct `ProxyResource`
- New field `SystemData` in struct `TrackedResource`
- New field `SystemData` in struct `RunOutput`
- New field `ExactStagingResourceGroup` in struct `ImageTemplateProperties`
- New field `StagingResourceGroup` in struct `ImageTemplateProperties`
- New field `Validate` in struct `ImageTemplateProperties`
- New field `SystemData` in struct `Resource`


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/virtualmachineimagebuilder/armvirtualmachineimagebuilder` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).