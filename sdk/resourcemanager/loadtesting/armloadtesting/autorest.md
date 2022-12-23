### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/65f5099082a02703fbca8c9863940ea83cdf0bb7/specification/loadtestservice/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/65f5099082a02703fbca8c9863940ea83cdf0bb7/specification/loadtestservice/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 1.0.2

# v1.0.0 and v1.0.1 has been retracted because of mistake. When this RP goes GA, it starts at version v1.0.2
```