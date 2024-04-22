### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
tag: package-composite-v3
require:
- https://github.com/Azure/azure-rest-api-specs/blob/b43974e07d3204c4b6f8396627f5430994a7f7c9/specification/security/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/b43974e07d3204c4b6f8396627f5430994a7f7c9/specification/security/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 0.13.0
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
```