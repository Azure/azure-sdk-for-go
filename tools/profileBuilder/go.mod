module github.com/Azure/azure-sdk-for-go/tools/profileBuilder

go 1.13

require (
	github.com/Azure/azure-sdk-for-go v54.2.1+incompatible
	github.com/Azure/azure-sdk-for-go/tools/internal v0.0.0
	github.com/spf13/cobra v1.1.3
	golang.org/x/tools v0.0.0-20191112195655-aa38f8e97acc
)

replace (
	github.com/Azure/azure-sdk-for-go => github.com/ArcturusZhang/azure-sdk-for-go v54.2.1+incompatible
	github.com/Azure/azure-sdk-for-go/tools/internal => ../internal
)
