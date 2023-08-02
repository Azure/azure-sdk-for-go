package service

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated_blob"
)

// PublicAccessType defines values for AccessType - private (default) or file or filesystem.
type PublicAccessType = generated_blob.PublicAccessType

const (
	File       PublicAccessType = generated_blob.PublicAccessTypeFile
	Filesystem PublicAccessType = generated_blob.PublicAccessTypeFilesystem
)
