# Code Generation - Azure Blob SDK for Golang

```bash
cd swagger
autorest README.md --use="https://github.com/Azure/autorest.go/releases/download/v4.0.0-preview.27/autorest-go-4.0.0-preview.27.tgz"
gofmt -w generated/*
```

### Settings

```yaml
clear-output-folder: true
license-header: MICROSOFT_MIT_NO_VERSION
input-file: "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/storage-dataplane-preview/specification/storage/data-plane/Microsoft.BlobStorage/preview/2019-12-12/blob.json"
module: "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
credential-scope: "https://storage.azure.com/.default"
output-folder: generated/
file-prefix: "zz_generated_"
modelerfour.lenient-model-deduplication: true
openapi-type: "data-plane"
verbose: true
security: AzureKey
module-version: "0.1.0"
```
