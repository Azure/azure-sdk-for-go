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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql/v2"
	"net/http"
	"net/url"
	"regexp"
)

// DatabaseEncryptionProtectorsServer is a fake server for instances of the armsql.DatabaseEncryptionProtectorsClient type.
type DatabaseEncryptionProtectorsServer struct {
	// BeginRevalidate is the fake for method DatabaseEncryptionProtectorsClient.BeginRevalidate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginRevalidate func(ctx context.Context, resourceGroupName string, serverName string, databaseName string, encryptionProtectorName armsql.EncryptionProtectorName, options *armsql.DatabaseEncryptionProtectorsClientBeginRevalidateOptions) (resp azfake.PollerResponder[armsql.DatabaseEncryptionProtectorsClientRevalidateResponse], errResp azfake.ErrorResponder)

	// BeginRevert is the fake for method DatabaseEncryptionProtectorsClient.BeginRevert
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginRevert func(ctx context.Context, resourceGroupName string, serverName string, databaseName string, encryptionProtectorName armsql.EncryptionProtectorName, options *armsql.DatabaseEncryptionProtectorsClientBeginRevertOptions) (resp azfake.PollerResponder[armsql.DatabaseEncryptionProtectorsClientRevertResponse], errResp azfake.ErrorResponder)
}

// NewDatabaseEncryptionProtectorsServerTransport creates a new instance of DatabaseEncryptionProtectorsServerTransport with the provided implementation.
// The returned DatabaseEncryptionProtectorsServerTransport instance is connected to an instance of armsql.DatabaseEncryptionProtectorsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewDatabaseEncryptionProtectorsServerTransport(srv *DatabaseEncryptionProtectorsServer) *DatabaseEncryptionProtectorsServerTransport {
	return &DatabaseEncryptionProtectorsServerTransport{
		srv:             srv,
		beginRevalidate: newTracker[azfake.PollerResponder[armsql.DatabaseEncryptionProtectorsClientRevalidateResponse]](),
		beginRevert:     newTracker[azfake.PollerResponder[armsql.DatabaseEncryptionProtectorsClientRevertResponse]](),
	}
}

// DatabaseEncryptionProtectorsServerTransport connects instances of armsql.DatabaseEncryptionProtectorsClient to instances of DatabaseEncryptionProtectorsServer.
// Don't use this type directly, use NewDatabaseEncryptionProtectorsServerTransport instead.
type DatabaseEncryptionProtectorsServerTransport struct {
	srv             *DatabaseEncryptionProtectorsServer
	beginRevalidate *tracker[azfake.PollerResponder[armsql.DatabaseEncryptionProtectorsClientRevalidateResponse]]
	beginRevert     *tracker[azfake.PollerResponder[armsql.DatabaseEncryptionProtectorsClientRevertResponse]]
}

// Do implements the policy.Transporter interface for DatabaseEncryptionProtectorsServerTransport.
func (d *DatabaseEncryptionProtectorsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "DatabaseEncryptionProtectorsClient.BeginRevalidate":
		resp, err = d.dispatchBeginRevalidate(req)
	case "DatabaseEncryptionProtectorsClient.BeginRevert":
		resp, err = d.dispatchBeginRevert(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (d *DatabaseEncryptionProtectorsServerTransport) dispatchBeginRevalidate(req *http.Request) (*http.Response, error) {
	if d.srv.BeginRevalidate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginRevalidate not implemented")}
	}
	beginRevalidate := d.beginRevalidate.get(req)
	if beginRevalidate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Sql/servers/(?P<serverName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/databases/(?P<databaseName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/encryptionProtector/(?P<encryptionProtectorName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/revalidate`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 5 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		serverNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serverName")])
		if err != nil {
			return nil, err
		}
		databaseNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("databaseName")])
		if err != nil {
			return nil, err
		}
		encryptionProtectorNameParam, err := parseWithCast(matches[regex.SubexpIndex("encryptionProtectorName")], func(v string) (armsql.EncryptionProtectorName, error) {
			p, unescapeErr := url.PathUnescape(v)
			if unescapeErr != nil {
				return "", unescapeErr
			}
			return armsql.EncryptionProtectorName(p), nil
		})
		if err != nil {
			return nil, err
		}
		respr, errRespr := d.srv.BeginRevalidate(req.Context(), resourceGroupNameParam, serverNameParam, databaseNameParam, encryptionProtectorNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginRevalidate = &respr
		d.beginRevalidate.add(req, beginRevalidate)
	}

	resp, err := server.PollerResponderNext(beginRevalidate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		d.beginRevalidate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginRevalidate) {
		d.beginRevalidate.remove(req)
	}

	return resp, nil
}

func (d *DatabaseEncryptionProtectorsServerTransport) dispatchBeginRevert(req *http.Request) (*http.Response, error) {
	if d.srv.BeginRevert == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginRevert not implemented")}
	}
	beginRevert := d.beginRevert.get(req)
	if beginRevert == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Sql/servers/(?P<serverName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/databases/(?P<databaseName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/encryptionProtector/(?P<encryptionProtectorName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/revert`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 5 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		serverNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serverName")])
		if err != nil {
			return nil, err
		}
		databaseNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("databaseName")])
		if err != nil {
			return nil, err
		}
		encryptionProtectorNameParam, err := parseWithCast(matches[regex.SubexpIndex("encryptionProtectorName")], func(v string) (armsql.EncryptionProtectorName, error) {
			p, unescapeErr := url.PathUnescape(v)
			if unescapeErr != nil {
				return "", unescapeErr
			}
			return armsql.EncryptionProtectorName(p), nil
		})
		if err != nil {
			return nil, err
		}
		respr, errRespr := d.srv.BeginRevert(req.Context(), resourceGroupNameParam, serverNameParam, databaseNameParam, encryptionProtectorNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginRevert = &respr
		d.beginRevert.add(req, beginRevert)
	}

	resp, err := server.PollerResponderNext(beginRevert, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		d.beginRevert.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginRevert) {
		d.beginRevert.remove(req)
	}

	return resp, nil
}
