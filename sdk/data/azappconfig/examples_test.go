//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig_test

import (
	"context"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig"
)

func ExampleNewClient() {
	credential, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	client, err := azappconfig.NewClient("https://my-app-config.azconfig.io", credential, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	_ = client // ignore

	// Output:
}

func ExampleNewClientFromConnectionString() {
	connectionString := os.Getenv("APPCONFIGURATION_CONNECTION_STRING")
	if connectionString == "" {
		return
	}

	client, err := azappconfig.NewClientFromConnectionString(connectionString, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	_ = client // ignore

	// Output:
}

func ExampleClient_AddSetting() {
	connectionString := os.Getenv("APPCONFIGURATION_CONNECTION_STRING")
	if connectionString == "" {
		return
	}

	client, err := azappconfig.NewClientFromConnectionString(connectionString, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	// Create configuration setting
	resp, err := client.AddSetting(context.TODO(), "example-key", to.Ptr("example-value"), &azappconfig.AddSettingOptions{
		Label: to.Ptr("example-label"),
	})

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	_ = resp // TODO: do something with resp

	// Output:
}

func ExampleClient_GetSetting() {
	connectionString := os.Getenv("APPCONFIGURATION_CONNECTION_STRING")
	if connectionString == "" {
		return
	}

	client, err := azappconfig.NewClientFromConnectionString(connectionString, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	// Get configuration setting
	resp, err := client.GetSetting(context.TODO(), "example-key", &azappconfig.GetSettingOptions{
		Label: to.Ptr("example-label"),
	})

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	_ = resp // TODO: do something with resp

	// Output:
}

func ExampleClient_SetSetting() {
	connectionString := os.Getenv("APPCONFIGURATION_CONNECTION_STRING")
	if connectionString == "" {
		return
	}

	client, err := azappconfig.NewClientFromConnectionString(connectionString, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	// Set configuration setting
	resp, err := client.SetSetting(context.TODO(), "example-key", to.Ptr("example-new-value"), &azappconfig.SetSettingOptions{
		Label: to.Ptr("example-label"),
	})

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	_ = resp // TODO: do something with resp

	// Output:
}

func ExampleClient_SetReadOnly() {
	connectionString := os.Getenv("APPCONFIGURATION_CONNECTION_STRING")
	if connectionString == "" {
		return
	}

	client, err := azappconfig.NewClientFromConnectionString(connectionString, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	// Set configuration setting read only
	resp, err := client.SetReadOnly(context.TODO(), "example-key", true, &azappconfig.SetReadOnlyOptions{
		Label: to.Ptr("example-label"),
	})

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	_ = resp // TODO: do something with resp

	// Remove read only status
	resp, err = client.SetReadOnly(context.TODO(), "example-key", false, &azappconfig.SetReadOnlyOptions{
		Label: to.Ptr("example-label"),
	})

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	_ = resp // TODO: do something with resp

	// Output:
}

func ExampleClient_NewListRevisionsPager() {
	connectionString := os.Getenv("APPCONFIGURATION_CONNECTION_STRING")
	if connectionString == "" {
		return
	}

	client, err := azappconfig.NewClientFromConnectionString(connectionString, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	pager := client.NewListRevisionsPager(azappconfig.SettingSelector{
		KeyFilter:   to.Ptr("*"),
		LabelFilter: to.Ptr("*"),
		Fields:      azappconfig.AllSettingFields(),
	}, nil)

	for pager.More() {
		page, err := pager.NextPage(context.TODO())

		if err != nil {
			//  TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: %s", err)
		}

		for _, setting := range page.Settings {
			// each page contains all of the returned settings revisions that match the provided [SettingSelector]

			_ = setting // ignore
		}
	}

	// Output:
}

func ExampleClient_DeleteSetting() {
	connectionString := os.Getenv("APPCONFIGURATION_CONNECTION_STRING")
	if connectionString == "" {
		return
	}

	client, err := azappconfig.NewClientFromConnectionString(connectionString, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	// Delete configuration setting
	resp, err := client.DeleteSetting(context.TODO(), "example-key", &azappconfig.DeleteSettingOptions{
		Label: to.Ptr("example-label"),
	})

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	_ = resp // TODO: do something with resp

	// Output:
}
