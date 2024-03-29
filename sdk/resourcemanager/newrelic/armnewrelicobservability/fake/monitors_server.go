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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/newrelic/armnewrelicobservability"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

// MonitorsServer is a fake server for instances of the armnewrelicobservability.MonitorsClient type.
type MonitorsServer struct {
	// BeginCreateOrUpdate is the fake for method MonitorsClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, monitorName string, resource armnewrelicobservability.NewRelicMonitorResource, options *armnewrelicobservability.MonitorsClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armnewrelicobservability.MonitorsClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method MonitorsClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, userEmail string, monitorName string, options *armnewrelicobservability.MonitorsClientBeginDeleteOptions) (resp azfake.PollerResponder[armnewrelicobservability.MonitorsClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method MonitorsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, monitorName string, options *armnewrelicobservability.MonitorsClientGetOptions) (resp azfake.Responder[armnewrelicobservability.MonitorsClientGetResponse], errResp azfake.ErrorResponder)

	// GetMetricRules is the fake for method MonitorsClient.GetMetricRules
	// HTTP status codes to indicate success: http.StatusOK
	GetMetricRules func(ctx context.Context, resourceGroupName string, monitorName string, request armnewrelicobservability.MetricsRequest, options *armnewrelicobservability.MonitorsClientGetMetricRulesOptions) (resp azfake.Responder[armnewrelicobservability.MonitorsClientGetMetricRulesResponse], errResp azfake.ErrorResponder)

	// GetMetricStatus is the fake for method MonitorsClient.GetMetricStatus
	// HTTP status codes to indicate success: http.StatusOK
	GetMetricStatus func(ctx context.Context, resourceGroupName string, monitorName string, request armnewrelicobservability.MetricsStatusRequest, options *armnewrelicobservability.MonitorsClientGetMetricStatusOptions) (resp azfake.Responder[armnewrelicobservability.MonitorsClientGetMetricStatusResponse], errResp azfake.ErrorResponder)

	// NewListAppServicesPager is the fake for method MonitorsClient.NewListAppServicesPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListAppServicesPager func(resourceGroupName string, monitorName string, request armnewrelicobservability.AppServicesGetRequest, options *armnewrelicobservability.MonitorsClientListAppServicesOptions) (resp azfake.PagerResponder[armnewrelicobservability.MonitorsClientListAppServicesResponse])

	// NewListByResourceGroupPager is the fake for method MonitorsClient.NewListByResourceGroupPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByResourceGroupPager func(resourceGroupName string, options *armnewrelicobservability.MonitorsClientListByResourceGroupOptions) (resp azfake.PagerResponder[armnewrelicobservability.MonitorsClientListByResourceGroupResponse])

	// NewListBySubscriptionPager is the fake for method MonitorsClient.NewListBySubscriptionPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListBySubscriptionPager func(options *armnewrelicobservability.MonitorsClientListBySubscriptionOptions) (resp azfake.PagerResponder[armnewrelicobservability.MonitorsClientListBySubscriptionResponse])

	// NewListHostsPager is the fake for method MonitorsClient.NewListHostsPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListHostsPager func(resourceGroupName string, monitorName string, request armnewrelicobservability.HostsGetRequest, options *armnewrelicobservability.MonitorsClientListHostsOptions) (resp azfake.PagerResponder[armnewrelicobservability.MonitorsClientListHostsResponse])

	// NewListLinkedResourcesPager is the fake for method MonitorsClient.NewListLinkedResourcesPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListLinkedResourcesPager func(resourceGroupName string, monitorName string, options *armnewrelicobservability.MonitorsClientListLinkedResourcesOptions) (resp azfake.PagerResponder[armnewrelicobservability.MonitorsClientListLinkedResourcesResponse])

	// NewListMonitoredResourcesPager is the fake for method MonitorsClient.NewListMonitoredResourcesPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListMonitoredResourcesPager func(resourceGroupName string, monitorName string, options *armnewrelicobservability.MonitorsClientListMonitoredResourcesOptions) (resp azfake.PagerResponder[armnewrelicobservability.MonitorsClientListMonitoredResourcesResponse])

	// SwitchBilling is the fake for method MonitorsClient.SwitchBilling
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	SwitchBilling func(ctx context.Context, resourceGroupName string, monitorName string, request armnewrelicobservability.SwitchBillingRequest, options *armnewrelicobservability.MonitorsClientSwitchBillingOptions) (resp azfake.Responder[armnewrelicobservability.MonitorsClientSwitchBillingResponse], errResp azfake.ErrorResponder)

	// Update is the fake for method MonitorsClient.Update
	// HTTP status codes to indicate success: http.StatusOK
	Update func(ctx context.Context, resourceGroupName string, monitorName string, properties armnewrelicobservability.NewRelicMonitorResourceUpdate, options *armnewrelicobservability.MonitorsClientUpdateOptions) (resp azfake.Responder[armnewrelicobservability.MonitorsClientUpdateResponse], errResp azfake.ErrorResponder)

	// VMHostPayload is the fake for method MonitorsClient.VMHostPayload
	// HTTP status codes to indicate success: http.StatusOK
	VMHostPayload func(ctx context.Context, resourceGroupName string, monitorName string, options *armnewrelicobservability.MonitorsClientVMHostPayloadOptions) (resp azfake.Responder[armnewrelicobservability.MonitorsClientVMHostPayloadResponse], errResp azfake.ErrorResponder)
}

// NewMonitorsServerTransport creates a new instance of MonitorsServerTransport with the provided implementation.
// The returned MonitorsServerTransport instance is connected to an instance of armnewrelicobservability.MonitorsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewMonitorsServerTransport(srv *MonitorsServer) *MonitorsServerTransport {
	return &MonitorsServerTransport{
		srv:                            srv,
		beginCreateOrUpdate:            newTracker[azfake.PollerResponder[armnewrelicobservability.MonitorsClientCreateOrUpdateResponse]](),
		beginDelete:                    newTracker[azfake.PollerResponder[armnewrelicobservability.MonitorsClientDeleteResponse]](),
		newListAppServicesPager:        newTracker[azfake.PagerResponder[armnewrelicobservability.MonitorsClientListAppServicesResponse]](),
		newListByResourceGroupPager:    newTracker[azfake.PagerResponder[armnewrelicobservability.MonitorsClientListByResourceGroupResponse]](),
		newListBySubscriptionPager:     newTracker[azfake.PagerResponder[armnewrelicobservability.MonitorsClientListBySubscriptionResponse]](),
		newListHostsPager:              newTracker[azfake.PagerResponder[armnewrelicobservability.MonitorsClientListHostsResponse]](),
		newListLinkedResourcesPager:    newTracker[azfake.PagerResponder[armnewrelicobservability.MonitorsClientListLinkedResourcesResponse]](),
		newListMonitoredResourcesPager: newTracker[azfake.PagerResponder[armnewrelicobservability.MonitorsClientListMonitoredResourcesResponse]](),
	}
}

// MonitorsServerTransport connects instances of armnewrelicobservability.MonitorsClient to instances of MonitorsServer.
// Don't use this type directly, use NewMonitorsServerTransport instead.
type MonitorsServerTransport struct {
	srv                            *MonitorsServer
	beginCreateOrUpdate            *tracker[azfake.PollerResponder[armnewrelicobservability.MonitorsClientCreateOrUpdateResponse]]
	beginDelete                    *tracker[azfake.PollerResponder[armnewrelicobservability.MonitorsClientDeleteResponse]]
	newListAppServicesPager        *tracker[azfake.PagerResponder[armnewrelicobservability.MonitorsClientListAppServicesResponse]]
	newListByResourceGroupPager    *tracker[azfake.PagerResponder[armnewrelicobservability.MonitorsClientListByResourceGroupResponse]]
	newListBySubscriptionPager     *tracker[azfake.PagerResponder[armnewrelicobservability.MonitorsClientListBySubscriptionResponse]]
	newListHostsPager              *tracker[azfake.PagerResponder[armnewrelicobservability.MonitorsClientListHostsResponse]]
	newListLinkedResourcesPager    *tracker[azfake.PagerResponder[armnewrelicobservability.MonitorsClientListLinkedResourcesResponse]]
	newListMonitoredResourcesPager *tracker[azfake.PagerResponder[armnewrelicobservability.MonitorsClientListMonitoredResourcesResponse]]
}

// Do implements the policy.Transporter interface for MonitorsServerTransport.
func (m *MonitorsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "MonitorsClient.BeginCreateOrUpdate":
		resp, err = m.dispatchBeginCreateOrUpdate(req)
	case "MonitorsClient.BeginDelete":
		resp, err = m.dispatchBeginDelete(req)
	case "MonitorsClient.Get":
		resp, err = m.dispatchGet(req)
	case "MonitorsClient.GetMetricRules":
		resp, err = m.dispatchGetMetricRules(req)
	case "MonitorsClient.GetMetricStatus":
		resp, err = m.dispatchGetMetricStatus(req)
	case "MonitorsClient.NewListAppServicesPager":
		resp, err = m.dispatchNewListAppServicesPager(req)
	case "MonitorsClient.NewListByResourceGroupPager":
		resp, err = m.dispatchNewListByResourceGroupPager(req)
	case "MonitorsClient.NewListBySubscriptionPager":
		resp, err = m.dispatchNewListBySubscriptionPager(req)
	case "MonitorsClient.NewListHostsPager":
		resp, err = m.dispatchNewListHostsPager(req)
	case "MonitorsClient.NewListLinkedResourcesPager":
		resp, err = m.dispatchNewListLinkedResourcesPager(req)
	case "MonitorsClient.NewListMonitoredResourcesPager":
		resp, err = m.dispatchNewListMonitoredResourcesPager(req)
	case "MonitorsClient.SwitchBilling":
		resp, err = m.dispatchSwitchBilling(req)
	case "MonitorsClient.Update":
		resp, err = m.dispatchUpdate(req)
	case "MonitorsClient.VMHostPayload":
		resp, err = m.dispatchVMHostPayload(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *MonitorsServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if m.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := m.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/NewRelic\.Observability/monitors/(?P<monitorName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armnewrelicobservability.NewRelicMonitorResource](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		monitorNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("monitorName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := m.srv.BeginCreateOrUpdate(req.Context(), resourceGroupNameParam, monitorNameParam, body, nil)
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

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		m.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		m.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (m *MonitorsServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if m.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := m.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/NewRelic\.Observability/monitors/(?P<monitorName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		userEmailParam, err := url.QueryUnescape(qp.Get("userEmail"))
		if err != nil {
			return nil, err
		}
		monitorNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("monitorName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := m.srv.BeginDelete(req.Context(), resourceGroupNameParam, userEmailParam, monitorNameParam, nil)
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

func (m *MonitorsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if m.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/NewRelic\.Observability/monitors/(?P<monitorName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	monitorNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("monitorName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := m.srv.Get(req.Context(), resourceGroupNameParam, monitorNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).NewRelicMonitorResource, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *MonitorsServerTransport) dispatchGetMetricRules(req *http.Request) (*http.Response, error) {
	if m.srv.GetMetricRules == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetMetricRules not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/NewRelic\.Observability/monitors/(?P<monitorName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/getMetricRules`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armnewrelicobservability.MetricsRequest](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	monitorNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("monitorName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := m.srv.GetMetricRules(req.Context(), resourceGroupNameParam, monitorNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).MetricRules, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *MonitorsServerTransport) dispatchGetMetricStatus(req *http.Request) (*http.Response, error) {
	if m.srv.GetMetricStatus == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetMetricStatus not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/NewRelic\.Observability/monitors/(?P<monitorName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/getMetricStatus`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armnewrelicobservability.MetricsStatusRequest](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	monitorNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("monitorName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := m.srv.GetMetricStatus(req.Context(), resourceGroupNameParam, monitorNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).MetricsStatusResponse, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *MonitorsServerTransport) dispatchNewListAppServicesPager(req *http.Request) (*http.Response, error) {
	if m.srv.NewListAppServicesPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListAppServicesPager not implemented")}
	}
	newListAppServicesPager := m.newListAppServicesPager.get(req)
	if newListAppServicesPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/NewRelic\.Observability/monitors/(?P<monitorName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/listAppServices`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armnewrelicobservability.AppServicesGetRequest](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		monitorNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("monitorName")])
		if err != nil {
			return nil, err
		}
		resp := m.srv.NewListAppServicesPager(resourceGroupNameParam, monitorNameParam, body, nil)
		newListAppServicesPager = &resp
		m.newListAppServicesPager.add(req, newListAppServicesPager)
		server.PagerResponderInjectNextLinks(newListAppServicesPager, req, func(page *armnewrelicobservability.MonitorsClientListAppServicesResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListAppServicesPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		m.newListAppServicesPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListAppServicesPager) {
		m.newListAppServicesPager.remove(req)
	}
	return resp, nil
}

func (m *MonitorsServerTransport) dispatchNewListByResourceGroupPager(req *http.Request) (*http.Response, error) {
	if m.srv.NewListByResourceGroupPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByResourceGroupPager not implemented")}
	}
	newListByResourceGroupPager := m.newListByResourceGroupPager.get(req)
	if newListByResourceGroupPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/NewRelic\.Observability/monitors`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		resp := m.srv.NewListByResourceGroupPager(resourceGroupNameParam, nil)
		newListByResourceGroupPager = &resp
		m.newListByResourceGroupPager.add(req, newListByResourceGroupPager)
		server.PagerResponderInjectNextLinks(newListByResourceGroupPager, req, func(page *armnewrelicobservability.MonitorsClientListByResourceGroupResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByResourceGroupPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		m.newListByResourceGroupPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByResourceGroupPager) {
		m.newListByResourceGroupPager.remove(req)
	}
	return resp, nil
}

func (m *MonitorsServerTransport) dispatchNewListBySubscriptionPager(req *http.Request) (*http.Response, error) {
	if m.srv.NewListBySubscriptionPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListBySubscriptionPager not implemented")}
	}
	newListBySubscriptionPager := m.newListBySubscriptionPager.get(req)
	if newListBySubscriptionPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/NewRelic\.Observability/monitors`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resp := m.srv.NewListBySubscriptionPager(nil)
		newListBySubscriptionPager = &resp
		m.newListBySubscriptionPager.add(req, newListBySubscriptionPager)
		server.PagerResponderInjectNextLinks(newListBySubscriptionPager, req, func(page *armnewrelicobservability.MonitorsClientListBySubscriptionResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListBySubscriptionPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		m.newListBySubscriptionPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListBySubscriptionPager) {
		m.newListBySubscriptionPager.remove(req)
	}
	return resp, nil
}

func (m *MonitorsServerTransport) dispatchNewListHostsPager(req *http.Request) (*http.Response, error) {
	if m.srv.NewListHostsPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListHostsPager not implemented")}
	}
	newListHostsPager := m.newListHostsPager.get(req)
	if newListHostsPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/NewRelic\.Observability/monitors/(?P<monitorName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/listHosts`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armnewrelicobservability.HostsGetRequest](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		monitorNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("monitorName")])
		if err != nil {
			return nil, err
		}
		resp := m.srv.NewListHostsPager(resourceGroupNameParam, monitorNameParam, body, nil)
		newListHostsPager = &resp
		m.newListHostsPager.add(req, newListHostsPager)
		server.PagerResponderInjectNextLinks(newListHostsPager, req, func(page *armnewrelicobservability.MonitorsClientListHostsResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListHostsPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		m.newListHostsPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListHostsPager) {
		m.newListHostsPager.remove(req)
	}
	return resp, nil
}

func (m *MonitorsServerTransport) dispatchNewListLinkedResourcesPager(req *http.Request) (*http.Response, error) {
	if m.srv.NewListLinkedResourcesPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListLinkedResourcesPager not implemented")}
	}
	newListLinkedResourcesPager := m.newListLinkedResourcesPager.get(req)
	if newListLinkedResourcesPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/NewRelic\.Observability/monitors/(?P<monitorName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/listLinkedResources`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		monitorNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("monitorName")])
		if err != nil {
			return nil, err
		}
		resp := m.srv.NewListLinkedResourcesPager(resourceGroupNameParam, monitorNameParam, nil)
		newListLinkedResourcesPager = &resp
		m.newListLinkedResourcesPager.add(req, newListLinkedResourcesPager)
		server.PagerResponderInjectNextLinks(newListLinkedResourcesPager, req, func(page *armnewrelicobservability.MonitorsClientListLinkedResourcesResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListLinkedResourcesPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		m.newListLinkedResourcesPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListLinkedResourcesPager) {
		m.newListLinkedResourcesPager.remove(req)
	}
	return resp, nil
}

func (m *MonitorsServerTransport) dispatchNewListMonitoredResourcesPager(req *http.Request) (*http.Response, error) {
	if m.srv.NewListMonitoredResourcesPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListMonitoredResourcesPager not implemented")}
	}
	newListMonitoredResourcesPager := m.newListMonitoredResourcesPager.get(req)
	if newListMonitoredResourcesPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/NewRelic\.Observability/monitors/(?P<monitorName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/monitoredResources`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		monitorNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("monitorName")])
		if err != nil {
			return nil, err
		}
		resp := m.srv.NewListMonitoredResourcesPager(resourceGroupNameParam, monitorNameParam, nil)
		newListMonitoredResourcesPager = &resp
		m.newListMonitoredResourcesPager.add(req, newListMonitoredResourcesPager)
		server.PagerResponderInjectNextLinks(newListMonitoredResourcesPager, req, func(page *armnewrelicobservability.MonitorsClientListMonitoredResourcesResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListMonitoredResourcesPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		m.newListMonitoredResourcesPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListMonitoredResourcesPager) {
		m.newListMonitoredResourcesPager.remove(req)
	}
	return resp, nil
}

func (m *MonitorsServerTransport) dispatchSwitchBilling(req *http.Request) (*http.Response, error) {
	if m.srv.SwitchBilling == nil {
		return nil, &nonRetriableError{errors.New("fake for method SwitchBilling not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/NewRelic\.Observability/monitors/(?P<monitorName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/switchBilling`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armnewrelicobservability.SwitchBillingRequest](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	monitorNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("monitorName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := m.srv.SwitchBilling(req.Context(), resourceGroupNameParam, monitorNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).NewRelicMonitorResource, req)
	if err != nil {
		return nil, err
	}
	if val := server.GetResponse(respr).RetryAfter; val != nil {
		resp.Header.Set("Retry-After", strconv.FormatInt(int64(*val), 10))
	}
	return resp, nil
}

func (m *MonitorsServerTransport) dispatchUpdate(req *http.Request) (*http.Response, error) {
	if m.srv.Update == nil {
		return nil, &nonRetriableError{errors.New("fake for method Update not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/NewRelic\.Observability/monitors/(?P<monitorName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armnewrelicobservability.NewRelicMonitorResourceUpdate](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	monitorNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("monitorName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := m.srv.Update(req.Context(), resourceGroupNameParam, monitorNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).NewRelicMonitorResource, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *MonitorsServerTransport) dispatchVMHostPayload(req *http.Request) (*http.Response, error) {
	if m.srv.VMHostPayload == nil {
		return nil, &nonRetriableError{errors.New("fake for method VMHostPayload not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/NewRelic\.Observability/monitors/(?P<monitorName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/vmHostPayloads`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	monitorNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("monitorName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := m.srv.VMHostPayload(req.Context(), resourceGroupNameParam, monitorNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).VMExtensionPayload, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
