### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/64ffad1a3042017e07f8a47df17d6acaa2c1e609/specification/containerservice/resource-manager/Microsoft.ContainerService/aks/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/64ffad1a3042017e07f8a47df17d6acaa2c1e609/specification/containerservice/resource-manager/Microsoft.ContainerService/aks/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 4.2.0-beta.2
tag: package-preview-2023-06
azcore-version: 1.8.0-beta.1
generate-fakes: true
inject-spans: true
```
