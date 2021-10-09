module github.com/Azure/azure-sdk-for-go/sdk/web/armweb

go 1.16

require (
	github.com/Azure/azure-sdk-for-go v57.1.0+incompatible
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.19.0
)

// To better align with the Azure SDK guidelines (https://azure.github.io/azure-sdk/general_introduction.html), we have decided to change the module path to "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/web/armweb". Therefore, we are deprecating the old module path (which is "github.com/Azure/azure-sdk-for-go/sdk/web/armweb") to avoid confusion.
retract [v0.1.0,v0.2.1]