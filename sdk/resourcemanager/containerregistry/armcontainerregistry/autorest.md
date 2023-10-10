### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/68d03f91ea7c30e1ab28fb9d35c13f81bc85b724/specification/containerregistry/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/68d03f91ea7c30e1ab28fb9d35c13f81bc85b724/specification/containerregistry/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 1.2.0-beta.2
azcore-version: 1.9.0-beta.1
generate-fakes: true
inject-spans: true
tag: package-2023-08-preview
```
