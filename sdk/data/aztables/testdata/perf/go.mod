module github.com/Azure/azure-sdk-for-go/sdk/data/aztables/testdata/perf

go 1.18

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.10.0
	github.com/Azure/azure-sdk-for-go/sdk/data/aztables v1.0.0
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.6.0
)

require (
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

replace github.com/Azure/azure-sdk-for-go/sdk/data/aztables => ../..
