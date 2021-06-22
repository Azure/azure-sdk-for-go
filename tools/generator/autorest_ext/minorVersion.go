package autorest_ext

import (
	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest"
	"github.com/Azure/azure-sdk-for-go/tools/generator/sdk"
)

func CanIncludeInMinor(r autorest.ChangelogResult) bool {
	if sdk.IsPreviewPackage(r.PackageName) {
		return true
	}
	return !r.Changelog.HasBreakingChanges()
}
