### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/17aa6a1314de5aafef059d9aa2229901df506e75/specification/network/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/17aa6a1314de5aafef059d9aa2229901df506e75/specification/network/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module: github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v3
module-version: 3.1.0-beta.2
tag: package-2022-11
azcore-version: 1.8.0-beta.1
generate-fakes: true
inject-spans: true
```
