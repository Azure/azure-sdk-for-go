//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package path

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
)

// CPKInfo contains a group of parameters for client provided encryption key.
type CPKInfo = blob.CPKInfo

// CPKScopeInfo contains a group of parameters for client provided encryption scope.
type CPKScopeInfo = blob.CPKScopeInfo

// AccessConditions identifies container-specific access conditions which you optionally set.
type AccessConditions = exported.PathAccessConditions

// CreateOptions contains the optional parameters when calling the Create operation.
type CreateOptions struct {
	AccessConditions  *AccessConditions
	ContinuationToken *string
	Permissions       *string
	Properties        *string
	SourceLeaseID     *string
	Umask             *string
}

// DeleteOptions contains the optional parameters when calling the Delete operation.
type DeleteOptions struct {
}

// SetAccessControlOptions contains the optional parameters when calling the SetAccessControl operation.
type SetAccessControlOptions struct {
}

// SetAccessControlRecursiveOptions contains the optional parameters when calling the SetAccessControlRecursive operation.
type SetAccessControlRecursiveOptions struct {
}

// GetPropertiesOptions contains the optional parameters when calling the GetProperties operation.
type GetPropertiesOptions = blob.GetPropertiesOptions

// SetMetadataOptions contains the optional parameters when calling the SetMetadata operation.
type SetMetadataOptions = blob.SetMetadataOptions

// SetHTTPHeadersOptions contains the optional parameters when calling the SetHTTPHeaders operation.
type SetHTTPHeadersOptions = blob.SetHTTPHeadersOptions
