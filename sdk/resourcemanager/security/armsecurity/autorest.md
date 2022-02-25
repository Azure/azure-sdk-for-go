### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
tag: package-composite-v3
require:
- https://github.com/Azure/azure-rest-api-specs/blob/4442d8de32ba14a9126cdadd8538da9d4539ff5e/specification/security/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/4442d8de32ba14a9126cdadd8538da9d4539ff5e/specification/security/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 0.4.0
modelerfour:
  lenient-model-deduplication: true
```