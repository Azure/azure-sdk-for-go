//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armbatch_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/batch/armbatch/v3"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e79d9ef3e065f2dcb6bd1db51e29c62a99dff5cb/specification/batch/resource-manager/Microsoft.Batch/stable/2024-07-01/examples/NspConfigurationsList.json
func ExampleNetworkSecurityPerimeterClient_NewListConfigurationsPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armbatch.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewNetworkSecurityPerimeterClient().NewListConfigurationsPager("default-azurebatch-japaneast", "sampleacct", nil)
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
		// page.NetworkSecurityPerimeterConfigurationListResult = armbatch.NetworkSecurityPerimeterConfigurationListResult{
		// 	Value: []*armbatch.NetworkSecurityPerimeterConfiguration{
		// 		{
		// 			Name: to.Ptr("00000000-0000-0000-0000-000000000000.sampleassociation"),
		// 			Type: to.Ptr("Microsoft.Batch/batchAccounts/networkSecurityPerimeterConfigurations"),
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/default-azurebatch-japaneast/providers/Microsoft.Batch/batchAccounts/sampleacct/networkSecurityPerimeterConfigurations/00000000-0000-0000-0000-000000000000.sampleassociation"),
		// 			Properties: &armbatch.NetworkSecurityPerimeterConfigurationProperties{
		// 				NetworkSecurityPerimeter: &armbatch.NetworkSecurityPerimeter{
		// 					ID: to.Ptr("/subscriptions/perimeterSubscriptionId/resourceGroups/perimeterResourceGroupName/providers/Microsoft.Network/networkSecurityPerimeters/perimeterName"),
		// 					Location: to.Ptr("perimeterLocation"),
		// 					PerimeterGUID: to.Ptr("00000000-0000-0000-0000-000000000000"),
		// 				},
		// 				Profile: &armbatch.NetworkSecurityProfile{
		// 					Name: to.Ptr("profileName"),
		// 					AccessRules: []*armbatch.AccessRule{
		// 						{
		// 							Name: to.Ptr("accessRule1"),
		// 							Properties: &armbatch.AccessRuleProperties{
		// 								AddressPrefixes: []*string{
		// 									to.Ptr("10.11.0.0/16"),
		// 									to.Ptr("10.10.1.0/24")},
		// 									Direction: to.Ptr(armbatch.AccessRuleDirectionInbound),
		// 									EmailAddresses: []*string{
		// 									},
		// 									FullyQualifiedDomainNames: []*string{
		// 									},
		// 									NetworkSecurityPerimeters: []*armbatch.NetworkSecurityPerimeter{
		// 									},
		// 									PhoneNumbers: []*string{
		// 									},
		// 									Subscriptions: []*armbatch.AccessRulePropertiesSubscriptionsItem{
		// 									},
		// 								},
		// 							},
		// 							{
		// 								Name: to.Ptr("accessRule2"),
		// 								Properties: &armbatch.AccessRuleProperties{
		// 									AddressPrefixes: []*string{
		// 									},
		// 									Direction: to.Ptr(armbatch.AccessRuleDirectionOutbound),
		// 									EmailAddresses: []*string{
		// 									},
		// 									FullyQualifiedDomainNames: []*string{
		// 										to.Ptr("paasrp1.contoso.org"),
		// 										to.Ptr("paasrp2.contoso.org")},
		// 										NetworkSecurityPerimeters: []*armbatch.NetworkSecurityPerimeter{
		// 										},
		// 										PhoneNumbers: []*string{
		// 										},
		// 										Subscriptions: []*armbatch.AccessRulePropertiesSubscriptionsItem{
		// 										},
		// 									},
		// 							}},
		// 							AccessRulesVersion: to.Ptr[int32](1),
		// 							DiagnosticSettingsVersion: to.Ptr[int32](1),
		// 							EnabledLogCategories: []*string{
		// 								to.Ptr("logCategory1"),
		// 								to.Ptr("logCategory2")},
		// 							},
		// 							ProvisioningIssues: []*armbatch.ProvisioningIssue{
		// 							},
		// 							ProvisioningState: to.Ptr(armbatch.NetworkSecurityPerimeterConfigurationProvisioningStateSucceeded),
		// 							ResourceAssociation: &armbatch.ResourceAssociation{
		// 								Name: to.Ptr("sampleassociation"),
		// 								AccessMode: to.Ptr(armbatch.ResourceAssociationAccessModeLearning),
		// 							},
		// 						},
		// 				}},
		// 			}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e79d9ef3e065f2dcb6bd1db51e29c62a99dff5cb/specification/batch/resource-manager/Microsoft.Batch/stable/2024-07-01/examples/NspConfigurationGet.json
func ExampleNetworkSecurityPerimeterClient_GetConfiguration() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armbatch.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewNetworkSecurityPerimeterClient().GetConfiguration(ctx, "default-azurebatch-japaneast", "sampleacct", "00000000-0000-0000-0000-000000000000.sampleassociation", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.NetworkSecurityPerimeterConfiguration = armbatch.NetworkSecurityPerimeterConfiguration{
	// 	Name: to.Ptr("00000000-0000-0000-0000-000000000000.sampleassociation"),
	// 	Type: to.Ptr("Microsoft.Batch/batchAccounts/networkSecurityPerimeterConfigurations"),
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/default-azurebatch-japaneast/providers/Microsoft.Batch/batchAccounts/sampleacct/networkSecurityPerimeterConfigurations/00000000-0000-0000-0000-000000000000.sampleassociation"),
	// 	Properties: &armbatch.NetworkSecurityPerimeterConfigurationProperties{
	// 		NetworkSecurityPerimeter: &armbatch.NetworkSecurityPerimeter{
	// 			ID: to.Ptr("/subscriptions/perimeterSubscriptionId/resourceGroups/perimeterResourceGroupName/providers/Microsoft.Network/networkSecurityPerimeters/perimeterName"),
	// 			Location: to.Ptr("perimeterLocation"),
	// 			PerimeterGUID: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 		},
	// 		Profile: &armbatch.NetworkSecurityProfile{
	// 			Name: to.Ptr("profileName"),
	// 			AccessRules: []*armbatch.AccessRule{
	// 				{
	// 					Name: to.Ptr("accessRule1"),
	// 					Properties: &armbatch.AccessRuleProperties{
	// 						AddressPrefixes: []*string{
	// 							to.Ptr("10.11.0.0/16"),
	// 							to.Ptr("10.10.1.0/24")},
	// 							Direction: to.Ptr(armbatch.AccessRuleDirectionInbound),
	// 							EmailAddresses: []*string{
	// 							},
	// 							FullyQualifiedDomainNames: []*string{
	// 							},
	// 							NetworkSecurityPerimeters: []*armbatch.NetworkSecurityPerimeter{
	// 							},
	// 							PhoneNumbers: []*string{
	// 							},
	// 							Subscriptions: []*armbatch.AccessRulePropertiesSubscriptionsItem{
	// 							},
	// 						},
	// 					},
	// 					{
	// 						Name: to.Ptr("accessRule2"),
	// 						Properties: &armbatch.AccessRuleProperties{
	// 							AddressPrefixes: []*string{
	// 							},
	// 							Direction: to.Ptr(armbatch.AccessRuleDirectionOutbound),
	// 							EmailAddresses: []*string{
	// 							},
	// 							FullyQualifiedDomainNames: []*string{
	// 								to.Ptr("paasrp1.contoso.org"),
	// 								to.Ptr("paasrp2.contoso.org")},
	// 								NetworkSecurityPerimeters: []*armbatch.NetworkSecurityPerimeter{
	// 								},
	// 								PhoneNumbers: []*string{
	// 								},
	// 								Subscriptions: []*armbatch.AccessRulePropertiesSubscriptionsItem{
	// 								},
	// 							},
	// 					}},
	// 					AccessRulesVersion: to.Ptr[int32](1),
	// 					DiagnosticSettingsVersion: to.Ptr[int32](1),
	// 					EnabledLogCategories: []*string{
	// 						to.Ptr("logCategory1"),
	// 						to.Ptr("logCategory2")},
	// 					},
	// 					ProvisioningIssues: []*armbatch.ProvisioningIssue{
	// 					},
	// 					ProvisioningState: to.Ptr(armbatch.NetworkSecurityPerimeterConfigurationProvisioningStateSucceeded),
	// 					ResourceAssociation: &armbatch.ResourceAssociation{
	// 						Name: to.Ptr("sampleassociation"),
	// 						AccessMode: to.Ptr(armbatch.ResourceAssociationAccessModeLearning),
	// 					},
	// 				},
	// 			}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e79d9ef3e065f2dcb6bd1db51e29c62a99dff5cb/specification/batch/resource-manager/Microsoft.Batch/stable/2024-07-01/examples/NspConfigurationReconcile.json
func ExampleNetworkSecurityPerimeterClient_BeginReconcileConfiguration() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armbatch.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewNetworkSecurityPerimeterClient().BeginReconcileConfiguration(ctx, "default-azurebatch-japaneast", "sampleacct", "00000000-0000-0000-0000-000000000000.sampleassociation", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}