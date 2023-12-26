### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/310a0100f5b020c1900c527a6aa70d21992f078a/specification/hdinsight/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/310a0100f5b020c1900c527a6aa70d21992f078a/specification/hdinsight/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 1.3.0-beta.1
directive:
- from: cluster.json
  where: $.definitions.Resource
  transform: >
    $["title"] = "Resource"
tag: package-2023-04-preview
```