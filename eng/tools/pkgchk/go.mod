module github.com/Azure/azure-sdk-for-go/eng/tools/pkgchk

go 1.13

require (
	github.com/Azure/azure-sdk-for-go v54.2.1+incompatible
	github.com/Azure/azure-sdk-for-go/eng/tools/internal v0.1.0
	github.com/spf13/cobra v1.1.3
)

replace github.com/Azure/azure-sdk-for-go/eng/tools/internal => ../internal
