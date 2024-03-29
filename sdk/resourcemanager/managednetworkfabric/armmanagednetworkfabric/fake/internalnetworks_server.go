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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/managednetworkfabric/armmanagednetworkfabric"
	"net/http"
	"net/url"
	"regexp"
)

// InternalNetworksServer is a fake server for instances of the armmanagednetworkfabric.InternalNetworksClient type.
type InternalNetworksServer struct {
	// BeginCreate is the fake for method InternalNetworksClient.BeginCreate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreate func(ctx context.Context, resourceGroupName string, l3IsolationDomainName string, internalNetworkName string, body armmanagednetworkfabric.InternalNetwork, options *armmanagednetworkfabric.InternalNetworksClientBeginCreateOptions) (resp azfake.PollerResponder[armmanagednetworkfabric.InternalNetworksClientCreateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method InternalNetworksClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, l3IsolationDomainName string, internalNetworkName string, options *armmanagednetworkfabric.InternalNetworksClientBeginDeleteOptions) (resp azfake.PollerResponder[armmanagednetworkfabric.InternalNetworksClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method InternalNetworksClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, l3IsolationDomainName string, internalNetworkName string, options *armmanagednetworkfabric.InternalNetworksClientGetOptions) (resp azfake.Responder[armmanagednetworkfabric.InternalNetworksClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByL3IsolationDomainPager is the fake for method InternalNetworksClient.NewListByL3IsolationDomainPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByL3IsolationDomainPager func(resourceGroupName string, l3IsolationDomainName string, options *armmanagednetworkfabric.InternalNetworksClientListByL3IsolationDomainOptions) (resp azfake.PagerResponder[armmanagednetworkfabric.InternalNetworksClientListByL3IsolationDomainResponse])

	// BeginUpdate is the fake for method InternalNetworksClient.BeginUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginUpdate func(ctx context.Context, resourceGroupName string, l3IsolationDomainName string, internalNetworkName string, body armmanagednetworkfabric.InternalNetworkPatch, options *armmanagednetworkfabric.InternalNetworksClientBeginUpdateOptions) (resp azfake.PollerResponder[armmanagednetworkfabric.InternalNetworksClientUpdateResponse], errResp azfake.ErrorResponder)

	// BeginUpdateAdministrativeState is the fake for method InternalNetworksClient.BeginUpdateAdministrativeState
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginUpdateAdministrativeState func(ctx context.Context, resourceGroupName string, l3IsolationDomainName string, internalNetworkName string, body armmanagednetworkfabric.UpdateAdministrativeState, options *armmanagednetworkfabric.InternalNetworksClientBeginUpdateAdministrativeStateOptions) (resp azfake.PollerResponder[armmanagednetworkfabric.InternalNetworksClientUpdateAdministrativeStateResponse], errResp azfake.ErrorResponder)

	// BeginUpdateBgpAdministrativeState is the fake for method InternalNetworksClient.BeginUpdateBgpAdministrativeState
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginUpdateBgpAdministrativeState func(ctx context.Context, resourceGroupName string, l3IsolationDomainName string, internalNetworkName string, body armmanagednetworkfabric.UpdateAdministrativeState, options *armmanagednetworkfabric.InternalNetworksClientBeginUpdateBgpAdministrativeStateOptions) (resp azfake.PollerResponder[armmanagednetworkfabric.InternalNetworksClientUpdateBgpAdministrativeStateResponse], errResp azfake.ErrorResponder)

	// BeginUpdateStaticRouteBfdAdministrativeState is the fake for method InternalNetworksClient.BeginUpdateStaticRouteBfdAdministrativeState
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginUpdateStaticRouteBfdAdministrativeState func(ctx context.Context, resourceGroupName string, l3IsolationDomainName string, internalNetworkName string, body armmanagednetworkfabric.UpdateAdministrativeState, options *armmanagednetworkfabric.InternalNetworksClientBeginUpdateStaticRouteBfdAdministrativeStateOptions) (resp azfake.PollerResponder[armmanagednetworkfabric.InternalNetworksClientUpdateStaticRouteBfdAdministrativeStateResponse], errResp azfake.ErrorResponder)
}

// NewInternalNetworksServerTransport creates a new instance of InternalNetworksServerTransport with the provided implementation.
// The returned InternalNetworksServerTransport instance is connected to an instance of armmanagednetworkfabric.InternalNetworksClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewInternalNetworksServerTransport(srv *InternalNetworksServer) *InternalNetworksServerTransport {
	return &InternalNetworksServerTransport{
		srv:                               srv,
		beginCreate:                       newTracker[azfake.PollerResponder[armmanagednetworkfabric.InternalNetworksClientCreateResponse]](),
		beginDelete:                       newTracker[azfake.PollerResponder[armmanagednetworkfabric.InternalNetworksClientDeleteResponse]](),
		newListByL3IsolationDomainPager:   newTracker[azfake.PagerResponder[armmanagednetworkfabric.InternalNetworksClientListByL3IsolationDomainResponse]](),
		beginUpdate:                       newTracker[azfake.PollerResponder[armmanagednetworkfabric.InternalNetworksClientUpdateResponse]](),
		beginUpdateAdministrativeState:    newTracker[azfake.PollerResponder[armmanagednetworkfabric.InternalNetworksClientUpdateAdministrativeStateResponse]](),
		beginUpdateBgpAdministrativeState: newTracker[azfake.PollerResponder[armmanagednetworkfabric.InternalNetworksClientUpdateBgpAdministrativeStateResponse]](),
		beginUpdateStaticRouteBfdAdministrativeState: newTracker[azfake.PollerResponder[armmanagednetworkfabric.InternalNetworksClientUpdateStaticRouteBfdAdministrativeStateResponse]](),
	}
}

// InternalNetworksServerTransport connects instances of armmanagednetworkfabric.InternalNetworksClient to instances of InternalNetworksServer.
// Don't use this type directly, use NewInternalNetworksServerTransport instead.
type InternalNetworksServerTransport struct {
	srv                                          *InternalNetworksServer
	beginCreate                                  *tracker[azfake.PollerResponder[armmanagednetworkfabric.InternalNetworksClientCreateResponse]]
	beginDelete                                  *tracker[azfake.PollerResponder[armmanagednetworkfabric.InternalNetworksClientDeleteResponse]]
	newListByL3IsolationDomainPager              *tracker[azfake.PagerResponder[armmanagednetworkfabric.InternalNetworksClientListByL3IsolationDomainResponse]]
	beginUpdate                                  *tracker[azfake.PollerResponder[armmanagednetworkfabric.InternalNetworksClientUpdateResponse]]
	beginUpdateAdministrativeState               *tracker[azfake.PollerResponder[armmanagednetworkfabric.InternalNetworksClientUpdateAdministrativeStateResponse]]
	beginUpdateBgpAdministrativeState            *tracker[azfake.PollerResponder[armmanagednetworkfabric.InternalNetworksClientUpdateBgpAdministrativeStateResponse]]
	beginUpdateStaticRouteBfdAdministrativeState *tracker[azfake.PollerResponder[armmanagednetworkfabric.InternalNetworksClientUpdateStaticRouteBfdAdministrativeStateResponse]]
}

// Do implements the policy.Transporter interface for InternalNetworksServerTransport.
func (i *InternalNetworksServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "InternalNetworksClient.BeginCreate":
		resp, err = i.dispatchBeginCreate(req)
	case "InternalNetworksClient.BeginDelete":
		resp, err = i.dispatchBeginDelete(req)
	case "InternalNetworksClient.Get":
		resp, err = i.dispatchGet(req)
	case "InternalNetworksClient.NewListByL3IsolationDomainPager":
		resp, err = i.dispatchNewListByL3IsolationDomainPager(req)
	case "InternalNetworksClient.BeginUpdate":
		resp, err = i.dispatchBeginUpdate(req)
	case "InternalNetworksClient.BeginUpdateAdministrativeState":
		resp, err = i.dispatchBeginUpdateAdministrativeState(req)
	case "InternalNetworksClient.BeginUpdateBgpAdministrativeState":
		resp, err = i.dispatchBeginUpdateBgpAdministrativeState(req)
	case "InternalNetworksClient.BeginUpdateStaticRouteBfdAdministrativeState":
		resp, err = i.dispatchBeginUpdateStaticRouteBfdAdministrativeState(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (i *InternalNetworksServerTransport) dispatchBeginCreate(req *http.Request) (*http.Response, error) {
	if i.srv.BeginCreate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreate not implemented")}
	}
	beginCreate := i.beginCreate.get(req)
	if beginCreate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ManagedNetworkFabric/l3IsolationDomains/(?P<l3IsolationDomainName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/internalNetworks/(?P<internalNetworkName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armmanagednetworkfabric.InternalNetwork](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		l3IsolationDomainNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("l3IsolationDomainName")])
		if err != nil {
			return nil, err
		}
		internalNetworkNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("internalNetworkName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := i.srv.BeginCreate(req.Context(), resourceGroupNameParam, l3IsolationDomainNameParam, internalNetworkNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreate = &respr
		i.beginCreate.add(req, beginCreate)
	}

	resp, err := server.PollerResponderNext(beginCreate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		i.beginCreate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreate) {
		i.beginCreate.remove(req)
	}

	return resp, nil
}

func (i *InternalNetworksServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if i.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := i.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ManagedNetworkFabric/l3IsolationDomains/(?P<l3IsolationDomainName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/internalNetworks/(?P<internalNetworkName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		l3IsolationDomainNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("l3IsolationDomainName")])
		if err != nil {
			return nil, err
		}
		internalNetworkNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("internalNetworkName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := i.srv.BeginDelete(req.Context(), resourceGroupNameParam, l3IsolationDomainNameParam, internalNetworkNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		i.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		i.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		i.beginDelete.remove(req)
	}

	return resp, nil
}

func (i *InternalNetworksServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if i.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ManagedNetworkFabric/l3IsolationDomains/(?P<l3IsolationDomainName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/internalNetworks/(?P<internalNetworkName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	l3IsolationDomainNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("l3IsolationDomainName")])
	if err != nil {
		return nil, err
	}
	internalNetworkNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("internalNetworkName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := i.srv.Get(req.Context(), resourceGroupNameParam, l3IsolationDomainNameParam, internalNetworkNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).InternalNetwork, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *InternalNetworksServerTransport) dispatchNewListByL3IsolationDomainPager(req *http.Request) (*http.Response, error) {
	if i.srv.NewListByL3IsolationDomainPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByL3IsolationDomainPager not implemented")}
	}
	newListByL3IsolationDomainPager := i.newListByL3IsolationDomainPager.get(req)
	if newListByL3IsolationDomainPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ManagedNetworkFabric/l3IsolationDomains/(?P<l3IsolationDomainName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/internalNetworks`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		l3IsolationDomainNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("l3IsolationDomainName")])
		if err != nil {
			return nil, err
		}
		resp := i.srv.NewListByL3IsolationDomainPager(resourceGroupNameParam, l3IsolationDomainNameParam, nil)
		newListByL3IsolationDomainPager = &resp
		i.newListByL3IsolationDomainPager.add(req, newListByL3IsolationDomainPager)
		server.PagerResponderInjectNextLinks(newListByL3IsolationDomainPager, req, func(page *armmanagednetworkfabric.InternalNetworksClientListByL3IsolationDomainResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByL3IsolationDomainPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		i.newListByL3IsolationDomainPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByL3IsolationDomainPager) {
		i.newListByL3IsolationDomainPager.remove(req)
	}
	return resp, nil
}

func (i *InternalNetworksServerTransport) dispatchBeginUpdate(req *http.Request) (*http.Response, error) {
	if i.srv.BeginUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdate not implemented")}
	}
	beginUpdate := i.beginUpdate.get(req)
	if beginUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ManagedNetworkFabric/l3IsolationDomains/(?P<l3IsolationDomainName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/internalNetworks/(?P<internalNetworkName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armmanagednetworkfabric.InternalNetworkPatch](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		l3IsolationDomainNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("l3IsolationDomainName")])
		if err != nil {
			return nil, err
		}
		internalNetworkNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("internalNetworkName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := i.srv.BeginUpdate(req.Context(), resourceGroupNameParam, l3IsolationDomainNameParam, internalNetworkNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginUpdate = &respr
		i.beginUpdate.add(req, beginUpdate)
	}

	resp, err := server.PollerResponderNext(beginUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		i.beginUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdate) {
		i.beginUpdate.remove(req)
	}

	return resp, nil
}

func (i *InternalNetworksServerTransport) dispatchBeginUpdateAdministrativeState(req *http.Request) (*http.Response, error) {
	if i.srv.BeginUpdateAdministrativeState == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdateAdministrativeState not implemented")}
	}
	beginUpdateAdministrativeState := i.beginUpdateAdministrativeState.get(req)
	if beginUpdateAdministrativeState == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ManagedNetworkFabric/l3IsolationDomains/(?P<l3IsolationDomainName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/internalNetworks/(?P<internalNetworkName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/updateAdministrativeState`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armmanagednetworkfabric.UpdateAdministrativeState](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		l3IsolationDomainNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("l3IsolationDomainName")])
		if err != nil {
			return nil, err
		}
		internalNetworkNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("internalNetworkName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := i.srv.BeginUpdateAdministrativeState(req.Context(), resourceGroupNameParam, l3IsolationDomainNameParam, internalNetworkNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginUpdateAdministrativeState = &respr
		i.beginUpdateAdministrativeState.add(req, beginUpdateAdministrativeState)
	}

	resp, err := server.PollerResponderNext(beginUpdateAdministrativeState, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		i.beginUpdateAdministrativeState.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdateAdministrativeState) {
		i.beginUpdateAdministrativeState.remove(req)
	}

	return resp, nil
}

func (i *InternalNetworksServerTransport) dispatchBeginUpdateBgpAdministrativeState(req *http.Request) (*http.Response, error) {
	if i.srv.BeginUpdateBgpAdministrativeState == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdateBgpAdministrativeState not implemented")}
	}
	beginUpdateBgpAdministrativeState := i.beginUpdateBgpAdministrativeState.get(req)
	if beginUpdateBgpAdministrativeState == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ManagedNetworkFabric/l3IsolationDomains/(?P<l3IsolationDomainName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/internalNetworks/(?P<internalNetworkName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/updateBgpAdministrativeState`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armmanagednetworkfabric.UpdateAdministrativeState](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		l3IsolationDomainNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("l3IsolationDomainName")])
		if err != nil {
			return nil, err
		}
		internalNetworkNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("internalNetworkName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := i.srv.BeginUpdateBgpAdministrativeState(req.Context(), resourceGroupNameParam, l3IsolationDomainNameParam, internalNetworkNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginUpdateBgpAdministrativeState = &respr
		i.beginUpdateBgpAdministrativeState.add(req, beginUpdateBgpAdministrativeState)
	}

	resp, err := server.PollerResponderNext(beginUpdateBgpAdministrativeState, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		i.beginUpdateBgpAdministrativeState.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdateBgpAdministrativeState) {
		i.beginUpdateBgpAdministrativeState.remove(req)
	}

	return resp, nil
}

func (i *InternalNetworksServerTransport) dispatchBeginUpdateStaticRouteBfdAdministrativeState(req *http.Request) (*http.Response, error) {
	if i.srv.BeginUpdateStaticRouteBfdAdministrativeState == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdateStaticRouteBfdAdministrativeState not implemented")}
	}
	beginUpdateStaticRouteBfdAdministrativeState := i.beginUpdateStaticRouteBfdAdministrativeState.get(req)
	if beginUpdateStaticRouteBfdAdministrativeState == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ManagedNetworkFabric/l3IsolationDomains/(?P<l3IsolationDomainName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/internalNetworks/(?P<internalNetworkName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/updateStaticRouteBfdAdministrativeState`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armmanagednetworkfabric.UpdateAdministrativeState](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		l3IsolationDomainNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("l3IsolationDomainName")])
		if err != nil {
			return nil, err
		}
		internalNetworkNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("internalNetworkName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := i.srv.BeginUpdateStaticRouteBfdAdministrativeState(req.Context(), resourceGroupNameParam, l3IsolationDomainNameParam, internalNetworkNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginUpdateStaticRouteBfdAdministrativeState = &respr
		i.beginUpdateStaticRouteBfdAdministrativeState.add(req, beginUpdateStaticRouteBfdAdministrativeState)
	}

	resp, err := server.PollerResponderNext(beginUpdateStaticRouteBfdAdministrativeState, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		i.beginUpdateStaticRouteBfdAdministrativeState.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdateStaticRouteBfdAdministrativeState) {
		i.beginUpdateStaticRouteBfdAdministrativeState.remove(req)
	}

	return resp, nil
}
