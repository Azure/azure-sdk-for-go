## Go

```yaml
title: WebPubSub
description: Azure Web PubSub client
clear-output-folder: false
slice-elements-byval: true
remove-non-reference-schema: true
go: true
input-file: https://github.com/Azure/azure-rest-api-specs/blob/052a4b8d50bfd5595a8b5b506015d18f2b65998d/specification/webpubsub/data-plane/WebPubSub/stable/2023-07-01/webpubsub.json
license-header: MICROSOFT_MIT_NO_VERSION
module: github.com/Azure/azure-sdk-for-go/sdk/messaging/azwebpubsub
openapi-type: "data-plane"
output-folder: ../azwebpubsub
use: "@autorest/go@4.0.0-preview.60"
directive:
    # Remove HealthAPI
    - from: swagger-document
      remove-operation: 'HealthApi_GetServiceStatus'
    # Rename enum WebPubSubPermission to Permission since the package name already contains WebPubSub.
    - from: 
        - constants.go
        - client.go
      where: $
      transform: return $.replace(/WebPubSubPermission/g, "Permission");
    # Make GenerateClientToken internal.
    - from: client.go
      where: $
      transform: return $.replace(/\bGenerateClientToken\b/g, "generateClientToken");
    # Make *Exists internal until SDK supports it.
    - from: client.go
      where: $
      transform: return $.replace(/\b(Group|Connection|User)Exists\b/g, function(match, group) { return group.toLowerCase() + "Exists";});
    # Make CheckPermission internal until SDK supports it, since it leverage 404 status code
    - from: client.go
      where: $
      transform: return $.replace(/\bCheckPermission\b/g, "checkPermission");
    # Add more properties to the client
    - from: client.go
      where: $
      transform: >-
        return $.replace(
            /(type Client struct[^}]+})/s, 
            "type Client struct {\n	internal *azcore.Client\n	endpoint string\n	key      *string\n}")
    # Add comments to type Permission
    - from: constants.go
      where: $
      transform: >-
        return $.replace(
            /type Permission string/s, 
            "// Permission contains the allowed permissions\ntype Permission string")
    # Add comments to InnerError
    - from: models.go
      where: $
      transform: >-
        return $.replace(
            /type InnerError struct/s, 
            "// InnerError - The inner error object\ntype InnerError struct")
    # delete unused error models
    - from: models.go
      where: $
      transform: return $.replace(/(?:\/\/.*\s)+type (?:ErrorDetail|InnerError).+\{(?:\s.+\s)+\}\s/g, "");
    - from: models_serde.go
      where: $
      transform: return $.replace(/(?:\/\/.*\s)+func \(\w \*?(?:ErrorDetail|InnerError)\).*\{\s(?:.+\s)+\}\s/g, "");
    # delete client name prefix from method options and response types
    - from:
        - client.go
        - models.go
        - models_serde.go
        - options.go
        - response_types.go
      where: $
      transform: return $.replace(/Client(\w+)((?:Options|Response))/g, "$1$2");
```