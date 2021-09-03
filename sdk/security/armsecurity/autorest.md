### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
require:
- https://github.com/Azure/azure-rest-api-specs/blob/87a56cc36600486d4ca312ecfbe09bf9b278fee4/specification/security/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/87a56cc36600486d4ca312ecfbe09bf9b278fee4/specification/security/resource-manager/readme.go.md
module-version: 0.1.0
directive:
  - from: alerts.json
    where: $.definitions.AlertSimulatorRequestProperties.properties.kind["x-ms-enum"].name
    transform: return "AlertKind"
  - from: externalSecuritySolutions.json
    where: $.definitions.AadConnectivityState
    transform: >
      $["x-ms-client-name"] = "AadConnectivityStateDummy"
  - from: externalSecuritySolutions.json
    where: $.definitions.ExternalSecuritySolutionKind
    transform: >
      $["x-ms-client-name"] = "ExternalSecuritySolutionKindDummy"
```