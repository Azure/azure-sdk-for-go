# Azure Blob Storage SDK for Golang

> see https://aka.ms/autorest

### Generation
```bash
cd swagger
autorest --clear-output-folder --use="https://github.com/Azure/autorest.go/releases/download/v4.0.0-preview.22/autorest-go-4.0.0-preview.22.tgz" --license-header=MICROSOFT_MIT_NO_VERSION  --input-file="https://raw.githubusercontent.com/Azure/azure-rest-api-specs/storage-dataplane-preview/specification/storage/data-plane/Microsoft.BlobStorage/preview/2020-02-10/blob.json"  --module="github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"  --credential-scope="https://storage.azure.com/.default"  --output-folder=$(pwd)/generated  --file-prefix="zz_generated_"  --modelerfour.lenient-model-deduplication --openapi-type="data-plane" --go --verbose --security=AzureKey
gofmt -w generated/*
```

### Settings
``` yaml
license-header: MICROSOFT_MIT_NO_VERSION  
input-file: "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/storage-dataplane-preview/specification/storage/data-plane/Microsoft.BlobStorage/preview/2020-02-10/blob.json"  
module: "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"  
credential-scope: "https://storage.azure.com/.default"  
output-folder: $(pwd)/generated  
file-prefix: "zz_generated_"  
modelerfour.lenient-model-deduplication: true 
openapi-type: "data-plane" 
go: true 
verbose: true 
security: AzureKey
```

#### TODO: Get rid of StorageError defined in the generated files since we define it.