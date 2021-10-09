// Deprecated: use github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault instead.
module github.com/Azure/azure-sdk-for-go/sdk/keyvault/armkeyvault

go 1.16

require (
	github.com/Azure/azure-sdk-for-go v57.0.0+incompatible
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.19.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v0.10.0
)

// use github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault instead.
retract [v0.0.1, v0.2.1]