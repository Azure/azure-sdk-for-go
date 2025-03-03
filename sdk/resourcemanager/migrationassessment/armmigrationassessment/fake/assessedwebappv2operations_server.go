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

// AssessedWebAppV2OperationsServer is a fake server for instances of the armmigrationassessment.AssessedWebAppV2OperationsClient type.
type AssessedWebAppV2OperationsServer struct {
	// Get is the fake for method AssessedWebAppV2OperationsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, projectName string, groupName string, assessmentName string, assessedWebAppName string, options *armmigrationassessment.AssessedWebAppV2OperationsClientGetOptions) (resp azfake.Responder[armmigrationassessment.AssessedWebAppV2OperationsClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByWebAppAssessmentV2Pager is the fake for method AssessedWebAppV2OperationsClient.NewListByWebAppAssessmentV2Pager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByWebAppAssessmentV2Pager func(resourceGroupName string, projectName string, groupName string, assessmentName string, options *armmigrationassessment.AssessedWebAppV2OperationsClientListByWebAppAssessmentV2Options) (resp azfake.PagerResponder[armmigrationassessment.AssessedWebAppV2OperationsClientListByWebAppAssessmentV2Response])
}

// NewAssessedWebAppV2OperationsServerTransport creates a new instance of AssessedWebAppV2OperationsServerTransport with the provided implementation.
// The returned AssessedWebAppV2OperationsServerTransport instance is connected to an instance of armmigrationassessment.AssessedWebAppV2OperationsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewAssessedWebAppV2OperationsServerTransport(srv *AssessedWebAppV2OperationsServer) *AssessedWebAppV2OperationsServerTransport {
	return &AssessedWebAppV2OperationsServerTransport{
		srv:                              srv,
		newListByWebAppAssessmentV2Pager: newTracker[azfake.PagerResponder[armmigrationassessment.AssessedWebAppV2OperationsClientListByWebAppAssessmentV2Response]](),
	}
}

// AssessedWebAppV2OperationsServerTransport connects instances of armmigrationassessment.AssessedWebAppV2OperationsClient to instances of AssessedWebAppV2OperationsServer.
// Don't use this type directly, use NewAssessedWebAppV2OperationsServerTransport instead.
type AssessedWebAppV2OperationsServerTransport struct {
	srv                              *AssessedWebAppV2OperationsServer
	newListByWebAppAssessmentV2Pager *tracker[azfake.PagerResponder[armmigrationassessment.AssessedWebAppV2OperationsClientListByWebAppAssessmentV2Response]]
}

// Do implements the policy.Transporter interface for AssessedWebAppV2OperationsServerTransport.
func (a *AssessedWebAppV2OperationsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return a.dispatchToMethodFake(req, method)
}

func (a *AssessedWebAppV2OperationsServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if assessedWebAppV2OperationsServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = assessedWebAppV2OperationsServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "AssessedWebAppV2OperationsClient.Get":
				res.resp, res.err = a.dispatchGet(req)
			case "AssessedWebAppV2OperationsClient.NewListByWebAppAssessmentV2Pager":
				res.resp, res.err = a.dispatchNewListByWebAppAssessmentV2Pager(req)
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

func (a *AssessedWebAppV2OperationsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if a.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Migrate/assessmentProjects/(?P<projectName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/groups/(?P<groupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/webAppAssessments/(?P<assessmentName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/assessedWebApps/(?P<assessedWebAppName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
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
	assessedWebAppNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("assessedWebAppName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := a.srv.Get(req.Context(), resourceGroupNameParam, projectNameParam, groupNameParam, assessmentNameParam, assessedWebAppNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).AssessedWebAppV2, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *AssessedWebAppV2OperationsServerTransport) dispatchNewListByWebAppAssessmentV2Pager(req *http.Request) (*http.Response, error) {
	if a.srv.NewListByWebAppAssessmentV2Pager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByWebAppAssessmentV2Pager not implemented")}
	}
	newListByWebAppAssessmentV2Pager := a.newListByWebAppAssessmentV2Pager.get(req)
	if newListByWebAppAssessmentV2Pager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Migrate/assessmentProjects/(?P<projectName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/groups/(?P<groupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/webAppAssessments/(?P<assessmentName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/assessedWebApps`
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
		var options *armmigrationassessment.AssessedWebAppV2OperationsClientListByWebAppAssessmentV2Options
		if filterParam != nil || pageSizeParam != nil || continuationTokenParam != nil || totalRecordCountParam != nil {
			options = &armmigrationassessment.AssessedWebAppV2OperationsClientListByWebAppAssessmentV2Options{
				Filter:            filterParam,
				PageSize:          pageSizeParam,
				ContinuationToken: continuationTokenParam,
				TotalRecordCount:  totalRecordCountParam,
			}
		}
		resp := a.srv.NewListByWebAppAssessmentV2Pager(resourceGroupNameParam, projectNameParam, groupNameParam, assessmentNameParam, options)
		newListByWebAppAssessmentV2Pager = &resp
		a.newListByWebAppAssessmentV2Pager.add(req, newListByWebAppAssessmentV2Pager)
		server.PagerResponderInjectNextLinks(newListByWebAppAssessmentV2Pager, req, func(page *armmigrationassessment.AssessedWebAppV2OperationsClientListByWebAppAssessmentV2Response, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByWebAppAssessmentV2Pager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		a.newListByWebAppAssessmentV2Pager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByWebAppAssessmentV2Pager) {
		a.newListByWebAppAssessmentV2Pager.remove(req)
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to AssessedWebAppV2OperationsServerTransport
var assessedWebAppV2OperationsServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
