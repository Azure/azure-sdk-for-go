//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsecurityinsights_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/securityinsights/armsecurityinsights/v2"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/securityinsights/resource-manager/Microsoft.SecurityInsights/preview/2022-05-01-preview/examples/threatintelligence/CreateThreatIntelligence.json
func ExampleThreatIntelligenceIndicatorClient_CreateIndicator() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armsecurityinsights.NewThreatIntelligenceIndicatorClient("bd794837-4d29-4647-9105-6339bfdb4e6a", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.CreateIndicator(ctx,
		"myRg",
		"myWorkspace",
		armsecurityinsights.ThreatIntelligenceIndicatorModel{
			Kind: to.Ptr(armsecurityinsights.ThreatIntelligenceResourceKindEnumIndicator),
			Properties: &armsecurityinsights.ThreatIntelligenceIndicatorProperties{
				Description:        to.Ptr("debugging indicators"),
				Confidence:         to.Ptr[int32](78),
				CreatedByRef:       to.Ptr("contoso@contoso.com"),
				DisplayName:        to.Ptr("new schema"),
				ExternalReferences: []*armsecurityinsights.ThreatIntelligenceExternalReference{},
				GranularMarkings:   []*armsecurityinsights.ThreatIntelligenceGranularMarkingModel{},
				KillChainPhases:    []*armsecurityinsights.ThreatIntelligenceKillChainPhase{},
				Labels:             []*string{},
				Modified:           to.Ptr(""),
				Pattern:            to.Ptr("[url:value = 'https://www.contoso.com']"),
				PatternType:        to.Ptr("url"),
				Revoked:            to.Ptr(false),
				Source:             to.Ptr("Azure Sentinel"),
				ThreatIntelligenceTags: []*string{
					to.Ptr("new schema")},
				ThreatTypes: []*string{
					to.Ptr("compromised")},
				ValidFrom:  to.Ptr("2021-09-15T17:44:00.114052Z"),
				ValidUntil: to.Ptr(""),
			},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/securityinsights/resource-manager/Microsoft.SecurityInsights/preview/2022-05-01-preview/examples/threatintelligence/GetThreatIntelligenceById.json
func ExampleThreatIntelligenceIndicatorClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armsecurityinsights.NewThreatIntelligenceIndicatorClient("bd794837-4d29-4647-9105-6339bfdb4e6a", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Get(ctx,
		"myRg",
		"myWorkspace",
		"e16ef847-962e-d7b6-9c8b-a33e4bd30e47",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/securityinsights/resource-manager/Microsoft.SecurityInsights/preview/2022-05-01-preview/examples/threatintelligence/UpdateThreatIntelligence.json
func ExampleThreatIntelligenceIndicatorClient_Create() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armsecurityinsights.NewThreatIntelligenceIndicatorClient("bd794837-4d29-4647-9105-6339bfdb4e6a", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Create(ctx,
		"myRg",
		"myWorkspace",
		"d9cd6f0b-96b9-3984-17cd-a779d1e15a93",
		armsecurityinsights.ThreatIntelligenceIndicatorModel{
			Kind: to.Ptr(armsecurityinsights.ThreatIntelligenceResourceKindEnumIndicator),
			Properties: &armsecurityinsights.ThreatIntelligenceIndicatorProperties{
				Description:        to.Ptr("debugging indicators"),
				Confidence:         to.Ptr[int32](78),
				CreatedByRef:       to.Ptr("contoso@contoso.com"),
				DisplayName:        to.Ptr("new schema"),
				ExternalReferences: []*armsecurityinsights.ThreatIntelligenceExternalReference{},
				GranularMarkings:   []*armsecurityinsights.ThreatIntelligenceGranularMarkingModel{},
				KillChainPhases:    []*armsecurityinsights.ThreatIntelligenceKillChainPhase{},
				Labels:             []*string{},
				Modified:           to.Ptr(""),
				Pattern:            to.Ptr("[url:value = 'https://www.contoso.com']"),
				PatternType:        to.Ptr("url"),
				Revoked:            to.Ptr(false),
				Source:             to.Ptr("Azure Sentinel"),
				ThreatIntelligenceTags: []*string{
					to.Ptr("new schema")},
				ThreatTypes: []*string{
					to.Ptr("compromised")},
				ValidFrom:  to.Ptr("2020-04-15T17:44:00.114052Z"),
				ValidUntil: to.Ptr(""),
			},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/securityinsights/resource-manager/Microsoft.SecurityInsights/preview/2022-05-01-preview/examples/threatintelligence/DeleteThreatIntelligence.json
func ExampleThreatIntelligenceIndicatorClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armsecurityinsights.NewThreatIntelligenceIndicatorClient("bd794837-4d29-4647-9105-6339bfdb4e6a", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.Delete(ctx,
		"myRg",
		"myWorkspace",
		"d9cd6f0b-96b9-3984-17cd-a779d1e15a93",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/securityinsights/resource-manager/Microsoft.SecurityInsights/preview/2022-05-01-preview/examples/threatintelligence/QueryThreatIntelligence.json
func ExampleThreatIntelligenceIndicatorClient_NewQueryIndicatorsPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armsecurityinsights.NewThreatIntelligenceIndicatorClient("bd794837-4d29-4647-9105-6339bfdb4e6a", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewQueryIndicatorsPager("myRg",
		"myWorkspace",
		armsecurityinsights.ThreatIntelligenceFilteringCriteria{
			MaxConfidence: to.Ptr[int32](80),
			MaxValidUntil: to.Ptr("2021-04-25T17:44:00.114052Z"),
			MinConfidence: to.Ptr[int32](25),
			MinValidUntil: to.Ptr("2021-04-05T17:44:00.114052Z"),
			PageSize:      to.Ptr[int32](100),
			SortBy: []*armsecurityinsights.ThreatIntelligenceSortingCriteria{
				{
					ItemKey:   to.Ptr("lastUpdatedTimeUtc"),
					SortOrder: to.Ptr(armsecurityinsights.ThreatIntelligenceSortingCriteriaEnumDescending),
				}},
			Sources: []*string{
				to.Ptr("Azure Sentinel")},
		},
		nil)
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range nextResult.Value {
			// TODO: use page item
			_ = v
		}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/securityinsights/resource-manager/Microsoft.SecurityInsights/preview/2022-05-01-preview/examples/threatintelligence/AppendTagsThreatIntelligence.json
func ExampleThreatIntelligenceIndicatorClient_AppendTags() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armsecurityinsights.NewThreatIntelligenceIndicatorClient("bd794837-4d29-4647-9105-6339bfdb4e6a", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.AppendTags(ctx,
		"myRg",
		"myWorkspace",
		"d9cd6f0b-96b9-3984-17cd-a779d1e15a93",
		armsecurityinsights.ThreatIntelligenceAppendTags{
			ThreatIntelligenceTags: []*string{
				to.Ptr("tag1"),
				to.Ptr("tag2")},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/securityinsights/resource-manager/Microsoft.SecurityInsights/preview/2022-05-01-preview/examples/threatintelligence/ReplaceTagsThreatIntelligence.json
func ExampleThreatIntelligenceIndicatorClient_ReplaceTags() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armsecurityinsights.NewThreatIntelligenceIndicatorClient("bd794837-4d29-4647-9105-6339bfdb4e6a", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.ReplaceTags(ctx,
		"myRg",
		"myWorkspace",
		"d9cd6f0b-96b9-3984-17cd-a779d1e15a93",
		armsecurityinsights.ThreatIntelligenceIndicatorModel{
			Etag: to.Ptr("\"0000262c-0000-0800-0000-5e9767060000\""),
			Kind: to.Ptr(armsecurityinsights.ThreatIntelligenceResourceKindEnumIndicator),
			Properties: &armsecurityinsights.ThreatIntelligenceIndicatorProperties{
				ThreatIntelligenceTags: []*string{
					to.Ptr("patching tags")},
			},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}
