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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement/v3"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

// AuthorizationServer is a fake server for instances of the armapimanagement.AuthorizationClient type.
type AuthorizationServer struct {
	// ConfirmConsentCode is the fake for method AuthorizationClient.ConfirmConsentCode
	// HTTP status codes to indicate success: http.StatusOK
	ConfirmConsentCode func(ctx context.Context, resourceGroupName string, serviceName string, authorizationProviderID string, authorizationID string, parameters armapimanagement.AuthorizationConfirmConsentCodeRequestContract, options *armapimanagement.AuthorizationClientConfirmConsentCodeOptions) (resp azfake.Responder[armapimanagement.AuthorizationClientConfirmConsentCodeResponse], errResp azfake.ErrorResponder)

	// CreateOrUpdate is the fake for method AuthorizationClient.CreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	CreateOrUpdate func(ctx context.Context, resourceGroupName string, serviceName string, authorizationProviderID string, authorizationID string, parameters armapimanagement.AuthorizationContract, options *armapimanagement.AuthorizationClientCreateOrUpdateOptions) (resp azfake.Responder[armapimanagement.AuthorizationClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// Delete is the fake for method AuthorizationClient.Delete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	Delete func(ctx context.Context, resourceGroupName string, serviceName string, authorizationProviderID string, authorizationID string, ifMatch string, options *armapimanagement.AuthorizationClientDeleteOptions) (resp azfake.Responder[armapimanagement.AuthorizationClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method AuthorizationClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, serviceName string, authorizationProviderID string, authorizationID string, options *armapimanagement.AuthorizationClientGetOptions) (resp azfake.Responder[armapimanagement.AuthorizationClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByAuthorizationProviderPager is the fake for method AuthorizationClient.NewListByAuthorizationProviderPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByAuthorizationProviderPager func(resourceGroupName string, serviceName string, authorizationProviderID string, options *armapimanagement.AuthorizationClientListByAuthorizationProviderOptions) (resp azfake.PagerResponder[armapimanagement.AuthorizationClientListByAuthorizationProviderResponse])
}

// NewAuthorizationServerTransport creates a new instance of AuthorizationServerTransport with the provided implementation.
// The returned AuthorizationServerTransport instance is connected to an instance of armapimanagement.AuthorizationClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewAuthorizationServerTransport(srv *AuthorizationServer) *AuthorizationServerTransport {
	return &AuthorizationServerTransport{
		srv:                                 srv,
		newListByAuthorizationProviderPager: newTracker[azfake.PagerResponder[armapimanagement.AuthorizationClientListByAuthorizationProviderResponse]](),
	}
}

// AuthorizationServerTransport connects instances of armapimanagement.AuthorizationClient to instances of AuthorizationServer.
// Don't use this type directly, use NewAuthorizationServerTransport instead.
type AuthorizationServerTransport struct {
	srv                                 *AuthorizationServer
	newListByAuthorizationProviderPager *tracker[azfake.PagerResponder[armapimanagement.AuthorizationClientListByAuthorizationProviderResponse]]
}

// Do implements the policy.Transporter interface for AuthorizationServerTransport.
func (a *AuthorizationServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return a.dispatchToMethodFake(req, method)
}

func (a *AuthorizationServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if authorizationServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = authorizationServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "AuthorizationClient.ConfirmConsentCode":
				res.resp, res.err = a.dispatchConfirmConsentCode(req)
			case "AuthorizationClient.CreateOrUpdate":
				res.resp, res.err = a.dispatchCreateOrUpdate(req)
			case "AuthorizationClient.Delete":
				res.resp, res.err = a.dispatchDelete(req)
			case "AuthorizationClient.Get":
				res.resp, res.err = a.dispatchGet(req)
			case "AuthorizationClient.NewListByAuthorizationProviderPager":
				res.resp, res.err = a.dispatchNewListByAuthorizationProviderPager(req)
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

func (a *AuthorizationServerTransport) dispatchConfirmConsentCode(req *http.Request) (*http.Response, error) {
	if a.srv.ConfirmConsentCode == nil {
		return nil, &nonRetriableError{errors.New("fake for method ConfirmConsentCode not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/authorizationProviders/(?P<authorizationProviderId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/authorizations/(?P<authorizationId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/confirmConsentCode`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 5 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armapimanagement.AuthorizationConfirmConsentCodeRequestContract](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	serviceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serviceName")])
	if err != nil {
		return nil, err
	}
	authorizationProviderIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("authorizationProviderId")])
	if err != nil {
		return nil, err
	}
	authorizationIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("authorizationId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := a.srv.ConfirmConsentCode(req.Context(), resourceGroupNameParam, serviceNameParam, authorizationProviderIDParam, authorizationIDParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.NewResponse(respContent, req, nil)
	if err != nil {
		return nil, err
	}
	if val := server.GetResponse(respr).ETag; val != nil {
		resp.Header.Set("ETag", *val)
	}
	return resp, nil
}

func (a *AuthorizationServerTransport) dispatchCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if a.srv.CreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method CreateOrUpdate not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/authorizationProviders/(?P<authorizationProviderId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/authorizations/(?P<authorizationId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 5 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armapimanagement.AuthorizationContract](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	serviceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serviceName")])
	if err != nil {
		return nil, err
	}
	authorizationProviderIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("authorizationProviderId")])
	if err != nil {
		return nil, err
	}
	authorizationIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("authorizationId")])
	if err != nil {
		return nil, err
	}
	ifMatchParam := getOptional(getHeaderValue(req.Header, "If-Match"))
	var options *armapimanagement.AuthorizationClientCreateOrUpdateOptions
	if ifMatchParam != nil {
		options = &armapimanagement.AuthorizationClientCreateOrUpdateOptions{
			IfMatch: ifMatchParam,
		}
	}
	respr, errRespr := a.srv.CreateOrUpdate(req.Context(), resourceGroupNameParam, serviceNameParam, authorizationProviderIDParam, authorizationIDParam, body, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusCreated}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).AuthorizationContract, req)
	if err != nil {
		return nil, err
	}
	if val := server.GetResponse(respr).ETag; val != nil {
		resp.Header.Set("ETag", *val)
	}
	return resp, nil
}

func (a *AuthorizationServerTransport) dispatchDelete(req *http.Request) (*http.Response, error) {
	if a.srv.Delete == nil {
		return nil, &nonRetriableError{errors.New("fake for method Delete not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/authorizationProviders/(?P<authorizationProviderId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/authorizations/(?P<authorizationId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 5 {
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
	authorizationProviderIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("authorizationProviderId")])
	if err != nil {
		return nil, err
	}
	authorizationIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("authorizationId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := a.srv.Delete(req.Context(), resourceGroupNameParam, serviceNameParam, authorizationProviderIDParam, authorizationIDParam, getHeaderValue(req.Header, "If-Match"), nil)
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

func (a *AuthorizationServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if a.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/authorizationProviders/(?P<authorizationProviderId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/authorizations/(?P<authorizationId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 5 {
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
	authorizationProviderIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("authorizationProviderId")])
	if err != nil {
		return nil, err
	}
	authorizationIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("authorizationId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := a.srv.Get(req.Context(), resourceGroupNameParam, serviceNameParam, authorizationProviderIDParam, authorizationIDParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).AuthorizationContract, req)
	if err != nil {
		return nil, err
	}
	if val := server.GetResponse(respr).ETag; val != nil {
		resp.Header.Set("ETag", *val)
	}
	return resp, nil
}

func (a *AuthorizationServerTransport) dispatchNewListByAuthorizationProviderPager(req *http.Request) (*http.Response, error) {
	if a.srv.NewListByAuthorizationProviderPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByAuthorizationProviderPager not implemented")}
	}
	newListByAuthorizationProviderPager := a.newListByAuthorizationProviderPager.get(req)
	if newListByAuthorizationProviderPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/authorizationProviders/(?P<authorizationProviderId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/authorizations`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		serviceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serviceName")])
		if err != nil {
			return nil, err
		}
		authorizationProviderIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("authorizationProviderId")])
		if err != nil {
			return nil, err
		}
		filterUnescaped, err := url.QueryUnescape(qp.Get("$filter"))
		if err != nil {
			return nil, err
		}
		filterParam := getOptional(filterUnescaped)
		topUnescaped, err := url.QueryUnescape(qp.Get("$top"))
		if err != nil {
			return nil, err
		}
		topParam, err := parseOptional(topUnescaped, func(v string) (int32, error) {
			p, parseErr := strconv.ParseInt(v, 10, 32)
			if parseErr != nil {
				return 0, parseErr
			}
			return int32(p), nil
		})
		if err != nil {
			return nil, err
		}
		skipUnescaped, err := url.QueryUnescape(qp.Get("$skip"))
		if err != nil {
			return nil, err
		}
		skipParam, err := parseOptional(skipUnescaped, func(v string) (int32, error) {
			p, parseErr := strconv.ParseInt(v, 10, 32)
			if parseErr != nil {
				return 0, parseErr
			}
			return int32(p), nil
		})
		if err != nil {
			return nil, err
		}
		var options *armapimanagement.AuthorizationClientListByAuthorizationProviderOptions
		if filterParam != nil || topParam != nil || skipParam != nil {
			options = &armapimanagement.AuthorizationClientListByAuthorizationProviderOptions{
				Filter: filterParam,
				Top:    topParam,
				Skip:   skipParam,
			}
		}
		resp := a.srv.NewListByAuthorizationProviderPager(resourceGroupNameParam, serviceNameParam, authorizationProviderIDParam, options)
		newListByAuthorizationProviderPager = &resp
		a.newListByAuthorizationProviderPager.add(req, newListByAuthorizationProviderPager)
		server.PagerResponderInjectNextLinks(newListByAuthorizationProviderPager, req, func(page *armapimanagement.AuthorizationClientListByAuthorizationProviderResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByAuthorizationProviderPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		a.newListByAuthorizationProviderPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByAuthorizationProviderPager) {
		a.newListByAuthorizationProviderPager.remove(req)
	}
	return resp, nil
}
