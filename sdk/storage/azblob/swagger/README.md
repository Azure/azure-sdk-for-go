# Code Generation - Azure Blob SDK for Golang

```bash
cd swagger
autorest README.md --use="https://github.com/Azure/autorest.go/releases/download/v4.0.0-preview.35/autorest-go-4.0.0-preview.35.tgz"
gofmt -w generated/*
```

### Settings

```yaml
clear-output-folder: true
license-header: MICROSOFT_MIT_NO_VERSION
input-file: "./blob.json"
module: "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
credential-scope: "https://storage.azure.com/.default"
output-folder: generated1/
file-prefix: "zz_generated_"
modelerfour.lenient-model-deduplication: true
openapi-type: "data-plane"
verbose: true
security: AzureKey
module-version: "0.3.0"
```

### Fix BlobMetadata.

``` yaml
directive:
- from: swagger-document
  where: $.definitions
  transform: >
    delete $.BlobMetadata["properties"];

```

### Remove DataLake stuff.

``` yaml
directive:
- from: swagger-document
  where: $["x-ms-paths"]
  transform: >
    for (const property in $)
    {
        if (property.includes('filesystem'))
        {
            delete $[property];
        }
    }
```

