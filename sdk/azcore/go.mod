module github.com/Azure/azure-sdk-for-go/sdk/azcore

go 1.26.0

require (
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.13.0
	github.com/stretchr/testify v1.12.1
)

require (
	go.yaml.in/yaml/v3 v3.0.6 // indirect
	golang.org/x/net v0.62.0 // indirect
	golang.org/x/text v0.43.0 // indirect
)

replace github.com/Azure/azure-sdk-for-go/sdk/internal => ../internal
