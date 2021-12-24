// Deprecated: use github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice instead.
module github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/web/armweb

go 1.16

require (
	github.com/Azure/azure-sdk-for-go v59.0.0+incompatible
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.20.0
)

// use github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice instead.
retract [v0.1.0,v0.2.1]