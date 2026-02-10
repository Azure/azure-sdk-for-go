### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/6312b1c8676b0973f86f078c1177dcb7510158df/specification/hdinsight/resource-manager/Microsoft.HDInsight/HDInsight/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/6312b1c8676b0973f86f078c1177dcb7510158df/specification/hdinsight/resource-manager/Microsoft.HDInsight/HDInsight/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 1.3.0-beta.4
modelerfour:
  lenient-model-deduplication: true
directive:
- from: cluster.json
  where: $.definitions.Resource
  transform: >
    $["title"] = "Resource"
tag: package-2025-01-preview
```