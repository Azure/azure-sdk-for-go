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
```
