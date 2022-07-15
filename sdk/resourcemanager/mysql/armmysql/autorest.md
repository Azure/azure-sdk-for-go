### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/c03c258c7a01a7d57b3110cc20e2e76752b6f2d6/specification/mysql/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/c03c258c7a01a7d57b3110cc20e2e76752b6f2d6/specification/mysql/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 1.0.1
package-singleservers: true
directive:
- from: Servers.json
  where: $.definitions.CloudError.properties.error
  transform: >
    $["description"] = undefined
```