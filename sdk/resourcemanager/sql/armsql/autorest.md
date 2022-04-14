### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/b9d36b704e582a2bd5677fedc813607e73963469/specification/sql/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/b9d36b704e582a2bd5677fedc813607e73963469/specification/sql/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 0.4.0
modelerfour:
  seal-single-value-enum-by-default: true
```