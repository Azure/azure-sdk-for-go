# Release History

## 0.2.0 (2021-10-28)
### Breaking Changes

- Function `NewContentTypeClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewGroupClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAPIExportClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewProductClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewLoggerClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewPolicyClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewQuotaByCounterKeysClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewTenantAccessGitClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewRegionClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewIssueClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAPIVersionSetClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAPIManagementOperationsClient` parameter(s) have been changed from `(*arm.Connection)` to `(azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewNamedValueClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewTenantAccessClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewSignInSettingsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewNetworkStatusClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAPIManagementServiceClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAPIReleaseClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewPolicyDescriptionClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewOutboundNetworkDependenciesEndpointsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewUserSubscriptionClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewUserConfirmationPasswordClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewEmailTemplateClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewProductPolicyClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewProductAPIClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewUserGroupClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAPIManagementServiceSKUsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewQuotaByPeriodKeysClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAPIDiagnosticClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewGatewayClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewGatewayHostnameConfigurationClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewCertificateClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewTenantSettingsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAPIManagementSKUsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewOpenIDConnectProviderClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewGatewayCertificateAuthorityClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAPIOperationPolicyClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAPISchemaClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewUserClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewNotificationRecipientEmailClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAPIOperationClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewCacheClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAPIIssueClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewReportsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAPIManagementClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAPIPolicyClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewGroupUserClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAPIClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewNotificationClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewContentItemClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewProductGroupClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewPortalRevisionClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewBackendClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewPrivateEndpointConnectionClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewTenantConfigurationClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAuthorizationServerClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewOperationClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewTagClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewDiagnosticClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAPIRevisionClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewGatewayAPIClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewSignUpSettingsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewDelegationSettingsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAPITagDescriptionClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAPIIssueCommentClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewSubscriptionClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewDeletedServicesClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewNotificationRecipientUserClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewIdentityProviderClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewProductSubscriptionsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewTagResourceClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewUserIdentitiesClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAPIProductClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewAPIIssueAttachmentClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewPortalSettingsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`

### New Content


Total 75 breaking change(s), 0 additive change(s).


## 0.1.0 (2021-10-08)
- To better align with the Azure SDK guidelines (https://azure.github.io/azure-sdk/general_introduction.html), we have decided to change the module path to "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement". Therefore, we are deprecating the old module path (which is "github.com/Azure/azure-sdk-for-go/sdk/apimanagement/armapimanagement") to avoid confusion.