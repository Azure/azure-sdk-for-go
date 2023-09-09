### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/bbe1ea8bf5aa6cfbfa8855e03dbb9a93f8266bcd/specification/containerregistry/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/bbe1ea8bf5aa6cfbfa8855e03dbb9a93f8266bcd/specification/containerregistry/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module:  github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry
module-version: 1.1.0-beta.4
azcore-version: 1.8.0-beta.1
generate-fakes: true
inject-spans: true
tag: package-2023-07
```
