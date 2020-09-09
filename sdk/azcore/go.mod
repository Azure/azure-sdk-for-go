module github.com/Azure/azure-sdk-for-go/sdk/azcore

require (
	github.com/Azure/azure-sdk-for-go/sdk/internal v0.2.3
	golang.org/x/net v0.0.0-20200904194848-62affa334b73
)

go 1.13

replace github.com/Azure/azure-sdk-for-go/sdk/internal => ../internal
