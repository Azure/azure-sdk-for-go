## Go

```yaml
clear-output-folder: false
export-clients: true
go: true
input-file: https://github.com/Azure/azure-rest-api-specs/blob/7452e1cc7db72fbc6cd9539b390d8b8e5c2a1864/specification/keyvault/data-plane/Microsoft.KeyVault/stable/7.5/certificates.json
license-header: MICROSOFT_MIT_NO_VERSION
module: github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates
openapi-type: "data-plane"
output-folder: ../azcertificates
override-client-name: Client
security: "AADToken"
security-scopes: "https://vault.azure.net/.default"
use: "@autorest/go@4.0.0-preview.59"
inject-spans: true
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
      to: ListCertificateProperties
  - rename-operation:
      from: GetCertificateIssuers
      to: ListIssuerProperties
  - rename-operation:
      from: GetCertificateVersions
      to: ListCertificatePropertiesVersions
  - rename-operation:
      from: GetDeletedCertificates
      to: ListDeletedCertificateProperties
  - rename-model:
      from: CertificateListResult
      to: CertificatePropertiesListResult
  - rename-model:
      from: DeletedCertificateListResult
      to: DeletedCertificatePropertiesListResult

  # remove redunant "certificate" from operation name
  - rename-operation:
      from: SetCertificateContacts
      to: SetContacts
  - rename-operation:
      from: GetCertificateContacts
      to: GetContacts
  - rename-operation:
      from: DeleteCertificateContacts
      to: DeleteContacts
  - rename-operation:
      from: SetCertificateIssuer
      to: SetIssuer
  - rename-operation:
      from: UpdateCertificateIssuer
      to: UpdateIssuer
  - rename-operation:
      from: GetCertificateIssuer
      to: GetIssuer
  - rename-operation:
      from: DeleteCertificateIssuer
      to: DeleteIssuer
  - rename-model:
      from: CertificateIssuerListResult
      to: IssuerPropertiesListResult
  - rename-model:
      from: UpdateCertificateIssuerParameters
      to: UpdateIssuerParameters
  - rename-model:
      from: SetCertificateIssuerParameters
      to: SetIssuerParameters

  # rename LifetimeAction
  - rename-model:
      from: Action
      to: LifetimeActionType
  - rename-model:
      from: Trigger
      to: LifetimeActionTrigger
  
  # rename CertificateBundle, CertificateItem, IssuerBundle
  - rename-model:
      from: CertificateBundle
      to: Certificate
  - rename-model:
      from: CertificateItem
      to: CertificateProperties
  - rename-model:
      from: DeletedCertificateBundle
      to: DeletedCertificate
  - rename-model:
      from: DeletedCertificateItem
      to: DeletedCertificateProperties
  - rename-model:
      from: IssuerBundle
      to: Issuer
  - rename-model:
      from: CertificateIssuerItem
      to: IssuerProperties
  - where-model: RestoreCertificateParameters
    transform: $.properties.value["x-ms-client-name"] = "CertificateBackup"

  # rename AdministratorDetails to AdministratorContact
  - rename-model:
      from: AdministratorDetails
      to: AdministratorContact
  - where-model: OrganizationDetails
    transform: $.properties.admin_details["x-ms-client-name"] = "AdminContacts"
  - where-model: 
      - Contact
      - AdministratorContact
    transform: $.properties.email["x-ms-client-name"] = "Email"

  # rename UPNs to UserPrincipalNames
  - where-model: SubjectAlternativeNames
    transform: $.properties.upns["x-ms-client-name"] = "UserPrincipalNames"

  # rename EKUs to EnhancedKeyUsage
  - where-model: X509CertificateProperties
    transform: $.properties.ekus["x-ms-client-name"] = "EnhancedKeyUsage"

  # capitalize acronyms
  - where-model: Certificate
    transform: $.properties.cer["x-ms-client-name"] = "CER"
  - where-model: Certificate
    transform: $.properties.kid["x-ms-client-name"] = "KID"
  - where-model: Certificate
    transform: $.properties.sid["x-ms-client-name"] = "SID"
  - where-model: CertificateOperation
    transform: $.properties.csr["x-ms-client-name"] = "CSR"

  # Remove MaxResults parameter
  - where: "$.paths..*"
    remove-parameter:
      in: query
      name: maxresults

  # remove JSONWeb prefix
  - from: 
      - models.go
      - constants.go
    where: $
    transform: return $.replace(/JSONWebKeyCurveName/g, "CurveName");
  - from: 
      - models.go
      - constants.go
    where: $
    transform: return $.replace(/JSONWebKeyType/g, "KeyType");

  # remove DeletionRecoveryLevel type
  - from: models.go
    where: $
    transform: return $.replace(/RecoveryLevel \*DeletionRecoveryLevel/g, "RecoveryLevel *string");
  - from: constants.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+type DeletionRecoveryLevel string/, "");
  - from: constants.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+func PossibleDeletionRecovery(?:.+\s)+\}/, "");
  - from: constants.go
    where: $
    transform: return $.replace(/const \(\n\s\/\/ DeletionRecoveryLevel(?:.+\s)+\)/, "");

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

  # delete the Attributes model defined in common.json (it's used only with allOf)
  - from: models.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+type Attributes.+\{(?:\s.+\s)+\}\s/, "");
  - from: models_serde.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+func \(a \*?Attributes\).*\{\s(?:.+\s)+\}\s/g, "");

  # delete the version path param check (version == "" is legal for Key Vault but indescribable by OpenAPI)
  - from: client.go
    where: $
    transform: return $.replace(/\sif certificateVersion == "" \{\s+.+certificateVersion cannot be empty"\)\s+\}\s/g, "");

  # delete client name prefix from method options and response types
  - from:
      - client.go
      - models.go
      - options.go
      - response_types.go
    where: $
    transform: return $.replace(/Client(\w+)((?:Options|Response))/g, "$1$2");

  # make cert IDs a convenience type so we can add parsing methods
  # (specifying models because others have "ID" fields whose values aren't key vault object identifiers)
  - from: models.go
    where: $
    transform: return $.replace(/(type (?:Deleted)?Certificate(?:Properties|Policy|Operation)? struct \{(?:\s.+\s)+\sID \*)string/g, "$1ID")
  - from: models.go
    where: $
    transform: return $.replace(/(type (?:Deleted)?Certificate struct \{(?:\s.+\s)+\sKID \*)string/g, "$1ID")
  - from: models.go
    where: $
    transform: return $.replace(/(type (?:Deleted)?Certificate struct \{(?:\s.+\s)+\sSID \*)string/g, "$1ID")

  # remove "certificate" prefix from some method parameter names
  - from: client.go
    where: $
    transform: return $.replace(/certificate((?:Name|Policy|Version)) string/g, (match) => { return match[0].toLowerCase() + match.substr(1); })

  # add doc comment
  - from: swagger-document
    where: $.definitions.X509CertificateProperties.properties.key_usage.items
    transform: $["description"] = "Defines how the certificate's key may be used."

  # remove redundant section of doc comment
  - from:
      - models.go
      - constants.go
    where: $
    transform: return $.replace(/For valid values, see JsonWebKeyCurveName./g, "");
```
