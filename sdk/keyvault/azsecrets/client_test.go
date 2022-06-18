//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsecrets_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/stretchr/testify/require"
)

// pollStatus calls a function until it stops returning a response error with the given status code.
// If this takes more than 2 minutes, it fails the test.
func pollStatus(t *testing.T, expectedStatus int, fn func() error) {
	var err error
	for i := 0; i < 12; i++ {
		err = fn()
		var respErr *azcore.ResponseError
		if !(errors.As(err, &respErr) && respErr.StatusCode == expectedStatus) {
			break
		}
		if i < 11 {
			recording.Sleep(10 * time.Second)
		}
	}
	require.NoError(t, err)
}

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
	pollStatus(t, 409, func() error {
		restoreResp, err = client.RestoreSecret(context.Background(), restoreParams, nil)
		return err
	})
	require.Equal(t, restoreResp.ID.Name(), name)
	require.Equal(t, setResp.ID, restoreResp.ID)

	// exercise otherwise unused unmarshalling code
	rp := azsecrets.RestoreSecretParameters{}
	data, err := restoreParams.MarshalJSON()
	require.NoError(t, err)
	err = rp.UnmarshalJSON(data)
	require.NoError(t, err)
	require.Equal(t, restoreParams, rp)
}

func TestCRUD(t *testing.T) {
	client := startTest(t)

	name := createRandomName(t, "secret")
	value := createRandomName(t, "value")

	// TODO: would be nice to promote value; it's actually required, and Key Vault's error message doesn't say so
	setParams := azsecrets.SetSecretParameters{
		ContentType: to.Ptr("big secret"),
		SecretAttributes: &azsecrets.SecretAttributes{
			Enabled:   to.Ptr(true),
			NotBefore: to.Ptr(time.Date(2030, 1, 1, 1, 1, 1, 0, time.UTC)),
		},
		Tags:  map[string]*string{"tag": to.Ptr("value")},
		Value: &value,
	}
	setResp, err := client.SetSecret(context.Background(), name, setParams, nil)
	require.NoError(t, err)
	require.Equal(t, setParams.ContentType, setResp.ContentType)
	// TODO: params has "SecretAttributes" field; response has "Attributes" field of same type
	require.Equal(t, setParams.SecretAttributes.Enabled, setResp.Attributes.Enabled)
	require.Equal(t, setParams.SecretAttributes.NotBefore.Unix(), setResp.Attributes.NotBefore.Unix())
	require.Equal(t, setParams.Tags, setResp.Tags)
	require.Equal(t, setParams.Value, setResp.Value)

	// TODO: would be nice to demote version
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

	updateParams := azsecrets.UpdateSecretParameters{
		SecretAttributes: &azsecrets.SecretAttributes{
			Expires: to.Ptr(time.Date(2040, 1, 1, 1, 1, 1, 0, time.UTC)),
		},
	}
	updateResp, err := client.UpdateSecret(context.Background(), name, setResp.ID.Version(), updateParams, nil)
	require.NoError(t, err)
	require.Equal(t, setParams.ContentType, updateResp.ContentType)
	require.Equal(t, setResp.ID, updateResp.ID)
	require.Equal(t, setParams.SecretAttributes.Enabled, updateResp.Attributes.Enabled)
	require.Equal(t, setParams.SecretAttributes.NotBefore.Unix(), updateResp.Attributes.NotBefore.Unix())
	require.Equal(t, setParams.Tags, updateResp.Tags)

	deleteResp, err := client.DeleteSecret(context.Background(), name, nil)
	require.NoError(t, err)
	require.Equal(t, setParams.ContentType, deleteResp.ContentType)
	require.Equal(t, setResp.ID, deleteResp.ID)
	require.Equal(t, setParams.SecretAttributes.Enabled, deleteResp.Attributes.Enabled)
	require.Equal(t, updateParams.SecretAttributes.Expires.Unix(), deleteResp.Attributes.Expires.Unix())
	require.Equal(t, setParams.SecretAttributes.NotBefore.Unix(), deleteResp.Attributes.NotBefore.Unix())
	require.Equal(t, setParams.Tags, deleteResp.Tags)
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

	_, err = client.PurgeDeletedSecret(context.Background(), name, nil)
	require.NoError(t, err)
}

func TestGetDeletedSecrets(t *testing.T) {
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
	// TODO: Get vs. List
	pager := client.NewGetDeletedSecretsPager(&azsecrets.GetDeletedSecretsOptions{MaxResults: to.Ptr(int32(1))})
	for pager.More() && len(expected) > 0 {
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		for _, secret := range page.Value {
			delete(expected, secret.ID.Name())
			if len(expected) == 0 {
				break
			}
		}
	}
	require.Empty(t, expected, "pager didn't return all expected secrets")
}

func TestGetSecrets(t *testing.T) {
	client := startTest(t)

	count := 4
	for i := 0; i < count; i++ {
		name := createRandomName(t, fmt.Sprintf("listsecrets%d", i))
		value := createRandomName(t, fmt.Sprintf("value%d", i))
		_, err := client.SetSecret(context.Background(), name, azsecrets.SetSecretParameters{Value: &value}, nil)
		require.NoError(t, err)
		defer cleanUpSecret(t, client, name)
	}

	// TODO: Get vs. List
	pager := client.NewGetSecretsPager(&azsecrets.GetSecretsOptions{MaxResults: to.Ptr(int32(1))})
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		for _, secret := range page.Value {
			if strings.HasPrefix(secret.ID.Name(), "listsecrets") {
				count--
			}
		}
	}
	require.Equal(t, count, 0)
}

func TestGetSecretVersions(t *testing.T) {
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

	// TODO: Get vs. List
	pager := client.NewGetSecretVersionsPager(name, &azsecrets.GetSecretVersionsOptions{MaxResults: to.Ptr(int32(1))})
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		for _, secret := range page.Value {
			require.NotNil(t, secret.ID)
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
	client := azsecrets.NewClient(fakeVaultURL, &FakeCredential{}, nil)
	expected := "parameter name cannot be empty"
	_, err := client.BackupSecret(context.Background(), "", nil)
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
	_, err = client.UpdateSecret(context.Background(), "", "", azsecrets.UpdateSecretParameters{}, nil)
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
