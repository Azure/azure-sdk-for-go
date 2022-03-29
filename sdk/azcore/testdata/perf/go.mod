module github.com/Azure/azure-sdk-for-go/sdk/azcore/testdata/perf

go 1.17

replace github.com/Azure/azure-sdk-for-go/sdk/azcore => ../../.

replace github.com/Azure/azure-sdk-for-go/sdk/internal => ../../../internal

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.21.0
	github.com/Azure/azure-sdk-for-go/sdk/internal v0.9.1
)

require (
	golang.org/x/net v0.0.0-20210610132358-84b48f89b13b // indirect
	golang.org/x/text v0.3.6 // indirect
)
