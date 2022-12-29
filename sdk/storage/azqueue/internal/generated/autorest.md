# Code Generation - Azure Queue SDK for Golang

### Settings

```yaml
go: true
clear-output-folder: false
version: "^3.0.0"
license-header: MICROSOFT_MIT_NO_VERSION
input-file: "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/main/specification/storage/data-plane/Microsoft.QueueStorage/preview/2018-03-28/queue.json"
credential-scope: "https://storage.azure.com/.default"
output-folder: ../generated
file-prefix: "zz_"
openapi-type: "data-plane"
verbose: true
security: AzureKey
modelerfour:
  group-parameters: false
  seal-single-value-enum-by-default: true
  lenient-model-deduplication: true
export-clients: true
use: "@autorest/go@4.0.0-preview.43"
```


### Remove QueueName from parameter list since it is not needed
``` yaml
directive:
- from: swagger-document
  where: $["x-ms-paths"]
  transform: >
    for (const property in $)
    {
        if (property.includes('/{queueName}/messages/{messageid}'))
        {
            $[property]["parameters"] = $[property]["parameters"].filter(function(param) { return (typeof param['$ref'] === "undefined") || (false == param['$ref'].endsWith("#/parameters/QueueName") && false == param['$ref'].endsWith("#/parameters/MessageId"))});
        }
        else if (property.includes('/{queueName}'))
        {
            $[property]["parameters"] = $[property]["parameters"].filter(function(param) { return (typeof param['$ref'] === "undefined") || (false == param['$ref'].endsWith("#/parameters/QueueName"))});
        }
    }
```

### Fix GeoReplication

``` yaml
directive:
- from: swagger-document
  where: $.definitions
  transform: >
    delete $.GeoReplication.properties.Status["x-ms-enum"];
    $.GeoReplication.properties.Status["x-ms-enum"] = {
        "name": "QueueGeoReplicationStatus",
        "modelAsString": false
    };
```

### Remove pager method (since we implement it ourselves on the client layer) and export various generated methods in service client to utilize them in higher layers

``` yaml
directive:
  - from: zz_service_client.go
    where: $
    transform: >-
      return $.
        replace(/func \(client \*ServiceClient\) NewListQueuesSegmentPager\(.+\/\/ listQueuesSegmentCreateRequest creates the ListQueuesSegment request/s, `// ListQueuesSegmentCreateRequest creates the ListQueuesFlatSegment ListQueuesSegment`).
        replace(/\(client \*ServiceClient\) listQueuesSegmentCreateRequest\(/, `(client *ServiceClient) ListQueuesSegmentCreateRequest(`).
        replace(/\(client \*ServiceClient\) listQueuesSegmentHandleResponse\(/, `(client *ServiceClient) ListQueuesSegmentHandleResponse(`);
```

### Fix function marshal and unmarshal signatures to be non-conflicting

``` yaml
directive:
  - from: zz_models_serde.go
    where: $
    transform: >-
      return $.
        replace(/d\s*\*xml\.Decoder/g, "dec *xml.Decoder").
        replace(/d\.DecodeElement\(/g, "dec.DecodeElement(").
        replace(/e\s*\*xml\.Encoder/g, "enc *xml.Encoder").
        replace(/e\.EncodeElement\(/g, "enc.EncodeElement(");
```