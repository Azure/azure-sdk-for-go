// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azadmin_test

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azadmin"
)

var settingsClient azadmin.SettingsClient

func ExampleNewSettingsClient() {
	vaultURL := "https://<TODO: your vault name>.managedhsm.azure.net/"
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client, err := azadmin.NewSettingsClient(vaultURL, cred, nil)
	if err != nil {
		// TODO: handle error
	}

	_ = client
}

func ExampleSettingsClient_GetSetting() {
	name := "AllowKeyManagementOperationsThroughARM"
	setting, err := settingsClient.GetSetting(context.TODO(), name, nil)
	if err != nil {
		// TODO: handle error
	}

	fmt.Printf("Setting Name: %s\tSetting Value: %s", *setting.Name, *setting.Value)
}

func ExampleSettingsClient_GetSettings() {
	settings, err := settingsClient.GetSettings(context.TODO(), nil)
	if err != nil {
		// TODO: handle error
	}

	for _, setting := range settings.Settings {
		fmt.Printf("Setting Name: %s\tSetting Value: %s", *setting.Name, *setting.Value)
	}
}

func ExampleSettingsClient_UpdateSetting() {
	name := "AllowKeyManagementOperationsThroughARM"
	parameters := azadmin.UpdateSettingRequest{Value: to.Ptr("true")}

	updatedSetting, err := settingsClient.UpdateSetting(context.TODO(), name, parameters, nil)
	if err != nil {
		// TODO: handle error
	}

	fmt.Printf("Setting Name: %s\tSetting Value: %s", *updatedSetting.Name, *updatedSetting.Value)
}
