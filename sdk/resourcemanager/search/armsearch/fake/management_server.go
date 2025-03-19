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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/search/armsearch"
	"net/http"
	"net/url"
	"regexp"
)

// ManagementServer is a fake server for instances of the armsearch.ManagementClient type.
type ManagementServer struct {
	// UsageBySubscriptionSKU is the fake for method ManagementClient.UsageBySubscriptionSKU
	// HTTP status codes to indicate success: http.StatusOK
	UsageBySubscriptionSKU func(ctx context.Context, location string, skuName string, searchManagementRequestOptions *armsearch.SearchManagementRequestOptions, options *armsearch.ManagementClientUsageBySubscriptionSKUOptions) (resp azfake.Responder[armsearch.ManagementClientUsageBySubscriptionSKUResponse], errResp azfake.ErrorResponder)
}

// NewManagementServerTransport creates a new instance of ManagementServerTransport with the provided implementation.
// The returned ManagementServerTransport instance is connected to an instance of armsearch.ManagementClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewManagementServerTransport(srv *ManagementServer) *ManagementServerTransport {
	return &ManagementServerTransport{srv: srv}
}

// ManagementServerTransport connects instances of armsearch.ManagementClient to instances of ManagementServer.
// Don't use this type directly, use NewManagementServerTransport instead.
type ManagementServerTransport struct {
	srv *ManagementServer
}

// Do implements the policy.Transporter interface for ManagementServerTransport.
func (m *ManagementServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return m.dispatchToMethodFake(req, method)
}

func (m *ManagementServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if managementServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = managementServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "ManagementClient.UsageBySubscriptionSKU":
				res.resp, res.err = m.dispatchUsageBySubscriptionSKU(req)
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

func (m *ManagementServerTransport) dispatchUsageBySubscriptionSKU(req *http.Request) (*http.Response, error) {
	if m.srv.UsageBySubscriptionSKU == nil {
		return nil, &nonRetriableError{errors.New("fake for method UsageBySubscriptionSKU not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Search/locations/(?P<location>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/usages/(?P<skuName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	clientRequestIDParam := getOptional(getHeaderValue(req.Header, "x-ms-client-request-id"))
	locationParam, err := url.PathUnescape(matches[regex.SubexpIndex("location")])
	if err != nil {
		return nil, err
	}
	skuNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("skuName")])
	if err != nil {
		return nil, err
	}
	var searchManagementRequestOptions *armsearch.SearchManagementRequestOptions
	if clientRequestIDParam != nil {
		searchManagementRequestOptions = &armsearch.SearchManagementRequestOptions{
			ClientRequestID: clientRequestIDParam,
		}
	}
	respr, errRespr := m.srv.UsageBySubscriptionSKU(req.Context(), locationParam, skuNameParam, searchManagementRequestOptions, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).QuotaUsageResult, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to ManagementServerTransport
var managementServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
