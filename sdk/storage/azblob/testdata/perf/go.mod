module github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/testdata/perf

go 1.17

replace github.com/Azure/azure-sdk-for-go/sdk/internal => ../../../../internal

replace github.com/Azure/azure-sdk-for-go/sdk/storage/azblob => ../../.

require github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v0.2.0

require github.com/Azure/azure-sdk-for-go/sdk/internal v0.8.3

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.21.0 // indirect
	golang.org/x/net v0.0.0-20210805182204-aaa1db679c0d // indirect
	golang.org/x/text v0.3.7 // indirect
)
