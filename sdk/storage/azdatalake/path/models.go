//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package path

type CreateOptions struct {
	ContinuationToken *string
	Permissions       *string
	Properties        *string
	SourceLeaseID     *string
	Umask             *string
}

type DeleteOptions struct {
}