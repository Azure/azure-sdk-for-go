module github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets

go 1.16

replace github.com/Azure/azure-sdk-for-go/sdk/internal => ../../internal

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.19.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v0.11.0
	github.com/Azure/azure-sdk-for-go/sdk/internal v0.7.1
	github.com/stretchr/testify v1.7.0
)
