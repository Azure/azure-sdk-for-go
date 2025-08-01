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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/providerhub/armproviderhub/v2"
	"net/http"
	"net/url"
	"regexp"
)

// CustomRolloutsServer is a fake server for instances of the armproviderhub.CustomRolloutsClient type.
type CustomRolloutsServer struct {
	// BeginCreateOrUpdate is the fake for method CustomRolloutsClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreateOrUpdate func(ctx context.Context, providerNamespace string, rolloutName string, properties armproviderhub.CustomRollout, options *armproviderhub.CustomRolloutsClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armproviderhub.CustomRolloutsClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// Delete is the fake for method CustomRolloutsClient.Delete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	Delete func(ctx context.Context, providerNamespace string, rolloutName string, options *armproviderhub.CustomRolloutsClientDeleteOptions) (resp azfake.Responder[armproviderhub.CustomRolloutsClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method CustomRolloutsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, providerNamespace string, rolloutName string, options *armproviderhub.CustomRolloutsClientGetOptions) (resp azfake.Responder[armproviderhub.CustomRolloutsClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByProviderRegistrationPager is the fake for method CustomRolloutsClient.NewListByProviderRegistrationPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByProviderRegistrationPager func(providerNamespace string, options *armproviderhub.CustomRolloutsClientListByProviderRegistrationOptions) (resp azfake.PagerResponder[armproviderhub.CustomRolloutsClientListByProviderRegistrationResponse])

	// Stop is the fake for method CustomRolloutsClient.Stop
	// HTTP status codes to indicate success: http.StatusOK
	Stop func(ctx context.Context, providerNamespace string, rolloutName string, options *armproviderhub.CustomRolloutsClientStopOptions) (resp azfake.Responder[armproviderhub.CustomRolloutsClientStopResponse], errResp azfake.ErrorResponder)
}

// NewCustomRolloutsServerTransport creates a new instance of CustomRolloutsServerTransport with the provided implementation.
// The returned CustomRolloutsServerTransport instance is connected to an instance of armproviderhub.CustomRolloutsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewCustomRolloutsServerTransport(srv *CustomRolloutsServer) *CustomRolloutsServerTransport {
	return &CustomRolloutsServerTransport{
		srv:                                srv,
		beginCreateOrUpdate:                newTracker[azfake.PollerResponder[armproviderhub.CustomRolloutsClientCreateOrUpdateResponse]](),
		newListByProviderRegistrationPager: newTracker[azfake.PagerResponder[armproviderhub.CustomRolloutsClientListByProviderRegistrationResponse]](),
	}
}

// CustomRolloutsServerTransport connects instances of armproviderhub.CustomRolloutsClient to instances of CustomRolloutsServer.
// Don't use this type directly, use NewCustomRolloutsServerTransport instead.
type CustomRolloutsServerTransport struct {
	srv                                *CustomRolloutsServer
	beginCreateOrUpdate                *tracker[azfake.PollerResponder[armproviderhub.CustomRolloutsClientCreateOrUpdateResponse]]
	newListByProviderRegistrationPager *tracker[azfake.PagerResponder[armproviderhub.CustomRolloutsClientListByProviderRegistrationResponse]]
}

// Do implements the policy.Transporter interface for CustomRolloutsServerTransport.
func (c *CustomRolloutsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return c.dispatchToMethodFake(req, method)
}

func (c *CustomRolloutsServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if customRolloutsServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = customRolloutsServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "CustomRolloutsClient.BeginCreateOrUpdate":
				res.resp, res.err = c.dispatchBeginCreateOrUpdate(req)
			case "CustomRolloutsClient.Delete":
				res.resp, res.err = c.dispatchDelete(req)
			case "CustomRolloutsClient.Get":
				res.resp, res.err = c.dispatchGet(req)
			case "CustomRolloutsClient.NewListByProviderRegistrationPager":
				res.resp, res.err = c.dispatchNewListByProviderRegistrationPager(req)
			case "CustomRolloutsClient.Stop":
				res.resp, res.err = c.dispatchStop(req)
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

func (c *CustomRolloutsServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if c.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := c.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ProviderHub/providerRegistrations/(?P<providerNamespace>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/customRollouts/(?P<rolloutName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armproviderhub.CustomRollout](req)
		if err != nil {
			return nil, err
		}
		providerNamespaceParam, err := url.PathUnescape(matches[regex.SubexpIndex("providerNamespace")])
		if err != nil {
			return nil, err
		}
		rolloutNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("rolloutName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := c.srv.BeginCreateOrUpdate(req.Context(), providerNamespaceParam, rolloutNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateOrUpdate = &respr
		c.beginCreateOrUpdate.add(req, beginCreateOrUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		c.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		c.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (c *CustomRolloutsServerTransport) dispatchDelete(req *http.Request) (*http.Response, error) {
	if c.srv.Delete == nil {
		return nil, &nonRetriableError{errors.New("fake for method Delete not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ProviderHub/providerRegistrations/(?P<providerNamespace>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/customRollouts/(?P<rolloutName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	providerNamespaceParam, err := url.PathUnescape(matches[regex.SubexpIndex("providerNamespace")])
	if err != nil {
		return nil, err
	}
	rolloutNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("rolloutName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := c.srv.Delete(req.Context(), providerNamespaceParam, rolloutNameParam, nil)
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

func (c *CustomRolloutsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if c.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ProviderHub/providerRegistrations/(?P<providerNamespace>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/customRollouts/(?P<rolloutName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	providerNamespaceParam, err := url.PathUnescape(matches[regex.SubexpIndex("providerNamespace")])
	if err != nil {
		return nil, err
	}
	rolloutNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("rolloutName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := c.srv.Get(req.Context(), providerNamespaceParam, rolloutNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).CustomRollout, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *CustomRolloutsServerTransport) dispatchNewListByProviderRegistrationPager(req *http.Request) (*http.Response, error) {
	if c.srv.NewListByProviderRegistrationPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByProviderRegistrationPager not implemented")}
	}
	newListByProviderRegistrationPager := c.newListByProviderRegistrationPager.get(req)
	if newListByProviderRegistrationPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ProviderHub/providerRegistrations/(?P<providerNamespace>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/customRollouts`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		providerNamespaceParam, err := url.PathUnescape(matches[regex.SubexpIndex("providerNamespace")])
		if err != nil {
			return nil, err
		}
		resp := c.srv.NewListByProviderRegistrationPager(providerNamespaceParam, nil)
		newListByProviderRegistrationPager = &resp
		c.newListByProviderRegistrationPager.add(req, newListByProviderRegistrationPager)
		server.PagerResponderInjectNextLinks(newListByProviderRegistrationPager, req, func(page *armproviderhub.CustomRolloutsClientListByProviderRegistrationResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByProviderRegistrationPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		c.newListByProviderRegistrationPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByProviderRegistrationPager) {
		c.newListByProviderRegistrationPager.remove(req)
	}
	return resp, nil
}

func (c *CustomRolloutsServerTransport) dispatchStop(req *http.Request) (*http.Response, error) {
	if c.srv.Stop == nil {
		return nil, &nonRetriableError{errors.New("fake for method Stop not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ProviderHub/providerRegistrations/(?P<providerNamespace>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/customRollouts/(?P<rolloutName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/stop`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	providerNamespaceParam, err := url.PathUnescape(matches[regex.SubexpIndex("providerNamespace")])
	if err != nil {
		return nil, err
	}
	rolloutNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("rolloutName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := c.srv.Stop(req.Context(), providerNamespaceParam, rolloutNameParam, nil)
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
	return resp, nil
}

// set this to conditionally intercept incoming requests to CustomRolloutsServerTransport
var customRolloutsServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
