### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
generate-fakes: true
modelerfour:
  lenient-model-deduplication: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/20077475fad69cd39ee2408a9c9835bd36a53be3/specification/compute/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/20077475fad69cd39ee2408a9c9835bd36a53be3/specification/compute/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module: github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4
module-version: 4.0.0
output-folder: $(go-sdk-folder)/sdk/resourcemanager/compute/armcompute
```
