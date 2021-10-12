### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/23b8c3e5ecc0a90bc89f93517d7f45ca0b6881d5/specification/sql/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/23b8c3e5ecc0a90bc89f93517d7f45ca0b6881d5/specification/sql/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 0.1.0
modelerfour:
  seal-single-value-enum-by-default: true
```