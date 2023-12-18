# Release History

## 2.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.0.0 (2023-08-25)
### Breaking Changes

- Function `*ContentItemClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, string, string, *ContentItemClientCreateOrUpdateOptions)` to `(context.Context, string, string, string, string, ContentItemContract, *ContentItemClientCreateOrUpdateOptions)`
- Function `*ContentTypeClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, string, *ContentTypeClientCreateOrUpdateOptions)` to `(context.Context, string, string, string, ContentTypeContract, *ContentTypeClientCreateOrUpdateOptions)`

### Features Added

- New enum type `AsyncResolverStatus` with values `AsyncResolverStatusFailed`, `AsyncResolverStatusInProgress`, `AsyncResolverStatusStarted`, `AsyncResolverStatusSucceeded`
- New enum type `AuthorizationType` with values `AuthorizationTypeOAuth2`
- New enum type `NatGatewayState` with values `NatGatewayStateDisabled`, `NatGatewayStateEnabled`
- New enum type `OAuth2GrantType` with values `OAuth2GrantTypeAuthorizationCode`, `OAuth2GrantTypeClientCredentials`
- New enum type `PolicyFragmentContentFormat` with values `PolicyFragmentContentFormatRawxml`, `PolicyFragmentContentFormatXML`
- New enum type `PortalSettingsCspMode` with values `PortalSettingsCspModeDisabled`, `PortalSettingsCspModeEnabled`, `PortalSettingsCspModeReportOnly`
- New enum type `TranslateRequiredQueryParametersConduct` with values `TranslateRequiredQueryParametersConductQuery`, `TranslateRequiredQueryParametersConductTemplate`
- New function `NewAPIWikiClient(string, azcore.TokenCredential, *arm.ClientOptions) (*APIWikiClient, error)`
- New function `*APIWikiClient.CreateOrUpdate(context.Context, string, string, string, WikiContract, *APIWikiClientCreateOrUpdateOptions) (APIWikiClientCreateOrUpdateResponse, error)`
- New function `*APIWikiClient.Delete(context.Context, string, string, string, string, *APIWikiClientDeleteOptions) (APIWikiClientDeleteResponse, error)`
- New function `*APIWikiClient.Get(context.Context, string, string, string, *APIWikiClientGetOptions) (APIWikiClientGetResponse, error)`
- New function `*APIWikiClient.GetEntityTag(context.Context, string, string, string, *APIWikiClientGetEntityTagOptions) (APIWikiClientGetEntityTagResponse, error)`
- New function `*APIWikiClient.Update(context.Context, string, string, string, string, WikiUpdateContract, *APIWikiClientUpdateOptions) (APIWikiClientUpdateResponse, error)`
- New function `NewAPIWikisClient(string, azcore.TokenCredential, *arm.ClientOptions) (*APIWikisClient, error)`
- New function `*APIWikisClient.NewListPager(string, string, string, *APIWikisClientListOptions) *runtime.Pager[APIWikisClientListResponse]`
- New function `NewAuthorizationAccessPolicyClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AuthorizationAccessPolicyClient, error)`
- New function `*AuthorizationAccessPolicyClient.CreateOrUpdate(context.Context, string, string, string, string, string, AuthorizationAccessPolicyContract, *AuthorizationAccessPolicyClientCreateOrUpdateOptions) (AuthorizationAccessPolicyClientCreateOrUpdateResponse, error)`
- New function `*AuthorizationAccessPolicyClient.Delete(context.Context, string, string, string, string, string, string, *AuthorizationAccessPolicyClientDeleteOptions) (AuthorizationAccessPolicyClientDeleteResponse, error)`
- New function `*AuthorizationAccessPolicyClient.Get(context.Context, string, string, string, string, string, *AuthorizationAccessPolicyClientGetOptions) (AuthorizationAccessPolicyClientGetResponse, error)`
- New function `*AuthorizationAccessPolicyClient.NewListByAuthorizationPager(string, string, string, string, *AuthorizationAccessPolicyClientListByAuthorizationOptions) *runtime.Pager[AuthorizationAccessPolicyClientListByAuthorizationResponse]`
- New function `NewAuthorizationClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AuthorizationClient, error)`
- New function `*AuthorizationClient.ConfirmConsentCode(context.Context, string, string, string, string, AuthorizationConfirmConsentCodeRequestContract, *AuthorizationClientConfirmConsentCodeOptions) (AuthorizationClientConfirmConsentCodeResponse, error)`
- New function `*AuthorizationClient.CreateOrUpdate(context.Context, string, string, string, string, AuthorizationContract, *AuthorizationClientCreateOrUpdateOptions) (AuthorizationClientCreateOrUpdateResponse, error)`
- New function `*AuthorizationClient.Delete(context.Context, string, string, string, string, string, *AuthorizationClientDeleteOptions) (AuthorizationClientDeleteResponse, error)`
- New function `*AuthorizationClient.Get(context.Context, string, string, string, string, *AuthorizationClientGetOptions) (AuthorizationClientGetResponse, error)`
- New function `*AuthorizationClient.NewListByAuthorizationProviderPager(string, string, string, *AuthorizationClientListByAuthorizationProviderOptions) *runtime.Pager[AuthorizationClientListByAuthorizationProviderResponse]`
- New function `NewAuthorizationLoginLinksClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AuthorizationLoginLinksClient, error)`
- New function `*AuthorizationLoginLinksClient.Post(context.Context, string, string, string, string, AuthorizationLoginRequestContract, *AuthorizationLoginLinksClientPostOptions) (AuthorizationLoginLinksClientPostResponse, error)`
- New function `NewAuthorizationProviderClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AuthorizationProviderClient, error)`
- New function `*AuthorizationProviderClient.CreateOrUpdate(context.Context, string, string, string, AuthorizationProviderContract, *AuthorizationProviderClientCreateOrUpdateOptions) (AuthorizationProviderClientCreateOrUpdateResponse, error)`
- New function `*AuthorizationProviderClient.Delete(context.Context, string, string, string, string, *AuthorizationProviderClientDeleteOptions) (AuthorizationProviderClientDeleteResponse, error)`
- New function `*AuthorizationProviderClient.Get(context.Context, string, string, string, *AuthorizationProviderClientGetOptions) (AuthorizationProviderClientGetResponse, error)`
- New function `*AuthorizationProviderClient.NewListByServicePager(string, string, *AuthorizationProviderClientListByServiceOptions) *runtime.Pager[AuthorizationProviderClientListByServiceResponse]`
- New function `*ClientFactory.NewAPIWikiClient() *APIWikiClient`
- New function `*ClientFactory.NewAPIWikisClient() *APIWikisClient`
- New function `*ClientFactory.NewAuthorizationAccessPolicyClient() *AuthorizationAccessPolicyClient`
- New function `*ClientFactory.NewAuthorizationClient() *AuthorizationClient`
- New function `*ClientFactory.NewAuthorizationLoginLinksClient() *AuthorizationLoginLinksClient`
- New function `*ClientFactory.NewAuthorizationProviderClient() *AuthorizationProviderClient`
- New function `*ClientFactory.NewDocumentationClient() *DocumentationClient`
- New function `*ClientFactory.NewGraphQLAPIResolverClient() *GraphQLAPIResolverClient`
- New function `*ClientFactory.NewGraphQLAPIResolverPolicyClient() *GraphQLAPIResolverPolicyClient`
- New function `*ClientFactory.NewPolicyFragmentClient() *PolicyFragmentClient`
- New function `*ClientFactory.NewPortalConfigClient() *PortalConfigClient`
- New function `*ClientFactory.NewProductWikiClient() *ProductWikiClient`
- New function `*ClientFactory.NewProductWikisClient() *ProductWikisClient`
- New function `NewDocumentationClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DocumentationClient, error)`
- New function `*DocumentationClient.CreateOrUpdate(context.Context, string, string, string, DocumentationContract, *DocumentationClientCreateOrUpdateOptions) (DocumentationClientCreateOrUpdateResponse, error)`
- New function `*DocumentationClient.Delete(context.Context, string, string, string, string, *DocumentationClientDeleteOptions) (DocumentationClientDeleteResponse, error)`
- New function `*DocumentationClient.Get(context.Context, string, string, string, *DocumentationClientGetOptions) (DocumentationClientGetResponse, error)`
- New function `*DocumentationClient.GetEntityTag(context.Context, string, string, string, *DocumentationClientGetEntityTagOptions) (DocumentationClientGetEntityTagResponse, error)`
- New function `*DocumentationClient.NewListByServicePager(string, string, *DocumentationClientListByServiceOptions) *runtime.Pager[DocumentationClientListByServiceResponse]`
- New function `*DocumentationClient.Update(context.Context, string, string, string, string, DocumentationUpdateContract, *DocumentationClientUpdateOptions) (DocumentationClientUpdateResponse, error)`
- New function `NewGraphQLAPIResolverClient(string, azcore.TokenCredential, *arm.ClientOptions) (*GraphQLAPIResolverClient, error)`
- New function `*GraphQLAPIResolverClient.CreateOrUpdate(context.Context, string, string, string, string, ResolverContract, *GraphQLAPIResolverClientCreateOrUpdateOptions) (GraphQLAPIResolverClientCreateOrUpdateResponse, error)`
- New function `*GraphQLAPIResolverClient.Delete(context.Context, string, string, string, string, string, *GraphQLAPIResolverClientDeleteOptions) (GraphQLAPIResolverClientDeleteResponse, error)`
- New function `*GraphQLAPIResolverClient.Get(context.Context, string, string, string, string, *GraphQLAPIResolverClientGetOptions) (GraphQLAPIResolverClientGetResponse, error)`
- New function `*GraphQLAPIResolverClient.GetEntityTag(context.Context, string, string, string, string, *GraphQLAPIResolverClientGetEntityTagOptions) (GraphQLAPIResolverClientGetEntityTagResponse, error)`
- New function `*GraphQLAPIResolverClient.NewListByAPIPager(string, string, string, *GraphQLAPIResolverClientListByAPIOptions) *runtime.Pager[GraphQLAPIResolverClientListByAPIResponse]`
- New function `*GraphQLAPIResolverClient.Update(context.Context, string, string, string, string, string, ResolverUpdateContract, *GraphQLAPIResolverClientUpdateOptions) (GraphQLAPIResolverClientUpdateResponse, error)`
- New function `NewGraphQLAPIResolverPolicyClient(string, azcore.TokenCredential, *arm.ClientOptions) (*GraphQLAPIResolverPolicyClient, error)`
- New function `*GraphQLAPIResolverPolicyClient.CreateOrUpdate(context.Context, string, string, string, string, PolicyIDName, PolicyContract, *GraphQLAPIResolverPolicyClientCreateOrUpdateOptions) (GraphQLAPIResolverPolicyClientCreateOrUpdateResponse, error)`
- New function `*GraphQLAPIResolverPolicyClient.Delete(context.Context, string, string, string, string, PolicyIDName, string, *GraphQLAPIResolverPolicyClientDeleteOptions) (GraphQLAPIResolverPolicyClientDeleteResponse, error)`
- New function `*GraphQLAPIResolverPolicyClient.Get(context.Context, string, string, string, string, PolicyIDName, *GraphQLAPIResolverPolicyClientGetOptions) (GraphQLAPIResolverPolicyClientGetResponse, error)`
- New function `*GraphQLAPIResolverPolicyClient.GetEntityTag(context.Context, string, string, string, string, PolicyIDName, *GraphQLAPIResolverPolicyClientGetEntityTagOptions) (GraphQLAPIResolverPolicyClientGetEntityTagResponse, error)`
- New function `*GraphQLAPIResolverPolicyClient.NewListByResolverPager(string, string, string, string, *GraphQLAPIResolverPolicyClientListByResolverOptions) *runtime.Pager[GraphQLAPIResolverPolicyClientListByResolverResponse]`
- New function `NewPolicyFragmentClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PolicyFragmentClient, error)`
- New function `*PolicyFragmentClient.BeginCreateOrUpdate(context.Context, string, string, string, PolicyFragmentContract, *PolicyFragmentClientBeginCreateOrUpdateOptions) (*runtime.Poller[PolicyFragmentClientCreateOrUpdateResponse], error)`
- New function `*PolicyFragmentClient.Delete(context.Context, string, string, string, string, *PolicyFragmentClientDeleteOptions) (PolicyFragmentClientDeleteResponse, error)`
- New function `*PolicyFragmentClient.Get(context.Context, string, string, string, *PolicyFragmentClientGetOptions) (PolicyFragmentClientGetResponse, error)`
- New function `*PolicyFragmentClient.GetEntityTag(context.Context, string, string, string, *PolicyFragmentClientGetEntityTagOptions) (PolicyFragmentClientGetEntityTagResponse, error)`
- New function `*PolicyFragmentClient.ListByService(context.Context, string, string, *PolicyFragmentClientListByServiceOptions) (PolicyFragmentClientListByServiceResponse, error)`
- New function `*PolicyFragmentClient.ListReferences(context.Context, string, string, string, *PolicyFragmentClientListReferencesOptions) (PolicyFragmentClientListReferencesResponse, error)`
- New function `NewPortalConfigClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PortalConfigClient, error)`
- New function `*PortalConfigClient.CreateOrUpdate(context.Context, string, string, string, string, PortalConfigContract, *PortalConfigClientCreateOrUpdateOptions) (PortalConfigClientCreateOrUpdateResponse, error)`
- New function `*PortalConfigClient.Get(context.Context, string, string, string, *PortalConfigClientGetOptions) (PortalConfigClientGetResponse, error)`
- New function `*PortalConfigClient.GetEntityTag(context.Context, string, string, string, *PortalConfigClientGetEntityTagOptions) (PortalConfigClientGetEntityTagResponse, error)`
- New function `*PortalConfigClient.ListByService(context.Context, string, string, *PortalConfigClientListByServiceOptions) (PortalConfigClientListByServiceResponse, error)`
- New function `*PortalConfigClient.Update(context.Context, string, string, string, string, PortalConfigContract, *PortalConfigClientUpdateOptions) (PortalConfigClientUpdateResponse, error)`
- New function `NewProductWikiClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ProductWikiClient, error)`
- New function `*ProductWikiClient.CreateOrUpdate(context.Context, string, string, string, WikiContract, *ProductWikiClientCreateOrUpdateOptions) (ProductWikiClientCreateOrUpdateResponse, error)`
- New function `*ProductWikiClient.Delete(context.Context, string, string, string, string, *ProductWikiClientDeleteOptions) (ProductWikiClientDeleteResponse, error)`
- New function `*ProductWikiClient.Get(context.Context, string, string, string, *ProductWikiClientGetOptions) (ProductWikiClientGetResponse, error)`
- New function `*ProductWikiClient.GetEntityTag(context.Context, string, string, string, *ProductWikiClientGetEntityTagOptions) (ProductWikiClientGetEntityTagResponse, error)`
- New function `*ProductWikiClient.Update(context.Context, string, string, string, string, WikiUpdateContract, *ProductWikiClientUpdateOptions) (ProductWikiClientUpdateResponse, error)`
- New function `NewProductWikisClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ProductWikisClient, error)`
- New function `*ProductWikisClient.NewListPager(string, string, string, *ProductWikisClientListOptions) *runtime.Pager[ProductWikisClientListResponse]`
- New function `*ServiceClient.BeginMigrateToStv2(context.Context, string, string, *ServiceClientBeginMigrateToStv2Options) (*runtime.Poller[ServiceClientMigrateToStv2Response], error)`
- New struct `AuthorizationAccessPolicyCollection`
- New struct `AuthorizationAccessPolicyContract`
- New struct `AuthorizationAccessPolicyContractProperties`
- New struct `AuthorizationCollection`
- New struct `AuthorizationConfirmConsentCodeRequestContract`
- New struct `AuthorizationContract`
- New struct `AuthorizationContractProperties`
- New struct `AuthorizationError`
- New struct `AuthorizationLoginRequestContract`
- New struct `AuthorizationLoginResponseContract`
- New struct `AuthorizationProviderCollection`
- New struct `AuthorizationProviderContract`
- New struct `AuthorizationProviderContractProperties`
- New struct `AuthorizationProviderOAuth2GrantTypes`
- New struct `AuthorizationProviderOAuth2Settings`
- New struct `DocumentationCollection`
- New struct `DocumentationContract`
- New struct `DocumentationContractProperties`
- New struct `DocumentationUpdateContract`
- New struct `PolicyFragmentCollection`
- New struct `PolicyFragmentContract`
- New struct `PolicyFragmentContractProperties`
- New struct `PortalConfigCollection`
- New struct `PortalConfigContract`
- New struct `PortalConfigCorsProperties`
- New struct `PortalConfigCspProperties`
- New struct `PortalConfigDelegationProperties`
- New struct `PortalConfigProperties`
- New struct `PortalConfigPropertiesSignin`
- New struct `PortalConfigPropertiesSignup`
- New struct `PortalConfigTermsOfServiceProperties`
- New struct `ProxyResource`
- New struct `ResolverCollection`
- New struct `ResolverContract`
- New struct `ResolverEntityBaseContract`
- New struct `ResolverResultContract`
- New struct `ResolverResultContractProperties`
- New struct `ResolverResultLogItemContract`
- New struct `ResolverUpdateContract`
- New struct `ResolverUpdateContractProperties`
- New struct `ResourceCollection`
- New struct `ResourceCollectionValueItem`
- New struct `WikiCollection`
- New struct `WikiContract`
- New struct `WikiContractProperties`
- New struct `WikiDocumentationContract`
- New struct `WikiUpdateContract`
- New field `TranslateRequiredQueryParametersConduct` in struct `APICreateOrUpdateProperties`
- New field `NatGatewayState`, `OutboundPublicIPAddresses` in struct `AdditionalLocation`
- New field `OAuth2AuthenticationSettings`, `OpenidAuthenticationSettings` in struct `AuthenticationSettingsContract`
- New field `UseInAPIDocumentation`, `UseInTestConsole` in struct `AuthorizationServerContractProperties`
- New field `UseInAPIDocumentation`, `UseInTestConsole` in struct `AuthorizationServerUpdateContractProperties`
- New field `Metrics` in struct `DiagnosticContractProperties`
- New field `ClientLibrary` in struct `IdentityProviderBaseParameters`
- New field `ClientLibrary` in struct `IdentityProviderContractProperties`
- New field `ClientLibrary` in struct `IdentityProviderCreateContractProperties`
- New field `ClientLibrary` in struct `IdentityProviderUpdateProperties`
- New field `UseInAPIDocumentation`, `UseInTestConsole` in struct `OpenidConnectProviderContractProperties`
- New field `UseInAPIDocumentation`, `UseInTestConsole` in struct `OpenidConnectProviderUpdateContractProperties`
- New field `NatGatewayState`, `OutboundPublicIPAddresses` in struct `ServiceBaseProperties`
- New field `NatGatewayState`, `OutboundPublicIPAddresses` in struct `ServiceProperties`
- New field `NatGatewayState`, `OutboundPublicIPAddresses` in struct `ServiceUpdateProperties`


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.

## 1.1.0 (2023-03-27)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).