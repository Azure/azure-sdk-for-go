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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7"
	"net/http"
	"net/url"
	"regexp"
)

// SnapshotsServer is a fake server for instances of the armcompute.SnapshotsClient type.
type SnapshotsServer struct {
	// BeginCreateOrUpdate is the fake for method SnapshotsClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, snapshotName string, snapshot armcompute.Snapshot, options *armcompute.SnapshotsClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armcompute.SnapshotsClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method SnapshotsClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, snapshotName string, options *armcompute.SnapshotsClientBeginDeleteOptions) (resp azfake.PollerResponder[armcompute.SnapshotsClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method SnapshotsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, snapshotName string, options *armcompute.SnapshotsClientGetOptions) (resp azfake.Responder[armcompute.SnapshotsClientGetResponse], errResp azfake.ErrorResponder)

	// BeginGrantAccess is the fake for method SnapshotsClient.BeginGrantAccess
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginGrantAccess func(ctx context.Context, resourceGroupName string, snapshotName string, grantAccessData armcompute.GrantAccessData, options *armcompute.SnapshotsClientBeginGrantAccessOptions) (resp azfake.PollerResponder[armcompute.SnapshotsClientGrantAccessResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method SnapshotsClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(options *armcompute.SnapshotsClientListOptions) (resp azfake.PagerResponder[armcompute.SnapshotsClientListResponse])

	// NewListByResourceGroupPager is the fake for method SnapshotsClient.NewListByResourceGroupPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByResourceGroupPager func(resourceGroupName string, options *armcompute.SnapshotsClientListByResourceGroupOptions) (resp azfake.PagerResponder[armcompute.SnapshotsClientListByResourceGroupResponse])

	// BeginRevokeAccess is the fake for method SnapshotsClient.BeginRevokeAccess
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginRevokeAccess func(ctx context.Context, resourceGroupName string, snapshotName string, options *armcompute.SnapshotsClientBeginRevokeAccessOptions) (resp azfake.PollerResponder[armcompute.SnapshotsClientRevokeAccessResponse], errResp azfake.ErrorResponder)

	// BeginUpdate is the fake for method SnapshotsClient.BeginUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginUpdate func(ctx context.Context, resourceGroupName string, snapshotName string, snapshot armcompute.SnapshotUpdate, options *armcompute.SnapshotsClientBeginUpdateOptions) (resp azfake.PollerResponder[armcompute.SnapshotsClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewSnapshotsServerTransport creates a new instance of SnapshotsServerTransport with the provided implementation.
// The returned SnapshotsServerTransport instance is connected to an instance of armcompute.SnapshotsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewSnapshotsServerTransport(srv *SnapshotsServer) *SnapshotsServerTransport {
	return &SnapshotsServerTransport{
		srv:                         srv,
		beginCreateOrUpdate:         newTracker[azfake.PollerResponder[armcompute.SnapshotsClientCreateOrUpdateResponse]](),
		beginDelete:                 newTracker[azfake.PollerResponder[armcompute.SnapshotsClientDeleteResponse]](),
		beginGrantAccess:            newTracker[azfake.PollerResponder[armcompute.SnapshotsClientGrantAccessResponse]](),
		newListPager:                newTracker[azfake.PagerResponder[armcompute.SnapshotsClientListResponse]](),
		newListByResourceGroupPager: newTracker[azfake.PagerResponder[armcompute.SnapshotsClientListByResourceGroupResponse]](),
		beginRevokeAccess:           newTracker[azfake.PollerResponder[armcompute.SnapshotsClientRevokeAccessResponse]](),
		beginUpdate:                 newTracker[azfake.PollerResponder[armcompute.SnapshotsClientUpdateResponse]](),
	}
}

// SnapshotsServerTransport connects instances of armcompute.SnapshotsClient to instances of SnapshotsServer.
// Don't use this type directly, use NewSnapshotsServerTransport instead.
type SnapshotsServerTransport struct {
	srv                         *SnapshotsServer
	beginCreateOrUpdate         *tracker[azfake.PollerResponder[armcompute.SnapshotsClientCreateOrUpdateResponse]]
	beginDelete                 *tracker[azfake.PollerResponder[armcompute.SnapshotsClientDeleteResponse]]
	beginGrantAccess            *tracker[azfake.PollerResponder[armcompute.SnapshotsClientGrantAccessResponse]]
	newListPager                *tracker[azfake.PagerResponder[armcompute.SnapshotsClientListResponse]]
	newListByResourceGroupPager *tracker[azfake.PagerResponder[armcompute.SnapshotsClientListByResourceGroupResponse]]
	beginRevokeAccess           *tracker[azfake.PollerResponder[armcompute.SnapshotsClientRevokeAccessResponse]]
	beginUpdate                 *tracker[azfake.PollerResponder[armcompute.SnapshotsClientUpdateResponse]]
}

// Do implements the policy.Transporter interface for SnapshotsServerTransport.
func (s *SnapshotsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return s.dispatchToMethodFake(req, method)
}

func (s *SnapshotsServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if snapshotsServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = snapshotsServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "SnapshotsClient.BeginCreateOrUpdate":
				res.resp, res.err = s.dispatchBeginCreateOrUpdate(req)
			case "SnapshotsClient.BeginDelete":
				res.resp, res.err = s.dispatchBeginDelete(req)
			case "SnapshotsClient.Get":
				res.resp, res.err = s.dispatchGet(req)
			case "SnapshotsClient.BeginGrantAccess":
				res.resp, res.err = s.dispatchBeginGrantAccess(req)
			case "SnapshotsClient.NewListPager":
				res.resp, res.err = s.dispatchNewListPager(req)
			case "SnapshotsClient.NewListByResourceGroupPager":
				res.resp, res.err = s.dispatchNewListByResourceGroupPager(req)
			case "SnapshotsClient.BeginRevokeAccess":
				res.resp, res.err = s.dispatchBeginRevokeAccess(req)
			case "SnapshotsClient.BeginUpdate":
				res.resp, res.err = s.dispatchBeginUpdate(req)
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

func (s *SnapshotsServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if s.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := s.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Compute/snapshots/(?P<snapshotName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armcompute.Snapshot](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		snapshotNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("snapshotName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := s.srv.BeginCreateOrUpdate(req.Context(), resourceGroupNameParam, snapshotNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateOrUpdate = &respr
		s.beginCreateOrUpdate.add(req, beginCreateOrUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		s.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		s.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (s *SnapshotsServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if s.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := s.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Compute/snapshots/(?P<snapshotName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		snapshotNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("snapshotName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := s.srv.BeginDelete(req.Context(), resourceGroupNameParam, snapshotNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		s.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		s.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		s.beginDelete.remove(req)
	}

	return resp, nil
}

func (s *SnapshotsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if s.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Compute/snapshots/(?P<snapshotName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	snapshotNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("snapshotName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.Get(req.Context(), resourceGroupNameParam, snapshotNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).Snapshot, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *SnapshotsServerTransport) dispatchBeginGrantAccess(req *http.Request) (*http.Response, error) {
	if s.srv.BeginGrantAccess == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginGrantAccess not implemented")}
	}
	beginGrantAccess := s.beginGrantAccess.get(req)
	if beginGrantAccess == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Compute/snapshots/(?P<snapshotName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/beginGetAccess`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armcompute.GrantAccessData](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		snapshotNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("snapshotName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := s.srv.BeginGrantAccess(req.Context(), resourceGroupNameParam, snapshotNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginGrantAccess = &respr
		s.beginGrantAccess.add(req, beginGrantAccess)
	}

	resp, err := server.PollerResponderNext(beginGrantAccess, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		s.beginGrantAccess.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginGrantAccess) {
		s.beginGrantAccess.remove(req)
	}

	return resp, nil
}

func (s *SnapshotsServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if s.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := s.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Compute/snapshots`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resp := s.srv.NewListPager(nil)
		newListPager = &resp
		s.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armcompute.SnapshotsClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		s.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		s.newListPager.remove(req)
	}
	return resp, nil
}

func (s *SnapshotsServerTransport) dispatchNewListByResourceGroupPager(req *http.Request) (*http.Response, error) {
	if s.srv.NewListByResourceGroupPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByResourceGroupPager not implemented")}
	}
	newListByResourceGroupPager := s.newListByResourceGroupPager.get(req)
	if newListByResourceGroupPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Compute/snapshots`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		resp := s.srv.NewListByResourceGroupPager(resourceGroupNameParam, nil)
		newListByResourceGroupPager = &resp
		s.newListByResourceGroupPager.add(req, newListByResourceGroupPager)
		server.PagerResponderInjectNextLinks(newListByResourceGroupPager, req, func(page *armcompute.SnapshotsClientListByResourceGroupResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByResourceGroupPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		s.newListByResourceGroupPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByResourceGroupPager) {
		s.newListByResourceGroupPager.remove(req)
	}
	return resp, nil
}

func (s *SnapshotsServerTransport) dispatchBeginRevokeAccess(req *http.Request) (*http.Response, error) {
	if s.srv.BeginRevokeAccess == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginRevokeAccess not implemented")}
	}
	beginRevokeAccess := s.beginRevokeAccess.get(req)
	if beginRevokeAccess == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Compute/snapshots/(?P<snapshotName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/endgetaccess`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		snapshotNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("snapshotName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := s.srv.BeginRevokeAccess(req.Context(), resourceGroupNameParam, snapshotNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginRevokeAccess = &respr
		s.beginRevokeAccess.add(req, beginRevokeAccess)
	}

	resp, err := server.PollerResponderNext(beginRevokeAccess, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		s.beginRevokeAccess.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginRevokeAccess) {
		s.beginRevokeAccess.remove(req)
	}

	return resp, nil
}

func (s *SnapshotsServerTransport) dispatchBeginUpdate(req *http.Request) (*http.Response, error) {
	if s.srv.BeginUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdate not implemented")}
	}
	beginUpdate := s.beginUpdate.get(req)
	if beginUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Compute/snapshots/(?P<snapshotName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armcompute.SnapshotUpdate](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		snapshotNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("snapshotName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := s.srv.BeginUpdate(req.Context(), resourceGroupNameParam, snapshotNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginUpdate = &respr
		s.beginUpdate.add(req, beginUpdate)
	}

	resp, err := server.PollerResponderNext(beginUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		s.beginUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdate) {
		s.beginUpdate.remove(req)
	}

	return resp, nil
}

// set this to conditionally intercept incoming requests to SnapshotsServerTransport
var snapshotsServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
