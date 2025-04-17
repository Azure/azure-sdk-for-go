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

// CapacityReservationsServer is a fake server for instances of the armcompute.CapacityReservationsClient type.
type CapacityReservationsServer struct {
	// BeginCreateOrUpdate is the fake for method CapacityReservationsClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, capacityReservationGroupName string, capacityReservationName string, parameters armcompute.CapacityReservation, options *armcompute.CapacityReservationsClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armcompute.CapacityReservationsClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method CapacityReservationsClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, capacityReservationGroupName string, capacityReservationName string, options *armcompute.CapacityReservationsClientBeginDeleteOptions) (resp azfake.PollerResponder[armcompute.CapacityReservationsClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method CapacityReservationsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, capacityReservationGroupName string, capacityReservationName string, options *armcompute.CapacityReservationsClientGetOptions) (resp azfake.Responder[armcompute.CapacityReservationsClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByCapacityReservationGroupPager is the fake for method CapacityReservationsClient.NewListByCapacityReservationGroupPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByCapacityReservationGroupPager func(resourceGroupName string, capacityReservationGroupName string, options *armcompute.CapacityReservationsClientListByCapacityReservationGroupOptions) (resp azfake.PagerResponder[armcompute.CapacityReservationsClientListByCapacityReservationGroupResponse])

	// BeginUpdate is the fake for method CapacityReservationsClient.BeginUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginUpdate func(ctx context.Context, resourceGroupName string, capacityReservationGroupName string, capacityReservationName string, parameters armcompute.CapacityReservationUpdate, options *armcompute.CapacityReservationsClientBeginUpdateOptions) (resp azfake.PollerResponder[armcompute.CapacityReservationsClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewCapacityReservationsServerTransport creates a new instance of CapacityReservationsServerTransport with the provided implementation.
// The returned CapacityReservationsServerTransport instance is connected to an instance of armcompute.CapacityReservationsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewCapacityReservationsServerTransport(srv *CapacityReservationsServer) *CapacityReservationsServerTransport {
	return &CapacityReservationsServerTransport{
		srv:                                    srv,
		beginCreateOrUpdate:                    newTracker[azfake.PollerResponder[armcompute.CapacityReservationsClientCreateOrUpdateResponse]](),
		beginDelete:                            newTracker[azfake.PollerResponder[armcompute.CapacityReservationsClientDeleteResponse]](),
		newListByCapacityReservationGroupPager: newTracker[azfake.PagerResponder[armcompute.CapacityReservationsClientListByCapacityReservationGroupResponse]](),
		beginUpdate:                            newTracker[azfake.PollerResponder[armcompute.CapacityReservationsClientUpdateResponse]](),
	}
}

// CapacityReservationsServerTransport connects instances of armcompute.CapacityReservationsClient to instances of CapacityReservationsServer.
// Don't use this type directly, use NewCapacityReservationsServerTransport instead.
type CapacityReservationsServerTransport struct {
	srv                                    *CapacityReservationsServer
	beginCreateOrUpdate                    *tracker[azfake.PollerResponder[armcompute.CapacityReservationsClientCreateOrUpdateResponse]]
	beginDelete                            *tracker[azfake.PollerResponder[armcompute.CapacityReservationsClientDeleteResponse]]
	newListByCapacityReservationGroupPager *tracker[azfake.PagerResponder[armcompute.CapacityReservationsClientListByCapacityReservationGroupResponse]]
	beginUpdate                            *tracker[azfake.PollerResponder[armcompute.CapacityReservationsClientUpdateResponse]]
}

// Do implements the policy.Transporter interface for CapacityReservationsServerTransport.
func (c *CapacityReservationsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return c.dispatchToMethodFake(req, method)
}

func (c *CapacityReservationsServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if capacityReservationsServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = capacityReservationsServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "CapacityReservationsClient.BeginCreateOrUpdate":
				res.resp, res.err = c.dispatchBeginCreateOrUpdate(req)
			case "CapacityReservationsClient.BeginDelete":
				res.resp, res.err = c.dispatchBeginDelete(req)
			case "CapacityReservationsClient.Get":
				res.resp, res.err = c.dispatchGet(req)
			case "CapacityReservationsClient.NewListByCapacityReservationGroupPager":
				res.resp, res.err = c.dispatchNewListByCapacityReservationGroupPager(req)
			case "CapacityReservationsClient.BeginUpdate":
				res.resp, res.err = c.dispatchBeginUpdate(req)
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

func (c *CapacityReservationsServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if c.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := c.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Compute/capacityReservationGroups/(?P<capacityReservationGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/capacityReservations/(?P<capacityReservationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armcompute.CapacityReservation](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		capacityReservationGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("capacityReservationGroupName")])
		if err != nil {
			return nil, err
		}
		capacityReservationNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("capacityReservationName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := c.srv.BeginCreateOrUpdate(req.Context(), resourceGroupNameParam, capacityReservationGroupNameParam, capacityReservationNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateOrUpdate = &respr
		c.beginCreateOrUpdate.add(req, beginCreateOrUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		c.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		c.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (c *CapacityReservationsServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if c.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := c.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Compute/capacityReservationGroups/(?P<capacityReservationGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/capacityReservations/(?P<capacityReservationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		capacityReservationGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("capacityReservationGroupName")])
		if err != nil {
			return nil, err
		}
		capacityReservationNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("capacityReservationName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := c.srv.BeginDelete(req.Context(), resourceGroupNameParam, capacityReservationGroupNameParam, capacityReservationNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		c.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		c.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		c.beginDelete.remove(req)
	}

	return resp, nil
}

func (c *CapacityReservationsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if c.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Compute/capacityReservationGroups/(?P<capacityReservationGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/capacityReservations/(?P<capacityReservationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
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
	capacityReservationGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("capacityReservationGroupName")])
	if err != nil {
		return nil, err
	}
	capacityReservationNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("capacityReservationName")])
	if err != nil {
		return nil, err
	}
	expandUnescaped, err := url.QueryUnescape(qp.Get("$expand"))
	if err != nil {
		return nil, err
	}
	expandParam := getOptional(armcompute.CapacityReservationInstanceViewTypes(expandUnescaped))
	var options *armcompute.CapacityReservationsClientGetOptions
	if expandParam != nil {
		options = &armcompute.CapacityReservationsClientGetOptions{
			Expand: expandParam,
		}
	}
	respr, errRespr := c.srv.Get(req.Context(), resourceGroupNameParam, capacityReservationGroupNameParam, capacityReservationNameParam, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).CapacityReservation, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *CapacityReservationsServerTransport) dispatchNewListByCapacityReservationGroupPager(req *http.Request) (*http.Response, error) {
	if c.srv.NewListByCapacityReservationGroupPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByCapacityReservationGroupPager not implemented")}
	}
	newListByCapacityReservationGroupPager := c.newListByCapacityReservationGroupPager.get(req)
	if newListByCapacityReservationGroupPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Compute/capacityReservationGroups/(?P<capacityReservationGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/capacityReservations`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		capacityReservationGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("capacityReservationGroupName")])
		if err != nil {
			return nil, err
		}
		resp := c.srv.NewListByCapacityReservationGroupPager(resourceGroupNameParam, capacityReservationGroupNameParam, nil)
		newListByCapacityReservationGroupPager = &resp
		c.newListByCapacityReservationGroupPager.add(req, newListByCapacityReservationGroupPager)
		server.PagerResponderInjectNextLinks(newListByCapacityReservationGroupPager, req, func(page *armcompute.CapacityReservationsClientListByCapacityReservationGroupResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByCapacityReservationGroupPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		c.newListByCapacityReservationGroupPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByCapacityReservationGroupPager) {
		c.newListByCapacityReservationGroupPager.remove(req)
	}
	return resp, nil
}

func (c *CapacityReservationsServerTransport) dispatchBeginUpdate(req *http.Request) (*http.Response, error) {
	if c.srv.BeginUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdate not implemented")}
	}
	beginUpdate := c.beginUpdate.get(req)
	if beginUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Compute/capacityReservationGroups/(?P<capacityReservationGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/capacityReservations/(?P<capacityReservationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armcompute.CapacityReservationUpdate](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		capacityReservationGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("capacityReservationGroupName")])
		if err != nil {
			return nil, err
		}
		capacityReservationNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("capacityReservationName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := c.srv.BeginUpdate(req.Context(), resourceGroupNameParam, capacityReservationGroupNameParam, capacityReservationNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginUpdate = &respr
		c.beginUpdate.add(req, beginUpdate)
	}

	resp, err := server.PollerResponderNext(beginUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		c.beginUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdate) {
		c.beginUpdate.remove(req)
	}

	return resp, nil
}

// set this to conditionally intercept incoming requests to CapacityReservationsServerTransport
var capacityReservationsServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
