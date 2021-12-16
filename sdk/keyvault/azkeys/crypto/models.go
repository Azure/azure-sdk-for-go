//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package crypto

// KeyOperationResult - The key operation result.
type KeyOperationResult struct {
	// READ-ONLY
	AdditionalAuthenticatedData []byte `json:"aad,omitempty" azure:"ro"`

	// READ-ONLY
	AuthenticationTag []byte `json:"tag,omitempty" azure:"ro"`

	// READ-ONLY
	IV []byte `json:"iv,omitempty" azure:"ro"`

	// READ-ONLY; Key identifier
	KeyID *string `json:"kid,omitempty" azure:"ro"`

	// READ-ONLY
	Result []byte `json:"value,omitempty" azure:"ro"`
}
