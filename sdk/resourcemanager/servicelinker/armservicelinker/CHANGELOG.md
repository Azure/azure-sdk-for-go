# Release History

## 2.0.0-beta.1 (2024-03-22)
### Breaking Changes

- Struct `LinkerList` has been removed
- Struct `SourceConfigurationResult` has been removed
- Field `SourceConfigurationResult` of struct `LinkerClientListConfigurationsResponse` has been removed
- Field `LinkerList` of struct `LinkerClientListResponse` has been removed

### Features Added

- New value `ActionTypeEnable`, `ActionTypeOptOut` added to enum type `ActionType`
- New value `AuthTypeAccessKey`, `AuthTypeEasyAuthMicrosoftEntraID`, `AuthTypeUserAccount` added to enum type `AuthType`
- New value `ClientTypeDapr`, `ClientTypeJmsSpringBoot`, `ClientTypeKafkaSpringBoot` added to enum type `ClientType`
- New value `TargetServiceTypeSelfHostedServer` added to enum type `TargetServiceType`
- New enum type `AccessKeyPermissions` with values `AccessKeyPermissionsListen`, `AccessKeyPermissionsManage`, `AccessKeyPermissionsRead`, `AccessKeyPermissionsSend`, `AccessKeyPermissionsWrite`
- New enum type `AllowType` with values `AllowTypeFalse`, `AllowTypeTrue`
- New enum type `AuthMode` with values `AuthModeOptInAllAuth`, `AuthModeOptOutAllAuth`
- New enum type `DaprBindingComponentDirection` with values `DaprBindingComponentDirectionInput`, `DaprBindingComponentDirectionOutput`
- New enum type `DaprMetadataRequired` with values `DaprMetadataRequiredFalse`, `DaprMetadataRequiredTrue`
- New enum type `DeleteOrUpdateBehavior` with values `DeleteOrUpdateBehaviorDefault`, `DeleteOrUpdateBehaviorForcedCleanup`
- New enum type `DryrunActionName` with values `DryrunActionNameCreateOrUpdate`
- New enum type `DryrunPrerequisiteResultType` with values `DryrunPrerequisiteResultTypeBasicError`, `DryrunPrerequisiteResultTypePermissionsMissing`
- New enum type `DryrunPreviewOperationType` with values `DryrunPreviewOperationTypeConfigAuth`, `DryrunPreviewOperationTypeConfigConnection`, `DryrunPreviewOperationTypeConfigNetwork`
- New enum type `LinkerConfigurationType` with values `LinkerConfigurationTypeDefault`, `LinkerConfigurationTypeKeyVaultSecret`
- New enum type `SecretSourceType` with values `SecretSourceTypeKeyVaultSecret`, `SecretSourceTypeRawValue`
- New function `*AccessKeyInfoBase.GetAuthInfoBase() *AuthInfoBase`
- New function `*BasicErrorDryrunPrerequisiteResult.GetDryrunPrerequisiteResult() *DryrunPrerequisiteResult`
- New function `*ClientFactory.NewConfigurationNamesClient() *ConfigurationNamesClient`
- New function `*ClientFactory.NewConnectorClient() *ConnectorClient`
- New function `*ClientFactory.NewLinkersClient() *LinkersClient`
- New function `NewConfigurationNamesClient(azcore.TokenCredential, *arm.ClientOptions) (*ConfigurationNamesClient, error)`
- New function `*ConfigurationNamesClient.NewListPager(*ConfigurationNamesClientListOptions) *runtime.Pager[ConfigurationNamesClientListResponse]`
- New function `NewConnectorClient(azcore.TokenCredential, *arm.ClientOptions) (*ConnectorClient, error)`
- New function `*ConnectorClient.BeginCreateDryrun(context.Context, string, string, string, string, DryrunResource, *ConnectorClientBeginCreateDryrunOptions) (*runtime.Poller[ConnectorClientCreateDryrunResponse], error)`
- New function `*ConnectorClient.BeginCreateOrUpdate(context.Context, string, string, string, string, LinkerResource, *ConnectorClientBeginCreateOrUpdateOptions) (*runtime.Poller[ConnectorClientCreateOrUpdateResponse], error)`
- New function `*ConnectorClient.BeginDelete(context.Context, string, string, string, string, *ConnectorClientBeginDeleteOptions) (*runtime.Poller[ConnectorClientDeleteResponse], error)`
- New function `*ConnectorClient.DeleteDryrun(context.Context, string, string, string, string, *ConnectorClientDeleteDryrunOptions) (ConnectorClientDeleteDryrunResponse, error)`
- New function `*ConnectorClient.GenerateConfigurations(context.Context, string, string, string, string, *ConnectorClientGenerateConfigurationsOptions) (ConnectorClientGenerateConfigurationsResponse, error)`
- New function `*ConnectorClient.Get(context.Context, string, string, string, string, *ConnectorClientGetOptions) (ConnectorClientGetResponse, error)`
- New function `*ConnectorClient.GetDryrun(context.Context, string, string, string, string, *ConnectorClientGetDryrunOptions) (ConnectorClientGetDryrunResponse, error)`
- New function `*ConnectorClient.NewListDryrunPager(string, string, string, *ConnectorClientListDryrunOptions) *runtime.Pager[ConnectorClientListDryrunResponse]`
- New function `*ConnectorClient.NewListPager(string, string, string, *ConnectorClientListOptions) *runtime.Pager[ConnectorClientListResponse]`
- New function `*ConnectorClient.BeginUpdate(context.Context, string, string, string, string, LinkerPatch, *ConnectorClientBeginUpdateOptions) (*runtime.Poller[ConnectorClientUpdateResponse], error)`
- New function `*ConnectorClient.BeginUpdateDryrun(context.Context, string, string, string, string, DryrunPatch, *ConnectorClientBeginUpdateDryrunOptions) (*runtime.Poller[ConnectorClientUpdateDryrunResponse], error)`
- New function `*ConnectorClient.BeginValidate(context.Context, string, string, string, string, *ConnectorClientBeginValidateOptions) (*runtime.Poller[ConnectorClientValidateResponse], error)`
- New function `*CreateOrUpdateDryrunParameters.GetDryrunParameters() *DryrunParameters`
- New function `*DryrunParameters.GetDryrunParameters() *DryrunParameters`
- New function `*DryrunPrerequisiteResult.GetDryrunPrerequisiteResult() *DryrunPrerequisiteResult`
- New function `*EasyAuthMicrosoftEntraIDAuthInfo.GetAuthInfoBase() *AuthInfoBase`
- New function `NewLinkersClient(azcore.TokenCredential, *arm.ClientOptions) (*LinkersClient, error)`
- New function `*LinkersClient.BeginCreateDryrun(context.Context, string, string, DryrunResource, *LinkersClientBeginCreateDryrunOptions) (*runtime.Poller[LinkersClientCreateDryrunResponse], error)`
- New function `*LinkersClient.DeleteDryrun(context.Context, string, string, *LinkersClientDeleteDryrunOptions) (LinkersClientDeleteDryrunResponse, error)`
- New function `*LinkersClient.GenerateConfigurations(context.Context, string, string, *LinkersClientGenerateConfigurationsOptions) (LinkersClientGenerateConfigurationsResponse, error)`
- New function `*LinkersClient.GetDryrun(context.Context, string, string, *LinkersClientGetDryrunOptions) (LinkersClientGetDryrunResponse, error)`
- New function `*LinkersClient.NewListDaprConfigurationsPager(string, *LinkersClientListDaprConfigurationsOptions) *runtime.Pager[LinkersClientListDaprConfigurationsResponse]`
- New function `*LinkersClient.NewListDryrunPager(string, *LinkersClientListDryrunOptions) *runtime.Pager[LinkersClientListDryrunResponse]`
- New function `*LinkersClient.BeginUpdateDryrun(context.Context, string, string, DryrunPatch, *LinkersClientBeginUpdateDryrunOptions) (*runtime.Poller[LinkersClientUpdateDryrunResponse], error)`
- New function `*PermissionsMissingDryrunPrerequisiteResult.GetDryrunPrerequisiteResult() *DryrunPrerequisiteResult`
- New function `*SelfHostedServer.GetTargetServiceBase() *TargetServiceBase`
- New function `*UserAccountAuthInfo.GetAuthInfoBase() *AuthInfoBase`
- New struct `AccessKeyInfoBase`
- New struct `BasicErrorDryrunPrerequisiteResult`
- New struct `ConfigurationInfo`
- New struct `ConfigurationName`
- New struct `ConfigurationNameItem`
- New struct `ConfigurationNameResult`
- New struct `ConfigurationNames`
- New struct `ConfigurationResult`
- New struct `ConfigurationStore`
- New struct `CreateOrUpdateDryrunParameters`
- New struct `DaprConfigurationList`
- New struct `DaprConfigurationProperties`
- New struct `DaprConfigurationResource`
- New struct `DaprMetadata`
- New struct `DaprProperties`
- New struct `DatabaseAADAuthInfo`
- New struct `DryrunList`
- New struct `DryrunOperationPreview`
- New struct `DryrunPatch`
- New struct `DryrunProperties`
- New struct `DryrunResource`
- New struct `EasyAuthMicrosoftEntraIDAuthInfo`
- New struct `FirewallRules`
- New struct `PermissionsMissingDryrunPrerequisiteResult`
- New struct `PublicNetworkSolution`
- New struct `ResourceList`
- New struct `SelfHostedServer`
- New struct `UserAccountAuthInfo`
- New anonymous field `ConfigurationResult` in struct `LinkerClientListConfigurationsResponse`
- New anonymous field `ResourceList` in struct `LinkerClientListResponse`
- New field `ConfigurationInfo`, `PublicNetworkSolution` in struct `LinkerProperties`
- New field `SystemData` in struct `ProxyResource`
- New field `SystemData` in struct `Resource`
- New field `AuthMode` in struct `SecretAuthInfo`
- New field `KeyVaultSecretName` in struct `SecretStore`
- New field `AuthMode`, `DeleteOrUpdateBehavior`, `Roles` in struct `ServicePrincipalCertificateAuthInfo`
- New field `AuthMode`, `DeleteOrUpdateBehavior`, `Roles`, `UserName` in struct `ServicePrincipalSecretAuthInfo`
- New field `ConfigType`, `Description`, `KeyVaultReferenceIdentity` in struct `SourceConfiguration`
- New field `AuthMode`, `DeleteOrUpdateBehavior`, `Roles`, `UserName` in struct `SystemAssignedIdentityAuthInfo`
- New field `AuthMode`, `DeleteOrUpdateBehavior`, `Roles`, `UserName` in struct `UserAssignedIdentityAuthInfo`
- New field `DeleteOrUpdateBehavior` in struct `VNetSolution`


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

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicelinker/armservicelinker` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).