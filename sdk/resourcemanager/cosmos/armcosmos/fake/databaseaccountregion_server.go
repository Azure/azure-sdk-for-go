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

// DatabaseAccountRegionServer is a fake server for instances of the armcosmos.DatabaseAccountRegionClient type.
type DatabaseAccountRegionServer struct {
	// NewListMetricsPager is the fake for method DatabaseAccountRegionClient.NewListMetricsPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListMetricsPager func(resourceGroupName string, accountName string, region string, filter string, options *armcosmos.DatabaseAccountRegionClientListMetricsOptions) (resp azfake.PagerResponder[armcosmos.DatabaseAccountRegionClientListMetricsResponse])
}

// NewDatabaseAccountRegionServerTransport creates a new instance of DatabaseAccountRegionServerTransport with the provided implementation.
// The returned DatabaseAccountRegionServerTransport instance is connected to an instance of armcosmos.DatabaseAccountRegionClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewDatabaseAccountRegionServerTransport(srv *DatabaseAccountRegionServer) *DatabaseAccountRegionServerTransport {
	return &DatabaseAccountRegionServerTransport{
		srv:                 srv,
		newListMetricsPager: newTracker[azfake.PagerResponder[armcosmos.DatabaseAccountRegionClientListMetricsResponse]](),
	}
}

// DatabaseAccountRegionServerTransport connects instances of armcosmos.DatabaseAccountRegionClient to instances of DatabaseAccountRegionServer.
// Don't use this type directly, use NewDatabaseAccountRegionServerTransport instead.
type DatabaseAccountRegionServerTransport struct {
	srv                 *DatabaseAccountRegionServer
	newListMetricsPager *tracker[azfake.PagerResponder[armcosmos.DatabaseAccountRegionClientListMetricsResponse]]
}

// Do implements the policy.Transporter interface for DatabaseAccountRegionServerTransport.
func (d *DatabaseAccountRegionServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return d.dispatchToMethodFake(req, method)
}

func (d *DatabaseAccountRegionServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if databaseAccountRegionServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = databaseAccountRegionServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "DatabaseAccountRegionClient.NewListMetricsPager":
				res.resp, res.err = d.dispatchNewListMetricsPager(req)
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

func (d *DatabaseAccountRegionServerTransport) dispatchNewListMetricsPager(req *http.Request) (*http.Response, error) {
	if d.srv.NewListMetricsPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListMetricsPager not implemented")}
	}
	newListMetricsPager := d.newListMetricsPager.get(req)
	if newListMetricsPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DocumentDB/databaseAccounts/(?P<accountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/region/(?P<region>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/metrics`
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
		accountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("accountName")])
		if err != nil {
			return nil, err
		}
		regionParam, err := url.PathUnescape(matches[regex.SubexpIndex("region")])
		if err != nil {
			return nil, err
		}
		filterParam, err := url.QueryUnescape(qp.Get("$filter"))
		if err != nil {
			return nil, err
		}
		resp := d.srv.NewListMetricsPager(resourceGroupNameParam, accountNameParam, regionParam, filterParam, nil)
		newListMetricsPager = &resp
		d.newListMetricsPager.add(req, newListMetricsPager)
	}
	resp, err := server.PagerResponderNext(newListMetricsPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		d.newListMetricsPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListMetricsPager) {
		d.newListMetricsPager.remove(req)
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to DatabaseAccountRegionServerTransport
var databaseAccountRegionServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
