# Autorest config for Azure Container Registry Go client

> see https://aka.ms/autorest

## Configuration

```yaml
input-file: https://github.com/Azure/azure-rest-api-specs/blob/c8d9a26a2857828e095903efa72512cf3a76c15d/specification/containerregistry/data-plane/Azure.ContainerRegistry/stable/2021-07-01/containerregistry.json
license-header: MICROSOFT_MIT_NO_VERSION
go: true
clear-output-folder: false
export-clients: true
openapi-type: "data-plane"
output-folder: ../azacr
use: "@autorest/go@4.0.0-preview.44"
honor-body-placement: true
remove-unreferenced-types: true
```

## Customizations

See the [AutoRest samples](https://github.com/Azure/autorest/tree/master/Samples/3b-custom-transformations)
for more about how we're customizing things.

### Remove response for "ContainerRegistry_DeleteRepository" operation

so that the generated code doesn't return a response for the deleted repository operation.

```yaml
directive:
  - from: swagger-document
    where: $["paths"]["/acr/v1/{name}"]
    transform: >
      delete $.delete["responses"]["202"].schema
```

### Remove response for "ContainerRegistryBlob_DeleteBlob" operation

so that the generated code doesn't return a response for the deleted blob operation.

```yaml
directive:
  - from: swagger-document
    where: $["paths"]["/v2/{name}/blobs/{digest}"]
    transform: >
      delete $.delete["responses"]["202"].schema
```

### Remove "Authentication_GetAcrAccessTokenFromLogin" operation

as the service team discourage using username/password to authenticate.

```yaml
directive:
  - from: swagger-document
    where: $["paths"]["/oauth2/token"]
    transform: >
      delete $.get
```

### Remove "definitions.TagAttributesBase.properties.signed"

as we don't have customer scenario using it.

```yaml
directive:
  - from: swagger-document
    where: $.definitions.TagAttributesBase
    transform: >
      delete $.properties.signed
```

### Remove "definitions.ManifestAttributesBase.properties.configMediaType"

as we don't have customer scenario using it.

```yaml
directive:
  - from: swagger-document
    where: $.definitions.ManifestAttributesBase
    transform: >
      delete $.properties.configMediaType
```

### Change "parameters.ApiVersionParameter.required" to true

so that the API version could be removed from client parameter.

```yaml
directive:
  - from: swagger-document
    where: $.parameters.ApiVersionParameter
    transform: >
      $.required = true
```

### Change NextLink client name to nextLink

```yaml
directive:
  from: swagger-document
  where: $.parameters.NextLink
  transform: >
    $["x-ms-client-name"] = "nextLink"
```

### Updates to OciManifest

```yaml
directive:
  from: swagger-document
  where: $.definitions.OCIManifest
  transform: >
    delete $["allOf"];
    $.properties["schemaVersion"] = {
      "type": "integer",
      "description": "Schema version"
    };
```

### Take stream as manifest body

```yaml
directive:
  from: swagger-document
  where: $.parameters.ManifestBody
  transform: >
    $.schema = {
      "type": "string",
      "format": "binary"
    }
```

### Change list order by param to enum
```yaml
directive:
  - from: containerregistry.json
    where: $.paths["/acr/v1/{name}/_tags"].get
    transform: >
      $.parameters.splice(3, 1);
      $.parameters.push({
        "name": "orderby",
        "x-ms-client-name": "OrderBy",
        "in": "query",
        "required": false,
        "x-ms-parameter-location": "method",
        "type": "string",
        "description": "Sort options for ordering tags in a collection.",
        "enum": [
          "none",
          "timedesc",
          "timeasc"
        ],
        "x-ms-enum": {
          "name": "ArtifactTagOrderBy",
          "values": [
            {
              "value": "none",
              "name": "None",
              "description": "Do not provide an orderby value in the request."
            },
            {
              "value": "timedesc",
              "name": "LastUpdatedOnDescending",
              "description": "Order tags by LastUpdatedOn field, from most recently updated to least recently updated."
            },
            {
              "value": "timeasc",
              "name": "LastUpdatedOnAscending",
              "description": "Order tags by LastUpdatedOn field, from least recently updated to most recently updated."
            }
          ]
        }
      });
  - from: containerregistry.json
    where: $.paths["/acr/v1/{name}/_manifests"]
    transform: >
      $.get.parameters.splice(3, 1);
      $.get.parameters.push({
        "name": "orderby",
        "x-ms-client-name": "OrderBy",
        "in": "query",
        "required": false,
        "x-ms-parameter-location": "method",
        "type": "string",
        "description": "Sort options for ordering manifests in a collection.",
        "enum": [
          "none",
          "timedesc",
          "timeasc"
        ],
        "x-ms-enum": {
          "name": "ArtifactManifestOrderBy",
          "values": [
            {
              "value": "none",
              "name": "None",
              "description": "Do not provide an orderby value in the request."
            },
            {
              "value": "timedesc",
              "name": "LastUpdatedOnDescending",
              "description": "Order manifests by LastUpdatedOn field, from most recently updated to least recently updated."
            },
            {
              "value": "timeasc",
              "name": "LastUpdatedOnAscending",
              "description": "Order manifest by LastUpdatedOn field, from least recently updated to most recently updated."
            }
          ]
        }
      });
```

### Rename paged operations from Get* to List*

```yaml
directive:
  - rename-operation:
      from: ContainerRegistry_GetManifests
      to: ContainerRegistry_ListManifests
  - rename-operation:
      from: ContainerRegistry_GetRepositories
      to: ContainerRegistry_ListRepositories
  - rename-operation:
      from: ContainerRegistry_GetTags
      to: ContainerRegistry_ListTags
```

### Change ContainerRegistry_CreateManifest behaviour

```yaml
directive:
  from: swagger-document
  where: $.paths["/v2/{name}/manifests/{reference}"].put
  transform: >
    $.consumes.push("application/vnd.oci.image.manifest.v1+json");
    delete $.responses["201"].schema;
```

### Change ContainerRegistry_GetManifest behaviour

```yaml
directive:
  from: swagger-document
  where: $.paths["/v2/{name}/manifests/{reference}"].get.responses["200"]
  transform: >
    $.schema = {
      type: "string",
      format: "file"
    };
    $.headers = {
      "Docker-Content-Digest": {
        "type": "string",
        "description": "Digest of the targeted content for the request."
      }
    };
```

### Remove generated constructors

```yaml
directive:
  - from: 
      - authentication_client.go
      - containerregistry_client.go
      - containerregistryblob_client.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+func New.+Client.+\{\s(?:.+\s)+\}\s/, "");
```

### Rename operations
```yaml
directive:
  - rename-operation:
      from: ContainerRegistry_GetProperties
      to: ContainerRegistry_GetRepositoryProperties
  - rename-operation:
      from: ContainerRegistry_UpdateProperties
      to: ContainerRegistry_UpdateRepositoryProperties
  - rename-operation:
      from: ContainerRegistry_UpdateTagAttributes
      to: ContainerRegistry_UpdateTagProperties
  - rename-operation:
      from: ContainerRegistry_CreateManifest
      to: ContainerRegistry_UploadManifest
```

### Rename parameter name
```yaml
directive:
  from: swagger-document
  where: $.parameters
  transform: >
    $.DigestReference["x-ms-client-name"] = "digest";
    $.TagReference["x-ms-client-name"] = "tag";
```

### Hide some of generated operation
```yaml
directive:
  - from:
      - containerregistry_client.go
    where: $
    transform: return $.replace(/DeleteManifest\(ctx/, "deleteManifest\(ctx").replace(/GetManifestProperties\(ctx/, "getManifestProperties\(ctx").replace(/UpdateManifestProperties\(ctx/, "updateManifestProperties\(ctx");
```

### Add 202 response to ContainerRegistryBlob_MountBlob
```yaml
directive:
  from: swagger-document
  where: $.paths["/v2/{name}/blobs/uploads/"]
  transform: >
    $.post["responses"]["202"] = $.post["responses"]["201"];
```

### Extract and add endpoint for nextLink
```yaml
directive:
  - from:
      - containerregistry_client.go
    where: $
    transform: return $.replaceAll(/result\.Link = &val/g, "val = runtime.JoinPaths(client.endpoint, extractNextLink(val))\n\t\tresult.Link = &val");
```