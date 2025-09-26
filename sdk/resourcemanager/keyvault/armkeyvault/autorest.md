### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/8a9dbb28e788355a47dc5bad3ea5f8da212b4bf6/specification/keyvault/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/8a9dbb28e788355a47dc5bad3ea5f8da212b4bf6/specification/keyvault/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 2.0.0
modelerfour:
  seal-single-value-enum-by-default: true
tag: package-2025-05-01
```