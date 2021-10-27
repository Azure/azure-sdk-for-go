// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestDefaultAzureCredential_GetTokenSuccess(t *testing.T) {
	env := map[string]string{"AZURE_TENANT_ID": tenantID, "AZURE_CLIENT_ID": clientID, "AZURE_CLIENT_SECRET": secret}
	setEnvironmentVariables(t, env)
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))

	cred, err := NewDefaultAzureCredential(&DefaultAzureCredentialOptions{AuthorityHost: AuthorityHost(srv.URL()), ClientOptions: policy.ClientOptions{Transport: srv}})
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{"scope"}})
	if err != nil {
		t.Fatalf("GetToken error: %v", err)
	}
}
