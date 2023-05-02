## Go

```yaml
title: MonitorIngestionClient
description: Azure Monitor Ingestion Go Client
generated-metadata: false

clear-output-folder: false
export-clients: true
go: true
input-file: https://github.com/Azure/azure-rest-api-specs/blob/f07297ce913bfc911470a86436e73c9aceec0587/specification/monitor/data-plane/ingestion/stable/2023-01-01/DataCollectionRules.json
license-header: MICROSOFT_MIT_NO_VERSION
module: github.com/Azure/azure-sdk-for-go/sdk/monitor/azingestion
openapi-type: "data-plane"
output-folder: ../azingestion
override-client-name: Client
security: "AADToken"
use: "@autorest/go@4.0.0-preview.46"
version: "^3.0.0"

directive:
  # delete unused model
  - remove-model: PendingCertificateSigningRequestResult

 # delete unused error models
  - from: models.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+type (?:ErrorResponse|ErrorDetail|ErrorAdditionalInfo).+\{(?:\s.+\s)+\}\s/g, "");

```