### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/89260be1a92c914b7b48af8e8f75938d5e76851d/specification/containerservice/resource-manager/Microsoft.ContainerService/aks/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/89260be1a92c914b7b48af8e8f75938d5e76851d/specification/containerservice/resource-manager/Microsoft.ContainerService/aks/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 4.5.0-beta.1
azcore-version: 1.9.0-beta.1
generate-fakes: true
inject-spans: true
tag: package-preview-2023-08
```
