// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package settings_test

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/settings"
)

var client settings.Client

func ExampleNewClient() {
	vaultURL := "https://<TODO: your vault name>.managedhsm.azure.net/"
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client, err := settings.NewClient(vaultURL, cred, nil)
	if err != nil {
		// TODO: handle error
	}

	_ = client
}

func ExampleClient_GetSetting() {
	name := "AllowKeyManagementOperationsThroughARM"
	setting, err := client.GetSetting(context.TODO(), name, nil)
	if err != nil {
		// TODO: handle error
	}

	fmt.Printf("Setting Name: %s\tSetting Value: %s", *setting.Name, *setting.Value)
}

func ExampleClient_GetSettings() {
	settings, err := client.GetSettings(context.TODO(), nil)
	if err != nil {
		// TODO: handle error
	}

	for _, setting := range settings.Settings {
		fmt.Printf("Setting Name: %s\tSetting Value: %s", *setting.Name, *setting.Value)
	}
}

func ExampleClient_UpdateSetting() {
	name := "AllowKeyManagementOperationsThroughARM"
	parameters := settings.UpdateSettingRequest{Value: to.Ptr("true")}

	updatedSetting, err := client.UpdateSetting(context.TODO(), name, parameters, nil)
	if err != nil {
		// TODO: handle error
	}

	fmt.Printf("Setting Name: %s\tSetting Value: %s", *updatedSetting.Name, *updatedSetting.Value)
}
