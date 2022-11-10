### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/68847d6ae901f0cb2efa62ae2c523ad8cf5c2ea3/specification/monitor/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/68847d6ae901f0cb2efa62ae2c523ad8cf5c2ea3/specification/monitor/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 0.9.0
directive:
  - from: swagger-document
    where: $["paths"]["/{resourceUri}/providers/Microsoft.Insights/metricDefinitions"]
    transform: >
      delete $.get
  - from: swagger-document
    where: $["paths"]["/{resourceUri}/providers/Microsoft.Insights/metrics"]
    transform: >
      delete $.get
  - from: swagger-document
    where: $["paths"]["/{resourceUri}/providers/microsoft.insights/metricNamespaces"]
    transform: >
      delete $.get
```
