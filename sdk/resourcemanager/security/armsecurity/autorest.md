### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
tag: package-composite-v3
require:
- https://github.com/Azure/azure-rest-api-specs/blob/24166cedd1055380f1b9d40df270bf51b287d7d9/specification/security/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/24166cedd1055380f1b9d40df270bf51b287d7d9/specification/security/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 0.10.0
modelerfour:
  lenient-model-deduplication: true
```