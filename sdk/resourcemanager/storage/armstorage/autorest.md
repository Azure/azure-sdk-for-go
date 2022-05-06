### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Alancere/azure-rest-api-specs/blob/mocktest_R0/specification/storage/resource-manager/readme.md
- https://github.com/Alancere/azure-rest-api-specs/blob/mocktest_R0/specification/storage/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 0.6.0
modelerfour:
  seal-single-value-enum-by-default: true
```