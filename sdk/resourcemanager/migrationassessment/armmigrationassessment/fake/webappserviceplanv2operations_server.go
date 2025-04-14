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

// WebAppServicePlanV2OperationsServer is a fake server for instances of the armmigrationassessment.WebAppServicePlanV2OperationsClient type.
type WebAppServicePlanV2OperationsServer struct {
	// Get is the fake for method WebAppServicePlanV2OperationsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, projectName string, groupName string, assessmentName string, webAppServicePlanName string, options *armmigrationassessment.WebAppServicePlanV2OperationsClientGetOptions) (resp azfake.Responder[armmigrationassessment.WebAppServicePlanV2OperationsClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByWebAppAssessmentV2Pager is the fake for method WebAppServicePlanV2OperationsClient.NewListByWebAppAssessmentV2Pager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByWebAppAssessmentV2Pager func(resourceGroupName string, projectName string, groupName string, assessmentName string, options *armmigrationassessment.WebAppServicePlanV2OperationsClientListByWebAppAssessmentV2Options) (resp azfake.PagerResponder[armmigrationassessment.WebAppServicePlanV2OperationsClientListByWebAppAssessmentV2Response])
}

// NewWebAppServicePlanV2OperationsServerTransport creates a new instance of WebAppServicePlanV2OperationsServerTransport with the provided implementation.
// The returned WebAppServicePlanV2OperationsServerTransport instance is connected to an instance of armmigrationassessment.WebAppServicePlanV2OperationsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewWebAppServicePlanV2OperationsServerTransport(srv *WebAppServicePlanV2OperationsServer) *WebAppServicePlanV2OperationsServerTransport {
	return &WebAppServicePlanV2OperationsServerTransport{
		srv:                              srv,
		newListByWebAppAssessmentV2Pager: newTracker[azfake.PagerResponder[armmigrationassessment.WebAppServicePlanV2OperationsClientListByWebAppAssessmentV2Response]](),
	}
}

// WebAppServicePlanV2OperationsServerTransport connects instances of armmigrationassessment.WebAppServicePlanV2OperationsClient to instances of WebAppServicePlanV2OperationsServer.
// Don't use this type directly, use NewWebAppServicePlanV2OperationsServerTransport instead.
type WebAppServicePlanV2OperationsServerTransport struct {
	srv                              *WebAppServicePlanV2OperationsServer
	newListByWebAppAssessmentV2Pager *tracker[azfake.PagerResponder[armmigrationassessment.WebAppServicePlanV2OperationsClientListByWebAppAssessmentV2Response]]
}

// Do implements the policy.Transporter interface for WebAppServicePlanV2OperationsServerTransport.
func (w *WebAppServicePlanV2OperationsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return w.dispatchToMethodFake(req, method)
}

func (w *WebAppServicePlanV2OperationsServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if webAppServicePlanV2OperationsServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = webAppServicePlanV2OperationsServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "WebAppServicePlanV2OperationsClient.Get":
				res.resp, res.err = w.dispatchGet(req)
			case "WebAppServicePlanV2OperationsClient.NewListByWebAppAssessmentV2Pager":
				res.resp, res.err = w.dispatchNewListByWebAppAssessmentV2Pager(req)
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

func (w *WebAppServicePlanV2OperationsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if w.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Migrate/assessmentProjects/(?P<projectName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/groups/(?P<groupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/webAppAssessments/(?P<assessmentName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/webAppServicePlans/(?P<webAppServicePlanName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 6 {
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
	groupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("groupName")])
	if err != nil {
		return nil, err
	}
	assessmentNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("assessmentName")])
	if err != nil {
		return nil, err
	}
	webAppServicePlanNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("webAppServicePlanName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := w.srv.Get(req.Context(), resourceGroupNameParam, projectNameParam, groupNameParam, assessmentNameParam, webAppServicePlanNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).WebAppServicePlanV2, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (w *WebAppServicePlanV2OperationsServerTransport) dispatchNewListByWebAppAssessmentV2Pager(req *http.Request) (*http.Response, error) {
	if w.srv.NewListByWebAppAssessmentV2Pager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByWebAppAssessmentV2Pager not implemented")}
	}
	newListByWebAppAssessmentV2Pager := w.newListByWebAppAssessmentV2Pager.get(req)
	if newListByWebAppAssessmentV2Pager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Migrate/assessmentProjects/(?P<projectName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/groups/(?P<groupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/webAppAssessments/(?P<assessmentName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/webAppServicePlans`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 5 {
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
		groupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("groupName")])
		if err != nil {
			return nil, err
		}
		assessmentNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("assessmentName")])
		if err != nil {
			return nil, err
		}
		var options *armmigrationassessment.WebAppServicePlanV2OperationsClientListByWebAppAssessmentV2Options
		if filterParam != nil || pageSizeParam != nil || continuationTokenParam != nil || totalRecordCountParam != nil {
			options = &armmigrationassessment.WebAppServicePlanV2OperationsClientListByWebAppAssessmentV2Options{
				Filter:            filterParam,
				PageSize:          pageSizeParam,
				ContinuationToken: continuationTokenParam,
				TotalRecordCount:  totalRecordCountParam,
			}
		}
		resp := w.srv.NewListByWebAppAssessmentV2Pager(resourceGroupNameParam, projectNameParam, groupNameParam, assessmentNameParam, options)
		newListByWebAppAssessmentV2Pager = &resp
		w.newListByWebAppAssessmentV2Pager.add(req, newListByWebAppAssessmentV2Pager)
		server.PagerResponderInjectNextLinks(newListByWebAppAssessmentV2Pager, req, func(page *armmigrationassessment.WebAppServicePlanV2OperationsClientListByWebAppAssessmentV2Response, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByWebAppAssessmentV2Pager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		w.newListByWebAppAssessmentV2Pager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByWebAppAssessmentV2Pager) {
		w.newListByWebAppAssessmentV2Pager.remove(req)
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to WebAppServicePlanV2OperationsServerTransport
var webAppServicePlanV2OperationsServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
