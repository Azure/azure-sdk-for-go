module github.com/Azure/azure-sdk-for-go/sdk/data/aztables/testdata/perf

go 1.23.0

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.20.0
	github.com/Azure/azure-sdk-for-go/sdk/data/aztables v1.4.0
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.11.2
)

require (
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/text v0.28.0 // indirect
)

replace github.com/Azure/azure-sdk-for-go/sdk/data/aztables => ../..
