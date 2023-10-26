## Go

```yaml
title: WebPubSub
description: Azure Web PubSub client
clear-output-folder: false
export-clients: true
go: true
input-file: https://github.com/Azure/azure-rest-api-specs/blob/main/specification/webpubsub/data-plane/WebPubSub/stable/2023-07-01/webpubsub.json #TODO: change to hash for commit
license-header: MICROSOFT_MIT_NO_VERSION
module: github.com/Azure/azure-sdk-for-go/sdk/messaging/azwebpubsub
openapi-type: "data-plane"
output-folder: ../azwebpubsub
use: "@autorest/go@4.0.0-preview.57"
```