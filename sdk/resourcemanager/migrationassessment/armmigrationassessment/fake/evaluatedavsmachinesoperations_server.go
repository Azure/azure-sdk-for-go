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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/migrationassessment/armmigrationassessment"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

// EvaluatedAvsMachinesOperationsServer is a fake server for instances of the armmigrationassessment.EvaluatedAvsMachinesOperationsClient type.
type EvaluatedAvsMachinesOperationsServer struct {
	// Get is the fake for method EvaluatedAvsMachinesOperationsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, projectName string, businessCaseName string, evaluatedAvsMachineName string, options *armmigrationassessment.EvaluatedAvsMachinesOperationsClientGetOptions) (resp azfake.Responder[armmigrationassessment.EvaluatedAvsMachinesOperationsClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByBusinessCasePager is the fake for method EvaluatedAvsMachinesOperationsClient.NewListByBusinessCasePager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByBusinessCasePager func(resourceGroupName string, projectName string, businessCaseName string, options *armmigrationassessment.EvaluatedAvsMachinesOperationsClientListByBusinessCaseOptions) (resp azfake.PagerResponder[armmigrationassessment.EvaluatedAvsMachinesOperationsClientListByBusinessCaseResponse])
}

// NewEvaluatedAvsMachinesOperationsServerTransport creates a new instance of EvaluatedAvsMachinesOperationsServerTransport with the provided implementation.
// The returned EvaluatedAvsMachinesOperationsServerTransport instance is connected to an instance of armmigrationassessment.EvaluatedAvsMachinesOperationsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewEvaluatedAvsMachinesOperationsServerTransport(srv *EvaluatedAvsMachinesOperationsServer) *EvaluatedAvsMachinesOperationsServerTransport {
	return &EvaluatedAvsMachinesOperationsServerTransport{
		srv:                        srv,
		newListByBusinessCasePager: newTracker[azfake.PagerResponder[armmigrationassessment.EvaluatedAvsMachinesOperationsClientListByBusinessCaseResponse]](),
	}
}

// EvaluatedAvsMachinesOperationsServerTransport connects instances of armmigrationassessment.EvaluatedAvsMachinesOperationsClient to instances of EvaluatedAvsMachinesOperationsServer.
// Don't use this type directly, use NewEvaluatedAvsMachinesOperationsServerTransport instead.
type EvaluatedAvsMachinesOperationsServerTransport struct {
	srv                        *EvaluatedAvsMachinesOperationsServer
	newListByBusinessCasePager *tracker[azfake.PagerResponder[armmigrationassessment.EvaluatedAvsMachinesOperationsClientListByBusinessCaseResponse]]
}

// Do implements the policy.Transporter interface for EvaluatedAvsMachinesOperationsServerTransport.
func (e *EvaluatedAvsMachinesOperationsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return e.dispatchToMethodFake(req, method)
}

func (e *EvaluatedAvsMachinesOperationsServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if evaluatedAvsMachinesOperationsServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = evaluatedAvsMachinesOperationsServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "EvaluatedAvsMachinesOperationsClient.Get":
				res.resp, res.err = e.dispatchGet(req)
			case "EvaluatedAvsMachinesOperationsClient.NewListByBusinessCasePager":
				res.resp, res.err = e.dispatchNewListByBusinessCasePager(req)
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

func (e *EvaluatedAvsMachinesOperationsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if e.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Migrate/assessmentProjects/(?P<projectName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/businessCases/(?P<businessCaseName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/evaluatedAvsMachines/(?P<evaluatedAvsMachineName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 5 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	projectNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("projectName")])
	if err != nil {
		return nil, err
	}
	businessCaseNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("businessCaseName")])
	if err != nil {
		return nil, err
	}
	evaluatedAvsMachineNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("evaluatedAvsMachineName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := e.srv.Get(req.Context(), resourceGroupNameParam, projectNameParam, businessCaseNameParam, evaluatedAvsMachineNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).EvaluatedAvsMachine, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (e *EvaluatedAvsMachinesOperationsServerTransport) dispatchNewListByBusinessCasePager(req *http.Request) (*http.Response, error) {
	if e.srv.NewListByBusinessCasePager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByBusinessCasePager not implemented")}
	}
	newListByBusinessCasePager := e.newListByBusinessCasePager.get(req)
	if newListByBusinessCasePager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Migrate/assessmentProjects/(?P<projectName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/businessCases/(?P<businessCaseName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/evaluatedAvsMachines`
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
		filterUnescaped, err := url.QueryUnescape(qp.Get("$filter"))
		if err != nil {
			return nil, err
		}
		filterParam := getOptional(filterUnescaped)
		pageSizeUnescaped, err := url.QueryUnescape(qp.Get("pageSize"))
		if err != nil {
			return nil, err
		}
		pageSizeParam, err := parseOptional(pageSizeUnescaped, func(v string) (int32, error) {
			p, parseErr := strconv.ParseInt(v, 10, 32)
			if parseErr != nil {
				return 0, parseErr
			}
			return int32(p), nil
		})
		if err != nil {
			return nil, err
		}
		continuationTokenUnescaped, err := url.QueryUnescape(qp.Get("continuationToken"))
		if err != nil {
			return nil, err
		}
		continuationTokenParam := getOptional(continuationTokenUnescaped)
		totalRecordCountUnescaped, err := url.QueryUnescape(qp.Get("totalRecordCount"))
		if err != nil {
			return nil, err
		}
		totalRecordCountParam, err := parseOptional(totalRecordCountUnescaped, func(v string) (int32, error) {
			p, parseErr := strconv.ParseInt(v, 10, 32)
			if parseErr != nil {
				return 0, parseErr
			}
			return int32(p), nil
		})
		if err != nil {
			return nil, err
		}
		projectNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("projectName")])
		if err != nil {
			return nil, err
		}
		businessCaseNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("businessCaseName")])
		if err != nil {
			return nil, err
		}
		var options *armmigrationassessment.EvaluatedAvsMachinesOperationsClientListByBusinessCaseOptions
		if filterParam != nil || pageSizeParam != nil || continuationTokenParam != nil || totalRecordCountParam != nil {
			options = &armmigrationassessment.EvaluatedAvsMachinesOperationsClientListByBusinessCaseOptions{
				Filter:            filterParam,
				PageSize:          pageSizeParam,
				ContinuationToken: continuationTokenParam,
				TotalRecordCount:  totalRecordCountParam,
			}
		}
		resp := e.srv.NewListByBusinessCasePager(resourceGroupNameParam, projectNameParam, businessCaseNameParam, options)
		newListByBusinessCasePager = &resp
		e.newListByBusinessCasePager.add(req, newListByBusinessCasePager)
		server.PagerResponderInjectNextLinks(newListByBusinessCasePager, req, func(page *armmigrationassessment.EvaluatedAvsMachinesOperationsClientListByBusinessCaseResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByBusinessCasePager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		e.newListByBusinessCasePager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByBusinessCasePager) {
		e.newListByBusinessCasePager.remove(req)
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to EvaluatedAvsMachinesOperationsServerTransport
var evaluatedAvsMachinesOperationsServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
