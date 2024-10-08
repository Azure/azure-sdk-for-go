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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/machinelearning/armmachinelearning/v4"
	"net/http"
	"net/url"
	"regexp"
)

// ManagedNetworkSettingsRuleServer is a fake server for instances of the armmachinelearning.ManagedNetworkSettingsRuleClient type.
type ManagedNetworkSettingsRuleServer struct {
	// BeginCreateOrUpdate is the fake for method ManagedNetworkSettingsRuleClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, workspaceName string, ruleName string, body armmachinelearning.OutboundRuleBasicResource, options *armmachinelearning.ManagedNetworkSettingsRuleClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armmachinelearning.ManagedNetworkSettingsRuleClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method ManagedNetworkSettingsRuleClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, workspaceName string, ruleName string, options *armmachinelearning.ManagedNetworkSettingsRuleClientBeginDeleteOptions) (resp azfake.PollerResponder[armmachinelearning.ManagedNetworkSettingsRuleClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method ManagedNetworkSettingsRuleClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, workspaceName string, ruleName string, options *armmachinelearning.ManagedNetworkSettingsRuleClientGetOptions) (resp azfake.Responder[armmachinelearning.ManagedNetworkSettingsRuleClientGetResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method ManagedNetworkSettingsRuleClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(resourceGroupName string, workspaceName string, options *armmachinelearning.ManagedNetworkSettingsRuleClientListOptions) (resp azfake.PagerResponder[armmachinelearning.ManagedNetworkSettingsRuleClientListResponse])
}

// NewManagedNetworkSettingsRuleServerTransport creates a new instance of ManagedNetworkSettingsRuleServerTransport with the provided implementation.
// The returned ManagedNetworkSettingsRuleServerTransport instance is connected to an instance of armmachinelearning.ManagedNetworkSettingsRuleClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewManagedNetworkSettingsRuleServerTransport(srv *ManagedNetworkSettingsRuleServer) *ManagedNetworkSettingsRuleServerTransport {
	return &ManagedNetworkSettingsRuleServerTransport{
		srv:                 srv,
		beginCreateOrUpdate: newTracker[azfake.PollerResponder[armmachinelearning.ManagedNetworkSettingsRuleClientCreateOrUpdateResponse]](),
		beginDelete:         newTracker[azfake.PollerResponder[armmachinelearning.ManagedNetworkSettingsRuleClientDeleteResponse]](),
		newListPager:        newTracker[azfake.PagerResponder[armmachinelearning.ManagedNetworkSettingsRuleClientListResponse]](),
	}
}

// ManagedNetworkSettingsRuleServerTransport connects instances of armmachinelearning.ManagedNetworkSettingsRuleClient to instances of ManagedNetworkSettingsRuleServer.
// Don't use this type directly, use NewManagedNetworkSettingsRuleServerTransport instead.
type ManagedNetworkSettingsRuleServerTransport struct {
	srv                 *ManagedNetworkSettingsRuleServer
	beginCreateOrUpdate *tracker[azfake.PollerResponder[armmachinelearning.ManagedNetworkSettingsRuleClientCreateOrUpdateResponse]]
	beginDelete         *tracker[azfake.PollerResponder[armmachinelearning.ManagedNetworkSettingsRuleClientDeleteResponse]]
	newListPager        *tracker[azfake.PagerResponder[armmachinelearning.ManagedNetworkSettingsRuleClientListResponse]]
}

// Do implements the policy.Transporter interface for ManagedNetworkSettingsRuleServerTransport.
func (m *ManagedNetworkSettingsRuleServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "ManagedNetworkSettingsRuleClient.BeginCreateOrUpdate":
		resp, err = m.dispatchBeginCreateOrUpdate(req)
	case "ManagedNetworkSettingsRuleClient.BeginDelete":
		resp, err = m.dispatchBeginDelete(req)
	case "ManagedNetworkSettingsRuleClient.Get":
		resp, err = m.dispatchGet(req)
	case "ManagedNetworkSettingsRuleClient.NewListPager":
		resp, err = m.dispatchNewListPager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *ManagedNetworkSettingsRuleServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if m.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := m.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.MachineLearningServices/workspaces/(?P<workspaceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/outboundRules/(?P<ruleName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armmachinelearning.OutboundRuleBasicResource](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		workspaceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("workspaceName")])
		if err != nil {
			return nil, err
		}
		ruleNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("ruleName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := m.srv.BeginCreateOrUpdate(req.Context(), resourceGroupNameParam, workspaceNameParam, ruleNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateOrUpdate = &respr
		m.beginCreateOrUpdate.add(req, beginCreateOrUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		m.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		m.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (m *ManagedNetworkSettingsRuleServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if m.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := m.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.MachineLearningServices/workspaces/(?P<workspaceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/outboundRules/(?P<ruleName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		workspaceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("workspaceName")])
		if err != nil {
			return nil, err
		}
		ruleNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("ruleName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := m.srv.BeginDelete(req.Context(), resourceGroupNameParam, workspaceNameParam, ruleNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		m.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		m.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		m.beginDelete.remove(req)
	}

	return resp, nil
}

func (m *ManagedNetworkSettingsRuleServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if m.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.MachineLearningServices/workspaces/(?P<workspaceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/outboundRules/(?P<ruleName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	workspaceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("workspaceName")])
	if err != nil {
		return nil, err
	}
	ruleNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("ruleName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := m.srv.Get(req.Context(), resourceGroupNameParam, workspaceNameParam, ruleNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).OutboundRuleBasicResource, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *ManagedNetworkSettingsRuleServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if m.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := m.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.MachineLearningServices/workspaces/(?P<workspaceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/outboundRules`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		workspaceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("workspaceName")])
		if err != nil {
			return nil, err
		}
		resp := m.srv.NewListPager(resourceGroupNameParam, workspaceNameParam, nil)
		newListPager = &resp
		m.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armmachinelearning.ManagedNetworkSettingsRuleClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		m.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		m.newListPager.remove(req)
	}
	return resp, nil
}
