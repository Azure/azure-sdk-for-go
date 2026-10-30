// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
)

// This file exposes Data Locality internals to the tests in package blob_test. It is compiled
// only into tests, so none of it is part of the public API.

// GetLayoutOptionsForTest is the options type of the unexported Get Blob Layout pager.
type GetLayoutOptionsForTest = getLayoutOptions

// GetLayoutPagerForTest returns the unexported Get Blob Layout pager of c.
func GetLayoutPagerForTest(c *Client, options *GetLayoutOptionsForTest) *runtime.Pager[generated.BlobClientGetLayoutResponse] {
	return c.getLayoutPager(options)
}

// WithLayoutEndpointForTest routes o to endpoint, as a managed download routes a range read to the
// endpoint its layout serves it from.
func WithLayoutEndpointForTest(o *DownloadStreamOptions, endpoint string) *DownloadStreamOptions {
	o.layoutEndpoint = endpoint
	return o
}
