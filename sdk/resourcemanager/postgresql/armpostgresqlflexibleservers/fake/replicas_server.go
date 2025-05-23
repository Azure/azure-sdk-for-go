// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package fake

import (
	"errors"
	"fmt"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake/server"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresqlflexibleservers/v5"
	"net/http"
	"net/url"
	"regexp"
)

// ReplicasServer is a fake server for instances of the armpostgresqlflexibleservers.ReplicasClient type.
type ReplicasServer struct {
	// NewListByServerPager is the fake for method ReplicasClient.NewListByServerPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByServerPager func(resourceGroupName string, serverName string, options *armpostgresqlflexibleservers.ReplicasClientListByServerOptions) (resp azfake.PagerResponder[armpostgresqlflexibleservers.ReplicasClientListByServerResponse])
}

// NewReplicasServerTransport creates a new instance of ReplicasServerTransport with the provided implementation.
// The returned ReplicasServerTransport instance is connected to an instance of armpostgresqlflexibleservers.ReplicasClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewReplicasServerTransport(srv *ReplicasServer) *ReplicasServerTransport {
	return &ReplicasServerTransport{
		srv:                  srv,
		newListByServerPager: newTracker[azfake.PagerResponder[armpostgresqlflexibleservers.ReplicasClientListByServerResponse]](),
	}
}

// ReplicasServerTransport connects instances of armpostgresqlflexibleservers.ReplicasClient to instances of ReplicasServer.
// Don't use this type directly, use NewReplicasServerTransport instead.
type ReplicasServerTransport struct {
	srv                  *ReplicasServer
	newListByServerPager *tracker[azfake.PagerResponder[armpostgresqlflexibleservers.ReplicasClientListByServerResponse]]
}

// Do implements the policy.Transporter interface for ReplicasServerTransport.
func (r *ReplicasServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return r.dispatchToMethodFake(req, method)
}

func (r *ReplicasServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if replicasServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = replicasServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "ReplicasClient.NewListByServerPager":
				res.resp, res.err = r.dispatchNewListByServerPager(req)
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

func (r *ReplicasServerTransport) dispatchNewListByServerPager(req *http.Request) (*http.Response, error) {
	if r.srv.NewListByServerPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByServerPager not implemented")}
	}
	newListByServerPager := r.newListByServerPager.get(req)
	if newListByServerPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DBforPostgreSQL/flexibleServers/(?P<serverName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/replicas`
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
		resp := r.srv.NewListByServerPager(resourceGroupNameParam, serverNameParam, nil)
		newListByServerPager = &resp
		r.newListByServerPager.add(req, newListByServerPager)
	}
	resp, err := server.PagerResponderNext(newListByServerPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		r.newListByServerPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByServerPager) {
		r.newListByServerPager.remove(req)
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to ReplicasServerTransport
var replicasServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
