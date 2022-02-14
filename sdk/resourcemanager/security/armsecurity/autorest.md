### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
tag: package-composite-v3
require:
- https://github.com/Azure/azure-rest-api-specs/blob/cf47fa91b882618a1043e3aeb5803b3a7397cd08/specification/security/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/cf47fa91b882618a1043e3aeb5803b3a7397cd08/specification/security/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 0.3.0
modelerfour:
  lenient-model-deduplication: true
```