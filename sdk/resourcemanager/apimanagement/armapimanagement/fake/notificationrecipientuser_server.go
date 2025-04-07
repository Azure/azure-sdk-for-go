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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement/v3"
	"net/http"
	"net/url"
	"regexp"
)

// NotificationRecipientUserServer is a fake server for instances of the armapimanagement.NotificationRecipientUserClient type.
type NotificationRecipientUserServer struct {
	// CheckEntityExists is the fake for method NotificationRecipientUserClient.CheckEntityExists
	// HTTP status codes to indicate success: http.StatusNoContent, http.StatusNotFound
	CheckEntityExists func(ctx context.Context, resourceGroupName string, serviceName string, notificationName armapimanagement.NotificationName, userID string, options *armapimanagement.NotificationRecipientUserClientCheckEntityExistsOptions) (resp azfake.Responder[armapimanagement.NotificationRecipientUserClientCheckEntityExistsResponse], errResp azfake.ErrorResponder)

	// CreateOrUpdate is the fake for method NotificationRecipientUserClient.CreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	CreateOrUpdate func(ctx context.Context, resourceGroupName string, serviceName string, notificationName armapimanagement.NotificationName, userID string, options *armapimanagement.NotificationRecipientUserClientCreateOrUpdateOptions) (resp azfake.Responder[armapimanagement.NotificationRecipientUserClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// Delete is the fake for method NotificationRecipientUserClient.Delete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	Delete func(ctx context.Context, resourceGroupName string, serviceName string, notificationName armapimanagement.NotificationName, userID string, options *armapimanagement.NotificationRecipientUserClientDeleteOptions) (resp azfake.Responder[armapimanagement.NotificationRecipientUserClientDeleteResponse], errResp azfake.ErrorResponder)

	// ListByNotification is the fake for method NotificationRecipientUserClient.ListByNotification
	// HTTP status codes to indicate success: http.StatusOK
	ListByNotification func(ctx context.Context, resourceGroupName string, serviceName string, notificationName armapimanagement.NotificationName, options *armapimanagement.NotificationRecipientUserClientListByNotificationOptions) (resp azfake.Responder[armapimanagement.NotificationRecipientUserClientListByNotificationResponse], errResp azfake.ErrorResponder)
}

// NewNotificationRecipientUserServerTransport creates a new instance of NotificationRecipientUserServerTransport with the provided implementation.
// The returned NotificationRecipientUserServerTransport instance is connected to an instance of armapimanagement.NotificationRecipientUserClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewNotificationRecipientUserServerTransport(srv *NotificationRecipientUserServer) *NotificationRecipientUserServerTransport {
	return &NotificationRecipientUserServerTransport{srv: srv}
}

// NotificationRecipientUserServerTransport connects instances of armapimanagement.NotificationRecipientUserClient to instances of NotificationRecipientUserServer.
// Don't use this type directly, use NewNotificationRecipientUserServerTransport instead.
type NotificationRecipientUserServerTransport struct {
	srv *NotificationRecipientUserServer
}

// Do implements the policy.Transporter interface for NotificationRecipientUserServerTransport.
func (n *NotificationRecipientUserServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return n.dispatchToMethodFake(req, method)
}

func (n *NotificationRecipientUserServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if notificationRecipientUserServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = notificationRecipientUserServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "NotificationRecipientUserClient.CheckEntityExists":
				res.resp, res.err = n.dispatchCheckEntityExists(req)
			case "NotificationRecipientUserClient.CreateOrUpdate":
				res.resp, res.err = n.dispatchCreateOrUpdate(req)
			case "NotificationRecipientUserClient.Delete":
				res.resp, res.err = n.dispatchDelete(req)
			case "NotificationRecipientUserClient.ListByNotification":
				res.resp, res.err = n.dispatchListByNotification(req)
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

func (n *NotificationRecipientUserServerTransport) dispatchCheckEntityExists(req *http.Request) (*http.Response, error) {
	if n.srv.CheckEntityExists == nil {
		return nil, &nonRetriableError{errors.New("fake for method CheckEntityExists not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/notifications/(?P<notificationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/recipientUsers/(?P<userId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 5 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	serviceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serviceName")])
	if err != nil {
		return nil, err
	}
	notificationNameParam, err := parseWithCast(matches[regex.SubexpIndex("notificationName")], func(v string) (armapimanagement.NotificationName, error) {
		p, unescapeErr := url.PathUnescape(v)
		if unescapeErr != nil {
			return "", unescapeErr
		}
		return armapimanagement.NotificationName(p), nil
	})
	if err != nil {
		return nil, err
	}
	userIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("userId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := n.srv.CheckEntityExists(req.Context(), resourceGroupNameParam, serviceNameParam, notificationNameParam, userIDParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusNoContent, http.StatusNotFound}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusNoContent, http.StatusNotFound", respContent.HTTPStatus)}
	}
	resp, err := server.NewResponse(respContent, req, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (n *NotificationRecipientUserServerTransport) dispatchCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if n.srv.CreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method CreateOrUpdate not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/notifications/(?P<notificationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/recipientUsers/(?P<userId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 5 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	serviceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serviceName")])
	if err != nil {
		return nil, err
	}
	notificationNameParam, err := parseWithCast(matches[regex.SubexpIndex("notificationName")], func(v string) (armapimanagement.NotificationName, error) {
		p, unescapeErr := url.PathUnescape(v)
		if unescapeErr != nil {
			return "", unescapeErr
		}
		return armapimanagement.NotificationName(p), nil
	})
	if err != nil {
		return nil, err
	}
	userIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("userId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := n.srv.CreateOrUpdate(req.Context(), resourceGroupNameParam, serviceNameParam, notificationNameParam, userIDParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusCreated}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).RecipientUserContract, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (n *NotificationRecipientUserServerTransport) dispatchDelete(req *http.Request) (*http.Response, error) {
	if n.srv.Delete == nil {
		return nil, &nonRetriableError{errors.New("fake for method Delete not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/notifications/(?P<notificationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/recipientUsers/(?P<userId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 5 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	serviceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serviceName")])
	if err != nil {
		return nil, err
	}
	notificationNameParam, err := parseWithCast(matches[regex.SubexpIndex("notificationName")], func(v string) (armapimanagement.NotificationName, error) {
		p, unescapeErr := url.PathUnescape(v)
		if unescapeErr != nil {
			return "", unescapeErr
		}
		return armapimanagement.NotificationName(p), nil
	})
	if err != nil {
		return nil, err
	}
	userIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("userId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := n.srv.Delete(req.Context(), resourceGroupNameParam, serviceNameParam, notificationNameParam, userIDParam, nil)
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

func (n *NotificationRecipientUserServerTransport) dispatchListByNotification(req *http.Request) (*http.Response, error) {
	if n.srv.ListByNotification == nil {
		return nil, &nonRetriableError{errors.New("fake for method ListByNotification not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/notifications/(?P<notificationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/recipientUsers`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	serviceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serviceName")])
	if err != nil {
		return nil, err
	}
	notificationNameParam, err := parseWithCast(matches[regex.SubexpIndex("notificationName")], func(v string) (armapimanagement.NotificationName, error) {
		p, unescapeErr := url.PathUnescape(v)
		if unescapeErr != nil {
			return "", unescapeErr
		}
		return armapimanagement.NotificationName(p), nil
	})
	if err != nil {
		return nil, err
	}
	respr, errRespr := n.srv.ListByNotification(req.Context(), resourceGroupNameParam, serviceNameParam, notificationNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).RecipientUserCollection, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to NotificationRecipientUserServerTransport
var notificationRecipientUserServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
