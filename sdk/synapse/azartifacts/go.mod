module github.com/Azure/azure-sdk-for-go/sdk/synapse/azartifacts

go 1.16

require (
	github.com/Azure/azure-sdk-for-go v54.3.0+incompatible
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.16.1
)

retract (
	v0.1.0 // https://github.com/Azure/azure-sdk-for-go/issues/16058
	v0.1.1 // contains retractions
)
