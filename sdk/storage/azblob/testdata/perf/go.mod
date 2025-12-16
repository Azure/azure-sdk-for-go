module github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/testdata/perf

go 1.24.0

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.20.0
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.11.2
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.6.3
)

require (
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/text v0.31.0 // indirect
)

replace github.com/Azure/azure-sdk-for-go/sdk/storage/azblob => ../../.
