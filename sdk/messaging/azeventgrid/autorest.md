## Go

``` yaml
title: EventGridClient
description: Azure Event Grid client
generated-metadata: false
clear-output-folder: false
export-clients: true
go: true
input-file: 
    # This was the commit that everyone used to generate their first official betas.
    - https://raw.githubusercontent.com/Azure/azure-rest-api-specs/947c9ce9b20900c6cbc8e95bc083e723d09a9c2c/specification/eventgrid/data-plane/Microsoft.EventGrid/preview/2023-06-01-preview/EventGrid.json
    # when we start using the .tsp file directly we can start referring to the compiled output.
    # ./tsp-output\@azure-tools\typespec-autorest\2023-06-01-preview\openapi.json
license-header: MICROSOFT_MIT_NO_VERSION
module: github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventgrid
openapi-type: "data-plane"
output-folder: ../azeventgrid
override-client-name: Client
security: "AADToken"
use: "@autorest/go@4.0.0-preview.46"
version: "^3.0.0"
directive:
  # we have to write a little wrapper code for this so we'll hide the public function
  # for now.
  - from: client.go
    where: $
    transform: return $.replace(/PublishCloudEvents\(/g, "internalPublishCloudEvents(");
  # make sure the casing of the properties is what compliant.
  - from: swagger-document
    where: $.definitions.CloudEvent.properties.specversion
    transform: $["x-ms-client-name"] = "SpecVersion"
  - from: swagger-document
    where: $.definitions.CloudEvent.properties.datacontenttype
    transform: $["x-ms-client-name"] = "DataContentType"
  - from: swagger-document
    where: $.definitions.CloudEvent.properties.dataschema
    transform: $["x-ms-client-name"] = "DataSchema"
  # make the endpoint a parameter of the client constructor
  - from: swagger-document
    where: $["x-ms-parameterized-host"]
    transform: $.parameters[0]["x-ms-parameter-location"] = "client"
  # delete client name prefix from method options and response types
  - from:
      - client.go
      - models.go
      - response_types.go
    where: $
    transform: return $.replace(/Client(\w+)((?:Options|Response))/g, "$1$2");
```
