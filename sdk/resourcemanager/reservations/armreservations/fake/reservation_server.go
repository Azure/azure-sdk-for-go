//go:build go1.18
// +build go1.18

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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/reservations/armreservations/v3"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

// ReservationServer is a fake server for instances of the armreservations.ReservationClient type.
type ReservationServer struct {
	// Archive is the fake for method ReservationClient.Archive
	// HTTP status codes to indicate success: http.StatusOK
	Archive func(ctx context.Context, reservationOrderID string, reservationID string, options *armreservations.ReservationClientArchiveOptions) (resp azfake.Responder[armreservations.ReservationClientArchiveResponse], errResp azfake.ErrorResponder)

	// BeginAvailableScopes is the fake for method ReservationClient.BeginAvailableScopes
	// HTTP status codes to indicate success: http.StatusOK
	BeginAvailableScopes func(ctx context.Context, reservationOrderID string, reservationID string, body armreservations.AvailableScopeRequest, options *armreservations.ReservationClientBeginAvailableScopesOptions) (resp azfake.PollerResponder[armreservations.ReservationClientAvailableScopesResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method ReservationClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, reservationOrderID string, reservationID string, options *armreservations.ReservationClientGetOptions) (resp azfake.Responder[armreservations.ReservationClientGetResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method ReservationClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(reservationOrderID string, options *armreservations.ReservationClientListOptions) (resp azfake.PagerResponder[armreservations.ReservationClientListResponse])

	// NewListAllPager is the fake for method ReservationClient.NewListAllPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListAllPager func(options *armreservations.ReservationClientListAllOptions) (resp azfake.PagerResponder[armreservations.ReservationClientListAllResponse])

	// NewListRevisionsPager is the fake for method ReservationClient.NewListRevisionsPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListRevisionsPager func(reservationOrderID string, reservationID string, options *armreservations.ReservationClientListRevisionsOptions) (resp azfake.PagerResponder[armreservations.ReservationClientListRevisionsResponse])

	// BeginMerge is the fake for method ReservationClient.BeginMerge
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginMerge func(ctx context.Context, reservationOrderID string, body armreservations.MergeRequest, options *armreservations.ReservationClientBeginMergeOptions) (resp azfake.PollerResponder[armreservations.ReservationClientMergeResponse], errResp azfake.ErrorResponder)

	// BeginSplit is the fake for method ReservationClient.BeginSplit
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginSplit func(ctx context.Context, reservationOrderID string, body armreservations.SplitRequest, options *armreservations.ReservationClientBeginSplitOptions) (resp azfake.PollerResponder[armreservations.ReservationClientSplitResponse], errResp azfake.ErrorResponder)

	// Unarchive is the fake for method ReservationClient.Unarchive
	// HTTP status codes to indicate success: http.StatusOK
	Unarchive func(ctx context.Context, reservationOrderID string, reservationID string, options *armreservations.ReservationClientUnarchiveOptions) (resp azfake.Responder[armreservations.ReservationClientUnarchiveResponse], errResp azfake.ErrorResponder)

	// BeginUpdate is the fake for method ReservationClient.BeginUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginUpdate func(ctx context.Context, reservationOrderID string, reservationID string, parameters armreservations.Patch, options *armreservations.ReservationClientBeginUpdateOptions) (resp azfake.PollerResponder[armreservations.ReservationClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewReservationServerTransport creates a new instance of ReservationServerTransport with the provided implementation.
// The returned ReservationServerTransport instance is connected to an instance of armreservations.ReservationClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewReservationServerTransport(srv *ReservationServer) *ReservationServerTransport {
	return &ReservationServerTransport{
		srv:                   srv,
		beginAvailableScopes:  newTracker[azfake.PollerResponder[armreservations.ReservationClientAvailableScopesResponse]](),
		newListPager:          newTracker[azfake.PagerResponder[armreservations.ReservationClientListResponse]](),
		newListAllPager:       newTracker[azfake.PagerResponder[armreservations.ReservationClientListAllResponse]](),
		newListRevisionsPager: newTracker[azfake.PagerResponder[armreservations.ReservationClientListRevisionsResponse]](),
		beginMerge:            newTracker[azfake.PollerResponder[armreservations.ReservationClientMergeResponse]](),
		beginSplit:            newTracker[azfake.PollerResponder[armreservations.ReservationClientSplitResponse]](),
		beginUpdate:           newTracker[azfake.PollerResponder[armreservations.ReservationClientUpdateResponse]](),
	}
}

// ReservationServerTransport connects instances of armreservations.ReservationClient to instances of ReservationServer.
// Don't use this type directly, use NewReservationServerTransport instead.
type ReservationServerTransport struct {
	srv                   *ReservationServer
	beginAvailableScopes  *tracker[azfake.PollerResponder[armreservations.ReservationClientAvailableScopesResponse]]
	newListPager          *tracker[azfake.PagerResponder[armreservations.ReservationClientListResponse]]
	newListAllPager       *tracker[azfake.PagerResponder[armreservations.ReservationClientListAllResponse]]
	newListRevisionsPager *tracker[azfake.PagerResponder[armreservations.ReservationClientListRevisionsResponse]]
	beginMerge            *tracker[azfake.PollerResponder[armreservations.ReservationClientMergeResponse]]
	beginSplit            *tracker[azfake.PollerResponder[armreservations.ReservationClientSplitResponse]]
	beginUpdate           *tracker[azfake.PollerResponder[armreservations.ReservationClientUpdateResponse]]
}

// Do implements the policy.Transporter interface for ReservationServerTransport.
func (r *ReservationServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "ReservationClient.Archive":
		resp, err = r.dispatchArchive(req)
	case "ReservationClient.BeginAvailableScopes":
		resp, err = r.dispatchBeginAvailableScopes(req)
	case "ReservationClient.Get":
		resp, err = r.dispatchGet(req)
	case "ReservationClient.NewListPager":
		resp, err = r.dispatchNewListPager(req)
	case "ReservationClient.NewListAllPager":
		resp, err = r.dispatchNewListAllPager(req)
	case "ReservationClient.NewListRevisionsPager":
		resp, err = r.dispatchNewListRevisionsPager(req)
	case "ReservationClient.BeginMerge":
		resp, err = r.dispatchBeginMerge(req)
	case "ReservationClient.BeginSplit":
		resp, err = r.dispatchBeginSplit(req)
	case "ReservationClient.Unarchive":
		resp, err = r.dispatchUnarchive(req)
	case "ReservationClient.BeginUpdate":
		resp, err = r.dispatchBeginUpdate(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *ReservationServerTransport) dispatchArchive(req *http.Request) (*http.Response, error) {
	if r.srv.Archive == nil {
		return nil, &nonRetriableError{errors.New("fake for method Archive not implemented")}
	}
	const regexStr = `/providers/Microsoft\.Capacity/reservationOrders/(?P<reservationOrderId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/reservations/(?P<reservationId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/archive`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	reservationOrderIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("reservationOrderId")])
	if err != nil {
		return nil, err
	}
	reservationIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("reservationId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := r.srv.Archive(req.Context(), reservationOrderIDParam, reservationIDParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.NewResponse(respContent, req, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *ReservationServerTransport) dispatchBeginAvailableScopes(req *http.Request) (*http.Response, error) {
	if r.srv.BeginAvailableScopes == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginAvailableScopes not implemented")}
	}
	beginAvailableScopes := r.beginAvailableScopes.get(req)
	if beginAvailableScopes == nil {
		const regexStr = `/providers/Microsoft\.Capacity/reservationOrders/(?P<reservationOrderId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/reservations/(?P<reservationId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/availableScopes`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armreservations.AvailableScopeRequest](req)
		if err != nil {
			return nil, err
		}
		reservationOrderIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("reservationOrderId")])
		if err != nil {
			return nil, err
		}
		reservationIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("reservationId")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := r.srv.BeginAvailableScopes(req.Context(), reservationOrderIDParam, reservationIDParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginAvailableScopes = &respr
		r.beginAvailableScopes.add(req, beginAvailableScopes)
	}

	resp, err := server.PollerResponderNext(beginAvailableScopes, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		r.beginAvailableScopes.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginAvailableScopes) {
		r.beginAvailableScopes.remove(req)
	}

	return resp, nil
}

func (r *ReservationServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if r.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/providers/Microsoft\.Capacity/reservationOrders/(?P<reservationOrderId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/reservations/(?P<reservationId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	qp := req.URL.Query()
	reservationOrderIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("reservationOrderId")])
	if err != nil {
		return nil, err
	}
	reservationIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("reservationId")])
	if err != nil {
		return nil, err
	}
	expandUnescaped, err := url.QueryUnescape(qp.Get("$expand"))
	if err != nil {
		return nil, err
	}
	expandParam := getOptional(expandUnescaped)
	var options *armreservations.ReservationClientGetOptions
	if expandParam != nil {
		options = &armreservations.ReservationClientGetOptions{
			Expand: expandParam,
		}
	}
	respr, errRespr := r.srv.Get(req.Context(), reservationOrderIDParam, reservationIDParam, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ReservationResponse, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *ReservationServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if r.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := r.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/providers/Microsoft\.Capacity/reservationOrders/(?P<reservationOrderId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/reservations`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		reservationOrderIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("reservationOrderId")])
		if err != nil {
			return nil, err
		}
		resp := r.srv.NewListPager(reservationOrderIDParam, nil)
		newListPager = &resp
		r.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armreservations.ReservationClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		r.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		r.newListPager.remove(req)
	}
	return resp, nil
}

func (r *ReservationServerTransport) dispatchNewListAllPager(req *http.Request) (*http.Response, error) {
	if r.srv.NewListAllPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListAllPager not implemented")}
	}
	newListAllPager := r.newListAllPager.get(req)
	if newListAllPager == nil {
		qp := req.URL.Query()
		filterUnescaped, err := url.QueryUnescape(qp.Get("$filter"))
		if err != nil {
			return nil, err
		}
		filterParam := getOptional(filterUnescaped)
		orderbyUnescaped, err := url.QueryUnescape(qp.Get("$orderby"))
		if err != nil {
			return nil, err
		}
		orderbyParam := getOptional(orderbyUnescaped)
		refreshSummaryUnescaped, err := url.QueryUnescape(qp.Get("refreshSummary"))
		if err != nil {
			return nil, err
		}
		refreshSummaryParam := getOptional(refreshSummaryUnescaped)
		skiptokenUnescaped, err := url.QueryUnescape(qp.Get("$skiptoken"))
		if err != nil {
			return nil, err
		}
		skiptokenParam, err := parseOptional(skiptokenUnescaped, func(v string) (float32, error) {
			p, parseErr := strconv.ParseFloat(v, 32)
			if parseErr != nil {
				return 0, parseErr
			}
			return float32(p), nil
		})
		if err != nil {
			return nil, err
		}
		selectedStateUnescaped, err := url.QueryUnescape(qp.Get("selectedState"))
		if err != nil {
			return nil, err
		}
		selectedStateParam := getOptional(selectedStateUnescaped)
		takeUnescaped, err := url.QueryUnescape(qp.Get("take"))
		if err != nil {
			return nil, err
		}
		takeParam, err := parseOptional(takeUnescaped, func(v string) (float32, error) {
			p, parseErr := strconv.ParseFloat(v, 32)
			if parseErr != nil {
				return 0, parseErr
			}
			return float32(p), nil
		})
		if err != nil {
			return nil, err
		}
		var options *armreservations.ReservationClientListAllOptions
		if filterParam != nil || orderbyParam != nil || refreshSummaryParam != nil || skiptokenParam != nil || selectedStateParam != nil || takeParam != nil {
			options = &armreservations.ReservationClientListAllOptions{
				Filter:         filterParam,
				Orderby:        orderbyParam,
				RefreshSummary: refreshSummaryParam,
				Skiptoken:      skiptokenParam,
				SelectedState:  selectedStateParam,
				Take:           takeParam,
			}
		}
		resp := r.srv.NewListAllPager(options)
		newListAllPager = &resp
		r.newListAllPager.add(req, newListAllPager)
		server.PagerResponderInjectNextLinks(newListAllPager, req, func(page *armreservations.ReservationClientListAllResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListAllPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		r.newListAllPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListAllPager) {
		r.newListAllPager.remove(req)
	}
	return resp, nil
}

func (r *ReservationServerTransport) dispatchNewListRevisionsPager(req *http.Request) (*http.Response, error) {
	if r.srv.NewListRevisionsPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListRevisionsPager not implemented")}
	}
	newListRevisionsPager := r.newListRevisionsPager.get(req)
	if newListRevisionsPager == nil {
		const regexStr = `/providers/Microsoft\.Capacity/reservationOrders/(?P<reservationOrderId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/reservations/(?P<reservationId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/revisions`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		reservationOrderIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("reservationOrderId")])
		if err != nil {
			return nil, err
		}
		reservationIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("reservationId")])
		if err != nil {
			return nil, err
		}
		resp := r.srv.NewListRevisionsPager(reservationOrderIDParam, reservationIDParam, nil)
		newListRevisionsPager = &resp
		r.newListRevisionsPager.add(req, newListRevisionsPager)
		server.PagerResponderInjectNextLinks(newListRevisionsPager, req, func(page *armreservations.ReservationClientListRevisionsResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListRevisionsPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		r.newListRevisionsPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListRevisionsPager) {
		r.newListRevisionsPager.remove(req)
	}
	return resp, nil
}

func (r *ReservationServerTransport) dispatchBeginMerge(req *http.Request) (*http.Response, error) {
	if r.srv.BeginMerge == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginMerge not implemented")}
	}
	beginMerge := r.beginMerge.get(req)
	if beginMerge == nil {
		const regexStr = `/providers/Microsoft\.Capacity/reservationOrders/(?P<reservationOrderId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/merge`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armreservations.MergeRequest](req)
		if err != nil {
			return nil, err
		}
		reservationOrderIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("reservationOrderId")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := r.srv.BeginMerge(req.Context(), reservationOrderIDParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginMerge = &respr
		r.beginMerge.add(req, beginMerge)
	}

	resp, err := server.PollerResponderNext(beginMerge, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		r.beginMerge.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginMerge) {
		r.beginMerge.remove(req)
	}

	return resp, nil
}

func (r *ReservationServerTransport) dispatchBeginSplit(req *http.Request) (*http.Response, error) {
	if r.srv.BeginSplit == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginSplit not implemented")}
	}
	beginSplit := r.beginSplit.get(req)
	if beginSplit == nil {
		const regexStr = `/providers/Microsoft\.Capacity/reservationOrders/(?P<reservationOrderId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/split`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armreservations.SplitRequest](req)
		if err != nil {
			return nil, err
		}
		reservationOrderIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("reservationOrderId")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := r.srv.BeginSplit(req.Context(), reservationOrderIDParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginSplit = &respr
		r.beginSplit.add(req, beginSplit)
	}

	resp, err := server.PollerResponderNext(beginSplit, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		r.beginSplit.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginSplit) {
		r.beginSplit.remove(req)
	}

	return resp, nil
}

func (r *ReservationServerTransport) dispatchUnarchive(req *http.Request) (*http.Response, error) {
	if r.srv.Unarchive == nil {
		return nil, &nonRetriableError{errors.New("fake for method Unarchive not implemented")}
	}
	const regexStr = `/providers/Microsoft\.Capacity/reservationOrders/(?P<reservationOrderId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/reservations/(?P<reservationId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/unarchive`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	reservationOrderIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("reservationOrderId")])
	if err != nil {
		return nil, err
	}
	reservationIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("reservationId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := r.srv.Unarchive(req.Context(), reservationOrderIDParam, reservationIDParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.NewResponse(respContent, req, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *ReservationServerTransport) dispatchBeginUpdate(req *http.Request) (*http.Response, error) {
	if r.srv.BeginUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdate not implemented")}
	}
	beginUpdate := r.beginUpdate.get(req)
	if beginUpdate == nil {
		const regexStr = `/providers/Microsoft\.Capacity/reservationOrders/(?P<reservationOrderId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/reservations/(?P<reservationId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armreservations.Patch](req)
		if err != nil {
			return nil, err
		}
		reservationOrderIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("reservationOrderId")])
		if err != nil {
			return nil, err
		}
		reservationIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("reservationId")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := r.srv.BeginUpdate(req.Context(), reservationOrderIDParam, reservationIDParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginUpdate = &respr
		r.beginUpdate.add(req, beginUpdate)
	}

	resp, err := server.PollerResponderNext(beginUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		r.beginUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdate) {
		r.beginUpdate.remove(req)
	}

	return resp, nil
}
