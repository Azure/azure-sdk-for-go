### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/b32e1896f30e6ea155449cb49719a6286e32b961/specification/storage/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/b32e1896f30e6ea155449cb49719a6286e32b961/specification/storage/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module: github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage
module-version: 1.4.0-beta.1
azcore-version: 1.7.0-beta.2
generate-fakes: true
inject-spans: true
modelerfour:
  seal-single-value-enum-by-default: true
```
