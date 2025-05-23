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

// BMSPrepareDataMoveOperationResultServer is a fake server for instances of the armrecoveryservicesbackup.BMSPrepareDataMoveOperationResultClient type.
type BMSPrepareDataMoveOperationResultServer struct {
	// Get is the fake for method BMSPrepareDataMoveOperationResultClient.Get
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	Get func(ctx context.Context, vaultName string, resourceGroupName string, operationID string, options *armrecoveryservicesbackup.BMSPrepareDataMoveOperationResultClientGetOptions) (resp azfake.Responder[armrecoveryservicesbackup.BMSPrepareDataMoveOperationResultClientGetResponse], errResp azfake.ErrorResponder)
}

// NewBMSPrepareDataMoveOperationResultServerTransport creates a new instance of BMSPrepareDataMoveOperationResultServerTransport with the provided implementation.
// The returned BMSPrepareDataMoveOperationResultServerTransport instance is connected to an instance of armrecoveryservicesbackup.BMSPrepareDataMoveOperationResultClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewBMSPrepareDataMoveOperationResultServerTransport(srv *BMSPrepareDataMoveOperationResultServer) *BMSPrepareDataMoveOperationResultServerTransport {
	return &BMSPrepareDataMoveOperationResultServerTransport{srv: srv}
}

// BMSPrepareDataMoveOperationResultServerTransport connects instances of armrecoveryservicesbackup.BMSPrepareDataMoveOperationResultClient to instances of BMSPrepareDataMoveOperationResultServer.
// Don't use this type directly, use NewBMSPrepareDataMoveOperationResultServerTransport instead.
type BMSPrepareDataMoveOperationResultServerTransport struct {
	srv *BMSPrepareDataMoveOperationResultServer
}

// Do implements the policy.Transporter interface for BMSPrepareDataMoveOperationResultServerTransport.
func (b *BMSPrepareDataMoveOperationResultServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return b.dispatchToMethodFake(req, method)
}

func (b *BMSPrepareDataMoveOperationResultServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if bmsPrepareDataMoveOperationResultServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = bmsPrepareDataMoveOperationResultServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "BMSPrepareDataMoveOperationResultClient.Get":
				res.resp, res.err = b.dispatchGet(req)
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

func (b *BMSPrepareDataMoveOperationResultServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if b.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.RecoveryServices/vaults/(?P<vaultName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/backupstorageconfig/vaultstorageconfig/operationResults/(?P<operationId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
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
	operationIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("operationId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := b.srv.Get(req.Context(), vaultNameParam, resourceGroupNameParam, operationIDParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusAccepted}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).VaultStorageConfigOperationResultResponseClassification, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to BMSPrepareDataMoveOperationResultServerTransport
var bmsPrepareDataMoveOperationResultServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
