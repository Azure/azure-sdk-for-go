module github.com/Azure/azure-sdk-for-go/eng/tools/azperf

go 1.17

replace github.com/Azure/azure-sdk-for-go/sdk/data/aztables => ../../../sdk/data/aztables

replace github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys => ../../../sdk/keyvault/azkeys

require github.com/spf13/cobra v1.2.1

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.20.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v0.12.0
	github.com/Azure/azure-sdk-for-go/sdk/data/aztables v0.4.0
	github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys v0.1.0
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v0.2.0
)

require (
	github.com/Azure/azure-sdk-for-go/sdk/internal v0.8.3 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/keyvault/internal v0.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/pkg/browser v0.0.0-20180916011732-0a3d74bf9ce4 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897 // indirect
	golang.org/x/net v0.0.0-20210805182204-aaa1db679c0d // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
