module github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/testdata/perf

go 1.18

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.11.1
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.8.0
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.3.2
)

require (
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/text v0.15.0 // indirect
)

replace github.com/Azure/azure-sdk-for-go/sdk/storage/azblob => ../../.
