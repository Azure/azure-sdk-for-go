//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package path

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
)

// SetAccessControlOptions contains the optional parameters when calling the SetAccessControl operation. dfs endpoint
type SetAccessControlOptions struct {
	Owner            *string
	Group            *string
	ACL              *string
	Permissions      *string
	AccessConditions *AccessConditions
}

// GetAccessControlOptions contains the optional parameters when calling the GetAccessControl operation.
type GetAccessControlOptions struct {
	UPN              *bool
	AccessConditions *AccessConditions
}

// SetAccessControlRecursiveOptions contains the optional parameters when calling the SetAccessControlRecursive operation.
type SetAccessControlRecursiveOptions struct {
	Owner       *string
	Group       *string
	ACL         *string
	Permissions *string
	//Mode              *SetAccessControlRecursiveMode
	BatchSize         *int32
	MaxBatches        *int32
	ContinueOnFailure *bool
	ContinuationToken *string
}

// UpdateAccessControlRecursiveOptions contains the optional parameters when calling the UpdateAccessControlRecursive operation.
type UpdateAccessControlRecursiveOptions struct {
	ACL               *string
	BatchSize         *int32
	MaxBatches        *int32
	ContinuationToken *string
	ContinueOnFailure *bool
}

// RemoveAccessControlRecursiveOptions contains the optional parameters when calling the RemoveAccessControlRecursive operation.
type RemoveAccessControlRecursiveOptions struct {
	ACL               *string
	ContinuationToken *string
	BatchSize         *int32
	MaxBatches        *int32
	ContinueOnFailure *bool
}

// CPKInfo contains a group of parameters for client provided encryption key.
type CPKInfo = blob.CPKInfo

// CPKScopeInfo contains a group of parameters for client provided encryption scope.
type CPKScopeInfo = blob.CPKScopeInfo

// AccessConditions identifies container-specific access conditions which you optionally set.
type AccessConditions = exported.PathAccessConditions

// HTTPHeaders contains the HTTP headers for path operations.
type HTTPHeaders = generated.PathHTTPHeaders

// SourceModifiedAccessConditions identifies the source path access conditions.
type SourceModifiedAccessConditions = generated.SourceModifiedAccessConditions

// SetMetadataOptions contains the optional parameters when calling the SetMetadata operation.
type SetMetadataOptions = blob.SetMetadataOptions

// SetHTTPHeadersOptions contains the optional parameters when calling the SetHTTPHeaders operation.
type SetHTTPHeadersOptions = blob.SetHTTPHeadersOptions
