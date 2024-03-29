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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/security/armsecurity"
	"net/http"
	"net/url"
	"regexp"
)

// AzureDevOpsOrgsServer is a fake server for instances of the armsecurity.AzureDevOpsOrgsClient type.
type AzureDevOpsOrgsServer struct {
	// BeginCreateOrUpdate is the fake for method AzureDevOpsOrgsClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, securityConnectorName string, orgName string, azureDevOpsOrg armsecurity.AzureDevOpsOrg, options *armsecurity.AzureDevOpsOrgsClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armsecurity.AzureDevOpsOrgsClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method AzureDevOpsOrgsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, securityConnectorName string, orgName string, options *armsecurity.AzureDevOpsOrgsClientGetOptions) (resp azfake.Responder[armsecurity.AzureDevOpsOrgsClientGetResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method AzureDevOpsOrgsClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(resourceGroupName string, securityConnectorName string, options *armsecurity.AzureDevOpsOrgsClientListOptions) (resp azfake.PagerResponder[armsecurity.AzureDevOpsOrgsClientListResponse])

	// ListAvailable is the fake for method AzureDevOpsOrgsClient.ListAvailable
	// HTTP status codes to indicate success: http.StatusOK
	ListAvailable func(ctx context.Context, resourceGroupName string, securityConnectorName string, options *armsecurity.AzureDevOpsOrgsClientListAvailableOptions) (resp azfake.Responder[armsecurity.AzureDevOpsOrgsClientListAvailableResponse], errResp azfake.ErrorResponder)

	// BeginUpdate is the fake for method AzureDevOpsOrgsClient.BeginUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginUpdate func(ctx context.Context, resourceGroupName string, securityConnectorName string, orgName string, azureDevOpsOrg armsecurity.AzureDevOpsOrg, options *armsecurity.AzureDevOpsOrgsClientBeginUpdateOptions) (resp azfake.PollerResponder[armsecurity.AzureDevOpsOrgsClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewAzureDevOpsOrgsServerTransport creates a new instance of AzureDevOpsOrgsServerTransport with the provided implementation.
// The returned AzureDevOpsOrgsServerTransport instance is connected to an instance of armsecurity.AzureDevOpsOrgsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewAzureDevOpsOrgsServerTransport(srv *AzureDevOpsOrgsServer) *AzureDevOpsOrgsServerTransport {
	return &AzureDevOpsOrgsServerTransport{
		srv:                 srv,
		beginCreateOrUpdate: newTracker[azfake.PollerResponder[armsecurity.AzureDevOpsOrgsClientCreateOrUpdateResponse]](),
		newListPager:        newTracker[azfake.PagerResponder[armsecurity.AzureDevOpsOrgsClientListResponse]](),
		beginUpdate:         newTracker[azfake.PollerResponder[armsecurity.AzureDevOpsOrgsClientUpdateResponse]](),
	}
}

// AzureDevOpsOrgsServerTransport connects instances of armsecurity.AzureDevOpsOrgsClient to instances of AzureDevOpsOrgsServer.
// Don't use this type directly, use NewAzureDevOpsOrgsServerTransport instead.
type AzureDevOpsOrgsServerTransport struct {
	srv                 *AzureDevOpsOrgsServer
	beginCreateOrUpdate *tracker[azfake.PollerResponder[armsecurity.AzureDevOpsOrgsClientCreateOrUpdateResponse]]
	newListPager        *tracker[azfake.PagerResponder[armsecurity.AzureDevOpsOrgsClientListResponse]]
	beginUpdate         *tracker[azfake.PollerResponder[armsecurity.AzureDevOpsOrgsClientUpdateResponse]]
}

// Do implements the policy.Transporter interface for AzureDevOpsOrgsServerTransport.
func (a *AzureDevOpsOrgsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "AzureDevOpsOrgsClient.BeginCreateOrUpdate":
		resp, err = a.dispatchBeginCreateOrUpdate(req)
	case "AzureDevOpsOrgsClient.Get":
		resp, err = a.dispatchGet(req)
	case "AzureDevOpsOrgsClient.NewListPager":
		resp, err = a.dispatchNewListPager(req)
	case "AzureDevOpsOrgsClient.ListAvailable":
		resp, err = a.dispatchListAvailable(req)
	case "AzureDevOpsOrgsClient.BeginUpdate":
		resp, err = a.dispatchBeginUpdate(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *AzureDevOpsOrgsServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if a.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := a.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Security/securityConnectors/(?P<securityConnectorName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/devops/default/azureDevOpsOrgs/(?P<orgName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armsecurity.AzureDevOpsOrg](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		securityConnectorNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("securityConnectorName")])
		if err != nil {
			return nil, err
		}
		orgNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("orgName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := a.srv.BeginCreateOrUpdate(req.Context(), resourceGroupNameParam, securityConnectorNameParam, orgNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateOrUpdate = &respr
		a.beginCreateOrUpdate.add(req, beginCreateOrUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		a.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		a.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (a *AzureDevOpsOrgsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if a.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Security/securityConnectors/(?P<securityConnectorName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/devops/default/azureDevOpsOrgs/(?P<orgName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	securityConnectorNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("securityConnectorName")])
	if err != nil {
		return nil, err
	}
	orgNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("orgName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := a.srv.Get(req.Context(), resourceGroupNameParam, securityConnectorNameParam, orgNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).AzureDevOpsOrg, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *AzureDevOpsOrgsServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if a.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := a.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Security/securityConnectors/(?P<securityConnectorName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/devops/default/azureDevOpsOrgs`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		securityConnectorNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("securityConnectorName")])
		if err != nil {
			return nil, err
		}
		resp := a.srv.NewListPager(resourceGroupNameParam, securityConnectorNameParam, nil)
		newListPager = &resp
		a.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armsecurity.AzureDevOpsOrgsClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		a.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		a.newListPager.remove(req)
	}
	return resp, nil
}

func (a *AzureDevOpsOrgsServerTransport) dispatchListAvailable(req *http.Request) (*http.Response, error) {
	if a.srv.ListAvailable == nil {
		return nil, &nonRetriableError{errors.New("fake for method ListAvailable not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Security/securityConnectors/(?P<securityConnectorName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/devops/default/listAvailableAzureDevOpsOrgs`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	securityConnectorNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("securityConnectorName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := a.srv.ListAvailable(req.Context(), resourceGroupNameParam, securityConnectorNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).AzureDevOpsOrgListResponse, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *AzureDevOpsOrgsServerTransport) dispatchBeginUpdate(req *http.Request) (*http.Response, error) {
	if a.srv.BeginUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdate not implemented")}
	}
	beginUpdate := a.beginUpdate.get(req)
	if beginUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Security/securityConnectors/(?P<securityConnectorName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/devops/default/azureDevOpsOrgs/(?P<orgName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armsecurity.AzureDevOpsOrg](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		securityConnectorNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("securityConnectorName")])
		if err != nil {
			return nil, err
		}
		orgNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("orgName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := a.srv.BeginUpdate(req.Context(), resourceGroupNameParam, securityConnectorNameParam, orgNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginUpdate = &respr
		a.beginUpdate.add(req, beginUpdate)
	}

	resp, err := server.PollerResponderNext(beginUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		a.beginUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdate) {
		a.beginUpdate.remove(req)
	}

	return resp, nil
}
