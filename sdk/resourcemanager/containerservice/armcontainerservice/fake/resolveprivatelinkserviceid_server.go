// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package fake

import (
	"context"
	"errors"
	"fmt"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake/server"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v7"
	"net/http"
	"net/url"
	"regexp"
)

// ResolvePrivateLinkServiceIDServer is a fake server for instances of the armcontainerservice.ResolvePrivateLinkServiceIDClient type.
type ResolvePrivateLinkServiceIDServer struct {
	// POST is the fake for method ResolvePrivateLinkServiceIDClient.POST
	// HTTP status codes to indicate success: http.StatusOK
	POST func(ctx context.Context, resourceGroupName string, resourceName string, parameters armcontainerservice.PrivateLinkResource, options *armcontainerservice.ResolvePrivateLinkServiceIDClientPOSTOptions) (resp azfake.Responder[armcontainerservice.ResolvePrivateLinkServiceIDClientPOSTResponse], errResp azfake.ErrorResponder)
}

// NewResolvePrivateLinkServiceIDServerTransport creates a new instance of ResolvePrivateLinkServiceIDServerTransport with the provided implementation.
// The returned ResolvePrivateLinkServiceIDServerTransport instance is connected to an instance of armcontainerservice.ResolvePrivateLinkServiceIDClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewResolvePrivateLinkServiceIDServerTransport(srv *ResolvePrivateLinkServiceIDServer) *ResolvePrivateLinkServiceIDServerTransport {
	return &ResolvePrivateLinkServiceIDServerTransport{srv: srv}
}

// ResolvePrivateLinkServiceIDServerTransport connects instances of armcontainerservice.ResolvePrivateLinkServiceIDClient to instances of ResolvePrivateLinkServiceIDServer.
// Don't use this type directly, use NewResolvePrivateLinkServiceIDServerTransport instead.
type ResolvePrivateLinkServiceIDServerTransport struct {
	srv *ResolvePrivateLinkServiceIDServer
}

// Do implements the policy.Transporter interface for ResolvePrivateLinkServiceIDServerTransport.
func (r *ResolvePrivateLinkServiceIDServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return r.dispatchToMethodFake(req, method)
}

func (r *ResolvePrivateLinkServiceIDServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if resolvePrivateLinkServiceIdServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = resolvePrivateLinkServiceIdServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "ResolvePrivateLinkServiceIDClient.POST":
				res.resp, res.err = r.dispatchPOST(req)
			default:
				res.err = fmt.Errorf("unhandled API %s", method)
			}

		}
		select {
		case resultChan <- res:
		case <-req.Context().Done():
		}
	}()

	select {
	case <-req.Context().Done():
		return nil, req.Context().Err()
	case res := <-resultChan:
		return res.resp, res.err
	}
}

func (r *ResolvePrivateLinkServiceIDServerTransport) dispatchPOST(req *http.Request) (*http.Response, error) {
	if r.srv.POST == nil {
		return nil, &nonRetriableError{errors.New("fake for method POST not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ContainerService/managedClusters/(?P<resourceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resolvePrivateLinkServiceId`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armcontainerservice.PrivateLinkResource](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	resourceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := r.srv.POST(req.Context(), resourceGroupNameParam, resourceNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).PrivateLinkResource, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to ResolvePrivateLinkServiceIDServerTransport
var resolvePrivateLinkServiceIdServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
