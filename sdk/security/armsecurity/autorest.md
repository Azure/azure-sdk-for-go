### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
require:
- https://github.com/Azure/azure-rest-api-specs/blob/dcf9fa24061fb4ac71fec8d054fb4f90e988925e/specification/security/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/dcf9fa24061fb4ac71fec8d054fb4f90e988925e/specification/security/resource-manager/readme.go.md
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