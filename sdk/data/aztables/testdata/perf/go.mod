module github.com/Azure/azure-sdk-for-go/sdk/data/aztables/testdata/perf

go 1.18

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.15.0
	github.com/Azure/azure-sdk-for-go/sdk/data/aztables v1.2.0
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.10.0
)

require (
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/text v0.19.0 // indirect
)

replace github.com/Azure/azure-sdk-for-go/sdk/data/aztables => ../..
