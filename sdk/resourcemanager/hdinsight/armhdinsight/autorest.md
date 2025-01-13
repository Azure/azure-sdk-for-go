### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/aec83a5f0ed56da4fd16fa027b9fa27edfa8988b/specification/hdinsight/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/aec83a5f0ed56da4fd16fa027b9fa27edfa8988b/specification/hdinsight/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 1.3.0-beta.2
directive:
- from: cluster.json
  where: $.definitions.Resource
  transform: >
    $["title"] = "Resource"
tag: package-2024-08-preview
```