# Release History

## 1.3.0 (2024-03-22)
### Features Added

- New function `*AccessClient.CreateRoleBinding(context.Context, string, string, AccessCreateRoleBindingRequestModel, *AccessClientCreateRoleBindingOptions) (AccessClientCreateRoleBindingResponse, error)`
- New function `*AccessClient.DeleteRoleBinding(context.Context, string, string, string, *AccessClientDeleteRoleBindingOptions) (AccessClientDeleteRoleBindingResponse, error)`
- New function `*AccessClient.ListRoleBindingNameList(context.Context, string, string, ListAccessRequestModel, *AccessClientListRoleBindingNameListOptions) (AccessClientListRoleBindingNameListResponse, error)`
- New function `*OrganizationClient.CreateAPIKey(context.Context, string, string, string, string, CreateAPIKeyModel, *OrganizationClientCreateAPIKeyOptions) (OrganizationClientCreateAPIKeyResponse, error)`
- New function `*OrganizationClient.DeleteClusterAPIKey(context.Context, string, string, string, *OrganizationClientDeleteClusterAPIKeyOptions) (OrganizationClientDeleteClusterAPIKeyResponse, error)`
- New function `*OrganizationClient.GetClusterAPIKey(context.Context, string, string, string, *OrganizationClientGetClusterAPIKeyOptions) (OrganizationClientGetClusterAPIKeyResponse, error)`
- New function `*OrganizationClient.GetClusterByID(context.Context, string, string, string, string, *OrganizationClientGetClusterByIDOptions) (OrganizationClientGetClusterByIDResponse, error)`
- New function `*OrganizationClient.GetEnvironmentByID(context.Context, string, string, string, *OrganizationClientGetEnvironmentByIDOptions) (OrganizationClientGetEnvironmentByIDResponse, error)`
- New function `*OrganizationClient.GetSchemaRegistryClusterByID(context.Context, string, string, string, string, *OrganizationClientGetSchemaRegistryClusterByIDOptions) (OrganizationClientGetSchemaRegistryClusterByIDResponse, error)`
- New function `*OrganizationClient.NewListClustersPager(string, string, string, *OrganizationClientListClustersOptions) *runtime.Pager[OrganizationClientListClustersResponse]`
- New function `*OrganizationClient.NewListEnvironmentsPager(string, string, *OrganizationClientListEnvironmentsOptions) *runtime.Pager[OrganizationClientListEnvironmentsResponse]`
- New function `*OrganizationClient.ListRegions(context.Context, string, string, ListAccessRequestModel, *OrganizationClientListRegionsOptions) (OrganizationClientListRegionsResponse, error)`
- New function `*OrganizationClient.NewListSchemaRegistryClustersPager(string, string, string, *OrganizationClientListSchemaRegistryClustersOptions) *runtime.Pager[OrganizationClientListSchemaRegistryClustersResponse]`
- New struct `APIKeyOwnerEntity`
- New struct `APIKeyProperties`
- New struct `APIKeyRecord`
- New struct `APIKeyResourceEntity`
- New struct `APIKeySpecEntity`
- New struct `AccessCreateRoleBindingRequestModel`
- New struct `AccessRoleBindingNameListSuccessResponse`
- New struct `ClusterProperties`
- New struct `CreateAPIKeyModel`
- New struct `EnvironmentProperties`
- New struct `GetEnvironmentsResponse`
- New struct `ListClustersSuccessResponse`
- New struct `ListRegionsSuccessResponse`
- New struct `ListSchemaRegistryClustersResponse`
- New struct `RegionProperties`
- New struct `RegionRecord`
- New struct `RegionSpecEntity`
- New struct `SCClusterByokEntity`
- New struct `SCClusterNetworkEnvironmentEntity`
- New struct `SCClusterRecord`
- New struct `SCClusterSpecEntity`
- New struct `SCConfluentListMetadata`
- New struct `SCEnvironmentRecord`
- New struct `SCMetadataEntity`
- New struct `SchemaRegistryClusterEnvironmentRegionEntity`
- New struct `SchemaRegistryClusterProperties`
- New struct `SchemaRegistryClusterRecord`
- New struct `SchemaRegistryClusterSpecEntity`
- New struct `SchemaRegistryClusterStatusEntity`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.
- New function `NewAccessClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AccessClient, error)`
- New function `*AccessClient.InviteUser(context.Context, string, string, AccessInviteUserAccountModel, *AccessClientInviteUserOptions) (AccessClientInviteUserResponse, error)`
- New function `*AccessClient.ListClusters(context.Context, string, string, ListAccessRequestModel, *AccessClientListClustersOptions) (AccessClientListClustersResponse, error)`
- New function `*AccessClient.ListEnvironments(context.Context, string, string, ListAccessRequestModel, *AccessClientListEnvironmentsOptions) (AccessClientListEnvironmentsResponse, error)`
- New function `*AccessClient.ListInvitations(context.Context, string, string, ListAccessRequestModel, *AccessClientListInvitationsOptions) (AccessClientListInvitationsResponse, error)`
- New function `*AccessClient.ListRoleBindings(context.Context, string, string, ListAccessRequestModel, *AccessClientListRoleBindingsOptions) (AccessClientListRoleBindingsResponse, error)`
- New function `*AccessClient.ListServiceAccounts(context.Context, string, string, ListAccessRequestModel, *AccessClientListServiceAccountsOptions) (AccessClientListServiceAccountsResponse, error)`
- New function `*AccessClient.ListUsers(context.Context, string, string, ListAccessRequestModel, *AccessClientListUsersOptions) (AccessClientListUsersResponse, error)`
- New function `*ClientFactory.NewAccessClient() *AccessClient`
- New function `*ValidationsClient.ValidateOrganizationV2(context.Context, string, string, OrganizationResource, *ValidationsClientValidateOrganizationV2Options) (ValidationsClientValidateOrganizationV2Response, error)`
- New struct `AccessInviteUserAccountModel`
- New struct `AccessInvitedUserDetails`
- New struct `AccessListClusterSuccessResponse`
- New struct `AccessListEnvironmentsSuccessResponse`
- New struct `AccessListInvitationsSuccessResponse`
- New struct `AccessListRoleBindingsSuccessResponse`
- New struct `AccessListServiceAccountsSuccessResponse`
- New struct `AccessListUsersSuccessResponse`
- New struct `ClusterByokEntity`
- New struct `ClusterConfigEntity`
- New struct `ClusterEnvironmentEntity`
- New struct `ClusterNetworkEntity`
- New struct `ClusterRecord`
- New struct `ClusterSpecEntity`
- New struct `ClusterStatusEntity`
- New struct `EnvironmentRecord`
- New struct `InvitationRecord`
- New struct `LinkOrganization`
- New struct `ListAccessRequestModel`
- New struct `ListMetadata`
- New struct `MetadataEntity`
- New struct `RoleBindingRecord`
- New struct `ServiceAccountRecord`
- New struct `UserRecord`
- New struct `ValidationResponse`
- New field `PrivateOfferID`, `PrivateOfferIDs`, `TermID` in struct `OfferDetail`
- New field `LinkOrganization` in struct `OrganizationResourceProperties`
- New field `AADEmail`, `UserPrincipalName` in struct `UserDetail`


## 1.1.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/confluent/armconfluent` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
