module github.com/Azure/azure-sdk-for-go/sdk/azcore/testdata/perf

go 1.18

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.11.1
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.9.0
)

require (
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/text v0.16.0 // indirect
)

replace github.com/Azure/azure-sdk-for-go/sdk/azcore => ../../
