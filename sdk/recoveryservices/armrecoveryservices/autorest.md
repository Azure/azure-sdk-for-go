### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
require:
- https://github.com/Azure/azure-rest-api-specs/blob/4c8162b0a1f7bbd46e9aedc0e19bbe181e549c4c/specification/recoveryservices/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/4c8162b0a1f7bbd46e9aedc0e19bbe181e549c4c/specification/recoveryservices/resource-manager/readme.go.md
module-version: 0.1.0
directive:
- from: vaults.json
  where: '$.paths["/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}/operationStatus/{operationId}"].get'
  transform: >
    $["operationId"] = "RecoveryServices_GetOperationStatus"
- from: vaults.json
  where: '$.paths["/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}/operationResults/{operationId}"].get'
  transform: >
    $["operationId"] = "RecoveryServices_GetOperationResult"
```