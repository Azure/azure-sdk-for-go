### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
tag: package-combine-2026-04
require:
- https://github.com/Azure/azure-rest-api-specs/blob/3c8605ab5c54e92d1a0cbf180d97c5d7ab99df60/specification/security/resource-manager/Microsoft.Security/Security/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/3c8605ab5c54e92d1a0cbf180d97c5d7ab99df60/specification/security/resource-manager/Microsoft.Security/Security/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 0.15.0
directive:
- from: externalSecuritySolutions.json
  where: $.definitions['ExternalSecuritySolutionKind']
  transform: >
      $ = {
        "type": "string",
        "description": "The kind of the external solution",
        "enum": [
          "CEF",
          "ATA",
          "AAD"
        ],
        "x-ms-enum": {
          "name": "ExternalSecuritySolutionKind",
          "modelAsString": true,
          "values": [
            {
              "value": "CEF"
            },
            {
              "value": "ATA"
            },
            {
              "value": "AAD"
            }
          ]
        }
      };
- from: externalSecuritySolutions.json
  where: $.definitions['ExternalSecuritySolution']
  transform: >
      $.properties['kind'] = {
        "$ref": "#/definitions/ExternalSecuritySolutionKind"
      };
      $.allOf = [
        {
          "$ref": "../../../common/v1/types.json#/definitions/Resource"
        },
        {
          "$ref": "../../../common/v1/types.json#/definitions/Location"
        }
      ]
- rename-model:
    from: SecurityStandard
    to: ArmSecurityStandard
- rename-model:
    from: SecurityStandards
    to: ArmSecurityStandards
- rename-model:
    from: SecurityStandardList
    to: ArmSecurityStandardList
- rename-model:
    from: SecurityStandardProperties
    to: ArmSecurityStandardProperties
- from: swagger-document
  where: '$.paths.*[?(@.operationId.startsWith("SecurityStandards_"))]'
  transform: >
    $["operationId"] = $["operationId"].replace("SecurityStandards_", "ArmSecurityStandards_");
- from: swagger-document
  where: '$.paths.*[?(@.operationId.startsWith("Standards_"))]'
  transform: >
    $["operationId"] = $["operationId"].replace("Standards_", "ArmStandards_");
```