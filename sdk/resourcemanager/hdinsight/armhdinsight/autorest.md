### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/c0f5f5f439ce6152ff3c078f9ba02f2549b2b58c/specification/hdinsight/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/c0f5f5f439ce6152ff3c078f9ba02f2549b2b58c/specification/hdinsight/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 0.1.0
directive:
- from: cluster.json
  where: $.definitions.Resource
  transform: >
    $["title"] = "Resource"
```