# Go

These settings apply only when `--go` is specified on the command line.

``` yaml
input-file:
# PR: https://github.com/Azure/azure-rest-api-specs/pull/27076/files
#- https://raw.githubusercontent.com/Azure/azure-rest-api-specs/18c24352ad4a2e0959c0b4ec1404c3a250912f8b/specification/ai/data-plane/OpenAI.Assistants/OpenApiV2/preview/2024-02-15-preview/assistants_generated.json
- ./testdata/generated/openapi.json
output-folder: ../azopenaiassistants
clear-output-folder: false
module: github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiassistants
license-header: MICROSOFT_MIT_NO_VERSION
openapi-type: data-plane
go: true
title: "OpenAIAssistants"
use: "@autorest/go@4.0.0-preview.52"
slice-elements-byval: true
# can't use this since it removes an innererror type that we want ()
# remove-non-reference-schema: true
```

## Transformations

Fix deployment and endpoint parameters so they show up in the right spots

``` yaml
directive:
  # Add x-ms-parameter-location to parameters in x-ms-parameterized-host
  - from: swagger-document
    where: $["x-ms-parameterized-host"].parameters.0
    transform: $["x-ms-parameter-location"] = "client";
```
