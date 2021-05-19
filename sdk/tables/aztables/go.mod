module github.com/Azure/azure-sdk-for-go/sdk/tables/aztables

go 1.13

replace github.com/Azure/azure-sdk-for-go/sdk/internal => ../../internal

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.16.1
	github.com/Azure/azure-sdk-for-go/sdk/internal v0.5.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/net v0.0.0-20210510120150-4163338589ed // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
