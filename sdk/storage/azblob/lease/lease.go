//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lease

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
)

type Client base.Client[generated.BlobClient]

func (c *Client) generated() *generated.BlobClient {
	return base.InnerClient((*base.Client[generated.BlobClient])(c))
}

func (c *Client) Acquire(ctx context.Context, options *AcquireOptions) (AcquireResponse, error) {
	options = shared.CopyOptions(options)
	if options.LeaseID == nil {
		id, err := uuid.New()
		if err != nil {
			return AcquireResponse{}, err
		}
		options.LeaseID = to.Ptr(id.String())
	}

	resp, err := c.generated().AcquireLease(ctx, &generated.BlobClientAcquireLeaseOptions{
		Duration:        options.Duration,
		ProposedLeaseID: options.LeaseID,
	}, options.ModifiedAccessConditions)
	if err != nil {
		return AcquireResponse{}, err
	}
	return AcquireResponse{AccessConditions: AccessConditions{LeaseID: resp.LeaseID}}, nil
}
