## Go

```yaml
clear-output-folder: false
export-clients: true
go: true
input-file: 
    - https://github.com/Azure/azure-rest-api-specs/blob/main/specification/keyvault/data-plane/Microsoft.KeyVault/stable/7.4/settings.json
license-header: MICROSOFT_MIT_NO_VERSION
openapi-type: "data-plane"
output-folder: ../settings
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

  # fix bug- change ListResult.Value to ListResult.Settings
  - where-model: SettingsListResult
    rename-property:
      from: value
      to: settings

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

  # fix bug- change ListResult.Value to ListResult.Settings

```