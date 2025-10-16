module github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/testdata/perf

go 1.24.0

toolchain go1.24.5

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.19.1
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.11.2
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.6.1
)

require (
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/text v0.30.0 // indirect
)

replace github.com/Azure/azure-sdk-for-go/sdk/storage/azblob => ../../.
