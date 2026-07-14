module github.com/Azure/azure-sdk-for-go/sdk/data/aztables/testdata/perf

go 1.25.0

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.21.1
	github.com/Azure/azure-sdk-for-go/sdk/data/aztables v1.4.1
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.12.0
)

require (
	golang.org/x/net v0.55.0 // indirect
	golang.org/x/text v0.37.0 // indirect
)

replace github.com/Azure/azure-sdk-for-go/sdk/data/aztables => ../..
