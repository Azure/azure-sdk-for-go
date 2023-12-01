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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hybridcompute/armhybridcompute/v2"
	"net/http"
	"net/url"
	"regexp"
)

// LicensesServer is a fake server for instances of the armhybridcompute.LicensesClient type.
type LicensesServer struct {
	// BeginCreateOrUpdate is the fake for method LicensesClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, licenseName string, parameters armhybridcompute.License, options *armhybridcompute.LicensesClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armhybridcompute.LicensesClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method LicensesClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, licenseName string, options *armhybridcompute.LicensesClientBeginDeleteOptions) (resp azfake.PollerResponder[armhybridcompute.LicensesClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method LicensesClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, licenseName string, options *armhybridcompute.LicensesClientGetOptions) (resp azfake.Responder[armhybridcompute.LicensesClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByResourceGroupPager is the fake for method LicensesClient.NewListByResourceGroupPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByResourceGroupPager func(resourceGroupName string, options *armhybridcompute.LicensesClientListByResourceGroupOptions) (resp azfake.PagerResponder[armhybridcompute.LicensesClientListByResourceGroupResponse])

	// NewListBySubscriptionPager is the fake for method LicensesClient.NewListBySubscriptionPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListBySubscriptionPager func(options *armhybridcompute.LicensesClientListBySubscriptionOptions) (resp azfake.PagerResponder[armhybridcompute.LicensesClientListBySubscriptionResponse])

	// BeginUpdate is the fake for method LicensesClient.BeginUpdate
	// HTTP status codes to indicate success: http.StatusOK
	BeginUpdate func(ctx context.Context, resourceGroupName string, licenseName string, parameters armhybridcompute.LicenseUpdate, options *armhybridcompute.LicensesClientBeginUpdateOptions) (resp azfake.PollerResponder[armhybridcompute.LicensesClientUpdateResponse], errResp azfake.ErrorResponder)

	// BeginValidateLicense is the fake for method LicensesClient.BeginValidateLicense
	// HTTP status codes to indicate success: http.StatusOK
	BeginValidateLicense func(ctx context.Context, parameters armhybridcompute.License, options *armhybridcompute.LicensesClientBeginValidateLicenseOptions) (resp azfake.PollerResponder[armhybridcompute.LicensesClientValidateLicenseResponse], errResp azfake.ErrorResponder)
}

// NewLicensesServerTransport creates a new instance of LicensesServerTransport with the provided implementation.
// The returned LicensesServerTransport instance is connected to an instance of armhybridcompute.LicensesClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewLicensesServerTransport(srv *LicensesServer) *LicensesServerTransport {
	return &LicensesServerTransport{
		srv:                         srv,
		beginCreateOrUpdate:         newTracker[azfake.PollerResponder[armhybridcompute.LicensesClientCreateOrUpdateResponse]](),
		beginDelete:                 newTracker[azfake.PollerResponder[armhybridcompute.LicensesClientDeleteResponse]](),
		newListByResourceGroupPager: newTracker[azfake.PagerResponder[armhybridcompute.LicensesClientListByResourceGroupResponse]](),
		newListBySubscriptionPager:  newTracker[azfake.PagerResponder[armhybridcompute.LicensesClientListBySubscriptionResponse]](),
		beginUpdate:                 newTracker[azfake.PollerResponder[armhybridcompute.LicensesClientUpdateResponse]](),
		beginValidateLicense:        newTracker[azfake.PollerResponder[armhybridcompute.LicensesClientValidateLicenseResponse]](),
	}
}

// LicensesServerTransport connects instances of armhybridcompute.LicensesClient to instances of LicensesServer.
// Don't use this type directly, use NewLicensesServerTransport instead.
type LicensesServerTransport struct {
	srv                         *LicensesServer
	beginCreateOrUpdate         *tracker[azfake.PollerResponder[armhybridcompute.LicensesClientCreateOrUpdateResponse]]
	beginDelete                 *tracker[azfake.PollerResponder[armhybridcompute.LicensesClientDeleteResponse]]
	newListByResourceGroupPager *tracker[azfake.PagerResponder[armhybridcompute.LicensesClientListByResourceGroupResponse]]
	newListBySubscriptionPager  *tracker[azfake.PagerResponder[armhybridcompute.LicensesClientListBySubscriptionResponse]]
	beginUpdate                 *tracker[azfake.PollerResponder[armhybridcompute.LicensesClientUpdateResponse]]
	beginValidateLicense        *tracker[azfake.PollerResponder[armhybridcompute.LicensesClientValidateLicenseResponse]]
}

// Do implements the policy.Transporter interface for LicensesServerTransport.
func (l *LicensesServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "LicensesClient.BeginCreateOrUpdate":
		resp, err = l.dispatchBeginCreateOrUpdate(req)
	case "LicensesClient.BeginDelete":
		resp, err = l.dispatchBeginDelete(req)
	case "LicensesClient.Get":
		resp, err = l.dispatchGet(req)
	case "LicensesClient.NewListByResourceGroupPager":
		resp, err = l.dispatchNewListByResourceGroupPager(req)
	case "LicensesClient.NewListBySubscriptionPager":
		resp, err = l.dispatchNewListBySubscriptionPager(req)
	case "LicensesClient.BeginUpdate":
		resp, err = l.dispatchBeginUpdate(req)
	case "LicensesClient.BeginValidateLicense":
		resp, err = l.dispatchBeginValidateLicense(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (l *LicensesServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if l.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := l.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HybridCompute/licenses/(?P<licenseName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armhybridcompute.License](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		licenseNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("licenseName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := l.srv.BeginCreateOrUpdate(req.Context(), resourceGroupNameParam, licenseNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateOrUpdate = &respr
		l.beginCreateOrUpdate.add(req, beginCreateOrUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		l.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		l.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (l *LicensesServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if l.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := l.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HybridCompute/licenses/(?P<licenseName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		licenseNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("licenseName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := l.srv.BeginDelete(req.Context(), resourceGroupNameParam, licenseNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		l.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusNoContent}, resp.StatusCode) {
		l.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		l.beginDelete.remove(req)
	}

	return resp, nil
}

func (l *LicensesServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if l.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HybridCompute/licenses/(?P<licenseName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	licenseNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("licenseName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := l.srv.Get(req.Context(), resourceGroupNameParam, licenseNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).License, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (l *LicensesServerTransport) dispatchNewListByResourceGroupPager(req *http.Request) (*http.Response, error) {
	if l.srv.NewListByResourceGroupPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByResourceGroupPager not implemented")}
	}
	newListByResourceGroupPager := l.newListByResourceGroupPager.get(req)
	if newListByResourceGroupPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HybridCompute/licenses`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		resp := l.srv.NewListByResourceGroupPager(resourceGroupNameParam, nil)
		newListByResourceGroupPager = &resp
		l.newListByResourceGroupPager.add(req, newListByResourceGroupPager)
		server.PagerResponderInjectNextLinks(newListByResourceGroupPager, req, func(page *armhybridcompute.LicensesClientListByResourceGroupResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByResourceGroupPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		l.newListByResourceGroupPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByResourceGroupPager) {
		l.newListByResourceGroupPager.remove(req)
	}
	return resp, nil
}

func (l *LicensesServerTransport) dispatchNewListBySubscriptionPager(req *http.Request) (*http.Response, error) {
	if l.srv.NewListBySubscriptionPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListBySubscriptionPager not implemented")}
	}
	newListBySubscriptionPager := l.newListBySubscriptionPager.get(req)
	if newListBySubscriptionPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HybridCompute/licenses`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resp := l.srv.NewListBySubscriptionPager(nil)
		newListBySubscriptionPager = &resp
		l.newListBySubscriptionPager.add(req, newListBySubscriptionPager)
		server.PagerResponderInjectNextLinks(newListBySubscriptionPager, req, func(page *armhybridcompute.LicensesClientListBySubscriptionResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListBySubscriptionPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		l.newListBySubscriptionPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListBySubscriptionPager) {
		l.newListBySubscriptionPager.remove(req)
	}
	return resp, nil
}

func (l *LicensesServerTransport) dispatchBeginUpdate(req *http.Request) (*http.Response, error) {
	if l.srv.BeginUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdate not implemented")}
	}
	beginUpdate := l.beginUpdate.get(req)
	if beginUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HybridCompute/licenses/(?P<licenseName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armhybridcompute.LicenseUpdate](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		licenseNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("licenseName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := l.srv.BeginUpdate(req.Context(), resourceGroupNameParam, licenseNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginUpdate = &respr
		l.beginUpdate.add(req, beginUpdate)
	}

	resp, err := server.PollerResponderNext(beginUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		l.beginUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdate) {
		l.beginUpdate.remove(req)
	}

	return resp, nil
}

func (l *LicensesServerTransport) dispatchBeginValidateLicense(req *http.Request) (*http.Response, error) {
	if l.srv.BeginValidateLicense == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginValidateLicense not implemented")}
	}
	beginValidateLicense := l.beginValidateLicense.get(req)
	if beginValidateLicense == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HybridCompute/validateLicense`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armhybridcompute.License](req)
		if err != nil {
			return nil, err
		}
		respr, errRespr := l.srv.BeginValidateLicense(req.Context(), body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginValidateLicense = &respr
		l.beginValidateLicense.add(req, beginValidateLicense)
	}

	resp, err := server.PollerResponderNext(beginValidateLicense, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		l.beginValidateLicense.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginValidateLicense) {
		l.beginValidateLicense.remove(req)
	}

	return resp, nil
}