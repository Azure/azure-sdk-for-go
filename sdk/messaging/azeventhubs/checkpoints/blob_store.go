// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package checkpoints

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/blob"
)

// BlobStore is a CheckpointStore implementation that uses Azure Blob storage.
type BlobStore struct {
	cc *blob.ContainerClient
}

// BlobStoreOptions contains optional parameters for the New, NewFromConnectionString and NewWithSharedKey
// functions
type BlobStoreOptions struct {
	azcore.ClientOptions
}

// NewBlobStore creates a checkpoint store that stores ownership and checkpoints in
// Azure Blob storage, using a container URL and a TokenCredential.
// NOTE: the container must exist before the checkpoint store can be used.
func NewBlobStore(containerURL string, cred azcore.TokenCredential, options *BlobStoreOptions) (*BlobStore, error) {
	cc, err := blob.NewContainerClient(containerURL, cred, toContainerClientOptions(options))

	if err != nil {
		return nil, err
	}

	return &BlobStore{
		cc: cc,
	}, nil
}

// NewBlobStoreFromConnectionString creates a checkpoint store that stores
// ownership and checkpoints in Azure Blob storage, using a storage account
// connection string.
// NOTE: the container must exist before the checkpoint store can be used.
func NewBlobStoreFromConnectionString(connectionString string, containerName string, options *BlobStoreOptions) (azeventhubs.CheckpointStore, error) {
	cc, err := blob.NewContainerClientFromConnectionString(connectionString, containerName, toContainerClientOptions(options))

	if err != nil {
		return nil, err
	}

	return &BlobStore{
		cc: cc,
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

		var expectedETag *string

		if po.ETag != "" {
			expectedETag = &po.ETag
		}

		lastModified, etag, err := b.setMetadata(ctx, blobName, newOwnershipBlobMetadata(po), expectedETag)

		if err != nil {
			var storageErr *blob.StorageError

			// we can fail to claim ownership and that's okay - it's expected that clients will
			// attempt to claim with whatever state they hold locally. If they fail it just means
			// someone else claimed ownership before them.
			if errors.As(err, &storageErr) && storageErr.ErrorCode == blob.StorageErrorCodeConditionNotMet {
				continue
			}

			return nil, err
		}

		newOwnership := po
		newOwnership.ETag = etag
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

	pager := b.cc.ListBlobsFlat(&blob.ContainerListBlobsFlatOptions{
		Prefix: &prefix,
		Include: []blob.ListBlobsIncludeItem{
			blob.ListBlobsIncludeItemMetadata,
		},
	})

	var checkpoints []azeventhubs.Checkpoint

	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

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

	if pager.Err() != nil {
		return nil, pager.Err()
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

	pager := b.cc.ListBlobsFlat(&blob.ContainerListBlobsFlatOptions{
		Prefix: &prefix,
		Include: []blob.ListBlobsIncludeItem{
			blob.ListBlobsIncludeItemMetadata,
		},
	})

	var ownerships []azeventhubs.Ownership

	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

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

	if pager.Err() != nil {
		return nil, pager.Err()
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

func isBlobNotFoundError(err error) bool {
	var storageErr *blob.StorageError
	return errors.As(err, &storageErr) && storageErr.StatusCode() == http.StatusNotFound
}

func (b *BlobStore) setMetadata(ctx context.Context, blobName string, blobMetadata map[string]string, etag *string) (*time.Time, string, error) {
	blobClient, err := b.cc.NewBlockBlobClient(blobName)

	if err != nil {
		return nil, "", err
	}

	if etag != nil {
		setMetadataResp, err := blobClient.SetMetadata(ctx, blobMetadata, &blob.BlobSetMetadataOptions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: etag,
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

		if !isBlobNotFoundError(err) {
			return nil, "", err
		}

		// in JS they check to see if the error is BlobNotFound. If it is, then they
		// do a full upload of a blob instead.
		uploadResp, err := blobClient.Upload(ctx, streaming.NopCloser(bytes.NewReader([]byte{})), &blob.BlockBlobUploadOptions{
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

func newCheckpointBlobMetadata(cpd azeventhubs.Checkpoint) map[string]string {
	m := map[string]string{}

	if cpd.SequenceNumber != nil {
		m["sequencenumber"] = strconv.FormatInt(*cpd.SequenceNumber, 10)
	}

	if cpd.Offset != nil {
		m["offset"] = strconv.FormatInt(*cpd.Offset, 10)
	}

	return m
}

func updateOwnership(b *blob.BlobItemInternal, destOwnership *azeventhubs.Ownership) error {
	if b == nil || b.Metadata == nil || b.Properties == nil {
		return fmt.Errorf("no ownership metadata for blob")
	}

	ownerID, ok := b.Metadata["ownerid"]

	if !ok || ownerID == nil {
		return errors.New("ownerid is missing from metadata")
	}

	destOwnership.OwnerID = *ownerID
	destOwnership.LastModifiedTime = *b.Properties.LastModified
	destOwnership.ETag = *b.Properties.Etag
	return nil
}

func newOwnershipBlobMetadata(od azeventhubs.Ownership) map[string]string {
	return map[string]string{
		"ownerid": od.OwnerID,
	}
}

func toContainerClientOptions(opts *BlobStoreOptions) *blob.ClientOptions {
	if opts == nil {
		return nil
	}

	return &blob.ClientOptions{
		Logging:          opts.Logging,
		Retry:            opts.Retry,
		Telemetry:        opts.Telemetry,
		Transport:        opts.Transport,
		PerCallPolicies:  opts.PerCallPolicies,
		PerRetryPolicies: opts.PerRetryPolicies,
	}
}
