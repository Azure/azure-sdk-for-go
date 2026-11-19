module github.com/Azure/azure-sdk-for-go/sdk/data/aztables/testdata/perf

go 1.26.0

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.23.3
	github.com/Azure/azure-sdk-for-go/sdk/data/aztables v1.4.1
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.13.0
)

require (
	golang.org/x/net v0.62.0 // indirect
	golang.org/x/text v0.43.0 // indirect
)

replace github.com/Azure/azure-sdk-for-go/sdk/data/aztables => ../..
