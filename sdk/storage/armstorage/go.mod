// Deprecated: use github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage instead.
module github.com/Azure/azure-sdk-for-go/sdk/storage/armstorage

go 1.16

require (
	github.com/Azure/azure-sdk-for-go v57.1.0+incompatible
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.19.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v0.10.0
)

// use github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage instead.
retract [v0.0.1, v0.2.1]