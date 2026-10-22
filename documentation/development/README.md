# Development Documentation

This folder contains details on how to generate a Go SDK using the TypeSpec code generator.

Development Flow:
* [Set Up](https://github.com/Azure/azure-sdk-for-go/tree/main/documentation/development/setup.md)
* [Generate SDKs](https://github.com/Azure/azure-sdk-for-go/tree/main/documentation/development/generate.md)
* [Test SDKs](https://github.com/Azure/azure-sdk-for-go/tree/main/documentation/development/testing.md)
* [Release SDKs](https://github.com/Azure/azure-sdk-for-go/tree/main/documentation/development/release.md)

Other folders and files:
* [Go mgmt. SDK release rotation runbook](https://github.com/Azure/azure-sdk-for-go/blob/main/documentation/development/ARM/mgmt-release-rotation-runbook.md)- Daily mail triage, workflow-result review, and PR approval for generated management SDKs
* [ARM](https://github.com/Azure/azure-sdk-for-go/tree/main/documentation/development/ARM)- Azure Resource Manager/management-plane SDK specific documentation (i.e. modules with the `arm` prefix)
* [breaking-changes](https://github.com/Azure/azure-sdk-for-go/tree/main/documentation/development/breaking-changes)- guide for breaking code changes plus how to deprecate a module
* [TROUBLESHOOTING.md](https://github.com/Azure/azure-sdk-for-go/tree/main/documentation/development/TROUBLESHOOTING.md)- troubleshooting code generation issues