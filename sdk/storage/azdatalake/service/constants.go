package service

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/filesystem"
)

// PublicAccessType defines values for AccessType - private (default) or file or filesystem.
type PublicAccessType = filesystem.PublicAccessType

const (
	File       PublicAccessType = filesystem.File
	Filesystem PublicAccessType = filesystem.Filesystem
)
