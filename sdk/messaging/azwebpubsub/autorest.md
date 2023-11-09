## Go

```yaml
title: WebPubSub
description: Azure Web PubSub client
clear-output-folder: false
export-clients: true
go: true
input-file: https://github.com/Azure/azure-rest-api-specs/blob/052a4b8d50bfd5595a8b5b506015d18f2b65998d/specification/webpubsub/data-plane/WebPubSub/stable/2023-07-01/webpubsub.json
license-header: MICROSOFT_MIT_NO_VERSION
module: github.com/Azure/azure-sdk-for-go/sdk/messaging/azwebpubsub
openapi-type: "data-plane"
output-folder: ../azwebpubsub
use: "@autorest/go@4.0.0-preview.57"
directive:
    # Make GenerateClientToken internal.
    - from: client.go
      where: $
      transform: return $.replace(/\bGenerateClientToken\b/g, "generateClientToken");
    # Make *Exists internal until SDK supports it.
    - from: client.go
      where: $
      transform: return $.replace(/\b(Group|Connection|User)Exists\b/g, function(match, group) { return group.toLowerCase() + "Exists";});
    # Add more properties to lient
    - from: client.go
      where: $
      transform: >-
        return $.replace(
            /(type Client struct[^}]+})/s, 
            "type Client struct {\n	internal *azcore.Client\n	endpoint string\n	hub string\n	key *string\n}")

```