// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package fake_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys/fake"
	"github.com/stretchr/testify/require"
)

var (
	keyName    = "keyName"
	keyVersion = "123"
	vault      = "https://fake.vault.azure.net/keys"
)

func getServer() fake.Server {
	return fake.Server{
		BackupKey: func(ctx context.Context, name string, options *azkeys.BackupKeyOptions) (resp azfake.Responder[azkeys.BackupKeyResponse], errResp azfake.ErrorResponder) {
			kvResp := azkeys.BackupKeyResponse{
				BackupKeyResult: azkeys.BackupKeyResult{
					Value: []byte("testing"),
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		CreateKey: func(ctx context.Context, name string, parameters azkeys.CreateKeyParameters, options *azkeys.CreateKeyOptions) (resp azfake.Responder[azkeys.CreateKeyResponse], errResp azfake.ErrorResponder) {
			kvResp := azkeys.CreateKeyResponse{
				KeyBundle: azkeys.KeyBundle{
					Attributes: &azkeys.KeyAttributes{Enabled: parameters.KeyAttributes.Enabled},
					Key: &azkeys.JSONWebKey{
						KID: to.Ptr(azkeys.ID(fmt.Sprintf("%s/%s/%s", vault, name, keyVersion))),
					},
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		Decrypt: func(ctx context.Context, name string, version string, parameters azkeys.KeyOperationParameters, options *azkeys.DecryptOptions) (resp azfake.Responder[azkeys.DecryptResponse], errResp azfake.ErrorResponder) {
			kvResp := azkeys.DecryptResponse{
				KeyOperationResult: azkeys.KeyOperationResult{
					KID: to.Ptr(azkeys.ID(fmt.Sprintf("%s/%s/%s", vault, name, version))),
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		DeleteKey: func(ctx context.Context, name string, options *azkeys.DeleteKeyOptions) (resp azfake.Responder[azkeys.DeleteKeyResponse], errResp azfake.ErrorResponder) {
			kvResp := azkeys.DeleteKeyResponse{
				DeletedKey: azkeys.DeletedKey{
					Key: &azkeys.JSONWebKey{
						KID: to.Ptr(azkeys.ID(fmt.Sprintf("%s/%s/%s", vault, name, keyVersion))),
					},
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		Encrypt: func(ctx context.Context, name string, version string, parameters azkeys.KeyOperationParameters, options *azkeys.EncryptOptions) (resp azfake.Responder[azkeys.EncryptResponse], errResp azfake.ErrorResponder) {
			kvResp := azkeys.EncryptResponse{
				KeyOperationResult: azkeys.KeyOperationResult{
					KID: to.Ptr(azkeys.ID(fmt.Sprintf("%s/%s/%s", vault, name, version))),
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		GetDeletedKey: func(ctx context.Context, name string, options *azkeys.GetDeletedKeyOptions) (resp azfake.Responder[azkeys.GetDeletedKeyResponse], errResp azfake.ErrorResponder) {
			kvResp := azkeys.GetDeletedKeyResponse{
				DeletedKey: azkeys.DeletedKey{
					Key: &azkeys.JSONWebKey{
						KID: to.Ptr(azkeys.ID(fmt.Sprintf("%s/%s/%s", vault, name, keyVersion))),
					},
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		GetKey: func(ctx context.Context, name string, version string, options *azkeys.GetKeyOptions) (resp azfake.Responder[azkeys.GetKeyResponse], errResp azfake.ErrorResponder) {
			kvResp := azkeys.GetKeyResponse{
				KeyBundle: azkeys.KeyBundle{
					Key: &azkeys.JSONWebKey{
						KID: to.Ptr(azkeys.ID(fmt.Sprintf("%s/%s/%s", vault, name, version))),
					},
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		GetKeyAttestation: func(ctx context.Context, name string, version string, options *azkeys.GetKeyAttestationOptions) (resp azfake.Responder[azkeys.GetKeyAttestationResponse], errResp azfake.ErrorResponder) {
			kvResp := azkeys.GetKeyAttestationResponse{
				KeyBundle: azkeys.KeyBundle{
					Key: &azkeys.JSONWebKey{
						KID: to.Ptr(azkeys.ID(fmt.Sprintf("%s/%s/%s", vault, name, version))),
					},
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		GetKeyRotationPolicy: func(ctx context.Context, name string, options *azkeys.GetKeyRotationPolicyOptions) (resp azfake.Responder[azkeys.GetKeyRotationPolicyResponse], errResp azfake.ErrorResponder) {
			kvResp := azkeys.GetKeyRotationPolicyResponse{
				KeyRotationPolicy: azkeys.KeyRotationPolicy{ID: to.Ptr("keyPolicy")},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		GetRandomBytes: func(ctx context.Context, parameters azkeys.GetRandomBytesParameters, options *azkeys.GetRandomBytesOptions) (resp azfake.Responder[azkeys.GetRandomBytesResponse], errResp azfake.ErrorResponder) {
			kvResp := azkeys.GetRandomBytesResponse{
				RandomBytes: azkeys.RandomBytes{
					Value: []byte("testing"),
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		ImportKey: func(ctx context.Context, name string, parameters azkeys.ImportKeyParameters, options *azkeys.ImportKeyOptions) (resp azfake.Responder[azkeys.ImportKeyResponse], errResp azfake.ErrorResponder) {
			kvResp := azkeys.ImportKeyResponse{
				KeyBundle: azkeys.KeyBundle{
					Key: &azkeys.JSONWebKey{
						KID: to.Ptr(azkeys.ID(fmt.Sprintf("%s/%s/%s", vault, name, keyVersion))),
					},
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		NewListDeletedKeyPropertiesPager: func(options *azkeys.ListDeletedKeyPropertiesOptions) (resp azfake.PagerResponder[azkeys.ListDeletedKeyPropertiesResponse]) {
			page1 := azkeys.ListDeletedKeyPropertiesResponse{
				DeletedKeyPropertiesListResult: azkeys.DeletedKeyPropertiesListResult{
					Value: []*azkeys.DeletedKeyProperties{
						{
							KID: to.Ptr(azkeys.ID(fmt.Sprintf("%s/%s/%s", vault, "keyName1", keyVersion))),
						},
						{
							KID: to.Ptr(azkeys.ID(fmt.Sprintf("%s/%s/%s", vault, "keyName2", keyVersion))),
						},
					},
				},
			}
			page2 := azkeys.ListDeletedKeyPropertiesResponse{
				DeletedKeyPropertiesListResult: azkeys.DeletedKeyPropertiesListResult{
					Value: []*azkeys.DeletedKeyProperties{
						{
							KID: to.Ptr(azkeys.ID(fmt.Sprintf("%s/%s/%s", vault, "keyName3", keyVersion))),
						},
						{
							KID: to.Ptr(azkeys.ID(fmt.Sprintf("%s/%s/%s", vault, "keyName4", keyVersion))),
						},
					},
				},
			}

			resp.AddPage(http.StatusOK, page1, nil)
			resp.AddPage(http.StatusOK, page2, nil)
			return
		},
		NewListKeyPropertiesPager: func(options *azkeys.ListKeyPropertiesOptions) (resp azfake.PagerResponder[azkeys.ListKeyPropertiesResponse]) {
			page1 := azkeys.ListKeyPropertiesResponse{
				KeyPropertiesListResult: azkeys.KeyPropertiesListResult{
					Value: []*azkeys.KeyProperties{
						{
							KID: to.Ptr(azkeys.ID(fmt.Sprintf("%s/%s/%s", vault, "keyName1", keyVersion))),
						},
						{
							KID: to.Ptr(azkeys.ID(fmt.Sprintf("%s/%s/%s", vault, "keyName2", keyVersion))),
						},
					},
				},
			}
			page2 := azkeys.ListKeyPropertiesResponse{
				KeyPropertiesListResult: azkeys.KeyPropertiesListResult{
					Value: []*azkeys.KeyProperties{
						{
							KID: to.Ptr(azkeys.ID(fmt.Sprintf("%s/%s/%s", vault, "keyName3", keyVersion))),
						},
						{
							KID: to.Ptr(azkeys.ID(fmt.Sprintf("%s/%s/%s", vault, "keyName4", keyVersion))),
						},
					},
				},
			}

			resp.AddPage(http.StatusOK, page1, nil)
			resp.AddPage(http.StatusOK, page2, nil)
			return
		},
		NewListKeyPropertiesVersionsPager: func(name string, options *azkeys.ListKeyPropertiesVersionsOptions) (resp azfake.PagerResponder[azkeys.ListKeyPropertiesVersionsResponse]) {
			page1 := azkeys.ListKeyPropertiesVersionsResponse{
				KeyPropertiesListResult: azkeys.KeyPropertiesListResult{
					Value: []*azkeys.KeyProperties{
						{
							KID: to.Ptr(azkeys.ID(fmt.Sprintf("%s/%s/%s", vault, name, "1"))),
						},
						{
							KID: to.Ptr(azkeys.ID(fmt.Sprintf("%s/%s/%s", vault, name, "2"))),
						},
					},
				},
			}
			page2 := azkeys.ListKeyPropertiesVersionsResponse{
				KeyPropertiesListResult: azkeys.KeyPropertiesListResult{
					Value: []*azkeys.KeyProperties{
						{
							KID: to.Ptr(azkeys.ID(fmt.Sprintf("%s/%s/%s", vault, name, keyVersion))),
						},
						{
							KID: to.Ptr(azkeys.ID(fmt.Sprintf("%s/%s/%s", vault, name, keyVersion))),
						},
					},
				},
			}

			resp.AddPage(http.StatusOK, page1, nil)
			resp.AddPage(http.StatusOK, page2, nil)
			return
		},
		PurgeDeletedKey: func(ctx context.Context, name string, options *azkeys.PurgeDeletedKeyOptions) (resp azfake.Responder[azkeys.PurgeDeletedKeyResponse], errResp azfake.ErrorResponder) {
			kvResp := azkeys.PurgeDeletedKeyResponse{}
			resp.SetResponse(http.StatusNoContent, kvResp, nil)
			return
		},
		RecoverDeletedKey: func(ctx context.Context, name string, options *azkeys.RecoverDeletedKeyOptions) (resp azfake.Responder[azkeys.RecoverDeletedKeyResponse], errResp azfake.ErrorResponder) {
			kvResp := azkeys.RecoverDeletedKeyResponse{
				KeyBundle: azkeys.KeyBundle{
					Key: &azkeys.JSONWebKey{
						KID: to.Ptr(azkeys.ID(fmt.Sprintf("%s/%s/%s", vault, name, keyVersion))),
					},
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		Release: func(ctx context.Context, name string, version string, parameters azkeys.ReleaseParameters, options *azkeys.ReleaseOptions) (resp azfake.Responder[azkeys.ReleaseResponse], errResp azfake.ErrorResponder) {
			kvResp := azkeys.ReleaseResponse{
				KeyReleaseResult: azkeys.KeyReleaseResult{
					Value: to.Ptr("test"),
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		RestoreKey: func(ctx context.Context, parameters azkeys.RestoreKeyParameters, options *azkeys.RestoreKeyOptions) (resp azfake.Responder[azkeys.RestoreKeyResponse], errResp azfake.ErrorResponder) {
			kvResp := azkeys.RestoreKeyResponse{
				KeyBundle: azkeys.KeyBundle{
					Key: &azkeys.JSONWebKey{
						KID: to.Ptr(azkeys.ID(fmt.Sprintf("%s/%s/%s", vault, keyName, keyVersion))),
					},
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		RotateKey: func(ctx context.Context, name string, options *azkeys.RotateKeyOptions) (resp azfake.Responder[azkeys.RotateKeyResponse], errResp azfake.ErrorResponder) {
			kvResp := azkeys.RotateKeyResponse{
				KeyBundle: azkeys.KeyBundle{
					Key: &azkeys.JSONWebKey{
						KID: to.Ptr(azkeys.ID(fmt.Sprintf("%s/%s/%s", vault, name, keyVersion))),
					},
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		Sign: func(ctx context.Context, name string, version string, parameters azkeys.SignParameters, options *azkeys.SignOptions) (resp azfake.Responder[azkeys.SignResponse], errResp azfake.ErrorResponder) {
			kvResp := azkeys.SignResponse{
				KeyOperationResult: azkeys.KeyOperationResult{
					KID: to.Ptr(azkeys.ID(fmt.Sprintf("%s/%s/%s", vault, name, version))),
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		UnwrapKey: func(ctx context.Context, name string, version string, parameters azkeys.KeyOperationParameters, options *azkeys.UnwrapKeyOptions) (resp azfake.Responder[azkeys.UnwrapKeyResponse], errResp azfake.ErrorResponder) {
			kvResp := azkeys.UnwrapKeyResponse{
				KeyOperationResult: azkeys.KeyOperationResult{
					KID: to.Ptr(azkeys.ID(fmt.Sprintf("%s/%s/%s", vault, name, version))),
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		UpdateKey: func(ctx context.Context, name string, version string, parameters azkeys.UpdateKeyParameters, options *azkeys.UpdateKeyOptions) (resp azfake.Responder[azkeys.UpdateKeyResponse], errResp azfake.ErrorResponder) {
			kvResp := azkeys.UpdateKeyResponse{
				KeyBundle: azkeys.KeyBundle{
					Key: &azkeys.JSONWebKey{
						KID: to.Ptr(azkeys.ID(fmt.Sprintf("%s/%s/%s", vault, name, keyVersion))),
					},
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		UpdateKeyRotationPolicy: func(ctx context.Context, name string, keyRotationPolicy azkeys.KeyRotationPolicy, options *azkeys.UpdateKeyRotationPolicyOptions) (resp azfake.Responder[azkeys.UpdateKeyRotationPolicyResponse], errResp azfake.ErrorResponder) {
			kvResp := azkeys.UpdateKeyRotationPolicyResponse{
				KeyRotationPolicy: azkeys.KeyRotationPolicy{ID: keyRotationPolicy.ID},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		Verify: func(ctx context.Context, name string, version string, parameters azkeys.VerifyParameters, options *azkeys.VerifyOptions) (resp azfake.Responder[azkeys.VerifyResponse], errResp azfake.ErrorResponder) {
			kvResp := azkeys.VerifyResponse{
				KeyVerifyResult: azkeys.KeyVerifyResult{
					Value: to.Ptr(true),
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		WrapKey: func(ctx context.Context, name string, version string, parameters azkeys.KeyOperationParameters, options *azkeys.WrapKeyOptions) (resp azfake.Responder[azkeys.WrapKeyResponse], errResp azfake.ErrorResponder) {
			kvResp := azkeys.WrapKeyResponse{
				KeyOperationResult: azkeys.KeyOperationResult{
					KID: to.Ptr(azkeys.ID(fmt.Sprintf("%s/%s/%s", vault, name, version))),
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
	}
}

func TestServer(t *testing.T) {
	fakeServer := getServer()

	client, err := azkeys.NewClient("https://fake-vault.vault.azure.net", &azfake.TokenCredential{}, &azkeys.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewServerTransport(&fakeServer),
		},
	})
	require.NoError(t, err)

	backupResp, err := client.BackupKey(context.Background(), keyName, nil)
	require.NoError(t, err)
	require.NotNil(t, backupResp.Value)

	createResp, err := client.CreateKey(context.Background(), keyName, azkeys.CreateKeyParameters{KeyAttributes: &azkeys.KeyAttributes{Enabled: to.Ptr(true)}}, nil)
	require.NoError(t, err)
	require.Equal(t, keyName, createResp.Key.KID.Name())
	require.True(t, *createResp.Attributes.Enabled)

	decryptResp, err := client.Decrypt(context.Background(), keyName, "", azkeys.KeyOperationParameters{}, nil)
	require.NoError(t, err)
	require.Equal(t, keyName, decryptResp.KID.Name())
	require.Empty(t, decryptResp.KID.Version())

	deleteResp, err := client.DeleteKey(context.Background(), keyName, nil)
	require.NoError(t, err)
	require.Equal(t, keyName, deleteResp.Key.KID.Name())

	encryptResp, err := client.Encrypt(context.Background(), keyName, keyVersion, azkeys.KeyOperationParameters{}, nil)
	require.NoError(t, err)
	require.Equal(t, keyName, encryptResp.KID.Name())
	require.Equal(t, keyVersion, encryptResp.KID.Version())

	getDeleteResp, err := client.GetDeletedKey(context.Background(), keyName, nil)
	require.NoError(t, err)
	require.Equal(t, keyName, getDeleteResp.Key.KID.Name())

	getResp, err := client.GetKey(context.Background(), keyName, keyVersion, nil)
	require.NoError(t, err)
	require.Equal(t, keyName, getResp.Key.KID.Name())
	require.Equal(t, keyVersion, getResp.Key.KID.Version())

	getAttestationResp, err := client.GetKeyAttestation(context.Background(), keyName, keyVersion, nil)
	require.NoError(t, err)
	require.Equal(t, keyName, getAttestationResp.Key.KID.Name())
	require.Equal(t, keyVersion, getAttestationResp.Key.KID.Version())

	getPolicyResp, err := client.GetKeyRotationPolicy(context.Background(), keyName, nil)
	require.NoError(t, err)
	require.NotEmpty(t, getPolicyResp.ID)

	getBytesResp, err := client.GetRandomBytes(context.Background(), azkeys.GetRandomBytesParameters{}, nil)
	require.NoError(t, err)
	require.NotEmpty(t, getBytesResp.Value)

	importResp, err := client.ImportKey(context.Background(), keyName, azkeys.ImportKeyParameters{}, nil)
	require.NoError(t, err)
	require.Equal(t, keyName, importResp.Key.KID.Name())

	deletedPager := client.NewListDeletedKeyPropertiesPager(nil)
	for deletedPager.More() {
		page, err := deletedPager.NextPage(context.Background())
		require.NoError(t, err)

		for _, key := range page.Value {
			require.Contains(t, key.KID.Name(), keyName)
		}
	}

	keyPager := client.NewListKeyPropertiesPager(nil)
	for keyPager.More() {
		page, err := keyPager.NextPage(context.Background())
		require.NoError(t, err)

		for _, key := range page.Value {
			require.Contains(t, key.KID.Name(), keyName)
		}
	}

	keyVersionsPager := client.NewListKeyPropertiesVersionsPager(keyName, nil)
	for keyVersionsPager.More() {
		page, err := keyVersionsPager.NextPage(context.Background())
		require.NoError(t, err)

		for _, key := range page.Value {
			require.Equal(t, key.KID.Name(), keyName)
		}
	}

	purgeResp, err := client.PurgeDeletedKey(context.Background(), keyName, nil)
	require.NoError(t, err)
	require.NotNil(t, purgeResp)

	recoverResp, err := client.RecoverDeletedKey(context.Background(), keyName, nil)
	require.NoError(t, err)
	require.Equal(t, keyName, recoverResp.Key.KID.Name())

	releaseResp, err := client.Release(context.Background(), keyName, keyVersion, azkeys.ReleaseParameters{}, nil)
	require.NoError(t, err)
	require.NotEmpty(t, releaseResp.Value)

	restoreResp, err := client.RestoreKey(context.Background(), azkeys.RestoreKeyParameters{}, nil)
	require.NoError(t, err)
	require.NotEmpty(t, restoreResp.Key.KID)

	rotateResp, err := client.RotateKey(context.Background(), keyName, nil)
	require.NoError(t, err)
	require.Equal(t, keyName, rotateResp.Key.KID.Name())

	signResp, err := client.Sign(context.Background(), keyName, keyVersion, azkeys.SignParameters{}, nil)
	require.NoError(t, err)
	require.Equal(t, keyName, signResp.KID.Name())
	require.Equal(t, keyVersion, signResp.KID.Version())

	unwrapResp, err := client.UnwrapKey(context.Background(), keyName, keyVersion, azkeys.KeyOperationParameters{}, nil)
	require.NoError(t, err)
	require.Equal(t, keyName, unwrapResp.KID.Name())
	require.Equal(t, keyVersion, unwrapResp.KID.Version())

	updateResp, err := client.UpdateKey(context.Background(), keyName, keyVersion, azkeys.UpdateKeyParameters{}, nil)
	require.NoError(t, err)
	require.Equal(t, keyName, updateResp.Key.KID.Name())
	require.Equal(t, keyVersion, updateResp.Key.KID.Version())

	updatePolicyResp, err := client.UpdateKeyRotationPolicy(context.Background(), keyName, azkeys.KeyRotationPolicy{ID: to.Ptr("test")}, nil)
	require.NoError(t, err)
	require.NotEmpty(t, updatePolicyResp.ID)

	verifyResp, err := client.Verify(context.Background(), keyName, keyVersion, azkeys.VerifyParameters{}, nil)
	require.NoError(t, err)
	require.True(t, *verifyResp.Value)

	wrapResp, err := client.WrapKey(context.Background(), keyName, keyVersion, azkeys.KeyOperationParameters{}, nil)
	require.NoError(t, err)
	require.Equal(t, keyName, wrapResp.KID.Name())
	require.Equal(t, keyVersion, wrapResp.KID.Version())

}
