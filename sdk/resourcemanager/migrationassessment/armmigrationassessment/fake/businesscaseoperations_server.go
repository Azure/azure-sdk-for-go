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
)

// BusinessCaseOperationsServer is a fake server for instances of the armmigrationassessment.BusinessCaseOperationsClient type.
type BusinessCaseOperationsServer struct {
	// BeginCompareSummary is the fake for method BusinessCaseOperationsClient.BeginCompareSummary
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginCompareSummary func(ctx context.Context, resourceGroupName string, projectName string, businessCaseName string, body any, options *armmigrationassessment.BusinessCaseOperationsClientBeginCompareSummaryOptions) (resp azfake.PollerResponder[armmigrationassessment.BusinessCaseOperationsClientCompareSummaryResponse], errResp azfake.ErrorResponder)

	// BeginCreate is the fake for method BusinessCaseOperationsClient.BeginCreate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreate func(ctx context.Context, resourceGroupName string, projectName string, businessCaseName string, resource armmigrationassessment.BusinessCase, options *armmigrationassessment.BusinessCaseOperationsClientBeginCreateOptions) (resp azfake.PollerResponder[armmigrationassessment.BusinessCaseOperationsClientCreateResponse], errResp azfake.ErrorResponder)

	// Delete is the fake for method BusinessCaseOperationsClient.Delete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	Delete func(ctx context.Context, resourceGroupName string, projectName string, businessCaseName string, options *armmigrationassessment.BusinessCaseOperationsClientDeleteOptions) (resp azfake.Responder[armmigrationassessment.BusinessCaseOperationsClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method BusinessCaseOperationsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, projectName string, businessCaseName string, options *armmigrationassessment.BusinessCaseOperationsClientGetOptions) (resp azfake.Responder[armmigrationassessment.BusinessCaseOperationsClientGetResponse], errResp azfake.ErrorResponder)

	// BeginGetReportDownloadURL is the fake for method BusinessCaseOperationsClient.BeginGetReportDownloadURL
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginGetReportDownloadURL func(ctx context.Context, resourceGroupName string, projectName string, businessCaseName string, body any, options *armmigrationassessment.BusinessCaseOperationsClientBeginGetReportDownloadURLOptions) (resp azfake.PollerResponder[armmigrationassessment.BusinessCaseOperationsClientGetReportDownloadURLResponse], errResp azfake.ErrorResponder)

	// NewListByAssessmentProjectPager is the fake for method BusinessCaseOperationsClient.NewListByAssessmentProjectPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByAssessmentProjectPager func(resourceGroupName string, projectName string, options *armmigrationassessment.BusinessCaseOperationsClientListByAssessmentProjectOptions) (resp azfake.PagerResponder[armmigrationassessment.BusinessCaseOperationsClientListByAssessmentProjectResponse])
}

// NewBusinessCaseOperationsServerTransport creates a new instance of BusinessCaseOperationsServerTransport with the provided implementation.
// The returned BusinessCaseOperationsServerTransport instance is connected to an instance of armmigrationassessment.BusinessCaseOperationsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewBusinessCaseOperationsServerTransport(srv *BusinessCaseOperationsServer) *BusinessCaseOperationsServerTransport {
	return &BusinessCaseOperationsServerTransport{
		srv:                             srv,
		beginCompareSummary:             newTracker[azfake.PollerResponder[armmigrationassessment.BusinessCaseOperationsClientCompareSummaryResponse]](),
		beginCreate:                     newTracker[azfake.PollerResponder[armmigrationassessment.BusinessCaseOperationsClientCreateResponse]](),
		beginGetReportDownloadURL:       newTracker[azfake.PollerResponder[armmigrationassessment.BusinessCaseOperationsClientGetReportDownloadURLResponse]](),
		newListByAssessmentProjectPager: newTracker[azfake.PagerResponder[armmigrationassessment.BusinessCaseOperationsClientListByAssessmentProjectResponse]](),
	}
}

// BusinessCaseOperationsServerTransport connects instances of armmigrationassessment.BusinessCaseOperationsClient to instances of BusinessCaseOperationsServer.
// Don't use this type directly, use NewBusinessCaseOperationsServerTransport instead.
type BusinessCaseOperationsServerTransport struct {
	srv                             *BusinessCaseOperationsServer
	beginCompareSummary             *tracker[azfake.PollerResponder[armmigrationassessment.BusinessCaseOperationsClientCompareSummaryResponse]]
	beginCreate                     *tracker[azfake.PollerResponder[armmigrationassessment.BusinessCaseOperationsClientCreateResponse]]
	beginGetReportDownloadURL       *tracker[azfake.PollerResponder[armmigrationassessment.BusinessCaseOperationsClientGetReportDownloadURLResponse]]
	newListByAssessmentProjectPager *tracker[azfake.PagerResponder[armmigrationassessment.BusinessCaseOperationsClientListByAssessmentProjectResponse]]
}

// Do implements the policy.Transporter interface for BusinessCaseOperationsServerTransport.
func (b *BusinessCaseOperationsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return b.dispatchToMethodFake(req, method)
}

func (b *BusinessCaseOperationsServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if businessCaseOperationsServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = businessCaseOperationsServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "BusinessCaseOperationsClient.BeginCompareSummary":
				res.resp, res.err = b.dispatchBeginCompareSummary(req)
			case "BusinessCaseOperationsClient.BeginCreate":
				res.resp, res.err = b.dispatchBeginCreate(req)
			case "BusinessCaseOperationsClient.Delete":
				res.resp, res.err = b.dispatchDelete(req)
			case "BusinessCaseOperationsClient.Get":
				res.resp, res.err = b.dispatchGet(req)
			case "BusinessCaseOperationsClient.BeginGetReportDownloadURL":
				res.resp, res.err = b.dispatchBeginGetReportDownloadURL(req)
			case "BusinessCaseOperationsClient.NewListByAssessmentProjectPager":
				res.resp, res.err = b.dispatchNewListByAssessmentProjectPager(req)
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

func (b *BusinessCaseOperationsServerTransport) dispatchBeginCompareSummary(req *http.Request) (*http.Response, error) {
	if b.srv.BeginCompareSummary == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCompareSummary not implemented")}
	}
	beginCompareSummary := b.beginCompareSummary.get(req)
	if beginCompareSummary == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Migrate/assessmentProjects/(?P<projectName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/businessCases/(?P<businessCaseName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/compareSummary`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[any](req)
		if err != nil {
			return nil, err
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
		respr, errRespr := b.srv.BeginCompareSummary(req.Context(), resourceGroupNameParam, projectNameParam, businessCaseNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCompareSummary = &respr
		b.beginCompareSummary.add(req, beginCompareSummary)
	}

	resp, err := server.PollerResponderNext(beginCompareSummary, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		b.beginCompareSummary.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCompareSummary) {
		b.beginCompareSummary.remove(req)
	}

	return resp, nil
}

func (b *BusinessCaseOperationsServerTransport) dispatchBeginCreate(req *http.Request) (*http.Response, error) {
	if b.srv.BeginCreate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreate not implemented")}
	}
	beginCreate := b.beginCreate.get(req)
	if beginCreate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Migrate/assessmentProjects/(?P<projectName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/businessCases/(?P<businessCaseName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armmigrationassessment.BusinessCase](req)
		if err != nil {
			return nil, err
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
		respr, errRespr := b.srv.BeginCreate(req.Context(), resourceGroupNameParam, projectNameParam, businessCaseNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreate = &respr
		b.beginCreate.add(req, beginCreate)
	}

	resp, err := server.PollerResponderNext(beginCreate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		b.beginCreate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreate) {
		b.beginCreate.remove(req)
	}

	return resp, nil
}

func (b *BusinessCaseOperationsServerTransport) dispatchDelete(req *http.Request) (*http.Response, error) {
	if b.srv.Delete == nil {
		return nil, &nonRetriableError{errors.New("fake for method Delete not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Migrate/assessmentProjects/(?P<projectName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/businessCases/(?P<businessCaseName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
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
	respr, errRespr := b.srv.Delete(req.Context(), resourceGroupNameParam, projectNameParam, businessCaseNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusNoContent}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusNoContent", respContent.HTTPStatus)}
	}
	resp, err := server.NewResponse(respContent, req, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (b *BusinessCaseOperationsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if b.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Migrate/assessmentProjects/(?P<projectName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/businessCases/(?P<businessCaseName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
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
	respr, errRespr := b.srv.Get(req.Context(), resourceGroupNameParam, projectNameParam, businessCaseNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).BusinessCase, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (b *BusinessCaseOperationsServerTransport) dispatchBeginGetReportDownloadURL(req *http.Request) (*http.Response, error) {
	if b.srv.BeginGetReportDownloadURL == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginGetReportDownloadURL not implemented")}
	}
	beginGetReportDownloadURL := b.beginGetReportDownloadURL.get(req)
	if beginGetReportDownloadURL == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Migrate/assessmentProjects/(?P<projectName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/businessCases/(?P<businessCaseName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/getReportDownloadUrl`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[any](req)
		if err != nil {
			return nil, err
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
		respr, errRespr := b.srv.BeginGetReportDownloadURL(req.Context(), resourceGroupNameParam, projectNameParam, businessCaseNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginGetReportDownloadURL = &respr
		b.beginGetReportDownloadURL.add(req, beginGetReportDownloadURL)
	}

	resp, err := server.PollerResponderNext(beginGetReportDownloadURL, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		b.beginGetReportDownloadURL.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginGetReportDownloadURL) {
		b.beginGetReportDownloadURL.remove(req)
	}

	return resp, nil
}

func (b *BusinessCaseOperationsServerTransport) dispatchNewListByAssessmentProjectPager(req *http.Request) (*http.Response, error) {
	if b.srv.NewListByAssessmentProjectPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByAssessmentProjectPager not implemented")}
	}
	newListByAssessmentProjectPager := b.newListByAssessmentProjectPager.get(req)
	if newListByAssessmentProjectPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Migrate/assessmentProjects/(?P<projectName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/businessCases`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
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
		resp := b.srv.NewListByAssessmentProjectPager(resourceGroupNameParam, projectNameParam, nil)
		newListByAssessmentProjectPager = &resp
		b.newListByAssessmentProjectPager.add(req, newListByAssessmentProjectPager)
		server.PagerResponderInjectNextLinks(newListByAssessmentProjectPager, req, func(page *armmigrationassessment.BusinessCaseOperationsClientListByAssessmentProjectResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByAssessmentProjectPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		b.newListByAssessmentProjectPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByAssessmentProjectPager) {
		b.newListByAssessmentProjectPager.remove(req)
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to BusinessCaseOperationsServerTransport
var businessCaseOperationsServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
