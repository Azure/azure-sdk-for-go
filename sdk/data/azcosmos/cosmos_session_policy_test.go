// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/assert"
)

func TestSessionPolicy(t *testing.T) {
	tests := []struct {
		name                     string
		isWriteOperation         bool
		overrideToken            bool
		responseStatus           int
		responseIncludesToken    bool
		expectSetToken           bool
		expectClearToken         bool
		expectSessionTokenHeader string
	}{
		{name: "ReadUsesToken", isWriteOperation: false, responseStatus: http.StatusOK, responseIncludesToken: false, expectSetToken: false, expectClearToken: false, expectSessionTokenHeader: "fake-cached-token"},
		{name: "ReadOverridingTokenIgnoresCache", isWriteOperation: false, overrideToken: true, responseStatus: http.StatusOK, responseIncludesToken: false, expectSetToken: false, expectClearToken: false, expectSessionTokenHeader: "fake-user-token"},
		{name: "SuccessfulWriteSetsToken", isWriteOperation: true, responseStatus: http.StatusOK, responseIncludesToken: true, expectSetToken: true, expectClearToken: false},
		{name: "FailedWriteClearsToken", isWriteOperation: true, responseStatus: http.StatusNotFound, responseIncludesToken: false, expectSetToken: false, expectClearToken: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			srv, close := mock.NewTLSServer()
			defer close()
			srv.SetResponse(mock.WithStatusCode(test.responseStatus))

			mockSessionContainer := &mockSessionContainer{}
			verifier := &sessionPolicyVerify{setSessionResponseHeaders: test.responseIncludesToken}
			pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{
				PerRetry: []policy.Policy{&sessionPolicy{sc: mockSessionContainer}, verifier}},
				&policy.ClientOptions{Transport: srv})

			req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			req.SetOperationValue(pipelineRequestOptions{
				isWriteOperation: test.isWriteOperation,
			})

			if test.overrideToken {
				req.Raw().Header.Set(cosmosHeaderSessionToken, "fake-user-token")
			}

			_, err = pl.Do(req)
			assert.NoError(t, err)

			assert.Equal(t, test.expectSetToken, mockSessionContainer.setSessionTokenInvoked)
			assert.Equal(t, test.expectClearToken, mockSessionContainer.clearSessionTokenInvoked)
			assert.Equal(t, test.expectSessionTokenHeader, verifier.sessionTokenHeaderSent)
		})
	}
}

type mockSessionContainer struct {
	clearSessionTokenInvoked bool
	getSessionTokenInvoked   bool
	setSessionTokenInvoked   bool
}

func (mc *mockSessionContainer) GetSessionToken(resourceAddress string) string {
	mc.getSessionTokenInvoked = true
	return "fake-cached-token"
}

func (mc *mockSessionContainer) SetSessionToken(resourceAddress string, containerRid string, sessionToken string) {
	mc.setSessionTokenInvoked = true
}

func (mc *mockSessionContainer) ClearSessionToken(resourceAddress string) {
	mc.clearSessionTokenInvoked = true
}

type sessionPolicyVerify struct {
	setSessionResponseHeaders bool
	sessionTokenHeaderSent    string
}

func (p *sessionPolicyVerify) Do(req *policy.Request) (*http.Response, error) {
	p.sessionTokenHeaderSent = req.Raw().Header.Get(cosmosHeaderSessionToken)
	resp, err := req.Next()
	if err != nil {
		return resp, err
	}

	if p.setSessionResponseHeaders {
		resp.Header.Set(cosmosHeaderSessionToken, "fake-response-token")
		resp.Header.Set(cosmosHeaderAltContentPath, "fake-alt-content-path")
		resp.Header.Set(cosmosHeaderContentPath, "fake-content-path")
	}

	return resp, nil
}
