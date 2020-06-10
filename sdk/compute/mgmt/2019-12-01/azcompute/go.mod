module github.com/Azure/azure-sdk-for-go/sdk/compute/mgmt/2019-12-01/azcompute

go 1.13

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.8.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v0.0.0-00010101000000-000000000000
	github.com/Azure/azure-sdk-for-go/sdk/network/mgmt/2020-03-01/aznetwork v0.0.0-00010101000000-000000000000
	github.com/Azure/azure-sdk-for-go/sdk/to v0.1.0
)

replace github.com/Azure/azure-sdk-for-go/sdk/azidentity => ../../../../azidentity

replace github.com/Azure/azure-sdk-for-go/sdk/network/mgmt/2020-03-01/aznetwork => ../../../../network/mgmt/2020-03-01/aznetwork
