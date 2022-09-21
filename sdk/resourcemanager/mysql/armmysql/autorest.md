### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- /mnt/vss/_work/1/s/azure-rest-api-specs/specification/mysql/resource-manager/readme.md
- /mnt/vss/_work/1/s/azure-rest-api-specs/specification/mysql/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 2.0.0
package-singleservers: true
directive:
- from: Servers.json
  where: $.definitions.CloudError.properties.error
  transform: >
    $["description"] = undefined
```