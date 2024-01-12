## Go

``` yaml
title: AzureEventGridSystemEvents
description: Azure Event Grid system events
generated-metadata: false
clear-output-folder: false
go: true
require: https://github.com/Azure/azure-rest-api-specs/blob/11bbc2b1df2e915a2227a6a1a48a27b9e67c3311/specification/eventgrid/data-plane/readme.md
license-header: MICROSOFT_MIT_NO_VERSION
openapi-type: "data-plane"
output-folder: ../systemevents
override-client-name: ClientDeleteMe
security: "AADToken"
use: "@autorest/go@4.0.0-preview.52"
version: "^3.0.0"
slice-elements-byval: true
remove-non-reference-schema: true
batch:
  - tag: package-2018-01
directive:
  - from: swagger-document
    where: $
    transform: $['paths'] = {}; return $;
  - from: swagger-document
    where: $.definitions.MediaJobOutput
    transform: >
      $.required.push("@odata.type");
      $["x-csharp-usage"] = "model,output";
  # reference azcore/messaging/CloudEvent
  - from: client.go
    where: $
    transform: return $.replace(/\[\]CloudEvent/g, "[]messaging.CloudEvent");
  - from: client.go
    where: $
    transform: return $.replace(/func \(client \*Client\) PublishCloudEventEvents\(/g, "func (client *Client) internalPublishCloudEventEvents(");  
  - from: swagger-document
    where: $.definitions.CloudEventEvent
    transform: $["x-ms-external"] = true
   # delete some models that look like they're system events...
  - from: models.go
    where: $
    transform: return $.replace(/\/\/ (SubscriptionDeletedEventData|SubscriptionValidationEventData|SubscriptionValidationResponse).+?\n}/gs, "")    
  - from: models_serde.go
    where: $    
    transform: |
      return $
        .replace(/\/\/ MarshalJSON implements the json.Marshaller interface for type (SubscriptionDeletedEventData|SubscriptionValidationEventData|SubscriptionValidationResponse).+?\n}/gs, "")
        .replace(/\/\/ UnmarshalJSON implements the json.Unmarshaller interface for type (SubscriptionDeletedEventData|SubscriptionValidationEventData|SubscriptionValidationResponse).+?\n}/gs, "");
  - from: 
      - models.go
      - client.go
      - response_types.go
      - options.go
    where: $
    transform: return $.replace(/CloudEventEvent/g, "CloudEvent");
  - from: 
      - models.go
      - models_serde.go
      - client.go
      - response_types.go
      - options.go
    where: $
    transform: return $.replace(/EventGridEvent/g, "Event");
  - from: 
      - client.go
    where: $
    transform: | 
      return $.replace(
        /(func \(client \*Client\) publishCloudEventsCreateRequest.+?)return req, nil/s, 
        '$1\nreq.Raw().Header.Set("Content-type", "application/cloudevents-batch+json; charset=utf-8")\nreturn req, nil');
  - from: 
      - models.go
      - models_serde.go
    where: $
    transform: |
      return $.replace(/ResourceActionCancelData/g, 'ResourceActionCancelEventData')
        .replace(/ResourceActionFailureData/g, 'ResourceActionFailureEventData')
        .replace(/ResourceActionSuccessData/g, 'ResourceActionSuccessEventData')
        .replace(/ResourceDeleteCancelData/g, 'ResourceDeleteCancelEventData')
        .replace(/ResourceDeleteFailureData/g, 'ResourceDeleteFailureEventData')
        .replace(/ResourceDeleteSuccessData/g, 'ResourceDeleteSuccessEventData')
        .replace(/ResourceWriteCancelData/g, 'ResourceWriteCancelEventData')
        .replace(/ResourceWriteFailureData/g, 'ResourceWriteFailureEventData')
        .replace(/ResourceWriteSuccessData/g, 'ResourceWriteSuccessEventData')


  # TODO:
  # missing:
  #
  #   subscriptiondeletedeventdata
  #   subscriptionvalidationeventdata
```
