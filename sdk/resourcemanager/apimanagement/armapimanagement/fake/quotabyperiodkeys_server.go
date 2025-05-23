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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement/v3"
	"net/http"
	"net/url"
	"regexp"
)

// QuotaByPeriodKeysServer is a fake server for instances of the armapimanagement.QuotaByPeriodKeysClient type.
type QuotaByPeriodKeysServer struct {
	// Get is the fake for method QuotaByPeriodKeysClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, serviceName string, quotaCounterKey string, quotaPeriodKey string, options *armapimanagement.QuotaByPeriodKeysClientGetOptions) (resp azfake.Responder[armapimanagement.QuotaByPeriodKeysClientGetResponse], errResp azfake.ErrorResponder)

	// Update is the fake for method QuotaByPeriodKeysClient.Update
	// HTTP status codes to indicate success: http.StatusOK
	Update func(ctx context.Context, resourceGroupName string, serviceName string, quotaCounterKey string, quotaPeriodKey string, parameters armapimanagement.QuotaCounterValueUpdateContract, options *armapimanagement.QuotaByPeriodKeysClientUpdateOptions) (resp azfake.Responder[armapimanagement.QuotaByPeriodKeysClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewQuotaByPeriodKeysServerTransport creates a new instance of QuotaByPeriodKeysServerTransport with the provided implementation.
// The returned QuotaByPeriodKeysServerTransport instance is connected to an instance of armapimanagement.QuotaByPeriodKeysClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewQuotaByPeriodKeysServerTransport(srv *QuotaByPeriodKeysServer) *QuotaByPeriodKeysServerTransport {
	return &QuotaByPeriodKeysServerTransport{srv: srv}
}

// QuotaByPeriodKeysServerTransport connects instances of armapimanagement.QuotaByPeriodKeysClient to instances of QuotaByPeriodKeysServer.
// Don't use this type directly, use NewQuotaByPeriodKeysServerTransport instead.
type QuotaByPeriodKeysServerTransport struct {
	srv *QuotaByPeriodKeysServer
}

// Do implements the policy.Transporter interface for QuotaByPeriodKeysServerTransport.
func (q *QuotaByPeriodKeysServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return q.dispatchToMethodFake(req, method)
}

func (q *QuotaByPeriodKeysServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if quotaByPeriodKeysServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = quotaByPeriodKeysServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "QuotaByPeriodKeysClient.Get":
				res.resp, res.err = q.dispatchGet(req)
			case "QuotaByPeriodKeysClient.Update":
				res.resp, res.err = q.dispatchUpdate(req)
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

func (q *QuotaByPeriodKeysServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if q.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/quotas/(?P<quotaCounterKey>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/periods/(?P<quotaPeriodKey>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
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
	quotaCounterKeyParam, err := url.PathUnescape(matches[regex.SubexpIndex("quotaCounterKey")])
	if err != nil {
		return nil, err
	}
	quotaPeriodKeyParam, err := url.PathUnescape(matches[regex.SubexpIndex("quotaPeriodKey")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := q.srv.Get(req.Context(), resourceGroupNameParam, serviceNameParam, quotaCounterKeyParam, quotaPeriodKeyParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).QuotaCounterContract, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (q *QuotaByPeriodKeysServerTransport) dispatchUpdate(req *http.Request) (*http.Response, error) {
	if q.srv.Update == nil {
		return nil, &nonRetriableError{errors.New("fake for method Update not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/quotas/(?P<quotaCounterKey>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/periods/(?P<quotaPeriodKey>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 5 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armapimanagement.QuotaCounterValueUpdateContract](req)
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
	quotaCounterKeyParam, err := url.PathUnescape(matches[regex.SubexpIndex("quotaCounterKey")])
	if err != nil {
		return nil, err
	}
	quotaPeriodKeyParam, err := url.PathUnescape(matches[regex.SubexpIndex("quotaPeriodKey")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := q.srv.Update(req.Context(), resourceGroupNameParam, serviceNameParam, quotaCounterKeyParam, quotaPeriodKeyParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).QuotaCounterContract, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to QuotaByPeriodKeysServerTransport
var quotaByPeriodKeysServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
