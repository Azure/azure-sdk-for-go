## Go

```yaml
title: MonitorIngestionClient
description: Azure Monitor Ingestion Go Client
generated-metadata: false

clear-output-folder: false
go: true
input-file: https://github.com/Azure/azure-rest-api-specs/blob/f07297ce913bfc911470a86436e73c9aceec0587/specification/monitor/data-plane/ingestion/stable/2023-01-01/DataCollectionRules.json
license-header: MICROSOFT_MIT_NO_VERSION
module: github.com/Azure/azure-sdk-for-go/sdk/monitor/azingest
openapi-type: "data-plane"
output-folder: ../azlogs
override-client-name: Client
security: "AADToken"
use: "@autorest/go@4.0.0-preview.61"
inject-spans: true
version: "^3.0.0"
rawjson-as-bytes: true

directive:
  # delete unused model
  - remove-model: PendingCertificateSigningRequestResult

  # remove x-ms-client-request-id, it's added in a pipeline policy
  - where-operation: Upload
    remove-parameter:
      in: header
      name: x-ms-client-request-id

  # rename parameter from "body" to "logs", "stream" to "streamName"
  - from: swagger-document
    where: $.paths..parameters..[?(@.name=='body')]
    transform: $["x-ms-client-name"] = "logs"
  - from: swagger-document
    where: $.paths..parameters..[?(@.name=='stream')]
    transform: $["x-ms-client-name"] = "streamName"

  # delete client name prefix from method options and response types
  - from:
      - client.go
      - options.go
      - response_types.go
    where: $
    transform: return $.replace(/Client(\w+)((?:Options|Response))/g, "$1$2");

  # update doc comments
  - from: swagger-document
    where: $.paths..parameters..[?(@.name=='Content-Encoding')]
    transform: $["description"] = "If the bytes of the \"logs\" parameter are already gzipped, set ContentEncoding to \"gzip\""
  - from: swagger-document
    where: $.paths./dataCollectionRules/{ruleId}/streams/{stream}.post
    transform: $["description"] = "Ingestion API used to directly ingest data using Data Collection Rules. Maximum size of of API call is 1 MB."

  # fix up body param to keep back-compat
  - from: swagger-document
    where: $.paths..parameters..[?(@.name=='body')]
    transform: |
      $['schema'] = { "type": "object" }
```