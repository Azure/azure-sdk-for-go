## Go

```yaml
clear-output-folder: false
export-clients: true
go: true
input-file: 
    - https://github.com/Azure/azure-rest-api-specs/blob/a2f6f742d088dcc712e67cb2745d8271eaa370ff/specification/keyvault/data-plane/Microsoft.KeyVault/preview/7.5-preview.1/settings.json
license-header: MICROSOFT_MIT_NO_VERSION
openapi-type: "data-plane"
output-folder: ../settings
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

  # delete unused error models
  - from: models.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+type (?:Error|KeyVaultError).+\{(?:\s.+\s)+\}\s/g, "");
  - from: models_serde.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+func \(\w \*?(?:Error|KeyVaultError)\).*\{\s(?:.+\s)+\}\s/g, "");

  # delete client name prefix from method options and response types
  - from:
      - client.go
      - models.go
      - options.go
      - response_types.go
      - options.go
    where: $
    transform: return $.replace(/Client(\w+)((?:Options|Response))/g, "$1$2");
  
  # add doc comment for Setting
  - from: swagger-document
    where: $.definitions.Setting
    transform: $["description"] = "A Key Vault setting."

  # remane SettingTypeEnum to SettingType
  - from: swagger-document
    where: $.definitions.Setting.properties.type.x-ms-enum
    transform: $["name"] = "SettingType"

  # fix up span names
  - from: client.go
    where: $
    transform: return $.replace(/StartSpan\(ctx, "Client/g, "StartSpan(ctx, \"settings.Client");
```
