module github.com/Azure/azure-sdk-for-go/sdk/samples/aztables

go 1.16

replace github.com/Azure/azure-sdk-for-go/sdk/data/aztables => ../../data/aztables

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.19.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v0.10.0
	github.com/Azure/azure-sdk-for-go/sdk/data/aztables v0.0.0-20210902203352-6fd8cebb673d
)
