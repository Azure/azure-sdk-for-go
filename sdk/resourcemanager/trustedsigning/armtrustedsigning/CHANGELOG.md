# Release History

## 0.2.0 (2025-06-16)
### Breaking Changes

- Type of `CodeSigningAccountPatchProperties.SKU` has been changed from `*AccountSKU` to `*AccountSKUPatch`
- Field `City`, `CommonName`, `Country`, `EnhancedKeyUsage`, `Organization`, `OrganizationUnit`, `PostalCode`, `State`, `StreetAddress` of struct `CertificateProfileProperties` has been removed

### Features Added

- New struct `AccountSKUPatch`
- New field `EnhancedKeyUsage` in struct `Certificate`


## 0.1.0 (2024-09-29)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/trustedsigning/armtrustedsigning` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).