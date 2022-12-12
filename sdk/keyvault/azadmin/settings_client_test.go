//go:build go1.18
// +build go1.18

package azadmin_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azadmin"
	"github.com/stretchr/testify/require"
)

// TODO- fix get settings, currently broken
/*func TestGetSettings(t *testing.T) {
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

}*/

func TestUpdateSetting(t *testing.T) {
	client := startSettingsTest(t)
	settingName := "AllowKeyManagementOperationsThroughARM"
	var updatedBool string

	res, err := client.GetSetting(context.Background(), "AllowKeyManagementOperationsThroughARM", nil)
	require.NoError(t, err)
	_ = res

	if *res.Value == "true" {
		updatedBool = "false"
	} else if *res.Value == "false" {
		updatedBool = "true"
	}

	updateSettingRequest := azadmin.UpdateSettingRequest{Value: to.Ptr(updatedBool)}
	testSerde(t, &updateSettingRequest)

	update, err := client.UpdateSetting(context.Background(), settingName, updateSettingRequest, nil)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, *res.Name, settingName)
	require.Equal(t, *res.Type, azadmin.SettingTypeEnumBoolean)

	require.NotEqual(t, res.Value, update.Value)
	_ = update
}

func TestGetSetting(t *testing.T) {
	client := startSettingsTest(t)
	settingName := "AllowKeyManagementOperationsThroughARM"

	res, err := client.GetSetting(context.Background(), settingName, nil)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, *res.Name, settingName)
	require.Equal(t, *res.Type, azadmin.SettingTypeEnumBoolean)
	require.NotNil(t, res.Value)
	testSerde(t, &res)
}
