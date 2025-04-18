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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservicesbackup/v4"
	"net/http"
	"net/url"
	"regexp"
)

// PrivateEndpointServer is a fake server for instances of the armrecoveryservicesbackup.PrivateEndpointClient type.
type PrivateEndpointServer struct {
	// GetOperationStatus is the fake for method PrivateEndpointClient.GetOperationStatus
	// HTTP status codes to indicate success: http.StatusOK
	GetOperationStatus func(ctx context.Context, vaultName string, resourceGroupName string, privateEndpointConnectionName string, operationID string, options *armrecoveryservicesbackup.PrivateEndpointClientGetOperationStatusOptions) (resp azfake.Responder[armrecoveryservicesbackup.PrivateEndpointClientGetOperationStatusResponse], errResp azfake.ErrorResponder)
}

// NewPrivateEndpointServerTransport creates a new instance of PrivateEndpointServerTransport with the provided implementation.
// The returned PrivateEndpointServerTransport instance is connected to an instance of armrecoveryservicesbackup.PrivateEndpointClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewPrivateEndpointServerTransport(srv *PrivateEndpointServer) *PrivateEndpointServerTransport {
	return &PrivateEndpointServerTransport{srv: srv}
}

// PrivateEndpointServerTransport connects instances of armrecoveryservicesbackup.PrivateEndpointClient to instances of PrivateEndpointServer.
// Don't use this type directly, use NewPrivateEndpointServerTransport instead.
type PrivateEndpointServerTransport struct {
	srv *PrivateEndpointServer
}

// Do implements the policy.Transporter interface for PrivateEndpointServerTransport.
func (p *PrivateEndpointServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return p.dispatchToMethodFake(req, method)
}

func (p *PrivateEndpointServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if privateEndpointServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = privateEndpointServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "PrivateEndpointClient.GetOperationStatus":
				res.resp, res.err = p.dispatchGetOperationStatus(req)
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

func (p *PrivateEndpointServerTransport) dispatchGetOperationStatus(req *http.Request) (*http.Response, error) {
	if p.srv.GetOperationStatus == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetOperationStatus not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.RecoveryServices/vaults/(?P<vaultName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/privateEndpointConnections/(?P<privateEndpointConnectionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/operationsStatus/(?P<operationId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 5 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	vaultNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("vaultName")])
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	privateEndpointConnectionNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("privateEndpointConnectionName")])
	if err != nil {
		return nil, err
	}
	operationIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("operationId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := p.srv.GetOperationStatus(req.Context(), vaultNameParam, resourceGroupNameParam, privateEndpointConnectionNameParam, operationIDParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).OperationStatus, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to PrivateEndpointServerTransport
var privateEndpointServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
