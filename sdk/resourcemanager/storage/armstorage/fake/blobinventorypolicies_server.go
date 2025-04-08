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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"net/http"
	"net/url"
	"regexp"
)

// BlobInventoryPoliciesServer is a fake server for instances of the armstorage.BlobInventoryPoliciesClient type.
type BlobInventoryPoliciesServer struct {
	// CreateOrUpdate is the fake for method BlobInventoryPoliciesClient.CreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK
	CreateOrUpdate func(ctx context.Context, resourceGroupName string, accountName string, blobInventoryPolicyName armstorage.BlobInventoryPolicyName, properties armstorage.BlobInventoryPolicy, options *armstorage.BlobInventoryPoliciesClientCreateOrUpdateOptions) (resp azfake.Responder[armstorage.BlobInventoryPoliciesClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// Delete is the fake for method BlobInventoryPoliciesClient.Delete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	Delete func(ctx context.Context, resourceGroupName string, accountName string, blobInventoryPolicyName armstorage.BlobInventoryPolicyName, options *armstorage.BlobInventoryPoliciesClientDeleteOptions) (resp azfake.Responder[armstorage.BlobInventoryPoliciesClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method BlobInventoryPoliciesClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, accountName string, blobInventoryPolicyName armstorage.BlobInventoryPolicyName, options *armstorage.BlobInventoryPoliciesClientGetOptions) (resp azfake.Responder[armstorage.BlobInventoryPoliciesClientGetResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method BlobInventoryPoliciesClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(resourceGroupName string, accountName string, options *armstorage.BlobInventoryPoliciesClientListOptions) (resp azfake.PagerResponder[armstorage.BlobInventoryPoliciesClientListResponse])
}

// NewBlobInventoryPoliciesServerTransport creates a new instance of BlobInventoryPoliciesServerTransport with the provided implementation.
// The returned BlobInventoryPoliciesServerTransport instance is connected to an instance of armstorage.BlobInventoryPoliciesClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewBlobInventoryPoliciesServerTransport(srv *BlobInventoryPoliciesServer) *BlobInventoryPoliciesServerTransport {
	return &BlobInventoryPoliciesServerTransport{
		srv:          srv,
		newListPager: newTracker[azfake.PagerResponder[armstorage.BlobInventoryPoliciesClientListResponse]](),
	}
}

// BlobInventoryPoliciesServerTransport connects instances of armstorage.BlobInventoryPoliciesClient to instances of BlobInventoryPoliciesServer.
// Don't use this type directly, use NewBlobInventoryPoliciesServerTransport instead.
type BlobInventoryPoliciesServerTransport struct {
	srv          *BlobInventoryPoliciesServer
	newListPager *tracker[azfake.PagerResponder[armstorage.BlobInventoryPoliciesClientListResponse]]
}

// Do implements the policy.Transporter interface for BlobInventoryPoliciesServerTransport.
func (b *BlobInventoryPoliciesServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return b.dispatchToMethodFake(req, method)
}

func (b *BlobInventoryPoliciesServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if blobInventoryPoliciesServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = blobInventoryPoliciesServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "BlobInventoryPoliciesClient.CreateOrUpdate":
				res.resp, res.err = b.dispatchCreateOrUpdate(req)
			case "BlobInventoryPoliciesClient.Delete":
				res.resp, res.err = b.dispatchDelete(req)
			case "BlobInventoryPoliciesClient.Get":
				res.resp, res.err = b.dispatchGet(req)
			case "BlobInventoryPoliciesClient.NewListPager":
				res.resp, res.err = b.dispatchNewListPager(req)
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

func (b *BlobInventoryPoliciesServerTransport) dispatchCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if b.srv.CreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method CreateOrUpdate not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Storage/storageAccounts/(?P<accountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/inventoryPolicies/(?P<blobInventoryPolicyName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armstorage.BlobInventoryPolicy](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	accountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("accountName")])
	if err != nil {
		return nil, err
	}
	blobInventoryPolicyNameParam, err := parseWithCast(matches[regex.SubexpIndex("blobInventoryPolicyName")], func(v string) (armstorage.BlobInventoryPolicyName, error) {
		p, unescapeErr := url.PathUnescape(v)
		if unescapeErr != nil {
			return "", unescapeErr
		}
		return armstorage.BlobInventoryPolicyName(p), nil
	})
	if err != nil {
		return nil, err
	}
	respr, errRespr := b.srv.CreateOrUpdate(req.Context(), resourceGroupNameParam, accountNameParam, blobInventoryPolicyNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).BlobInventoryPolicy, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (b *BlobInventoryPoliciesServerTransport) dispatchDelete(req *http.Request) (*http.Response, error) {
	if b.srv.Delete == nil {
		return nil, &nonRetriableError{errors.New("fake for method Delete not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Storage/storageAccounts/(?P<accountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/inventoryPolicies/(?P<blobInventoryPolicyName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	accountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("accountName")])
	if err != nil {
		return nil, err
	}
	blobInventoryPolicyNameParam, err := parseWithCast(matches[regex.SubexpIndex("blobInventoryPolicyName")], func(v string) (armstorage.BlobInventoryPolicyName, error) {
		p, unescapeErr := url.PathUnescape(v)
		if unescapeErr != nil {
			return "", unescapeErr
		}
		return armstorage.BlobInventoryPolicyName(p), nil
	})
	if err != nil {
		return nil, err
	}
	respr, errRespr := b.srv.Delete(req.Context(), resourceGroupNameParam, accountNameParam, blobInventoryPolicyNameParam, nil)
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

func (b *BlobInventoryPoliciesServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if b.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Storage/storageAccounts/(?P<accountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/inventoryPolicies/(?P<blobInventoryPolicyName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	accountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("accountName")])
	if err != nil {
		return nil, err
	}
	blobInventoryPolicyNameParam, err := parseWithCast(matches[regex.SubexpIndex("blobInventoryPolicyName")], func(v string) (armstorage.BlobInventoryPolicyName, error) {
		p, unescapeErr := url.PathUnescape(v)
		if unescapeErr != nil {
			return "", unescapeErr
		}
		return armstorage.BlobInventoryPolicyName(p), nil
	})
	if err != nil {
		return nil, err
	}
	respr, errRespr := b.srv.Get(req.Context(), resourceGroupNameParam, accountNameParam, blobInventoryPolicyNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).BlobInventoryPolicy, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (b *BlobInventoryPoliciesServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if b.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := b.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Storage/storageAccounts/(?P<accountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/inventoryPolicies`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		accountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("accountName")])
		if err != nil {
			return nil, err
		}
		resp := b.srv.NewListPager(resourceGroupNameParam, accountNameParam, nil)
		newListPager = &resp
		b.newListPager.add(req, newListPager)
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		b.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		b.newListPager.remove(req)
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to BlobInventoryPoliciesServerTransport
var blobInventoryPoliciesServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
