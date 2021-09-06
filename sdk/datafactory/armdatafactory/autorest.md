### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
require:
- https://github.com/Azure/azure-rest-api-specs/blob/5d09c12c024fa7efbaca6a95b9741a46a886fe6f/specification/datafactory/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/5d09c12c024fa7efbaca6a95b9741a46a886fe6f/specification/datafactory/resource-manager/readme.go.md
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