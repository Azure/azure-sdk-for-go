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

// ServerThreatProtectionSettingsServer is a fake server for instances of the armpostgresqlflexibleservers.ServerThreatProtectionSettingsClient type.
type ServerThreatProtectionSettingsServer struct {
	// BeginCreateOrUpdate is the fake for method ServerThreatProtectionSettingsClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated, http.StatusAccepted
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, serverName string, threatProtectionName armpostgresqlflexibleservers.ThreatProtectionName, parameters armpostgresqlflexibleservers.ServerThreatProtectionSettingsModel, options *armpostgresqlflexibleservers.ServerThreatProtectionSettingsClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armpostgresqlflexibleservers.ServerThreatProtectionSettingsClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method ServerThreatProtectionSettingsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, serverName string, threatProtectionName armpostgresqlflexibleservers.ThreatProtectionName, options *armpostgresqlflexibleservers.ServerThreatProtectionSettingsClientGetOptions) (resp azfake.Responder[armpostgresqlflexibleservers.ServerThreatProtectionSettingsClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByServerPager is the fake for method ServerThreatProtectionSettingsClient.NewListByServerPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByServerPager func(resourceGroupName string, serverName string, options *armpostgresqlflexibleservers.ServerThreatProtectionSettingsClientListByServerOptions) (resp azfake.PagerResponder[armpostgresqlflexibleservers.ServerThreatProtectionSettingsClientListByServerResponse])
}

// NewServerThreatProtectionSettingsServerTransport creates a new instance of ServerThreatProtectionSettingsServerTransport with the provided implementation.
// The returned ServerThreatProtectionSettingsServerTransport instance is connected to an instance of armpostgresqlflexibleservers.ServerThreatProtectionSettingsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewServerThreatProtectionSettingsServerTransport(srv *ServerThreatProtectionSettingsServer) *ServerThreatProtectionSettingsServerTransport {
	return &ServerThreatProtectionSettingsServerTransport{
		srv:                  srv,
		beginCreateOrUpdate:  newTracker[azfake.PollerResponder[armpostgresqlflexibleservers.ServerThreatProtectionSettingsClientCreateOrUpdateResponse]](),
		newListByServerPager: newTracker[azfake.PagerResponder[armpostgresqlflexibleservers.ServerThreatProtectionSettingsClientListByServerResponse]](),
	}
}

// ServerThreatProtectionSettingsServerTransport connects instances of armpostgresqlflexibleservers.ServerThreatProtectionSettingsClient to instances of ServerThreatProtectionSettingsServer.
// Don't use this type directly, use NewServerThreatProtectionSettingsServerTransport instead.
type ServerThreatProtectionSettingsServerTransport struct {
	srv                  *ServerThreatProtectionSettingsServer
	beginCreateOrUpdate  *tracker[azfake.PollerResponder[armpostgresqlflexibleservers.ServerThreatProtectionSettingsClientCreateOrUpdateResponse]]
	newListByServerPager *tracker[azfake.PagerResponder[armpostgresqlflexibleservers.ServerThreatProtectionSettingsClientListByServerResponse]]
}

// Do implements the policy.Transporter interface for ServerThreatProtectionSettingsServerTransport.
func (s *ServerThreatProtectionSettingsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return s.dispatchToMethodFake(req, method)
}

func (s *ServerThreatProtectionSettingsServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if serverThreatProtectionSettingsServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = serverThreatProtectionSettingsServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "ServerThreatProtectionSettingsClient.BeginCreateOrUpdate":
				res.resp, res.err = s.dispatchBeginCreateOrUpdate(req)
			case "ServerThreatProtectionSettingsClient.Get":
				res.resp, res.err = s.dispatchGet(req)
			case "ServerThreatProtectionSettingsClient.NewListByServerPager":
				res.resp, res.err = s.dispatchNewListByServerPager(req)
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

func (s *ServerThreatProtectionSettingsServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if s.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := s.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DBforPostgreSQL/flexibleServers/(?P<serverName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/advancedThreatProtectionSettings/(?P<threatProtectionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armpostgresqlflexibleservers.ServerThreatProtectionSettingsModel](req)
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
		threatProtectionNameParam, err := parseWithCast(matches[regex.SubexpIndex("threatProtectionName")], func(v string) (armpostgresqlflexibleservers.ThreatProtectionName, error) {
			p, unescapeErr := url.PathUnescape(v)
			if unescapeErr != nil {
				return "", unescapeErr
			}
			return armpostgresqlflexibleservers.ThreatProtectionName(p), nil
		})
		if err != nil {
			return nil, err
		}
		respr, errRespr := s.srv.BeginCreateOrUpdate(req.Context(), resourceGroupNameParam, serverNameParam, threatProtectionNameParam, body, nil)
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

	if !contains([]int{http.StatusOK, http.StatusCreated, http.StatusAccepted}, resp.StatusCode) {
		s.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		s.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (s *ServerThreatProtectionSettingsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if s.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DBforPostgreSQL/flexibleServers/(?P<serverName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/advancedThreatProtectionSettings/(?P<threatProtectionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
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
	threatProtectionNameParam, err := parseWithCast(matches[regex.SubexpIndex("threatProtectionName")], func(v string) (armpostgresqlflexibleservers.ThreatProtectionName, error) {
		p, unescapeErr := url.PathUnescape(v)
		if unescapeErr != nil {
			return "", unescapeErr
		}
		return armpostgresqlflexibleservers.ThreatProtectionName(p), nil
	})
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.Get(req.Context(), resourceGroupNameParam, serverNameParam, threatProtectionNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ServerThreatProtectionSettingsModel, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *ServerThreatProtectionSettingsServerTransport) dispatchNewListByServerPager(req *http.Request) (*http.Response, error) {
	if s.srv.NewListByServerPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByServerPager not implemented")}
	}
	newListByServerPager := s.newListByServerPager.get(req)
	if newListByServerPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DBforPostgreSQL/flexibleServers/(?P<serverName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/advancedThreatProtectionSettings`
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
		resp := s.srv.NewListByServerPager(resourceGroupNameParam, serverNameParam, nil)
		newListByServerPager = &resp
		s.newListByServerPager.add(req, newListByServerPager)
		server.PagerResponderInjectNextLinks(newListByServerPager, req, func(page *armpostgresqlflexibleservers.ServerThreatProtectionSettingsClientListByServerResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByServerPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		s.newListByServerPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByServerPager) {
		s.newListByServerPager.remove(req)
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to ServerThreatProtectionSettingsServerTransport
var serverThreatProtectionSettingsServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
