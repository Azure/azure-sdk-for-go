### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
require:
- https://github.com/Azure/azure-rest-api-specs/blob/5d09c12c024fa7efbaca6a95b9741a46a886fe6f/specification/cdn/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/5d09c12c024fa7efbaca6a95b9741a46a886fe6f/specification/cdn/resource-manager/readme.go.md
module-version: 0.1.0
directive:
- from: cdn.json
  where: $.definitions.DeliveryRuleAction
  transform: >
    $["x-ms-client-name"] = "DeliveryRuleActionType"
```