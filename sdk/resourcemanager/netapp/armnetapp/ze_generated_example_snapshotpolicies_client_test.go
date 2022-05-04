//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armnetapp_test

import (
	"context"
	"log"

	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/netapp/armnetapp"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/netapp/resource-manager/Microsoft.NetApp/stable/2021-10-01/examples/SnapshotPolicies_List.json
func ExampleSnapshotPoliciesClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armnetapp.NewSnapshotPoliciesClient("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListPager("<resource-group-name>",
		"<account-name>",
		nil)
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
			return
		}
		for _, v := range nextResult.Value {
			// TODO: use page item
			_ = v
		}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/netapp/resource-manager/Microsoft.NetApp/stable/2021-10-01/examples/SnapshotPolicies_Get.json
func ExampleSnapshotPoliciesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armnetapp.NewSnapshotPoliciesClient("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Get(ctx,
		"<resource-group-name>",
		"<account-name>",
		"<snapshot-policy-name>",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/netapp/resource-manager/Microsoft.NetApp/stable/2021-10-01/examples/SnapshotPolicies_Create.json
func ExampleSnapshotPoliciesClient_Create() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armnetapp.NewSnapshotPoliciesClient("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Create(ctx,
		"<resource-group-name>",
		"<account-name>",
		"<snapshot-policy-name>",
		armnetapp.SnapshotPolicy{
			Location: to.Ptr("<location>"),
			Properties: &armnetapp.SnapshotPolicyProperties{
				DailySchedule: &armnetapp.DailySchedule{
					Hour:            to.Ptr[int32](14),
					Minute:          to.Ptr[int32](30),
					SnapshotsToKeep: to.Ptr[int32](4),
				},
				Enabled: to.Ptr(true),
				HourlySchedule: &armnetapp.HourlySchedule{
					Minute:          to.Ptr[int32](50),
					SnapshotsToKeep: to.Ptr[int32](2),
				},
				MonthlySchedule: &armnetapp.MonthlySchedule{
					DaysOfMonth:     to.Ptr("<days-of-month>"),
					Hour:            to.Ptr[int32](14),
					Minute:          to.Ptr[int32](15),
					SnapshotsToKeep: to.Ptr[int32](5),
				},
				WeeklySchedule: &armnetapp.WeeklySchedule{
					Day:             to.Ptr("<day>"),
					Hour:            to.Ptr[int32](14),
					Minute:          to.Ptr[int32](45),
					SnapshotsToKeep: to.Ptr[int32](3),
				},
			},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/netapp/resource-manager/Microsoft.NetApp/stable/2021-10-01/examples/SnapshotPolicies_Update.json
func ExampleSnapshotPoliciesClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armnetapp.NewSnapshotPoliciesClient("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginUpdate(ctx,
		"<resource-group-name>",
		"<account-name>",
		"<snapshot-policy-name>",
		armnetapp.SnapshotPolicyPatch{
			Location: to.Ptr("<location>"),
			Properties: &armnetapp.SnapshotPolicyProperties{
				DailySchedule: &armnetapp.DailySchedule{
					Hour:            to.Ptr[int32](14),
					Minute:          to.Ptr[int32](30),
					SnapshotsToKeep: to.Ptr[int32](4),
				},
				Enabled: to.Ptr(true),
				HourlySchedule: &armnetapp.HourlySchedule{
					Minute:          to.Ptr[int32](50),
					SnapshotsToKeep: to.Ptr[int32](2),
				},
				MonthlySchedule: &armnetapp.MonthlySchedule{
					DaysOfMonth:     to.Ptr("<days-of-month>"),
					Hour:            to.Ptr[int32](14),
					Minute:          to.Ptr[int32](15),
					SnapshotsToKeep: to.Ptr[int32](5),
				},
				WeeklySchedule: &armnetapp.WeeklySchedule{
					Day:             to.Ptr("<day>"),
					Hour:            to.Ptr[int32](14),
					Minute:          to.Ptr[int32](45),
					SnapshotsToKeep: to.Ptr[int32](3),
				},
			},
		},
		&armnetapp.SnapshotPoliciesClientBeginUpdateOptions{ResumeToken: ""})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/netapp/resource-manager/Microsoft.NetApp/stable/2021-10-01/examples/SnapshotPolicies_Delete.json
func ExampleSnapshotPoliciesClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armnetapp.NewSnapshotPoliciesClient("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginDelete(ctx,
		"<resource-group-name>",
		"<account-name>",
		"<snapshot-policy-name>",
		&armnetapp.SnapshotPoliciesClientBeginDeleteOptions{ResumeToken: ""})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/netapp/resource-manager/Microsoft.NetApp/stable/2021-10-01/examples/SnapshotPolicies_ListVolumes.json
func ExampleSnapshotPoliciesClient_ListVolumes() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armnetapp.NewSnapshotPoliciesClient("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.ListVolumes(ctx,
		"<resource-group-name>",
		"<account-name>",
		"<snapshot-policy-name>",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}
