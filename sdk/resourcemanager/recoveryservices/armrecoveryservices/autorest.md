### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/69eacf00a36d565d3220d5dd6f4a5293664f1ae9/specification/recoveryservices/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/69eacf00a36d565d3220d5dd6f4a5293664f1ae9/specification/recoveryservices/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 0.2.0
directive:
- from: vaults.json
  where: '$.paths["/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}/operationStatus/{operationId}"].get'
  transform: >
    $["operationId"] = "Operations_OperationStatus_Get"
- from: vaults.json
  where: '$.paths["/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}/operationResults/{operationId}"].get'
  transform: >
    $["operationId"] = "Operations_GetOperationResult"
```