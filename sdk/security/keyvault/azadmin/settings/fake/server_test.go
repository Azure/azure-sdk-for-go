// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package fake_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/settings"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/settings/fake"
	"github.com/stretchr/testify/require"
)

var (
	settingName   = "settingName"
	settingValue  = "settingValue"
	settingName2  = "settingName2"
	settingValue2 = "settingValue2"
)

func getServer() fake.Server {
	return fake.Server{
		GetSetting: func(ctx context.Context, settingName string, options *settings.GetSettingOptions) (resp azfake.Responder[settings.GetSettingResponse], errResp azfake.ErrorResponder) {
			kvResp := settings.GetSettingResponse{Setting: settings.Setting{
				Name:  to.Ptr(settingName),
				Value: to.Ptr(settingValue),
			}}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		GetSettings: func(ctx context.Context, options *settings.GetSettingsOptions) (resp azfake.Responder[settings.GetSettingsResponse], errResp azfake.ErrorResponder) {
			kvResp := settings.GetSettingsResponse{
				ListResult: settings.ListResult{
					Settings: []*settings.Setting{
						{Name: &settingName, Value: &settingValue},
						{Name: &settingName2, Value: &settingValue2},
					},
				}}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		UpdateSetting: func(ctx context.Context, settingName string, parameters settings.UpdateSettingRequest, options *settings.UpdateSettingOptions) (resp azfake.Responder[settings.UpdateSettingResponse], errResp azfake.ErrorResponder) {
			kvResp := settings.UpdateSettingResponse{Setting: settings.Setting{
				Name:  to.Ptr(settingName),
				Value: parameters.Value,
			}}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
	}
}

func TestServer(t *testing.T) {
	fakeServer := getServer()

	client, err := settings.NewClient("https://fake-vault.vault.azure.net", &azfake.TokenCredential{}, &settings.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewServerTransport(&fakeServer),
		},
	})
	require.NoError(t, err)

	// get setting
	settingResp, err := client.GetSetting(context.TODO(), settingName, nil)
	require.NoError(t, err)
	require.Equal(t, settingName, *settingResp.Name)

	// get settings
	settingsResp, err := client.GetSettings(context.TODO(), nil)
	require.NoError(t, err)
	require.Len(t, settingsResp.Settings, 2)
	require.Equal(t, *settingsResp.Settings[0].Name, settingName)
	require.Equal(t, *settingsResp.Settings[1].Name, settingName2)

	// update setting
	updateResp, err := client.UpdateSetting(context.TODO(), settingName, settings.UpdateSettingRequest{Value: &settingValue2}, nil)
	require.NoError(t, err)
	require.Equal(t, settingName, *updateResp.Name)
	require.Equal(t, settingValue2, *updateResp.Value)
}
