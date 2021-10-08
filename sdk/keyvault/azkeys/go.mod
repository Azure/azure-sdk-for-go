module github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys

go 1.16

replace github.com/Azure/azure-sdk-for-go/sdk/internal => ../../internal

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.19.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v0.11.0
	github.com/Azure/azure-sdk-for-go/sdk/internal v0.7.1
	github.com/stretchr/testify v1.7.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dnaeon/go-vcr v1.2.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
