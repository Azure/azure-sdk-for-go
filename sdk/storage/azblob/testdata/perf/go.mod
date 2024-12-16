module github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/testdata/perf

go 1.18

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.16.0
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.10.0
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.5.0
)

require (
	golang.org/x/net v0.32.0 // indirect
	golang.org/x/text v0.21.0 // indirect
)

replace github.com/Azure/azure-sdk-for-go/sdk/storage/azblob => ../../.
