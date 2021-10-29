# Release History

## 0.2.0 (2021-10-29)
### Breaking Changes

- Function `NewPrivateEndpointConnectionsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewClustersClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewEventHubsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewOperationsClient` parameter(s) have been changed from `(*arm.Connection)` to `(azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewDisasterRecoveryConfigsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewNamespacesClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewPrivateLinkResourcesClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewConfigurationClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewConsumerGroupsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NetworkRuleSet.MarshalJSON` has been removed
- Function `AuthorizationRule.MarshalJSON` has been removed
- Function `ConsumerGroup.MarshalJSON` has been removed
- Function `Eventhub.MarshalJSON` has been removed
- Function `ArmDisasterRecovery.MarshalJSON` has been removed
- Function `PrivateEndpointConnection.MarshalJSON` has been removed
- Field `Resource` of struct `AuthorizationRule` has been removed
- Field `Resource` of struct `ArmDisasterRecovery` has been removed
- Field `Resource` of struct `ConsumerGroup` has been removed
- Field `Code` of struct `ErrorResponse` has been removed
- Field `Message` of struct `ErrorResponse` has been removed
- Field `Resource` of struct `Eventhub` has been removed
- Field `Resource` of struct `PrivateEndpointConnection` has been removed
- Field `Resource` of struct `NetworkRuleSet` has been removed

### New Content

- New const `SchemaCompatibilityForward`
- New const `SchemaCompatibilityBackward`
- New const `SchemaTypeAvro`
- New const `SchemaTypeUnknown`
- New const `SchemaCompatibilityNone`
- New function `SchemaType.ToPtr() *SchemaType`
- New function `ErrorDetail.MarshalJSON() ([]byte, error)`
- New function `*SchemaGroupProperties.UnmarshalJSON([]byte) error`
- New function `*SchemaRegistryClient.Get(context.Context, string, string, string, *SchemaRegistryGetOptions) (SchemaRegistryGetResponse, error)`
- New function `*SchemaRegistryClient.Delete(context.Context, string, string, string, *SchemaRegistryDeleteOptions) (SchemaRegistryDeleteResponse, error)`
- New function `SchemaGroupProperties.MarshalJSON() ([]byte, error)`
- New function `*SchemaRegistryListByNamespacePager.NextPage(context.Context) bool`
- New function `*NamespacesClient.ListNetworkRuleSet(context.Context, string, string, *NamespacesListNetworkRuleSetOptions) (NamespacesListNetworkRuleSetResponse, error)`
- New function `NetworkRuleSetListResult.MarshalJSON() ([]byte, error)`
- New function `PossibleSchemaTypeValues() []SchemaType`
- New function `*SchemaRegistryClient.ListByNamespace(string, string, *SchemaRegistryListByNamespaceOptions) *SchemaRegistryListByNamespacePager`
- New function `*SchemaRegistryListByNamespacePager.Err() error`
- New function `NewSchemaRegistryClient(string, azcore.TokenCredential, *arm.ClientOptions) *SchemaRegistryClient`
- New function `*SchemaRegistryClient.CreateOrUpdate(context.Context, string, string, string, SchemaGroup, *SchemaRegistryCreateOrUpdateOptions) (SchemaRegistryCreateOrUpdateResponse, error)`
- New function `PossibleSchemaCompatibilityValues() []SchemaCompatibility`
- New function `SchemaGroupListResult.MarshalJSON() ([]byte, error)`
- New function `SchemaCompatibility.ToPtr() *SchemaCompatibility`
- New function `*SchemaRegistryListByNamespacePager.PageResponse() SchemaRegistryListByNamespaceResponse`
- New struct `ErrorAdditionalInfo`
- New struct `ErrorDetail`
- New struct `NamespacesListNetworkRuleSetOptions`
- New struct `NamespacesListNetworkRuleSetResponse`
- New struct `NamespacesListNetworkRuleSetResult`
- New struct `NetworkRuleSetListResult`
- New struct `ProxyResource`
- New struct `SchemaGroup`
- New struct `SchemaGroupListResult`
- New struct `SchemaGroupProperties`
- New struct `SchemaRegistryClient`
- New struct `SchemaRegistryCreateOrUpdateOptions`
- New struct `SchemaRegistryCreateOrUpdateResponse`
- New struct `SchemaRegistryCreateOrUpdateResult`
- New struct `SchemaRegistryDeleteOptions`
- New struct `SchemaRegistryDeleteResponse`
- New struct `SchemaRegistryGetOptions`
- New struct `SchemaRegistryGetResponse`
- New struct `SchemaRegistryGetResult`
- New struct `SchemaRegistryListByNamespaceOptions`
- New struct `SchemaRegistryListByNamespacePager`
- New struct `SchemaRegistryListByNamespaceResponse`
- New struct `SchemaRegistryListByNamespaceResult`
- New field `DataLakeAccountName` in struct `DestinationProperties`
- New field `DataLakeFolderPath` in struct `DestinationProperties`
- New field `DataLakeSubscriptionID` in struct `DestinationProperties`
- New anonymous field `ProxyResource` in struct `PrivateEndpointConnection`
- New field `InnerError` in struct `ErrorResponse`
- New field `Description` in struct `OperationDisplay`
- New anonymous field `ProxyResource` in struct `Eventhub`
- New field `AlternateName` in struct `EHNamespaceProperties`
- New anonymous field `ProxyResource` in struct `ConsumerGroup`
- New anonymous field `ProxyResource` in struct `ArmDisasterRecovery`
- New anonymous field `ProxyResource` in struct `AuthorizationRule`
- New field `Properties` in struct `Operation`
- New field `IsDataAction` in struct `Operation`
- New field `Origin` in struct `Operation`
- New anonymous field `ProxyResource` in struct `NetworkRuleSet`

Total 22 breaking change(s), 80 additive change(s).


## 0.1.0 (2021-10-08)
- To better align with the Azure SDK guidelines (https://azure.github.io/azure-sdk/general_introduction.html), we have decided to change the module path to "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/eventhub/armeventhub". Therefore, we are deprecating the old module path (which is "github.com/Azure/azure-sdk-for-go/sdk/eventhub/armeventhub") to avoid confusion.