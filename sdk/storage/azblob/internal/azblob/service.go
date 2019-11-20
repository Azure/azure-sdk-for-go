// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// Service contains all methods for operations on the service.
type Service struct{}

// ListContainersCreateRequest prepares the ListContainersSegment request.
func (Service) ListContainersCreateRequest(u *url.URL, p azcore.Pipeline, options *ListContainersOptions) *azcore.Request {
	req := p.NewRequest(http.MethodGet, *u)
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
	req.Header.Set("x-ms-version", "2018-11-09")
	return req
}

// ListContainersCreateResponse handles the response to the ListContainersSegment request.
func (Service) ListContainersCreateResponse(resp *azcore.Response) (*ListContainersPage, error) {
	if err := resp.CheckStatusCode(http.StatusOK); err != nil {
		return nil, err
	}
	result := &ListContainersPage{}
	return result, resp.UnmarshalAsXML(result)
}
