### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/69eacf00a36d565d3220d5dd6f4a5293664f1ae9/specification/keyvault/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/69eacf00a36d565d3220d5dd6f4a5293664f1ae9/specification/keyvault/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 0.2.1
modelerfour:
  seal-single-value-enum-by-default: true
```