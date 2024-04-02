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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/9ac34f238dd6b9071f486b57e9f9f1a0c43ec6f6/specification/cdn/resource-manager/Microsoft.Cdn/stable/2024-02-01/examples/AFDCustomDomains_ListByProfile.json
func ExampleAFDCustomDomainsClient_NewListByProfilePager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcdn.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewAFDCustomDomainsClient().NewListByProfilePager("RG", "profile1", nil)
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
		// page.AFDDomainListResult = armcdn.AFDDomainListResult{
		// 	Value: []*armcdn.AFDDomain{
		// 		{
		// 			Name: to.Ptr("domain1"),
		// 			Type: to.Ptr("Microsoft.Cdn/profiles/customdomains"),
		// 			ID: to.Ptr("/subscriptions/subid/resourcegroups/RG/providers/Microsoft.Cdn/profiles/profile1/customdomains/domain1"),
		// 			Properties: &armcdn.AFDDomainProperties{
		// 				AzureDNSZone: &armcdn.ResourceReference{
		// 					ID: to.Ptr(""),
		// 				},
		// 				PreValidatedCustomDomainResourceID: &armcdn.ResourceReference{
		// 					ID: to.Ptr(""),
		// 				},
		// 				ProfileName: to.Ptr("profile1"),
		// 				TLSSettings: &armcdn.AFDDomainHTTPSParameters{
		// 					CertificateType: to.Ptr(armcdn.AfdCertificateTypeManagedCertificate),
		// 					MinimumTLSVersion: to.Ptr(armcdn.AfdMinimumTLSVersionTLS12),
		// 					Secret: &armcdn.ResourceReference{
		// 						ID: to.Ptr(""),
		// 					},
		// 				},
		// 				DeploymentStatus: to.Ptr(armcdn.DeploymentStatusNotStarted),
		// 				ProvisioningState: to.Ptr(armcdn.AfdProvisioningStateSucceeded),
		// 				DomainValidationState: to.Ptr(armcdn.DomainValidationStatePending),
		// 				HostName: to.Ptr("www.contoso.com"),
		// 				ValidationProperties: &armcdn.DomainValidationProperties{
		// 					ExpirationDate: to.Ptr("2009-06-15T13:45:43.0000000Z"),
		// 					ValidationToken: to.Ptr("8c9912db-c615-4eeb-8465"),
		// 				},
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/9ac34f238dd6b9071f486b57e9f9f1a0c43ec6f6/specification/cdn/resource-manager/Microsoft.Cdn/stable/2024-02-01/examples/AFDCustomDomains_Get.json
func ExampleAFDCustomDomainsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcdn.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewAFDCustomDomainsClient().Get(ctx, "RG", "profile1", "domain1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.AFDDomain = armcdn.AFDDomain{
	// 	Name: to.Ptr("domain1"),
	// 	Type: to.Ptr("Microsoft.Cdn/profiles/customdomains"),
	// 	ID: to.Ptr("/subscriptions/subid/resourcegroups/RG/providers/Microsoft.Cdn/profiles/profile1/customdomains/domain1"),
	// 	Properties: &armcdn.AFDDomainProperties{
	// 		AzureDNSZone: &armcdn.ResourceReference{
	// 			ID: to.Ptr(""),
	// 		},
	// 		PreValidatedCustomDomainResourceID: &armcdn.ResourceReference{
	// 			ID: to.Ptr(""),
	// 		},
	// 		ProfileName: to.Ptr("profile1"),
	// 		TLSSettings: &armcdn.AFDDomainHTTPSParameters{
	// 			CertificateType: to.Ptr(armcdn.AfdCertificateTypeManagedCertificate),
	// 			MinimumTLSVersion: to.Ptr(armcdn.AfdMinimumTLSVersionTLS12),
	// 			Secret: &armcdn.ResourceReference{
	// 				ID: to.Ptr(""),
	// 			},
	// 		},
	// 		DeploymentStatus: to.Ptr(armcdn.DeploymentStatusNotStarted),
	// 		ProvisioningState: to.Ptr(armcdn.AfdProvisioningStateSucceeded),
	// 		DomainValidationState: to.Ptr(armcdn.DomainValidationStatePending),
	// 		HostName: to.Ptr("www.contoso.com"),
	// 		ValidationProperties: &armcdn.DomainValidationProperties{
	// 			ExpirationDate: to.Ptr("2009-06-15T13:45:43.0000000Z"),
	// 			ValidationToken: to.Ptr("8c9912db-c615-4eeb-8465"),
	// 		},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/9ac34f238dd6b9071f486b57e9f9f1a0c43ec6f6/specification/cdn/resource-manager/Microsoft.Cdn/stable/2024-02-01/examples/AFDCustomDomains_Create.json
func ExampleAFDCustomDomainsClient_BeginCreate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcdn.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewAFDCustomDomainsClient().BeginCreate(ctx, "RG", "profile1", "domain1", armcdn.AFDDomain{
		Properties: &armcdn.AFDDomainProperties{
			AzureDNSZone: &armcdn.ResourceReference{
				ID: to.Ptr(""),
			},
			TLSSettings: &armcdn.AFDDomainHTTPSParameters{
				CertificateType:   to.Ptr(armcdn.AfdCertificateTypeManagedCertificate),
				MinimumTLSVersion: to.Ptr(armcdn.AfdMinimumTLSVersionTLS12),
			},
			HostName: to.Ptr("www.someDomain.net"),
		},
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
	// res.AFDDomain = armcdn.AFDDomain{
	// 	Name: to.Ptr("domain1"),
	// 	Type: to.Ptr("Microsoft.Cdn/profiles/customdomains"),
	// 	ID: to.Ptr("/subscriptions/subid/resourcegroups/RG/providers/Microsoft.Cdn/profiles/profile1/customdomains/domain1"),
	// 	Properties: &armcdn.AFDDomainProperties{
	// 		AzureDNSZone: &armcdn.ResourceReference{
	// 			ID: to.Ptr(""),
	// 		},
	// 		PreValidatedCustomDomainResourceID: &armcdn.ResourceReference{
	// 			ID: to.Ptr(""),
	// 		},
	// 		ProfileName: to.Ptr("profile1"),
	// 		TLSSettings: &armcdn.AFDDomainHTTPSParameters{
	// 			CertificateType: to.Ptr(armcdn.AfdCertificateTypeManagedCertificate),
	// 			MinimumTLSVersion: to.Ptr(armcdn.AfdMinimumTLSVersionTLS12),
	// 			Secret: &armcdn.ResourceReference{
	// 				ID: to.Ptr(""),
	// 			},
	// 		},
	// 		DeploymentStatus: to.Ptr(armcdn.DeploymentStatusNotStarted),
	// 		ProvisioningState: to.Ptr(armcdn.AfdProvisioningStateSucceeded),
	// 		DomainValidationState: to.Ptr(armcdn.DomainValidationStateSubmitting),
	// 		HostName: to.Ptr("www.contoso.com"),
	// 		ValidationProperties: &armcdn.DomainValidationProperties{
	// 			ExpirationDate: to.Ptr(""),
	// 			ValidationToken: to.Ptr(""),
	// 		},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/9ac34f238dd6b9071f486b57e9f9f1a0c43ec6f6/specification/cdn/resource-manager/Microsoft.Cdn/stable/2024-02-01/examples/AFDCustomDomains_Update.json
func ExampleAFDCustomDomainsClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcdn.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewAFDCustomDomainsClient().BeginUpdate(ctx, "RG", "profile1", "domain1", armcdn.AFDDomainUpdateParameters{
		Properties: &armcdn.AFDDomainUpdatePropertiesParameters{
			AzureDNSZone: &armcdn.ResourceReference{
				ID: to.Ptr(""),
			},
			TLSSettings: &armcdn.AFDDomainHTTPSParameters{
				CertificateType:   to.Ptr(armcdn.AfdCertificateTypeCustomerCertificate),
				MinimumTLSVersion: to.Ptr(armcdn.AfdMinimumTLSVersionTLS12),
			},
		},
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
	// res.AFDDomain = armcdn.AFDDomain{
	// 	Name: to.Ptr("domain1"),
	// 	Type: to.Ptr("Microsoft.Cdn/profiles/customdomains"),
	// 	ID: to.Ptr("/subscriptions/subid/resourcegroups/RG/providers/Microsoft.Cdn/profiles/profile1/customdomains/domain1"),
	// 	Properties: &armcdn.AFDDomainProperties{
	// 		AzureDNSZone: &armcdn.ResourceReference{
	// 			ID: to.Ptr(""),
	// 		},
	// 		PreValidatedCustomDomainResourceID: &armcdn.ResourceReference{
	// 			ID: to.Ptr(""),
	// 		},
	// 		ProfileName: to.Ptr("profile1"),
	// 		TLSSettings: &armcdn.AFDDomainHTTPSParameters{
	// 			CertificateType: to.Ptr(armcdn.AfdCertificateTypeManagedCertificate),
	// 			MinimumTLSVersion: to.Ptr(armcdn.AfdMinimumTLSVersionTLS12),
	// 			Secret: &armcdn.ResourceReference{
	// 				ID: to.Ptr("/subscriptions/subid/resourcegroups/RG/providers/Microsoft.Cdn/profiles/profile1/secrets/mysecert"),
	// 			},
	// 		},
	// 		DeploymentStatus: to.Ptr(armcdn.DeploymentStatusNotStarted),
	// 		ProvisioningState: to.Ptr(armcdn.AfdProvisioningStateSucceeded),
	// 		DomainValidationState: to.Ptr(armcdn.DomainValidationStateApproved),
	// 		HostName: to.Ptr("www.contoso.com"),
	// 		ValidationProperties: &armcdn.DomainValidationProperties{
	// 			ExpirationDate: to.Ptr("2009-06-15T13:45:43.0000000Z"),
	// 			ValidationToken: to.Ptr("8c9912db-c615-4eeb-8465"),
	// 		},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/9ac34f238dd6b9071f486b57e9f9f1a0c43ec6f6/specification/cdn/resource-manager/Microsoft.Cdn/stable/2024-02-01/examples/AFDCustomDomains_Delete.json
func ExampleAFDCustomDomainsClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcdn.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewAFDCustomDomainsClient().BeginDelete(ctx, "RG", "profile1", "domain1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/9ac34f238dd6b9071f486b57e9f9f1a0c43ec6f6/specification/cdn/resource-manager/Microsoft.Cdn/stable/2024-02-01/examples/AFDCustomDomains_RefreshValidationToken.json
func ExampleAFDCustomDomainsClient_BeginRefreshValidationToken() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcdn.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewAFDCustomDomainsClient().BeginRefreshValidationToken(ctx, "RG", "profile1", "domain1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}
