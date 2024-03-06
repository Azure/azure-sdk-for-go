## Go

```yaml
clear-output-folder: false
export-clients: true
go: true
input-file: https://github.com/Azure/azure-rest-api-specs/blob/7452e1cc7db72fbc6cd9539b390d8b8e5c2a1864/specification/keyvault/data-plane/Microsoft.KeyVault/stable/7.5/rbac.json
license-header: MICROSOFT_MIT_NO_VERSION
openapi-type: "data-plane"
output-folder: ../rbac
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

    # rename role definition and role assignment operations so they will generate as one access control client
  - rename-operation:
      from: RoleDefinitions_Delete
      to: DeleteRoleDefinition
  - rename-operation:
      from: RoleAssignments_Delete
      to: DeleteRoleAssignment
  - rename-operation:
      from: RoleDefinitions_CreateOrUpdate
      to: CreateOrUpdateRoleDefinition
  - rename-operation:
      from: RoleAssignments_Create
      to: CreateRoleAssignment
  - rename-operation:
      from: RoleDefinitions_Get
      to: GetRoleDefinition
  - rename-operation:
      from: RoleAssignments_Get
      to: GetRoleAssignment
  - rename-operation:
      from: RoleDefinitions_List
      to: ListRoleDefinitions
  - rename-operation:
      from: RoleAssignments_ListForScope
      to: ListRoleAssignments

  # delete unused error models
  - from: models.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+type (?:Error|KeyVaultError).+\{(?:\s.+\s)+\}\s/g, "");
  - from: models_serde.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+func \(\w \*?(?:Error|KeyVaultError)\).*\{\s(?:.+\s)+\}\s/g, "");

  # delete unused filter models
  - from: models.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+type (?:RoleAssignmentFilter|RoleDefinitionFilter).+\{(?:\s.+\s)+\}\s/g, "");
  - from: models_serde.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+func \(\w \*?(?:RoleAssignmentFilter|RoleDefinitionFilter)\).*\{\s(?:.+\s)+\}\s/g, "");

  # change type of scope parameter from string to RoleScope
  - from: client.go
    where: $
    transform:  return $.replace(/scope string/g, "scope RoleScope");
  - from: client.go
    where: $
    transform:  return $.replace(/scope\)/g, "string(scope))");
  
  # delete client name prefix from method options and response types
  - from:
      - client.go
      - models.go
      - options.go
      - response_types.go
      - options.go
    where: $
    transform: return $.replace(/Client(\w+)((?:Options|Response))/g, "$1$2");

  # fix up span names
  - from: client.go
    where: $
    transform: return $.replace(/StartSpan\(ctx, "Client/g, "StartSpan(ctx, \"rbac.Client");
```