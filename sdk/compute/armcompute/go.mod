module github.com/Azure/azure-sdk-for-go/sdk/compute/armcompute

go 1.16

require (
	github.com/Azure/azure-sdk-for-go v57.0.0+incompatible
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.19.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v0.10.0
)

retract (
	v0.3.1
	v0.3.0
	v0.2.0
	v0.1.0
)
