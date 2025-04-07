// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package fake

import (
	"context"
	"errors"
	"fmt"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake/server"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hybridconnectivity/armhybridconnectivity"
	"net/http"
	"net/url"
	"regexp"
)

// PublicCloudConnectorsServer is a fake server for instances of the armhybridconnectivity.PublicCloudConnectorsClient type.
type PublicCloudConnectorsServer struct {
	// BeginCreateOrUpdate is the fake for method PublicCloudConnectorsClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, publicCloudConnector string, resource armhybridconnectivity.PublicCloudConnector, options *armhybridconnectivity.PublicCloudConnectorsClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armhybridconnectivity.PublicCloudConnectorsClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// Delete is the fake for method PublicCloudConnectorsClient.Delete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	Delete func(ctx context.Context, resourceGroupName string, publicCloudConnector string, options *armhybridconnectivity.PublicCloudConnectorsClientDeleteOptions) (resp azfake.Responder[armhybridconnectivity.PublicCloudConnectorsClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method PublicCloudConnectorsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, publicCloudConnector string, options *armhybridconnectivity.PublicCloudConnectorsClientGetOptions) (resp azfake.Responder[armhybridconnectivity.PublicCloudConnectorsClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByResourceGroupPager is the fake for method PublicCloudConnectorsClient.NewListByResourceGroupPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByResourceGroupPager func(resourceGroupName string, options *armhybridconnectivity.PublicCloudConnectorsClientListByResourceGroupOptions) (resp azfake.PagerResponder[armhybridconnectivity.PublicCloudConnectorsClientListByResourceGroupResponse])

	// NewListBySubscriptionPager is the fake for method PublicCloudConnectorsClient.NewListBySubscriptionPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListBySubscriptionPager func(options *armhybridconnectivity.PublicCloudConnectorsClientListBySubscriptionOptions) (resp azfake.PagerResponder[armhybridconnectivity.PublicCloudConnectorsClientListBySubscriptionResponse])

	// BeginTestPermissions is the fake for method PublicCloudConnectorsClient.BeginTestPermissions
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginTestPermissions func(ctx context.Context, resourceGroupName string, publicCloudConnector string, options *armhybridconnectivity.PublicCloudConnectorsClientBeginTestPermissionsOptions) (resp azfake.PollerResponder[armhybridconnectivity.PublicCloudConnectorsClientTestPermissionsResponse], errResp azfake.ErrorResponder)

	// Update is the fake for method PublicCloudConnectorsClient.Update
	// HTTP status codes to indicate success: http.StatusOK
	Update func(ctx context.Context, resourceGroupName string, publicCloudConnector string, properties armhybridconnectivity.PublicCloudConnector, options *armhybridconnectivity.PublicCloudConnectorsClientUpdateOptions) (resp azfake.Responder[armhybridconnectivity.PublicCloudConnectorsClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewPublicCloudConnectorsServerTransport creates a new instance of PublicCloudConnectorsServerTransport with the provided implementation.
// The returned PublicCloudConnectorsServerTransport instance is connected to an instance of armhybridconnectivity.PublicCloudConnectorsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewPublicCloudConnectorsServerTransport(srv *PublicCloudConnectorsServer) *PublicCloudConnectorsServerTransport {
	return &PublicCloudConnectorsServerTransport{
		srv:                         srv,
		beginCreateOrUpdate:         newTracker[azfake.PollerResponder[armhybridconnectivity.PublicCloudConnectorsClientCreateOrUpdateResponse]](),
		newListByResourceGroupPager: newTracker[azfake.PagerResponder[armhybridconnectivity.PublicCloudConnectorsClientListByResourceGroupResponse]](),
		newListBySubscriptionPager:  newTracker[azfake.PagerResponder[armhybridconnectivity.PublicCloudConnectorsClientListBySubscriptionResponse]](),
		beginTestPermissions:        newTracker[azfake.PollerResponder[armhybridconnectivity.PublicCloudConnectorsClientTestPermissionsResponse]](),
	}
}

// PublicCloudConnectorsServerTransport connects instances of armhybridconnectivity.PublicCloudConnectorsClient to instances of PublicCloudConnectorsServer.
// Don't use this type directly, use NewPublicCloudConnectorsServerTransport instead.
type PublicCloudConnectorsServerTransport struct {
	srv                         *PublicCloudConnectorsServer
	beginCreateOrUpdate         *tracker[azfake.PollerResponder[armhybridconnectivity.PublicCloudConnectorsClientCreateOrUpdateResponse]]
	newListByResourceGroupPager *tracker[azfake.PagerResponder[armhybridconnectivity.PublicCloudConnectorsClientListByResourceGroupResponse]]
	newListBySubscriptionPager  *tracker[azfake.PagerResponder[armhybridconnectivity.PublicCloudConnectorsClientListBySubscriptionResponse]]
	beginTestPermissions        *tracker[azfake.PollerResponder[armhybridconnectivity.PublicCloudConnectorsClientTestPermissionsResponse]]
}

// Do implements the policy.Transporter interface for PublicCloudConnectorsServerTransport.
func (p *PublicCloudConnectorsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return p.dispatchToMethodFake(req, method)
}

func (p *PublicCloudConnectorsServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if publicCloudConnectorsServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = publicCloudConnectorsServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "PublicCloudConnectorsClient.BeginCreateOrUpdate":
				res.resp, res.err = p.dispatchBeginCreateOrUpdate(req)
			case "PublicCloudConnectorsClient.Delete":
				res.resp, res.err = p.dispatchDelete(req)
			case "PublicCloudConnectorsClient.Get":
				res.resp, res.err = p.dispatchGet(req)
			case "PublicCloudConnectorsClient.NewListByResourceGroupPager":
				res.resp, res.err = p.dispatchNewListByResourceGroupPager(req)
			case "PublicCloudConnectorsClient.NewListBySubscriptionPager":
				res.resp, res.err = p.dispatchNewListBySubscriptionPager(req)
			case "PublicCloudConnectorsClient.BeginTestPermissions":
				res.resp, res.err = p.dispatchBeginTestPermissions(req)
			case "PublicCloudConnectorsClient.Update":
				res.resp, res.err = p.dispatchUpdate(req)
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

func (p *PublicCloudConnectorsServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if p.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := p.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HybridConnectivity/publicCloudConnectors/(?P<publicCloudConnector>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armhybridconnectivity.PublicCloudConnector](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		publicCloudConnectorParam, err := url.PathUnescape(matches[regex.SubexpIndex("publicCloudConnector")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := p.srv.BeginCreateOrUpdate(req.Context(), resourceGroupNameParam, publicCloudConnectorParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateOrUpdate = &respr
		p.beginCreateOrUpdate.add(req, beginCreateOrUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		p.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		p.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (p *PublicCloudConnectorsServerTransport) dispatchDelete(req *http.Request) (*http.Response, error) {
	if p.srv.Delete == nil {
		return nil, &nonRetriableError{errors.New("fake for method Delete not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HybridConnectivity/publicCloudConnectors/(?P<publicCloudConnector>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	publicCloudConnectorParam, err := url.PathUnescape(matches[regex.SubexpIndex("publicCloudConnector")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := p.srv.Delete(req.Context(), resourceGroupNameParam, publicCloudConnectorParam, nil)
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

func (p *PublicCloudConnectorsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if p.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HybridConnectivity/publicCloudConnectors/(?P<publicCloudConnector>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	publicCloudConnectorParam, err := url.PathUnescape(matches[regex.SubexpIndex("publicCloudConnector")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := p.srv.Get(req.Context(), resourceGroupNameParam, publicCloudConnectorParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).PublicCloudConnector, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (p *PublicCloudConnectorsServerTransport) dispatchNewListByResourceGroupPager(req *http.Request) (*http.Response, error) {
	if p.srv.NewListByResourceGroupPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByResourceGroupPager not implemented")}
	}
	newListByResourceGroupPager := p.newListByResourceGroupPager.get(req)
	if newListByResourceGroupPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HybridConnectivity/publicCloudConnectors`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		resp := p.srv.NewListByResourceGroupPager(resourceGroupNameParam, nil)
		newListByResourceGroupPager = &resp
		p.newListByResourceGroupPager.add(req, newListByResourceGroupPager)
		server.PagerResponderInjectNextLinks(newListByResourceGroupPager, req, func(page *armhybridconnectivity.PublicCloudConnectorsClientListByResourceGroupResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByResourceGroupPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		p.newListByResourceGroupPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByResourceGroupPager) {
		p.newListByResourceGroupPager.remove(req)
	}
	return resp, nil
}

func (p *PublicCloudConnectorsServerTransport) dispatchNewListBySubscriptionPager(req *http.Request) (*http.Response, error) {
	if p.srv.NewListBySubscriptionPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListBySubscriptionPager not implemented")}
	}
	newListBySubscriptionPager := p.newListBySubscriptionPager.get(req)
	if newListBySubscriptionPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HybridConnectivity/publicCloudConnectors`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resp := p.srv.NewListBySubscriptionPager(nil)
		newListBySubscriptionPager = &resp
		p.newListBySubscriptionPager.add(req, newListBySubscriptionPager)
		server.PagerResponderInjectNextLinks(newListBySubscriptionPager, req, func(page *armhybridconnectivity.PublicCloudConnectorsClientListBySubscriptionResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListBySubscriptionPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		p.newListBySubscriptionPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListBySubscriptionPager) {
		p.newListBySubscriptionPager.remove(req)
	}
	return resp, nil
}

func (p *PublicCloudConnectorsServerTransport) dispatchBeginTestPermissions(req *http.Request) (*http.Response, error) {
	if p.srv.BeginTestPermissions == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginTestPermissions not implemented")}
	}
	beginTestPermissions := p.beginTestPermissions.get(req)
	if beginTestPermissions == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HybridConnectivity/publicCloudConnectors/(?P<publicCloudConnector>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/testPermissions`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		publicCloudConnectorParam, err := url.PathUnescape(matches[regex.SubexpIndex("publicCloudConnector")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := p.srv.BeginTestPermissions(req.Context(), resourceGroupNameParam, publicCloudConnectorParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginTestPermissions = &respr
		p.beginTestPermissions.add(req, beginTestPermissions)
	}

	resp, err := server.PollerResponderNext(beginTestPermissions, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		p.beginTestPermissions.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginTestPermissions) {
		p.beginTestPermissions.remove(req)
	}

	return resp, nil
}

func (p *PublicCloudConnectorsServerTransport) dispatchUpdate(req *http.Request) (*http.Response, error) {
	if p.srv.Update == nil {
		return nil, &nonRetriableError{errors.New("fake for method Update not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HybridConnectivity/publicCloudConnectors/(?P<publicCloudConnector>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armhybridconnectivity.PublicCloudConnector](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	publicCloudConnectorParam, err := url.PathUnescape(matches[regex.SubexpIndex("publicCloudConnector")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := p.srv.Update(req.Context(), resourceGroupNameParam, publicCloudConnectorParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).PublicCloudConnector, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to PublicCloudConnectorsServerTransport
var publicCloudConnectorsServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
