## Go

```yaml
clear-output-folder: false
export-clients: true
go: true
input-file: 
    - https://github.com/Azure/azure-rest-api-specs/blob/main/specification/keyvault/data-plane/Microsoft.KeyVault/preview/7.4-preview.1/rbac.json
    #- https://github.com/Azure/azure-rest-api-specs/blob/main/specification/keyvault/data-plane/Microsoft.KeyVault/preview/7.4-preview.1/backuprestore.json
    #- https://github.com/Azure/azure-rest-api-specs/blob/main/specification/keyvault/data-plane/Microsoft.KeyVault/preview/7.4-preview.1/settings.json
license-header: MICROSOFT_MIT_NO_VERSION
module: github.com/Azure/azure-sdk-for-go/sdk/keyvault/azadmin
openapi-type: "data-plane"
output-folder: ../azadmin
override-client-name: Client
security: "AADToken"
security-scopes: "https://vault.azure.net/.default"
use: "@autorest/go@4.0.0-preview.44"
version: "^3.0.0"

#directive:
  #- rename-operation:
      #from: RoleDefinitions_Delete
      #to: AccessControl_DeleteRoleDefinitions
  #- rename-operation:
      #from: RoleAssignments_Delete
      #to: AccessControl_DeleteRoleAssignments
```