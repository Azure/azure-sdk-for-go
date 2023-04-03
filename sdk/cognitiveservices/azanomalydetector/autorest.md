## Go

These settings apply only when `--go` is specified on the command line.

``` yaml
input-file:
- https://github.com/azure/azure-rest-api-specs/blob/main/specification/cognitiveservices/data-plane/AnomalyDetector/stable/v1.1/openapi.json
output-folder: ../azanomalydetector
clear-output-folder: false
module: github.com/Azure/azure-sdk-for-go/sdk/cognitiveservices/azanomalydetector
license-header: MICROSOFT_MIT_NO_VERSION
openapi-type: data-plane
go: true
use: "@autorest/go@4.0.0-preview.47"
```

### Temporary transformations

``` yaml
directive:

  # Add x-ms-parameter-location to parameters in x-ms-parameterized-host
  - from: swagger-document
    where: $.x-ms-parameterized-host.parameters.0
    transform: $["x-ms-parameter-location"] = "client"
  - from: swagger-document
    where: $.x-ms-parameterized-host.parameters.1
    transform: $["x-ms-parameter-location"] = "client"

  # delete APIVersion schema -- it conflicts with server parameter
  - remove-model: APIVersion

  # Rename body parameters to body
  - from: swagger-document
    where: $.paths..parameters..[?(@.name=='options' && @.in=='body')]
    transform: $["name"] = "body"

```
