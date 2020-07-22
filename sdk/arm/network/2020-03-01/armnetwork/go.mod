module github.com/Azure/azure-sdk-for-go/sdk/arm/network/2020-03-01/armnetwork

go 1.13

require (
	github.com/Azure/azure-sdk-for-go/sdk/armcore v0.1.0
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.9.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v0.0.0-00010101000000-000000000000
	github.com/Azure/azure-sdk-for-go/sdk/to v0.1.0
)

replace github.com/Azure/azure-sdk-for-go/sdk/azidentity => ../../../../azidentity

replace github.com/Azure/azure-sdk-for-go/sdk/arm/resources/2019-05-01/armresources => ../../../resources/2019-05-01/armresources
