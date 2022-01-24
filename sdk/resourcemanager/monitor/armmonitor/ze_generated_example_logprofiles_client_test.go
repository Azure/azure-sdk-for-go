//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmonitor_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
)

// x-ms-original-file: specification/monitor/resource-manager/Microsoft.Insights/stable/2016-03-01/examples/deleteLogProfile.json
func ExampleLogProfilesClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armmonitor.NewLogProfilesClient("<subscription-id>", cred, nil)
	_, err = client.Delete(ctx,
		"<log-profile-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
}

// x-ms-original-file: specification/monitor/resource-manager/Microsoft.Insights/stable/2016-03-01/examples/getLogProfile.json
func ExampleLogProfilesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armmonitor.NewLogProfilesClient("<subscription-id>", cred, nil)
	res, err := client.Get(ctx,
		"<log-profile-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.LogProfilesClientGetResult)
}

// x-ms-original-file: specification/monitor/resource-manager/Microsoft.Insights/stable/2016-03-01/examples/createOrUpdateLogProfile.json
func ExampleLogProfilesClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armmonitor.NewLogProfilesClient("<subscription-id>", cred, nil)
	res, err := client.CreateOrUpdate(ctx,
		"<log-profile-name>",
		armmonitor.LogProfileResource{
			Location: to.StringPtr("<location>"),
			Tags:     map[string]*string{},
			Properties: &armmonitor.LogProfileProperties{
				Categories: []*string{
					to.StringPtr("Write"),
					to.StringPtr("Delete"),
					to.StringPtr("Action")},
				Locations: []*string{
					to.StringPtr("global")},
				RetentionPolicy: &armmonitor.RetentionPolicy{
					Days:    to.Int32Ptr(3),
					Enabled: to.BoolPtr(true),
				},
				ServiceBusRuleID: to.StringPtr("<service-bus-rule-id>"),
				StorageAccountID: to.StringPtr("<storage-account-id>"),
			},
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.LogProfilesClientCreateOrUpdateResult)
}

// x-ms-original-file: specification/monitor/resource-manager/Microsoft.Insights/stable/2016-03-01/examples/patchLogProfile.json
func ExampleLogProfilesClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armmonitor.NewLogProfilesClient("<subscription-id>", cred, nil)
	res, err := client.Update(ctx,
		"<log-profile-name>",
		armmonitor.LogProfileResourcePatch{
			Properties: &armmonitor.LogProfileProperties{
				Categories: []*string{
					to.StringPtr("Write"),
					to.StringPtr("Delete"),
					to.StringPtr("Action")},
				Locations: []*string{
					to.StringPtr("global")},
				RetentionPolicy: &armmonitor.RetentionPolicy{
					Days:    to.Int32Ptr(3),
					Enabled: to.BoolPtr(true),
				},
				ServiceBusRuleID: to.StringPtr("<service-bus-rule-id>"),
				StorageAccountID: to.StringPtr("<storage-account-id>"),
			},
			Tags: map[string]*string{
				"key1": to.StringPtr("value1"),
			},
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.LogProfilesClientUpdateResult)
}

// x-ms-original-file: specification/monitor/resource-manager/Microsoft.Insights/stable/2016-03-01/examples/listLogProfile.json
func ExampleLogProfilesClient_List() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armmonitor.NewLogProfilesClient("<subscription-id>", cred, nil)
	res, err := client.List(ctx,
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.LogProfilesClientListResult)
}
