module github.com/Azure/azure-sdk-for-go/sdk/compute/armcompute

go 1.16

require (
	github.com/Azure/azure-sdk-for-go v57.0.0+incompatible
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.19.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v0.10.0
)

// To better align with the Azure SDK guidelines (https://azure.github.io/azure-sdk/general_introduction.html), we have decided to change the module path to "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute". Therefore, we are deprecating the old module path (which is "github.com/Azure/azure-sdk-for-go/sdk/compute/armcompute") to avoid confusion.
retract [v0.0.1, v0.3.1]