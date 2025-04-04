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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/netapp/armnetapp/v7"
	"net/http"
	"net/url"
	"regexp"
)

// BackupVaultsServer is a fake server for instances of the armnetapp.BackupVaultsClient type.
type BackupVaultsServer struct {
	// BeginCreateOrUpdate is the fake for method BackupVaultsClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, accountName string, backupVaultName string, body armnetapp.BackupVault, options *armnetapp.BackupVaultsClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armnetapp.BackupVaultsClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method BackupVaultsClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, accountName string, backupVaultName string, options *armnetapp.BackupVaultsClientBeginDeleteOptions) (resp azfake.PollerResponder[armnetapp.BackupVaultsClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method BackupVaultsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, accountName string, backupVaultName string, options *armnetapp.BackupVaultsClientGetOptions) (resp azfake.Responder[armnetapp.BackupVaultsClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByNetAppAccountPager is the fake for method BackupVaultsClient.NewListByNetAppAccountPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByNetAppAccountPager func(resourceGroupName string, accountName string, options *armnetapp.BackupVaultsClientListByNetAppAccountOptions) (resp azfake.PagerResponder[armnetapp.BackupVaultsClientListByNetAppAccountResponse])

	// BeginUpdate is the fake for method BackupVaultsClient.BeginUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginUpdate func(ctx context.Context, resourceGroupName string, accountName string, backupVaultName string, body armnetapp.BackupVaultPatch, options *armnetapp.BackupVaultsClientBeginUpdateOptions) (resp azfake.PollerResponder[armnetapp.BackupVaultsClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewBackupVaultsServerTransport creates a new instance of BackupVaultsServerTransport with the provided implementation.
// The returned BackupVaultsServerTransport instance is connected to an instance of armnetapp.BackupVaultsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewBackupVaultsServerTransport(srv *BackupVaultsServer) *BackupVaultsServerTransport {
	return &BackupVaultsServerTransport{
		srv:                         srv,
		beginCreateOrUpdate:         newTracker[azfake.PollerResponder[armnetapp.BackupVaultsClientCreateOrUpdateResponse]](),
		beginDelete:                 newTracker[azfake.PollerResponder[armnetapp.BackupVaultsClientDeleteResponse]](),
		newListByNetAppAccountPager: newTracker[azfake.PagerResponder[armnetapp.BackupVaultsClientListByNetAppAccountResponse]](),
		beginUpdate:                 newTracker[azfake.PollerResponder[armnetapp.BackupVaultsClientUpdateResponse]](),
	}
}

// BackupVaultsServerTransport connects instances of armnetapp.BackupVaultsClient to instances of BackupVaultsServer.
// Don't use this type directly, use NewBackupVaultsServerTransport instead.
type BackupVaultsServerTransport struct {
	srv                         *BackupVaultsServer
	beginCreateOrUpdate         *tracker[azfake.PollerResponder[armnetapp.BackupVaultsClientCreateOrUpdateResponse]]
	beginDelete                 *tracker[azfake.PollerResponder[armnetapp.BackupVaultsClientDeleteResponse]]
	newListByNetAppAccountPager *tracker[azfake.PagerResponder[armnetapp.BackupVaultsClientListByNetAppAccountResponse]]
	beginUpdate                 *tracker[azfake.PollerResponder[armnetapp.BackupVaultsClientUpdateResponse]]
}

// Do implements the policy.Transporter interface for BackupVaultsServerTransport.
func (b *BackupVaultsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return b.dispatchToMethodFake(req, method)
}

func (b *BackupVaultsServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if backupVaultsServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = backupVaultsServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "BackupVaultsClient.BeginCreateOrUpdate":
				res.resp, res.err = b.dispatchBeginCreateOrUpdate(req)
			case "BackupVaultsClient.BeginDelete":
				res.resp, res.err = b.dispatchBeginDelete(req)
			case "BackupVaultsClient.Get":
				res.resp, res.err = b.dispatchGet(req)
			case "BackupVaultsClient.NewListByNetAppAccountPager":
				res.resp, res.err = b.dispatchNewListByNetAppAccountPager(req)
			case "BackupVaultsClient.BeginUpdate":
				res.resp, res.err = b.dispatchBeginUpdate(req)
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

func (b *BackupVaultsServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if b.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := b.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.NetApp/netAppAccounts/(?P<accountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/backupVaults/(?P<backupVaultName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armnetapp.BackupVault](req)
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
		backupVaultNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("backupVaultName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := b.srv.BeginCreateOrUpdate(req.Context(), resourceGroupNameParam, accountNameParam, backupVaultNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateOrUpdate = &respr
		b.beginCreateOrUpdate.add(req, beginCreateOrUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		b.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		b.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (b *BackupVaultsServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if b.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := b.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.NetApp/netAppAccounts/(?P<accountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/backupVaults/(?P<backupVaultName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
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
		backupVaultNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("backupVaultName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := b.srv.BeginDelete(req.Context(), resourceGroupNameParam, accountNameParam, backupVaultNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		b.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		b.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		b.beginDelete.remove(req)
	}

	return resp, nil
}

func (b *BackupVaultsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if b.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.NetApp/netAppAccounts/(?P<accountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/backupVaults/(?P<backupVaultName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
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
	backupVaultNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("backupVaultName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := b.srv.Get(req.Context(), resourceGroupNameParam, accountNameParam, backupVaultNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).BackupVault, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (b *BackupVaultsServerTransport) dispatchNewListByNetAppAccountPager(req *http.Request) (*http.Response, error) {
	if b.srv.NewListByNetAppAccountPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByNetAppAccountPager not implemented")}
	}
	newListByNetAppAccountPager := b.newListByNetAppAccountPager.get(req)
	if newListByNetAppAccountPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.NetApp/netAppAccounts/(?P<accountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/backupVaults`
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
		resp := b.srv.NewListByNetAppAccountPager(resourceGroupNameParam, accountNameParam, nil)
		newListByNetAppAccountPager = &resp
		b.newListByNetAppAccountPager.add(req, newListByNetAppAccountPager)
		server.PagerResponderInjectNextLinks(newListByNetAppAccountPager, req, func(page *armnetapp.BackupVaultsClientListByNetAppAccountResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByNetAppAccountPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		b.newListByNetAppAccountPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByNetAppAccountPager) {
		b.newListByNetAppAccountPager.remove(req)
	}
	return resp, nil
}

func (b *BackupVaultsServerTransport) dispatchBeginUpdate(req *http.Request) (*http.Response, error) {
	if b.srv.BeginUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdate not implemented")}
	}
	beginUpdate := b.beginUpdate.get(req)
	if beginUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.NetApp/netAppAccounts/(?P<accountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/backupVaults/(?P<backupVaultName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armnetapp.BackupVaultPatch](req)
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
		backupVaultNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("backupVaultName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := b.srv.BeginUpdate(req.Context(), resourceGroupNameParam, accountNameParam, backupVaultNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginUpdate = &respr
		b.beginUpdate.add(req, beginUpdate)
	}

	resp, err := server.PollerResponderNext(beginUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		b.beginUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdate) {
		b.beginUpdate.remove(req)
	}

	return resp, nil
}

// set this to conditionally intercept incoming requests to BackupVaultsServerTransport
var backupVaultsServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
