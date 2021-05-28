module github.com/Azure/azure-sdk-for-go/sdk/armcore

go 1.14

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.14.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v0.9.1
	github.com/Azure/azure-sdk-for-go/sdk/internal v0.5.0
)

replace github.com/Azure/azure-sdk-for-go/sdk/azidentity => ../azidentity/