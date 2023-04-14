### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
tag: package-composite-v3
require:
- https://github.com/Azure/azure-rest-api-specs/blob/af3f7994582c0cbd61a48b636907ad2ac95d332c/specification/security/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/af3f7994582c0cbd61a48b636907ad2ac95d332c/specification/security/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 0.10.1
modelerfour:
  lenient-model-deduplication: true
```