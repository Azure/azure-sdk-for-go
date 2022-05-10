module github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/testdata/perf

go 1.18

replace github.com/Azure/azure-sdk-for-go/sdk/internal => ../../../../internal

replace github.com/Azure/azure-sdk-for-go/sdk/storage/azfile => ../../.

require github.com/Azure/azure-sdk-for-go/sdk/storage/azfile v0.1.0

require github.com/Azure/azure-sdk-for-go/sdk/internal v0.9.2

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.23.0 // indirect
	golang.org/x/net v0.0.0-20211015210444-4f30a5c0130f // indirect
	golang.org/x/text v0.3.7 // indirect
)
