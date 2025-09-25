//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig_test

import (
	"context"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
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

func ExampleClient_BeginCreateSnapshot() {
	connectionString := os.Getenv("APPCONFIGURATION_CONNECTION_STRING")
	if connectionString == "" {
		return
	}

	client, err := azappconfig.NewClientFromConnectionString(connectionString, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	snapshotName := "example-snapshot"

	filter := []azappconfig.SettingFilter{
		{
			// TODO: Update the following line with your application specific filter logic
			KeyFilter:   to.Ptr("*"),
			LabelFilter: to.Ptr("*"),
		},
	}

	_, err = client.BeginCreateSnapshot(context.TODO(), snapshotName, filter, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
}

func ExampleClient_ArchiveSnapshot() {
	connectionString := os.Getenv("APPCONFIGURATION_CONNECTION_STRING")
	if connectionString == "" {
		return
	}

	client, err := azappconfig.NewClientFromConnectionString(connectionString, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	snapshotName := "existing-snapshot-example"

	_, err = client.ArchiveSnapshot(context.TODO(), snapshotName, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
}

func ExampleClient_RecoverSnapshot() {
	connectionString := os.Getenv("APPCONFIGURATION_CONNECTION_STRING")
	if connectionString == "" {
		return
	}

	client, err := azappconfig.NewClientFromConnectionString(connectionString, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	snapshotName := "existing-snapshot-example"

	_, err = client.RecoverSnapshot(context.TODO(), snapshotName, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
}

func ExampleClient_NewListSnapshotsPager() {
	connectionString := os.Getenv("APPCONFIGURATION_CONNECTION_STRING")
	if connectionString == "" {
		return
	}

	client, err := azappconfig.NewClientFromConnectionString(connectionString, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	snapshotPager := client.NewListSnapshotsPager(nil)

	for snapshotPager.More() {
		snapshotPage, err := snapshotPager.NextPage(context.TODO())

		if err != nil {
			//  TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: %s", err)
		}

		for _, snapshot := range snapshotPage.Snapshots {
			// TODO: implement your application specific logic here
			_ = snapshot
		}
	}
}

func ExampleClient_NewListSettingsForSnapshotPager() {
	connectionString := os.Getenv("APPCONFIGURATION_CONNECTION_STRING")
	if connectionString == "" {
		return
	}

	client, err := azappconfig.NewClientFromConnectionString(connectionString, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	snapshotName := "existing-snapshot-example"

	snapshotPager := client.NewListSettingsForSnapshotPager(snapshotName, nil)

	for snapshotPager.More() {
		snapshotPage, err := snapshotPager.NextPage(context.TODO())

		if err != nil {
			//  TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: %s", err)
		}

		for _, setting := range snapshotPage.Settings {
			// TODO: implement your application specific logic here
			_ = setting
		}
	}
}

func ExampleClient_GetSnapshot() {
	connectionString := os.Getenv("APPCONFIGURATION_CONNECTION_STRING")
	if connectionString == "" {
		return
	}

	client, err := azappconfig.NewClientFromConnectionString(connectionString, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	snapshotName := "snapshot-example"

	snapshot, err := client.GetSnapshot(context.TODO(), snapshotName, nil)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
	}

	_ = snapshot // TODO: do something with snapshot
}

func ExampleClient_NewListSettingsPager_matchConditions() {
	connectionString := os.Getenv("APPCONFIGURATION_CONNECTION_STRING")
	if connectionString == "" {
		return
	}

	client, err := azappconfig.NewClientFromConnectionString(connectionString, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	// matchConditions will contain an ETag for each page of settings returned
	matchConditions := []azcore.MatchConditions{}

	pager := client.NewListSettingsPager(azappconfig.SettingSelector{}, nil)
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			//  TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: %s", err)
		}

		matchConditions = append(matchConditions, azcore.MatchConditions{
			// filter out any pages that haven't changed since they were last retrieved
			IfNoneMatch: page.ETag,
		})
	}

	pager = client.NewListSettingsPager(azappconfig.SettingSelector{}, &azappconfig.ListSettingsOptions{
		MatchConditions: matchConditions,
	})
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			//  TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: %s", err)
		}

		// if the values per page haven't changed, page.Settings will be empty
		_ = page.Settings
	}
}

func ExampleClient_NewListSettingsPager_usingTags() {
	connectionString := os.Getenv("APPCONFIGURATION_CONNECTION_STRING")
	if connectionString == "" {
		return
	}

	client, err := azappconfig.NewClientFromConnectionString(connectionString, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	// First, create a configuration setting with tags
	_, err = client.AddSetting(context.Background(), "endpoint", to.Ptr("https://beta.endpoint.com"), &azappconfig.AddSettingOptions{
		Label: to.Ptr("beta"),
		Tags: map[string]*string{
			"someKey": to.Ptr("someValue"),
		},
	})
	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	// To gather all the information available for settings grouped by a specific tag,
	// use a setting selector that filters for settings with the "someKey=someValue" tag.
	// This will retrieve all the Configuration Settings in the store that satisfy that condition.
	selector := azappconfig.SettingSelector{
		TagsFilter: []string{"someKey=someValue"},
	}

	pager := client.NewListSettingsPager(selector, nil)
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			//  TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: %s", err)
		}

		for _, setting := range page.Settings {
			// Process each setting that matches the tag filter
			_ = setting // TODO: do something with setting
		}
	}
}
