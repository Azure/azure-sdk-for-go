module github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/testdata/perf

go 1.23.0

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.18.1
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.11.1
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.6.1
)

require (
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/text v0.27.0 // indirect
)

replace github.com/Azure/azure-sdk-for-go/sdk/storage/azblob => ../../.
