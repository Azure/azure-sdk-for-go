### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
require:
- https://github.com/Azure/azure-rest-api-specs/blob/d5e70e3c12490a8c980b890cb611e85bbbae5858/specification/mysql/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/d5e70e3c12490a8c980b890cb611e85bbbae5858/specification/mysql/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 0.1.0
package-singleservers: true
directive:
- from: Servers.json
  where: $.definitions.CloudError.properties.error
  transform: >
    $["description"] = undefined
```