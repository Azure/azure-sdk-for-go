## Go

```yaml
clear-output-folder: false
export-clients: true
go: true
input-file: 
    - https://github.com/Azure/azure-rest-api-specs/blob/a2f6f742d088dcc712e67cb2745d8271eaa370ff/specification/keyvault/data-plane/Microsoft.KeyVault/preview/7.5-preview.1/backuprestore.json
license-header: MICROSOFT_MIT_NO_VERSION
openapi-type: "data-plane"
output-folder: ../backup
override-client-name: Client
security: "AADToken"
security-scopes: "https://vault.azure.net/.default"
use: "@autorest/go@4.0.0-preview.59"
inject-spans: true
version: "^3.0.0"

directive:

  # make vault URL a parameter of the client constructor
  - from: swagger-document
    where: $["x-ms-parameterized-host"]
    transform: $.parameters[0]["x-ms-parameter-location"] = "client"

  # rename restore operation from BeginFullRestoreOperation to BeginFullRestore
  - rename-operation:
      from: FullRestoreOperation
      to: FullRestore
  - rename-operation:
      from: SelectiveKeyRestoreOperation
      to: SelectiveKeyRestore

  # Fix SASToken names
  - rename-model: 
      from: SASTokenParameter
      to: SASTokenParameters
  - from: swagger-document
    where: $.definitions.RestoreOperationParameters.properties.sasTokenParameters
    transform: $["x-ms-client-name"] = "SASTokenParameters"
  - from: swagger-document
    where: $.definitions.SelectiveKeyRestoreOperationParameters.properties.sasTokenParameters
    transform: $["x-ms-client-name"] = "SASTokenParameters"

  # make SASToken parameter required
  - from: swagger-document
    where: $.paths["/backup"].post.parameters[0]
    transform: $["required"] = true

  # delete backup and restore status methods
  - from: swagger-document
    where: $["paths"]
    transform: >
        delete $["/backup/{jobId}/pending"];
  - from: swagger-document
    where: $["paths"]
    transform: >
        delete $["/restore/{jobId}/pending"];

  # delete unused error models
  - from: models.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+type (?:Error|KeyVaultError).+\{(?:\s.+\s)+\}\s/g, "");
  - from: models_serde.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+func \(\w \*?(?:Error|KeyVaultError)\).*\{\s(?:.+\s)+\}\s/g, "");
  - from: models.go
    where: $
    transform: return $.replace(/Error \*Error/g, "Error *ErrorInfo");

  # modify Restore to use implementation with custom poller handler
  - from: client.go
    where: $
    transform:  return $.replace(/\[ClientFullRestoreResponse\], error\) \{\s(?:.+\s)+\}/, "[ClientFullRestoreResponse], error) {return client.beginFullRestore(ctx, restoreBlobDetails, options)}");
  - from: client.go
    where: $
    transform:  return $.replace(/\[ClientSelectiveKeyRestoreResponse\], error\) \{\s(?:.+\s)+\}/, "[ClientSelectiveKeyRestoreResponse], error) {return client.beginSelectiveKeyRestore(ctx, keyName, restoreBlobDetails, options)}");

  # delete client name prefix from method options and response types
  - from:
      - client.go
      - models.go
      - options.go
      - response_types.go
      - options.go
    where: $
    transform: return $.replace(/Client(\w+)((?:Options|Response))/g, "$1$2");

  # add doc comments for models with missing descriptions
  - from: swagger-document
    where: $.definitions.SASTokenParameters
    transform: $["description"] = "Contains the information required to access blob storage."
  - from: swagger-document
    where: $.definitions.RestoreOperationParameters
    transform: $["description"] = "Parameters for the restore operation"
  - from: swagger-document
    where: $.definitions.SelectiveKeyRestoreOperationParameters
    transform: $["description"] = "Parameters for the selective restore operation"

  # fix up span names
  - from: client.go
    where: $
    transform: return $.replace(/StartSpan\(ctx, "Client/g, "StartSpan(ctx, \"backup.Client");
```