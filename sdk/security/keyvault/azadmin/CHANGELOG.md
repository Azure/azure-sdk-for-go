## Release History

### 1.2.1 (Unreleased)

#### Features Added
* Added API Version support. Users can now change the default API Version by setting ClientOptions.APIVersion

#### Breaking Changes

#### Bugs Fixed

#### Other Changes

### 1.2.0 (2024-10-21)

#### Features Added
* Added CAE support
* Client requests tokens from the Vault's tenant, overriding any credential default
  (thanks @francescomari)

### 1.1.0 (2024-02-13)

#### Other Changes
* Upgraded to API service version `7.5`
* Upgraded dependencies

### 1.1.0-beta.1 (2023-11-09)

#### Features Added
* Managed Identity can now be used in place of a SAS token to access the blob storage resource when performing backup and restore operations.

#### Other Changes
* Upgraded service version to `7.5-preview.1`
* Updated to latest version of `azcore`.
* Enabled spans for distributed tracing.

### 1.0.1 (2023-08-24)

#### Other Changes
* Upgraded dependencies 

### 1.0.0 (2023-07-17)

#### Features Added
* First stable release of the azadmin module

### 0.3.0 (2023-06-08)

### Breatking Changes
* Renamed `SASTokenParameter` to `SASTokenParameters`
* Renamed `RestoreOperationParameters.SasTokenParameters` to `RestoreOperationParameters.SASTokenParameters`

### Other Changes
* Updated dependencies

### 0.2.0 (2023-04-13)

#### Breaking Changes
* Renamed `ServerError` to `ErrorInfo`

### 0.1.0 (2023-03-13)
* This is the initial release of the `azadmin` library
