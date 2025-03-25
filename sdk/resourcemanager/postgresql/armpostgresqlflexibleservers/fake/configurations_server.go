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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresqlflexibleservers/v4"
	"net/http"
	"net/url"
	"regexp"
)

// ConfigurationsServer is a fake server for instances of the armpostgresqlflexibleservers.ConfigurationsClient type.
type ConfigurationsServer struct {
	// Get is the fake for method ConfigurationsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, serverName string, configurationName string, options *armpostgresqlflexibleservers.ConfigurationsClientGetOptions) (resp azfake.Responder[armpostgresqlflexibleservers.ConfigurationsClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByServerPager is the fake for method ConfigurationsClient.NewListByServerPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByServerPager func(resourceGroupName string, serverName string, options *armpostgresqlflexibleservers.ConfigurationsClientListByServerOptions) (resp azfake.PagerResponder[armpostgresqlflexibleservers.ConfigurationsClientListByServerResponse])

	// BeginPut is the fake for method ConfigurationsClient.BeginPut
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated, http.StatusAccepted
	BeginPut func(ctx context.Context, resourceGroupName string, serverName string, configurationName string, parameters armpostgresqlflexibleservers.Configuration, options *armpostgresqlflexibleservers.ConfigurationsClientBeginPutOptions) (resp azfake.PollerResponder[armpostgresqlflexibleservers.ConfigurationsClientPutResponse], errResp azfake.ErrorResponder)

	// BeginUpdate is the fake for method ConfigurationsClient.BeginUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated, http.StatusAccepted
	BeginUpdate func(ctx context.Context, resourceGroupName string, serverName string, configurationName string, parameters armpostgresqlflexibleservers.ConfigurationForUpdate, options *armpostgresqlflexibleservers.ConfigurationsClientBeginUpdateOptions) (resp azfake.PollerResponder[armpostgresqlflexibleservers.ConfigurationsClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewConfigurationsServerTransport creates a new instance of ConfigurationsServerTransport with the provided implementation.
// The returned ConfigurationsServerTransport instance is connected to an instance of armpostgresqlflexibleservers.ConfigurationsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewConfigurationsServerTransport(srv *ConfigurationsServer) *ConfigurationsServerTransport {
	return &ConfigurationsServerTransport{
		srv:                  srv,
		newListByServerPager: newTracker[azfake.PagerResponder[armpostgresqlflexibleservers.ConfigurationsClientListByServerResponse]](),
		beginPut:             newTracker[azfake.PollerResponder[armpostgresqlflexibleservers.ConfigurationsClientPutResponse]](),
		beginUpdate:          newTracker[azfake.PollerResponder[armpostgresqlflexibleservers.ConfigurationsClientUpdateResponse]](),
	}
}

// ConfigurationsServerTransport connects instances of armpostgresqlflexibleservers.ConfigurationsClient to instances of ConfigurationsServer.
// Don't use this type directly, use NewConfigurationsServerTransport instead.
type ConfigurationsServerTransport struct {
	srv                  *ConfigurationsServer
	newListByServerPager *tracker[azfake.PagerResponder[armpostgresqlflexibleservers.ConfigurationsClientListByServerResponse]]
	beginPut             *tracker[azfake.PollerResponder[armpostgresqlflexibleservers.ConfigurationsClientPutResponse]]
	beginUpdate          *tracker[azfake.PollerResponder[armpostgresqlflexibleservers.ConfigurationsClientUpdateResponse]]
}

// Do implements the policy.Transporter interface for ConfigurationsServerTransport.
func (c *ConfigurationsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return c.dispatchToMethodFake(req, method)
}

func (c *ConfigurationsServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if configurationsServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = configurationsServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "ConfigurationsClient.Get":
				res.resp, res.err = c.dispatchGet(req)
			case "ConfigurationsClient.NewListByServerPager":
				res.resp, res.err = c.dispatchNewListByServerPager(req)
			case "ConfigurationsClient.BeginPut":
				res.resp, res.err = c.dispatchBeginPut(req)
			case "ConfigurationsClient.BeginUpdate":
				res.resp, res.err = c.dispatchBeginUpdate(req)
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

func (c *ConfigurationsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if c.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DBforPostgreSQL/flexibleServers/(?P<serverName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/configurations/(?P<configurationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	serverNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serverName")])
	if err != nil {
		return nil, err
	}
	configurationNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("configurationName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := c.srv.Get(req.Context(), resourceGroupNameParam, serverNameParam, configurationNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).Configuration, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ConfigurationsServerTransport) dispatchNewListByServerPager(req *http.Request) (*http.Response, error) {
	if c.srv.NewListByServerPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByServerPager not implemented")}
	}
	newListByServerPager := c.newListByServerPager.get(req)
	if newListByServerPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DBforPostgreSQL/flexibleServers/(?P<serverName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/configurations`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		serverNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serverName")])
		if err != nil {
			return nil, err
		}
		resp := c.srv.NewListByServerPager(resourceGroupNameParam, serverNameParam, nil)
		newListByServerPager = &resp
		c.newListByServerPager.add(req, newListByServerPager)
		server.PagerResponderInjectNextLinks(newListByServerPager, req, func(page *armpostgresqlflexibleservers.ConfigurationsClientListByServerResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByServerPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		c.newListByServerPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByServerPager) {
		c.newListByServerPager.remove(req)
	}
	return resp, nil
}

func (c *ConfigurationsServerTransport) dispatchBeginPut(req *http.Request) (*http.Response, error) {
	if c.srv.BeginPut == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginPut not implemented")}
	}
	beginPut := c.beginPut.get(req)
	if beginPut == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DBforPostgreSQL/flexibleServers/(?P<serverName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/configurations/(?P<configurationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armpostgresqlflexibleservers.Configuration](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		serverNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serverName")])
		if err != nil {
			return nil, err
		}
		configurationNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("configurationName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := c.srv.BeginPut(req.Context(), resourceGroupNameParam, serverNameParam, configurationNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginPut = &respr
		c.beginPut.add(req, beginPut)
	}

	resp, err := server.PollerResponderNext(beginPut, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated, http.StatusAccepted}, resp.StatusCode) {
		c.beginPut.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginPut) {
		c.beginPut.remove(req)
	}

	return resp, nil
}

func (c *ConfigurationsServerTransport) dispatchBeginUpdate(req *http.Request) (*http.Response, error) {
	if c.srv.BeginUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdate not implemented")}
	}
	beginUpdate := c.beginUpdate.get(req)
	if beginUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DBforPostgreSQL/flexibleServers/(?P<serverName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/configurations/(?P<configurationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armpostgresqlflexibleservers.ConfigurationForUpdate](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		serverNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serverName")])
		if err != nil {
			return nil, err
		}
		configurationNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("configurationName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := c.srv.BeginUpdate(req.Context(), resourceGroupNameParam, serverNameParam, configurationNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginUpdate = &respr
		c.beginUpdate.add(req, beginUpdate)
	}

	resp, err := server.PollerResponderNext(beginUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated, http.StatusAccepted}, resp.StatusCode) {
		c.beginUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdate) {
		c.beginUpdate.remove(req)
	}

	return resp, nil
}

// set this to conditionally intercept incoming requests to ConfigurationsServerTransport
var configurationsServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
