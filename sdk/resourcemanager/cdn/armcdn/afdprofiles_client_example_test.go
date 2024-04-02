//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armcdn_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cdn/armcdn/v2"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/9ac34f238dd6b9071f486b57e9f9f1a0c43ec6f6/specification/cdn/resource-manager/Microsoft.Cdn/stable/2024-02-01/examples/AFDProfiles_CheckEndpointNameAvailability.json
func ExampleAFDProfilesClient_CheckEndpointNameAvailability() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcdn.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewAFDProfilesClient().CheckEndpointNameAvailability(ctx, "myResourceGroup", "profile1", armcdn.CheckEndpointNameAvailabilityInput{
		Name:                              to.Ptr("sampleName"),
		Type:                              to.Ptr(armcdn.ResourceTypeMicrosoftCdnProfilesAfdEndpoints),
		AutoGeneratedDomainNameLabelScope: to.Ptr(armcdn.AutoGeneratedDomainNameLabelScopeTenantReuse),
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.CheckEndpointNameAvailabilityOutput = armcdn.CheckEndpointNameAvailabilityOutput{
	// 	AvailableHostname: to.Ptr(""),
	// 	Message: to.Ptr("Name not available"),
	// 	NameAvailable: to.Ptr(false),
	// 	Reason: to.Ptr("Name is already in use"),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/9ac34f238dd6b9071f486b57e9f9f1a0c43ec6f6/specification/cdn/resource-manager/Microsoft.Cdn/stable/2024-02-01/examples/AFDProfiles_ListResourceUsage.json
func ExampleAFDProfilesClient_NewListResourceUsagePager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcdn.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewAFDProfilesClient().NewListResourceUsagePager("RG", "profile1", nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range page.Value {
			// You could use page here. We use blank identifier for just demo purposes.
			_ = v
		}
		// If the HTTP response code is 200 as defined in example definition, your page structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
		// page.UsagesListResult = armcdn.UsagesListResult{
		// 	Value: []*armcdn.Usage{
		// 		{
		// 			Name: &armcdn.UsageName{
		// 				LocalizedValue: to.Ptr("afdendpoint"),
		// 				Value: to.Ptr("afdendpoint"),
		// 			},
		// 			CurrentValue: to.Ptr[int64](0),
		// 			ID: to.Ptr("/subscriptions/subid/resourcegroups/RG/providers/Microsoft.Cdn/profiles/profile1/afdendpoints/endpoint1"),
		// 			Limit: to.Ptr[int64](25),
		// 			Unit: to.Ptr(armcdn.UsageUnitCount),
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/9ac34f238dd6b9071f486b57e9f9f1a0c43ec6f6/specification/cdn/resource-manager/Microsoft.Cdn/stable/2024-02-01/examples/AFDProfiles_CheckHostNameAvailability.json
func ExampleAFDProfilesClient_CheckHostNameAvailability() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcdn.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewAFDProfilesClient().CheckHostNameAvailability(ctx, "RG", "profile1", armcdn.CheckHostNameAvailabilityInput{
		HostName: to.Ptr("www.someDomain.net"),
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.CheckNameAvailabilityOutput = armcdn.CheckNameAvailabilityOutput{
	// 	Message: to.Ptr("The hostname 'www.someDomain.net' is already owned by another profile."),
	// 	NameAvailable: to.Ptr(false),
	// 	Reason: to.Ptr("Conflict"),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/9ac34f238dd6b9071f486b57e9f9f1a0c43ec6f6/specification/cdn/resource-manager/Microsoft.Cdn/stable/2024-02-01/examples/AFDProfiles_ValidateSecret.json
func ExampleAFDProfilesClient_ValidateSecret() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcdn.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewAFDProfilesClient().ValidateSecret(ctx, "RG", "profile1", armcdn.ValidateSecretInput{
		SecretSource: &armcdn.ResourceReference{
			ID: to.Ptr("/subscriptions/subid/resourcegroups/RG/providers/Microsoft.KeyVault/vault/kvName/certificate/certName"),
		},
		SecretType: to.Ptr(armcdn.SecretTypeCustomerCertificate),
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.ValidateSecretOutput = armcdn.ValidateSecretOutput{
	// 	Status: to.Ptr(armcdn.StatusValid),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/9ac34f238dd6b9071f486b57e9f9f1a0c43ec6f6/specification/cdn/resource-manager/Microsoft.Cdn/stable/2024-02-01/examples/AFDProfiles_Upgrade.json
func ExampleAFDProfilesClient_BeginUpgrade() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcdn.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewAFDProfilesClient().BeginUpgrade(ctx, "RG", "profile1", armcdn.ProfileUpgradeParameters{
		WafMappingList: []*armcdn.ProfileChangeSKUWafMapping{
			{
				ChangeToWafPolicy: &armcdn.ResourceReference{
					ID: to.Ptr("/subscriptions/subid/resourcegroups/RG/providers/Microsoft.Network/frontdoorwebapplicationfirewallpolicies/waf2"),
				},
				SecurityPolicyName: to.Ptr("securityPolicy1"),
			}},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Profile = armcdn.Profile{
	// 	Name: to.Ptr("profile1"),
	// 	Type: to.Ptr("Microsoft.Cdn/profiles"),
	// 	ID: to.Ptr("/subscriptions/subid/resourcegroups/RG/providers/Microsoft.Cdn/profiles/profile1"),
	// 	Location: to.Ptr("Global"),
	// 	Tags: map[string]*string{
	// 	},
	// 	Kind: to.Ptr("frontdoor"),
	// 	Properties: &armcdn.ProfileProperties{
	// 		ExtendedProperties: map[string]*string{
	// 		},
	// 		FrontDoorID: to.Ptr("id"),
	// 		OriginResponseTimeoutSeconds: to.Ptr[int32](60),
	// 		ProvisioningState: to.Ptr(armcdn.ProfileProvisioningStateSucceeded),
	// 		ResourceState: to.Ptr(armcdn.ProfileResourceState("Enabled")),
	// 	},
	// 	SKU: &armcdn.SKU{
	// 		Name: to.Ptr(armcdn.SKUNameStandardAzureFrontDoor),
	// 	},
	// }
}
