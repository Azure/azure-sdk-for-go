module github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys

go 1.16

replace github.com/Azure/azure-sdk-for-go/sdk/keyvault/internal => ../internal

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.20.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v0.12.0
	github.com/Azure/azure-sdk-for-go/sdk/internal v0.8.3
	github.com/Azure/azure-sdk-for-go/sdk/keyvault/internal v0.1.0
	github.com/stretchr/testify v1.7.0
)
