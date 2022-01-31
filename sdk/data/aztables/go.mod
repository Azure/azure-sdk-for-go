module github.com/Azure/azure-sdk-for-go/sdk/data/aztables

go 1.16

replace github.com/Azure/azure-sdk-for-go/sdk/internal => ../../internal

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.21.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v0.13.0
	github.com/Azure/azure-sdk-for-go/sdk/internal v0.9.0
	github.com/stretchr/testify v1.7.0
)
