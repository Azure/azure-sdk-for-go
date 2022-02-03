module github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets

go 1.16

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.21.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v0.13.0
	github.com/Azure/azure-sdk-for-go/sdk/internal v0.9.0
	github.com/Azure/azure-sdk-for-go/sdk/keyvault/internal v0.2.0
	github.com/stretchr/testify v1.7.0
)

replace github.com/Azure/azure-sdk-for-go/sdk/keyvault/internal => ../internal
