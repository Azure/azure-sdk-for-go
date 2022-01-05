### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- /home/vsts/work/1/s/azure-rest-api-specs/specification/hdinsight/resource-manager/readme.md
- /home/vsts/work/1/s/azure-rest-api-specs/specification/hdinsight/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 0.1.1
directive:
- from: cluster.json
  where: $.definitions.Resource
  transform: >
    $["title"] = "Resource"
```