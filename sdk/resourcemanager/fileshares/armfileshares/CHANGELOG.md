# Release History

## 1.0.0 (2026-10-29)
### Features Added

- New value `MediaTierHDD` added to enum type `MediaTier`
- New enum type `InfrastructureEncryption` with values `InfrastructureEncryptionDisabled`, `InfrastructureEncryptionEnabled`
- New enum type `KeyEncryptionKeyIdentityType` with values `KeyEncryptionKeyIdentityTypeDelegatedResourceIdentity`, `KeyEncryptionKeyIdentityTypeSystemAssignedIdentity`, `KeyEncryptionKeyIdentityTypeUserAssignedIdentity`
- New enum type `ManagedServiceIdentityType` with values `ManagedServiceIdentityTypeNone`, `ManagedServiceIdentityTypeSystemAssigned`, `ManagedServiceIdentityTypeSystemAssignedUserAssigned`, `ManagedServiceIdentityTypeUserAssigned`
- New enum type `RedundancyConversionStatus` with values `RedundancyConversionStatusComplete`, `RedundancyConversionStatusFailed`, `RedundancyConversionStatusInProgress`, `RedundancyConversionStatusSubmittedForConversion`
- New enum type `ScaleMode` with values `ScaleModeCompatibility`, `ScaleModeScale`
- New function `*Client.GetRedundancyConversionStatus(ctx context.Context, resourceGroupName string, resourceName string, options *ClientGetRedundancyConversionStatusOptions) (ClientGetRedundancyConversionStatusResponse, error)`
- New function `*Client.StartRedundancyConversion(ctx context.Context, resourceGroupName string, resourceName string, body StartRedundancyConversionRequest, options *ClientStartRedundancyConversionOptions) (ClientStartRedundancyConversionResponse, error)`
- New struct `CustomerManagedKeyEncryption`
- New struct `ErrorAdditionalInfo`
- New struct `ErrorDetail`
- New struct `FileShareEncryption`
- New struct `FileShareLimitsRequest`
- New struct `FileShareOfferInput`
- New struct `KeyEncryptionKeyIdentity`
- New struct `ManagedServiceIdentity`
- New struct `MigrationState`
- New struct `RedundancyConversionStatusOutput`
- New struct `RedundancyConversionStatusResponse`
- New struct `StartRedundancyConversionInput`
- New struct `StartRedundancyConversionRequest`
- New struct `UserAssignedIdentity`
- New field `Identity` in struct `FileShare`
- New field `Encryption`, `MigrationState`, `ScaleMode` in struct `FileShareProperties`
- New field `MediaTier`, `Protocol`, `ScaleMode` in struct `FileShareProvisioningRecommendationInput`
- New field `Identity` in struct `FileShareUpdate`
- New field `Encryption` in struct `FileShareUpdateProperties`
- New field `Body` in struct `InformationalOperationsClientGetLimitsOptions`


## 0.1.0 (2026-05-07)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/fileshares/armfileshares` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).