//go:build go1.18
// +build go1.18

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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v3"
	"net/http"
	"regexp"
)

// IPAllocationsServer is a fake server for instances of the armnetwork.IPAllocationsClient type.
type IPAllocationsServer struct {
	// BeginCreateOrUpdate is the fake for method IPAllocationsClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, ipAllocationName string, parameters armnetwork.IPAllocation, options *armnetwork.IPAllocationsClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armnetwork.IPAllocationsClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method IPAllocationsClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, ipAllocationName string, options *armnetwork.IPAllocationsClientBeginDeleteOptions) (resp azfake.PollerResponder[armnetwork.IPAllocationsClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method IPAllocationsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, ipAllocationName string, options *armnetwork.IPAllocationsClientGetOptions) (resp azfake.Responder[armnetwork.IPAllocationsClientGetResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method IPAllocationsClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(options *armnetwork.IPAllocationsClientListOptions) (resp azfake.PagerResponder[armnetwork.IPAllocationsClientListResponse])

	// NewListByResourceGroupPager is the fake for method IPAllocationsClient.NewListByResourceGroupPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByResourceGroupPager func(resourceGroupName string, options *armnetwork.IPAllocationsClientListByResourceGroupOptions) (resp azfake.PagerResponder[armnetwork.IPAllocationsClientListByResourceGroupResponse])

	// UpdateTags is the fake for method IPAllocationsClient.UpdateTags
	// HTTP status codes to indicate success: http.StatusOK
	UpdateTags func(ctx context.Context, resourceGroupName string, ipAllocationName string, parameters armnetwork.TagsObject, options *armnetwork.IPAllocationsClientUpdateTagsOptions) (resp azfake.Responder[armnetwork.IPAllocationsClientUpdateTagsResponse], errResp azfake.ErrorResponder)
}

// NewIPAllocationsServerTransport creates a new instance of IPAllocationsServerTransport with the provided implementation.
// The returned IPAllocationsServerTransport instance is connected to an instance of armnetwork.IPAllocationsClient by way of the
// undefined.Transporter field.
func NewIPAllocationsServerTransport(srv *IPAllocationsServer) *IPAllocationsServerTransport {
	return &IPAllocationsServerTransport{srv: srv}
}

// IPAllocationsServerTransport connects instances of armnetwork.IPAllocationsClient to instances of IPAllocationsServer.
// Don't use this type directly, use NewIPAllocationsServerTransport instead.
type IPAllocationsServerTransport struct {
	srv                         *IPAllocationsServer
	beginCreateOrUpdate         *azfake.PollerResponder[armnetwork.IPAllocationsClientCreateOrUpdateResponse]
	beginDelete                 *azfake.PollerResponder[armnetwork.IPAllocationsClientDeleteResponse]
	newListPager                *azfake.PagerResponder[armnetwork.IPAllocationsClientListResponse]
	newListByResourceGroupPager *azfake.PagerResponder[armnetwork.IPAllocationsClientListByResourceGroupResponse]
}

// Do implements the policy.Transporter interface for IPAllocationsServerTransport.
func (i *IPAllocationsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "IPAllocationsClient.BeginCreateOrUpdate":
		resp, err = i.dispatchBeginCreateOrUpdate(req)
	case "IPAllocationsClient.BeginDelete":
		resp, err = i.dispatchBeginDelete(req)
	case "IPAllocationsClient.Get":
		resp, err = i.dispatchGet(req)
	case "IPAllocationsClient.NewListPager":
		resp, err = i.dispatchNewListPager(req)
	case "IPAllocationsClient.NewListByResourceGroupPager":
		resp, err = i.dispatchNewListByResourceGroupPager(req)
	case "IPAllocationsClient.UpdateTags":
		resp, err = i.dispatchUpdateTags(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (i *IPAllocationsServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if i.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("method BeginCreateOrUpdate not implemented")}
	}
	if i.beginCreateOrUpdate == nil {
		const regexStr = "/subscriptions/(?P<subscriptionId>[a-zA-Z0-9-_]+)/resourceGroups/(?P<resourceGroupName>[a-zA-Z0-9-_]+)/providers/Microsoft.Network/IpAllocations/(?P<ipAllocationName>[a-zA-Z0-9-_]+)"
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.Path)
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armnetwork.IPAllocation](req)
		if err != nil {
			return nil, err
		}
		respr, errRespr := i.srv.BeginCreateOrUpdate(req.Context(), matches[regex.SubexpIndex("resourceGroupName")], matches[regex.SubexpIndex("ipAllocationName")], body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		i.beginCreateOrUpdate = &respr
	}

	resp, err := server.PollerResponderNext(i.beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(i.beginCreateOrUpdate) {
		i.beginCreateOrUpdate = nil
	}

	return resp, nil
}

func (i *IPAllocationsServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if i.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("method BeginDelete not implemented")}
	}
	if i.beginDelete == nil {
		const regexStr = "/subscriptions/(?P<subscriptionId>[a-zA-Z0-9-_]+)/resourceGroups/(?P<resourceGroupName>[a-zA-Z0-9-_]+)/providers/Microsoft.Network/IpAllocations/(?P<ipAllocationName>[a-zA-Z0-9-_]+)"
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.Path)
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		respr, errRespr := i.srv.BeginDelete(req.Context(), matches[regex.SubexpIndex("resourceGroupName")], matches[regex.SubexpIndex("ipAllocationName")], nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		i.beginDelete = &respr
	}

	resp, err := server.PollerResponderNext(i.beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(i.beginDelete) {
		i.beginDelete = nil
	}

	return resp, nil
}

func (i *IPAllocationsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if i.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("method Get not implemented")}
	}
	const regexStr = "/subscriptions/(?P<subscriptionId>[a-zA-Z0-9-_]+)/resourceGroups/(?P<resourceGroupName>[a-zA-Z0-9-_]+)/providers/Microsoft.Network/IpAllocations/(?P<ipAllocationName>[a-zA-Z0-9-_]+)"
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.Path)
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	qp := req.URL.Query()
	expandParam := getOptional(qp.Get("$expand"))
	var options *armnetwork.IPAllocationsClientGetOptions
	if expandParam != nil {
		options = &armnetwork.IPAllocationsClientGetOptions{
			Expand: expandParam,
		}
	}
	respr, errRespr := i.srv.Get(req.Context(), matches[regex.SubexpIndex("resourceGroupName")], matches[regex.SubexpIndex("ipAllocationName")], options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).IPAllocation, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *IPAllocationsServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if i.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("method NewListPager not implemented")}
	}
	if i.newListPager == nil {
		const regexStr = "/subscriptions/(?P<subscriptionId>[a-zA-Z0-9-_]+)/providers/Microsoft.Network/IpAllocations"
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.Path)
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resp := i.srv.NewListPager(nil)
		i.newListPager = &resp
		server.PagerResponderInjectNextLinks(i.newListPager, req, func(page *armnetwork.IPAllocationsClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(i.newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(i.newListPager) {
		i.newListPager = nil
	}
	return resp, nil
}

func (i *IPAllocationsServerTransport) dispatchNewListByResourceGroupPager(req *http.Request) (*http.Response, error) {
	if i.srv.NewListByResourceGroupPager == nil {
		return nil, &nonRetriableError{errors.New("method NewListByResourceGroupPager not implemented")}
	}
	if i.newListByResourceGroupPager == nil {
		const regexStr = "/subscriptions/(?P<subscriptionId>[a-zA-Z0-9-_]+)/resourceGroups/(?P<resourceGroupName>[a-zA-Z0-9-_]+)/providers/Microsoft.Network/IpAllocations"
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.Path)
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resp := i.srv.NewListByResourceGroupPager(matches[regex.SubexpIndex("resourceGroupName")], nil)
		i.newListByResourceGroupPager = &resp
		server.PagerResponderInjectNextLinks(i.newListByResourceGroupPager, req, func(page *armnetwork.IPAllocationsClientListByResourceGroupResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(i.newListByResourceGroupPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(i.newListByResourceGroupPager) {
		i.newListByResourceGroupPager = nil
	}
	return resp, nil
}

func (i *IPAllocationsServerTransport) dispatchUpdateTags(req *http.Request) (*http.Response, error) {
	if i.srv.UpdateTags == nil {
		return nil, &nonRetriableError{errors.New("method UpdateTags not implemented")}
	}
	const regexStr = "/subscriptions/(?P<subscriptionId>[a-zA-Z0-9-_]+)/resourceGroups/(?P<resourceGroupName>[a-zA-Z0-9-_]+)/providers/Microsoft.Network/IpAllocations/(?P<ipAllocationName>[a-zA-Z0-9-_]+)"
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.Path)
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armnetwork.TagsObject](req)
	if err != nil {
		return nil, err
	}
	respr, errRespr := i.srv.UpdateTags(req.Context(), matches[regex.SubexpIndex("resourceGroupName")], matches[regex.SubexpIndex("ipAllocationName")], body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).IPAllocation, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
