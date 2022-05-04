//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsecurity_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/security/armsecurity"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/security/resource-manager/Microsoft.Security/stable/2020-01-01/examples/ApplicationWhitelistings/GetAdaptiveApplicationControlsSubscription_example.json
func ExampleAdaptiveApplicationControlsClient_List() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armsecurity.NewAdaptiveApplicationControlsClient("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.List(ctx,
		&armsecurity.AdaptiveApplicationControlsClientListOptions{IncludePathRecommendations: to.Ptr(true),
			Summary: to.Ptr(false),
		})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/security/resource-manager/Microsoft.Security/stable/2020-01-01/examples/ApplicationWhitelistings/GetAdaptiveApplicationControlsGroup_example.json
func ExampleAdaptiveApplicationControlsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armsecurity.NewAdaptiveApplicationControlsClient("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Get(ctx,
		"<asc-location>",
		"<group-name>",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/security/resource-manager/Microsoft.Security/stable/2020-01-01/examples/ApplicationWhitelistings/PutAdaptiveApplicationControls_example.json
func ExampleAdaptiveApplicationControlsClient_Put() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armsecurity.NewAdaptiveApplicationControlsClient("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Put(ctx,
		"<asc-location>",
		"<group-name>",
		armsecurity.AdaptiveApplicationControlGroup{
			Properties: &armsecurity.AdaptiveApplicationControlGroupData{
				EnforcementMode: to.Ptr(armsecurity.EnforcementModeAudit),
				PathRecommendations: []*armsecurity.PathRecommendation{
					{
						Type:                to.Ptr(armsecurity.RecommendationType("PublisherSignature")),
						Path:                to.Ptr("<path>"),
						Action:              to.Ptr(armsecurity.RecommendationActionRecommended),
						Common:              to.Ptr(true),
						ConfigurationStatus: to.Ptr(armsecurity.ConfigurationStatusConfigured),
						FileType:            to.Ptr(armsecurity.FileTypeExe),
						PublisherInfo: &armsecurity.PublisherInfo{
							BinaryName:    to.Ptr("<binary-name>"),
							ProductName:   to.Ptr("<product-name>"),
							PublisherName: to.Ptr("<publisher-name>"),
							Version:       to.Ptr("<version>"),
						},
						UserSids: []*string{
							to.Ptr("S-1-1-0")},
						Usernames: []*armsecurity.UserRecommendation{
							{
								RecommendationAction: to.Ptr(armsecurity.RecommendationActionRecommended),
								Username:             to.Ptr("<username>"),
							}},
					},
					{
						Type:                to.Ptr(armsecurity.RecommendationType("ProductSignature")),
						Path:                to.Ptr("<path>"),
						Action:              to.Ptr(armsecurity.RecommendationActionRecommended),
						Common:              to.Ptr(true),
						ConfigurationStatus: to.Ptr(armsecurity.ConfigurationStatusConfigured),
						FileType:            to.Ptr(armsecurity.FileTypeExe),
						PublisherInfo: &armsecurity.PublisherInfo{
							BinaryName:    to.Ptr("<binary-name>"),
							ProductName:   to.Ptr("<product-name>"),
							PublisherName: to.Ptr("<publisher-name>"),
							Version:       to.Ptr("<version>"),
						},
						UserSids: []*string{
							to.Ptr("S-1-1-0")},
						Usernames: []*armsecurity.UserRecommendation{
							{
								RecommendationAction: to.Ptr(armsecurity.RecommendationActionRecommended),
								Username:             to.Ptr("<username>"),
							}},
					},
					{
						Type:                to.Ptr(armsecurity.RecommendationType("PublisherSignature")),
						Path:                to.Ptr("<path>"),
						Action:              to.Ptr(armsecurity.RecommendationActionRecommended),
						Common:              to.Ptr(true),
						ConfigurationStatus: to.Ptr(armsecurity.ConfigurationStatusConfigured),
						FileType:            to.Ptr(armsecurity.FileTypeExe),
						PublisherInfo: &armsecurity.PublisherInfo{
							BinaryName:    to.Ptr("<binary-name>"),
							ProductName:   to.Ptr("<product-name>"),
							PublisherName: to.Ptr("<publisher-name>"),
							Version:       to.Ptr("<version>"),
						},
						UserSids: []*string{
							to.Ptr("S-1-1-0")},
						Usernames: []*armsecurity.UserRecommendation{
							{
								RecommendationAction: to.Ptr(armsecurity.RecommendationActionRecommended),
								Username:             to.Ptr("<username>"),
							}},
					},
					{
						Type:   to.Ptr(armsecurity.RecommendationType("File")),
						Path:   to.Ptr("<path>"),
						Action: to.Ptr(armsecurity.RecommendationActionAdd),
						Common: to.Ptr(true),
					}},
				ProtectionMode: &armsecurity.ProtectionMode{
					Exe:    to.Ptr(armsecurity.EnforcementModeAudit),
					Msi:    to.Ptr(armsecurity.EnforcementModeNone),
					Script: to.Ptr(armsecurity.EnforcementModeNone),
				},
				VMRecommendations: []*armsecurity.VMRecommendation{
					{
						ConfigurationStatus:  to.Ptr(armsecurity.ConfigurationStatusConfigured),
						EnforcementSupport:   to.Ptr(armsecurity.EnforcementSupportSupported),
						RecommendationAction: to.Ptr(armsecurity.RecommendationActionRecommended),
						ResourceID:           to.Ptr("<resource-id>"),
					},
					{
						ConfigurationStatus:  to.Ptr(armsecurity.ConfigurationStatusConfigured),
						EnforcementSupport:   to.Ptr(armsecurity.EnforcementSupportSupported),
						RecommendationAction: to.Ptr(armsecurity.RecommendationActionRecommended),
						ResourceID:           to.Ptr("<resource-id>"),
					}},
			},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/security/resource-manager/Microsoft.Security/stable/2020-01-01/examples/ApplicationWhitelistings/DeleteAdaptiveApplicationControls_example.json
func ExampleAdaptiveApplicationControlsClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armsecurity.NewAdaptiveApplicationControlsClient("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.Delete(ctx,
		"<asc-location>",
		"<group-name>",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}
