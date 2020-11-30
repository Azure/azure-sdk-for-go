Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewDeletedSasDefinitionListResultPage` parameter(s) have been changed from `(func(context.Context, DeletedSasDefinitionListResult) (DeletedSasDefinitionListResult, error))` to `(DeletedSasDefinitionListResult, func(context.Context, DeletedSasDefinitionListResult) (DeletedSasDefinitionListResult, error))`
- Function `NewKeyListResultPage` parameter(s) have been changed from `(func(context.Context, KeyListResult) (KeyListResult, error))` to `(KeyListResult, func(context.Context, KeyListResult) (KeyListResult, error))`
- Function `NewCertificateIssuerListResultPage` parameter(s) have been changed from `(func(context.Context, CertificateIssuerListResult) (CertificateIssuerListResult, error))` to `(CertificateIssuerListResult, func(context.Context, CertificateIssuerListResult) (CertificateIssuerListResult, error))`
- Function `NewSasDefinitionListResultPage` parameter(s) have been changed from `(func(context.Context, SasDefinitionListResult) (SasDefinitionListResult, error))` to `(SasDefinitionListResult, func(context.Context, SasDefinitionListResult) (SasDefinitionListResult, error))`
- Function `NewDeletedStorageListResultPage` parameter(s) have been changed from `(func(context.Context, DeletedStorageListResult) (DeletedStorageListResult, error))` to `(DeletedStorageListResult, func(context.Context, DeletedStorageListResult) (DeletedStorageListResult, error))`
- Function `NewSecretListResultPage` parameter(s) have been changed from `(func(context.Context, SecretListResult) (SecretListResult, error))` to `(SecretListResult, func(context.Context, SecretListResult) (SecretListResult, error))`
- Function `NewDeletedCertificateListResultPage` parameter(s) have been changed from `(func(context.Context, DeletedCertificateListResult) (DeletedCertificateListResult, error))` to `(DeletedCertificateListResult, func(context.Context, DeletedCertificateListResult) (DeletedCertificateListResult, error))`
- Function `NewDeletedKeyListResultPage` parameter(s) have been changed from `(func(context.Context, DeletedKeyListResult) (DeletedKeyListResult, error))` to `(DeletedKeyListResult, func(context.Context, DeletedKeyListResult) (DeletedKeyListResult, error))`
- Function `NewDeletedSecretListResultPage` parameter(s) have been changed from `(func(context.Context, DeletedSecretListResult) (DeletedSecretListResult, error))` to `(DeletedSecretListResult, func(context.Context, DeletedSecretListResult) (DeletedSecretListResult, error))`
- Function `NewCertificateListResultPage` parameter(s) have been changed from `(func(context.Context, CertificateListResult) (CertificateListResult, error))` to `(CertificateListResult, func(context.Context, CertificateListResult) (CertificateListResult, error))`
- Function `NewStorageListResultPage` parameter(s) have been changed from `(func(context.Context, StorageListResult) (StorageListResult, error))` to `(StorageListResult, func(context.Context, StorageListResult) (StorageListResult, error))`

## New Content

- New const `CustomizedRecoverablePurgeable`
- New const `CustomizedRecoverableProtectedSubscription`
- New const `CustomizedRecoverable`
- New function `IssuerBundle.MarshalJSON() ([]byte, error)`
- New function `CertificatePolicy.MarshalJSON() ([]byte, error)`
- New function `IssuerAttributes.MarshalJSON() ([]byte, error)`
- New function `Attributes.MarshalJSON() ([]byte, error)`
- New function `SecretAttributes.MarshalJSON() ([]byte, error)`
- New function `SasDefinitionAttributes.MarshalJSON() ([]byte, error)`
- New function `CertificateOperation.MarshalJSON() ([]byte, error)`
- New function `CertificateAttributes.MarshalJSON() ([]byte, error)`
- New function `Contacts.MarshalJSON() ([]byte, error)`
- New function `KeyAttributes.MarshalJSON() ([]byte, error)`
- New function `StorageAccountAttributes.MarshalJSON() ([]byte, error)`
