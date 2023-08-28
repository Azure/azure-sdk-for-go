### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/c53808ba54beef57059371708f1fa6949a11a280/specification/compute/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/c53808ba54beef57059371708f1fa6949a11a280/specification/compute/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module: github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5
module-version: 5.2.0-beta.1
tag: package-2023-01-02
azcore-version: 1.8.0-beta.1
generate-fakes: true
inject-spans: true
```
