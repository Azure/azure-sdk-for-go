module github.com/Azure/azure-sdk-for-go/sdk/tables/aztable

go 1.16

replace github.com/Azure/azure-sdk-for-go/sdk/internal => ../../internal

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.18.1
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v0.9.3
	github.com/Azure/azure-sdk-for-go/sdk/internal v0.7.0
	github.com/Azure/azure-sdk-for-go/sdk/to v0.1.4
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
