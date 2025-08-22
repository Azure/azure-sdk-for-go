# Release History

## 2.0.0-beta.1 (2025-08-21)
### Breaking Changes

- Type alias `EmailChannelAuthMethod` type has been changed from `float32` to `int32`
- Struct `ConnectionItemName` has been removed
- Struct `Error` has been removed
- Struct `ErrorBody` has been removed
- Struct `PrivateLinkResourceBase` has been removed
- Struct `Resource` has been removed

### Features Added

- New value `PublicNetworkAccessSecuredByPerimeter` added to enum type `PublicNetworkAccess`
- New enum type `AccessMode` with values `AccessModeAudit`, `AccessModeEnforced`, `AccessModeLearning`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `NspAccessRuleDirection` with values `NspAccessRuleDirectionInbound`, `NspAccessRuleDirectionOutbound`
- New enum type `ProvisioningState` with values `ProvisioningStateAccepted`, `ProvisioningStateCreating`, `ProvisioningStateDeleting`, `ProvisioningStateFailed`, `ProvisioningStateSucceeded`, `ProvisioningStateUpdating`
- New enum type `Severity` with values `SeverityError`, `SeverityWarning`
- New function `*ClientFactory.NewNetworkSecurityPerimeterConfigurationsClient() *NetworkSecurityPerimeterConfigurationsClient`
- New function `NewNetworkSecurityPerimeterConfigurationsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NetworkSecurityPerimeterConfigurationsClient, error)`
- New function `*NetworkSecurityPerimeterConfigurationsClient.Get(context.Context, string, string, string, *NetworkSecurityPerimeterConfigurationsClientGetOptions) (NetworkSecurityPerimeterConfigurationsClientGetResponse, error)`
- New function `*NetworkSecurityPerimeterConfigurationsClient.NewListPager(string, string, *NetworkSecurityPerimeterConfigurationsClientListOptions) *runtime.Pager[NetworkSecurityPerimeterConfigurationsClientListResponse]`
- New function `*NetworkSecurityPerimeterConfigurationsClient.BeginReconcile(context.Context, string, string, string, *NetworkSecurityPerimeterConfigurationsClientBeginReconcileOptions) (*runtime.Poller[NetworkSecurityPerimeterConfigurationsClientReconcileResponse], error)`
- New struct `NetworkSecurityPerimeter`
- New struct `NetworkSecurityPerimeterConfiguration`
- New struct `NetworkSecurityPerimeterConfigurationList`
- New struct `NetworkSecurityPerimeterConfigurationProperties`
- New struct `NspAccessRule`
- New struct `NspAccessRuleProperties`
- New struct `NspAccessRulePropertiesSubscriptionsItem`
- New struct `Profile`
- New struct `ProvisioningIssue`
- New struct `ProvisioningIssueProperties`
- New struct `ResourceAssociation`
- New struct `SystemData`
- New field `SystemData` in struct `Bot`
- New field `SystemData` in struct `BotChannel`
- New field `NetworkSecurityPerimeterConfigurations` in struct `BotProperties`
- New field `SystemData` in struct `ConnectionSetting`
- New field `ID`, `Name` in struct `ConnectionSettingProperties`
- New field `SystemData` in struct `ListChannelWithKeysResponse`
- New field `SystemData` in struct `PrivateEndpointConnection`
- New field `NextLink` in struct `PrivateEndpointConnectionListResult`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.

## 1.1.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2023-01-13)
### Breaking Changes

- Field `ID` of struct `ConnectionSettingProperties` has been removed
- Field `Name` of struct `ConnectionSettingProperties` has been removed
- Field `CallingWebHook` of struct `MsTeamsChannelProperties` has been removed

### Features Added

- New value `ChannelNameAcsChatChannel`, `ChannelNameM365Extensions`, `ChannelNameOmnichannel`, `ChannelNameSearchAssistant`, `ChannelNameTelephonyChannel` added to type alias `ChannelName`
- New type alias `EmailChannelAuthMethod` with values `EmailChannelAuthMethodGraph`, `EmailChannelAuthMethodPassword`
- New function `*AcsChatChannel.GetChannel() *Channel`
- New function `NewEmailClient(string, azcore.TokenCredential, *arm.ClientOptions) (*EmailClient, error)`
- New function `*EmailClient.CreateSignInURL(context.Context, string, string, *EmailClientCreateSignInURLOptions) (EmailClientCreateSignInURLResponse, error)`
- New function `*M365Extensions.GetChannel() *Channel`
- New function `*Omnichannel.GetChannel() *Channel`
- New function `*OutlookChannel.GetChannel() *Channel`
- New function `NewQnAMakerEndpointKeysClient(string, azcore.TokenCredential, *arm.ClientOptions) (*QnAMakerEndpointKeysClient, error)`
- New function `*QnAMakerEndpointKeysClient.Get(context.Context, QnAMakerEndpointKeysRequestBody, *QnAMakerEndpointKeysClientGetOptions) (QnAMakerEndpointKeysClientGetResponse, error)`
- New function `*SearchAssistant.GetChannel() *Channel`
- New function `*TelephonyChannel.GetChannel() *Channel`
- New struct `AcsChatChannel`
- New struct `CreateEmailSignInURLResponse`
- New struct `CreateEmailSignInURLResponseProperties`
- New struct `EmailClient`
- New struct `M365Extensions`
- New struct `Omnichannel`
- New struct `OutlookChannel`
- New struct `QnAMakerEndpointKeysClient`
- New struct `QnAMakerEndpointKeysRequestBody`
- New struct `QnAMakerEndpointKeysResponse`
- New struct `SearchAssistant`
- New struct `TelephonyChannel`
- New struct `TelephonyChannelProperties`
- New struct `TelephonyChannelResourceAPIConfiguration`
- New struct `TelephonyPhoneNumbers`
- New field `TenantID` in struct `BotProperties`
- New field `RequireTermsAgreement` in struct `ChannelSettings`
- New field `AbsCode` in struct `CheckNameAvailabilityResponseBody`
- New field `ExtensionKey1` in struct `DirectLineChannelProperties`
- New field `ExtensionKey2` in struct `DirectLineChannelProperties`
- New field `AppID` in struct `DirectLineSite`
- New field `ETag` in struct `DirectLineSite`
- New field `IsDetailedLoggingEnabled` in struct `DirectLineSite`
- New field `IsEndpointParametersEnabled` in struct `DirectLineSite`
- New field `IsNoStorageEnabled` in struct `DirectLineSite`
- New field `IsTokenEnabled` in struct `DirectLineSite`
- New field `IsWebChatSpeechEnabled` in struct `DirectLineSite`
- New field `IsWebchatPreviewEnabled` in struct `DirectLineSite`
- New field `TenantID` in struct `DirectLineSite`
- New field `CognitiveServiceResourceID` in struct `DirectLineSpeechChannelProperties`
- New field `AuthMethod` in struct `EmailChannelProperties`
- New field `MagicCode` in struct `EmailChannelProperties`
- New field `CallingWebhook` in struct `MsTeamsChannelProperties`
- New field `GroupIDs` in struct `PrivateEndpointConnectionProperties`
- New field `AppID` in struct `Site`
- New field `IsDetailedLoggingEnabled` in struct `Site`
- New field `IsEndpointParametersEnabled` in struct `Site`
- New field `IsNoStorageEnabled` in struct `Site`
- New field `IsWebChatSpeechEnabled` in struct `Site`
- New field `TenantID` in struct `Site`
- New field `AppID` in struct `WebChatSite`
- New field `ETag` in struct `WebChatSite`
- New field `IsBlockUserUploadEnabled` in struct `WebChatSite`
- New field `IsDetailedLoggingEnabled` in struct `WebChatSite`
- New field `IsEndpointParametersEnabled` in struct `WebChatSite`
- New field `IsNoStorageEnabled` in struct `WebChatSite`
- New field `IsSecureSiteEnabled` in struct `WebChatSite`
- New field `IsTokenEnabled` in struct `WebChatSite`
- New field `IsV1Enabled` in struct `WebChatSite`
- New field `IsV3Enabled` in struct `WebChatSite`
- New field `IsWebChatSpeechEnabled` in struct `WebChatSite`
- New field `TenantID` in struct `WebChatSite`
- New field `TrustedOrigins` in struct `WebChatSite`


## 0.5.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/botservice/armbotservice` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.5.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).