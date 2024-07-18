module github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/testdata/perf

go 1.18

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.12.0
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.9.1
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.4.0
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.13.0
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.10.0
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.4.0
)

require (
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/text v0.16.0 // indirect
)

replace github.com/Azure/azure-sdk-for-go/sdk/storage/azblob => ../../.
