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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hardwaresecuritymodules/armhardwaresecuritymodules/v2"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

// DedicatedHsmServer is a fake server for instances of the armhardwaresecuritymodules.DedicatedHsmClient type.
type DedicatedHsmServer struct {
	// BeginCreateOrUpdate is the fake for method DedicatedHsmClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, name string, parameters armhardwaresecuritymodules.DedicatedHsm, options *armhardwaresecuritymodules.DedicatedHsmClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armhardwaresecuritymodules.DedicatedHsmClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method DedicatedHsmClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, name string, options *armhardwaresecuritymodules.DedicatedHsmClientBeginDeleteOptions) (resp azfake.PollerResponder[armhardwaresecuritymodules.DedicatedHsmClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method DedicatedHsmClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, name string, options *armhardwaresecuritymodules.DedicatedHsmClientGetOptions) (resp azfake.Responder[armhardwaresecuritymodules.DedicatedHsmClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByResourceGroupPager is the fake for method DedicatedHsmClient.NewListByResourceGroupPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByResourceGroupPager func(resourceGroupName string, options *armhardwaresecuritymodules.DedicatedHsmClientListByResourceGroupOptions) (resp azfake.PagerResponder[armhardwaresecuritymodules.DedicatedHsmClientListByResourceGroupResponse])

	// NewListBySubscriptionPager is the fake for method DedicatedHsmClient.NewListBySubscriptionPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListBySubscriptionPager func(options *armhardwaresecuritymodules.DedicatedHsmClientListBySubscriptionOptions) (resp azfake.PagerResponder[armhardwaresecuritymodules.DedicatedHsmClientListBySubscriptionResponse])

	// NewListOutboundNetworkDependenciesEndpointsPager is the fake for method DedicatedHsmClient.NewListOutboundNetworkDependenciesEndpointsPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListOutboundNetworkDependenciesEndpointsPager func(resourceGroupName string, name string, options *armhardwaresecuritymodules.DedicatedHsmClientListOutboundNetworkDependenciesEndpointsOptions) (resp azfake.PagerResponder[armhardwaresecuritymodules.DedicatedHsmClientListOutboundNetworkDependenciesEndpointsResponse])

	// BeginUpdate is the fake for method DedicatedHsmClient.BeginUpdate
	// HTTP status codes to indicate success: http.StatusOK
	BeginUpdate func(ctx context.Context, resourceGroupName string, name string, parameters armhardwaresecuritymodules.DedicatedHsmPatchParameters, options *armhardwaresecuritymodules.DedicatedHsmClientBeginUpdateOptions) (resp azfake.PollerResponder[armhardwaresecuritymodules.DedicatedHsmClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewDedicatedHsmServerTransport creates a new instance of DedicatedHsmServerTransport with the provided implementation.
// The returned DedicatedHsmServerTransport instance is connected to an instance of armhardwaresecuritymodules.DedicatedHsmClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewDedicatedHsmServerTransport(srv *DedicatedHsmServer) *DedicatedHsmServerTransport {
	return &DedicatedHsmServerTransport{
		srv:                         srv,
		beginCreateOrUpdate:         newTracker[azfake.PollerResponder[armhardwaresecuritymodules.DedicatedHsmClientCreateOrUpdateResponse]](),
		beginDelete:                 newTracker[azfake.PollerResponder[armhardwaresecuritymodules.DedicatedHsmClientDeleteResponse]](),
		newListByResourceGroupPager: newTracker[azfake.PagerResponder[armhardwaresecuritymodules.DedicatedHsmClientListByResourceGroupResponse]](),
		newListBySubscriptionPager:  newTracker[azfake.PagerResponder[armhardwaresecuritymodules.DedicatedHsmClientListBySubscriptionResponse]](),
		newListOutboundNetworkDependenciesEndpointsPager: newTracker[azfake.PagerResponder[armhardwaresecuritymodules.DedicatedHsmClientListOutboundNetworkDependenciesEndpointsResponse]](),
		beginUpdate: newTracker[azfake.PollerResponder[armhardwaresecuritymodules.DedicatedHsmClientUpdateResponse]](),
	}
}

// DedicatedHsmServerTransport connects instances of armhardwaresecuritymodules.DedicatedHsmClient to instances of DedicatedHsmServer.
// Don't use this type directly, use NewDedicatedHsmServerTransport instead.
type DedicatedHsmServerTransport struct {
	srv                                              *DedicatedHsmServer
	beginCreateOrUpdate                              *tracker[azfake.PollerResponder[armhardwaresecuritymodules.DedicatedHsmClientCreateOrUpdateResponse]]
	beginDelete                                      *tracker[azfake.PollerResponder[armhardwaresecuritymodules.DedicatedHsmClientDeleteResponse]]
	newListByResourceGroupPager                      *tracker[azfake.PagerResponder[armhardwaresecuritymodules.DedicatedHsmClientListByResourceGroupResponse]]
	newListBySubscriptionPager                       *tracker[azfake.PagerResponder[armhardwaresecuritymodules.DedicatedHsmClientListBySubscriptionResponse]]
	newListOutboundNetworkDependenciesEndpointsPager *tracker[azfake.PagerResponder[armhardwaresecuritymodules.DedicatedHsmClientListOutboundNetworkDependenciesEndpointsResponse]]
	beginUpdate                                      *tracker[azfake.PollerResponder[armhardwaresecuritymodules.DedicatedHsmClientUpdateResponse]]
}

// Do implements the policy.Transporter interface for DedicatedHsmServerTransport.
func (d *DedicatedHsmServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "DedicatedHsmClient.BeginCreateOrUpdate":
		resp, err = d.dispatchBeginCreateOrUpdate(req)
	case "DedicatedHsmClient.BeginDelete":
		resp, err = d.dispatchBeginDelete(req)
	case "DedicatedHsmClient.Get":
		resp, err = d.dispatchGet(req)
	case "DedicatedHsmClient.NewListByResourceGroupPager":
		resp, err = d.dispatchNewListByResourceGroupPager(req)
	case "DedicatedHsmClient.NewListBySubscriptionPager":
		resp, err = d.dispatchNewListBySubscriptionPager(req)
	case "DedicatedHsmClient.NewListOutboundNetworkDependenciesEndpointsPager":
		resp, err = d.dispatchNewListOutboundNetworkDependenciesEndpointsPager(req)
	case "DedicatedHsmClient.BeginUpdate":
		resp, err = d.dispatchBeginUpdate(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (d *DedicatedHsmServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if d.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := d.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HardwareSecurityModules/dedicatedHSMs/(?P<name>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armhardwaresecuritymodules.DedicatedHsm](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		nameParam, err := url.PathUnescape(matches[regex.SubexpIndex("name")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := d.srv.BeginCreateOrUpdate(req.Context(), resourceGroupNameParam, nameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateOrUpdate = &respr
		d.beginCreateOrUpdate.add(req, beginCreateOrUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		d.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		d.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (d *DedicatedHsmServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if d.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := d.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HardwareSecurityModules/dedicatedHSMs/(?P<name>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		nameParam, err := url.PathUnescape(matches[regex.SubexpIndex("name")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := d.srv.BeginDelete(req.Context(), resourceGroupNameParam, nameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		d.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		d.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		d.beginDelete.remove(req)
	}

	return resp, nil
}

func (d *DedicatedHsmServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if d.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HardwareSecurityModules/dedicatedHSMs/(?P<name>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	nameParam, err := url.PathUnescape(matches[regex.SubexpIndex("name")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := d.srv.Get(req.Context(), resourceGroupNameParam, nameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).DedicatedHsm, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DedicatedHsmServerTransport) dispatchNewListByResourceGroupPager(req *http.Request) (*http.Response, error) {
	if d.srv.NewListByResourceGroupPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByResourceGroupPager not implemented")}
	}
	newListByResourceGroupPager := d.newListByResourceGroupPager.get(req)
	if newListByResourceGroupPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HardwareSecurityModules/dedicatedHSMs`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		topUnescaped, err := url.QueryUnescape(qp.Get("$top"))
		if err != nil {
			return nil, err
		}
		topParam, err := parseOptional(topUnescaped, func(v string) (int32, error) {
			p, parseErr := strconv.ParseInt(v, 10, 32)
			if parseErr != nil {
				return 0, parseErr
			}
			return int32(p), nil
		})
		if err != nil {
			return nil, err
		}
		var options *armhardwaresecuritymodules.DedicatedHsmClientListByResourceGroupOptions
		if topParam != nil {
			options = &armhardwaresecuritymodules.DedicatedHsmClientListByResourceGroupOptions{
				Top: topParam,
			}
		}
		resp := d.srv.NewListByResourceGroupPager(resourceGroupNameParam, options)
		newListByResourceGroupPager = &resp
		d.newListByResourceGroupPager.add(req, newListByResourceGroupPager)
		server.PagerResponderInjectNextLinks(newListByResourceGroupPager, req, func(page *armhardwaresecuritymodules.DedicatedHsmClientListByResourceGroupResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByResourceGroupPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		d.newListByResourceGroupPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByResourceGroupPager) {
		d.newListByResourceGroupPager.remove(req)
	}
	return resp, nil
}

func (d *DedicatedHsmServerTransport) dispatchNewListBySubscriptionPager(req *http.Request) (*http.Response, error) {
	if d.srv.NewListBySubscriptionPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListBySubscriptionPager not implemented")}
	}
	newListBySubscriptionPager := d.newListBySubscriptionPager.get(req)
	if newListBySubscriptionPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HardwareSecurityModules/dedicatedHSMs`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		topUnescaped, err := url.QueryUnescape(qp.Get("$top"))
		if err != nil {
			return nil, err
		}
		topParam, err := parseOptional(topUnescaped, func(v string) (int32, error) {
			p, parseErr := strconv.ParseInt(v, 10, 32)
			if parseErr != nil {
				return 0, parseErr
			}
			return int32(p), nil
		})
		if err != nil {
			return nil, err
		}
		var options *armhardwaresecuritymodules.DedicatedHsmClientListBySubscriptionOptions
		if topParam != nil {
			options = &armhardwaresecuritymodules.DedicatedHsmClientListBySubscriptionOptions{
				Top: topParam,
			}
		}
		resp := d.srv.NewListBySubscriptionPager(options)
		newListBySubscriptionPager = &resp
		d.newListBySubscriptionPager.add(req, newListBySubscriptionPager)
		server.PagerResponderInjectNextLinks(newListBySubscriptionPager, req, func(page *armhardwaresecuritymodules.DedicatedHsmClientListBySubscriptionResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListBySubscriptionPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		d.newListBySubscriptionPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListBySubscriptionPager) {
		d.newListBySubscriptionPager.remove(req)
	}
	return resp, nil
}

func (d *DedicatedHsmServerTransport) dispatchNewListOutboundNetworkDependenciesEndpointsPager(req *http.Request) (*http.Response, error) {
	if d.srv.NewListOutboundNetworkDependenciesEndpointsPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListOutboundNetworkDependenciesEndpointsPager not implemented")}
	}
	newListOutboundNetworkDependenciesEndpointsPager := d.newListOutboundNetworkDependenciesEndpointsPager.get(req)
	if newListOutboundNetworkDependenciesEndpointsPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HardwareSecurityModules/dedicatedHSMs/(?P<name>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/outboundNetworkDependenciesEndpoints`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		nameParam, err := url.PathUnescape(matches[regex.SubexpIndex("name")])
		if err != nil {
			return nil, err
		}
		resp := d.srv.NewListOutboundNetworkDependenciesEndpointsPager(resourceGroupNameParam, nameParam, nil)
		newListOutboundNetworkDependenciesEndpointsPager = &resp
		d.newListOutboundNetworkDependenciesEndpointsPager.add(req, newListOutboundNetworkDependenciesEndpointsPager)
		server.PagerResponderInjectNextLinks(newListOutboundNetworkDependenciesEndpointsPager, req, func(page *armhardwaresecuritymodules.DedicatedHsmClientListOutboundNetworkDependenciesEndpointsResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListOutboundNetworkDependenciesEndpointsPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		d.newListOutboundNetworkDependenciesEndpointsPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListOutboundNetworkDependenciesEndpointsPager) {
		d.newListOutboundNetworkDependenciesEndpointsPager.remove(req)
	}
	return resp, nil
}

func (d *DedicatedHsmServerTransport) dispatchBeginUpdate(req *http.Request) (*http.Response, error) {
	if d.srv.BeginUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdate not implemented")}
	}
	beginUpdate := d.beginUpdate.get(req)
	if beginUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HardwareSecurityModules/dedicatedHSMs/(?P<name>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armhardwaresecuritymodules.DedicatedHsmPatchParameters](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		nameParam, err := url.PathUnescape(matches[regex.SubexpIndex("name")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := d.srv.BeginUpdate(req.Context(), resourceGroupNameParam, nameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginUpdate = &respr
		d.beginUpdate.add(req, beginUpdate)
	}

	resp, err := server.PollerResponderNext(beginUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		d.beginUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdate) {
		d.beginUpdate.remove(req)
	}

	return resp, nil
}
