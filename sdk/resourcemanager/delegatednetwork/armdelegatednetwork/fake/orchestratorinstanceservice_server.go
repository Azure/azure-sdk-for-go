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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/delegatednetwork/armdelegatednetwork"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

// OrchestratorInstanceServiceServer is a fake server for instances of the armdelegatednetwork.OrchestratorInstanceServiceClient type.
type OrchestratorInstanceServiceServer struct {
	// BeginCreate is the fake for method OrchestratorInstanceServiceClient.BeginCreate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreate func(ctx context.Context, resourceGroupName string, resourceName string, parameters armdelegatednetwork.Orchestrator, options *armdelegatednetwork.OrchestratorInstanceServiceClientBeginCreateOptions) (resp azfake.PollerResponder[armdelegatednetwork.OrchestratorInstanceServiceClientCreateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method OrchestratorInstanceServiceClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, resourceName string, options *armdelegatednetwork.OrchestratorInstanceServiceClientBeginDeleteOptions) (resp azfake.PollerResponder[armdelegatednetwork.OrchestratorInstanceServiceClientDeleteResponse], errResp azfake.ErrorResponder)

	// GetDetails is the fake for method OrchestratorInstanceServiceClient.GetDetails
	// HTTP status codes to indicate success: http.StatusOK
	GetDetails func(ctx context.Context, resourceGroupName string, resourceName string, options *armdelegatednetwork.OrchestratorInstanceServiceClientGetDetailsOptions) (resp azfake.Responder[armdelegatednetwork.OrchestratorInstanceServiceClientGetDetailsResponse], errResp azfake.ErrorResponder)

	// NewListByResourceGroupPager is the fake for method OrchestratorInstanceServiceClient.NewListByResourceGroupPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByResourceGroupPager func(resourceGroupName string, options *armdelegatednetwork.OrchestratorInstanceServiceClientListByResourceGroupOptions) (resp azfake.PagerResponder[armdelegatednetwork.OrchestratorInstanceServiceClientListByResourceGroupResponse])

	// NewListBySubscriptionPager is the fake for method OrchestratorInstanceServiceClient.NewListBySubscriptionPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListBySubscriptionPager func(options *armdelegatednetwork.OrchestratorInstanceServiceClientListBySubscriptionOptions) (resp azfake.PagerResponder[armdelegatednetwork.OrchestratorInstanceServiceClientListBySubscriptionResponse])

	// Patch is the fake for method OrchestratorInstanceServiceClient.Patch
	// HTTP status codes to indicate success: http.StatusOK
	Patch func(ctx context.Context, resourceGroupName string, resourceName string, parameters armdelegatednetwork.OrchestratorResourceUpdateParameters, options *armdelegatednetwork.OrchestratorInstanceServiceClientPatchOptions) (resp azfake.Responder[armdelegatednetwork.OrchestratorInstanceServiceClientPatchResponse], errResp azfake.ErrorResponder)
}

// NewOrchestratorInstanceServiceServerTransport creates a new instance of OrchestratorInstanceServiceServerTransport with the provided implementation.
// The returned OrchestratorInstanceServiceServerTransport instance is connected to an instance of armdelegatednetwork.OrchestratorInstanceServiceClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewOrchestratorInstanceServiceServerTransport(srv *OrchestratorInstanceServiceServer) *OrchestratorInstanceServiceServerTransport {
	return &OrchestratorInstanceServiceServerTransport{
		srv:                         srv,
		beginCreate:                 newTracker[azfake.PollerResponder[armdelegatednetwork.OrchestratorInstanceServiceClientCreateResponse]](),
		beginDelete:                 newTracker[azfake.PollerResponder[armdelegatednetwork.OrchestratorInstanceServiceClientDeleteResponse]](),
		newListByResourceGroupPager: newTracker[azfake.PagerResponder[armdelegatednetwork.OrchestratorInstanceServiceClientListByResourceGroupResponse]](),
		newListBySubscriptionPager:  newTracker[azfake.PagerResponder[armdelegatednetwork.OrchestratorInstanceServiceClientListBySubscriptionResponse]](),
	}
}

// OrchestratorInstanceServiceServerTransport connects instances of armdelegatednetwork.OrchestratorInstanceServiceClient to instances of OrchestratorInstanceServiceServer.
// Don't use this type directly, use NewOrchestratorInstanceServiceServerTransport instead.
type OrchestratorInstanceServiceServerTransport struct {
	srv                         *OrchestratorInstanceServiceServer
	beginCreate                 *tracker[azfake.PollerResponder[armdelegatednetwork.OrchestratorInstanceServiceClientCreateResponse]]
	beginDelete                 *tracker[azfake.PollerResponder[armdelegatednetwork.OrchestratorInstanceServiceClientDeleteResponse]]
	newListByResourceGroupPager *tracker[azfake.PagerResponder[armdelegatednetwork.OrchestratorInstanceServiceClientListByResourceGroupResponse]]
	newListBySubscriptionPager  *tracker[azfake.PagerResponder[armdelegatednetwork.OrchestratorInstanceServiceClientListBySubscriptionResponse]]
}

// Do implements the policy.Transporter interface for OrchestratorInstanceServiceServerTransport.
func (o *OrchestratorInstanceServiceServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "OrchestratorInstanceServiceClient.BeginCreate":
		resp, err = o.dispatchBeginCreate(req)
	case "OrchestratorInstanceServiceClient.BeginDelete":
		resp, err = o.dispatchBeginDelete(req)
	case "OrchestratorInstanceServiceClient.GetDetails":
		resp, err = o.dispatchGetDetails(req)
	case "OrchestratorInstanceServiceClient.NewListByResourceGroupPager":
		resp, err = o.dispatchNewListByResourceGroupPager(req)
	case "OrchestratorInstanceServiceClient.NewListBySubscriptionPager":
		resp, err = o.dispatchNewListBySubscriptionPager(req)
	case "OrchestratorInstanceServiceClient.Patch":
		resp, err = o.dispatchPatch(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (o *OrchestratorInstanceServiceServerTransport) dispatchBeginCreate(req *http.Request) (*http.Response, error) {
	if o.srv.BeginCreate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreate not implemented")}
	}
	beginCreate := o.beginCreate.get(req)
	if beginCreate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DelegatedNetwork/orchestrators/(?P<resourceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armdelegatednetwork.Orchestrator](req)
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
		respr, errRespr := o.srv.BeginCreate(req.Context(), resourceGroupNameParam, resourceNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreate = &respr
		o.beginCreate.add(req, beginCreate)
	}

	resp, err := server.PollerResponderNext(beginCreate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		o.beginCreate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreate) {
		o.beginCreate.remove(req)
	}

	return resp, nil
}

func (o *OrchestratorInstanceServiceServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if o.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := o.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DelegatedNetwork/orchestrators/(?P<resourceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
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
		resourceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceName")])
		if err != nil {
			return nil, err
		}
		forceDeleteUnescaped, err := url.QueryUnescape(qp.Get("forceDelete"))
		if err != nil {
			return nil, err
		}
		forceDeleteParam, err := parseOptional(forceDeleteUnescaped, strconv.ParseBool)
		if err != nil {
			return nil, err
		}
		var options *armdelegatednetwork.OrchestratorInstanceServiceClientBeginDeleteOptions
		if forceDeleteParam != nil {
			options = &armdelegatednetwork.OrchestratorInstanceServiceClientBeginDeleteOptions{
				ForceDelete: forceDeleteParam,
			}
		}
		respr, errRespr := o.srv.BeginDelete(req.Context(), resourceGroupNameParam, resourceNameParam, options)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		o.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		o.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		o.beginDelete.remove(req)
	}

	return resp, nil
}

func (o *OrchestratorInstanceServiceServerTransport) dispatchGetDetails(req *http.Request) (*http.Response, error) {
	if o.srv.GetDetails == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetDetails not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DelegatedNetwork/orchestrators/(?P<resourceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	resourceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := o.srv.GetDetails(req.Context(), resourceGroupNameParam, resourceNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).Orchestrator, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *OrchestratorInstanceServiceServerTransport) dispatchNewListByResourceGroupPager(req *http.Request) (*http.Response, error) {
	if o.srv.NewListByResourceGroupPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByResourceGroupPager not implemented")}
	}
	newListByResourceGroupPager := o.newListByResourceGroupPager.get(req)
	if newListByResourceGroupPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DelegatedNetwork/orchestrators`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		resp := o.srv.NewListByResourceGroupPager(resourceGroupNameParam, nil)
		newListByResourceGroupPager = &resp
		o.newListByResourceGroupPager.add(req, newListByResourceGroupPager)
		server.PagerResponderInjectNextLinks(newListByResourceGroupPager, req, func(page *armdelegatednetwork.OrchestratorInstanceServiceClientListByResourceGroupResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByResourceGroupPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		o.newListByResourceGroupPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByResourceGroupPager) {
		o.newListByResourceGroupPager.remove(req)
	}
	return resp, nil
}

func (o *OrchestratorInstanceServiceServerTransport) dispatchNewListBySubscriptionPager(req *http.Request) (*http.Response, error) {
	if o.srv.NewListBySubscriptionPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListBySubscriptionPager not implemented")}
	}
	newListBySubscriptionPager := o.newListBySubscriptionPager.get(req)
	if newListBySubscriptionPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DelegatedNetwork/orchestrators`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resp := o.srv.NewListBySubscriptionPager(nil)
		newListBySubscriptionPager = &resp
		o.newListBySubscriptionPager.add(req, newListBySubscriptionPager)
		server.PagerResponderInjectNextLinks(newListBySubscriptionPager, req, func(page *armdelegatednetwork.OrchestratorInstanceServiceClientListBySubscriptionResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListBySubscriptionPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		o.newListBySubscriptionPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListBySubscriptionPager) {
		o.newListBySubscriptionPager.remove(req)
	}
	return resp, nil
}

func (o *OrchestratorInstanceServiceServerTransport) dispatchPatch(req *http.Request) (*http.Response, error) {
	if o.srv.Patch == nil {
		return nil, &nonRetriableError{errors.New("fake for method Patch not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DelegatedNetwork/orchestrators/(?P<resourceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armdelegatednetwork.OrchestratorResourceUpdateParameters](req)
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
	respr, errRespr := o.srv.Patch(req.Context(), resourceGroupNameParam, resourceNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).Orchestrator, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
