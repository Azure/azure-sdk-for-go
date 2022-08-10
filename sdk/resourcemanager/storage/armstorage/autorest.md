### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/6555955d8dac9d3a91ff5eb1740b5af4c7294307/specification/storage/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/6555955d8dac9d3a91ff5eb1740b5af4c7294307/specification/storage/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 2.0.0
modelerfour:
  seal-single-value-enum-by-default: true
```