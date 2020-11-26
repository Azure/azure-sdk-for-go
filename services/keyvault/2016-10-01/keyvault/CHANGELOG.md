
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewDeletedCertificateListResultPage` signature has been changed from `(func(context.Context, DeletedCertificateListResult) (DeletedCertificateListResult, error))` to `(DeletedCertificateListResult,func(context.Context, DeletedCertificateListResult) (DeletedCertificateListResult, error))`
- Function `NewSasDefinitionListResultPage` signature has been changed from `(func(context.Context, SasDefinitionListResult) (SasDefinitionListResult, error))` to `(SasDefinitionListResult,func(context.Context, SasDefinitionListResult) (SasDefinitionListResult, error))`
- Function `NewCertificateListResultPage` signature has been changed from `(func(context.Context, CertificateListResult) (CertificateListResult, error))` to `(CertificateListResult,func(context.Context, CertificateListResult) (CertificateListResult, error))`
- Function `NewCertificateIssuerListResultPage` signature has been changed from `(func(context.Context, CertificateIssuerListResult) (CertificateIssuerListResult, error))` to `(CertificateIssuerListResult,func(context.Context, CertificateIssuerListResult) (CertificateIssuerListResult, error))`
- Function `NewSecretListResultPage` signature has been changed from `(func(context.Context, SecretListResult) (SecretListResult, error))` to `(SecretListResult,func(context.Context, SecretListResult) (SecretListResult, error))`
- Function `NewDeletedSecretListResultPage` signature has been changed from `(func(context.Context, DeletedSecretListResult) (DeletedSecretListResult, error))` to `(DeletedSecretListResult,func(context.Context, DeletedSecretListResult) (DeletedSecretListResult, error))`
- Function `NewDeletedKeyListResultPage` signature has been changed from `(func(context.Context, DeletedKeyListResult) (DeletedKeyListResult, error))` to `(DeletedKeyListResult,func(context.Context, DeletedKeyListResult) (DeletedKeyListResult, error))`
- Function `NewKeyListResultPage` signature has been changed from `(func(context.Context, KeyListResult) (KeyListResult, error))` to `(KeyListResult,func(context.Context, KeyListResult) (KeyListResult, error))`
- Function `NewStorageListResultPage` signature has been changed from `(func(context.Context, StorageListResult) (StorageListResult, error))` to `(StorageListResult,func(context.Context, StorageListResult) (StorageListResult, error))`

## New Content

- Function `IssuerAttributes.MarshalJSON() ([]byte,error)` is added
- Function `CertificateAttributes.MarshalJSON() ([]byte,error)` is added
- Function `SasDefinitionAttributes.MarshalJSON() ([]byte,error)` is added
- Function `Attributes.MarshalJSON() ([]byte,error)` is added
- Function `StorageAccountAttributes.MarshalJSON() ([]byte,error)` is added
- Function `SecretAttributes.MarshalJSON() ([]byte,error)` is added
- Function `Contacts.MarshalJSON() ([]byte,error)` is added
- Function `CertificateOperation.MarshalJSON() ([]byte,error)` is added
- Function `KeyAttributes.MarshalJSON() ([]byte,error)` is added
- Function `IssuerBundle.MarshalJSON() ([]byte,error)` is added
- Function `CertificatePolicy.MarshalJSON() ([]byte,error)` is added

