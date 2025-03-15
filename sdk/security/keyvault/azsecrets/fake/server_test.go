// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package fake_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets/fake"
	"github.com/stretchr/testify/require"
)

var (
	secretName  = "secretName"
	secretValue = "secretValue"
	version     = "123"
)

func getServer() fake.Server {
	return fake.Server{
		BackupSecret: func(ctx context.Context, name string, options *azsecrets.BackupSecretOptions) (resp azfake.Responder[azsecrets.BackupSecretResponse], errResp azfake.ErrorResponder) {
			kvResp := azsecrets.BackupSecretResponse{BackupSecretResult: azsecrets.BackupSecretResult{
				Value: []byte("backup"),
			}}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		DeleteSecret: func(ctx context.Context, name string, options *azsecrets.DeleteSecretOptions) (resp azfake.Responder[azsecrets.DeleteSecretResponse], errResp azfake.ErrorResponder) {
			kvResp := azsecrets.DeleteSecretResponse{DeletedSecret: azsecrets.DeletedSecret{
				ID: to.Ptr(azsecrets.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/secrets/%s/%s", name, version))),
			}}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		GetDeletedSecret: func(ctx context.Context, name string, options *azsecrets.GetDeletedSecretOptions) (resp azfake.Responder[azsecrets.GetDeletedSecretResponse], errResp azfake.ErrorResponder) {
			kvResp := azsecrets.GetDeletedSecretResponse{DeletedSecret: azsecrets.DeletedSecret{
				ID: to.Ptr(azsecrets.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/secrets/%s/%s", name, version))),
			}}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		GetSecret: func(ctx context.Context, name string, version string, options *azsecrets.GetSecretOptions) (resp azfake.Responder[azsecrets.GetSecretResponse], errResp azfake.ErrorResponder) {
			kvResp := azsecrets.GetSecretResponse{Secret: azsecrets.Secret{
				ID: to.Ptr(azsecrets.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/secrets/%s/%s", name, version))),
			}}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		NewListDeletedSecretPropertiesPager: func(options *azsecrets.ListDeletedSecretPropertiesOptions) (resp azfake.PagerResponder[azsecrets.ListDeletedSecretPropertiesResponse]) {
			page1 := azsecrets.ListDeletedSecretPropertiesResponse{
				DeletedSecretPropertiesListResult: azsecrets.DeletedSecretPropertiesListResult{
					Value: []*azsecrets.DeletedSecretProperties{
						{
							ID: to.Ptr(azsecrets.ID("https://fake-vault.vault.azure.net/secrets/secret1")),
						},
						{
							ID: to.Ptr(azsecrets.ID("https://fake-vault.vault.azure.net/secrets/secret2")),
						},
					},
				},
			}
			page2 := azsecrets.ListDeletedSecretPropertiesResponse{
				DeletedSecretPropertiesListResult: azsecrets.DeletedSecretPropertiesListResult{
					Value: []*azsecrets.DeletedSecretProperties{
						{
							ID: to.Ptr(azsecrets.ID("https://fake-vault.vault.azure.net/secrets/secret3")),
						},
						{
							ID: to.Ptr(azsecrets.ID("https://fake-vault.vault.azure.net/secrets/secret4")),
						},
					},
				},
			}

			resp.AddPage(http.StatusOK, page1, nil)
			resp.AddPage(http.StatusOK, page2, nil)
			return
		},
		NewListSecretPropertiesPager: func(options *azsecrets.ListSecretPropertiesOptions) (resp azfake.PagerResponder[azsecrets.ListSecretPropertiesResponse]) {
			page1 := azsecrets.ListSecretPropertiesResponse{
				SecretPropertiesListResult: azsecrets.SecretPropertiesListResult{
					Value: []*azsecrets.SecretProperties{
						{
							ID: to.Ptr(azsecrets.ID("https://fake-vault.vault.azure.net/secrets/secret1")),
						},
						{
							ID: to.Ptr(azsecrets.ID("https://fake-vault.vault.azure.net/secrets/secret2")),
						},
					},
				},
			}
			page2 := azsecrets.ListSecretPropertiesResponse{
				SecretPropertiesListResult: azsecrets.SecretPropertiesListResult{
					Value: []*azsecrets.SecretProperties{
						{
							ID: to.Ptr(azsecrets.ID("https://fake-vault.vault.azure.net/secrets/secret3")),
						},
						{
							ID: to.Ptr(azsecrets.ID("https://fake-vault.vault.azure.net/secrets/secret4")),
						},
					},
				},
			}

			resp.AddPage(http.StatusOK, page1, nil)
			resp.AddPage(http.StatusOK, page2, nil)
			return
		},
		NewListSecretPropertiesVersionsPager: func(name string, options *azsecrets.ListSecretPropertiesVersionsOptions) (resp azfake.PagerResponder[azsecrets.ListSecretPropertiesVersionsResponse]) {
			page1 := azsecrets.ListSecretPropertiesVersionsResponse{
				SecretPropertiesListResult: azsecrets.SecretPropertiesListResult{
					Value: []*azsecrets.SecretProperties{
						{
							ID: to.Ptr(azsecrets.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/secrets/%s/1", name))),
						},
						{
							ID: to.Ptr(azsecrets.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/secrets/%s/2", name))),
						},
					},
				},
			}
			page2 := azsecrets.ListSecretPropertiesVersionsResponse{
				SecretPropertiesListResult: azsecrets.SecretPropertiesListResult{
					Value: []*azsecrets.SecretProperties{
						{
							ID: to.Ptr(azsecrets.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/secrets/%s/3", name))),
						},
						{
							ID: to.Ptr(azsecrets.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/secrets/%s/4", name))),
						},
					},
				},
			}

			resp.AddPage(http.StatusOK, page1, nil)
			resp.AddPage(http.StatusOK, page2, nil)
			return
		},
		PurgeDeletedSecret: func(ctx context.Context, name string, options *azsecrets.PurgeDeletedSecretOptions) (resp azfake.Responder[azsecrets.PurgeDeletedSecretResponse], errResp azfake.ErrorResponder) {
			resp.SetResponse(http.StatusNoContent, azsecrets.PurgeDeletedSecretResponse{}, nil)
			return
		},
		RecoverDeletedSecret: func(ctx context.Context, name string, options *azsecrets.RecoverDeletedSecretOptions) (resp azfake.Responder[azsecrets.RecoverDeletedSecretResponse], errResp azfake.ErrorResponder) {
			kvResp := azsecrets.RecoverDeletedSecretResponse{Secret: azsecrets.Secret{
				ID: to.Ptr(azsecrets.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/secrets/%s/%s", name, version))),
			}}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		RestoreSecret: func(ctx context.Context, parameters azsecrets.RestoreSecretParameters, options *azsecrets.RestoreSecretOptions) (resp azfake.Responder[azsecrets.RestoreSecretResponse], errResp azfake.ErrorResponder) {
			kvResp := azsecrets.RestoreSecretResponse{Secret: azsecrets.Secret{
				ID: to.Ptr(azsecrets.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/secrets/%s/%s", secretName, version))),
			}}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		SetSecret: func(ctx context.Context, name string, parameters azsecrets.SetSecretParameters, options *azsecrets.SetSecretOptions) (resp azfake.Responder[azsecrets.SetSecretResponse], errResp azfake.ErrorResponder) {
			kvResp := azsecrets.SetSecretResponse{Secret: azsecrets.Secret{
				ID: to.Ptr(azsecrets.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/secrets/%s/%s", name, version))),
			}}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		UpdateSecretProperties: func(ctx context.Context, name string, version string, parameters azsecrets.UpdateSecretPropertiesParameters, options *azsecrets.UpdateSecretPropertiesOptions) (resp azfake.Responder[azsecrets.UpdateSecretPropertiesResponse], errResp azfake.ErrorResponder) {
			kvResp := azsecrets.UpdateSecretPropertiesResponse{
				Secret: azsecrets.Secret{
					ID: to.Ptr(azsecrets.ID(fmt.Sprintf("https://fake-vault.vault.azure.net/secrets/%s/%s", name, version))),
					Attributes: &azsecrets.SecretAttributes{
						Expires: to.Ptr(time.Date(2040, 1, 1, 1, 1, 1, 0, time.UTC)),
					},
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
	}
}

func TestServer(t *testing.T) {
	fakeServer := getServer()

	client, err := azsecrets.NewClient("https://fake-vault.vault.azure.net", &azfake.TokenCredential{}, &azsecrets.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewServerTransport(&fakeServer),
		},
	})
	require.NoError(t, err)

	// set secret
	setResp, err := client.SetSecret(context.TODO(), secretName, azsecrets.SetSecretParameters{Value: to.Ptr(secretValue)}, nil)
	require.NoError(t, err)
	require.Equal(t, secretName, setResp.ID.Name())
	require.Equal(t, version, setResp.ID.Version())

	// get secret
	getResp, err := client.GetSecret(context.TODO(), secretName, "", nil)
	require.NoError(t, err)
	require.Equal(t, secretName, getResp.ID.Name())
	require.Empty(t, getResp.ID.Version())

	// update secret properties
	date := time.Date(2040, 1, 1, 1, 1, 1, 0, time.UTC)
	updateResp, err := client.UpdateSecretProperties(context.TODO(), secretName, version, azsecrets.UpdateSecretPropertiesParameters{
		SecretAttributes: &azsecrets.SecretAttributes{
			Expires: to.Ptr(date),
		},
	}, nil)
	require.NoError(t, err)
	require.Equal(t, secretName, updateResp.ID.Name())
	require.Equal(t, version, updateResp.ID.Version())

	// new list secret properties versions
	versionsPager := client.NewListSecretPropertiesVersionsPager(secretName, nil)
	for versionsPager.More() {
		page, err := versionsPager.NextPage(context.TODO())
		require.NoError(t, err)

		for _, secret := range page.Value {
			require.NotNil(t, secret.ID)
			require.Equal(t, secret.ID.Name(), secretName)
		}
	}

	// new list secret properties pager
	pager := client.NewListSecretPropertiesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		require.NoError(t, err)

		for _, secret := range page.Value {
			require.NotNil(t, secret.ID)
			require.Contains(t, secret.ID.Name(), "secret")
		}
	}

	// delete secret
	deletedResp, err := client.DeleteSecret(context.TODO(), secretName, nil)
	require.NoError(t, err)
	require.Equal(t, secretName, deletedResp.ID.Name())
	require.Equal(t, version, deletedResp.ID.Version())

	// get deleted secret
	getDeletedResp, err := client.GetDeletedSecret(context.TODO(), secretName, nil)
	require.NoError(t, err)
	require.Equal(t, secretName, getDeletedResp.ID.Name())
	require.Equal(t, version, getDeletedResp.ID.Version())

	// new list deleted secret properties pager
	deletedPager := client.NewListDeletedSecretPropertiesPager(nil)
	for deletedPager.More() {
		page, err := deletedPager.NextPage(context.TODO())
		require.NoError(t, err)

		for _, secret := range page.Value {
			require.NotNil(t, secret.ID)
			require.Contains(t, secret.ID.Name(), "secret")
		}
	}

	// recover deleted secret
	recoverResp, err := client.RecoverDeletedSecret(context.TODO(), secretName, nil)
	require.NoError(t, err)
	require.Equal(t, secretName, recoverResp.ID.Name())
	require.Equal(t, version, recoverResp.ID.Version())

	// purge secret
	_, err = client.PurgeDeletedSecret(context.TODO(), secretName, nil)
	require.NoError(t, err)

	// backup secret
	backupResp, err := client.BackupSecret(context.TODO(), secretName, nil)
	require.NoError(t, err)
	require.NotNil(t, backupResp.Value)

	// restore secret
	restoreResp, err := client.RestoreSecret(context.TODO(), azsecrets.RestoreSecretParameters{SecretBackup: backupResp.Value}, nil)
	require.NoError(t, err)
	require.Equal(t, secretName, restoreResp.ID.Name())
	require.Equal(t, version, restoreResp.ID.Version())
}
