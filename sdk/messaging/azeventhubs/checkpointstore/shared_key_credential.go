// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package checkpointstore

import "github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/blob"

// SharedKeyCredential contains an Azure Storage account's name and its primary or secondary key.
type SharedKeyCredential struct {
	key *blob.SharedKeyCredential
}

// NewSharedKeyCredential creates an immutable SharedKeyCredential.
func NewSharedKeyCredential(accountName string, accountKey string) (*SharedKeyCredential, error) {
	key, err := blob.NewSharedKeyCredential(accountName, accountKey)

	if err != nil {
		return nil, err
	}

	return &SharedKeyCredential{
		key: key,
	}, nil
}

// SetAccountKey replaces the existing account key with the specified account key.
func (c *SharedKeyCredential) SetAccountKey(accountKey string) error {
	return c.key.SetAccountKey(accountKey)
}
