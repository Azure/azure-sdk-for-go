## Go

``` yaml
clear-output-folder: false
export-clients: true
go: true
input-file: https://github.com/Azure/azure-rest-api-specs/blob/dba6ed1f03bda88ac6884c0a883246446cc72495/specification/operationalinsights/data-plane/Microsoft.OperationalInsights/preview/2021-05-19_Preview/OperationalInsights.json
license-header: MICROSOFT_MIT_NO_VERSION
module: github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery
openapi-type: "data-plane"
output-folder: ../azquery
override-client-name: Client
security: "AADToken"
security-scopes:  "https://api.loganalytics.io/.default"
use: "@autorest/go@4.0.0-preview.43"
version: "^3.0.0"

directive:
  # delete metadata endpoint
  - from: swagger-document
    where: $["paths"]
    transform: >
        delete $["/workspaces/{workspaceId}/metadata"];

  # delete metadata operations
  - remove-operation: "Metadata_Post"
  - remove-operation: "Metadata_Get"
  - remove-operation: "Query_Get"

  # delete metadata models
  - remove-model: metadataResults
  - remove-model: metadataCategory
  - remove-model: metadataSolution
  - remove-model: metadataResourceType
  - remove-model: metadataTable
  - remove-model: metadataFunction
  - remove-model: metadataQuery
  - remove-model: metadataApplication
  - remove-model: metadataWorkspace
  - remove-model: metadataResource
  - remove-model: metadataPermissions

  # delete generated constructor
  - from: client.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+func NewClient.+\{\s(?:.+\s)+\}\s/, "");