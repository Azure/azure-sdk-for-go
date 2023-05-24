//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/path"
	"time"
)

// CreateOptions contains the optional parameters when calling the Create operation. dfs endpoint
type CreateOptions struct {
	AccessConditions *AccessConditions
	Metadata         map[string]*string
	CPKInfo          *CPKInfo
	HTTPHeaders      *HTTPHeaders
	//PathExpiryOptions              *ExpiryOptions
	ExpiresOn       *time.Time
	LeaseDuration   *time.Duration
	ProposedLeaseID *string
	Permissions     *string
	Umask           *string
	Owner           *string
	Group           *string
	ACL             *string
}

// DeleteOptions contains the optional parameters when calling the Delete operation. dfs endpoint
type DeleteOptions struct {
	// used to distinguish between dir or file deletion
	//Recursive        *bool
	AccessConditions *AccessConditions
}

type RenameOptions struct {
	SourceModifiedAccessConditions *SourceModifiedAccessConditions
	AccessConditions               *AccessConditions
}

type GetPropertiesOptions struct {
	AccessConditions *AccessConditions
	CPKInfo          *CPKInfo
}

type SetExpiryOptions struct {
	ExpiresOn *time.Time
}

// ===================================== PATH IMPORTS ===========================================

// CPKInfo contains a group of parameters for client provided encryption key.
type CPKInfo = path.CPKInfo

// CPKScopeInfo contains a group of parameters for client provided encryption scope.
type CPKScopeInfo = path.CPKScopeInfo

// AccessConditions identifies container-specific access conditions which you optionally set.
type AccessConditions = path.AccessConditions

// HTTPHeaders contains the HTTP headers for path operations.
type HTTPHeaders = path.HTTPHeaders

// SourceModifiedAccessConditions identifies the source path access conditions.
type SourceModifiedAccessConditions = path.SourceModifiedAccessConditions

type SetAccessControlRecursiveOptions = path.SetAccessControlRecursiveOptions

type SetMetadataOptions = path.SetMetadataOptions

type SetHTTPHeadersOptions = path.SetHTTPHeadersOptions

type RemoveAccessControlRecursiveOptions = path.RemoveAccessControlRecursiveOptions

type UpdateAccessControlRecursiveOptions = path.UpdateAccessControlRecursiveOptions

type SetAccessControlOptions = path.SetAccessControlOptions
