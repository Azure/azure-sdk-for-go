### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/124789ad0942fcafded1c1dbd6d2a703b23d10c7/specification/storage/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/124789ad0942fcafded1c1dbd6d2a703b23d10c7/specification/storage/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 1.1.0
modelerfour:
  seal-single-value-enum-by-default: true
```