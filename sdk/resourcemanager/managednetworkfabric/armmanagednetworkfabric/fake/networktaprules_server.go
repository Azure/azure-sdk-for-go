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

// NetworkTapRulesServer is a fake server for instances of the armmanagednetworkfabric.NetworkTapRulesClient type.
type NetworkTapRulesServer struct {
	// BeginCreate is the fake for method NetworkTapRulesClient.BeginCreate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreate func(ctx context.Context, resourceGroupName string, networkTapRuleName string, body armmanagednetworkfabric.NetworkTapRule, options *armmanagednetworkfabric.NetworkTapRulesClientBeginCreateOptions) (resp azfake.PollerResponder[armmanagednetworkfabric.NetworkTapRulesClientCreateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method NetworkTapRulesClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, networkTapRuleName string, options *armmanagednetworkfabric.NetworkTapRulesClientBeginDeleteOptions) (resp azfake.PollerResponder[armmanagednetworkfabric.NetworkTapRulesClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method NetworkTapRulesClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, networkTapRuleName string, options *armmanagednetworkfabric.NetworkTapRulesClientGetOptions) (resp azfake.Responder[armmanagednetworkfabric.NetworkTapRulesClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByResourceGroupPager is the fake for method NetworkTapRulesClient.NewListByResourceGroupPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByResourceGroupPager func(resourceGroupName string, options *armmanagednetworkfabric.NetworkTapRulesClientListByResourceGroupOptions) (resp azfake.PagerResponder[armmanagednetworkfabric.NetworkTapRulesClientListByResourceGroupResponse])

	// NewListBySubscriptionPager is the fake for method NetworkTapRulesClient.NewListBySubscriptionPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListBySubscriptionPager func(options *armmanagednetworkfabric.NetworkTapRulesClientListBySubscriptionOptions) (resp azfake.PagerResponder[armmanagednetworkfabric.NetworkTapRulesClientListBySubscriptionResponse])

	// BeginResync is the fake for method NetworkTapRulesClient.BeginResync
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginResync func(ctx context.Context, resourceGroupName string, networkTapRuleName string, options *armmanagednetworkfabric.NetworkTapRulesClientBeginResyncOptions) (resp azfake.PollerResponder[armmanagednetworkfabric.NetworkTapRulesClientResyncResponse], errResp azfake.ErrorResponder)

	// BeginUpdate is the fake for method NetworkTapRulesClient.BeginUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginUpdate func(ctx context.Context, resourceGroupName string, networkTapRuleName string, body armmanagednetworkfabric.NetworkTapRulePatch, options *armmanagednetworkfabric.NetworkTapRulesClientBeginUpdateOptions) (resp azfake.PollerResponder[armmanagednetworkfabric.NetworkTapRulesClientUpdateResponse], errResp azfake.ErrorResponder)

	// BeginUpdateAdministrativeState is the fake for method NetworkTapRulesClient.BeginUpdateAdministrativeState
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginUpdateAdministrativeState func(ctx context.Context, resourceGroupName string, networkTapRuleName string, body armmanagednetworkfabric.UpdateAdministrativeState, options *armmanagednetworkfabric.NetworkTapRulesClientBeginUpdateAdministrativeStateOptions) (resp azfake.PollerResponder[armmanagednetworkfabric.NetworkTapRulesClientUpdateAdministrativeStateResponse], errResp azfake.ErrorResponder)

	// BeginValidateConfiguration is the fake for method NetworkTapRulesClient.BeginValidateConfiguration
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginValidateConfiguration func(ctx context.Context, resourceGroupName string, networkTapRuleName string, options *armmanagednetworkfabric.NetworkTapRulesClientBeginValidateConfigurationOptions) (resp azfake.PollerResponder[armmanagednetworkfabric.NetworkTapRulesClientValidateConfigurationResponse], errResp azfake.ErrorResponder)
}

// NewNetworkTapRulesServerTransport creates a new instance of NetworkTapRulesServerTransport with the provided implementation.
// The returned NetworkTapRulesServerTransport instance is connected to an instance of armmanagednetworkfabric.NetworkTapRulesClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewNetworkTapRulesServerTransport(srv *NetworkTapRulesServer) *NetworkTapRulesServerTransport {
	return &NetworkTapRulesServerTransport{
		srv:                            srv,
		beginCreate:                    newTracker[azfake.PollerResponder[armmanagednetworkfabric.NetworkTapRulesClientCreateResponse]](),
		beginDelete:                    newTracker[azfake.PollerResponder[armmanagednetworkfabric.NetworkTapRulesClientDeleteResponse]](),
		newListByResourceGroupPager:    newTracker[azfake.PagerResponder[armmanagednetworkfabric.NetworkTapRulesClientListByResourceGroupResponse]](),
		newListBySubscriptionPager:     newTracker[azfake.PagerResponder[armmanagednetworkfabric.NetworkTapRulesClientListBySubscriptionResponse]](),
		beginResync:                    newTracker[azfake.PollerResponder[armmanagednetworkfabric.NetworkTapRulesClientResyncResponse]](),
		beginUpdate:                    newTracker[azfake.PollerResponder[armmanagednetworkfabric.NetworkTapRulesClientUpdateResponse]](),
		beginUpdateAdministrativeState: newTracker[azfake.PollerResponder[armmanagednetworkfabric.NetworkTapRulesClientUpdateAdministrativeStateResponse]](),
		beginValidateConfiguration:     newTracker[azfake.PollerResponder[armmanagednetworkfabric.NetworkTapRulesClientValidateConfigurationResponse]](),
	}
}

// NetworkTapRulesServerTransport connects instances of armmanagednetworkfabric.NetworkTapRulesClient to instances of NetworkTapRulesServer.
// Don't use this type directly, use NewNetworkTapRulesServerTransport instead.
type NetworkTapRulesServerTransport struct {
	srv                            *NetworkTapRulesServer
	beginCreate                    *tracker[azfake.PollerResponder[armmanagednetworkfabric.NetworkTapRulesClientCreateResponse]]
	beginDelete                    *tracker[azfake.PollerResponder[armmanagednetworkfabric.NetworkTapRulesClientDeleteResponse]]
	newListByResourceGroupPager    *tracker[azfake.PagerResponder[armmanagednetworkfabric.NetworkTapRulesClientListByResourceGroupResponse]]
	newListBySubscriptionPager     *tracker[azfake.PagerResponder[armmanagednetworkfabric.NetworkTapRulesClientListBySubscriptionResponse]]
	beginResync                    *tracker[azfake.PollerResponder[armmanagednetworkfabric.NetworkTapRulesClientResyncResponse]]
	beginUpdate                    *tracker[azfake.PollerResponder[armmanagednetworkfabric.NetworkTapRulesClientUpdateResponse]]
	beginUpdateAdministrativeState *tracker[azfake.PollerResponder[armmanagednetworkfabric.NetworkTapRulesClientUpdateAdministrativeStateResponse]]
	beginValidateConfiguration     *tracker[azfake.PollerResponder[armmanagednetworkfabric.NetworkTapRulesClientValidateConfigurationResponse]]
}

// Do implements the policy.Transporter interface for NetworkTapRulesServerTransport.
func (n *NetworkTapRulesServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "NetworkTapRulesClient.BeginCreate":
		resp, err = n.dispatchBeginCreate(req)
	case "NetworkTapRulesClient.BeginDelete":
		resp, err = n.dispatchBeginDelete(req)
	case "NetworkTapRulesClient.Get":
		resp, err = n.dispatchGet(req)
	case "NetworkTapRulesClient.NewListByResourceGroupPager":
		resp, err = n.dispatchNewListByResourceGroupPager(req)
	case "NetworkTapRulesClient.NewListBySubscriptionPager":
		resp, err = n.dispatchNewListBySubscriptionPager(req)
	case "NetworkTapRulesClient.BeginResync":
		resp, err = n.dispatchBeginResync(req)
	case "NetworkTapRulesClient.BeginUpdate":
		resp, err = n.dispatchBeginUpdate(req)
	case "NetworkTapRulesClient.BeginUpdateAdministrativeState":
		resp, err = n.dispatchBeginUpdateAdministrativeState(req)
	case "NetworkTapRulesClient.BeginValidateConfiguration":
		resp, err = n.dispatchBeginValidateConfiguration(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (n *NetworkTapRulesServerTransport) dispatchBeginCreate(req *http.Request) (*http.Response, error) {
	if n.srv.BeginCreate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreate not implemented")}
	}
	beginCreate := n.beginCreate.get(req)
	if beginCreate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ManagedNetworkFabric/networkTapRules/(?P<networkTapRuleName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armmanagednetworkfabric.NetworkTapRule](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		networkTapRuleNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("networkTapRuleName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := n.srv.BeginCreate(req.Context(), resourceGroupNameParam, networkTapRuleNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreate = &respr
		n.beginCreate.add(req, beginCreate)
	}

	resp, err := server.PollerResponderNext(beginCreate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		n.beginCreate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreate) {
		n.beginCreate.remove(req)
	}

	return resp, nil
}

func (n *NetworkTapRulesServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if n.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := n.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ManagedNetworkFabric/networkTapRules/(?P<networkTapRuleName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		networkTapRuleNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("networkTapRuleName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := n.srv.BeginDelete(req.Context(), resourceGroupNameParam, networkTapRuleNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		n.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		n.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		n.beginDelete.remove(req)
	}

	return resp, nil
}

func (n *NetworkTapRulesServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if n.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ManagedNetworkFabric/networkTapRules/(?P<networkTapRuleName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	networkTapRuleNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("networkTapRuleName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := n.srv.Get(req.Context(), resourceGroupNameParam, networkTapRuleNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).NetworkTapRule, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (n *NetworkTapRulesServerTransport) dispatchNewListByResourceGroupPager(req *http.Request) (*http.Response, error) {
	if n.srv.NewListByResourceGroupPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByResourceGroupPager not implemented")}
	}
	newListByResourceGroupPager := n.newListByResourceGroupPager.get(req)
	if newListByResourceGroupPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ManagedNetworkFabric/networkTapRules`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		resp := n.srv.NewListByResourceGroupPager(resourceGroupNameParam, nil)
		newListByResourceGroupPager = &resp
		n.newListByResourceGroupPager.add(req, newListByResourceGroupPager)
		server.PagerResponderInjectNextLinks(newListByResourceGroupPager, req, func(page *armmanagednetworkfabric.NetworkTapRulesClientListByResourceGroupResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByResourceGroupPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		n.newListByResourceGroupPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByResourceGroupPager) {
		n.newListByResourceGroupPager.remove(req)
	}
	return resp, nil
}

func (n *NetworkTapRulesServerTransport) dispatchNewListBySubscriptionPager(req *http.Request) (*http.Response, error) {
	if n.srv.NewListBySubscriptionPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListBySubscriptionPager not implemented")}
	}
	newListBySubscriptionPager := n.newListBySubscriptionPager.get(req)
	if newListBySubscriptionPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ManagedNetworkFabric/networkTapRules`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resp := n.srv.NewListBySubscriptionPager(nil)
		newListBySubscriptionPager = &resp
		n.newListBySubscriptionPager.add(req, newListBySubscriptionPager)
		server.PagerResponderInjectNextLinks(newListBySubscriptionPager, req, func(page *armmanagednetworkfabric.NetworkTapRulesClientListBySubscriptionResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListBySubscriptionPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		n.newListBySubscriptionPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListBySubscriptionPager) {
		n.newListBySubscriptionPager.remove(req)
	}
	return resp, nil
}

func (n *NetworkTapRulesServerTransport) dispatchBeginResync(req *http.Request) (*http.Response, error) {
	if n.srv.BeginResync == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginResync not implemented")}
	}
	beginResync := n.beginResync.get(req)
	if beginResync == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ManagedNetworkFabric/networkTapRules/(?P<networkTapRuleName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resync`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		networkTapRuleNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("networkTapRuleName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := n.srv.BeginResync(req.Context(), resourceGroupNameParam, networkTapRuleNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginResync = &respr
		n.beginResync.add(req, beginResync)
	}

	resp, err := server.PollerResponderNext(beginResync, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		n.beginResync.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginResync) {
		n.beginResync.remove(req)
	}

	return resp, nil
}

func (n *NetworkTapRulesServerTransport) dispatchBeginUpdate(req *http.Request) (*http.Response, error) {
	if n.srv.BeginUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdate not implemented")}
	}
	beginUpdate := n.beginUpdate.get(req)
	if beginUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ManagedNetworkFabric/networkTapRules/(?P<networkTapRuleName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armmanagednetworkfabric.NetworkTapRulePatch](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		networkTapRuleNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("networkTapRuleName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := n.srv.BeginUpdate(req.Context(), resourceGroupNameParam, networkTapRuleNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginUpdate = &respr
		n.beginUpdate.add(req, beginUpdate)
	}

	resp, err := server.PollerResponderNext(beginUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		n.beginUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdate) {
		n.beginUpdate.remove(req)
	}

	return resp, nil
}

func (n *NetworkTapRulesServerTransport) dispatchBeginUpdateAdministrativeState(req *http.Request) (*http.Response, error) {
	if n.srv.BeginUpdateAdministrativeState == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdateAdministrativeState not implemented")}
	}
	beginUpdateAdministrativeState := n.beginUpdateAdministrativeState.get(req)
	if beginUpdateAdministrativeState == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ManagedNetworkFabric/networkTapRules/(?P<networkTapRuleName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/updateAdministrativeState`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
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
		networkTapRuleNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("networkTapRuleName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := n.srv.BeginUpdateAdministrativeState(req.Context(), resourceGroupNameParam, networkTapRuleNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginUpdateAdministrativeState = &respr
		n.beginUpdateAdministrativeState.add(req, beginUpdateAdministrativeState)
	}

	resp, err := server.PollerResponderNext(beginUpdateAdministrativeState, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		n.beginUpdateAdministrativeState.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdateAdministrativeState) {
		n.beginUpdateAdministrativeState.remove(req)
	}

	return resp, nil
}

func (n *NetworkTapRulesServerTransport) dispatchBeginValidateConfiguration(req *http.Request) (*http.Response, error) {
	if n.srv.BeginValidateConfiguration == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginValidateConfiguration not implemented")}
	}
	beginValidateConfiguration := n.beginValidateConfiguration.get(req)
	if beginValidateConfiguration == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ManagedNetworkFabric/networkTapRules/(?P<networkTapRuleName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/validateConfiguration`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		networkTapRuleNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("networkTapRuleName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := n.srv.BeginValidateConfiguration(req.Context(), resourceGroupNameParam, networkTapRuleNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginValidateConfiguration = &respr
		n.beginValidateConfiguration.add(req, beginValidateConfiguration)
	}

	resp, err := server.PollerResponderNext(beginValidateConfiguration, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		n.beginValidateConfiguration.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginValidateConfiguration) {
		n.beginValidateConfiguration.remove(req)
	}

	return resp, nil
}