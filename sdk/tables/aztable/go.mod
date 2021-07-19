module github.com/Azure/azure-sdk-for-go/sdk/tables/aztable

go 1.13

replace github.com/Azure/azure-sdk-for-go/sdk/internal => ../../internal

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.16.2
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v0.9.2 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v0.5.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/net v0.0.0-20210525063256-abc453219eb5 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
