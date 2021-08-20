### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
require:
- https://github.com/Azure/azure-rest-api-specs/blob/406c7b7d4633491f3b4cdb11e91bbe1045068dce/specification/cdn/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/406c7b7d4633491f3b4cdb11e91bbe1045068dce/specification/cdn/resource-manager/readme.go.md
module-version: 0.1.0
directive:
- from: cdn.json
  where: $.definitions.DeliveryRuleAction
  transform: >
    $["x-ms-client-name"] = "DeliveryRuleActionType"
```