//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package settings_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/settings"
	"github.com/stretchr/testify/require"
)

func TestGetSetting(t *testing.T) {
	client := startSettingsTest(t)
	settingName := "AllowKeyManagementOperationsThroughARM"

	res, err := client.GetSetting(context.Background(), settingName, nil)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, *res.Name, settingName)
	require.Equal(t, *res.Type, settings.SettingTypeBoolean)
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
	require.Equal(t, settings.SettingTypeBoolean, *res.Type)

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

	for i := 0; i < 4; i++ {
		res, err = client.UpdateSetting(context.Background(), "invalid name", settings.UpdateSettingRequest{Value: to.Ptr("true")}, nil)
		var httpErr *azcore.ResponseError
		// if correct error is returned, break from the loop and check for correctness
		if errors.As(err, &httpErr) && httpErr.StatusCode == 400 {
			break
		}
		// else sleep for 30 seconds and try again
		recording.Sleep(30 * time.Second)
	}
	require.Error(t, err)
	require.Nil(t, res.Name)
	require.Nil(t, res.Type)
	require.Nil(t, res.Value)
	var httpErr *azcore.ResponseError
	require.ErrorAs(t, err, &httpErr)
}
