module github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/testdata/perf

go 1.25.0

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.22.0
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.12.0
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.7.0
)

require (
	golang.org/x/net v0.55.0 // indirect
	golang.org/x/text v0.38.0 // indirect
)

replace github.com/Azure/azure-sdk-for-go/sdk/storage/azblob => ../../.
