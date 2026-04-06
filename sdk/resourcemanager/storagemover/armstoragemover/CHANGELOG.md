# Release History

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