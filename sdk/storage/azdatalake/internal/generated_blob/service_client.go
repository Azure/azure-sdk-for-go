// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated_blob

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

func (client *ServiceClient) Endpoint() string {
	return client.url
}

func (client *ServiceClient) InternalClient() *azcore.Client {
	return client.internal
}

// ServiceClient contains the methods for the Service group.
// Don't use this type directly, use a constructor function instead.
type ServiceClient struct {
	internal *azcore.Client
	url      string
}

// NewServiceClient creates a new instance of ServiceClient with the specified values.
//   - endpoint - The URL of the service account, container, or blob that is the target of the desired operation.
//   - azClient - azcore.Client is a basic HTTP client. It consists of a pipeline and tracing provider.
func NewServiceClient(endpoint string, azClient *azcore.Client) *ServiceClient {
	client := &ServiceClient{
		internal: azClient,
		url:      endpoint,
	}
	return client
}

// NewListContainersSegmentPager - The List Containers Segment operation returns a list of the containers under the specified
// account
//
// Generated from API version 2026-06-06
//   - options - ServiceClientListContainersSegmentOptions contains the optional parameters for the ServiceClient.NewListContainersSegmentPager
//     method.
//
// listContainersSegmentCreateRequest creates the ListContainersSegment request.
func (client *ServiceClient) ListContainersSegmentCreateRequest(ctx context.Context, options *ServiceClientListContainersSegmentOptions) (*policy.Request, error) {
	req, err := runtime.NewRequest(ctx, http.MethodGet, client.url)
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("comp", "list")
	if options != nil && options.Include != nil {
		reqQP.Set("include", strings.Join(strings.Fields(strings.Trim(fmt.Sprint(options.Include), "[]")), ","))
	}
	if options != nil && options.Marker != nil {
		reqQP.Set("marker", *options.Marker)
	}
	if options != nil && options.Maxresults != nil {
		reqQP.Set("maxresults", strconv.FormatInt(int64(*options.Maxresults), 10))
	}
	if options != nil && options.Prefix != nil {
		reqQP.Set("prefix", *options.Prefix)
	}
	if options != nil && options.Timeout != nil {
		reqQP.Set("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
	}
	req.Raw().URL.RawQuery = strings.ReplaceAll(reqQP.Encode(), "+", "%20")
	req.Raw().Header["Accept"] = []string{"application/xml"}
	if options != nil && options.RequestID != nil {
		req.Raw().Header["x-ms-client-request-id"] = []string{*options.RequestID}
	}
	req.Raw().Header["x-ms-version"] = []string{"2026-06-06"}
	return req, nil
}

// listContainersSegmentHandleResponse handles the ListContainersSegment response.
func (client *ServiceClient) ListContainersSegmentHandleResponse(resp *http.Response) (ServiceClientListFileSystemsSegmentResponse, error) {
	result := ServiceClientListFileSystemsSegmentResponse{}
	if val := resp.Header.Get("x-ms-client-request-id"); val != "" {
		result.ClientRequestID = &val
	}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	if err := runtime.UnmarshalAsXML(resp, &result.ListFileSystemsSegmentResponse); err != nil {
		return ServiceClientListFileSystemsSegmentResponse{}, err
	}
	return result, nil
}
