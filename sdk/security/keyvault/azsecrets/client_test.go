// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsecrets_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	azcred "github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/stretchr/testify/require"
)

func TestBackupRestore(t *testing.T) {
	client := startTest(t)

	name := createRandomName(t, "testbackupsecret")
	value := createRandomName(t, "value")

	setResp, err := client.SetSecret(context.Background(), name, azsecrets.SetSecretParameters{Value: &value}, nil)
	require.NoError(t, err)
	defer cleanUpSecret(t, client, name)

	backupResp, err := client.BackupSecret(context.Background(), name, nil)
	require.NoError(t, err)
	require.Greater(t, len(backupResp.Value), 0)
	testSerde(t, &backupResp.BackupSecretResult)

	_, err = client.DeleteSecret(context.Background(), name, nil)
	require.NoError(t, err)
	pollStatus(t, 404, func() error {
		_, err := client.GetDeletedSecret(context.Background(), name, nil)
		return err
	})

	_, err = client.PurgeDeletedSecret(context.Background(), name, nil)
	require.NoError(t, err)

	var restoreResp azsecrets.RestoreSecretResponse
	restoreParams := azsecrets.RestoreSecretParameters{backupResp.Value}
	testSerde(t, &restoreParams)
	pollStatus(t, 409, func() error {
		restoreResp, err = client.RestoreSecret(context.Background(), restoreParams, nil)
		return err
	})
	require.Equal(t, restoreResp.ID.Name(), name)
	require.Equal(t, setResp.ID, restoreResp.ID)
}

func TestCRUD(t *testing.T) {
	client := startTest(t)

	name := createRandomName(t, "secret")
	value := createRandomName(t, "value")

	setParams := azsecrets.SetSecretParameters{
		ContentType: to.Ptr("big secret"),
		SecretAttributes: &azsecrets.SecretAttributes{
			Enabled:   to.Ptr(true),
			NotBefore: to.Ptr(time.Date(2030, 1, 1, 1, 1, 1, 0, time.UTC)),
		},
		Tags:  map[string]*string{"tag": to.Ptr("value")},
		Value: &value,
	}
	testSerde(t, &setParams)
	setResp, err := client.SetSecret(context.Background(), name, setParams, nil)
	require.NoError(t, err)
	require.Equal(t, setParams.ContentType, setResp.ContentType)
	require.Equal(t, setParams.SecretAttributes.Enabled, setResp.Attributes.Enabled)
	require.Equal(t, setParams.SecretAttributes.NotBefore.Unix(), setResp.Attributes.NotBefore.Unix())
	require.Equal(t, setParams.Tags, setResp.Tags)
	require.Equal(t, setParams.Value, setResp.Value)
	require.Equal(t, name, setResp.ID.Name())
	require.NotEmpty(t, setResp.ID.Version())
	testSerde(t, &setResp.Secret)

	getResp, err := client.GetSecret(context.Background(), setResp.ID.Name(), "", nil)
	require.NoError(t, err)
	require.Equal(t, setParams.ContentType, getResp.ContentType)
	require.NotNil(t, setResp.ID)
	require.Equal(t, setResp.ID, getResp.ID)
	require.Equal(t, setResp.ID.Name(), getResp.ID.Name())
	require.Equal(t, setResp.ID.Version(), getResp.ID.Version())
	require.Equal(t, setParams.SecretAttributes.Enabled, getResp.Attributes.Enabled)
	require.Equal(t, setParams.SecretAttributes.NotBefore.Unix(), getResp.Attributes.NotBefore.Unix())
	require.Equal(t, setParams.Tags, getResp.Tags)
	require.Equal(t, setParams.Value, getResp.Value)

	updateParams := azsecrets.UpdateSecretPropertiesParameters{
		SecretAttributes: &azsecrets.SecretAttributes{
			Expires: to.Ptr(time.Date(2040, 1, 1, 1, 1, 1, 0, time.UTC)),
		},
	}
	testSerde(t, &updateParams)
	updateResp, err := client.UpdateSecretProperties(context.Background(), name, setResp.ID.Version(), updateParams, nil)
	require.NoError(t, err)
	require.Equal(t, setParams.ContentType, updateResp.ContentType)
	require.Equal(t, setResp.ID, updateResp.ID)
	require.Equal(t, setParams.SecretAttributes.Enabled, updateResp.Attributes.Enabled)
	require.Equal(t, setParams.SecretAttributes.NotBefore.Unix(), updateResp.Attributes.NotBefore.Unix())
	require.Equal(t, setParams.Tags, updateResp.Tags)
	require.Equal(t, setResp.ID.Version(), updateResp.ID.Version())

	deleteResp, err := client.DeleteSecret(context.Background(), name, nil)
	require.NoError(t, err)
	require.Equal(t, setParams.ContentType, deleteResp.ContentType)
	require.Equal(t, setResp.ID, deleteResp.ID)
	require.Equal(t, setParams.SecretAttributes.Enabled, deleteResp.Attributes.Enabled)
	require.Equal(t, updateParams.SecretAttributes.Expires.Unix(), deleteResp.Attributes.Expires.Unix())
	require.Equal(t, setParams.SecretAttributes.NotBefore.Unix(), deleteResp.Attributes.NotBefore.Unix())
	require.Equal(t, setParams.Tags, deleteResp.Tags)
	require.Equal(t, name, deleteResp.ID.Name())
	require.Equal(t, updateResp.ID.Version(), deleteResp.ID.Version())
	testSerde(t, &deleteResp.DeletedSecret)
	pollStatus(t, 404, func() error {
		_, err := client.GetDeletedSecret(context.Background(), name, nil)
		return err
	})

	getDeletedResp, err := client.GetDeletedSecret(context.Background(), name, nil)
	require.NoError(t, err)
	require.Equal(t, setParams.ContentType, getDeletedResp.ContentType)
	require.Equal(t, setParams.SecretAttributes.Enabled, getDeletedResp.Attributes.Enabled)
	require.Equal(t, updateParams.SecretAttributes.Expires.Unix(), getDeletedResp.Attributes.Expires.Unix())
	require.Equal(t, setParams.SecretAttributes.NotBefore.Unix(), getDeletedResp.Attributes.NotBefore.Unix())
	require.Equal(t, setParams.Tags, getDeletedResp.Tags)
	require.Equal(t, name, getDeletedResp.ID.Name())
	require.Equal(t, setResp.ID.Version(), getDeletedResp.ID.Version())

	_, err = client.PurgeDeletedSecret(context.Background(), name, nil)
	require.NoError(t, err)
}

func TestDisableChallengeResourceVerification(t *testing.T) {
	authResource := `"Bearer authorization="https://login.microsoftonline.com/tenant", resource="%s""`
	authScope := `"Bearer authorization="https://login.microsoftonline.com/tenant", scope="%s""`
	vaultURL := "https://fakevault.vault.azure.net"
	for _, test := range []struct {
		challenge, resource string
		disableVerify, err  bool
	}{
		// happy path: resource matches requested vault's host (vault.azure.net)
		{challenge: authResource, resource: "https://vault.azure.net"},
		{challenge: authScope, resource: "https://vault.azure.net/.default"},
		{challenge: authResource, resource: "https://vault.azure.net", disableVerify: true},
		{challenge: authScope, resource: "https://vault.azure.net/.default", disableVerify: true},

		// error cases: resource/scope doesn't match the requested vault's host (vault.azure.net)
		{challenge: authResource, resource: "https://vault.azure.cn", err: true},
		{challenge: authResource, resource: "https://myvault.azure.net", err: true},
		{challenge: authScope, resource: "https://vault.azure.cn/.default", err: true},
		{challenge: authScope, resource: "https://myvault.azure.net/.default", err: true},

		// the policy shouldn't return errors for the above error cases when verification is disabled
		{challenge: authResource, resource: "https://vault.azure.cn", disableVerify: true},
		{challenge: authResource, resource: "https://myvault.azure.net", disableVerify: true},
		{challenge: authScope, resource: "https://vault.azure.cn/.default", disableVerify: true},
		{challenge: authScope, resource: "https://myvault.azure.net/.default", disableVerify: true},
	} {
		t.Run("", func(t *testing.T) {
			srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer close()
			srv.AppendResponse(mock.WithStatusCode(401), mock.WithHeader("WWW-Authenticate", fmt.Sprintf(test.challenge, test.resource)))
			srv.AppendResponse(mock.WithStatusCode(200), mock.WithBody([]byte(`{"value":[]}`)))
			options := &azsecrets.ClientOptions{
				ClientOptions: policy.ClientOptions{
					Transport: srv,
				},
				DisableChallengeResourceVerification: test.disableVerify,
			}
			client, err := azsecrets.NewClient(vaultURL, &azcred.Fake{}, options)
			require.NoError(t, err)
			pager := client.NewListSecretPropertiesPager(nil)
			_, err = pager.NextPage(context.Background())
			if test.err {
				require.Error(t, err)
				require.Contains(t, err.Error(), "challenge resource")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestID(t *testing.T) {
	for _, test := range []struct{ ID, name, version string }{
		{"https://foo.vault.azure.net/secrets/name/version", "name", "version"},
		{"https://foo.vault.azure.net/secrets/name", "name", ""},
	} {
		t.Run(test.ID, func(t *testing.T) {
			ID := azsecrets.ID(test.ID)
			require.Equal(t, test.name, ID.Name())
			require.Equal(t, test.version, ID.Version())
		})
	}
}

func TestListDeletedSecrets(t *testing.T) {
	client := startTest(t)

	secret1 := createRandomName(t, "secret1")
	value1 := createRandomName(t, "value1")
	secret2 := createRandomName(t, "secret2")
	value2 := createRandomName(t, "value2")

	_, err := client.SetSecret(context.Background(), secret1, azsecrets.SetSecretParameters{Value: &value1}, nil)
	require.NoError(t, err)
	_, err = client.DeleteSecret(context.Background(), secret1, nil)
	require.NoError(t, err)
	_, err = client.SetSecret(context.Background(), secret2, azsecrets.SetSecretParameters{Value: &value2}, nil)
	require.NoError(t, err)
	_, err = client.DeleteSecret(context.Background(), secret2, nil)
	require.NoError(t, err)
	defer func() {
		_, err := client.PurgeDeletedSecret(context.Background(), secret1, nil)
		require.NoError(t, err)
		_, err = client.PurgeDeletedSecret(context.Background(), secret2, nil)
		require.NoError(t, err)
	}()

	pollStatus(t, 404, func() error {
		_, err := client.GetDeletedSecret(context.Background(), secret1, nil)
		return err
	})
	pollStatus(t, 404, func() error {
		_, err := client.GetDeletedSecret(context.Background(), secret2, nil)
		return err
	})

	expected := map[string]struct{}{secret1: {}, secret2: {}}
	pager := client.NewListDeletedSecretPropertiesPager(nil)
	for pager.More() && len(expected) > 0 {
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		testSerde(t, &page.DeletedSecretPropertiesListResult)
		for _, secret := range page.Value {
			testSerde(t, secret)
			delete(expected, secret.ID.Name())
			if len(expected) == 0 {
				break
			}
		}
	}
	require.Empty(t, expected, "pager didn't return all expected secrets")
}

func TestListSecrets(t *testing.T) {
	client := startTest(t)

	count := 4
	for i := 0; i < count; i++ {
		name := createRandomName(t, fmt.Sprintf("listsecrets%d", i))
		value := createRandomName(t, fmt.Sprintf("value%d", i))
		_, err := client.SetSecret(context.Background(), name, azsecrets.SetSecretParameters{Value: &value}, nil)
		require.NoError(t, err)
		defer cleanUpSecret(t, client, name)
	}

	pager := client.NewListSecretPropertiesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		testSerde(t, &page.SecretPropertiesListResult)
		for _, secret := range page.Value {
			testSerde(t, secret)
			if strings.HasPrefix(secret.ID.Name(), "listsecrets") {
				count--
			}
		}
	}
	require.Equal(t, count, 0)
}

func TestListSecretVersions(t *testing.T) {
	client := startTest(t)

	name := createRandomName(t, "listversions")
	commonParams := azsecrets.SetSecretParameters{
		ContentType: to.Ptr("content-type"),
		Tags:        map[string]*string{"tag": to.Ptr("value")},
		SecretAttributes: &azsecrets.SecretAttributes{
			Expires:   to.Ptr(time.Date(2050, 1, 1, 1, 1, 1, 0, time.UTC)),
			NotBefore: to.Ptr(time.Date(2040, 1, 1, 1, 1, 1, 0, time.UTC)),
		},
	}
	count := 3
	for i := 0; i < count; i++ {
		params := commonParams
		params.Value = to.Ptr(fmt.Sprintf("value%d", i))
		res, err := client.SetSecret(context.Background(), name, params, nil)
		require.Equal(t, params.Value, res.Value)
		require.NoError(t, err)
	}
	defer cleanUpSecret(t, client, name)

	pager := client.NewListSecretPropertiesVersionsPager(name, nil)
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		testSerde(t, &page.SecretPropertiesListResult)
		for i, secret := range page.Value {
			testSerde(t, secret)
			if i > 0 {
				require.NotEqual(t, page.Value[i-1].ID.Version(), secret.ID.Version())
			}
			require.NotNil(t, secret.ID)
			require.Equal(t, name, secret.ID.Name())
			if strings.HasPrefix(secret.ID.Name(), name) {
				count--
				require.Equal(t, commonParams.ContentType, secret.ContentType)
				require.Equal(t, commonParams.SecretAttributes.Expires.Unix(), secret.Attributes.Expires.Unix())
				require.Equal(t, commonParams.SecretAttributes.NotBefore.Unix(), secret.Attributes.NotBefore.Unix())
				require.Equal(t, commonParams.Tags, secret.Tags)
			}
		}
	}
	require.Equal(t, count, 0)
}

func TestNameRequired(t *testing.T) {
	client, err := azsecrets.NewClient(fakeVaultURL, &azcred.Fake{}, nil)
	require.NoError(t, err)
	expected := "parameter name cannot be empty"
	_, err = client.BackupSecret(context.Background(), "", nil)
	require.EqualError(t, err, expected)
	_, err = client.DeleteSecret(context.Background(), "", nil)
	require.EqualError(t, err, expected)
	_, err = client.GetDeletedSecret(context.Background(), "", nil)
	require.EqualError(t, err, expected)
	_, err = client.GetSecret(context.Background(), "", "", nil)
	require.EqualError(t, err, expected)
	_, err = client.PurgeDeletedSecret(context.Background(), "", nil)
	require.EqualError(t, err, expected)
	_, err = client.RecoverDeletedSecret(context.Background(), "", nil)
	require.EqualError(t, err, expected)
	_, err = client.SetSecret(context.Background(), "", azsecrets.SetSecretParameters{}, nil)
	require.EqualError(t, err, expected)
	_, err = client.UpdateSecretProperties(context.Background(), "", "", azsecrets.UpdateSecretPropertiesParameters{}, nil)
	require.EqualError(t, err, expected)
}

func TestRecover(t *testing.T) {
	client := startTest(t)

	name := createRandomName(t, "secret")
	value := createRandomName(t, "value")

	setResp, err := client.SetSecret(context.Background(), name, azsecrets.SetSecretParameters{Value: &value}, nil)
	require.NoError(t, err)
	defer cleanUpSecret(t, client, name)
	require.Equal(t, value, *setResp.Value)

	_, err = client.DeleteSecret(context.Background(), name, nil)
	require.NoError(t, err)

	pollStatus(t, 404, func() error {
		_, err := client.GetDeletedSecret(context.Background(), name, nil)
		return err
	})

	recoverResp, err := client.RecoverDeletedSecret(context.Background(), name, nil)
	require.NoError(t, err)
	require.Equal(t, setResp.ID, recoverResp.ID)

	var getResp azsecrets.GetSecretResponse
	pollStatus(t, 404, func() error {
		getResp, err = client.GetSecret(context.Background(), name, "", nil)
		return err
	})
	require.Equal(t, value, *getResp.Value)
	require.Equal(t, setResp.Attributes, getResp.Attributes)
	require.Equal(t, setResp.ID, getResp.ID)
	require.Equal(t, setResp.ContentType, getResp.ContentType)
}

func TestAPIVersion(t *testing.T) {
	apiVersion := "2025-07-01"
	var requireVersion = func(req *http.Request) bool {
		version := req.URL.Query().Get("api-version")
		require.Equal(t, version, apiVersion)
		return true
	}
	srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer close()
	srv.AppendResponse(
		mock.WithStatusCode(200),
		mock.WithPredicate(requireVersion),
	)
	srv.AppendResponse(mock.WithStatusCode(http.StatusInternalServerError))

	opts := &azsecrets.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport:  srv,
			APIVersion: apiVersion,
		},
	}
	client, err := azsecrets.NewClient(vaultURL, &azcred.Fake{}, opts)
	require.NoError(t, err)

	_, err = client.GetSecret(context.Background(), "name", "", nil)
	require.NoError(t, err)
}

func TestPreviousVersion(t *testing.T) {
	tests := []struct {
		name     string
		jsonBody string
		expected *string
	}{
		{name: "no previousVersion", jsonBody: `{"kid":"https://vault/secret/secret-name"}`, expected: nil},
		{name: "with previousVersion", jsonBody: `{"kid":"https://vault/secret/secret-name","previousVersion":"dummyVersion"}`, expected: to.Ptr("dummyVersion")},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var s azsecrets.Secret
			err := json.Unmarshal([]byte(tc.jsonBody), &s)
			require.NoError(t, err)
			if tc.expected == nil {
				require.Nil(t, s.PreviousVersion)
			} else {
				require.NotNil(t, s.PreviousVersion)
				require.Equal(t, *tc.expected, *s.PreviousVersion)
			}

			var ds azsecrets.DeletedSecret
			err = json.Unmarshal([]byte(tc.jsonBody), &ds)
			require.NoError(t, err)
			if tc.expected == nil {
				require.Nil(t, ds.PreviousVersion)
			} else {
				require.NotNil(t, ds.PreviousVersion)
				require.Equal(t, *tc.expected, *ds.PreviousVersion)
			}
		})
	}
}
