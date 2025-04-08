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

// VmwareCollectorsOperationsServer is a fake server for instances of the armmigrationassessment.VmwareCollectorsOperationsClient type.
type VmwareCollectorsOperationsServer struct {
	// BeginCreate is the fake for method VmwareCollectorsOperationsClient.BeginCreate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreate func(ctx context.Context, resourceGroupName string, projectName string, vmWareCollectorName string, resource armmigrationassessment.VmwareCollector, options *armmigrationassessment.VmwareCollectorsOperationsClientBeginCreateOptions) (resp azfake.PollerResponder[armmigrationassessment.VmwareCollectorsOperationsClientCreateResponse], errResp azfake.ErrorResponder)

	// Delete is the fake for method VmwareCollectorsOperationsClient.Delete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	Delete func(ctx context.Context, resourceGroupName string, projectName string, vmWareCollectorName string, options *armmigrationassessment.VmwareCollectorsOperationsClientDeleteOptions) (resp azfake.Responder[armmigrationassessment.VmwareCollectorsOperationsClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method VmwareCollectorsOperationsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, projectName string, vmWareCollectorName string, options *armmigrationassessment.VmwareCollectorsOperationsClientGetOptions) (resp azfake.Responder[armmigrationassessment.VmwareCollectorsOperationsClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByAssessmentProjectPager is the fake for method VmwareCollectorsOperationsClient.NewListByAssessmentProjectPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByAssessmentProjectPager func(resourceGroupName string, projectName string, options *armmigrationassessment.VmwareCollectorsOperationsClientListByAssessmentProjectOptions) (resp azfake.PagerResponder[armmigrationassessment.VmwareCollectorsOperationsClientListByAssessmentProjectResponse])
}

// NewVmwareCollectorsOperationsServerTransport creates a new instance of VmwareCollectorsOperationsServerTransport with the provided implementation.
// The returned VmwareCollectorsOperationsServerTransport instance is connected to an instance of armmigrationassessment.VmwareCollectorsOperationsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewVmwareCollectorsOperationsServerTransport(srv *VmwareCollectorsOperationsServer) *VmwareCollectorsOperationsServerTransport {
	return &VmwareCollectorsOperationsServerTransport{
		srv:                             srv,
		beginCreate:                     newTracker[azfake.PollerResponder[armmigrationassessment.VmwareCollectorsOperationsClientCreateResponse]](),
		newListByAssessmentProjectPager: newTracker[azfake.PagerResponder[armmigrationassessment.VmwareCollectorsOperationsClientListByAssessmentProjectResponse]](),
	}
}

// VmwareCollectorsOperationsServerTransport connects instances of armmigrationassessment.VmwareCollectorsOperationsClient to instances of VmwareCollectorsOperationsServer.
// Don't use this type directly, use NewVmwareCollectorsOperationsServerTransport instead.
type VmwareCollectorsOperationsServerTransport struct {
	srv                             *VmwareCollectorsOperationsServer
	beginCreate                     *tracker[azfake.PollerResponder[armmigrationassessment.VmwareCollectorsOperationsClientCreateResponse]]
	newListByAssessmentProjectPager *tracker[azfake.PagerResponder[armmigrationassessment.VmwareCollectorsOperationsClientListByAssessmentProjectResponse]]
}

// Do implements the policy.Transporter interface for VmwareCollectorsOperationsServerTransport.
func (v *VmwareCollectorsOperationsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return v.dispatchToMethodFake(req, method)
}

func (v *VmwareCollectorsOperationsServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if vmwareCollectorsOperationsServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = vmwareCollectorsOperationsServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "VmwareCollectorsOperationsClient.BeginCreate":
				res.resp, res.err = v.dispatchBeginCreate(req)
			case "VmwareCollectorsOperationsClient.Delete":
				res.resp, res.err = v.dispatchDelete(req)
			case "VmwareCollectorsOperationsClient.Get":
				res.resp, res.err = v.dispatchGet(req)
			case "VmwareCollectorsOperationsClient.NewListByAssessmentProjectPager":
				res.resp, res.err = v.dispatchNewListByAssessmentProjectPager(req)
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

func (v *VmwareCollectorsOperationsServerTransport) dispatchBeginCreate(req *http.Request) (*http.Response, error) {
	if v.srv.BeginCreate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreate not implemented")}
	}
	beginCreate := v.beginCreate.get(req)
	if beginCreate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Migrate/assessmentProjects/(?P<projectName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/vmwarecollectors/(?P<vmWareCollectorName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armmigrationassessment.VmwareCollector](req)
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
		vmWareCollectorNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("vmWareCollectorName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := v.srv.BeginCreate(req.Context(), resourceGroupNameParam, projectNameParam, vmWareCollectorNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreate = &respr
		v.beginCreate.add(req, beginCreate)
	}

	resp, err := server.PollerResponderNext(beginCreate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		v.beginCreate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreate) {
		v.beginCreate.remove(req)
	}

	return resp, nil
}

func (v *VmwareCollectorsOperationsServerTransport) dispatchDelete(req *http.Request) (*http.Response, error) {
	if v.srv.Delete == nil {
		return nil, &nonRetriableError{errors.New("fake for method Delete not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Migrate/assessmentProjects/(?P<projectName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/vmwarecollectors/(?P<vmWareCollectorName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
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
	vmWareCollectorNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("vmWareCollectorName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := v.srv.Delete(req.Context(), resourceGroupNameParam, projectNameParam, vmWareCollectorNameParam, nil)
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

func (v *VmwareCollectorsOperationsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if v.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Migrate/assessmentProjects/(?P<projectName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/vmwarecollectors/(?P<vmWareCollectorName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
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
	vmWareCollectorNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("vmWareCollectorName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := v.srv.Get(req.Context(), resourceGroupNameParam, projectNameParam, vmWareCollectorNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).VmwareCollector, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (v *VmwareCollectorsOperationsServerTransport) dispatchNewListByAssessmentProjectPager(req *http.Request) (*http.Response, error) {
	if v.srv.NewListByAssessmentProjectPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByAssessmentProjectPager not implemented")}
	}
	newListByAssessmentProjectPager := v.newListByAssessmentProjectPager.get(req)
	if newListByAssessmentProjectPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Migrate/assessmentProjects/(?P<projectName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/vmwarecollectors`
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
		resp := v.srv.NewListByAssessmentProjectPager(resourceGroupNameParam, projectNameParam, nil)
		newListByAssessmentProjectPager = &resp
		v.newListByAssessmentProjectPager.add(req, newListByAssessmentProjectPager)
		server.PagerResponderInjectNextLinks(newListByAssessmentProjectPager, req, func(page *armmigrationassessment.VmwareCollectorsOperationsClientListByAssessmentProjectResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByAssessmentProjectPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		v.newListByAssessmentProjectPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByAssessmentProjectPager) {
		v.newListByAssessmentProjectPager.remove(req)
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to VmwareCollectorsOperationsServerTransport
var vmwareCollectorsOperationsServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
