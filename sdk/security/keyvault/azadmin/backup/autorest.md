## Go

```yaml
clear-output-folder: false
export-clients: true
go: true
input-file: 
    - https://github.com/Azure/azure-rest-api-specs/blob/main/specification/keyvault/data-plane/Microsoft.KeyVault/stable/7.4/backuprestore.json
license-header: MICROSOFT_MIT_NO_VERSION
openapi-type: "data-plane"
output-folder: ../backup
override-client-name: Client
security: "AADToken"
security-scopes: "https://vault.azure.net/.default"
use: "@autorest/go@4.0.0-preview.44"
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

  # delete generated client constructor
  - from: client.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+func NewClient.+\{\s(?:.+\s)+\}\s/, "");

  # delete unused error models
  - from: models.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+type (?:Error|KeyVaultError).+\{(?:\s.+\s)+\}\s/g, "");
  - from: models_serde.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+func \(\w \*?(?:Error|KeyVaultError)\).*\{\s(?:.+\s)+\}\s/g, "");
  - from: models.go
    where: $
    transform: return $.replace(/Error \*Error/g, "Error *ServerError");

  # modify Restore to use implementation with custom poller handler
  - from: client.go
    where: $
    transform:  return $.replace(/\[ClientFullRestoreResponse\], error\) \{\s(?:.+\s)+\}/, "[ClientFullRestoreResponse], error) {return client.beginFullRestore(ctx, restoreBlobDetails, options)}");
  - from: client.go
    where: $
    transform:  return $.replace(/\[ClientSelectiveKeyRestoreResponse\], error\) \{\s(?:.+\s)+\}/, "[ClientSelectiveKeyRestoreResponse], error) {return client.beginSelectiveKeyRestore(ctx, keyName, restoreBlobDetails, options)}");
```