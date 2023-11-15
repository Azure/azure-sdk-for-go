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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mobilenetwork/armmobilenetwork/v3"
	"net/http"
	"net/url"
	"regexp"
)

// SimPoliciesServer is a fake server for instances of the armmobilenetwork.SimPoliciesClient type.
type SimPoliciesServer struct {
	// BeginCreateOrUpdate is the fake for method SimPoliciesClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, mobileNetworkName string, simPolicyName string, parameters armmobilenetwork.SimPolicy, options *armmobilenetwork.SimPoliciesClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armmobilenetwork.SimPoliciesClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method SimPoliciesClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, mobileNetworkName string, simPolicyName string, options *armmobilenetwork.SimPoliciesClientBeginDeleteOptions) (resp azfake.PollerResponder[armmobilenetwork.SimPoliciesClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method SimPoliciesClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, mobileNetworkName string, simPolicyName string, options *armmobilenetwork.SimPoliciesClientGetOptions) (resp azfake.Responder[armmobilenetwork.SimPoliciesClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByMobileNetworkPager is the fake for method SimPoliciesClient.NewListByMobileNetworkPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByMobileNetworkPager func(resourceGroupName string, mobileNetworkName string, options *armmobilenetwork.SimPoliciesClientListByMobileNetworkOptions) (resp azfake.PagerResponder[armmobilenetwork.SimPoliciesClientListByMobileNetworkResponse])

	// UpdateTags is the fake for method SimPoliciesClient.UpdateTags
	// HTTP status codes to indicate success: http.StatusOK
	UpdateTags func(ctx context.Context, resourceGroupName string, mobileNetworkName string, simPolicyName string, parameters armmobilenetwork.TagsObject, options *armmobilenetwork.SimPoliciesClientUpdateTagsOptions) (resp azfake.Responder[armmobilenetwork.SimPoliciesClientUpdateTagsResponse], errResp azfake.ErrorResponder)
}

// NewSimPoliciesServerTransport creates a new instance of SimPoliciesServerTransport with the provided implementation.
// The returned SimPoliciesServerTransport instance is connected to an instance of armmobilenetwork.SimPoliciesClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewSimPoliciesServerTransport(srv *SimPoliciesServer) *SimPoliciesServerTransport {
	return &SimPoliciesServerTransport{
		srv:                         srv,
		beginCreateOrUpdate:         newTracker[azfake.PollerResponder[armmobilenetwork.SimPoliciesClientCreateOrUpdateResponse]](),
		beginDelete:                 newTracker[azfake.PollerResponder[armmobilenetwork.SimPoliciesClientDeleteResponse]](),
		newListByMobileNetworkPager: newTracker[azfake.PagerResponder[armmobilenetwork.SimPoliciesClientListByMobileNetworkResponse]](),
	}
}

// SimPoliciesServerTransport connects instances of armmobilenetwork.SimPoliciesClient to instances of SimPoliciesServer.
// Don't use this type directly, use NewSimPoliciesServerTransport instead.
type SimPoliciesServerTransport struct {
	srv                         *SimPoliciesServer
	beginCreateOrUpdate         *tracker[azfake.PollerResponder[armmobilenetwork.SimPoliciesClientCreateOrUpdateResponse]]
	beginDelete                 *tracker[azfake.PollerResponder[armmobilenetwork.SimPoliciesClientDeleteResponse]]
	newListByMobileNetworkPager *tracker[azfake.PagerResponder[armmobilenetwork.SimPoliciesClientListByMobileNetworkResponse]]
}

// Do implements the policy.Transporter interface for SimPoliciesServerTransport.
func (s *SimPoliciesServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "SimPoliciesClient.BeginCreateOrUpdate":
		resp, err = s.dispatchBeginCreateOrUpdate(req)
	case "SimPoliciesClient.BeginDelete":
		resp, err = s.dispatchBeginDelete(req)
	case "SimPoliciesClient.Get":
		resp, err = s.dispatchGet(req)
	case "SimPoliciesClient.NewListByMobileNetworkPager":
		resp, err = s.dispatchNewListByMobileNetworkPager(req)
	case "SimPoliciesClient.UpdateTags":
		resp, err = s.dispatchUpdateTags(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *SimPoliciesServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if s.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := s.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.MobileNetwork/mobileNetworks/(?P<mobileNetworkName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/simPolicies/(?P<simPolicyName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armmobilenetwork.SimPolicy](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		mobileNetworkNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("mobileNetworkName")])
		if err != nil {
			return nil, err
		}
		simPolicyNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("simPolicyName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := s.srv.BeginCreateOrUpdate(req.Context(), resourceGroupNameParam, mobileNetworkNameParam, simPolicyNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateOrUpdate = &respr
		s.beginCreateOrUpdate.add(req, beginCreateOrUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		s.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		s.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (s *SimPoliciesServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if s.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := s.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.MobileNetwork/mobileNetworks/(?P<mobileNetworkName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/simPolicies/(?P<simPolicyName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		mobileNetworkNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("mobileNetworkName")])
		if err != nil {
			return nil, err
		}
		simPolicyNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("simPolicyName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := s.srv.BeginDelete(req.Context(), resourceGroupNameParam, mobileNetworkNameParam, simPolicyNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		s.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		s.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		s.beginDelete.remove(req)
	}

	return resp, nil
}

func (s *SimPoliciesServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if s.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.MobileNetwork/mobileNetworks/(?P<mobileNetworkName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/simPolicies/(?P<simPolicyName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	mobileNetworkNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("mobileNetworkName")])
	if err != nil {
		return nil, err
	}
	simPolicyNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("simPolicyName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.Get(req.Context(), resourceGroupNameParam, mobileNetworkNameParam, simPolicyNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).SimPolicy, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *SimPoliciesServerTransport) dispatchNewListByMobileNetworkPager(req *http.Request) (*http.Response, error) {
	if s.srv.NewListByMobileNetworkPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByMobileNetworkPager not implemented")}
	}
	newListByMobileNetworkPager := s.newListByMobileNetworkPager.get(req)
	if newListByMobileNetworkPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.MobileNetwork/mobileNetworks/(?P<mobileNetworkName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/simPolicies`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		mobileNetworkNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("mobileNetworkName")])
		if err != nil {
			return nil, err
		}
		resp := s.srv.NewListByMobileNetworkPager(resourceGroupNameParam, mobileNetworkNameParam, nil)
		newListByMobileNetworkPager = &resp
		s.newListByMobileNetworkPager.add(req, newListByMobileNetworkPager)
		server.PagerResponderInjectNextLinks(newListByMobileNetworkPager, req, func(page *armmobilenetwork.SimPoliciesClientListByMobileNetworkResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByMobileNetworkPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		s.newListByMobileNetworkPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByMobileNetworkPager) {
		s.newListByMobileNetworkPager.remove(req)
	}
	return resp, nil
}

func (s *SimPoliciesServerTransport) dispatchUpdateTags(req *http.Request) (*http.Response, error) {
	if s.srv.UpdateTags == nil {
		return nil, &nonRetriableError{errors.New("fake for method UpdateTags not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.MobileNetwork/mobileNetworks/(?P<mobileNetworkName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/simPolicies/(?P<simPolicyName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armmobilenetwork.TagsObject](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	mobileNetworkNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("mobileNetworkName")])
	if err != nil {
		return nil, err
	}
	simPolicyNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("simPolicyName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.UpdateTags(req.Context(), resourceGroupNameParam, mobileNetworkNameParam, simPolicyNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).SimPolicy, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
