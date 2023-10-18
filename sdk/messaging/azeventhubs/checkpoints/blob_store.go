// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package checkpoints

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

// BlobStore is a CheckpointStore implementation that uses Azure Blob storage.
type BlobStore struct {
	cc *container.Client
}

// BlobStoreOptions contains optional parameters for the New, NewFromConnectionString and NewWithSharedKey
// functions
type BlobStoreOptions struct {
	azcore.ClientOptions
}

// NewBlobStore creates a checkpoint store that stores ownership and checkpoints in
// Azure Blob storage.
// NOTE: the container must exist before the checkpoint store can be used.
func NewBlobStore(containerClient *container.Client, options *BlobStoreOptions) (*BlobStore, error) {
	return &BlobStore{
		cc: containerClient,
	}, nil
}

// ClaimOwnership attempts to claim ownership of the partitions in partitionOwnership and returns
// the actual partitions that were claimed.
//
// If we fail to claim ownership because of another update then it will be omitted from the
// returned slice of [Ownership]'s. It is not considered an error.
func (b *BlobStore) ClaimOwnership(ctx context.Context, partitionOwnership []azeventhubs.Ownership, options *azeventhubs.ClaimOwnershipOptions) ([]azeventhubs.Ownership, error) {
	var ownerships []azeventhubs.Ownership

	// TODO: in parallel?
	for _, po := range partitionOwnership {
		legacyBlobName, err := nameForOwnershipBlob(po)

		if err != nil {
			return nil, err
		}

		correctBlobName := strings.ToLower(legacyBlobName)

		lastModified, etag, err := b.setOwnershipMetadata(ctx, correctBlobName, po)

		if err != nil {
			if bloberror.HasCode(err,
				bloberror.ConditionNotMet,     // updated before we could update it
				bloberror.BlobAlreadyExists) { // created before we could create it

				log.Writef(azeventhubs.EventConsumer, "[%s] skipping %s because: %s", po.OwnerID, po.PartitionID, err)
				continue
			}

			return nil, err
		}

		newOwnership := po
		newOwnership.ETag = &etag
		newOwnership.LastModifiedTime = *lastModified

		ownerships = append(ownerships, newOwnership)

		if legacyBlobName != correctBlobName {
			// Kill the legacy blob since it was written with mixed-case, which is incorrect.
			bc := b.cc.NewBlobClient(legacyBlobName)

			// Best effort to delete the old blob if it exists.
			// If we fail here it doesn't break anything - future List* operations will just not show this older blob.
			if _, err = bc.Delete(ctx, nil); err != nil {
				log.Writef(azeventhubs.EventConsumer, "[%s] failed to delete legacy blob (%s): %s", po.OwnerID, legacyBlobName, err)
			}
		}
	}

	return ownerships, nil
}

// ListCheckpoints lists all the available checkpoints.
func (b *BlobStore) ListCheckpoints(ctx context.Context, fullyQualifiedNamespace string, eventHubName string, consumerGroup string, options *azeventhubs.ListCheckpointsOptions) ([]azeventhubs.Checkpoint, error) {
	legacyPrefix, err := prefixForCheckpointBlobs(azeventhubs.Checkpoint{
		FullyQualifiedNamespace: fullyQualifiedNamespace,
		EventHubName:            eventHubName,
		ConsumerGroup:           consumerGroup,
	})

	if err != nil {
		return nil, err
	}

	correctPrefix := strings.ToLower(legacyPrefix)

	checkpointsMap := map[string]azeventhubs.Checkpoint{}

	addCheckpointsToMap := func(prefix string, isLegacy bool) error {
		pager := b.cc.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
			Prefix: &prefix,
			Include: container.ListBlobsInclude{
				Metadata: true,
			},
		})

		for pager.More() {
			resp, err := pager.NextPage(ctx)

			if err != nil {
				return err
			}

			for _, blob := range resp.Segment.BlobItems {
				partitionID := partitionIDRegexp.FindString(*blob.Name)

				cp := azeventhubs.Checkpoint{
					FullyQualifiedNamespace: fullyQualifiedNamespace,
					EventHubName:            eventHubName,
					ConsumerGroup:           consumerGroup,
					PartitionID:             partitionID,
				}

				if err := copyCheckpointPropsFromMetadata(blob.Metadata, &cp); err != nil {
					return err
				}

				if !isLegacy {
					cp.ConsumerGroup = strings.ToLower(cp.ConsumerGroup)
					cp.EventHubName = strings.ToLower(cp.EventHubName)
					cp.FullyQualifiedNamespace = strings.ToLower(cp.FullyQualifiedNamespace)
					cp.PartitionID = strings.ToLower(cp.PartitionID)
				}

				checkpointsMap[strings.ToLower(*blob.Name)] = cp
			}
		}

		return nil
	}

	if correctPrefix != legacyPrefix {
		// we've got a casing difference - we need to make sure we look at both possible
		// casings because of our previous bug where we didn't ToLower() the blob name by
		// default.
		if err := addCheckpointsToMap(legacyPrefix, true); err != nil {
			return nil, err
		}
	}

	if err := addCheckpointsToMap(correctPrefix, false); err != nil {
		return nil, err
	}

	var checkpoints []azeventhubs.Checkpoint

	for _, cp := range checkpointsMap {
		checkpoints = append(checkpoints, cp)
	}

	return checkpoints, nil
}

var partitionIDRegexp = regexp.MustCompile("[^/]+?$")

// ListOwnership lists all ownerships.
func (b *BlobStore) ListOwnership(ctx context.Context, fullyQualifiedNamespace string, eventHubName string, consumerGroup string, options *azeventhubs.ListOwnershipOptions) ([]azeventhubs.Ownership, error) {
	legacyPrefix, err := prefixForOwnershipBlobs(azeventhubs.Ownership{
		FullyQualifiedNamespace: fullyQualifiedNamespace,
		EventHubName:            eventHubName,
		ConsumerGroup:           consumerGroup,
		// ignore partition ID as this is wildcarded.
	})

	if err != nil {
		return nil, err
	}

	blobsMap := map[string]azeventhubs.Ownership{}

	addOwnershipsToMap := func(prefix string, isLegacy bool) error {
		pager := b.cc.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
			Prefix: &prefix,
			Include: container.ListBlobsInclude{
				Metadata: true,
			},
		})

		for pager.More() {
			resp, err := pager.NextPage(ctx)

			if err != nil {
				return err
			}

			for _, blob := range resp.Segment.BlobItems {
				partitionID := partitionIDRegexp.FindString(*blob.Name)

				o := azeventhubs.Ownership{
					FullyQualifiedNamespace: fullyQualifiedNamespace,
					EventHubName:            eventHubName,
					ConsumerGroup:           consumerGroup,
					PartitionID:             partitionID,
				}

				if err := copyOwnershipPropsFromBlob(blob, &o); err != nil {
					return err
				}

				if isLegacy {
					// we clear the etag if we're loading up the "old" ownership blobs, where we incorrectly
					// wrote them using mixed-case, which was a bug. This makes it so any consumers of this
					// ownership record understand they need to create a new blob, not update it.
					o.ETag = nil
				} else {
					o.ConsumerGroup = strings.ToLower(o.ConsumerGroup)
					o.EventHubName = strings.ToLower(o.EventHubName)
					o.FullyQualifiedNamespace = strings.ToLower(o.FullyQualifiedNamespace)
					o.PartitionID = strings.ToLower(o.PartitionID)
				}

				blobsMap[strings.ToLower(*blob.Name)] = o
			}
		}

		return nil
	}

	correctPrefix := strings.ToLower(legacyPrefix)

	if legacyPrefix != correctPrefix {
		// we've got a casing difference - we need to make sure we look at both possible
		// casings because of our previous bug where we didn't ToLower() the blob name by
		// default.
		if err := addOwnershipsToMap(legacyPrefix, true); err != nil {
			return nil, err
		}
	}

	if err := addOwnershipsToMap(correctPrefix, false); err != nil {
		return nil, err
	}

	var ownerships []azeventhubs.Ownership

	for _, o := range blobsMap {
		ownerships = append(ownerships, o)
	}

	return ownerships, nil
}

// SetCheckpoint updates a specific checkpoint with a sequence and offset.
//
// NOTE: This function doesn't attempt to prevent simultaneous checkpoint updates - ownership is assumed.
func (b *BlobStore) SetCheckpoint(ctx context.Context, checkpoint azeventhubs.Checkpoint, options *azeventhubs.SetCheckpointOptions) error {
	legacyBlobName, err := nameForCheckpointBlob(checkpoint)

	if err != nil {
		return err
	}

	correctBlobName := strings.ToLower(legacyBlobName)
	_, _, err = b.setCheckpointMetadata(ctx, correctBlobName, checkpoint)

	if legacyBlobName != correctBlobName {
		// delete the old blob
		bc := b.cc.NewBlobClient(legacyBlobName)

		if _, err := bc.Delete(ctx, nil); err != nil {
			log.Writef(azeventhubs.EventConsumer, "Failed to delete legacy blob (%s): %s", legacyBlobName, err)
		}
	}

	return err
}

func (b *BlobStore) setOwnershipMetadata(ctx context.Context, blobName string, ownership azeventhubs.Ownership) (*time.Time, azcore.ETag, error) {
	blobMetadata := newOwnershipBlobMetadata(ownership)
	blobClient := b.cc.NewBlockBlobClient(blobName)

	if ownership.ETag != nil {
		log.Writef(azeventhubs.EventConsumer, "[%s] claiming ownership for %s with etag %s", ownership.OwnerID, ownership.PartitionID, string(*ownership.ETag))
		setMetadataResp, err := blobClient.SetMetadata(ctx, blobMetadata, &blob.SetMetadataOptions{
			AccessConditions: &blob.AccessConditions{
				ModifiedAccessConditions: &blob.ModifiedAccessConditions{
					IfMatch: ownership.ETag,
				},
			},
		})

		if err != nil {
			return nil, "", err
		}

		return setMetadataResp.LastModified, *setMetadataResp.ETag, nil
	}

	log.Writef(azeventhubs.EventConsumer, "[%s] claiming ownership for %s with NO etags", ownership.PartitionID, ownership.OwnerID)
	uploadResp, err := blobClient.Upload(ctx, streaming.NopCloser(bytes.NewReader([]byte{})), &blockblob.UploadOptions{
		Metadata: blobMetadata,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: to.Ptr(azcore.ETag("*")),
			},
		},
	})

	if err != nil {
		return nil, "", err
	}

	return uploadResp.LastModified, *uploadResp.ETag, nil
}

// setCheckpointMetadata sets the metadata for a checkpoint, falling back to creating
// the blob if it doesn't already exist.
//
// NOTE: unlike [setOwnershipMetadata] this function doesn't attempt to prevent simultaneous
// checkpoint updates - ownership is assumed.
func (b *BlobStore) setCheckpointMetadata(ctx context.Context, blobName string, checkpoint azeventhubs.Checkpoint) (*time.Time, azcore.ETag, error) {
	blobMetadata := newCheckpointBlobMetadata(checkpoint)
	blobClient := b.cc.NewBlockBlobClient(blobName)

	setMetadataResp, err := blobClient.SetMetadata(ctx, blobMetadata, nil)

	if err == nil {
		return setMetadataResp.LastModified, *setMetadataResp.ETag, nil
	}

	if !bloberror.HasCode(err, bloberror.BlobNotFound) {
		return nil, "", err
	}

	uploadResp, err := blobClient.Upload(ctx, streaming.NopCloser(bytes.NewReader([]byte{})), &blockblob.UploadOptions{
		Metadata: blobMetadata,
	})

	if err != nil {
		return nil, "", err
	}

	return uploadResp.LastModified, *uploadResp.ETag, nil
}

func nameForCheckpointBlob(a azeventhubs.Checkpoint) (string, error) {
	if a.FullyQualifiedNamespace == "" || a.EventHubName == "" || a.ConsumerGroup == "" || a.PartitionID == "" {
		return "", errors.New("missing fields for blob name")
	}

	// checkpoint: fully-qualified-namespace/event-hub-name/consumer-group/checkpoint/partition-id
	return fmt.Sprintf("%s/%s/%s/checkpoint/%s", a.FullyQualifiedNamespace, a.EventHubName, a.ConsumerGroup, a.PartitionID), nil
}

func prefixForCheckpointBlobs(a azeventhubs.Checkpoint) (string, error) {
	if a.FullyQualifiedNamespace == "" || a.EventHubName == "" || a.ConsumerGroup == "" {
		return "", errors.New("missing fields for blob prefix")
	}

	// checkpoint: fully-qualified-namespace/event-hub-name/consumer-group/checkpoint/
	return fmt.Sprintf("%s/%s/%s/checkpoint/", a.FullyQualifiedNamespace, a.EventHubName, a.ConsumerGroup), nil
}

func nameForOwnershipBlob(a azeventhubs.Ownership) (string, error) {
	if a.FullyQualifiedNamespace == "" || a.EventHubName == "" || a.ConsumerGroup == "" || a.PartitionID == "" {
		return "", errors.New("missing fields for blob name")
	}

	// ownership : fully-qualified-namespace/event-hub-name/consumer-group/ownership/partition-id
	return fmt.Sprintf("%s/%s/%s/ownership/%s", a.FullyQualifiedNamespace, a.EventHubName, a.ConsumerGroup, a.PartitionID), nil
}

func prefixForOwnershipBlobs(a azeventhubs.Ownership) (string, error) {
	if a.FullyQualifiedNamespace == "" || a.EventHubName == "" || a.ConsumerGroup == "" {
		return "", errors.New("missing fields for blob prefix")
	}

	return fmt.Sprintf("%s/%s/%s/ownership/", a.FullyQualifiedNamespace, a.EventHubName, a.ConsumerGroup), nil
}

func copyCheckpointPropsFromMetadata(metadata map[string]*string, destCheckpoint *azeventhubs.Checkpoint) error {
	if metadata == nil {
		return fmt.Errorf("no checkpoint metadata for blob")
	}

	sequenceNumberStr, ok := metadata["sequencenumber"]

	if !ok || sequenceNumberStr == nil {
		return errors.New("sequencenumber is missing from metadata")
	}

	sequenceNumber, err := strconv.ParseInt(*sequenceNumberStr, 10, 64)

	if err != nil {
		return fmt.Errorf("sequencenumber could not be parsed as an int64: %s", err.Error())
	}

	offsetStr, ok := metadata["offset"]

	if !ok || offsetStr == nil {
		return errors.New("offset is missing from metadata")
	}

	offset, err := strconv.ParseInt(*offsetStr, 10, 64)

	if err != nil {
		return fmt.Errorf("offset could not be parsed as an int64: %s", err.Error())
	}

	destCheckpoint.Offset = &offset
	destCheckpoint.SequenceNumber = &sequenceNumber
	return nil
}

func newCheckpointBlobMetadata(cpd azeventhubs.Checkpoint) map[string]*string {
	m := map[string]*string{}

	if cpd.SequenceNumber != nil {
		m["sequencenumber"] = to.Ptr(strconv.FormatInt(*cpd.SequenceNumber, 10))
	}

	if cpd.Offset != nil {
		m["offset"] = to.Ptr(strconv.FormatInt(*cpd.Offset, 10))
	}

	return m
}

func copyOwnershipPropsFromBlob(b *container.BlobItem, destOwnership *azeventhubs.Ownership) error {
	if b == nil || b.Metadata == nil || b.Properties == nil {
		return fmt.Errorf("no ownership metadata for blob")
	}

	ownerID, ok := b.Metadata["ownerid"]

	if !ok || ownerID == nil {
		return errors.New("ownerid is missing from metadata")
	}

	destOwnership.OwnerID = *ownerID
	destOwnership.LastModifiedTime = *b.Properties.LastModified
	destOwnership.ETag = b.Properties.ETag
	return nil
}

func newOwnershipBlobMetadata(od azeventhubs.Ownership) map[string]*string {
	return map[string]*string{
		"ownerid": &od.OwnerID,
	}
}
