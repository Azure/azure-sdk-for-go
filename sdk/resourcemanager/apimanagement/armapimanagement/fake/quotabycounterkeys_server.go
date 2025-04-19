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

// QuotaByCounterKeysServer is a fake server for instances of the armapimanagement.QuotaByCounterKeysClient type.
type QuotaByCounterKeysServer struct {
	// ListByService is the fake for method QuotaByCounterKeysClient.ListByService
	// HTTP status codes to indicate success: http.StatusOK
	ListByService func(ctx context.Context, resourceGroupName string, serviceName string, quotaCounterKey string, options *armapimanagement.QuotaByCounterKeysClientListByServiceOptions) (resp azfake.Responder[armapimanagement.QuotaByCounterKeysClientListByServiceResponse], errResp azfake.ErrorResponder)

	// Update is the fake for method QuotaByCounterKeysClient.Update
	// HTTP status codes to indicate success: http.StatusOK
	Update func(ctx context.Context, resourceGroupName string, serviceName string, quotaCounterKey string, parameters armapimanagement.QuotaCounterValueUpdateContract, options *armapimanagement.QuotaByCounterKeysClientUpdateOptions) (resp azfake.Responder[armapimanagement.QuotaByCounterKeysClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewQuotaByCounterKeysServerTransport creates a new instance of QuotaByCounterKeysServerTransport with the provided implementation.
// The returned QuotaByCounterKeysServerTransport instance is connected to an instance of armapimanagement.QuotaByCounterKeysClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewQuotaByCounterKeysServerTransport(srv *QuotaByCounterKeysServer) *QuotaByCounterKeysServerTransport {
	return &QuotaByCounterKeysServerTransport{srv: srv}
}

// QuotaByCounterKeysServerTransport connects instances of armapimanagement.QuotaByCounterKeysClient to instances of QuotaByCounterKeysServer.
// Don't use this type directly, use NewQuotaByCounterKeysServerTransport instead.
type QuotaByCounterKeysServerTransport struct {
	srv *QuotaByCounterKeysServer
}

// Do implements the policy.Transporter interface for QuotaByCounterKeysServerTransport.
func (q *QuotaByCounterKeysServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return q.dispatchToMethodFake(req, method)
}

func (q *QuotaByCounterKeysServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if quotaByCounterKeysServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = quotaByCounterKeysServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "QuotaByCounterKeysClient.ListByService":
				res.resp, res.err = q.dispatchListByService(req)
			case "QuotaByCounterKeysClient.Update":
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

func (q *QuotaByCounterKeysServerTransport) dispatchListByService(req *http.Request) (*http.Response, error) {
	if q.srv.ListByService == nil {
		return nil, &nonRetriableError{errors.New("fake for method ListByService not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/quotas/(?P<quotaCounterKey>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
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
	respr, errRespr := q.srv.ListByService(req.Context(), resourceGroupNameParam, serviceNameParam, quotaCounterKeyParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).QuotaCounterCollection, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (q *QuotaByCounterKeysServerTransport) dispatchUpdate(req *http.Request) (*http.Response, error) {
	if q.srv.Update == nil {
		return nil, &nonRetriableError{errors.New("fake for method Update not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/quotas/(?P<quotaCounterKey>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
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
	respr, errRespr := q.srv.Update(req.Context(), resourceGroupNameParam, serviceNameParam, quotaCounterKeyParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).QuotaCounterCollection, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to QuotaByCounterKeysServerTransport
var quotaByCounterKeysServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
