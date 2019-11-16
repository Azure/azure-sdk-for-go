// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	scope = "https://storage.azure.com/.default"
)

// ServiceVersion enumerates the possible service API versions.
type ServiceVersion string

const (
	// V2019_02_02 specifies the 2019-02-02 service API version.
	V2019_02_02 = "2019-02-02"
)

//ServiceClient provides access to Azure Storage service resources and blob containers.
type ServiceClient struct {
	u *url.URL
	p azcore.Pipeline
	o ServiceClientOptions
}

// ServiceClientOptions contains client configuration options for connecting to blob storage.
type ServiceClientOptions struct {
	azcore.PipelineOptions
	// Version is the service API version to use when making requests.
	// Leaving this empty will default to the latest API version.
	Version ServiceVersion

	// SecondaryEndpoint is the secondary storage endpoint that can be read from if the storage
	// account is enabled for RA-GRS.
	// If this value is set, the secondary endpoint will be used for GET or HEAD requests during retries.
	// If the status of the response from the secondary endpoint is a 404, then subsequent retries for
	// the request will not use the secondary endpoint again, as this indicates that the resource
	// may not have propagated there yet. Otherwise, subsequent retries will alternate back and forth
	// between the primary and secondary endpoint.
	SecondaryEndpoint string
}

func setDefaultClientOptions(s *ServiceClientOptions) *ServiceClientOptions {
	if s == nil {
		s = &ServiceClientOptions{}
	}
	if len(s.Version) == 0 {
		s.Version = V2019_02_02
	}
	return s
}

// NewServiceClient creates a new instance of ServiceClient with the specified values.
// Passing nil for cred will default to anonymous access.
// Pass nil for options to use the default configuration.
func NewServiceClient(endpoint string, cred azcore.Credential, options *ServiceClientOptions) (*ServiceClient, error) {
	if cred == nil {
		cred = azcore.AnonymousCredential()
	}
	options = setDefaultClientOptions(options)
	p := azcore.NewPipeline(options.HTTPTransport,
		azcore.NewTelemetryPolicy(options.Telemetry),
		azcore.NewUniqueRequestIDPolicy(),
		azcore.NewRetryPolicy(options.Retry),
		cred.AuthenticationPolicy(azcore.AuthenticationPolicyOptions{Scopes: []string{scope}}),
		azcore.NewRequestLogPolicy(options.LogOptions))
	return NewServiceClientWithPipeline(endpoint, p, options)
}

// NewServiceClientWithPipeline creates a new instance of ServiceClient with the specified values.
func NewServiceClientWithPipeline(endpoint string, pipeline azcore.Pipeline, options *ServiceClientOptions) (*ServiceClient, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	options = setDefaultClientOptions(options)
	return &ServiceClient{u: u, p: pipeline, o: *options}, nil
}

// ServiceVersion returns the service API version used when making requests.
func (c *ServiceClient) ServiceVersion() ServiceVersion {
	return c.o.Version
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

// listContainersPreparer prepares the ListContainersSegment request.
func (c *ServiceClient) listContainersPreparer(options *ListContainersOptions) *azcore.Request {
	req := c.p.NewRequest(http.MethodGet, *c.u)
	if options != nil {
		if options.Prefix != nil && len(*options.Prefix) > 0 {
			req.SetQueryParam("prefix", *options.Prefix)
		}
		if options.Marker != nil && len(*options.Marker) > 0 {
			req.SetQueryParam("marker", *options.Marker)
		}
		if options.Maxresults != nil {
			req.SetQueryParam("maxresults", strconv.FormatInt(int64(*options.Maxresults), 10))
		}
		if options.Include != ListContainersIncludeNone {
			req.SetQueryParam("include", string(options.Include))
		}
		if options.Timeout != nil {
			req.SetQueryParam("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
		}
		if options.RequestID != nil {
			req.Header.Set("x-ms-client-request-id", *options.RequestID)
		}
	}
	req.SetQueryParam("comp", "list")
	req.Header.Set("x-ms-version", string(c.ServiceVersion()))
	return req
}

// listContainersResponder handles the response to the ListContainersSegment request.
func (c *ServiceClient) listContainersResponder(resp *azcore.Response) (*ListContainersPage, error) {
	if err := resp.CheckStatusCode(http.StatusOK); err != nil {
		return nil, err
	}
	result := &ListContainersPage{response: resp}
	return result, resp.UnmarshalAsXML(result)
}
