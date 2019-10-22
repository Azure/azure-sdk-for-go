module github.com/Azure/azure-sdk-for-go/sdk/azidentity

replace github.com/Azure/azure-sdk-for-go/sdk/azcore => ../azcore

go 1.13

require (
	github.com/Azure/azure-amqp-common-go/v2 v2.1.0
	github.com/Azure/azure-event-hubs-go v1.3.1
	github.com/Azure/azure-event-hubs-go/v2 v2.0.3
	github.com/Azure/azure-sdk-for-go v34.0.0+incompatible
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.0.0-00010101000000-000000000000
	github.com/Azure/go-autorest/autorest v0.9.1 // indirect
	github.com/Azure/go-autorest/autorest/azure/auth v0.3.0 // indirect
	golang.org/x/tools v0.0.0-20191002212750-6fe9ea94a73d // indirect
)
