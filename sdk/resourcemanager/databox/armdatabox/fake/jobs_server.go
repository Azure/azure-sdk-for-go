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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/databox/armdatabox/v2"
	"net/http"
	"net/url"
	"regexp"
)

// JobsServer is a fake server for instances of the armdatabox.JobsClient type.
type JobsServer struct {
	// BookShipmentPickUp is the fake for method JobsClient.BookShipmentPickUp
	// HTTP status codes to indicate success: http.StatusOK
	BookShipmentPickUp func(ctx context.Context, resourceGroupName string, jobName string, shipmentPickUpRequest armdatabox.ShipmentPickUpRequest, options *armdatabox.JobsClientBookShipmentPickUpOptions) (resp azfake.Responder[armdatabox.JobsClientBookShipmentPickUpResponse], errResp azfake.ErrorResponder)

	// Cancel is the fake for method JobsClient.Cancel
	// HTTP status codes to indicate success: http.StatusNoContent
	Cancel func(ctx context.Context, resourceGroupName string, jobName string, cancellationReason armdatabox.CancellationReason, options *armdatabox.JobsClientCancelOptions) (resp azfake.Responder[armdatabox.JobsClientCancelResponse], errResp azfake.ErrorResponder)

	// BeginCreate is the fake for method JobsClient.BeginCreate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginCreate func(ctx context.Context, resourceGroupName string, jobName string, jobResource armdatabox.JobResource, options *armdatabox.JobsClientBeginCreateOptions) (resp azfake.PollerResponder[armdatabox.JobsClientCreateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method JobsClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, jobName string, options *armdatabox.JobsClientBeginDeleteOptions) (resp azfake.PollerResponder[armdatabox.JobsClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method JobsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, jobName string, options *armdatabox.JobsClientGetOptions) (resp azfake.Responder[armdatabox.JobsClientGetResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method JobsClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(options *armdatabox.JobsClientListOptions) (resp azfake.PagerResponder[armdatabox.JobsClientListResponse])

	// NewListByResourceGroupPager is the fake for method JobsClient.NewListByResourceGroupPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByResourceGroupPager func(resourceGroupName string, options *armdatabox.JobsClientListByResourceGroupOptions) (resp azfake.PagerResponder[armdatabox.JobsClientListByResourceGroupResponse])

	// NewListCredentialsPager is the fake for method JobsClient.NewListCredentialsPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListCredentialsPager func(resourceGroupName string, jobName string, options *armdatabox.JobsClientListCredentialsOptions) (resp azfake.PagerResponder[armdatabox.JobsClientListCredentialsResponse])

	// MarkDevicesShipped is the fake for method JobsClient.MarkDevicesShipped
	// HTTP status codes to indicate success: http.StatusNoContent
	MarkDevicesShipped func(ctx context.Context, jobName string, resourceGroupName string, markDevicesShippedRequest armdatabox.MarkDevicesShippedRequest, options *armdatabox.JobsClientMarkDevicesShippedOptions) (resp azfake.Responder[armdatabox.JobsClientMarkDevicesShippedResponse], errResp azfake.ErrorResponder)

	// BeginUpdate is the fake for method JobsClient.BeginUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginUpdate func(ctx context.Context, resourceGroupName string, jobName string, jobResourceUpdateParameter armdatabox.JobResourceUpdateParameter, options *armdatabox.JobsClientBeginUpdateOptions) (resp azfake.PollerResponder[armdatabox.JobsClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewJobsServerTransport creates a new instance of JobsServerTransport with the provided implementation.
// The returned JobsServerTransport instance is connected to an instance of armdatabox.JobsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewJobsServerTransport(srv *JobsServer) *JobsServerTransport {
	return &JobsServerTransport{
		srv:                         srv,
		beginCreate:                 newTracker[azfake.PollerResponder[armdatabox.JobsClientCreateResponse]](),
		beginDelete:                 newTracker[azfake.PollerResponder[armdatabox.JobsClientDeleteResponse]](),
		newListPager:                newTracker[azfake.PagerResponder[armdatabox.JobsClientListResponse]](),
		newListByResourceGroupPager: newTracker[azfake.PagerResponder[armdatabox.JobsClientListByResourceGroupResponse]](),
		newListCredentialsPager:     newTracker[azfake.PagerResponder[armdatabox.JobsClientListCredentialsResponse]](),
		beginUpdate:                 newTracker[azfake.PollerResponder[armdatabox.JobsClientUpdateResponse]](),
	}
}

// JobsServerTransport connects instances of armdatabox.JobsClient to instances of JobsServer.
// Don't use this type directly, use NewJobsServerTransport instead.
type JobsServerTransport struct {
	srv                         *JobsServer
	beginCreate                 *tracker[azfake.PollerResponder[armdatabox.JobsClientCreateResponse]]
	beginDelete                 *tracker[azfake.PollerResponder[armdatabox.JobsClientDeleteResponse]]
	newListPager                *tracker[azfake.PagerResponder[armdatabox.JobsClientListResponse]]
	newListByResourceGroupPager *tracker[azfake.PagerResponder[armdatabox.JobsClientListByResourceGroupResponse]]
	newListCredentialsPager     *tracker[azfake.PagerResponder[armdatabox.JobsClientListCredentialsResponse]]
	beginUpdate                 *tracker[azfake.PollerResponder[armdatabox.JobsClientUpdateResponse]]
}

// Do implements the policy.Transporter interface for JobsServerTransport.
func (j *JobsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return j.dispatchToMethodFake(req, method)
}

func (j *JobsServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if jobsServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = jobsServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "JobsClient.BookShipmentPickUp":
				res.resp, res.err = j.dispatchBookShipmentPickUp(req)
			case "JobsClient.Cancel":
				res.resp, res.err = j.dispatchCancel(req)
			case "JobsClient.BeginCreate":
				res.resp, res.err = j.dispatchBeginCreate(req)
			case "JobsClient.BeginDelete":
				res.resp, res.err = j.dispatchBeginDelete(req)
			case "JobsClient.Get":
				res.resp, res.err = j.dispatchGet(req)
			case "JobsClient.NewListPager":
				res.resp, res.err = j.dispatchNewListPager(req)
			case "JobsClient.NewListByResourceGroupPager":
				res.resp, res.err = j.dispatchNewListByResourceGroupPager(req)
			case "JobsClient.NewListCredentialsPager":
				res.resp, res.err = j.dispatchNewListCredentialsPager(req)
			case "JobsClient.MarkDevicesShipped":
				res.resp, res.err = j.dispatchMarkDevicesShipped(req)
			case "JobsClient.BeginUpdate":
				res.resp, res.err = j.dispatchBeginUpdate(req)
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

func (j *JobsServerTransport) dispatchBookShipmentPickUp(req *http.Request) (*http.Response, error) {
	if j.srv.BookShipmentPickUp == nil {
		return nil, &nonRetriableError{errors.New("fake for method BookShipmentPickUp not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DataBox/jobs/(?P<jobName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/bookShipmentPickUp`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armdatabox.ShipmentPickUpRequest](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	jobNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("jobName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := j.srv.BookShipmentPickUp(req.Context(), resourceGroupNameParam, jobNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ShipmentPickUpResponse, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (j *JobsServerTransport) dispatchCancel(req *http.Request) (*http.Response, error) {
	if j.srv.Cancel == nil {
		return nil, &nonRetriableError{errors.New("fake for method Cancel not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DataBox/jobs/(?P<jobName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/cancel`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armdatabox.CancellationReason](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	jobNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("jobName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := j.srv.Cancel(req.Context(), resourceGroupNameParam, jobNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusNoContent}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusNoContent", respContent.HTTPStatus)}
	}
	resp, err := server.NewResponse(respContent, req, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (j *JobsServerTransport) dispatchBeginCreate(req *http.Request) (*http.Response, error) {
	if j.srv.BeginCreate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreate not implemented")}
	}
	beginCreate := j.beginCreate.get(req)
	if beginCreate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DataBox/jobs/(?P<jobName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armdatabox.JobResource](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		jobNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("jobName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := j.srv.BeginCreate(req.Context(), resourceGroupNameParam, jobNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreate = &respr
		j.beginCreate.add(req, beginCreate)
	}

	resp, err := server.PollerResponderNext(beginCreate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		j.beginCreate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreate) {
		j.beginCreate.remove(req)
	}

	return resp, nil
}

func (j *JobsServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if j.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := j.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DataBox/jobs/(?P<jobName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		jobNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("jobName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := j.srv.BeginDelete(req.Context(), resourceGroupNameParam, jobNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		j.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		j.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		j.beginDelete.remove(req)
	}

	return resp, nil
}

func (j *JobsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if j.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DataBox/jobs/(?P<jobName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	qp := req.URL.Query()
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	jobNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("jobName")])
	if err != nil {
		return nil, err
	}
	expandUnescaped, err := url.QueryUnescape(qp.Get("$expand"))
	if err != nil {
		return nil, err
	}
	expandParam := getOptional(expandUnescaped)
	var options *armdatabox.JobsClientGetOptions
	if expandParam != nil {
		options = &armdatabox.JobsClientGetOptions{
			Expand: expandParam,
		}
	}
	respr, errRespr := j.srv.Get(req.Context(), resourceGroupNameParam, jobNameParam, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).JobResource, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (j *JobsServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if j.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := j.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DataBox/jobs`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		skipTokenUnescaped, err := url.QueryUnescape(qp.Get("$skipToken"))
		if err != nil {
			return nil, err
		}
		skipTokenParam := getOptional(skipTokenUnescaped)
		var options *armdatabox.JobsClientListOptions
		if skipTokenParam != nil {
			options = &armdatabox.JobsClientListOptions{
				SkipToken: skipTokenParam,
			}
		}
		resp := j.srv.NewListPager(options)
		newListPager = &resp
		j.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armdatabox.JobsClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		j.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		j.newListPager.remove(req)
	}
	return resp, nil
}

func (j *JobsServerTransport) dispatchNewListByResourceGroupPager(req *http.Request) (*http.Response, error) {
	if j.srv.NewListByResourceGroupPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByResourceGroupPager not implemented")}
	}
	newListByResourceGroupPager := j.newListByResourceGroupPager.get(req)
	if newListByResourceGroupPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DataBox/jobs`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		skipTokenUnescaped, err := url.QueryUnescape(qp.Get("$skipToken"))
		if err != nil {
			return nil, err
		}
		skipTokenParam := getOptional(skipTokenUnescaped)
		var options *armdatabox.JobsClientListByResourceGroupOptions
		if skipTokenParam != nil {
			options = &armdatabox.JobsClientListByResourceGroupOptions{
				SkipToken: skipTokenParam,
			}
		}
		resp := j.srv.NewListByResourceGroupPager(resourceGroupNameParam, options)
		newListByResourceGroupPager = &resp
		j.newListByResourceGroupPager.add(req, newListByResourceGroupPager)
		server.PagerResponderInjectNextLinks(newListByResourceGroupPager, req, func(page *armdatabox.JobsClientListByResourceGroupResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByResourceGroupPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		j.newListByResourceGroupPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByResourceGroupPager) {
		j.newListByResourceGroupPager.remove(req)
	}
	return resp, nil
}

func (j *JobsServerTransport) dispatchNewListCredentialsPager(req *http.Request) (*http.Response, error) {
	if j.srv.NewListCredentialsPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListCredentialsPager not implemented")}
	}
	newListCredentialsPager := j.newListCredentialsPager.get(req)
	if newListCredentialsPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DataBox/jobs/(?P<jobName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/listCredentials`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		jobNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("jobName")])
		if err != nil {
			return nil, err
		}
		resp := j.srv.NewListCredentialsPager(resourceGroupNameParam, jobNameParam, nil)
		newListCredentialsPager = &resp
		j.newListCredentialsPager.add(req, newListCredentialsPager)
	}
	resp, err := server.PagerResponderNext(newListCredentialsPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		j.newListCredentialsPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListCredentialsPager) {
		j.newListCredentialsPager.remove(req)
	}
	return resp, nil
}

func (j *JobsServerTransport) dispatchMarkDevicesShipped(req *http.Request) (*http.Response, error) {
	if j.srv.MarkDevicesShipped == nil {
		return nil, &nonRetriableError{errors.New("fake for method MarkDevicesShipped not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DataBox/jobs/(?P<jobName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/markDevicesShipped`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armdatabox.MarkDevicesShippedRequest](req)
	if err != nil {
		return nil, err
	}
	jobNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("jobName")])
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := j.srv.MarkDevicesShipped(req.Context(), jobNameParam, resourceGroupNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusNoContent}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusNoContent", respContent.HTTPStatus)}
	}
	resp, err := server.NewResponse(respContent, req, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (j *JobsServerTransport) dispatchBeginUpdate(req *http.Request) (*http.Response, error) {
	if j.srv.BeginUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdate not implemented")}
	}
	beginUpdate := j.beginUpdate.get(req)
	if beginUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DataBox/jobs/(?P<jobName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armdatabox.JobResourceUpdateParameter](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		jobNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("jobName")])
		if err != nil {
			return nil, err
		}
		ifMatchParam := getOptional(getHeaderValue(req.Header, "If-Match"))
		var options *armdatabox.JobsClientBeginUpdateOptions
		if ifMatchParam != nil {
			options = &armdatabox.JobsClientBeginUpdateOptions{
				IfMatch: ifMatchParam,
			}
		}
		respr, errRespr := j.srv.BeginUpdate(req.Context(), resourceGroupNameParam, jobNameParam, body, options)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginUpdate = &respr
		j.beginUpdate.add(req, beginUpdate)
	}

	resp, err := server.PollerResponderNext(beginUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		j.beginUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdate) {
		j.beginUpdate.remove(req)
	}

	return resp, nil
}

// set this to conditionally intercept incoming requests to JobsServerTransport
var jobsServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
