## Go

``` yaml
title: AzureEventGridSystemEvents
description: Azure Event Grid system events
generated-metadata: false
clear-output-folder: false
go: true
require: https://github.com/Azure/azure-rest-api-specs/blob/64819c695760764afa059d799fc7320d3fee33de/specification/eventgrid/data-plane/readme.md
license-header: MICROSOFT_MIT_NO_VERSION
openapi-type: "data-plane"
output-folder: ../azsystemevents
override-client-name: ClientDeleteMe
security: "AADToken"
use: "@autorest/go@4.0.0-preview.63"
version: "^3.0.0"
module: github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/azsystemevents
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
      - responses.go
      - options.go
    where: $
    transform: return $.replace(/CloudEventEvent/g, "CloudEvent");
  - from: 
      - models.go
      - models_serde.go
      - responses.go
      - options.go
    where: $
    transform: return $.replace(/Schema of the Data property of an EventGridEvent/g, "Schema of the Data property of an CloudEvent/EventGridEvent");
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

Remove models that are only used as base models, but aren't needed. The OpenAPI emitter
just duplicates the fields into each child, rather than embedding.  So these aren't needed.

```yaml
directive:
  - from: models.go
    where: $
    transform: return $.replace(/\/\/ (AvsClusterEventData|AvsPrivateCloudEventData|AvsScriptExecutionEventData) - .+?\n}\n/gs, "");
  - from: models_serde.go
    where: $
    transform: |
      for (let name of ["AvsClusterEventData", "AvsPrivateCloudEventData", "AvsScriptExecutionEventData"]) {
        // ex:                '// MarshalJSON implements the json.Marshaller interface for type AvsScriptExecutionEventData.'
        const marshalPrefix = `// MarshalJSON implements the json\.Marshaller interface for type ${name}.+?\n}\n`;
        // ex:                  '// UnmarshalJSON implements the json.Unmarshaller interface for type AvsClusterEventData.'
        const unmarshalPrefix = `// UnmarshalJSON implements the json\.Unmarshaller interface for type ${name}.+?\n}\n`;

        $ = $.replace(new RegExp(marshalPrefix, "gs"), "");
        $ = $.replace(new RegExp(unmarshalPrefix, "gs"), "");
      }      
      return $;
```

Fix acronyms so they match our naming convention.

```yaml
directive:
  - from: 
      - models.go
      - models_serde.go
    where: $
    debug: true
    transform: |
      const acronyms = ["Acs", "Avs", "Iot"];
      for (let acr of acronyms) {
        // ex:
        // '// AcsChatMessageDeletedEventData - Schema'
        // 'type AcsChatMessageDeletedEventData struct'
        // 'Participants []AcsChatThreadParticipantProperties'
        // 'ParticipantRemoved *AcsChatThreadParticipantProperties'
        const re = new RegExp(`([ *\\]])${acr}([A-Za-z0-9]+?(?:EventData|Properties))`, "sg");
        $ = $.replace(re, `$1${acr.toUpperCase()}$2`);
      }
      return $;
```
