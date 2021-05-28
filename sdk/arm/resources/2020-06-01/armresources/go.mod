module github.com/Azure/azure-sdk-for-go/sdk/arm/resources/2020-06-01/armresources

go 1.13

require (
	github.com/Azure/azure-sdk-for-go/sdk/armcore v0.6.0
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.14.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v0.9.1
	github.com/Azure/azure-sdk-for-go/sdk/internal v0.5.0
	github.com/Azure/azure-sdk-for-go/sdk/to v0.1.2
)

replace github.com/Azure/azure-sdk-for-go/sdk/azidentity => ../../../../azidentity/

replace github.com/Azure/azure-sdk-for-go/sdk/armcore => ../../../../armcore
