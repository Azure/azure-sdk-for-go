## Go

``` yaml

clear-output-folder: false
export-clients: true
go: true
input-file: C:\Users\arpitsaxena\source\repos\azure-rest-api-specs-pr\specification\pki\Microsoft.Pki\preview\2022-09-01-preview\pki.json
license-header: MICROSOFT_MIT_NO_VERSION
module: github.com/Azure/azure-sdk-for-go/sdk/identity/azcertificatemanager
openapi-type: "data-plane"
output-folder: C:\Users\arpitsaxena\Documents\repos\Azure\azure-sdk-for-go\sdk\identity\azcertificatemanager
use: "@autorest/go"

directive:
  # delete unused error models
  - from: models.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+type (?:Error).+\{(?:\s.+\s)+\}\s/g, "");
  - from: models_serde.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+func \(\w \*?(?:Error)\).*\{\s(?:.+\s)+\}\s/g, "");

  # delete generated constructor and client
  - from: client.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+func NewClient.+\{\s(?:.+\s)+\}\s/, "");
  - from: client.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+type Client struct.+\{\s(?:.+\s)+\}\s/, "");
