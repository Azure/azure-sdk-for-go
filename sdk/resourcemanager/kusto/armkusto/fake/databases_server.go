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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/kusto/armkusto/v2"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

// DatabasesServer is a fake server for instances of the armkusto.DatabasesClient type.
type DatabasesServer struct {
	// AddPrincipals is the fake for method DatabasesClient.AddPrincipals
	// HTTP status codes to indicate success: http.StatusOK
	AddPrincipals func(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, databasePrincipalsToAdd armkusto.DatabasePrincipalListRequest, options *armkusto.DatabasesClientAddPrincipalsOptions) (resp azfake.Responder[armkusto.DatabasesClientAddPrincipalsResponse], errResp azfake.ErrorResponder)

	// CheckNameAvailability is the fake for method DatabasesClient.CheckNameAvailability
	// HTTP status codes to indicate success: http.StatusOK
	CheckNameAvailability func(ctx context.Context, resourceGroupName string, clusterName string, resourceName armkusto.CheckNameRequest, options *armkusto.DatabasesClientCheckNameAvailabilityOptions) (resp azfake.Responder[armkusto.DatabasesClientCheckNameAvailabilityResponse], errResp azfake.ErrorResponder)

	// BeginCreateOrUpdate is the fake for method DatabasesClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated, http.StatusAccepted
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, parameters armkusto.DatabaseClassification, options *armkusto.DatabasesClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armkusto.DatabasesClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method DatabasesClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, options *armkusto.DatabasesClientBeginDeleteOptions) (resp azfake.PollerResponder[armkusto.DatabasesClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method DatabasesClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, options *armkusto.DatabasesClientGetOptions) (resp azfake.Responder[armkusto.DatabasesClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByClusterPager is the fake for method DatabasesClient.NewListByClusterPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByClusterPager func(resourceGroupName string, clusterName string, options *armkusto.DatabasesClientListByClusterOptions) (resp azfake.PagerResponder[armkusto.DatabasesClientListByClusterResponse])

	// NewListPrincipalsPager is the fake for method DatabasesClient.NewListPrincipalsPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPrincipalsPager func(resourceGroupName string, clusterName string, databaseName string, options *armkusto.DatabasesClientListPrincipalsOptions) (resp azfake.PagerResponder[armkusto.DatabasesClientListPrincipalsResponse])

	// RemovePrincipals is the fake for method DatabasesClient.RemovePrincipals
	// HTTP status codes to indicate success: http.StatusOK
	RemovePrincipals func(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, databasePrincipalsToRemove armkusto.DatabasePrincipalListRequest, options *armkusto.DatabasesClientRemovePrincipalsOptions) (resp azfake.Responder[armkusto.DatabasesClientRemovePrincipalsResponse], errResp azfake.ErrorResponder)

	// BeginUpdate is the fake for method DatabasesClient.BeginUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated, http.StatusAccepted
	BeginUpdate func(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, parameters armkusto.DatabaseClassification, options *armkusto.DatabasesClientBeginUpdateOptions) (resp azfake.PollerResponder[armkusto.DatabasesClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewDatabasesServerTransport creates a new instance of DatabasesServerTransport with the provided implementation.
// The returned DatabasesServerTransport instance is connected to an instance of armkusto.DatabasesClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewDatabasesServerTransport(srv *DatabasesServer) *DatabasesServerTransport {
	return &DatabasesServerTransport{
		srv:                    srv,
		beginCreateOrUpdate:    newTracker[azfake.PollerResponder[armkusto.DatabasesClientCreateOrUpdateResponse]](),
		beginDelete:            newTracker[azfake.PollerResponder[armkusto.DatabasesClientDeleteResponse]](),
		newListByClusterPager:  newTracker[azfake.PagerResponder[armkusto.DatabasesClientListByClusterResponse]](),
		newListPrincipalsPager: newTracker[azfake.PagerResponder[armkusto.DatabasesClientListPrincipalsResponse]](),
		beginUpdate:            newTracker[azfake.PollerResponder[armkusto.DatabasesClientUpdateResponse]](),
	}
}

// DatabasesServerTransport connects instances of armkusto.DatabasesClient to instances of DatabasesServer.
// Don't use this type directly, use NewDatabasesServerTransport instead.
type DatabasesServerTransport struct {
	srv                    *DatabasesServer
	beginCreateOrUpdate    *tracker[azfake.PollerResponder[armkusto.DatabasesClientCreateOrUpdateResponse]]
	beginDelete            *tracker[azfake.PollerResponder[armkusto.DatabasesClientDeleteResponse]]
	newListByClusterPager  *tracker[azfake.PagerResponder[armkusto.DatabasesClientListByClusterResponse]]
	newListPrincipalsPager *tracker[azfake.PagerResponder[armkusto.DatabasesClientListPrincipalsResponse]]
	beginUpdate            *tracker[azfake.PollerResponder[armkusto.DatabasesClientUpdateResponse]]
}

// Do implements the policy.Transporter interface for DatabasesServerTransport.
func (d *DatabasesServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "DatabasesClient.AddPrincipals":
		resp, err = d.dispatchAddPrincipals(req)
	case "DatabasesClient.CheckNameAvailability":
		resp, err = d.dispatchCheckNameAvailability(req)
	case "DatabasesClient.BeginCreateOrUpdate":
		resp, err = d.dispatchBeginCreateOrUpdate(req)
	case "DatabasesClient.BeginDelete":
		resp, err = d.dispatchBeginDelete(req)
	case "DatabasesClient.Get":
		resp, err = d.dispatchGet(req)
	case "DatabasesClient.NewListByClusterPager":
		resp, err = d.dispatchNewListByClusterPager(req)
	case "DatabasesClient.NewListPrincipalsPager":
		resp, err = d.dispatchNewListPrincipalsPager(req)
	case "DatabasesClient.RemovePrincipals":
		resp, err = d.dispatchRemovePrincipals(req)
	case "DatabasesClient.BeginUpdate":
		resp, err = d.dispatchBeginUpdate(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (d *DatabasesServerTransport) dispatchAddPrincipals(req *http.Request) (*http.Response, error) {
	if d.srv.AddPrincipals == nil {
		return nil, &nonRetriableError{errors.New("fake for method AddPrincipals not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Kusto/clusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/databases/(?P<databaseName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/addPrincipals`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armkusto.DatabasePrincipalListRequest](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	clusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("clusterName")])
	if err != nil {
		return nil, err
	}
	databaseNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("databaseName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := d.srv.AddPrincipals(req.Context(), resourceGroupNameParam, clusterNameParam, databaseNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).DatabasePrincipalListResult, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DatabasesServerTransport) dispatchCheckNameAvailability(req *http.Request) (*http.Response, error) {
	if d.srv.CheckNameAvailability == nil {
		return nil, &nonRetriableError{errors.New("fake for method CheckNameAvailability not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Kusto/clusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/checkNameAvailability`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armkusto.CheckNameRequest](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	clusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("clusterName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := d.srv.CheckNameAvailability(req.Context(), resourceGroupNameParam, clusterNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).CheckNameResult, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DatabasesServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if d.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := d.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Kusto/clusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/databases/(?P<databaseName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		raw, err := readRequestBody(req)
		if err != nil {
			return nil, err
		}
		body, err := unmarshalDatabaseClassification(raw)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		clusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("clusterName")])
		if err != nil {
			return nil, err
		}
		databaseNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("databaseName")])
		if err != nil {
			return nil, err
		}
		callerRoleUnescaped, err := url.QueryUnescape(qp.Get("callerRole"))
		if err != nil {
			return nil, err
		}
		callerRoleParam := getOptional(armkusto.CallerRole(callerRoleUnescaped))
		var options *armkusto.DatabasesClientBeginCreateOrUpdateOptions
		if callerRoleParam != nil {
			options = &armkusto.DatabasesClientBeginCreateOrUpdateOptions{
				CallerRole: callerRoleParam,
			}
		}
		respr, errRespr := d.srv.BeginCreateOrUpdate(req.Context(), resourceGroupNameParam, clusterNameParam, databaseNameParam, body, options)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateOrUpdate = &respr
		d.beginCreateOrUpdate.add(req, beginCreateOrUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated, http.StatusAccepted}, resp.StatusCode) {
		d.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		d.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (d *DatabasesServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if d.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := d.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Kusto/clusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/databases/(?P<databaseName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		clusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("clusterName")])
		if err != nil {
			return nil, err
		}
		databaseNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("databaseName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := d.srv.BeginDelete(req.Context(), resourceGroupNameParam, clusterNameParam, databaseNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		d.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		d.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		d.beginDelete.remove(req)
	}

	return resp, nil
}

func (d *DatabasesServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if d.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Kusto/clusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/databases/(?P<databaseName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	clusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("clusterName")])
	if err != nil {
		return nil, err
	}
	databaseNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("databaseName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := d.srv.Get(req.Context(), resourceGroupNameParam, clusterNameParam, databaseNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).DatabaseClassification, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DatabasesServerTransport) dispatchNewListByClusterPager(req *http.Request) (*http.Response, error) {
	if d.srv.NewListByClusterPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByClusterPager not implemented")}
	}
	newListByClusterPager := d.newListByClusterPager.get(req)
	if newListByClusterPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Kusto/clusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/databases`
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
		clusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("clusterName")])
		if err != nil {
			return nil, err
		}
		topUnescaped, err := url.QueryUnescape(qp.Get("$top"))
		if err != nil {
			return nil, err
		}
		topParam, err := parseOptional(topUnescaped, func(v string) (int32, error) {
			p, parseErr := strconv.ParseInt(v, 10, 32)
			if parseErr != nil {
				return 0, parseErr
			}
			return int32(p), nil
		})
		if err != nil {
			return nil, err
		}
		skiptokenUnescaped, err := url.QueryUnescape(qp.Get("$skiptoken"))
		if err != nil {
			return nil, err
		}
		skiptokenParam := getOptional(skiptokenUnescaped)
		var options *armkusto.DatabasesClientListByClusterOptions
		if topParam != nil || skiptokenParam != nil {
			options = &armkusto.DatabasesClientListByClusterOptions{
				Top:       topParam,
				Skiptoken: skiptokenParam,
			}
		}
		resp := d.srv.NewListByClusterPager(resourceGroupNameParam, clusterNameParam, options)
		newListByClusterPager = &resp
		d.newListByClusterPager.add(req, newListByClusterPager)
		server.PagerResponderInjectNextLinks(newListByClusterPager, req, func(page *armkusto.DatabasesClientListByClusterResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByClusterPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		d.newListByClusterPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByClusterPager) {
		d.newListByClusterPager.remove(req)
	}
	return resp, nil
}

func (d *DatabasesServerTransport) dispatchNewListPrincipalsPager(req *http.Request) (*http.Response, error) {
	if d.srv.NewListPrincipalsPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPrincipalsPager not implemented")}
	}
	newListPrincipalsPager := d.newListPrincipalsPager.get(req)
	if newListPrincipalsPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Kusto/clusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/databases/(?P<databaseName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/listPrincipals`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		clusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("clusterName")])
		if err != nil {
			return nil, err
		}
		databaseNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("databaseName")])
		if err != nil {
			return nil, err
		}
		resp := d.srv.NewListPrincipalsPager(resourceGroupNameParam, clusterNameParam, databaseNameParam, nil)
		newListPrincipalsPager = &resp
		d.newListPrincipalsPager.add(req, newListPrincipalsPager)
	}
	resp, err := server.PagerResponderNext(newListPrincipalsPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		d.newListPrincipalsPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPrincipalsPager) {
		d.newListPrincipalsPager.remove(req)
	}
	return resp, nil
}

func (d *DatabasesServerTransport) dispatchRemovePrincipals(req *http.Request) (*http.Response, error) {
	if d.srv.RemovePrincipals == nil {
		return nil, &nonRetriableError{errors.New("fake for method RemovePrincipals not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Kusto/clusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/databases/(?P<databaseName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/removePrincipals`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armkusto.DatabasePrincipalListRequest](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	clusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("clusterName")])
	if err != nil {
		return nil, err
	}
	databaseNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("databaseName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := d.srv.RemovePrincipals(req.Context(), resourceGroupNameParam, clusterNameParam, databaseNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).DatabasePrincipalListResult, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DatabasesServerTransport) dispatchBeginUpdate(req *http.Request) (*http.Response, error) {
	if d.srv.BeginUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdate not implemented")}
	}
	beginUpdate := d.beginUpdate.get(req)
	if beginUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Kusto/clusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/databases/(?P<databaseName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		raw, err := readRequestBody(req)
		if err != nil {
			return nil, err
		}
		body, err := unmarshalDatabaseClassification(raw)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		clusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("clusterName")])
		if err != nil {
			return nil, err
		}
		databaseNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("databaseName")])
		if err != nil {
			return nil, err
		}
		callerRoleUnescaped, err := url.QueryUnescape(qp.Get("callerRole"))
		if err != nil {
			return nil, err
		}
		callerRoleParam := getOptional(armkusto.CallerRole(callerRoleUnescaped))
		var options *armkusto.DatabasesClientBeginUpdateOptions
		if callerRoleParam != nil {
			options = &armkusto.DatabasesClientBeginUpdateOptions{
				CallerRole: callerRoleParam,
			}
		}
		respr, errRespr := d.srv.BeginUpdate(req.Context(), resourceGroupNameParam, clusterNameParam, databaseNameParam, body, options)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginUpdate = &respr
		d.beginUpdate.add(req, beginUpdate)
	}

	resp, err := server.PollerResponderNext(beginUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated, http.StatusAccepted}, resp.StatusCode) {
		d.beginUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdate) {
		d.beginUpdate.remove(req)
	}

	return resp, nil
}
