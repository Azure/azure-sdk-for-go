### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/d830271b241897bad300f9275ff7b27d502aa6c5/specification/storage/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/d830271b241897bad300f9275ff7b27d502aa6c5/specification/storage/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 0.4.0
modelerfour:
  seal-single-value-enum-by-default: true
```