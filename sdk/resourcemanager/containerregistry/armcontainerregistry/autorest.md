### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/969fd0c2634fbcc1975d7abe3749330a5145a97c/specification/containerregistry/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/969fd0c2634fbcc1975d7abe3749330a5145a97c/specification/containerregistry/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module:  github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry
module-version: 1.1.0-beta.3
azcore-version: 1.8.0-beta.1
generate-fakes: true
inject-spans: true
```
