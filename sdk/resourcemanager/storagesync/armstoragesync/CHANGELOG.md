# Release History

## 2.0.0-beta.1 (2026-03-19)
### Breaking Changes

- Function `*ServerEndpointsClient.BeginUpdate` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, storageSyncServiceName string, syncGroupName string, serverEndpointName string, options *ServerEndpointsClientBeginUpdateOptions)` to `(ctx context.Context, resourceGroupName string, storageSyncServiceName string, syncGroupName string, serverEndpointName string, parameters ServerEndpointUpdateParameters, options *ServerEndpointsClientBeginUpdateOptions)`
- Function `*ServicesClient.BeginUpdate` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, storageSyncServiceName string, options *ServicesClientBeginUpdateOptions)` to `(ctx context.Context, resourceGroupName string, storageSyncServiceName string, parameters ServiceUpdateParameters, options *ServicesClientBeginUpdateOptions)`
- Enum `ProgressType` has been removed
- Enum `Reason` has been removed
- Struct `Error` has been removed
- Struct `OperationDisplayResource` has been removed
- Struct `ProxyResource` has been removed
- Struct `Resource` has been removed
- Struct `ResourcesMoveInfo` has been removed
- Struct `SubscriptionState` has been removed
- Struct `TrackedResource` has been removed
- Field `XMSCorrelationRequestID`, `XMSRequestID` of struct `CloudEndpointsClientGetResponse` has been removed
- Field `XMSCorrelationRequestID`, `XMSRequestID` of struct `CloudEndpointsClientListBySyncGroupResponse` has been removed
- Field `XMSCorrelationRequestID`, `XMSRequestID` of struct `CloudEndpointsClientRestoreheartbeatResponse` has been removed
- Field `XMSCorrelationRequestID`, `XMSRequestID` of struct `MicrosoftStorageSyncClientLocationOperationStatusResponse` has been removed
- Field `XMSCorrelationRequestID`, `XMSRequestID` of struct `OperationStatusClientGetResponse` has been removed
- Field `XMSCorrelationRequestID`, `XMSRequestID` of struct `OperationsClientListResponse` has been removed
- Field `XMSCorrelationRequestID`, `XMSRequestID` of struct `PrivateEndpointConnectionsClientListByStorageSyncServiceResponse` has been removed
- Field `XMSCorrelationRequestID`, `XMSRequestID` of struct `RegisteredServersClientGetResponse` has been removed
- Field `XMSCorrelationRequestID`, `XMSRequestID` of struct `RegisteredServersClientListByStorageSyncServiceResponse` has been removed
- Field `Parameters` of struct `ServerEndpointsClientBeginUpdateOptions` has been removed
- Field `XMSCorrelationRequestID`, `XMSRequestID` of struct `ServerEndpointsClientGetResponse` has been removed
- Field `XMSCorrelationRequestID`, `XMSRequestID` of struct `ServerEndpointsClientListBySyncGroupResponse` has been removed
- Field `Parameters` of struct `ServicesClientBeginUpdateOptions` has been removed
- Field `XMSCorrelationRequestID`, `XMSRequestID` of struct `ServicesClientGetResponse` has been removed
- Field `XMSCorrelationRequestID`, `XMSRequestID` of struct `ServicesClientListByResourceGroupResponse` has been removed
- Field `XMSCorrelationRequestID`, `XMSRequestID` of struct `ServicesClientListBySubscriptionResponse` has been removed
- Field `XMSCorrelationRequestID`, `XMSRequestID` of struct `SyncGroupsClientCreateResponse` has been removed
- Field `XMSCorrelationRequestID`, `XMSRequestID` of struct `SyncGroupsClientDeleteResponse` has been removed
- Field `XMSCorrelationRequestID`, `XMSRequestID` of struct `SyncGroupsClientGetResponse` has been removed
- Field `XMSCorrelationRequestID`, `XMSRequestID` of struct `SyncGroupsClientListByStorageSyncServiceResponse` has been removed
- Field `XMSCorrelationRequestID`, `XMSRequestID` of struct `WorkflowsClientAbortResponse` has been removed
- Field `XMSCorrelationRequestID`, `XMSRequestID` of struct `WorkflowsClientGetResponse` has been removed
- Field `XMSCorrelationRequestID`, `XMSRequestID` of struct `WorkflowsClientListByStorageSyncServiceResponse` has been removed

### Features Added

- New enum type `CloudTieringLowDiskModeState` with values `CloudTieringLowDiskModeStateDisabled`, `CloudTieringLowDiskModeStateEnabled`
- New enum type `ManagedServiceIdentityType` with values `ManagedServiceIdentityTypeNone`, `ManagedServiceIdentityTypeSystemAssigned`, `ManagedServiceIdentityTypeSystemAssignedUserAssigned`, `ManagedServiceIdentityTypeUserAssigned`
- New enum type `ServerAuthType` with values `ServerAuthTypeCertificate`, `ServerAuthTypeManagedIdentity`
- New enum type `ServerProvisioningStatus` with values `ServerProvisioningStatusError`, `ServerProvisioningStatusInProgress`, `ServerProvisioningStatusNotStarted`, `ServerProvisioningStatusReadySyncFunctional`, `ServerProvisioningStatusReadySyncNotFunctional`
- New function `*CloudEndpointsClient.AfsShareMetadataCertificatePublicKeys(ctx context.Context, resourceGroupName string, storageSyncServiceName string, syncGroupName string, cloudEndpointName string, options *CloudEndpointsClientAfsShareMetadataCertificatePublicKeysOptions) (CloudEndpointsClientAfsShareMetadataCertificatePublicKeysResponse, error)`
- New function `*RegisteredServersClient.BeginUpdate(ctx context.Context, resourceGroupName string, storageSyncServiceName string, serverID string, parameters RegisteredServerUpdateParameters, options *RegisteredServersClientBeginUpdateOptions) (*runtime.Poller[RegisteredServersClientUpdateResponse], error)`
- New struct `CloudEndpointAfsShareMetadataCertificatePublicKeys`
- New struct `CloudTieringLowDiskMode`
- New struct `ManagedServiceIdentity`
- New struct `RegisteredServerUpdateParameters`
- New struct `RegisteredServerUpdateProperties`
- New struct `ServerEndpointProvisioningStatus`
- New struct `ServerEndpointProvisioningStepStatus`
- New struct `UserAssignedIdentity`
- New field `NextLink` in struct `CloudEndpointArray`
- New field `CorrelationRequestID`, `RequestID` in struct `CloudEndpointsClientGetResponse`
- New field `CorrelationRequestID`, `RequestID` in struct `CloudEndpointsClientListBySyncGroupResponse`
- New field `CorrelationRequestID`, `RequestID` in struct `CloudEndpointsClientRestoreheartbeatResponse`
- New field `CorrelationRequestID`, `RequestID` in struct `MicrosoftStorageSyncClientLocationOperationStatusResponse`
- New field `LockAggregationType` in struct `OperationResourceMetricSpecification`
- New field `CorrelationRequestID`, `RequestID` in struct `OperationStatusClientGetResponse`
- New field `CorrelationRequestID`, `RequestID` in struct `OperationsClientListResponse`
- New field `NextLink` in struct `PrivateEndpointConnectionListResult`
- New field `GroupIDs` in struct `PrivateEndpointConnectionProperties`
- New field `CorrelationRequestID`, `RequestID` in struct `PrivateEndpointConnectionsClientListByStorageSyncServiceResponse`
- New field `NextLink` in struct `RegisteredServerArray`
- New field `ApplicationID`, `Identity` in struct `RegisteredServerCreateParametersProperties`
- New field `ActiveAuthType`, `ApplicationID`, `Identity`, `LatestApplicationID` in struct `RegisteredServerProperties`
- New field `CorrelationRequestID`, `RequestID` in struct `RegisteredServersClientGetResponse`
- New field `CorrelationRequestID`, `RequestID` in struct `RegisteredServersClientListByStorageSyncServiceResponse`
- New field `NextLink` in struct `ServerEndpointArray`
- New field `LowDiskMode` in struct `ServerEndpointCloudTieringStatus`
- New field `ServerEndpointProvisioningStatus` in struct `ServerEndpointProperties`
- New field `CorrelationRequestID`, `RequestID` in struct `ServerEndpointsClientGetResponse`
- New field `CorrelationRequestID`, `RequestID` in struct `ServerEndpointsClientListBySyncGroupResponse`
- New field `Identity` in struct `Service`
- New field `NextLink` in struct `ServiceArray`
- New field `ID`, `Identity`, `Name`, `SystemData`, `Type` in struct `ServiceCreateParameters`
- New field `UseIdentity` in struct `ServiceCreateParametersProperties`
- New field `UseIdentity` in struct `ServiceProperties`
- New field `Identity` in struct `ServiceUpdateParameters`
- New field `UseIdentity` in struct `ServiceUpdateProperties`
- New field `CorrelationRequestID`, `RequestID` in struct `ServicesClientGetResponse`
- New field `CorrelationRequestID`, `RequestID` in struct `ServicesClientListByResourceGroupResponse`
- New field `CorrelationRequestID`, `RequestID` in struct `ServicesClientListBySubscriptionResponse`
- New field `NextLink` in struct `SyncGroupArray`
- New field `CorrelationRequestID`, `RequestID` in struct `SyncGroupsClientCreateResponse`
- New field `CorrelationRequestID`, `RequestID` in struct `SyncGroupsClientDeleteResponse`
- New field `CorrelationRequestID`, `RequestID` in struct `SyncGroupsClientGetResponse`
- New field `CorrelationRequestID`, `RequestID` in struct `SyncGroupsClientListByStorageSyncServiceResponse`
- New field `NextLink` in struct `WorkflowArray`
- New field `CorrelationRequestID`, `RequestID` in struct `WorkflowsClientAbortResponse`
- New field `CorrelationRequestID`, `RequestID` in struct `WorkflowsClientGetResponse`
- New field `CorrelationRequestID`, `RequestID` in struct `WorkflowsClientListByStorageSyncServiceResponse`


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

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagesync/armstoragesync` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).