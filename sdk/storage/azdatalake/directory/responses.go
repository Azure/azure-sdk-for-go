package directory

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/path"
)

// CreateResponse contains the response fields for the Create operation.
type CreateResponse = generated.PathClientCreateResponse

// DeleteResponse contains the response fields for the Delete operation.
type DeleteResponse = generated.PathClientDeleteResponse

// SetAccessControlResponse contains the response fields for the SetAccessControl operation.
type SetAccessControlResponse = path.SetAccessControlResponse

// SetAccessControlRecursiveResponse contains the response fields for the SetAccessControlRecursive operation.
type SetAccessControlRecursiveResponse = path.SetAccessControlRecursiveResponse

// UpdateAccessControlRecursiveResponse contains the response fields for the UpdateAccessControlRecursive operation.
type UpdateAccessControlRecursiveResponse = path.SetAccessControlRecursiveResponse

// RemoveAccessControlRecursiveResponse contains the response fields for the RemoveAccessControlRecursive operation.
type RemoveAccessControlRecursiveResponse = path.SetAccessControlRecursiveResponse

// GetPropertiesResponse contains the response fields for the GetProperties operation.
type GetPropertiesResponse = blob.GetPropertiesResponse

// SetMetadataResponse contains the response fields for the SetMetadata operation.
type SetMetadataResponse = path.SetMetadataResponse

// SetHTTPHeadersResponse contains the response fields for the SetHTTPHeaders operation.
type SetHTTPHeadersResponse = path.SetHTTPHeadersResponse
