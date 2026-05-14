module github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/testdata/perf

go 1.25.0

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.21.1
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.12.0
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.6.4
)

require (
	golang.org/x/net v0.54.0 // indirect
	golang.org/x/text v0.37.0 // indirect
)

replace github.com/Azure/azure-sdk-for-go/sdk/storage/azblob => ../../.
