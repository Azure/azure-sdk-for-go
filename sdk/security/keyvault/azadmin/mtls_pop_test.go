// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azadmin_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/backup"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/ekm"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/rbac"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/settings"
	"github.com/stretchr/testify/require"
)

type popCredential struct{ t *testing.T }

func (c popCredential) GetToken(_ context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
	require.True(c.t, tro.EnableMTLSPoP, "Managed HSM must signal PoP without a client enable flag")
	return azcore.AccessToken{Token: "bearer", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

func TestClientsSignalMTLSPoP(t *testing.T) {
	const endpoint = "https://test.managedhsm.azure.net"
	for _, name := range []string{"backup", "ekm", "rbac", "settings"} {
		t.Run(name, func(t *testing.T) {
			srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer close()
			srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized),
				mock.WithHeader("WWW-Authenticate", `Bearer authorization="https://login.microsoftonline.com/tenant", resource="https://managedhsm.azure.net"`))
			if name == "backup" {
				srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted),
					mock.WithHeader("Azure-AsyncOperation", endpoint+"/backup/status/job"),
					mock.WithBody([]byte(`{"status":"InProgress"}`)))
			} else {
				srv.AppendResponse(mock.WithBody([]byte(`{}`)))
			}
			options := azcore.ClientOptions{Transport: srv}
			switch name {
			case "backup":
				client, err := backup.NewClient(endpoint, popCredential{t}, &backup.ClientOptions{ClientOptions: options})
				require.NoError(t, err)
				_, err = client.BeginFullBackup(context.Background(), backup.SASTokenParameters{}, nil)
				require.NoError(t, err)
			case "ekm":
				client, err := ekm.NewClient(endpoint, popCredential{t}, &ekm.ClientOptions{ClientOptions: options})
				require.NoError(t, err)
				_, err = client.GetEkmConnection(context.Background(), nil)
				require.NoError(t, err)
			case "rbac":
				client, err := rbac.NewClient(endpoint, popCredential{t}, &rbac.ClientOptions{ClientOptions: options})
				require.NoError(t, err)
				_, err = client.GetRoleDefinition(context.Background(), rbac.RoleScopeGlobal, "role", nil)
				require.NoError(t, err)
			case "settings":
				client, err := settings.NewClient(endpoint, popCredential{t}, &settings.ClientOptions{ClientOptions: options})
				require.NoError(t, err)
				_, err = client.GetSettings(context.Background(), nil)
				require.NoError(t, err)
			}
			require.Equal(t, 2, srv.Requests())
		})
	}
}
