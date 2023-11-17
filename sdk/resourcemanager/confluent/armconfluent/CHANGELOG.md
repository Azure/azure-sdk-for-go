# Release History

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
