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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"
	"net/http"
	"net/url"
	"regexp"
)

// LoadBalancersServer is a fake server for instances of the armnetwork.LoadBalancersClient type.
type LoadBalancersServer struct {
	// BeginCreateOrUpdate is the fake for method LoadBalancersClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, loadBalancerName string, parameters armnetwork.LoadBalancer, options *armnetwork.LoadBalancersClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armnetwork.LoadBalancersClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method LoadBalancersClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, loadBalancerName string, options *armnetwork.LoadBalancersClientBeginDeleteOptions) (resp azfake.PollerResponder[armnetwork.LoadBalancersClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method LoadBalancersClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, loadBalancerName string, options *armnetwork.LoadBalancersClientGetOptions) (resp azfake.Responder[armnetwork.LoadBalancersClientGetResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method LoadBalancersClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(resourceGroupName string, options *armnetwork.LoadBalancersClientListOptions) (resp azfake.PagerResponder[armnetwork.LoadBalancersClientListResponse])

	// NewListAllPager is the fake for method LoadBalancersClient.NewListAllPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListAllPager func(options *armnetwork.LoadBalancersClientListAllOptions) (resp azfake.PagerResponder[armnetwork.LoadBalancersClientListAllResponse])

	// BeginListInboundNatRulePortMappings is the fake for method LoadBalancersClient.BeginListInboundNatRulePortMappings
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginListInboundNatRulePortMappings func(ctx context.Context, groupName string, loadBalancerName string, backendPoolName string, parameters armnetwork.QueryInboundNatRulePortMappingRequest, options *armnetwork.LoadBalancersClientBeginListInboundNatRulePortMappingsOptions) (resp azfake.PollerResponder[armnetwork.LoadBalancersClientListInboundNatRulePortMappingsResponse], errResp azfake.ErrorResponder)

	// BeginSwapPublicIPAddresses is the fake for method LoadBalancersClient.BeginSwapPublicIPAddresses
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginSwapPublicIPAddresses func(ctx context.Context, location string, parameters armnetwork.LoadBalancerVipSwapRequest, options *armnetwork.LoadBalancersClientBeginSwapPublicIPAddressesOptions) (resp azfake.PollerResponder[armnetwork.LoadBalancersClientSwapPublicIPAddressesResponse], errResp azfake.ErrorResponder)

	// UpdateTags is the fake for method LoadBalancersClient.UpdateTags
	// HTTP status codes to indicate success: http.StatusOK
	UpdateTags func(ctx context.Context, resourceGroupName string, loadBalancerName string, parameters armnetwork.TagsObject, options *armnetwork.LoadBalancersClientUpdateTagsOptions) (resp azfake.Responder[armnetwork.LoadBalancersClientUpdateTagsResponse], errResp azfake.ErrorResponder)
}

// NewLoadBalancersServerTransport creates a new instance of LoadBalancersServerTransport with the provided implementation.
// The returned LoadBalancersServerTransport instance is connected to an instance of armnetwork.LoadBalancersClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewLoadBalancersServerTransport(srv *LoadBalancersServer) *LoadBalancersServerTransport {
	return &LoadBalancersServerTransport{
		srv:                                 srv,
		beginCreateOrUpdate:                 newTracker[azfake.PollerResponder[armnetwork.LoadBalancersClientCreateOrUpdateResponse]](),
		beginDelete:                         newTracker[azfake.PollerResponder[armnetwork.LoadBalancersClientDeleteResponse]](),
		newListPager:                        newTracker[azfake.PagerResponder[armnetwork.LoadBalancersClientListResponse]](),
		newListAllPager:                     newTracker[azfake.PagerResponder[armnetwork.LoadBalancersClientListAllResponse]](),
		beginListInboundNatRulePortMappings: newTracker[azfake.PollerResponder[armnetwork.LoadBalancersClientListInboundNatRulePortMappingsResponse]](),
		beginSwapPublicIPAddresses:          newTracker[azfake.PollerResponder[armnetwork.LoadBalancersClientSwapPublicIPAddressesResponse]](),
	}
}

// LoadBalancersServerTransport connects instances of armnetwork.LoadBalancersClient to instances of LoadBalancersServer.
// Don't use this type directly, use NewLoadBalancersServerTransport instead.
type LoadBalancersServerTransport struct {
	srv                                 *LoadBalancersServer
	beginCreateOrUpdate                 *tracker[azfake.PollerResponder[armnetwork.LoadBalancersClientCreateOrUpdateResponse]]
	beginDelete                         *tracker[azfake.PollerResponder[armnetwork.LoadBalancersClientDeleteResponse]]
	newListPager                        *tracker[azfake.PagerResponder[armnetwork.LoadBalancersClientListResponse]]
	newListAllPager                     *tracker[azfake.PagerResponder[armnetwork.LoadBalancersClientListAllResponse]]
	beginListInboundNatRulePortMappings *tracker[azfake.PollerResponder[armnetwork.LoadBalancersClientListInboundNatRulePortMappingsResponse]]
	beginSwapPublicIPAddresses          *tracker[azfake.PollerResponder[armnetwork.LoadBalancersClientSwapPublicIPAddressesResponse]]
}

// Do implements the policy.Transporter interface for LoadBalancersServerTransport.
func (l *LoadBalancersServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "LoadBalancersClient.BeginCreateOrUpdate":
		resp, err = l.dispatchBeginCreateOrUpdate(req)
	case "LoadBalancersClient.BeginDelete":
		resp, err = l.dispatchBeginDelete(req)
	case "LoadBalancersClient.Get":
		resp, err = l.dispatchGet(req)
	case "LoadBalancersClient.NewListPager":
		resp, err = l.dispatchNewListPager(req)
	case "LoadBalancersClient.NewListAllPager":
		resp, err = l.dispatchNewListAllPager(req)
	case "LoadBalancersClient.BeginListInboundNatRulePortMappings":
		resp, err = l.dispatchBeginListInboundNatRulePortMappings(req)
	case "LoadBalancersClient.BeginSwapPublicIPAddresses":
		resp, err = l.dispatchBeginSwapPublicIPAddresses(req)
	case "LoadBalancersClient.UpdateTags":
		resp, err = l.dispatchUpdateTags(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (l *LoadBalancersServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if l.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := l.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft.Network/loadBalancers/(?P<loadBalancerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armnetwork.LoadBalancer](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		loadBalancerNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("loadBalancerName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := l.srv.BeginCreateOrUpdate(req.Context(), resourceGroupNameUnescaped, loadBalancerNameUnescaped, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateOrUpdate = &respr
		l.beginCreateOrUpdate.add(req, beginCreateOrUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		l.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		l.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (l *LoadBalancersServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if l.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := l.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft.Network/loadBalancers/(?P<loadBalancerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		loadBalancerNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("loadBalancerName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := l.srv.BeginDelete(req.Context(), resourceGroupNameUnescaped, loadBalancerNameUnescaped, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		l.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		l.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		l.beginDelete.remove(req)
	}

	return resp, nil
}

func (l *LoadBalancersServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if l.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft.Network/loadBalancers/(?P<loadBalancerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	qp := req.URL.Query()
	resourceGroupNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	loadBalancerNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("loadBalancerName")])
	if err != nil {
		return nil, err
	}
	expandUnescaped, err := url.QueryUnescape(qp.Get("$expand"))
	if err != nil {
		return nil, err
	}
	expandParam := getOptional(expandUnescaped)
	var options *armnetwork.LoadBalancersClientGetOptions
	if expandParam != nil {
		options = &armnetwork.LoadBalancersClientGetOptions{
			Expand: expandParam,
		}
	}
	respr, errRespr := l.srv.Get(req.Context(), resourceGroupNameUnescaped, loadBalancerNameUnescaped, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).LoadBalancer, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (l *LoadBalancersServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if l.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := l.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft.Network/loadBalancers`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		resp := l.srv.NewListPager(resourceGroupNameUnescaped, nil)
		newListPager = &resp
		l.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armnetwork.LoadBalancersClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		l.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		l.newListPager.remove(req)
	}
	return resp, nil
}

func (l *LoadBalancersServerTransport) dispatchNewListAllPager(req *http.Request) (*http.Response, error) {
	if l.srv.NewListAllPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListAllPager not implemented")}
	}
	newListAllPager := l.newListAllPager.get(req)
	if newListAllPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft.Network/loadBalancers`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resp := l.srv.NewListAllPager(nil)
		newListAllPager = &resp
		l.newListAllPager.add(req, newListAllPager)
		server.PagerResponderInjectNextLinks(newListAllPager, req, func(page *armnetwork.LoadBalancersClientListAllResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListAllPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		l.newListAllPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListAllPager) {
		l.newListAllPager.remove(req)
	}
	return resp, nil
}

func (l *LoadBalancersServerTransport) dispatchBeginListInboundNatRulePortMappings(req *http.Request) (*http.Response, error) {
	if l.srv.BeginListInboundNatRulePortMappings == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginListInboundNatRulePortMappings not implemented")}
	}
	beginListInboundNatRulePortMappings := l.beginListInboundNatRulePortMappings.get(req)
	if beginListInboundNatRulePortMappings == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<groupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft.Network/loadBalancers/(?P<loadBalancerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/backendAddressPools/(?P<backendPoolName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/queryInboundNatRulePortMapping`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armnetwork.QueryInboundNatRulePortMappingRequest](req)
		if err != nil {
			return nil, err
		}
		groupNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("groupName")])
		if err != nil {
			return nil, err
		}
		loadBalancerNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("loadBalancerName")])
		if err != nil {
			return nil, err
		}
		backendPoolNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("backendPoolName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := l.srv.BeginListInboundNatRulePortMappings(req.Context(), groupNameUnescaped, loadBalancerNameUnescaped, backendPoolNameUnescaped, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginListInboundNatRulePortMappings = &respr
		l.beginListInboundNatRulePortMappings.add(req, beginListInboundNatRulePortMappings)
	}

	resp, err := server.PollerResponderNext(beginListInboundNatRulePortMappings, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		l.beginListInboundNatRulePortMappings.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginListInboundNatRulePortMappings) {
		l.beginListInboundNatRulePortMappings.remove(req)
	}

	return resp, nil
}

func (l *LoadBalancersServerTransport) dispatchBeginSwapPublicIPAddresses(req *http.Request) (*http.Response, error) {
	if l.srv.BeginSwapPublicIPAddresses == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginSwapPublicIPAddresses not implemented")}
	}
	beginSwapPublicIPAddresses := l.beginSwapPublicIPAddresses.get(req)
	if beginSwapPublicIPAddresses == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft.Network/locations/(?P<location>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/setLoadBalancerFrontendPublicIpAddresses`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armnetwork.LoadBalancerVipSwapRequest](req)
		if err != nil {
			return nil, err
		}
		locationUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("location")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := l.srv.BeginSwapPublicIPAddresses(req.Context(), locationUnescaped, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginSwapPublicIPAddresses = &respr
		l.beginSwapPublicIPAddresses.add(req, beginSwapPublicIPAddresses)
	}

	resp, err := server.PollerResponderNext(beginSwapPublicIPAddresses, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		l.beginSwapPublicIPAddresses.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginSwapPublicIPAddresses) {
		l.beginSwapPublicIPAddresses.remove(req)
	}

	return resp, nil
}

func (l *LoadBalancersServerTransport) dispatchUpdateTags(req *http.Request) (*http.Response, error) {
	if l.srv.UpdateTags == nil {
		return nil, &nonRetriableError{errors.New("fake for method UpdateTags not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft.Network/loadBalancers/(?P<loadBalancerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armnetwork.TagsObject](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	loadBalancerNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("loadBalancerName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := l.srv.UpdateTags(req.Context(), resourceGroupNameUnescaped, loadBalancerNameUnescaped, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).LoadBalancer, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
