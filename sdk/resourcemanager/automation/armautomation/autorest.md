### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/432872fac1d0f8edcae98a0e8504afc0ee302710/specification/automation/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/432872fac1d0f8edcae98a0e8504afc0ee302710/specification/automation/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 0.9.0
directive:
  - where-operation: DscConfiguration_CreateOrUpdate
    transform: delete $['x-ms-examples']
  - where-operation: DscConfiguration_Update
    transform: delete $['x-ms-examples']
```