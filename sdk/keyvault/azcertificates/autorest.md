## Go

```yaml
clear-output-folder: false
export-clients: true
go: true
input-file: https://github.com/Azure/azure-rest-api-specs/blob/37cd8dfac3c570a24bb645b31c012d12efb760df/specification/keyvault/data-plane/Microsoft.KeyVault/stable/7.3/certificates.json
license-header: MICROSOFT_MIT_NO_VERSION
module: github.com/Azure/azure-sdk-for-go/sdk/keyvault/azcertificates
openapi-type: "data-plane"
output-folder: ../azcertificates
override-client-name: Client
security: "AADToken"
security-scopes: "https://vault.azure.net/.default"
use: "@autorest/go@4.0.0-preview.43"
version: "^3.0.0"

directive:
  # delete unused model
  - remove-model: PendingCertificateSigningRequestResult

  # make vault URL a parameter of the client constructor
  - from: swagger-document
    where: $["x-ms-parameterized-host"]
    transform: $.parameters[0]["x-ms-parameter-location"] = "client"

  # rename parameter models to match their methods
  - rename-model:
      from: CertificateCreateParameters
      to: CreateCertificateParameters
  - rename-model:
      from: CertificateImportParameters
      to: ImportCertificateParameters
  - rename-model:
      from: CertificateIssuerSetParameters
      to: SetCertificateIssuerParameters
  - rename-model:
      from: CertificateIssuerUpdateParameters
      to: UpdateCertificateIssuerParameters
  - rename-model:
      from: CertificateMergeParameters
      to: MergeCertificateParameters
  - rename-model:
      from: CertificateOperationUpdateParameter
      to: UpdateCertificateOperationParameter
  - rename-model:
      from: CertificateRestoreParameters
      to: RestoreCertificateParameters
  - rename-model:
      from: CertificateUpdateParameters
      to: UpdateCertificateParameters

  # rename paged operations from Get* to List*
  - rename-operation:
      from: GetCertificates
      to: ListCertificates
  - rename-operation:
      from: GetCertificateIssuers
      to: ListCertificateIssuers
  - rename-operation:
      from: GetCertificateVersions
      to: ListCertificateVersions
  - rename-operation:
      from: GetDeletedCertificates
      to: ListDeletedCertificates

  # Maxresults -> MaxResults
  - from: swagger-document
    where: $.paths..parameters..[?(@.name=='maxresults')]
    transform: $["x-ms-client-name"] = "MaxResults"

  # capitalize acronyms
  - where-model: CertificateBundle
    transform: $.properties.cer["x-ms-client-name"] = "CER"
  - where-model: CertificateBundle
    transform: $.properties.kid["x-ms-client-name"] = "KID"
  - where-model: CertificateBundle
    transform: $.properties.sid["x-ms-client-name"] = "SID"
  - where-model: CertificateOperation
    transform: $.properties.csr["x-ms-client-name"] = "CSR"
  - where-model: SubjectAlternativeNames
    transform: $.properties.upns["x-ms-client-name"] = "UPNs"
  - where-model: X509CertificateProperties
    transform: $.properties.ekus["x-ms-client-name"] = "EKUs"

  # delete unused KeyVaultError
  - from: models.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+type KeyVaultError.+\{(?:\s.+\s)+\}\s/g, "");
  - from: models_serde.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+func \(\w \*?KeyVaultError\).*\{\s(?:.+\s)+\}\s/g, "");

  # delete the Attributes model defined in common.json (it's used only with allOf)
  - from: models.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+type Attributes.+\{(?:\s.+\s)+\}\s/, "");
  - from: models_serde.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+func \(a \*?Attributes\).*\{\s(?:.+\s)+\}\s/g, "");

  # delete generated constructor
  - from: client.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+func NewClient.+\{\s(?:.+\s)+\}\s/, "");

  # delete the version path param check (version == "" is legal for Key Vault but indescribable by OpenAPI)
  - from: client.go
    where: $
    transform: return $.replace(/\sif certificateVersion == "" \{\s+.+certificateVersion cannot be empty"\)\s+\}\s/g, "");

  # delete client name prefix from method options and response types
  - from:
      - client.go
      - models.go
      - response_types.go
    where: $
    transform: return $.replace(/Client(\w+)((?:Options|Response))/g, "$1$2");

  # make cert IDs a convenience type so we can add parsing methods
  # (specifying models because others have "ID" fields whose values aren't cert IDs)
  - from: models.go
    where: $
    transform: return $.replace(/(type (?:Deleted)?Certificate(?:Bundle|Item) struct \{(?:\s.+\s)+\sID \*)string/g, "$1ID")

  # remove "certificate" prefix from some method parameter names
  - from: client.go
  - where: $
  - transform: return $.replace(/certificate((?:Name|Policy|Version)) string/g, (match) => { return match[0].toLowerCase() + match.substr(1); })
```
