Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewKeyListResultPage` parameter(s) have been changed from `(func(context.Context, KeyListResult) (KeyListResult, error))` to `(KeyListResult, func(context.Context, KeyListResult) (KeyListResult, error))`
- Function `NewCertificateIssuerListResultPage` parameter(s) have been changed from `(func(context.Context, CertificateIssuerListResult) (CertificateIssuerListResult, error))` to `(CertificateIssuerListResult, func(context.Context, CertificateIssuerListResult) (CertificateIssuerListResult, error))`
- Function `NewSecretListResultPage` parameter(s) have been changed from `(func(context.Context, SecretListResult) (SecretListResult, error))` to `(SecretListResult, func(context.Context, SecretListResult) (SecretListResult, error))`
- Function `NewCertificateListResultPage` parameter(s) have been changed from `(func(context.Context, CertificateListResult) (CertificateListResult, error))` to `(CertificateListResult, func(context.Context, CertificateListResult) (CertificateListResult, error))`

## New Content

- New function `KeyAttributes.MarshalJSON() ([]byte, error)`
- New function `IssuerAttributes.MarshalJSON() ([]byte, error)`
- New function `CertificatePolicy.MarshalJSON() ([]byte, error)`
- New function `SecretAttributes.MarshalJSON() ([]byte, error)`
- New function `CertificateOperation.MarshalJSON() ([]byte, error)`
- New function `Contacts.MarshalJSON() ([]byte, error)`
- New function `Attributes.MarshalJSON() ([]byte, error)`
- New function `CertificateAttributes.MarshalJSON() ([]byte, error)`
- New function `IssuerBundle.MarshalJSON() ([]byte, error)`
