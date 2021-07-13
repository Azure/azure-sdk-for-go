# Unreleased

## Breaking Changes

### Struct Changes

#### Removed Struct Fields

1. Export.BlobListblobPath

### Signature Changes

#### Struct Fields

1. DriveStatus.PercentComplete changed type from *int32 to *int64
1. JobDetails.DeliveryPackage changed type from *PackageInfomation to *DeliveryPackageInformation
1. UpdateJobParametersProperties.DeliveryPackage changed type from *PackageInfomation to *DeliveryPackageInformation

## Additive Changes

### New Constants

1. CreatedByType.Application
1. CreatedByType.Key
1. CreatedByType.ManagedIdentity
1. CreatedByType.User
1. EncryptionKekType.CustomerManaged
1. EncryptionKekType.MicrosoftManaged
1. IdentityType.None
1. IdentityType.SystemAssigned
1. IdentityType.UserAssigned

### New Funcs

1. IdentityDetails.MarshalJSON() ([]byte, error)
1. PossibleCreatedByTypeValues() []CreatedByType
1. PossibleEncryptionKekTypeValues() []EncryptionKekType
1. PossibleIdentityTypeValues() []IdentityType
1. ShippingInformation.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. DeliveryPackageInformation
1. EncryptionKeyDetails
1. IdentityDetails
1. SystemData

#### New Struct Fields

1. Export.BlobListBlobPath
1. JobDetails.EncryptionKey
1. JobResponse.Identity
1. JobResponse.SystemData
1. LocationProperties.AdditionalShippingInformation
1. ShippingInformation.AdditionalInformation
