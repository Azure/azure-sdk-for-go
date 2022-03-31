//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package crypto

// KeyOperationResult - The key operation result.
type KeyOperationResult struct {
	// Additional data to authenticate but not encrypt/decrypt when using authenticated crypto algorithms.
	// READ-ONLY
	AuthData []byte `json:"aad,omitempty" azure:"ro"`

	// The tag to authenticate when performing decryption with an authenticated algorithm.
	// READ-ONLY
	AuthTag []byte `json:"tag,omitempty" azure:"ro"`

	// READ-ONLY
	IV []byte `json:"iv,omitempty" azure:"ro"`

	// READ-ONLY; Key identifier
	KeyID *string `json:"kid,omitempty" azure:"ro"`

	// READ-ONLY
	Result []byte `json:"value,omitempty" azure:"ro"`
}
