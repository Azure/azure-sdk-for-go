//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armpolicyinsights_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/policyinsights/armpolicyinsights"
)

// x-ms-original-file: specification/policyinsights/resource-manager/Microsoft.PolicyInsights/stable/2021-10-01/examples/Remediations_ListDeploymentsManagementGroupScope.json
func ExampleRemediationsClient_ListDeploymentsAtManagementGroup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpolicyinsights.NewRemediationsClient("<subscription-id>", cred, nil)
	pager := client.ListDeploymentsAtManagementGroup(armpolicyinsights.Enum0("Microsoft.Management"),
		"<management-group-id>",
		"<remediation-name>",
		&armpolicyinsights.QueryOptions{Top: nil,
			Filter:    nil,
			OrderBy:   nil,
			Select:    nil,
			From:      nil,
			To:        nil,
			Apply:     nil,
			SkipToken: nil,
			Expand:    nil,
		})
	for {
		nextResult := pager.NextPage(ctx)
		if err := pager.Err(); err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		if !nextResult {
			break
		}
		for _, v := range pager.PageResponse().Value {
			log.Printf("Pager result: %#v\n", v)
		}
	}
}

// x-ms-original-file: specification/policyinsights/resource-manager/Microsoft.PolicyInsights/stable/2021-10-01/examples/Remediations_CancelManagementGroupScope.json
func ExampleRemediationsClient_CancelAtManagementGroup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpolicyinsights.NewRemediationsClient("<subscription-id>", cred, nil)
	res, err := client.CancelAtManagementGroup(ctx,
		armpolicyinsights.Enum0("Microsoft.Management"),
		"<management-group-id>",
		"<remediation-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.RemediationsClientCancelAtManagementGroupResult)
}

// x-ms-original-file: specification/policyinsights/resource-manager/Microsoft.PolicyInsights/stable/2021-10-01/examples/Remediations_ListManagementGroupScope.json
func ExampleRemediationsClient_ListForManagementGroup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpolicyinsights.NewRemediationsClient("<subscription-id>", cred, nil)
	pager := client.ListForManagementGroup(armpolicyinsights.Enum0("Microsoft.Management"),
		"<management-group-id>",
		&armpolicyinsights.QueryOptions{Top: nil,
			Filter:    nil,
			OrderBy:   nil,
			Select:    nil,
			From:      nil,
			To:        nil,
			Apply:     nil,
			SkipToken: nil,
			Expand:    nil,
		})
	for {
		nextResult := pager.NextPage(ctx)
		if err := pager.Err(); err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		if !nextResult {
			break
		}
		for _, v := range pager.PageResponse().Value {
			log.Printf("Pager result: %#v\n", v)
		}
	}
}

// x-ms-original-file: specification/policyinsights/resource-manager/Microsoft.PolicyInsights/stable/2021-10-01/examples/Remediations_CreateManagementGroupScope.json
func ExampleRemediationsClient_CreateOrUpdateAtManagementGroup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpolicyinsights.NewRemediationsClient("<subscription-id>", cred, nil)
	res, err := client.CreateOrUpdateAtManagementGroup(ctx,
		armpolicyinsights.Enum0("Microsoft.Management"),
		"<management-group-id>",
		"<remediation-name>",
		armpolicyinsights.Remediation{
			Properties: &armpolicyinsights.RemediationProperties{
				PolicyAssignmentID: to.StringPtr("<policy-assignment-id>"),
			},
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.RemediationsClientCreateOrUpdateAtManagementGroupResult)
}

// x-ms-original-file: specification/policyinsights/resource-manager/Microsoft.PolicyInsights/stable/2021-10-01/examples/Remediations_GetManagementGroupScope.json
func ExampleRemediationsClient_GetAtManagementGroup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpolicyinsights.NewRemediationsClient("<subscription-id>", cred, nil)
	res, err := client.GetAtManagementGroup(ctx,
		armpolicyinsights.Enum0("Microsoft.Management"),
		"<management-group-id>",
		"<remediation-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.RemediationsClientGetAtManagementGroupResult)
}

// x-ms-original-file: specification/policyinsights/resource-manager/Microsoft.PolicyInsights/stable/2021-10-01/examples/Remediations_DeleteManagementGroupScope.json
func ExampleRemediationsClient_DeleteAtManagementGroup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpolicyinsights.NewRemediationsClient("<subscription-id>", cred, nil)
	res, err := client.DeleteAtManagementGroup(ctx,
		armpolicyinsights.Enum0("Microsoft.Management"),
		"<management-group-id>",
		"<remediation-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.RemediationsClientDeleteAtManagementGroupResult)
}

// x-ms-original-file: specification/policyinsights/resource-manager/Microsoft.PolicyInsights/stable/2021-10-01/examples/Remediations_ListDeploymentsSubscriptionScope.json
func ExampleRemediationsClient_ListDeploymentsAtSubscription() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpolicyinsights.NewRemediationsClient("<subscription-id>", cred, nil)
	pager := client.ListDeploymentsAtSubscription("<remediation-name>",
		&armpolicyinsights.QueryOptions{Top: nil,
			Filter:    nil,
			OrderBy:   nil,
			Select:    nil,
			From:      nil,
			To:        nil,
			Apply:     nil,
			SkipToken: nil,
			Expand:    nil,
		})
	for {
		nextResult := pager.NextPage(ctx)
		if err := pager.Err(); err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		if !nextResult {
			break
		}
		for _, v := range pager.PageResponse().Value {
			log.Printf("Pager result: %#v\n", v)
		}
	}
}

// x-ms-original-file: specification/policyinsights/resource-manager/Microsoft.PolicyInsights/stable/2021-10-01/examples/Remediations_CancelSubscriptionScope.json
func ExampleRemediationsClient_CancelAtSubscription() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpolicyinsights.NewRemediationsClient("<subscription-id>", cred, nil)
	res, err := client.CancelAtSubscription(ctx,
		"<remediation-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.RemediationsClientCancelAtSubscriptionResult)
}

// x-ms-original-file: specification/policyinsights/resource-manager/Microsoft.PolicyInsights/stable/2021-10-01/examples/Remediations_ListSubscriptionScope.json
func ExampleRemediationsClient_ListForSubscription() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpolicyinsights.NewRemediationsClient("<subscription-id>", cred, nil)
	pager := client.ListForSubscription(&armpolicyinsights.QueryOptions{Top: nil,
		Filter:    nil,
		OrderBy:   nil,
		Select:    nil,
		From:      nil,
		To:        nil,
		Apply:     nil,
		SkipToken: nil,
		Expand:    nil,
	})
	for {
		nextResult := pager.NextPage(ctx)
		if err := pager.Err(); err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		if !nextResult {
			break
		}
		for _, v := range pager.PageResponse().Value {
			log.Printf("Pager result: %#v\n", v)
		}
	}
}

// x-ms-original-file: specification/policyinsights/resource-manager/Microsoft.PolicyInsights/stable/2021-10-01/examples/Remediations_CreateSubscriptionScope.json
func ExampleRemediationsClient_CreateOrUpdateAtSubscription() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpolicyinsights.NewRemediationsClient("<subscription-id>", cred, nil)
	res, err := client.CreateOrUpdateAtSubscription(ctx,
		"<remediation-name>",
		armpolicyinsights.Remediation{
			Properties: &armpolicyinsights.RemediationProperties{
				PolicyAssignmentID: to.StringPtr("<policy-assignment-id>"),
			},
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.RemediationsClientCreateOrUpdateAtSubscriptionResult)
}

// x-ms-original-file: specification/policyinsights/resource-manager/Microsoft.PolicyInsights/stable/2021-10-01/examples/Remediations_GetSubscriptionScope.json
func ExampleRemediationsClient_GetAtSubscription() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpolicyinsights.NewRemediationsClient("<subscription-id>", cred, nil)
	res, err := client.GetAtSubscription(ctx,
		"<remediation-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.RemediationsClientGetAtSubscriptionResult)
}

// x-ms-original-file: specification/policyinsights/resource-manager/Microsoft.PolicyInsights/stable/2021-10-01/examples/Remediations_DeleteSubscriptionScope.json
func ExampleRemediationsClient_DeleteAtSubscription() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpolicyinsights.NewRemediationsClient("<subscription-id>", cred, nil)
	res, err := client.DeleteAtSubscription(ctx,
		"<remediation-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.RemediationsClientDeleteAtSubscriptionResult)
}

// x-ms-original-file: specification/policyinsights/resource-manager/Microsoft.PolicyInsights/stable/2021-10-01/examples/Remediations_ListDeploymentsResourceGroupScope.json
func ExampleRemediationsClient_ListDeploymentsAtResourceGroup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpolicyinsights.NewRemediationsClient("<subscription-id>", cred, nil)
	pager := client.ListDeploymentsAtResourceGroup("<resource-group-name>",
		"<remediation-name>",
		&armpolicyinsights.QueryOptions{Top: nil,
			Filter:    nil,
			OrderBy:   nil,
			Select:    nil,
			From:      nil,
			To:        nil,
			Apply:     nil,
			SkipToken: nil,
			Expand:    nil,
		})
	for {
		nextResult := pager.NextPage(ctx)
		if err := pager.Err(); err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		if !nextResult {
			break
		}
		for _, v := range pager.PageResponse().Value {
			log.Printf("Pager result: %#v\n", v)
		}
	}
}

// x-ms-original-file: specification/policyinsights/resource-manager/Microsoft.PolicyInsights/stable/2021-10-01/examples/Remediations_CancelResourceGroupScope.json
func ExampleRemediationsClient_CancelAtResourceGroup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpolicyinsights.NewRemediationsClient("<subscription-id>", cred, nil)
	res, err := client.CancelAtResourceGroup(ctx,
		"<resource-group-name>",
		"<remediation-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.RemediationsClientCancelAtResourceGroupResult)
}

// x-ms-original-file: specification/policyinsights/resource-manager/Microsoft.PolicyInsights/stable/2021-10-01/examples/Remediations_ListResourceGroupScope.json
func ExampleRemediationsClient_ListForResourceGroup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpolicyinsights.NewRemediationsClient("<subscription-id>", cred, nil)
	pager := client.ListForResourceGroup("<resource-group-name>",
		&armpolicyinsights.QueryOptions{Top: nil,
			Filter:    nil,
			OrderBy:   nil,
			Select:    nil,
			From:      nil,
			To:        nil,
			Apply:     nil,
			SkipToken: nil,
			Expand:    nil,
		})
	for {
		nextResult := pager.NextPage(ctx)
		if err := pager.Err(); err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		if !nextResult {
			break
		}
		for _, v := range pager.PageResponse().Value {
			log.Printf("Pager result: %#v\n", v)
		}
	}
}

// x-ms-original-file: specification/policyinsights/resource-manager/Microsoft.PolicyInsights/stable/2021-10-01/examples/Remediations_CreateResourceGroupScope.json
func ExampleRemediationsClient_CreateOrUpdateAtResourceGroup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpolicyinsights.NewRemediationsClient("<subscription-id>", cred, nil)
	res, err := client.CreateOrUpdateAtResourceGroup(ctx,
		"<resource-group-name>",
		"<remediation-name>",
		armpolicyinsights.Remediation{
			Properties: &armpolicyinsights.RemediationProperties{
				PolicyAssignmentID: to.StringPtr("<policy-assignment-id>"),
			},
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.RemediationsClientCreateOrUpdateAtResourceGroupResult)
}

// x-ms-original-file: specification/policyinsights/resource-manager/Microsoft.PolicyInsights/stable/2021-10-01/examples/Remediations_GetResourceGroupScope.json
func ExampleRemediationsClient_GetAtResourceGroup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpolicyinsights.NewRemediationsClient("<subscription-id>", cred, nil)
	res, err := client.GetAtResourceGroup(ctx,
		"<resource-group-name>",
		"<remediation-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.RemediationsClientGetAtResourceGroupResult)
}

// x-ms-original-file: specification/policyinsights/resource-manager/Microsoft.PolicyInsights/stable/2021-10-01/examples/Remediations_DeleteResourceGroupScope.json
func ExampleRemediationsClient_DeleteAtResourceGroup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpolicyinsights.NewRemediationsClient("<subscription-id>", cred, nil)
	res, err := client.DeleteAtResourceGroup(ctx,
		"<resource-group-name>",
		"<remediation-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.RemediationsClientDeleteAtResourceGroupResult)
}

// x-ms-original-file: specification/policyinsights/resource-manager/Microsoft.PolicyInsights/stable/2021-10-01/examples/Remediations_ListDeploymentsResourceScope.json
func ExampleRemediationsClient_ListDeploymentsAtResource() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpolicyinsights.NewRemediationsClient("<subscription-id>", cred, nil)
	pager := client.ListDeploymentsAtResource("<resource-id>",
		"<remediation-name>",
		&armpolicyinsights.QueryOptions{Top: nil,
			Filter:    nil,
			OrderBy:   nil,
			Select:    nil,
			From:      nil,
			To:        nil,
			Apply:     nil,
			SkipToken: nil,
			Expand:    nil,
		})
	for {
		nextResult := pager.NextPage(ctx)
		if err := pager.Err(); err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		if !nextResult {
			break
		}
		for _, v := range pager.PageResponse().Value {
			log.Printf("Pager result: %#v\n", v)
		}
	}
}

// x-ms-original-file: specification/policyinsights/resource-manager/Microsoft.PolicyInsights/stable/2021-10-01/examples/Remediations_CancelResourceScope.json
func ExampleRemediationsClient_CancelAtResource() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpolicyinsights.NewRemediationsClient("<subscription-id>", cred, nil)
	res, err := client.CancelAtResource(ctx,
		"<resource-id>",
		"<remediation-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.RemediationsClientCancelAtResourceResult)
}

// x-ms-original-file: specification/policyinsights/resource-manager/Microsoft.PolicyInsights/stable/2021-10-01/examples/Remediations_ListResourceScope.json
func ExampleRemediationsClient_ListForResource() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpolicyinsights.NewRemediationsClient("<subscription-id>", cred, nil)
	pager := client.ListForResource("<resource-id>",
		&armpolicyinsights.QueryOptions{Top: nil,
			Filter:    nil,
			OrderBy:   nil,
			Select:    nil,
			From:      nil,
			To:        nil,
			Apply:     nil,
			SkipToken: nil,
			Expand:    nil,
		})
	for {
		nextResult := pager.NextPage(ctx)
		if err := pager.Err(); err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		if !nextResult {
			break
		}
		for _, v := range pager.PageResponse().Value {
			log.Printf("Pager result: %#v\n", v)
		}
	}
}

// x-ms-original-file: specification/policyinsights/resource-manager/Microsoft.PolicyInsights/stable/2021-10-01/examples/Remediations_CreateResourceScope.json
func ExampleRemediationsClient_CreateOrUpdateAtResource() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpolicyinsights.NewRemediationsClient("<subscription-id>", cred, nil)
	res, err := client.CreateOrUpdateAtResource(ctx,
		"<resource-id>",
		"<remediation-name>",
		armpolicyinsights.Remediation{
			Properties: &armpolicyinsights.RemediationProperties{
				PolicyAssignmentID: to.StringPtr("<policy-assignment-id>"),
			},
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.RemediationsClientCreateOrUpdateAtResourceResult)
}

// x-ms-original-file: specification/policyinsights/resource-manager/Microsoft.PolicyInsights/stable/2021-10-01/examples/Remediations_GetResourceScope.json
func ExampleRemediationsClient_GetAtResource() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpolicyinsights.NewRemediationsClient("<subscription-id>", cred, nil)
	res, err := client.GetAtResource(ctx,
		"<resource-id>",
		"<remediation-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.RemediationsClientGetAtResourceResult)
}

// x-ms-original-file: specification/policyinsights/resource-manager/Microsoft.PolicyInsights/stable/2021-10-01/examples/Remediations_DeleteResourceScope.json
func ExampleRemediationsClient_DeleteAtResource() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpolicyinsights.NewRemediationsClient("<subscription-id>", cred, nil)
	res, err := client.DeleteAtResource(ctx,
		"<resource-id>",
		"<remediation-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.RemediationsClientDeleteAtResourceResult)
}
