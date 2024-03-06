// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package checkpoints

import (
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/stretchr/testify/require"
)

func TestBlobStore_copyOwnershipPropsFromBlob(t *testing.T) {
	t.Run("MetadataWorkaround", func(t *testing.T) {
		now := time.Now()
		blobItem := container.BlobItem{
			Properties: &container.BlobProperties{
				ETag:         to.Ptr(azcore.ETag([]byte{1, 2, 3})),
				LastModified: &now,
			},
		}
		ownership := &azeventhubs.Ownership{}
		err := copyOwnershipPropsFromBlob(&blobItem, ownership)
		require.NoError(t, err)

		// this is the workaround - if the metadata dictionary is empty then we
		// just give you back an empty owner ID
		require.Empty(t, ownership.OwnerID)
		require.Equal(t, ownership.ETag, to.Ptr(azcore.ETag([]byte{1, 2, 3})))
		require.Equal(t, now, ownership.LastModifiedTime)
	})

	t.Run("WithMetadataAndOwnerID", func(t *testing.T) {
		now := time.Now()
		blobItem := container.BlobItem{
			Properties: &container.BlobProperties{
				ETag:         to.Ptr(azcore.ETag([]byte{1, 2, 3})),
				LastModified: &now,
			},
			Metadata: map[string]*string{
				"ownerid": to.Ptr("owner id"),
			},
		}
		ownership := &azeventhubs.Ownership{}
		err := copyOwnershipPropsFromBlob(&blobItem, ownership)
		require.NoError(t, err)

		require.Equal(t, "owner id", ownership.OwnerID)
		require.Equal(t, ownership.ETag, to.Ptr(azcore.ETag([]byte{1, 2, 3})))
		require.Equal(t, now, ownership.LastModifiedTime)
	})

	t.Run("WithMetadataNilOwnerID", func(t *testing.T) {
		now := time.Now()
		blobItem := container.BlobItem{
			Properties: &container.BlobProperties{
				ETag:         to.Ptr(azcore.ETag([]byte{1, 2, 3})),
				LastModified: &now,
			},
			Metadata: map[string]*string{
				// In the future this is what I'd expect to see.
				"ownerid": nil,
			},
		}
		ownership := &azeventhubs.Ownership{}
		err := copyOwnershipPropsFromBlob(&blobItem, ownership)
		require.NoError(t, err)

		require.Empty(t, ownership.OwnerID)
		require.Equal(t, ownership.ETag, to.Ptr(azcore.ETag([]byte{1, 2, 3})))
		require.Equal(t, now, ownership.LastModifiedTime)
	})

	t.Run("WithMetadataNoOwnerIDFails", func(t *testing.T) {
		now := time.Now()
		blobItem := container.BlobItem{
			Properties: &container.BlobProperties{
				ETag:         to.Ptr(azcore.ETag([]byte{1, 2, 3})),
				LastModified: &now,
			},
			Metadata: map[string]*string{}, // having metadata but no ownerid is incorrectly formed
		}
		ownership := &azeventhubs.Ownership{}
		err := copyOwnershipPropsFromBlob(&blobItem, ownership)
		require.EqualError(t, err, "ownerid is missing from metadata")
	})

}
