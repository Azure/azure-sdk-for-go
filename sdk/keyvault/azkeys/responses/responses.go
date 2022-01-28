//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package responses

import (
	"context"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/models"
)

// KeyVaultClientCreateKeyResponse contains the response from method KeyVaultClient.CreateKey.
type CreateKey struct {
	models.KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// CreateECKeyResponse contains the response from method Client.CreateECKey.
type CreateECKey struct {
	models.KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// CreateOCTKeyResponse contains the response from method Client.CreateOCTKey.
type CreateOCTKey struct {
	models.KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// CreateRSAKeyResponse contains the response from method Client.CreateRSAKey.
type CreateRSAKey struct {
	models.KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// ListKeysPage contains the current page of results for the Client.ListSecrets operation
type ListKeysPage struct {
	// READ-ONLY; The URL to get the next set of keys.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// READ-ONLY; A response message containing a list of keys in the key vault along with a link to the next page of keys.
	Keys []*models.KeyItem `json:"value,omitempty" azure:"ro"`
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// GetKeyResponse contains the response for the Client.GetResponse method
type GetKey struct {
	models.KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// GetDeletedKeyResponse contains the response from a Client.GetDeletedKey
type GetDeletedKey struct {
	models.DeletedKeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PurgeDeletedKeyResponse contains the response from method Client.PurgeDeletedKey.
type PurgeDeletedKey struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// DeletedKeyResponse contains the response for a Client.BeginDeleteKey operation.
type DeleteKey struct {
	models.DeletedKeyBundle
	// RawResponse holds the underlying HTTP response
	RawResponse *http.Response
}

// BeginDeleteKey contains the response from the Client.BeginDeleteKey method
type BeginDeleteKey struct {
	// Poller contains an initialized WidgetPoller
	Poller *DeleteKeyPoller

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// PollUntilDone will poll the service endpoint until a terminal state is reached or an error occurs
func (b *BeginDeleteKey) PollUntilDone(ctx context.Context, freq time.Duration) (DeleteKey, error) {
	return b.Poller.pollUntilDone(ctx, freq)
}

// BackupKeyResponse contains the response from the Client.BackupKey method
type BackupKey struct {
	// READ-ONLY; The backup blob containing the backed up key.
	Value []byte `json:"value,omitempty" azure:"ro"`

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// RecoverDeletedKeyResponse is the response object for the Client.RecoverDeletedKey operation.
type RecoverDeletedKey struct {
	models.KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// BeginRecoverDeletedKey contains the response of the Client.BeginRecoverDeletedKey operations
type BeginRecoverDeletedKey struct {
	// Poller contains an initialized RecoverDeletedKeyPoller
	Poller *RecoverDeletedKeyPoller

	// RawResponse cotains the underlying HTTP response
	RawResponse *http.Response
}

// PollUntilDone will poll the service endpoint until a terminal state is reached or an error occurs
func (b *BeginRecoverDeletedKey) PollUntilDone(ctx context.Context, freq time.Duration) (RecoverDeletedKey, error) {
	return b.Poller.pollUntilDone(ctx, freq)
}

// UpdateKeyPropertiesResponse contains the response for the Client.UpdateKeyProperties method
type UpdateKeyProperties struct {
	models.KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// ListDeletedKeysPage holds the data for a single page.
type ListDeletedKeysPage struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response

	// READ-ONLY; The URL to get the next set of deleted keys.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// READ-ONLY; A response message containing a list of the deleted keys in the vault along with a link to the next page of deleted keys
	DeletedKeys []*models.DeletedKeyItem `json:"value,omitempty" azure:"ro"`
}

// ListKeyVersionsPage contains the current page from a ListKeyVersionsPager.PageResponse method
type ListKeyVersionsPage struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response

	// READ-ONLY; The URL to get the next set of keys.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// READ-ONLY; A response message containing a list of keys in the key vault along with a link to the next page of keys.
	Keys []models.KeyItem `json:"value,omitempty" azure:"ro"`
}

// RestoreKeyBackupResponse contains the response object for the Client.RestoreKeyBackup operation.
type RestoreKeyBackup struct {
	models.KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// ImportKeyResponse contains the response of the Client.ImportKey method
type ImportKey struct {
	models.KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// GetRandomBytesResponse is the response struct for the Client.GetRandomBytes function.
type GetRandomBytes struct {
	// The bytes encoded as a base64url string.
	Value []byte `json:"value,omitempty"`

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

type RotateKey struct {
	models.KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// GetKeyRotationPolicyResponse contains the response struct for the Client.GetKeyRotationPolicy function
type GetKeyRotationPolicy struct {
	models.KeyRotationPolicy
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// ReleaseKeyResponse contains the response of Client.ReleaseKey
type ReleaseKey struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response

	// READ-ONLY; A signed object containing the released key.
	Value *string `json:"value,omitempty" azure:"ro"`
}

// UpdateKeyRotationPolicyResponse contains the response for the Client.UpdateKeyRotationPolicy function
type UpdateKeyRotationPolicy struct {
	models.KeyRotationPolicy
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}
