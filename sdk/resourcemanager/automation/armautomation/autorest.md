### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/947878feb7a528ff4344f8fe67ab85522afcbb79/specification/automation/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/947878feb7a528ff4344f8fe67ab85522afcbb79/specification/automation/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
tag: package-2024-10-23
module-version: 1.0.0
directive:
  - where-operation: DscConfiguration_CreateOrUpdate
    transform: delete $['x-ms-examples']
  - where-operation: DscConfiguration_Update
    transform: delete $['x-ms-examples']
```