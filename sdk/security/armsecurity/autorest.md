### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
require:
- https://github.com/Azure/azure-rest-api-specs/blob/4c8162b0a1f7bbd46e9aedc0e19bbe181e549c4c/specification/security/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/4c8162b0a1f7bbd46e9aedc0e19bbe181e549c4c/specification/security/resource-manager/readme.go.md
module-version: 0.1.0
directive:
  - from: externalSecuritySolutions.json
    where: $.definitions.AadConnectivityState
    transform: >
      $["x-ms-client-name"] = "AadConnectivityStateDummy"
  - from: externalSecuritySolutions.json
    where: $.definitions.ExternalSecuritySolutionKind
    transform: >
      $["x-ms-client-name"] = "ExternalSecuritySolutionKindDummy"
```