### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
tag: package-composite-v3
require:
- https://github.com/Azure/azure-rest-api-specs/blob/fc3c409825d6a31ba909a88ea34dcc13165edc9c/specification/security/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/fc3c409825d6a31ba909a88ea34dcc13165edc9c/specification/security/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 0.9.0
modelerfour:
  lenient-model-deduplication: true
```