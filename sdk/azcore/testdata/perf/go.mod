module github.com/Azure/azure-sdk-for-go/sdk/azcore/testdata/perf

go 1.23.0

toolchain go1.24.0

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.17.0
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.10.0
)

require (
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/text v0.23.0 // indirect
)

replace github.com/Azure/azure-sdk-for-go/sdk/azcore => ../../
