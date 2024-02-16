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
output-folder: ../azsystemevents
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
```

```yaml
directive:
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
```

```yaml
directive:
  - from: 
    - models.go
    - models_serde.go
    where: $
    transform: | 
      return $
        .replace(/ChannelLatencyMs \*string/g, "ChannelLatencyMS *string")
        .replace(/m.ChannelLatencyMs/g, "m.ChannelLatencyMS");
```

```yaml
directive:
  - from: constants.go
    where: $
    transform: return $.replace(/EventGridMqttClientDisconnectionReason/g, "EventGridMQTTClientDisconnectionReason")
  - from: models.go
    where: $
    transform: return $.replace(/DisconnectionReason \*EventGridMqttClientDisconnectionReason/, "DisconnectionReason *EventGridMQTTClientDisconnectionReason")
```

Manually map n/a to the zero value for the type for `MediaLiveEventChannelArchiveHeartbeatEventData` and `MediaLiveEventIngestHeartbeatEventData`

```yaml
directive:
  - from: models_serde.go
    where: $
    transform: |
      return $
        .replace(/(\s+err = unpopulate\(val, "IngestDriftValue", &m.IngestDriftValue\))/, "$1\nfixNAValue(&m.IngestDriftValue)")
        .replace(/(\s+err = unpopulate\(val, "ChannelLatencyMs", &m.ChannelLatencyMS\))/, "$1\nfixNAValue(&m.ChannelLatencyMS)");
```

Rename `AcsRouterWorkerSelector.TTLSeconds` to `TimeToLive`

```yaml
directive:
  - from: 
    - models.go
    where: $
    transform: return $.replace(/TTLSeconds \*float32/g, "TimeToLive *float32");
  - from: 
    - models_serde.go
    where: $
    transform: return $.replace(/a\.TTLSeconds/g, "a.TimeToLive");
```

Fix the EventGridEvent deserialization to work similar to CloudEvent.

```yaml
directive:
  - from: models_serde.go
    where: $
    transform: |
      return $.replace(/err = unpopulate\(val, "Data", &e.Data\)/, "e.Data = []byte(val)")
```
