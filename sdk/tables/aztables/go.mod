module github.com/Azure/azure-sdk-for-go/sdk/tables/aztables

go 1.13

replace github.com/Azure/azure-sdk-for-go/sdk/internal => ../../internal

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.14.2
	github.com/Azure/azure-sdk-for-go/sdk/internal v0.5.0
	github.com/google/uuid v1.2.0 // indirect
	github.com/stretchr/testify v1.7.0
)
