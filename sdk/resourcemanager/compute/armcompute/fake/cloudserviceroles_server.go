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
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7"
	"net/http"
	"net/url"
	"regexp"
)

// CloudServiceRolesServer is a fake server for instances of the armcompute.CloudServiceRolesClient type.
type CloudServiceRolesServer struct {
	// Get is the fake for method CloudServiceRolesClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, roleName string, resourceGroupName string, cloudServiceName string, options *armcompute.CloudServiceRolesClientGetOptions) (resp azfake.Responder[armcompute.CloudServiceRolesClientGetResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method CloudServiceRolesClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(resourceGroupName string, cloudServiceName string, options *armcompute.CloudServiceRolesClientListOptions) (resp azfake.PagerResponder[armcompute.CloudServiceRolesClientListResponse])
}

// NewCloudServiceRolesServerTransport creates a new instance of CloudServiceRolesServerTransport with the provided implementation.
// The returned CloudServiceRolesServerTransport instance is connected to an instance of armcompute.CloudServiceRolesClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewCloudServiceRolesServerTransport(srv *CloudServiceRolesServer) *CloudServiceRolesServerTransport {
	return &CloudServiceRolesServerTransport{
		srv:          srv,
		newListPager: newTracker[azfake.PagerResponder[armcompute.CloudServiceRolesClientListResponse]](),
	}
}

// CloudServiceRolesServerTransport connects instances of armcompute.CloudServiceRolesClient to instances of CloudServiceRolesServer.
// Don't use this type directly, use NewCloudServiceRolesServerTransport instead.
type CloudServiceRolesServerTransport struct {
	srv          *CloudServiceRolesServer
	newListPager *tracker[azfake.PagerResponder[armcompute.CloudServiceRolesClientListResponse]]
}

// Do implements the policy.Transporter interface for CloudServiceRolesServerTransport.
func (c *CloudServiceRolesServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return c.dispatchToMethodFake(req, method)
}

func (c *CloudServiceRolesServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if cloudServiceRolesServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = cloudServiceRolesServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "CloudServiceRolesClient.Get":
				res.resp, res.err = c.dispatchGet(req)
			case "CloudServiceRolesClient.NewListPager":
				res.resp, res.err = c.dispatchNewListPager(req)
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

func (c *CloudServiceRolesServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if c.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Compute/cloudServices/(?P<cloudServiceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/roles/(?P<roleName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	roleNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("roleName")])
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	cloudServiceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("cloudServiceName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := c.srv.Get(req.Context(), roleNameParam, resourceGroupNameParam, cloudServiceNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).CloudServiceRole, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *CloudServiceRolesServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if c.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := c.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Compute/cloudServices/(?P<cloudServiceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/roles`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		cloudServiceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("cloudServiceName")])
		if err != nil {
			return nil, err
		}
		resp := c.srv.NewListPager(resourceGroupNameParam, cloudServiceNameParam, nil)
		newListPager = &resp
		c.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armcompute.CloudServiceRolesClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		c.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		c.newListPager.remove(req)
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to CloudServiceRolesServerTransport
var cloudServiceRolesServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
