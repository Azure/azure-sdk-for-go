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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry"
	"net/http"
	"net/url"
	"regexp"
)

// WebhooksServer is a fake server for instances of the armcontainerregistry.WebhooksClient type.
type WebhooksServer struct {
	// BeginCreate is the fake for method WebhooksClient.BeginCreate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreate func(ctx context.Context, resourceGroupName string, registryName string, webhookName string, webhookCreateParameters armcontainerregistry.WebhookCreateParameters, options *armcontainerregistry.WebhooksClientBeginCreateOptions) (resp azfake.PollerResponder[armcontainerregistry.WebhooksClientCreateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method WebhooksClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, registryName string, webhookName string, options *armcontainerregistry.WebhooksClientBeginDeleteOptions) (resp azfake.PollerResponder[armcontainerregistry.WebhooksClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method WebhooksClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, registryName string, webhookName string, options *armcontainerregistry.WebhooksClientGetOptions) (resp azfake.Responder[armcontainerregistry.WebhooksClientGetResponse], errResp azfake.ErrorResponder)

	// GetCallbackConfig is the fake for method WebhooksClient.GetCallbackConfig
	// HTTP status codes to indicate success: http.StatusOK
	GetCallbackConfig func(ctx context.Context, resourceGroupName string, registryName string, webhookName string, options *armcontainerregistry.WebhooksClientGetCallbackConfigOptions) (resp azfake.Responder[armcontainerregistry.WebhooksClientGetCallbackConfigResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method WebhooksClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(resourceGroupName string, registryName string, options *armcontainerregistry.WebhooksClientListOptions) (resp azfake.PagerResponder[armcontainerregistry.WebhooksClientListResponse])

	// NewListEventsPager is the fake for method WebhooksClient.NewListEventsPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListEventsPager func(resourceGroupName string, registryName string, webhookName string, options *armcontainerregistry.WebhooksClientListEventsOptions) (resp azfake.PagerResponder[armcontainerregistry.WebhooksClientListEventsResponse])

	// Ping is the fake for method WebhooksClient.Ping
	// HTTP status codes to indicate success: http.StatusOK
	Ping func(ctx context.Context, resourceGroupName string, registryName string, webhookName string, options *armcontainerregistry.WebhooksClientPingOptions) (resp azfake.Responder[armcontainerregistry.WebhooksClientPingResponse], errResp azfake.ErrorResponder)

	// BeginUpdate is the fake for method WebhooksClient.BeginUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginUpdate func(ctx context.Context, resourceGroupName string, registryName string, webhookName string, webhookUpdateParameters armcontainerregistry.WebhookUpdateParameters, options *armcontainerregistry.WebhooksClientBeginUpdateOptions) (resp azfake.PollerResponder[armcontainerregistry.WebhooksClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewWebhooksServerTransport creates a new instance of WebhooksServerTransport with the provided implementation.
// The returned WebhooksServerTransport instance is connected to an instance of armcontainerregistry.WebhooksClient by way of the
// undefined.Transporter field.
func NewWebhooksServerTransport(srv *WebhooksServer) *WebhooksServerTransport {
	return &WebhooksServerTransport{srv: srv}
}

// WebhooksServerTransport connects instances of armcontainerregistry.WebhooksClient to instances of WebhooksServer.
// Don't use this type directly, use NewWebhooksServerTransport instead.
type WebhooksServerTransport struct {
	srv                *WebhooksServer
	beginCreate        *azfake.PollerResponder[armcontainerregistry.WebhooksClientCreateResponse]
	beginDelete        *azfake.PollerResponder[armcontainerregistry.WebhooksClientDeleteResponse]
	newListPager       *azfake.PagerResponder[armcontainerregistry.WebhooksClientListResponse]
	newListEventsPager *azfake.PagerResponder[armcontainerregistry.WebhooksClientListEventsResponse]
	beginUpdate        *azfake.PollerResponder[armcontainerregistry.WebhooksClientUpdateResponse]
}

// Do implements the policy.Transporter interface for WebhooksServerTransport.
func (w *WebhooksServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "WebhooksClient.BeginCreate":
		resp, err = w.dispatchBeginCreate(req)
	case "WebhooksClient.BeginDelete":
		resp, err = w.dispatchBeginDelete(req)
	case "WebhooksClient.Get":
		resp, err = w.dispatchGet(req)
	case "WebhooksClient.GetCallbackConfig":
		resp, err = w.dispatchGetCallbackConfig(req)
	case "WebhooksClient.NewListPager":
		resp, err = w.dispatchNewListPager(req)
	case "WebhooksClient.NewListEventsPager":
		resp, err = w.dispatchNewListEventsPager(req)
	case "WebhooksClient.Ping":
		resp, err = w.dispatchPing(req)
	case "WebhooksClient.BeginUpdate":
		resp, err = w.dispatchBeginUpdate(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (w *WebhooksServerTransport) dispatchBeginCreate(req *http.Request) (*http.Response, error) {
	if w.srv.BeginCreate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreate not implemented")}
	}
	if w.beginCreate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft.ContainerRegistry/registries/(?P<registryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/webhooks/(?P<webhookName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armcontainerregistry.WebhookCreateParameters](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		registryNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("registryName")])
		if err != nil {
			return nil, err
		}
		webhookNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("webhookName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := w.srv.BeginCreate(req.Context(), resourceGroupNameUnescaped, registryNameUnescaped, webhookNameUnescaped, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		w.beginCreate = &respr
	}

	resp, err := server.PollerResponderNext(w.beginCreate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(w.beginCreate) {
		w.beginCreate = nil
	}

	return resp, nil
}

func (w *WebhooksServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if w.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	if w.beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft.ContainerRegistry/registries/(?P<registryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/webhooks/(?P<webhookName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		registryNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("registryName")])
		if err != nil {
			return nil, err
		}
		webhookNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("webhookName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := w.srv.BeginDelete(req.Context(), resourceGroupNameUnescaped, registryNameUnescaped, webhookNameUnescaped, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		w.beginDelete = &respr
	}

	resp, err := server.PollerResponderNext(w.beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(w.beginDelete) {
		w.beginDelete = nil
	}

	return resp, nil
}

func (w *WebhooksServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if w.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft.ContainerRegistry/registries/(?P<registryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/webhooks/(?P<webhookName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	registryNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("registryName")])
	if err != nil {
		return nil, err
	}
	webhookNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("webhookName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := w.srv.Get(req.Context(), resourceGroupNameUnescaped, registryNameUnescaped, webhookNameUnescaped, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).Webhook, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (w *WebhooksServerTransport) dispatchGetCallbackConfig(req *http.Request) (*http.Response, error) {
	if w.srv.GetCallbackConfig == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetCallbackConfig not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft.ContainerRegistry/registries/(?P<registryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/webhooks/(?P<webhookName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/getCallbackConfig`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	registryNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("registryName")])
	if err != nil {
		return nil, err
	}
	webhookNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("webhookName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := w.srv.GetCallbackConfig(req.Context(), resourceGroupNameUnescaped, registryNameUnescaped, webhookNameUnescaped, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).CallbackConfig, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (w *WebhooksServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if w.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	if w.newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft.ContainerRegistry/registries/(?P<registryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/webhooks`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		registryNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("registryName")])
		if err != nil {
			return nil, err
		}
		resp := w.srv.NewListPager(resourceGroupNameUnescaped, registryNameUnescaped, nil)
		w.newListPager = &resp
		server.PagerResponderInjectNextLinks(w.newListPager, req, func(page *armcontainerregistry.WebhooksClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(w.newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(w.newListPager) {
		w.newListPager = nil
	}
	return resp, nil
}

func (w *WebhooksServerTransport) dispatchNewListEventsPager(req *http.Request) (*http.Response, error) {
	if w.srv.NewListEventsPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListEventsPager not implemented")}
	}
	if w.newListEventsPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft.ContainerRegistry/registries/(?P<registryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/webhooks/(?P<webhookName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/listEvents`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		registryNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("registryName")])
		if err != nil {
			return nil, err
		}
		webhookNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("webhookName")])
		if err != nil {
			return nil, err
		}
		resp := w.srv.NewListEventsPager(resourceGroupNameUnescaped, registryNameUnescaped, webhookNameUnescaped, nil)
		w.newListEventsPager = &resp
		server.PagerResponderInjectNextLinks(w.newListEventsPager, req, func(page *armcontainerregistry.WebhooksClientListEventsResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(w.newListEventsPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(w.newListEventsPager) {
		w.newListEventsPager = nil
	}
	return resp, nil
}

func (w *WebhooksServerTransport) dispatchPing(req *http.Request) (*http.Response, error) {
	if w.srv.Ping == nil {
		return nil, &nonRetriableError{errors.New("fake for method Ping not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft.ContainerRegistry/registries/(?P<registryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/webhooks/(?P<webhookName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/ping`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	registryNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("registryName")])
	if err != nil {
		return nil, err
	}
	webhookNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("webhookName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := w.srv.Ping(req.Context(), resourceGroupNameUnescaped, registryNameUnescaped, webhookNameUnescaped, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).EventInfo, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (w *WebhooksServerTransport) dispatchBeginUpdate(req *http.Request) (*http.Response, error) {
	if w.srv.BeginUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdate not implemented")}
	}
	if w.beginUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft.ContainerRegistry/registries/(?P<registryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/webhooks/(?P<webhookName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armcontainerregistry.WebhookUpdateParameters](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		registryNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("registryName")])
		if err != nil {
			return nil, err
		}
		webhookNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("webhookName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := w.srv.BeginUpdate(req.Context(), resourceGroupNameUnescaped, registryNameUnescaped, webhookNameUnescaped, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		w.beginUpdate = &respr
	}

	resp, err := server.PollerResponderNext(w.beginUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(w.beginUpdate) {
		w.beginUpdate = nil
	}

	return resp, nil
}
