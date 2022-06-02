### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
tag: package-2020-10-01
require:
- https://github.com/Azure/azure-rest-api-specs/blob/08894fa8d66cb44dc62a73f7a09530f905985fa3/specification/authorization/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/08894fa8d66cb44dc62a73f7a09530f905985fa3/specification/authorization/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 1.0.0
modelerfour:
  lenient-model-deduplication: true
```