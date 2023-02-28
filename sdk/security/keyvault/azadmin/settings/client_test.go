//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package settings_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/settings"
	"github.com/stretchr/testify/require"
)

const fakeHsmURL = "https://fakehsm.managedhsm.azure.net/"

var (
	credential azcore.TokenCredential
	hsmURL     string
)

func TestMain(m *testing.M) {
	if recording.GetRecordMode() != recording.PlaybackMode {
		hsmURL = os.Getenv("AZURE_MANAGEDHSM_URL")
	}
	if hsmURL == "" {
		if recording.GetRecordMode() != recording.PlaybackMode {
			panic("no value for AZURE_MANAGEDHSM_URL")
		}
		hsmURL = fakeHsmURL
	}

	err := recording.ResetProxy(nil)
	if err != nil {
		panic(err)
	}
	if recording.GetRecordMode() == recording.PlaybackMode {
		credential = &FakeCredential{}
	} else {
		tenantID := lookupEnvVar("KEYVAULT_TENANT_ID")
		clientID := lookupEnvVar("KEYVAULT_CLIENT_ID")
		secret := lookupEnvVar("KEYVAULT_CLIENT_SECRET")
		credential, err = azidentity.NewClientSecretCredential(tenantID, clientID, secret, nil)
		if err != nil {
			panic(err)
		}
	}
	if recording.GetRecordMode() == recording.RecordingMode {
		err := recording.AddGeneralRegexSanitizer(fakeHsmURL, hsmURL, nil)
		if err != nil {
			panic(err)
		}
		defer func() {
			err := recording.ResetProxy(nil)
			if err != nil {
				panic(err)
			}
		}()
	}
	code := m.Run()
	os.Exit(code)
}

func startRecording(t *testing.T) {
	err := recording.Start(t, "sdk/security/keyvault/azadmin/testdata", nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	})
}

func startSettingsTest(t *testing.T) *settings.Client {
	startRecording(t)
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	opts := &settings.ClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}}
	client, err := settings.NewClient(hsmURL, credential, opts)
	require.NoError(t, err)
	return client
}

func lookupEnvVar(s string) string {
	ret, ok := os.LookupEnv(s)
	if !ok {
		panic(fmt.Sprintf("Could not find env var: '%s'", s))
	}
	return ret
}

type FakeCredential struct{}

func (f *FakeCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "faketoken", ExpiresOn: time.Now().Add(time.Hour).UTC()}, nil
}

type serdeModel interface {
	json.Marshaler
	json.Unmarshaler
}

func testSerde[T serdeModel](t *testing.T, model T) {
	data, err := model.MarshalJSON()
	require.NoError(t, err)
	err = model.UnmarshalJSON(data)
	require.NoError(t, err)

	// testing unmarshal error scenarios
	var data2 []byte
	err = model.UnmarshalJSON(data2)
	require.Error(t, err)

	m := regexp.MustCompile(":.*$")
	modifiedData := m.ReplaceAllString(string(data), ":false}")
	if modifiedData != "{}" {
		data3 := []byte(modifiedData)
		err = model.UnmarshalJSON(data3)
		require.Error(t, err)
	}
}

func TestGetSetting(t *testing.T) {
	client := startSettingsTest(t)
	settingName := "AllowKeyManagementOperationsThroughARM"

	res, err := client.GetSetting(context.Background(), settingName, nil)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, *res.Name, settingName)
	require.Equal(t, *res.Type, settings.SettingTypeEnumBoolean)
	require.NotNil(t, res.Value)
	testSerde(t, &res)
}

func TestGetSetting_InvalidSettingName(t *testing.T) {
	client := startSettingsTest(t)

	res, err := client.GetSetting(context.Background(), "", nil)
	require.Error(t, err, "parameter settingName cannot be empty")
	require.Nil(t, res.Name)
	require.Nil(t, res.Type)
	require.Nil(t, res.Value)

	res, err = client.GetSetting(context.Background(), "invalid name", nil)
	require.Error(t, err)
	require.Nil(t, res.Name)
	require.Nil(t, res.Type)
	require.Nil(t, res.Value)
	var httpErr *azcore.ResponseError
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, "UnknownError", httpErr.ErrorCode)
	require.Equal(t, 500, httpErr.StatusCode)
}

func TestGetSettings(t *testing.T) {
	client := startSettingsTest(t)

	res, err := client.GetSettings(context.Background(), nil)
	require.NoError(t, err)
	require.NotNil(t, res)

	for _, setting := range res.Settings {
		require.NotNil(t, setting.Name)
		require.NotNil(t, setting.Type)
		require.NotNil(t, setting.Value)
	}

	testSerde(t, &res)

}

func TestUpdateSetting(t *testing.T) {
	client := startSettingsTest(t)
	settingName := "AllowKeyManagementOperationsThroughARM"
	var updatedBool string

	res, err := client.GetSetting(context.Background(), settingName, nil)
	require.NoError(t, err)
	_ = res

	if *res.Value == "true" {
		updatedBool = "false"
	} else if *res.Value == "false" {
		updatedBool = "true"
	}

	updateSettingRequest := settings.UpdateSettingRequest{Value: to.Ptr(updatedBool)}
	testSerde(t, &updateSettingRequest)

	update, err := client.UpdateSetting(context.Background(), settingName, updateSettingRequest, nil)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, settingName, *res.Name)
	require.Equal(t, settings.SettingTypeEnumBoolean, *res.Type)

	require.NotEqual(t, res.Value, update.Value)
	_ = update
}

func TestUpdateSetting_InvalidSettingName(t *testing.T) {
	client := startSettingsTest(t)

	res, err := client.UpdateSetting(context.Background(), "", settings.UpdateSettingRequest{Value: to.Ptr("true")}, nil)
	require.Error(t, err, "parameter settingName cannot be empty")
	require.Nil(t, res.Name)
	require.Nil(t, res.Type)
	require.Nil(t, res.Value)

	res, err = client.UpdateSetting(context.Background(), "invalid name", settings.UpdateSettingRequest{Value: to.Ptr("true")}, nil)
	require.Error(t, err)
	require.Nil(t, res.Name)
	require.Nil(t, res.Type)
	require.Nil(t, res.Value)
	var httpErr *azcore.ResponseError
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, "Nocontentprovided", httpErr.ErrorCode)
	require.Equal(t, 400, httpErr.StatusCode)
}

func TestUpdateSetting_InvalidUpdateSettingRequest(t *testing.T) {
	client := startSettingsTest(t)

	res, err := client.UpdateSetting(context.Background(), "AllowKeyManagementOperationsThroughARM", settings.UpdateSettingRequest{Value: to.Ptr("invalid")}, nil)
	require.Error(t, err)
	require.Nil(t, res.Name)
	require.Nil(t, res.Type)
	require.Nil(t, res.Value)
	var httpErr *azcore.ResponseError
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, "Nocontentprovided", httpErr.ErrorCode)
	require.Equal(t, 400, httpErr.StatusCode)
}
