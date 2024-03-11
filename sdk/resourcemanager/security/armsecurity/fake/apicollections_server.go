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

// APICollectionsServer is a fake server for instances of the armsecurity.APICollectionsClient type.
type APICollectionsServer struct {
	// GetByAzureAPIManagementService is the fake for method APICollectionsClient.GetByAzureAPIManagementService
	// HTTP status codes to indicate success: http.StatusOK
	GetByAzureAPIManagementService func(ctx context.Context, resourceGroupName string, serviceName string, apiID string, options *armsecurity.APICollectionsClientGetByAzureAPIManagementServiceOptions) (resp azfake.Responder[armsecurity.APICollectionsClientGetByAzureAPIManagementServiceResponse], errResp azfake.ErrorResponder)

	// NewListByAzureAPIManagementServicePager is the fake for method APICollectionsClient.NewListByAzureAPIManagementServicePager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByAzureAPIManagementServicePager func(resourceGroupName string, serviceName string, options *armsecurity.APICollectionsClientListByAzureAPIManagementServiceOptions) (resp azfake.PagerResponder[armsecurity.APICollectionsClientListByAzureAPIManagementServiceResponse])

	// NewListByResourceGroupPager is the fake for method APICollectionsClient.NewListByResourceGroupPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByResourceGroupPager func(resourceGroupName string, options *armsecurity.APICollectionsClientListByResourceGroupOptions) (resp azfake.PagerResponder[armsecurity.APICollectionsClientListByResourceGroupResponse])

	// NewListBySubscriptionPager is the fake for method APICollectionsClient.NewListBySubscriptionPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListBySubscriptionPager func(options *armsecurity.APICollectionsClientListBySubscriptionOptions) (resp azfake.PagerResponder[armsecurity.APICollectionsClientListBySubscriptionResponse])

	// OffboardAzureAPIManagementAPI is the fake for method APICollectionsClient.OffboardAzureAPIManagementAPI
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	OffboardAzureAPIManagementAPI func(ctx context.Context, resourceGroupName string, serviceName string, apiID string, options *armsecurity.APICollectionsClientOffboardAzureAPIManagementAPIOptions) (resp azfake.Responder[armsecurity.APICollectionsClientOffboardAzureAPIManagementAPIResponse], errResp azfake.ErrorResponder)

	// BeginOnboardAzureAPIManagementAPI is the fake for method APICollectionsClient.BeginOnboardAzureAPIManagementAPI
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginOnboardAzureAPIManagementAPI func(ctx context.Context, resourceGroupName string, serviceName string, apiID string, options *armsecurity.APICollectionsClientBeginOnboardAzureAPIManagementAPIOptions) (resp azfake.PollerResponder[armsecurity.APICollectionsClientOnboardAzureAPIManagementAPIResponse], errResp azfake.ErrorResponder)
}

// NewAPICollectionsServerTransport creates a new instance of APICollectionsServerTransport with the provided implementation.
// The returned APICollectionsServerTransport instance is connected to an instance of armsecurity.APICollectionsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewAPICollectionsServerTransport(srv *APICollectionsServer) *APICollectionsServerTransport {
	return &APICollectionsServerTransport{
		srv:                                     srv,
		newListByAzureAPIManagementServicePager: newTracker[azfake.PagerResponder[armsecurity.APICollectionsClientListByAzureAPIManagementServiceResponse]](),
		newListByResourceGroupPager:             newTracker[azfake.PagerResponder[armsecurity.APICollectionsClientListByResourceGroupResponse]](),
		newListBySubscriptionPager:              newTracker[azfake.PagerResponder[armsecurity.APICollectionsClientListBySubscriptionResponse]](),
		beginOnboardAzureAPIManagementAPI:       newTracker[azfake.PollerResponder[armsecurity.APICollectionsClientOnboardAzureAPIManagementAPIResponse]](),
	}
}

// APICollectionsServerTransport connects instances of armsecurity.APICollectionsClient to instances of APICollectionsServer.
// Don't use this type directly, use NewAPICollectionsServerTransport instead.
type APICollectionsServerTransport struct {
	srv                                     *APICollectionsServer
	newListByAzureAPIManagementServicePager *tracker[azfake.PagerResponder[armsecurity.APICollectionsClientListByAzureAPIManagementServiceResponse]]
	newListByResourceGroupPager             *tracker[azfake.PagerResponder[armsecurity.APICollectionsClientListByResourceGroupResponse]]
	newListBySubscriptionPager              *tracker[azfake.PagerResponder[armsecurity.APICollectionsClientListBySubscriptionResponse]]
	beginOnboardAzureAPIManagementAPI       *tracker[azfake.PollerResponder[armsecurity.APICollectionsClientOnboardAzureAPIManagementAPIResponse]]
}

// Do implements the policy.Transporter interface for APICollectionsServerTransport.
func (a *APICollectionsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "APICollectionsClient.GetByAzureAPIManagementService":
		resp, err = a.dispatchGetByAzureAPIManagementService(req)
	case "APICollectionsClient.NewListByAzureAPIManagementServicePager":
		resp, err = a.dispatchNewListByAzureAPIManagementServicePager(req)
	case "APICollectionsClient.NewListByResourceGroupPager":
		resp, err = a.dispatchNewListByResourceGroupPager(req)
	case "APICollectionsClient.NewListBySubscriptionPager":
		resp, err = a.dispatchNewListBySubscriptionPager(req)
	case "APICollectionsClient.OffboardAzureAPIManagementAPI":
		resp, err = a.dispatchOffboardAzureAPIManagementAPI(req)
	case "APICollectionsClient.BeginOnboardAzureAPIManagementAPI":
		resp, err = a.dispatchBeginOnboardAzureAPIManagementAPI(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *APICollectionsServerTransport) dispatchGetByAzureAPIManagementService(req *http.Request) (*http.Response, error) {
	if a.srv.GetByAzureAPIManagementService == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetByAzureAPIManagementService not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Security/apiCollections/(?P<apiId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	serviceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serviceName")])
	if err != nil {
		return nil, err
	}
	apiIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("apiId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := a.srv.GetByAzureAPIManagementService(req.Context(), resourceGroupNameParam, serviceNameParam, apiIDParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).APICollection, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *APICollectionsServerTransport) dispatchNewListByAzureAPIManagementServicePager(req *http.Request) (*http.Response, error) {
	if a.srv.NewListByAzureAPIManagementServicePager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByAzureAPIManagementServicePager not implemented")}
	}
	newListByAzureAPIManagementServicePager := a.newListByAzureAPIManagementServicePager.get(req)
	if newListByAzureAPIManagementServicePager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Security/apiCollections`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		serviceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serviceName")])
		if err != nil {
			return nil, err
		}
		resp := a.srv.NewListByAzureAPIManagementServicePager(resourceGroupNameParam, serviceNameParam, nil)
		newListByAzureAPIManagementServicePager = &resp
		a.newListByAzureAPIManagementServicePager.add(req, newListByAzureAPIManagementServicePager)
		server.PagerResponderInjectNextLinks(newListByAzureAPIManagementServicePager, req, func(page *armsecurity.APICollectionsClientListByAzureAPIManagementServiceResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByAzureAPIManagementServicePager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		a.newListByAzureAPIManagementServicePager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByAzureAPIManagementServicePager) {
		a.newListByAzureAPIManagementServicePager.remove(req)
	}
	return resp, nil
}

func (a *APICollectionsServerTransport) dispatchNewListByResourceGroupPager(req *http.Request) (*http.Response, error) {
	if a.srv.NewListByResourceGroupPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByResourceGroupPager not implemented")}
	}
	newListByResourceGroupPager := a.newListByResourceGroupPager.get(req)
	if newListByResourceGroupPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Security/apiCollections`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		resp := a.srv.NewListByResourceGroupPager(resourceGroupNameParam, nil)
		newListByResourceGroupPager = &resp
		a.newListByResourceGroupPager.add(req, newListByResourceGroupPager)
		server.PagerResponderInjectNextLinks(newListByResourceGroupPager, req, func(page *armsecurity.APICollectionsClientListByResourceGroupResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByResourceGroupPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		a.newListByResourceGroupPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByResourceGroupPager) {
		a.newListByResourceGroupPager.remove(req)
	}
	return resp, nil
}

func (a *APICollectionsServerTransport) dispatchNewListBySubscriptionPager(req *http.Request) (*http.Response, error) {
	if a.srv.NewListBySubscriptionPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListBySubscriptionPager not implemented")}
	}
	newListBySubscriptionPager := a.newListBySubscriptionPager.get(req)
	if newListBySubscriptionPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Security/apiCollections`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resp := a.srv.NewListBySubscriptionPager(nil)
		newListBySubscriptionPager = &resp
		a.newListBySubscriptionPager.add(req, newListBySubscriptionPager)
		server.PagerResponderInjectNextLinks(newListBySubscriptionPager, req, func(page *armsecurity.APICollectionsClientListBySubscriptionResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListBySubscriptionPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		a.newListBySubscriptionPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListBySubscriptionPager) {
		a.newListBySubscriptionPager.remove(req)
	}
	return resp, nil
}

func (a *APICollectionsServerTransport) dispatchOffboardAzureAPIManagementAPI(req *http.Request) (*http.Response, error) {
	if a.srv.OffboardAzureAPIManagementAPI == nil {
		return nil, &nonRetriableError{errors.New("fake for method OffboardAzureAPIManagementAPI not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Security/apiCollections/(?P<apiId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	serviceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serviceName")])
	if err != nil {
		return nil, err
	}
	apiIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("apiId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := a.srv.OffboardAzureAPIManagementAPI(req.Context(), resourceGroupNameParam, serviceNameParam, apiIDParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusNoContent}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusNoContent", respContent.HTTPStatus)}
	}
	resp, err := server.NewResponse(respContent, req, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *APICollectionsServerTransport) dispatchBeginOnboardAzureAPIManagementAPI(req *http.Request) (*http.Response, error) {
	if a.srv.BeginOnboardAzureAPIManagementAPI == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginOnboardAzureAPIManagementAPI not implemented")}
	}
	beginOnboardAzureAPIManagementAPI := a.beginOnboardAzureAPIManagementAPI.get(req)
	if beginOnboardAzureAPIManagementAPI == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Security/apiCollections/(?P<apiId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		serviceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serviceName")])
		if err != nil {
			return nil, err
		}
		apiIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("apiId")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := a.srv.BeginOnboardAzureAPIManagementAPI(req.Context(), resourceGroupNameParam, serviceNameParam, apiIDParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginOnboardAzureAPIManagementAPI = &respr
		a.beginOnboardAzureAPIManagementAPI.add(req, beginOnboardAzureAPIManagementAPI)
	}

	resp, err := server.PollerResponderNext(beginOnboardAzureAPIManagementAPI, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		a.beginOnboardAzureAPIManagementAPI.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginOnboardAzureAPIManagementAPI) {
		a.beginOnboardAzureAPIManagementAPI.remove(req)
	}

	return resp, nil
}