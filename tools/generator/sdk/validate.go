package sdk

import (
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/internal/utils"
)

func IsPreviewPackage(pkgName string) bool {
	return strings.HasPrefix(utils.NormalizePath(pkgName), "services/preview/")
}
