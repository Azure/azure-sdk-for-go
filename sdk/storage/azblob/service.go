// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azinternal "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/azblob"
)

const (
	scope = "https://storage.azure.com/.default"
)

type ServiceClient struct {
	s azinternal.Service
	u *url.URL
	p azcore.Pipeline
}

func NewServiceClient(endpoint string, cred azcore.Credential, options azcore.PipelineOptions) (*ServiceClient, error) {
	p := azcore.NewPipeline(options.HTTPClient,
		azcore.NewTelemetryPolicy(options.Telemetry),
		azcore.NewUniqueRequestIDPolicy(),
		azcore.NewRetryPolicy(options.Retry),
		cred.AuthenticationPolicy(azcore.AuthenticationPolicyOptions{Options: azcore.TokenRequestOptions{Scopes: []string{scope}}}),
		azcore.NewRequestLogPolicy(options.LogOptions))
	return NewServiceClientWithPipeline(endpoint, p)
}

func NewServiceClientWithPipeline(endpoint string, p azcore.Pipeline) (*ServiceClient, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	return &ServiceClient{u: u, p: p}, nil
}

// ListContainers the List Containers Segment operation returns a list of the containers under the specified
// account
//
// prefix is filters the results to return only containers whose name begins with the specified prefix. marker is a
// string value that identifies the portion of the list of containers to be returned with the next listing operation.
// The operation returns the NextMarker value within the response body if the listing operation did not return all
// containers remaining to be listed with the current page. The NextMarker value can be used as the value for the
// marker parameter in a subsequent call to request the next page of list items. The marker value is opaque to the
// client. maxresults is specifies the maximum number of containers to return. If the request does not specify
// maxresults, or specifies a value greater than 5000, the server will return up to 5000 items. Note that if the
// listing operation crosses a partition boundary, then the service will return a continuation token for retrieving the
// remainder of the results. For this reason, it is possible that the service will return fewer results than specified
// by maxresults, or than the default of 5000. include is include this parameter to specify that the container's
// metadata be returned as part of the response body. timeout is the timeout parameter is expressed in seconds. For
// more information, see <a
// href="https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/setting-timeouts-for-blob-service-operations">Setting
// Timeouts for Blob Service Operations.</a> requestID is provides a client-generated, opaque value with a 1 KB
// character limit that is recorded in the analytics logs when storage analytics logging is enabled.
func (c *ServiceClient) ListContainers(options *ListContainersOptions) *ListContainersIterator {
	if options == nil {
		options = &ListContainersOptions{}
	}
	return &ListContainersIterator{
		client: c,
		op:     options,
	}
}
