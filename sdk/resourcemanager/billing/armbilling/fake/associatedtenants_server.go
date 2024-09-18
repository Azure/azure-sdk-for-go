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

// AssociatedTenantsServer is a fake server for instances of the armbilling.AssociatedTenantsClient type.
type AssociatedTenantsServer struct {
	// BeginCreateOrUpdate is the fake for method AssociatedTenantsClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreateOrUpdate func(ctx context.Context, billingAccountName string, associatedTenantName string, parameters armbilling.AssociatedTenant, options *armbilling.AssociatedTenantsClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armbilling.AssociatedTenantsClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method AssociatedTenantsClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, billingAccountName string, associatedTenantName string, options *armbilling.AssociatedTenantsClientBeginDeleteOptions) (resp azfake.PollerResponder[armbilling.AssociatedTenantsClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method AssociatedTenantsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, billingAccountName string, associatedTenantName string, options *armbilling.AssociatedTenantsClientGetOptions) (resp azfake.Responder[armbilling.AssociatedTenantsClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByBillingAccountPager is the fake for method AssociatedTenantsClient.NewListByBillingAccountPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByBillingAccountPager func(billingAccountName string, options *armbilling.AssociatedTenantsClientListByBillingAccountOptions) (resp azfake.PagerResponder[armbilling.AssociatedTenantsClientListByBillingAccountResponse])
}

// NewAssociatedTenantsServerTransport creates a new instance of AssociatedTenantsServerTransport with the provided implementation.
// The returned AssociatedTenantsServerTransport instance is connected to an instance of armbilling.AssociatedTenantsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewAssociatedTenantsServerTransport(srv *AssociatedTenantsServer) *AssociatedTenantsServerTransport {
	return &AssociatedTenantsServerTransport{
		srv:                          srv,
		beginCreateOrUpdate:          newTracker[azfake.PollerResponder[armbilling.AssociatedTenantsClientCreateOrUpdateResponse]](),
		beginDelete:                  newTracker[azfake.PollerResponder[armbilling.AssociatedTenantsClientDeleteResponse]](),
		newListByBillingAccountPager: newTracker[azfake.PagerResponder[armbilling.AssociatedTenantsClientListByBillingAccountResponse]](),
	}
}

// AssociatedTenantsServerTransport connects instances of armbilling.AssociatedTenantsClient to instances of AssociatedTenantsServer.
// Don't use this type directly, use NewAssociatedTenantsServerTransport instead.
type AssociatedTenantsServerTransport struct {
	srv                          *AssociatedTenantsServer
	beginCreateOrUpdate          *tracker[azfake.PollerResponder[armbilling.AssociatedTenantsClientCreateOrUpdateResponse]]
	beginDelete                  *tracker[azfake.PollerResponder[armbilling.AssociatedTenantsClientDeleteResponse]]
	newListByBillingAccountPager *tracker[azfake.PagerResponder[armbilling.AssociatedTenantsClientListByBillingAccountResponse]]
}

// Do implements the policy.Transporter interface for AssociatedTenantsServerTransport.
func (a *AssociatedTenantsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "AssociatedTenantsClient.BeginCreateOrUpdate":
		resp, err = a.dispatchBeginCreateOrUpdate(req)
	case "AssociatedTenantsClient.BeginDelete":
		resp, err = a.dispatchBeginDelete(req)
	case "AssociatedTenantsClient.Get":
		resp, err = a.dispatchGet(req)
	case "AssociatedTenantsClient.NewListByBillingAccountPager":
		resp, err = a.dispatchNewListByBillingAccountPager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *AssociatedTenantsServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if a.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := a.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/providers/Microsoft\.Billing/billingAccounts/(?P<billingAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/associatedTenants/(?P<associatedTenantName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armbilling.AssociatedTenant](req)
		if err != nil {
			return nil, err
		}
		billingAccountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("billingAccountName")])
		if err != nil {
			return nil, err
		}
		associatedTenantNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("associatedTenantName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := a.srv.BeginCreateOrUpdate(req.Context(), billingAccountNameParam, associatedTenantNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateOrUpdate = &respr
		a.beginCreateOrUpdate.add(req, beginCreateOrUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		a.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		a.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (a *AssociatedTenantsServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if a.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := a.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/providers/Microsoft\.Billing/billingAccounts/(?P<billingAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/associatedTenants/(?P<associatedTenantName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		billingAccountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("billingAccountName")])
		if err != nil {
			return nil, err
		}
		associatedTenantNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("associatedTenantName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := a.srv.BeginDelete(req.Context(), billingAccountNameParam, associatedTenantNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		a.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		a.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		a.beginDelete.remove(req)
	}

	return resp, nil
}

func (a *AssociatedTenantsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if a.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/providers/Microsoft\.Billing/billingAccounts/(?P<billingAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/associatedTenants/(?P<associatedTenantName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	billingAccountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("billingAccountName")])
	if err != nil {
		return nil, err
	}
	associatedTenantNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("associatedTenantName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := a.srv.Get(req.Context(), billingAccountNameParam, associatedTenantNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).AssociatedTenant, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *AssociatedTenantsServerTransport) dispatchNewListByBillingAccountPager(req *http.Request) (*http.Response, error) {
	if a.srv.NewListByBillingAccountPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByBillingAccountPager not implemented")}
	}
	newListByBillingAccountPager := a.newListByBillingAccountPager.get(req)
	if newListByBillingAccountPager == nil {
		const regexStr = `/providers/Microsoft\.Billing/billingAccounts/(?P<billingAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/associatedTenants`
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
		includeRevokedUnescaped, err := url.QueryUnescape(qp.Get("includeRevoked"))
		if err != nil {
			return nil, err
		}
		includeRevokedParam, err := parseOptional(includeRevokedUnescaped, strconv.ParseBool)
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
		topUnescaped, err := url.QueryUnescape(qp.Get("top"))
		if err != nil {
			return nil, err
		}
		topParam, err := parseOptional(topUnescaped, func(v string) (int64, error) {
			p, parseErr := strconv.ParseInt(v, 10, 64)
			if parseErr != nil {
				return 0, parseErr
			}
			return p, nil
		})
		if err != nil {
			return nil, err
		}
		skipUnescaped, err := url.QueryUnescape(qp.Get("skip"))
		if err != nil {
			return nil, err
		}
		skipParam, err := parseOptional(skipUnescaped, func(v string) (int64, error) {
			p, parseErr := strconv.ParseInt(v, 10, 64)
			if parseErr != nil {
				return 0, parseErr
			}
			return p, nil
		})
		if err != nil {
			return nil, err
		}
		countUnescaped, err := url.QueryUnescape(qp.Get("count"))
		if err != nil {
			return nil, err
		}
		countParam, err := parseOptional(countUnescaped, strconv.ParseBool)
		if err != nil {
			return nil, err
		}
		searchUnescaped, err := url.QueryUnescape(qp.Get("search"))
		if err != nil {
			return nil, err
		}
		searchParam := getOptional(searchUnescaped)
		var options *armbilling.AssociatedTenantsClientListByBillingAccountOptions
		if includeRevokedParam != nil || filterParam != nil || orderByParam != nil || topParam != nil || skipParam != nil || countParam != nil || searchParam != nil {
			options = &armbilling.AssociatedTenantsClientListByBillingAccountOptions{
				IncludeRevoked: includeRevokedParam,
				Filter:         filterParam,
				OrderBy:        orderByParam,
				Top:            topParam,
				Skip:           skipParam,
				Count:          countParam,
				Search:         searchParam,
			}
		}
		resp := a.srv.NewListByBillingAccountPager(billingAccountNameParam, options)
		newListByBillingAccountPager = &resp
		a.newListByBillingAccountPager.add(req, newListByBillingAccountPager)
		server.PagerResponderInjectNextLinks(newListByBillingAccountPager, req, func(page *armbilling.AssociatedTenantsClientListByBillingAccountResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByBillingAccountPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		a.newListByBillingAccountPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByBillingAccountPager) {
		a.newListByBillingAccountPager.remove(req)
	}
	return resp, nil
}