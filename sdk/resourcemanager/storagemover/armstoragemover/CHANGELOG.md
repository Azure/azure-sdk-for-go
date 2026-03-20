# Release History

## 2.4.0 (2026-03-17)
### Features Added

- New value `CredentialTypeAzureKeyVaultS3WithHMAC` added to enum type `CredentialType`
- New value `EndpointTypeS3WithHmac` added to enum type `EndpointType`
- New enum type `ConnectionStatus` with values `ConnectionStatusApproved`, `ConnectionStatusDisconnected`, `ConnectionStatusPending`, `ConnectionStatusRejected`, `ConnectionStatusStale`
- New enum type `DataIntegrityValidation` with values `DataIntegrityValidationNone`, `DataIntegrityValidationSaveFileMD5`, `DataIntegrityValidationSaveVerifyFileMD5`
- New enum type `EndpointKind` with values `EndpointKindSource`, `EndpointKindTarget`
- New enum type `Frequency` with values `FrequencyDaily`, `FrequencyMonthly`, `FrequencyOnetime`, `FrequencyWeekly`
- New enum type `S3WithHmacSourceType` with values `S3WithHmacSourceTypeBACKBLAZE`, `S3WithHmacSourceTypeCLOUDFLARE`, `S3WithHmacSourceTypeGCS`, `S3WithHmacSourceTypeIBM`, `S3WithHmacSourceTypeMINIO`
- New enum type `TriggerType` with values `TriggerTypeManual`, `TriggerTypeScheduled`
- New function `*AzureKeyVaultS3WithHmacCredentials.GetCredentials() *Credentials`
- New function `*ClientFactory.NewConnectionsClient() *ConnectionsClient`
- New function `NewConnectionsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ConnectionsClient, error)`
- New function `*ConnectionsClient.CreateOrUpdate(ctx context.Context, resourceGroupName string, storageMoverName string, connectionName string, connection Connection, options *ConnectionsClientCreateOrUpdateOptions) (ConnectionsClientCreateOrUpdateResponse, error)`
- New function `*ConnectionsClient.BeginDelete(ctx context.Context, resourceGroupName string, storageMoverName string, connectionName string, options *ConnectionsClientBeginDeleteOptions) (*runtime.Poller[ConnectionsClientDeleteResponse], error)`
- New function `*ConnectionsClient.Get(ctx context.Context, resourceGroupName string, storageMoverName string, connectionName string, options *ConnectionsClientGetOptions) (ConnectionsClientGetResponse, error)`
- New function `*ConnectionsClient.NewListPager(resourceGroupName string, storageMoverName string, options *ConnectionsClientListOptions) *runtime.Pager[ConnectionsClientListResponse]`
- New function `*S3WithHmacEndpointProperties.GetEndpointBaseProperties() *EndpointBaseProperties`
- New function `*S3WithHmacEndpointUpdateProperties.GetEndpointBaseUpdateProperties() *EndpointBaseUpdateProperties`
- New struct `AzureKeyVaultS3WithHmacCredentials`
- New struct `Connection`
- New struct `ConnectionList`
- New struct `ConnectionProperties`
- New struct `JobRunWarning`
- New struct `S3WithHmacEndpointProperties`
- New struct `S3WithHmacEndpointUpdateProperties`
- New struct `ScheduleInfo`
- New field `EndpointKind` in struct `AzureMultiCloudConnectorEndpointProperties`
- New field `EndpointKind` in struct `AzureStorageBlobContainerEndpointProperties`
- New field `EndpointKind` in struct `AzureStorageNfsFileShareEndpointProperties`
- New field `EndpointKind` in struct `AzureStorageSmbFileShareEndpointProperties`
- New field `Connections`, `DataIntegrityValidation`, `PreservePermissions`, `Schedule` in struct `JobDefinitionProperties`
- New field `Connections`, `DataIntegrityValidation` in struct `JobDefinitionUpdateProperties`
- New field `ScheduledExecutionTime`, `TriggerType`, `Warnings` in struct `JobRunProperties`
- New field `EndpointKind` in struct `NfsMountEndpointProperties`
- New field `EndpointKind` in struct `SmbMountEndpointProperties`


## 2.3.0 (2025-08-29)
### Features Added

- New value `EndpointTypeAzureMultiCloudConnector`, `EndpointTypeAzureStorageNfsFileShare` added to enum type `EndpointType`
- New enum type `JobType` with values `JobTypeCloudToCloud`, `JobTypeOnPremToCloud`
- New enum type `ManagedServiceIdentityType` with values `ManagedServiceIdentityTypeNone`, `ManagedServiceIdentityTypeSystemAssigned`, `ManagedServiceIdentityTypeSystemAssignedUserAssigned`, `ManagedServiceIdentityTypeUserAssigned`
- New function `*AzureMultiCloudConnectorEndpointProperties.GetEndpointBaseProperties() *EndpointBaseProperties`
- New function `*AzureMultiCloudConnectorEndpointUpdateProperties.GetEndpointBaseUpdateProperties() *EndpointBaseUpdateProperties`
- New function `*AzureStorageNfsFileShareEndpointProperties.GetEndpointBaseProperties() *EndpointBaseProperties`
- New function `*AzureStorageNfsFileShareEndpointUpdateProperties.GetEndpointBaseUpdateProperties() *EndpointBaseUpdateProperties`
- New struct `AzureMultiCloudConnectorEndpointProperties`
- New struct `AzureMultiCloudConnectorEndpointUpdateProperties`
- New struct `AzureStorageNfsFileShareEndpointProperties`
- New struct `AzureStorageNfsFileShareEndpointUpdateProperties`
- New struct `JobDefinitionPropertiesSourceTargetMap`
- New struct `ManagedServiceIdentity`
- New struct `SourceEndpoint`
- New struct `SourceEndpointProperties`
- New struct `SourceTargetMap`
- New struct `TargetEndpoint`
- New struct `TargetEndpointProperties`
- New struct `UserAssignedIdentity`
- New field `Identity` in struct `Endpoint`
- New field `Identity` in struct `EndpointBaseUpdateParameters`
- New field `JobType`, `SourceTargetMap` in struct `JobDefinitionProperties`


## 2.2.0 (2024-06-21)
### Features Added

- New value `JobRunStatusPausedByBandwidthManagement` added to enum type `JobRunStatus`
- New value `ProvisioningStateCanceled`, `ProvisioningStateDeleting`, `ProvisioningStateFailed` added to enum type `ProvisioningState`
- New enum type `DayOfWeek` with values `DayOfWeekFriday`, `DayOfWeekMonday`, `DayOfWeekSaturday`, `DayOfWeekSunday`, `DayOfWeekThursday`, `DayOfWeekTuesday`, `DayOfWeekWednesday`
- New enum type `Minute` with values `MinuteThirty`, `MinuteZero`
- New struct `Time`
- New struct `UploadLimitSchedule`
- New struct `UploadLimitWeeklyRecurrence`
- New field `TimeZone`, `UploadLimitSchedule` in struct `AgentProperties`
- New field `UploadLimitSchedule` in struct `AgentUpdateProperties`


## 2.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.0.0 (2023-10-27)
### Breaking Changes

- Type of `EndpointBaseUpdateParameters.Properties` has been changed from `*EndpointBaseUpdateProperties` to `EndpointBaseUpdatePropertiesClassification`

### Features Added

- New value `EndpointTypeAzureStorageSmbFileShare`, `EndpointTypeSmbMount` added to enum type `EndpointType`
- New enum type `CredentialType` with values `CredentialTypeAzureKeyVaultSmb`
- New function `*AzureKeyVaultSmbCredentials.GetCredentials() *Credentials`
- New function `*AzureStorageBlobContainerEndpointUpdateProperties.GetEndpointBaseUpdateProperties() *EndpointBaseUpdateProperties`
- New function `*AzureStorageSmbFileShareEndpointProperties.GetEndpointBaseProperties() *EndpointBaseProperties`
- New function `*AzureStorageSmbFileShareEndpointUpdateProperties.GetEndpointBaseUpdateProperties() *EndpointBaseUpdateProperties`
- New function `*Credentials.GetCredentials() *Credentials`
- New function `*EndpointBaseUpdateProperties.GetEndpointBaseUpdateProperties() *EndpointBaseUpdateProperties`
- New function `*SmbMountEndpointProperties.GetEndpointBaseProperties() *EndpointBaseProperties`
- New function `*SmbMountEndpointUpdateProperties.GetEndpointBaseUpdateProperties() *EndpointBaseUpdateProperties`
- New function `*NfsMountEndpointUpdateProperties.GetEndpointBaseUpdateProperties() *EndpointBaseUpdateProperties`
- New struct `AzureKeyVaultSmbCredentials`
- New struct `AzureStorageBlobContainerEndpointUpdateProperties`
- New struct `AzureStorageSmbFileShareEndpointProperties`
- New struct `AzureStorageSmbFileShareEndpointUpdateProperties`
- New struct `NfsMountEndpointUpdateProperties`
- New struct `SmbMountEndpointProperties`
- New struct `SmbMountEndpointUpdateProperties`


## 2.0.0-beta.1 (2023-07-28)
### Breaking Changes

- Type of `EndpointBaseUpdateParameters.Properties` has been changed from `*EndpointBaseUpdateProperties` to `EndpointBaseUpdatePropertiesClassification`
### Features Added
- New value `EndpointTypeAzureStorageSmbFileShare`, `EndpointTypeSmbMount` added to enum type `EndpointType`
- New enum type `CredentialType` with values `CredentialTypeAzureKeyVaultSmb`
- New function `*AzureKeyVaultSmbCredentials.GetCredentials() *Credentials`
- New function `*AzureStorageBlobContainerEndpointUpdateProperties.GetEndpointBaseUpdateProperties() *EndpointBaseUpdateProperties`
- New function `*AzureStorageSmbFileShareEndpointProperties.GetEndpointBaseProperties() *EndpointBaseProperties`
- New function `*AzureStorageSmbFileShareEndpointUpdateProperties.GetEndpointBaseUpdateProperties() *EndpointBaseUpdateProperties`
- New function `*Credentials.GetCredentials() *Credentials`
- New function `*EndpointBaseUpdateProperties.GetEndpointBaseUpdateProperties() *EndpointBaseUpdateProperties`
- New function `*SmbMountEndpointProperties.GetEndpointBaseProperties() *EndpointBaseProperties`
- New function `*SmbMountEndpointUpdateProperties.GetEndpointBaseUpdateProperties() *EndpointBaseUpdateProperties`
- New function `*NfsMountEndpointUpdateProperties.GetEndpointBaseUpdateProperties() *EndpointBaseUpdateProperties`
- New struct `AzureKeyVaultSmbCredentials`
- New struct `AzureStorageBlobContainerEndpointUpdateProperties`
- New struct `AzureStorageSmbFileShareEndpointProperties`
- New struct `AzureStorageSmbFileShareEndpointUpdateProperties`
- New struct `NfsMountEndpointUpdateProperties`
- New struct `SmbMountEndpointProperties`
- New struct `SmbMountEndpointUpdateProperties`


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2023-03-07)
### Other Changes

- Release stable version.

## 0.1.0 (2023-02-24)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagemover/armstoragemover` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).