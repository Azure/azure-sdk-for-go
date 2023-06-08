### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/5d2adf9b7fda669b4a2538c65e937ee74fe3f966/specification/compute/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/5d2adf9b7fda669b4a2538c65e937ee74fe3f966/specification/compute/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 5.1.0-beta.1
tag: package-2023-03-01
azcore-version: 1.7.0-beta.2
generate-fakes: true
inject-spans: true
```
