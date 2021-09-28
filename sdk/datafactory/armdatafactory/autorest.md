### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
require:
- https://github.com/Azure/azure-rest-api-specs/blob/4c8162b0a1f7bbd46e9aedc0e19bbe181e549c4c/specification/datafactory/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/4c8162b0a1f7bbd46e9aedc0e19bbe181e549c4c/specification/datafactory/resource-manager/readme.go.md
module-version: 0.1.0
directive:
 - from: DataFlow.json
   where: $.definitions.DataFlow
   transform: >
     $["required"] = ["type"]
 - from: ManagedPrivateEndpoint.json
   where: $.definitions.ManagedPrivateEndpoint
   transform: >
     $["required"] = ["type"]
 - from: ManagedVirtualNetwork.json
   where: $.definitions.ManagedVirtualNetwork
   transform: >
     $["required"] = ["type"]
 - from: Trigger.json
   where: $.definitions.BlobEventTypes
   transform: >
     $["x-ms-client-name"] = "BlobEventType"
```