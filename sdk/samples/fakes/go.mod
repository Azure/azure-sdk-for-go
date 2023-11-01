module fakes

go 1.18

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.9.0-beta.1
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5 v5.3.0-beta.2
)

require (
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.4.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/text v0.13.0 // indirect
)

replace github.com/Azure/azure-sdk-for-go/sdk/azcore => ../../azcore
