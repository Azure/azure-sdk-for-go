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
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
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
func (b *BlobStore) ClaimOwnership(ctx context.Context, partitionOwnership []azeventhubs.Ownership, options *azeventhubs.ClaimOwnershipOptions) ([]azeventhubs.Ownership, error) {
	var ownerships []azeventhubs.Ownership

	// TODO: in parallel?
	for _, po := range partitionOwnership {
		blobName, err := nameForOwnershipBlob(po)

		if err != nil {
			return nil, err
		}
		lastModified, etag, err := b.setMetadata(ctx, blobName, newOwnershipBlobMetadata(po), po.ETag)

		if err != nil {
			if bloberror.HasCode(err, bloberror.ConditionNotMet) {
				// we can fail to claim ownership and that's okay - it's expected that clients will
				// attempt to claim with whatever state they hold locally. If they fail it just means
				// someone else claimed ownership before them.
				continue
			}

			return nil, err
		}

		newOwnership := po
		newOwnership.ETag = &etag
		newOwnership.LastModifiedTime = *lastModified

		ownerships = append(ownerships, newOwnership)
	}

	return ownerships, nil
}

// ListCheckpoints lists all the available checkpoints.
func (b *BlobStore) ListCheckpoints(ctx context.Context, fullyQualifiedNamespace string, eventHubName string, consumerGroup string, options *azeventhubs.ListCheckpointsOptions) ([]azeventhubs.Checkpoint, error) {
	prefix, err := prefixForCheckpointBlobs(azeventhubs.Checkpoint{
		FullyQualifiedNamespace: fullyQualifiedNamespace,
		EventHubName:            eventHubName,
		ConsumerGroup:           consumerGroup,
	})

	if err != nil {
		return nil, err
	}

	pager := b.cc.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Prefix: &prefix,
		Include: container.ListBlobsInclude{
			Metadata: true,
		},
	})

	var checkpoints []azeventhubs.Checkpoint

	for pager.More() {
		resp, err := pager.NextPage(ctx)

		if err != nil {
			return nil, err
		}

		for _, blob := range resp.Segment.BlobItems {
			partitionID := partitionIDRegexp.FindString(*blob.Name)

			cp := azeventhubs.Checkpoint{
				FullyQualifiedNamespace: fullyQualifiedNamespace,
				EventHubName:            eventHubName,
				ConsumerGroup:           consumerGroup,
				PartitionID:             partitionID,
			}

			if err := updateCheckpoint(blob.Metadata, &cp); err != nil {
				return nil, err
			}

			checkpoints = append(checkpoints, cp)
		}
	}

	return checkpoints, nil
}

var partitionIDRegexp = regexp.MustCompile("[^/]+?$")

// ListOwnership lists all ownerships.
func (b *BlobStore) ListOwnership(ctx context.Context, fullyQualifiedNamespace string, eventHubName string, consumerGroup string, options *azeventhubs.ListOwnershipOptions) ([]azeventhubs.Ownership, error) {
	prefix, err := prefixForOwnershipBlobs(azeventhubs.Ownership{
		FullyQualifiedNamespace: fullyQualifiedNamespace,
		EventHubName:            eventHubName,
		ConsumerGroup:           consumerGroup,
		// ignore partition ID as this is wildcarded.
	})

	if err != nil {
		return nil, err
	}

	pager := b.cc.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Prefix: &prefix,
		Include: container.ListBlobsInclude{
			Metadata: true,
		},
	})

	var ownerships []azeventhubs.Ownership

	for pager.More() {
		resp, err := pager.NextPage(ctx)

		if err != nil {
			return nil, err
		}

		for _, blob := range resp.Segment.BlobItems {
			partitionID := partitionIDRegexp.FindString(*blob.Name)

			o := azeventhubs.Ownership{
				FullyQualifiedNamespace: fullyQualifiedNamespace,
				EventHubName:            eventHubName,
				ConsumerGroup:           consumerGroup,
				PartitionID:             partitionID,
			}

			if err := updateOwnership(blob, &o); err != nil {
				return nil, err
			}

			ownerships = append(ownerships, o)
		}
	}

	return ownerships, nil
}

// UpdateCheckpoint updates a specific checkpoint with a sequence and offset.
func (b *BlobStore) UpdateCheckpoint(ctx context.Context, checkpoint azeventhubs.Checkpoint, options *azeventhubs.UpdateCheckpointOptions) error {
	blobName, err := nameForCheckpointBlob(checkpoint)

	if err != nil {
		return err
	}

	_, _, err = b.setMetadata(ctx, blobName, newCheckpointBlobMetadata(checkpoint), nil)
	return err
}

func (b *BlobStore) setMetadata(ctx context.Context, blobName string, blobMetadata map[string]*string, etag *azcore.ETag) (*time.Time, azcore.ETag, error) {
	blobClient := b.cc.NewBlockBlobClient(blobName)

	if etag != nil {
		setMetadataResp, err := blobClient.SetMetadata(ctx, blobMetadata, &blob.SetMetadataOptions{
			AccessConditions: &blob.AccessConditions{
				ModifiedAccessConditions: &blob.ModifiedAccessConditions{
					IfMatch: etag,
				},
			},
		})

		if err != nil {
			return nil, "", err
		}

		return setMetadataResp.LastModified, *setMetadataResp.ETag, nil
	} else {
		setMetadataResp, err := blobClient.SetMetadata(ctx, blobMetadata, nil)

		if err == nil {
			return setMetadataResp.LastModified, *setMetadataResp.ETag, nil
		}

		if !bloberror.HasCode(err, bloberror.BlobNotFound) {
			return nil, "", err
		}

		// in JS they check to see if the error is BlobNotFound. If it is, then they
		// do a full upload of a blob instead.
		uploadResp, err := blobClient.Upload(ctx, streaming.NopCloser(bytes.NewReader([]byte{})), &blockblob.UploadOptions{
			Metadata: blobMetadata,
		})

		if err != nil {
			return nil, "", err
		}

		return uploadResp.LastModified, *uploadResp.ETag, nil
	}
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

	// ownership : fully-qualified-namespace/event-hub-name/consumer-group/ownership/
	return fmt.Sprintf("%s/%s/%s/ownership/", a.FullyQualifiedNamespace, a.EventHubName, a.ConsumerGroup), nil
}

func updateCheckpoint(metadata map[string]*string, destCheckpoint *azeventhubs.Checkpoint) error {
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

func updateOwnership(b *container.BlobItem, destOwnership *azeventhubs.Ownership) error {
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
