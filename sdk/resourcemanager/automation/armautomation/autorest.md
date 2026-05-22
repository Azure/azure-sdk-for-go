### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/0f9539236cbea0cd9ca5dc0bde00d15a039fe22d/specification/automation/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/0f9539236cbea0cd9ca5dc0bde00d15a039fe22d/specification/automation/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 1.0.0
tag: package-2024-10-23
directive:
  - where-operation: DscConfiguration_CreateOrUpdate
    transform: delete $['x-ms-examples']
  - where-operation: DscConfiguration_Update
    transform: delete $['x-ms-examples']
```