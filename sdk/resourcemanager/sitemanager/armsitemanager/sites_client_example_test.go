// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package armsitemanager_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sitemanager/armsitemanager"
	"log"
)

// Generated from example definition: 2025-03-01-preview/Sites_CreateOrUpdate_MaximumSet_Gen.json
func ExampleSitesClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armsitemanager.NewClientFactory("0154f7fe-df09-4981-bf82-7ad5c1f596eb", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewSitesClient().BeginCreateOrUpdate(ctx, "rgsites", "string", armsitemanager.Site{
		Properties: &armsitemanager.SiteProperties{
			DisplayName: to.Ptr("string"),
			Labels: map[string]*string{
				"key8188": to.Ptr("mcgnu"),
			},
			Description: to.Ptr("enxcmpvfvadbapo"),
			SiteAddress: &armsitemanager.SiteAddressProperties{
				StreetAddress1:  to.Ptr("fodimymrxbhrfslsmzfhmitn"),
				StreetAddress2:  to.Ptr("widjg"),
				City:            to.Ptr("zkcbzjkatafo"),
				StateOrProvince: to.Ptr("wk"),
				Country:         to.Ptr("xeevcfvimlfzsfuxtyujw"),
				PostalCode:      to.Ptr("qbrhqk"),
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
	// res = armsitemanager.SitesClientCreateOrUpdateResponse{
	// 	Site: &armsitemanager.Site{
	// 		Properties: &armsitemanager.SiteProperties{
	// 			DisplayName: to.Ptr("string"),
	// 			Labels: map[string]*string{
	// 				"key8188": to.Ptr("mcgnu"),
	// 			},
	// 			Description: to.Ptr("enxcmpvfvadbapo"),
	// 			SiteAddress: &armsitemanager.SiteAddressProperties{
	// 				StreetAddress1: to.Ptr("fodimymrxbhrfslsmzfhmitn"),
	// 				StreetAddress2: to.Ptr("widjg"),
	// 				City: to.Ptr("zkcbzjkatafo"),
	// 				StateOrProvince: to.Ptr("wk"),
	// 				Country: to.Ptr("xeevcfvimlfzsfuxtyujw"),
	// 				PostalCode: to.Ptr("qbrhqk"),
	// 			},
	// 			ProvisioningState: to.Ptr(armsitemanager.ResourceProvisioningStateSucceeded),
	// 		},
	// 		ID: to.Ptr("/providers/Microsoft.Management/serviceGroups/SGSites/providers/Microsoft.Edge/Sites/Rome"),
	// 		Name: to.Ptr("string"),
	// 		Type: to.Ptr("string"),
	// 		SystemData: &armsitemanager.SystemData{
	// 			CreatedBy: to.Ptr("julxbiyjzi"),
	// 			CreatedByType: to.Ptr(armsitemanager.CreatedByTypeUser),
	// 			CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-12-30T07:53:03.972Z"); return t}()),
	// 			LastModifiedBy: to.Ptr("bceneuzzvzqmiocbrfef"),
	// 			LastModifiedByType: to.Ptr(armsitemanager.CreatedByTypeUser),
	// 			LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-12-30T07:53:03.972Z"); return t}()),
	// 		},
	// 	},
	// }
}

// Generated from example definition: 2025-03-01-preview/Sites_Delete_MaximumSet_Gen.json
func ExampleSitesClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armsitemanager.NewClientFactory("0154f7fe-df09-4981-bf82-7ad5c1f596eb", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSitesClient().Delete(ctx, "rgsites", "string", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armsitemanager.SitesClientDeleteResponse{
	// }
}

// Generated from example definition: 2025-03-01-preview/Sites_Get_MaximumSet_Gen.json
func ExampleSitesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armsitemanager.NewClientFactory("0154f7fe-df09-4981-bf82-7ad5c1f596eb", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSitesClient().Get(ctx, "rgsites", "string", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armsitemanager.SitesClientGetResponse{
	// 	Site: &armsitemanager.Site{
	// 		Properties: &armsitemanager.SiteProperties{
	// 			DisplayName: to.Ptr("string"),
	// 			Labels: map[string]*string{
	// 				"key8188": to.Ptr("mcgnu"),
	// 			},
	// 			Description: to.Ptr("enxcmpvfvadbapo"),
	// 			SiteAddress: &armsitemanager.SiteAddressProperties{
	// 				StreetAddress1: to.Ptr("fodimymrxbhrfslsmzfhmitn"),
	// 				StreetAddress2: to.Ptr("widjg"),
	// 				City: to.Ptr("zkcbzjkatafo"),
	// 				StateOrProvince: to.Ptr("wk"),
	// 				Country: to.Ptr("xeevcfvimlfzsfuxtyujw"),
	// 				PostalCode: to.Ptr("qbrhqk"),
	// 			},
	// 			ProvisioningState: to.Ptr(armsitemanager.ResourceProvisioningStateSucceeded),
	// 		},
	// 		ID: to.Ptr("/providers/Microsoft.Management/serviceGroups/SGSites/providers/Microsoft.Edge/Sites/Rome"),
	// 		Name: to.Ptr("string"),
	// 		Type: to.Ptr("string"),
	// 		SystemData: &armsitemanager.SystemData{
	// 			CreatedBy: to.Ptr("julxbiyjzi"),
	// 			CreatedByType: to.Ptr(armsitemanager.CreatedByTypeUser),
	// 			CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-12-30T07:53:03.972Z"); return t}()),
	// 			LastModifiedBy: to.Ptr("bceneuzzvzqmiocbrfef"),
	// 			LastModifiedByType: to.Ptr(armsitemanager.CreatedByTypeUser),
	// 			LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-12-30T07:53:03.972Z"); return t}()),
	// 		},
	// 	},
	// }
}

// Generated from example definition: 2025-03-01-preview/Sites_ListByResourceGroup_MaximumSet_Gen.json
func ExampleSitesClient_NewListByResourceGroupPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armsitemanager.NewClientFactory("0154f7fe-df09-4981-bf82-7ad5c1f596eb", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewSitesClient().NewListByResourceGroupPager("rgsites", nil)
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
		// page = armsitemanager.SitesClientListByResourceGroupResponse{
		// 	SiteListResult: armsitemanager.SiteListResult{
		// 		Value: []*armsitemanager.Site{
		// 			{
		// 				Properties: &armsitemanager.SiteProperties{
		// 					DisplayName: to.Ptr("string"),
		// 					ProvisioningState: to.Ptr(armsitemanager.ResourceProvisioningStateSucceeded),
		// 					Description: to.Ptr("mazbpkzbkvvntk"),
		// 					SiteAddress: &armsitemanager.SiteAddressProperties{
		// 						StreetAddress1: to.Ptr("fodimymrxbhrfslsmzfhmitn"),
		// 						StreetAddress2: to.Ptr("widjg"),
		// 						City: to.Ptr("zkcbzjkatafo"),
		// 						StateOrProvince: to.Ptr("wk"),
		// 						Country: to.Ptr("xeevcfvimlfzsfuxtyujw"),
		// 						PostalCode: to.Ptr("qbrhqk"),
		// 					},
		// 					Labels: map[string]*string{
		// 						"key8188": to.Ptr("mcgnu"),
		// 					},
		// 				},
		// 				ID: to.Ptr("/providers/Microsoft.Management/serviceGroups/SGSites/providers/Microsoft.Edge/Sites/Rome"),
		// 				Name: to.Ptr("string"),
		// 				Type: to.Ptr("string"),
		// 				SystemData: &armsitemanager.SystemData{
		// 					CreatedBy: to.Ptr("julxbiyjzi"),
		// 					CreatedByType: to.Ptr(armsitemanager.CreatedByTypeUser),
		// 					CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-12-30T07:53:03.972Z"); return t}()),
		// 					LastModifiedBy: to.Ptr("bceneuzzvzqmiocbrfef"),
		// 					LastModifiedByType: to.Ptr(armsitemanager.CreatedByTypeUser),
		// 					LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-12-30T07:53:03.972Z"); return t}()),
		// 				},
		// 			},
		// 		},
		// 		NextLink: to.Ptr("https://microsoft.com/a"),
		// 	},
		// }
	}
}

// Generated from example definition: 2025-03-01-preview/Sites_Update_MaximumSet_Gen.json
func ExampleSitesClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armsitemanager.NewClientFactory("0154f7fe-df09-4981-bf82-7ad5c1f596eb", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSitesClient().Update(ctx, "rgsites", "string", armsitemanager.SiteUpdate{
		Properties: &armsitemanager.SiteUpdateProperties{
			DisplayName: to.Ptr("string"),
			Description: to.Ptr("zztr"),
			SiteAddress: &armsitemanager.SiteAddressProperties{
				StreetAddress1:  to.Ptr("fodimymrxbhrfslsmzfhmitn"),
				StreetAddress2:  to.Ptr("widjg"),
				City:            to.Ptr("zkcbzjkatafo"),
				StateOrProvince: to.Ptr("wk"),
				Country:         to.Ptr("xeevcfvimlfzsfuxtyujw"),
				PostalCode:      to.Ptr("qbrhqk"),
			},
			Labels: map[string]*string{
				"key9939": to.Ptr("jdlzxcvcfqmruq"),
			},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armsitemanager.SitesClientUpdateResponse{
	// 	Site: &armsitemanager.Site{
	// 		Properties: &armsitemanager.SiteProperties{
	// 			DisplayName: to.Ptr("string"),
	// 			Labels: map[string]*string{
	// 				"key8188": to.Ptr("mcgnu"),
	// 			},
	// 			Description: to.Ptr("enxcmpvfvadbapo"),
	// 			SiteAddress: &armsitemanager.SiteAddressProperties{
	// 				StreetAddress1: to.Ptr("fodimymrxbhrfslsmzfhmitn"),
	// 				StreetAddress2: to.Ptr("widjg"),
	// 				City: to.Ptr("zkcbzjkatafo"),
	// 				StateOrProvince: to.Ptr("wk"),
	// 				Country: to.Ptr("xeevcfvimlfzsfuxtyujw"),
	// 				PostalCode: to.Ptr("qbrhqk"),
	// 			},
	// 			ProvisioningState: to.Ptr(armsitemanager.ResourceProvisioningStateSucceeded),
	// 		},
	// 		ID: to.Ptr("/providers/Microsoft.Management/serviceGroups/SGSites/providers/Microsoft.Edge/Sites/Rome"),
	// 		Name: to.Ptr("string"),
	// 		Type: to.Ptr("string"),
	// 		SystemData: &armsitemanager.SystemData{
	// 			CreatedBy: to.Ptr("julxbiyjzi"),
	// 			CreatedByType: to.Ptr(armsitemanager.CreatedByTypeUser),
	// 			CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-12-30T07:53:03.972Z"); return t}()),
	// 			LastModifiedBy: to.Ptr("bceneuzzvzqmiocbrfef"),
	// 			LastModifiedByType: to.Ptr(armsitemanager.CreatedByTypeUser),
	// 			LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-12-30T07:53:03.972Z"); return t}()),
	// 		},
	// 	},
	// }
}
