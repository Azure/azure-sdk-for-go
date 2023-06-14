### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/f36175f4c54eeec5b6d409406e131dadb540546a/specification/containerservice/resource-manager/Microsoft.ContainerService/aks/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/f36175f4c54eeec5b6d409406e131dadb540546a/specification/containerservice/resource-manager/Microsoft.ContainerService/aks/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module: github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v4
module-version: 4.1.0-beta.3
tag: package-preview-2023-05
azcore-version: 1.7.0-beta.2
generate-fakes: true
inject-spans: true
```
