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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos/v3"
	"net/http"
	"net/url"
	"regexp"
)

// RestorableSQLDatabasesServer is a fake server for instances of the armcosmos.RestorableSQLDatabasesClient type.
type RestorableSQLDatabasesServer struct {
	// NewListPager is the fake for method RestorableSQLDatabasesClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(location string, instanceID string, options *armcosmos.RestorableSQLDatabasesClientListOptions) (resp azfake.PagerResponder[armcosmos.RestorableSQLDatabasesClientListResponse])
}

// NewRestorableSQLDatabasesServerTransport creates a new instance of RestorableSQLDatabasesServerTransport with the provided implementation.
// The returned RestorableSQLDatabasesServerTransport instance is connected to an instance of armcosmos.RestorableSQLDatabasesClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewRestorableSQLDatabasesServerTransport(srv *RestorableSQLDatabasesServer) *RestorableSQLDatabasesServerTransport {
	return &RestorableSQLDatabasesServerTransport{
		srv:          srv,
		newListPager: newTracker[azfake.PagerResponder[armcosmos.RestorableSQLDatabasesClientListResponse]](),
	}
}

// RestorableSQLDatabasesServerTransport connects instances of armcosmos.RestorableSQLDatabasesClient to instances of RestorableSQLDatabasesServer.
// Don't use this type directly, use NewRestorableSQLDatabasesServerTransport instead.
type RestorableSQLDatabasesServerTransport struct {
	srv          *RestorableSQLDatabasesServer
	newListPager *tracker[azfake.PagerResponder[armcosmos.RestorableSQLDatabasesClientListResponse]]
}

// Do implements the policy.Transporter interface for RestorableSQLDatabasesServerTransport.
func (r *RestorableSQLDatabasesServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return r.dispatchToMethodFake(req, method)
}

func (r *RestorableSQLDatabasesServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if restorableSqlDatabasesServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = restorableSqlDatabasesServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "RestorableSQLDatabasesClient.NewListPager":
				res.resp, res.err = r.dispatchNewListPager(req)
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

func (r *RestorableSQLDatabasesServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if r.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := r.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DocumentDB/locations/(?P<location>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/restorableDatabaseAccounts/(?P<instanceId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/restorableSqlDatabases`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		locationParam, err := url.PathUnescape(matches[regex.SubexpIndex("location")])
		if err != nil {
			return nil, err
		}
		instanceIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("instanceId")])
		if err != nil {
			return nil, err
		}
		resp := r.srv.NewListPager(locationParam, instanceIDParam, nil)
		newListPager = &resp
		r.newListPager.add(req, newListPager)
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		r.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		r.newListPager.remove(req)
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to RestorableSQLDatabasesServerTransport
var restorableSqlDatabasesServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
