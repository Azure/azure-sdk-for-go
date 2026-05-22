# Release History

## 2.0.0-beta.1 (2026-02-10)
### Breaking Changes

- Function `*MarketplaceAgreementsClient.Create` parameter(s) have been changed from `(ctx context.Context, options *MarketplaceAgreementsClientCreateOptions)` to `(ctx context.Context, body AgreementResource, options *MarketplaceAgreementsClientCreateOptions)`
- Function `*OrganizationClient.BeginCreate` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, organizationName string, options *OrganizationClientBeginCreateOptions)` to `(ctx context.Context, resourceGroupName string, organizationName string, body OrganizationResource, options *OrganizationClientBeginCreateOptions)`
- Function `*OrganizationClient.Update` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, organizationName string, options *OrganizationClientUpdateOptions)` to `(ctx context.Context, resourceGroupName string, organizationName string, body OrganizationResourceUpdate, options *OrganizationClientUpdateOptions)`
- Struct `ErrorResponseBody` has been removed
- Struct `ResourceProviderDefaultErrorResponse` has been removed
- Struct `SCConfluentListMetadata` has been removed
- Field `Body` of struct `MarketplaceAgreementsClientCreateOptions` has been removed
- Field `Body` of struct `OrganizationClientBeginCreateOptions` has been removed
- Field `Body` of struct `OrganizationClientUpdateOptions` has been removed

### Features Added

- New enum type `AuthType` with values `AuthTypeKAFKAAPIKEY`, `AuthTypeSERVICEACCOUNT`
- New enum type `ConnectorClass` with values `ConnectorClassAZUREBLOBSINK`, `ConnectorClassAZUREBLOBSOURCE`, `ConnectorClassAZURECOSMOSV2SINK`, `ConnectorClassAZURECOSMOSV2SOURCE`
- New enum type `ConnectorServiceType` with values `ConnectorServiceTypeAzureBlobStorageSinkConnector`, `ConnectorServiceTypeAzureBlobStorageSourceConnector`, `ConnectorServiceTypeAzureCosmosDBSinkConnector`, `ConnectorServiceTypeAzureCosmosDBSourceConnector`, `ConnectorServiceTypeAzureSynapseAnalyticsSinkConnector`
- New enum type `ConnectorStatus` with values `ConnectorStatusFAILED`, `ConnectorStatusPAUSED`, `ConnectorStatusPROVISIONING`, `ConnectorStatusRUNNING`
- New enum type `ConnectorType` with values `ConnectorTypeSINK`, `ConnectorTypeSOURCE`
- New enum type `DataFormatType` with values `DataFormatTypeAVRO`, `DataFormatTypeBYTES`, `DataFormatTypeJSON`, `DataFormatTypePROTOBUF`, `DataFormatTypeSTRING`
- New enum type `Package` with values `PackageADVANCED`, `PackageESSENTIALS`
- New enum type `PartnerConnectorType` with values `PartnerConnectorTypeKafkaAzureBlobStorageSink`, `PartnerConnectorTypeKafkaAzureBlobStorageSource`, `PartnerConnectorTypeKafkaAzureCosmosDBSink`, `PartnerConnectorTypeKafkaAzureCosmosDBSource`, `PartnerConnectorTypeKafkaAzureSynapseAnalyticsSink`
- New function `*AzureBlobStorageSinkConnectorServiceInfo.GetConnectorServiceTypeInfoBase() *ConnectorServiceTypeInfoBase`
- New function `*AzureBlobStorageSourceConnectorServiceInfo.GetConnectorServiceTypeInfoBase() *ConnectorServiceTypeInfoBase`
- New function `*AzureCosmosDBSinkConnectorServiceInfo.GetConnectorServiceTypeInfoBase() *ConnectorServiceTypeInfoBase`
- New function `*AzureCosmosDBSourceConnectorServiceInfo.GetConnectorServiceTypeInfoBase() *ConnectorServiceTypeInfoBase`
- New function `*AzureSynapseAnalyticsSinkConnectorServiceInfo.GetConnectorServiceTypeInfoBase() *ConnectorServiceTypeInfoBase`
- New function `*ClientFactory.NewClusterClient() *ClusterClient`
- New function `*ClientFactory.NewConnectorClient() *ConnectorClient`
- New function `*ClientFactory.NewEnvironmentClient() *EnvironmentClient`
- New function `*ClientFactory.NewTopicsClient() *TopicsClient`
- New function `NewClusterClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ClusterClient, error)`
- New function `*ClusterClient.CreateOrUpdate(ctx context.Context, resourceGroupName string, organizationName string, environmentID string, clusterID string, body SCClusterRecord, options *ClusterClientCreateOrUpdateOptions) (ClusterClientCreateOrUpdateResponse, error)`
- New function `*ClusterClient.BeginDelete(ctx context.Context, resourceGroupName string, organizationName string, environmentID string, clusterID string, options *ClusterClientBeginDeleteOptions) (*runtime.Poller[ClusterClientDeleteResponse], error)`
- New function `NewConnectorClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ConnectorClient, error)`
- New function `*ConnectorClient.CreateOrUpdate(ctx context.Context, resourceGroupName string, organizationName string, environmentID string, clusterID string, connectorName string, body ConnectorResource, options *ConnectorClientCreateOrUpdateOptions) (ConnectorClientCreateOrUpdateResponse, error)`
- New function `*ConnectorClient.BeginDelete(ctx context.Context, resourceGroupName string, organizationName string, environmentID string, clusterID string, connectorName string, options *ConnectorClientBeginDeleteOptions) (*runtime.Poller[ConnectorClientDeleteResponse], error)`
- New function `*ConnectorClient.Get(ctx context.Context, resourceGroupName string, organizationName string, environmentID string, clusterID string, connectorName string, options *ConnectorClientGetOptions) (ConnectorClientGetResponse, error)`
- New function `*ConnectorClient.NewListPager(resourceGroupName string, organizationName string, environmentID string, clusterID string, options *ConnectorClientListOptions) *runtime.Pager[ConnectorClientListResponse]`
- New function `*ConnectorServiceTypeInfoBase.GetConnectorServiceTypeInfoBase() *ConnectorServiceTypeInfoBase`
- New function `NewEnvironmentClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*EnvironmentClient, error)`
- New function `*EnvironmentClient.CreateOrUpdate(ctx context.Context, resourceGroupName string, organizationName string, environmentID string, body SCEnvironmentRecord, options *EnvironmentClientCreateOrUpdateOptions) (EnvironmentClientCreateOrUpdateResponse, error)`
- New function `*EnvironmentClient.BeginDelete(ctx context.Context, resourceGroupName string, organizationName string, environmentID string, options *EnvironmentClientBeginDeleteOptions) (*runtime.Poller[EnvironmentClientDeleteResponse], error)`
- New function `*KafkaAzureBlobStorageSinkConnectorInfo.GetPartnerInfoBase() *PartnerInfoBase`
- New function `*KafkaAzureBlobStorageSourceConnectorInfo.GetPartnerInfoBase() *PartnerInfoBase`
- New function `*KafkaAzureCosmosDBSinkConnectorInfo.GetPartnerInfoBase() *PartnerInfoBase`
- New function `*KafkaAzureCosmosDBSourceConnectorInfo.GetPartnerInfoBase() *PartnerInfoBase`
- New function `*KafkaAzureSynapseAnalyticsSinkConnectorInfo.GetPartnerInfoBase() *PartnerInfoBase`
- New function `*PartnerInfoBase.GetPartnerInfoBase() *PartnerInfoBase`
- New function `NewTopicsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*TopicsClient, error)`
- New function `*TopicsClient.Create(ctx context.Context, resourceGroupName string, organizationName string, environmentID string, clusterID string, topicName string, body TopicRecord, options *TopicsClientCreateOptions) (TopicsClientCreateResponse, error)`
- New function `*TopicsClient.BeginDelete(ctx context.Context, resourceGroupName string, organizationName string, environmentID string, clusterID string, topicName string, options *TopicsClientBeginDeleteOptions) (*runtime.Poller[TopicsClientDeleteResponse], error)`
- New function `*TopicsClient.Get(ctx context.Context, resourceGroupName string, organizationName string, environmentID string, clusterID string, topicName string, options *TopicsClientGetOptions) (TopicsClientGetResponse, error)`
- New function `*TopicsClient.NewListPager(resourceGroupName string, organizationName string, environmentID string, clusterID string, options *TopicsClientListOptions) *runtime.Pager[TopicsClientListResponse]`
- New struct `AzureBlobStorageSinkConnectorServiceInfo`
- New struct `AzureBlobStorageSourceConnectorServiceInfo`
- New struct `AzureCosmosDBSinkConnectorServiceInfo`
- New struct `AzureCosmosDBSourceConnectorServiceInfo`
- New struct `AzureSynapseAnalyticsSinkConnectorServiceInfo`
- New struct `ConnectorInfoBase`
- New struct `ConnectorResource`
- New struct `ConnectorResourceProperties`
- New struct `KafkaAzureBlobStorageSinkConnectorInfo`
- New struct `KafkaAzureBlobStorageSourceConnectorInfo`
- New struct `KafkaAzureCosmosDBSinkConnectorInfo`
- New struct `KafkaAzureCosmosDBSourceConnectorInfo`
- New struct `KafkaAzureSynapseAnalyticsSinkConnectorInfo`
- New struct `ListConnectorsSuccessResponse`
- New struct `ListTopicsSuccessResponse`
- New struct `StreamGovernanceConfig`
- New struct `TopicMetadataEntity`
- New struct `TopicProperties`
- New struct `TopicRecord`
- New struct `TopicsInputConfig`
- New struct `TopicsRelatedLink`
- New field `StreamGovernanceConfig` in struct `EnvironmentProperties`
- New field `SystemData`, `Type` in struct `SCClusterRecord`
- New field `Package` in struct `SCClusterSpecEntity`
- New field `SystemData`, `Type` in struct `SCEnvironmentRecord`


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
