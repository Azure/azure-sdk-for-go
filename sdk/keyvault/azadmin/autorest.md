## Go

```yaml
clear-output-folder: false
export-clients: true
go: true
input-file: 
    - https://github.com/Azure/azure-rest-api-specs/blob/main/specification/keyvault/data-plane/Microsoft.KeyVault/preview/7.4-preview.1/rbac.json
    - https://github.com/Azure/azure-rest-api-specs/blob/main/specification/keyvault/data-plane/Microsoft.KeyVault/preview/7.4-preview.1/backuprestore.json
    - https://github.com/Azure/azure-rest-api-specs/blob/main/specification/keyvault/data-plane/Microsoft.KeyVault/preview/7.4-preview.1/settings.json
license-header: MICROSOFT_MIT_NO_VERSION
module: github.com/Azure/azure-sdk-for-go/sdk/keyvault/azadmin
openapi-type: "data-plane"
output-folder: ../azadmin
override-client-name: BackupClient
security: "AADToken"
security-scopes: "https://vault.azure.net/.default"
use: "@autorest/go@4.0.0-preview.44"
version: "^3.0.0"

directive:

  # make vault URL a parameter of the client constructor
  - from: swagger-document
    where: $["x-ms-parameterized-host"]
    transform: $.parameters[0]["x-ms-parameter-location"] = "client"

  # rename role definition and role assignment operations so they will generate as one access control client
  - rename-operation:
      from: RoleDefinitions_Delete
      to: AccessControl_DeleteRoleDefinition
  - rename-operation:
      from: RoleAssignments_Delete
      to: AccessControl_DeleteRoleAssignment
  - rename-operation:
      from: RoleDefinitions_CreateOrUpdate
      to: AccessControl_CreateOrUpdateRoleDefinition
  - rename-operation:
      from: RoleAssignments_Create
      to: AccessControl_CreateRoleAssignment
  - rename-operation:
      from: RoleDefinitions_Get
      to: AccessControl_GetRoleDefinition
  - rename-operation:
      from: RoleAssignments_Get
      to: AccessControl_GetRoleAssignment
  - rename-operation:
      from: RoleDefinitions_List
      to: AccessControl_ListRoleDefinitions
  - rename-operation:
      from: RoleAssignments_ListForScope
      to: AccessControl_ListRoleAssignments

    # rename setting operations to generate as their own client
  - rename-operation:
      from: GetSetting
      to: Settings_GetSetting
  - rename-operation:
      from: GetSettings
      to: Settings_GetSettings
  - rename-operation:
      from: UpdateSetting
      to: Settings_UpdateSetting

  # delete generated client constructor
  - from: accesscontrol_client.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+func NewAccessControlClient.+\{\s(?:.+\s)+\}\s/, "");
  - from: backup_client.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+func NewBackupClient.+\{\s(?:.+\s)+\}\s/, "");
  - from: settings_client.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+func NewSettingsClient.+\{\s(?:.+\s)+\}\s/, "");

 
  # change type of scope parameter from string to RoleScope
  - from: accesscontrol_client.go
    where: $
    transform:  return $.replace(/scope string/g, "scope RoleScope");
  - from: accesscontrol_client.go
    where: $
    transform:  return $.replace(/scope\)/g, "string(scope))");

```