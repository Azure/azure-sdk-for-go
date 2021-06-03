module github.com/Azure/azure-sdk-for-go/sdk/armcore

go 1.14

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.14.0
	github.com/Azure/azure-sdk-for-go/sdk/internal v0.5.1
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b // indirect
)

replace github.com/Azure/azure-sdk-for-go/sdk/azidentity => ../azidentity/

replace github.com/Azure/azure-sdk-for-go/sdk/azcore => ../azcore/
