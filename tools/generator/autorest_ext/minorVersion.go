package autorest_ext

import (
	"github.com/Azure/azure-sdk-for-go/tools/generator/sdk"
	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest"
)

func CanIncludeInMinor(r autorest.ChangelogResult) bool {
	if sdk.IsPreviewPackage(r.PackageName) {
		return true
	}
	return !r.Changelog.HasBreakingChanges()
}
