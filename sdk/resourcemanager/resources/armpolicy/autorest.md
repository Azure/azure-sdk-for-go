### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/acc18877f6e4e6050d23d03e8216c2f43baeaee2/specification/resources/resource-manager/Microsoft.Authorization/policy/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/acc18877f6e4e6050d23d03e8216c2f43baeaee2/specification/resources/resource-manager/Microsoft.Authorization/policy/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 0.11.0
package-policy: true
tag: package-policy-2025-03-go
modelerfour:
  lenient-model-deduplication: true
```