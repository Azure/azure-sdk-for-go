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
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake/server"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/billing/armbilling"
)

// SavingsPlanOrdersServer is a fake server for instances of the armbilling.SavingsPlanOrdersClient type.
type SavingsPlanOrdersServer struct {
	// GetByBillingAccount is the fake for method SavingsPlanOrdersClient.GetByBillingAccount
	// HTTP status codes to indicate success: http.StatusOK
	GetByBillingAccount func(ctx context.Context, billingAccountName string, savingsPlanOrderID string, options *armbilling.SavingsPlanOrdersClientGetByBillingAccountOptions) (resp azfake.Responder[armbilling.SavingsPlanOrdersClientGetByBillingAccountResponse], errResp azfake.ErrorResponder)

	// NewListByBillingAccountPager is the fake for method SavingsPlanOrdersClient.NewListByBillingAccountPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByBillingAccountPager func(billingAccountName string, options *armbilling.SavingsPlanOrdersClientListByBillingAccountOptions) (resp azfake.PagerResponder[armbilling.SavingsPlanOrdersClientListByBillingAccountResponse])
}

// NewSavingsPlanOrdersServerTransport creates a new instance of SavingsPlanOrdersServerTransport with the provided implementation.
// The returned SavingsPlanOrdersServerTransport instance is connected to an instance of armbilling.SavingsPlanOrdersClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewSavingsPlanOrdersServerTransport(srv *SavingsPlanOrdersServer) *SavingsPlanOrdersServerTransport {
	return &SavingsPlanOrdersServerTransport{
		srv:                          srv,
		newListByBillingAccountPager: newTracker[azfake.PagerResponder[armbilling.SavingsPlanOrdersClientListByBillingAccountResponse]](),
	}
}

// SavingsPlanOrdersServerTransport connects instances of armbilling.SavingsPlanOrdersClient to instances of SavingsPlanOrdersServer.
// Don't use this type directly, use NewSavingsPlanOrdersServerTransport instead.
type SavingsPlanOrdersServerTransport struct {
	srv                          *SavingsPlanOrdersServer
	newListByBillingAccountPager *tracker[azfake.PagerResponder[armbilling.SavingsPlanOrdersClientListByBillingAccountResponse]]
}

// Do implements the policy.Transporter interface for SavingsPlanOrdersServerTransport.
func (s *SavingsPlanOrdersServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "SavingsPlanOrdersClient.GetByBillingAccount":
		resp, err = s.dispatchGetByBillingAccount(req)
	case "SavingsPlanOrdersClient.NewListByBillingAccountPager":
		resp, err = s.dispatchNewListByBillingAccountPager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *SavingsPlanOrdersServerTransport) dispatchGetByBillingAccount(req *http.Request) (*http.Response, error) {
	if s.srv.GetByBillingAccount == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetByBillingAccount not implemented")}
	}
	const regexStr = `/providers/Microsoft\.Billing/billingAccounts/(?P<billingAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/savingsPlanOrders/(?P<savingsPlanOrderId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	qp := req.URL.Query()
	billingAccountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("billingAccountName")])
	if err != nil {
		return nil, err
	}
	savingsPlanOrderIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("savingsPlanOrderId")])
	if err != nil {
		return nil, err
	}
	expandUnescaped, err := url.QueryUnescape(qp.Get("expand"))
	if err != nil {
		return nil, err
	}
	expandParam := getOptional(expandUnescaped)
	var options *armbilling.SavingsPlanOrdersClientGetByBillingAccountOptions
	if expandParam != nil {
		options = &armbilling.SavingsPlanOrdersClientGetByBillingAccountOptions{
			Expand: expandParam,
		}
	}
	respr, errRespr := s.srv.GetByBillingAccount(req.Context(), billingAccountNameParam, savingsPlanOrderIDParam, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).SavingsPlanOrderModel, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *SavingsPlanOrdersServerTransport) dispatchNewListByBillingAccountPager(req *http.Request) (*http.Response, error) {
	if s.srv.NewListByBillingAccountPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByBillingAccountPager not implemented")}
	}
	newListByBillingAccountPager := s.newListByBillingAccountPager.get(req)
	if newListByBillingAccountPager == nil {
		const regexStr = `/providers/Microsoft\.Billing/billingAccounts/(?P<billingAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/savingsPlanOrders`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		billingAccountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("billingAccountName")])
		if err != nil {
			return nil, err
		}
		filterUnescaped, err := url.QueryUnescape(qp.Get("filter"))
		if err != nil {
			return nil, err
		}
		filterParam := getOptional(filterUnescaped)
		orderByUnescaped, err := url.QueryUnescape(qp.Get("orderBy"))
		if err != nil {
			return nil, err
		}
		orderByParam := getOptional(orderByUnescaped)
		skiptokenUnescaped, err := url.QueryUnescape(qp.Get("skiptoken"))
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
		var options *armbilling.SavingsPlanOrdersClientListByBillingAccountOptions
		if filterParam != nil || orderByParam != nil || skiptokenParam != nil {
			options = &armbilling.SavingsPlanOrdersClientListByBillingAccountOptions{
				Filter:    filterParam,
				OrderBy:   orderByParam,
				Skiptoken: skiptokenParam,
			}
		}
		resp := s.srv.NewListByBillingAccountPager(billingAccountNameParam, options)
		newListByBillingAccountPager = &resp
		s.newListByBillingAccountPager.add(req, newListByBillingAccountPager)
		server.PagerResponderInjectNextLinks(newListByBillingAccountPager, req, func(page *armbilling.SavingsPlanOrdersClientListByBillingAccountResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByBillingAccountPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		s.newListByBillingAccountPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByBillingAccountPager) {
		s.newListByBillingAccountPager.remove(req)
	}
	return resp, nil
}
