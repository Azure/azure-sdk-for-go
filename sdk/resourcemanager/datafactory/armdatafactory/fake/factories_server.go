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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datafactory/armdatafactory/v10"
	"net/http"
	"net/url"
	"regexp"
)

// FactoriesServer is a fake server for instances of the armdatafactory.FactoriesClient type.
type FactoriesServer struct {
	// ConfigureFactoryRepo is the fake for method FactoriesClient.ConfigureFactoryRepo
	// HTTP status codes to indicate success: http.StatusOK
	ConfigureFactoryRepo func(ctx context.Context, locationID string, factoryRepoUpdate armdatafactory.FactoryRepoUpdate, options *armdatafactory.FactoriesClientConfigureFactoryRepoOptions) (resp azfake.Responder[armdatafactory.FactoriesClientConfigureFactoryRepoResponse], errResp azfake.ErrorResponder)

	// CreateOrUpdate is the fake for method FactoriesClient.CreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK
	CreateOrUpdate func(ctx context.Context, resourceGroupName string, factoryName string, factory armdatafactory.Factory, options *armdatafactory.FactoriesClientCreateOrUpdateOptions) (resp azfake.Responder[armdatafactory.FactoriesClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// Delete is the fake for method FactoriesClient.Delete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	Delete func(ctx context.Context, resourceGroupName string, factoryName string, options *armdatafactory.FactoriesClientDeleteOptions) (resp azfake.Responder[armdatafactory.FactoriesClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method FactoriesClient.Get
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNotModified
	Get func(ctx context.Context, resourceGroupName string, factoryName string, options *armdatafactory.FactoriesClientGetOptions) (resp azfake.Responder[armdatafactory.FactoriesClientGetResponse], errResp azfake.ErrorResponder)

	// GetDataPlaneAccess is the fake for method FactoriesClient.GetDataPlaneAccess
	// HTTP status codes to indicate success: http.StatusOK
	GetDataPlaneAccess func(ctx context.Context, resourceGroupName string, factoryName string, policy armdatafactory.UserAccessPolicy, options *armdatafactory.FactoriesClientGetDataPlaneAccessOptions) (resp azfake.Responder[armdatafactory.FactoriesClientGetDataPlaneAccessResponse], errResp azfake.ErrorResponder)

	// GetGitHubAccessToken is the fake for method FactoriesClient.GetGitHubAccessToken
	// HTTP status codes to indicate success: http.StatusOK
	GetGitHubAccessToken func(ctx context.Context, resourceGroupName string, factoryName string, gitHubAccessTokenRequest armdatafactory.GitHubAccessTokenRequest, options *armdatafactory.FactoriesClientGetGitHubAccessTokenOptions) (resp azfake.Responder[armdatafactory.FactoriesClientGetGitHubAccessTokenResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method FactoriesClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(options *armdatafactory.FactoriesClientListOptions) (resp azfake.PagerResponder[armdatafactory.FactoriesClientListResponse])

	// NewListByResourceGroupPager is the fake for method FactoriesClient.NewListByResourceGroupPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByResourceGroupPager func(resourceGroupName string, options *armdatafactory.FactoriesClientListByResourceGroupOptions) (resp azfake.PagerResponder[armdatafactory.FactoriesClientListByResourceGroupResponse])

	// Update is the fake for method FactoriesClient.Update
	// HTTP status codes to indicate success: http.StatusOK
	Update func(ctx context.Context, resourceGroupName string, factoryName string, factoryUpdateParameters armdatafactory.FactoryUpdateParameters, options *armdatafactory.FactoriesClientUpdateOptions) (resp azfake.Responder[armdatafactory.FactoriesClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewFactoriesServerTransport creates a new instance of FactoriesServerTransport with the provided implementation.
// The returned FactoriesServerTransport instance is connected to an instance of armdatafactory.FactoriesClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewFactoriesServerTransport(srv *FactoriesServer) *FactoriesServerTransport {
	return &FactoriesServerTransport{
		srv:                         srv,
		newListPager:                newTracker[azfake.PagerResponder[armdatafactory.FactoriesClientListResponse]](),
		newListByResourceGroupPager: newTracker[azfake.PagerResponder[armdatafactory.FactoriesClientListByResourceGroupResponse]](),
	}
}

// FactoriesServerTransport connects instances of armdatafactory.FactoriesClient to instances of FactoriesServer.
// Don't use this type directly, use NewFactoriesServerTransport instead.
type FactoriesServerTransport struct {
	srv                         *FactoriesServer
	newListPager                *tracker[azfake.PagerResponder[armdatafactory.FactoriesClientListResponse]]
	newListByResourceGroupPager *tracker[azfake.PagerResponder[armdatafactory.FactoriesClientListByResourceGroupResponse]]
}

// Do implements the policy.Transporter interface for FactoriesServerTransport.
func (f *FactoriesServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return f.dispatchToMethodFake(req, method)
}

func (f *FactoriesServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if factoriesServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = factoriesServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "FactoriesClient.ConfigureFactoryRepo":
				res.resp, res.err = f.dispatchConfigureFactoryRepo(req)
			case "FactoriesClient.CreateOrUpdate":
				res.resp, res.err = f.dispatchCreateOrUpdate(req)
			case "FactoriesClient.Delete":
				res.resp, res.err = f.dispatchDelete(req)
			case "FactoriesClient.Get":
				res.resp, res.err = f.dispatchGet(req)
			case "FactoriesClient.GetDataPlaneAccess":
				res.resp, res.err = f.dispatchGetDataPlaneAccess(req)
			case "FactoriesClient.GetGitHubAccessToken":
				res.resp, res.err = f.dispatchGetGitHubAccessToken(req)
			case "FactoriesClient.NewListPager":
				res.resp, res.err = f.dispatchNewListPager(req)
			case "FactoriesClient.NewListByResourceGroupPager":
				res.resp, res.err = f.dispatchNewListByResourceGroupPager(req)
			case "FactoriesClient.Update":
				res.resp, res.err = f.dispatchUpdate(req)
			default:
				res.err = fmt.Errorf("unhandled API %s", method)
			}

		}
		select {
		case resultChan <- res:
		case <-req.Context().Done():
		}
	}()

	select {
	case <-req.Context().Done():
		return nil, req.Context().Err()
	case res := <-resultChan:
		return res.resp, res.err
	}
}

func (f *FactoriesServerTransport) dispatchConfigureFactoryRepo(req *http.Request) (*http.Response, error) {
	if f.srv.ConfigureFactoryRepo == nil {
		return nil, &nonRetriableError{errors.New("fake for method ConfigureFactoryRepo not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DataFactory/locations/(?P<locationId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/configureFactoryRepo`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armdatafactory.FactoryRepoUpdate](req)
	if err != nil {
		return nil, err
	}
	locationIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("locationId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := f.srv.ConfigureFactoryRepo(req.Context(), locationIDParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).Factory, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (f *FactoriesServerTransport) dispatchCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if f.srv.CreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method CreateOrUpdate not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DataFactory/factories/(?P<factoryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armdatafactory.Factory](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	factoryNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("factoryName")])
	if err != nil {
		return nil, err
	}
	ifMatchParam := getOptional(getHeaderValue(req.Header, "If-Match"))
	var options *armdatafactory.FactoriesClientCreateOrUpdateOptions
	if ifMatchParam != nil {
		options = &armdatafactory.FactoriesClientCreateOrUpdateOptions{
			IfMatch: ifMatchParam,
		}
	}
	respr, errRespr := f.srv.CreateOrUpdate(req.Context(), resourceGroupNameParam, factoryNameParam, body, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).Factory, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (f *FactoriesServerTransport) dispatchDelete(req *http.Request) (*http.Response, error) {
	if f.srv.Delete == nil {
		return nil, &nonRetriableError{errors.New("fake for method Delete not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DataFactory/factories/(?P<factoryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	factoryNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("factoryName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := f.srv.Delete(req.Context(), resourceGroupNameParam, factoryNameParam, nil)
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

func (f *FactoriesServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if f.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DataFactory/factories/(?P<factoryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	factoryNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("factoryName")])
	if err != nil {
		return nil, err
	}
	ifNoneMatchParam := getOptional(getHeaderValue(req.Header, "If-None-Match"))
	var options *armdatafactory.FactoriesClientGetOptions
	if ifNoneMatchParam != nil {
		options = &armdatafactory.FactoriesClientGetOptions{
			IfNoneMatch: ifNoneMatchParam,
		}
	}
	respr, errRespr := f.srv.Get(req.Context(), resourceGroupNameParam, factoryNameParam, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusNotModified}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusNotModified", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).Factory, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (f *FactoriesServerTransport) dispatchGetDataPlaneAccess(req *http.Request) (*http.Response, error) {
	if f.srv.GetDataPlaneAccess == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetDataPlaneAccess not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DataFactory/factories/(?P<factoryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/getDataPlaneAccess`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armdatafactory.UserAccessPolicy](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	factoryNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("factoryName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := f.srv.GetDataPlaneAccess(req.Context(), resourceGroupNameParam, factoryNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).AccessPolicyResponse, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (f *FactoriesServerTransport) dispatchGetGitHubAccessToken(req *http.Request) (*http.Response, error) {
	if f.srv.GetGitHubAccessToken == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetGitHubAccessToken not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DataFactory/factories/(?P<factoryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/getGitHubAccessToken`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armdatafactory.GitHubAccessTokenRequest](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	factoryNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("factoryName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := f.srv.GetGitHubAccessToken(req.Context(), resourceGroupNameParam, factoryNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).GitHubAccessTokenResponse, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (f *FactoriesServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if f.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := f.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DataFactory/factories`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resp := f.srv.NewListPager(nil)
		newListPager = &resp
		f.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armdatafactory.FactoriesClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		f.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		f.newListPager.remove(req)
	}
	return resp, nil
}

func (f *FactoriesServerTransport) dispatchNewListByResourceGroupPager(req *http.Request) (*http.Response, error) {
	if f.srv.NewListByResourceGroupPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByResourceGroupPager not implemented")}
	}
	newListByResourceGroupPager := f.newListByResourceGroupPager.get(req)
	if newListByResourceGroupPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DataFactory/factories`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		resp := f.srv.NewListByResourceGroupPager(resourceGroupNameParam, nil)
		newListByResourceGroupPager = &resp
		f.newListByResourceGroupPager.add(req, newListByResourceGroupPager)
		server.PagerResponderInjectNextLinks(newListByResourceGroupPager, req, func(page *armdatafactory.FactoriesClientListByResourceGroupResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByResourceGroupPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		f.newListByResourceGroupPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByResourceGroupPager) {
		f.newListByResourceGroupPager.remove(req)
	}
	return resp, nil
}

func (f *FactoriesServerTransport) dispatchUpdate(req *http.Request) (*http.Response, error) {
	if f.srv.Update == nil {
		return nil, &nonRetriableError{errors.New("fake for method Update not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DataFactory/factories/(?P<factoryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armdatafactory.FactoryUpdateParameters](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	factoryNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("factoryName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := f.srv.Update(req.Context(), resourceGroupNameParam, factoryNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).Factory, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to FactoriesServerTransport
var factoriesServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
