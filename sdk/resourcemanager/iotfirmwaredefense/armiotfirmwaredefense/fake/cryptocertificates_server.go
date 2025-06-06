// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package fake

import (
	"errors"
	"fmt"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake/server"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/iotfirmwaredefense/armiotfirmwaredefense/v2"
	"net/http"
	"net/url"
	"regexp"
)

// CryptoCertificatesServer is a fake server for instances of the armiotfirmwaredefense.CryptoCertificatesClient type.
type CryptoCertificatesServer struct {
	// NewListByFirmwarePager is the fake for method CryptoCertificatesClient.NewListByFirmwarePager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByFirmwarePager func(resourceGroupName string, workspaceName string, firmwareID string, options *armiotfirmwaredefense.CryptoCertificatesClientListByFirmwareOptions) (resp azfake.PagerResponder[armiotfirmwaredefense.CryptoCertificatesClientListByFirmwareResponse])
}

// NewCryptoCertificatesServerTransport creates a new instance of CryptoCertificatesServerTransport with the provided implementation.
// The returned CryptoCertificatesServerTransport instance is connected to an instance of armiotfirmwaredefense.CryptoCertificatesClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewCryptoCertificatesServerTransport(srv *CryptoCertificatesServer) *CryptoCertificatesServerTransport {
	return &CryptoCertificatesServerTransport{
		srv:                    srv,
		newListByFirmwarePager: newTracker[azfake.PagerResponder[armiotfirmwaredefense.CryptoCertificatesClientListByFirmwareResponse]](),
	}
}

// CryptoCertificatesServerTransport connects instances of armiotfirmwaredefense.CryptoCertificatesClient to instances of CryptoCertificatesServer.
// Don't use this type directly, use NewCryptoCertificatesServerTransport instead.
type CryptoCertificatesServerTransport struct {
	srv                    *CryptoCertificatesServer
	newListByFirmwarePager *tracker[azfake.PagerResponder[armiotfirmwaredefense.CryptoCertificatesClientListByFirmwareResponse]]
}

// Do implements the policy.Transporter interface for CryptoCertificatesServerTransport.
func (c *CryptoCertificatesServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return c.dispatchToMethodFake(req, method)
}

func (c *CryptoCertificatesServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if cryptoCertificatesServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = cryptoCertificatesServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "CryptoCertificatesClient.NewListByFirmwarePager":
				res.resp, res.err = c.dispatchNewListByFirmwarePager(req)
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

func (c *CryptoCertificatesServerTransport) dispatchNewListByFirmwarePager(req *http.Request) (*http.Response, error) {
	if c.srv.NewListByFirmwarePager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByFirmwarePager not implemented")}
	}
	newListByFirmwarePager := c.newListByFirmwarePager.get(req)
	if newListByFirmwarePager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.IoTFirmwareDefense/workspaces/(?P<workspaceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/firmwares/(?P<firmwareId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/cryptoCertificates`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		workspaceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("workspaceName")])
		if err != nil {
			return nil, err
		}
		firmwareIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("firmwareId")])
		if err != nil {
			return nil, err
		}
		resp := c.srv.NewListByFirmwarePager(resourceGroupNameParam, workspaceNameParam, firmwareIDParam, nil)
		newListByFirmwarePager = &resp
		c.newListByFirmwarePager.add(req, newListByFirmwarePager)
		server.PagerResponderInjectNextLinks(newListByFirmwarePager, req, func(page *armiotfirmwaredefense.CryptoCertificatesClientListByFirmwareResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByFirmwarePager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		c.newListByFirmwarePager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByFirmwarePager) {
		c.newListByFirmwarePager.remove(req)
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to CryptoCertificatesServerTransport
var cryptoCertificatesServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
