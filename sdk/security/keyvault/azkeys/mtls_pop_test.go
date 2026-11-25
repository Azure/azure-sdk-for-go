// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azkeys_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/stretchr/testify/require"
)

type popCredential struct{ t *testing.T }

func (c popCredential) GetToken(_ context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
	require.True(c.t, tro.EnableMTLSPoP, "Key Vault must signal PoP without a client enable flag")
	return azcore.AccessToken{Token: "bearer", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

func TestClientSignalsMTLSPoP(t *testing.T) {
	srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized),
		mock.WithHeader("WWW-Authenticate", `Bearer authorization="https://login.microsoftonline.com/tenant", resource="https://vault.azure.net"`))
	srv.AppendResponse(mock.WithBody([]byte(`{}`)))
	client, err := azkeys.NewClient("https://test.vault.azure.net", popCredential{t}, &azkeys.ClientOptions{
		ClientOptions: azcore.ClientOptions{Transport: srv},
	})
	require.NoError(t, err)
	_, err = client.GetKey(context.Background(), "key", "", nil)
	require.NoError(t, err)
	require.Equal(t, 2, srv.Requests())
}
