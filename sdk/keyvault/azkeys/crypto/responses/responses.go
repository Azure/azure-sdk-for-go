//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package responses

import (
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/crypto/models"
)

type EncryptResponse struct {
	models.KeyOperationResult
}

type DecryptResponse struct {
	models.KeyOperationResult
}

// WrapKeyResponse contains the response for the Client.WrapKey method
type WrapKeyResponse struct {
	models.KeyOperationResult
}

// UnwrapKeyResponse contains the response for the Client.UnwrapKey method
type UnwrapKeyResponse struct {
	models.KeyOperationResult
}

// SignResponse contains the response for the Client.Sign method.
type SignResponse struct {
	models.KeyOperationResult
}

// VerifyResponse contains the response for the Client.Verify method
type VerifyResponse struct {
	// READ-ONLY; True if the signature is verified, otherwise false.
	IsValid *bool `json:"value,omitempty" azure:"ro"`
}
