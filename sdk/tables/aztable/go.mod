module github.com/Azure/azure-sdk-for-go/sdk/tables/aztable

go 1.13

replace github.com/Azure/azure-sdk-for-go/sdk/internal => ../../internal
replace github.com/Azure/azure-sdk-for-go/sdk/azidentity => ../../azidentity
replace github.com/Azure/azure-sdk-for-go/sdk/azcore => ../../azcore
replace github.com/Azure/azure-sdk-for-go/sdk/tables/aztable/internal => ./internal

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.16.2
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v0.9.2 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v0.5.2
	github.com/Azure/azure-sdk-for-go/sdk/tables/aztable/internal v0.1.0
	github.com/Azure/azure-sdk-for-go/sdk/to v0.1.4 // indirect
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
