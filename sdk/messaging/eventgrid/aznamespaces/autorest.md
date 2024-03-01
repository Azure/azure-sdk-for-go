## Go

``` yaml
title: EventGridClient
description: Azure Event Grid client
generated-metadata: false
clear-output-folder: false
go: true
input-file: 
    - https://raw.githubusercontent.com/Azure/azure-rest-api-specs/2264262e0c7575a794cc395609d2342c7e598149/specification/eventgrid/data-plane/Microsoft.EventGrid/preview/2023-10-01-preview/EventGrid.json
license-header: MICROSOFT_MIT_NO_VERSION
module: github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/aznamespaces
openapi-type: "data-plane"
output-folder: ../aznamespaces
override-client-name: Client
security: "AADToken"
use: "@autorest/go@4.0.0-preview.63"
version: "^3.0.0"
slice-elements-byval: true
remove-non-reference-schema: true
```

Make sure the content type is setup properly for publishing single and multiple events.

```yaml
directive:
  - from: 
    - client.go
    where: $
    transform: | 
      return $.replace(
        /(func \(client \*Client\) publishCloudEventsCreateRequest.+?)return req, nil/s, 
        '$1\nreq.Raw().Header.Set("Content-type", "application/cloudevents-batch+json; charset=utf-8")\nreturn req, nil');
  - from: 
    - client.go
    where: $
    transform: | 
      return $.replace(
        /(func \(client \*Client\) publishCloudEventCreateRequest.+?)return req, nil/s, 
        '$1\nreq.Raw().Header.Set("Content-type", "application/cloudevents+json; charset=utf-8")\nreturn req, nil');        
```

Fix the error type so it's a bit more presentable, and looks like an error for this package.

```yaml
directive:
  - from: swagger-document
    where: $.definitions["Azure.Core.Foundations.Error"]
    debug: true
    transform: |
      $.properties = { 
        code: $.properties["code"],
        message: {
          ...$.properties["message"],
          "x-ms-client-name": "InternalErrorMessageRename"
        },
      };
      $["x-ms-client-name"] = "Error";
  
  - from: swagger-document
    where: $.definitions
    transform: delete $["Azure.Core.Foundations.InnerError"];

  - from: 
    - models.go
    - models_serde.go
    where: $
    transform: return $.replace(/InternalErrorMessageRename/g, "message");

  - from: 
    - models.go
    - models_serde.go
    where: $
    transform: |
      return $
        .replace(/\/\/ AzureCoreFoundationsErrorResponse.+?\n}/gs, "")
        .replace(/\/\/ MarshalJSON implements the json\.Marshaller interface for type AzureCoreFoundationsErrorResponse\..+?\n}/gs, "")
        .replace(/\/\/ UnmarshalJSON implements the json\.Unmarshaller interface for type AzureCoreFoundationsErrorResponse\..+?\n}/gs, "");
```

Trim out the 'Interface any' for types that are empty.

```yaml
directive:
  - from: responses.go
    where: $
    transform: $.replace(/\s+\/\/ Anything\s+Interface any/sg, "$1");
```

For functions that have empty responses (ie, PublishCloudEvent 
and PublishCloudEvents) we can remove the schema attribute, which cleans 
up the PublishCloudEventResponse/PublishCloudEventsResponse
so they don't have a vestigial `Interface any` field.

```yaml
directive:
  # remove the 'Interface any' that's generated for an empty response object.
  - from:
      - swagger-document
    where: $["x-ms-paths"]["/topics/{topicName}:publish?_overload=publishCloudEvents"].post.responses["200"]
    transform: delete $["schema"];
  - from:
      - swagger-document
    where: $["paths"]["/topics/{topicName}:publish"].post.responses["200"]
    transform: delete $["schema"];
```

Use azcore's CloudEvent type instead of a locally generated version.

```yaml
directive:
  # replace references to the "generated" CloudEvent to the actual version in azcore/messaging
  - from:
      - client.go
      - models.go
      - responses.go
      - options.go
    where: $
    transform: |
      return $.replace(/\[\]CloudEvent/g, "[]messaging.CloudEvent")
              .replace(/\*CloudEvent/g, "messaging.CloudEvent")
              .replace(/event CloudEvent/g, "event messaging.CloudEvent")
  - from: swagger-document
    where: $.definitions.CloudEvent
    transform: $["x-ms-external"] = true
  # make the endpoint a parameter of the client constructor
  - from: swagger-document
    where: $["x-ms-parameterized-host"]
    transform: $.parameters[0]["x-ms-parameter-location"] = "client"
  # delete client name prefix from method options and response types
  - from:
      - client.go
      - models.go
      - responses.go
      - options.go
    where: $
    transform: return $.replace(/Client(\w+)((?:Options|Response))/g, "$1$2");
```

Fix incorrect string formatting for the "release with delay"

```yaml
directive:
  - from: swagger-document
    where: $.paths["/topics/{topicName}/eventsubscriptions/{eventSubscriptionName}:release"].post.parameters[3]
    transform: $.type = "integer";
  - from: client.go
    where: $
    transform: return $.replace(/string\(\*options.ReleaseDelayInSeconds\)/g, "fmt.Sprintf(\"%d\", *options.ReleaseDelayInSeconds)")
```

Add doc for ReleaseDelay enum

```yaml
directive:
  - from: constants.go
    where: $
    transform: return $.replace(/type ReleaseDelay int32/, "// ReleaseDelay indicates how long the service should delay before releasing an event.\ntype ReleaseDelay int32")
```

We want to flatten out the settlement arg functions so we'll internalize
them and do the flattening in custom code.

```yaml
directive:
  # Rename the functions so they're internal

  - from: client.go
    where: $
    transform: return $.replace(/func \(client \*Client\) RejectCloudEvents\(/, "func \(client \*Client\) internalRejectCloudEvents(")
  - from: client.go
    where: $
    transform: return $.replace(/func \(client \*Client\) AcknowledgeCloudEvents\(/, "func \(client \*Client\) internalAcknowledgeCloudEvents(")
  - from: client.go
    where: $
    transform: return $.replace(/func \(client \*Client\) ReleaseCloudEvents\(/, "func \(client \*Client\) internalReleaseCloudEvents(")
  - from: client.go
    where: $
    transform: return $.replace(/func \(client \*Client\) RenewCloudEventLocks\(/, "func \(client \*Client\) internalRenewCloudEventLocks(")

  # Rename the old param bags to be internal as well

  - from: 
    - client.go
    - models.go
    - models_serde.go
    where: $
    transform: return $.replace(/\bReleaseOptions\b/g, "releaseOptions")
  - from: 
    - client.go
    - models.go
    - models_serde.go
    where: $
    transform: return $.replace(/\bRejectOptions\b/g, "rejectOptions")
  - from: 
    - client.go
    - models.go
    - models_serde.go
    where: $
    transform: return $.replace(/\bAcknowledgeOptions\b/g, "acknowledgeOptions")
  - from: 
    - client.go
    - models.go
    - models_serde.go
    where: $
    transform: return $.replace(/\RenewLockOptions\b/g, "renewLockOptions")
```
