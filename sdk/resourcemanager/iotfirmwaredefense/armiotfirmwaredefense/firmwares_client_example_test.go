// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package armiotfirmwaredefense_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/iotfirmwaredefense/armiotfirmwaredefense/v2"
	"log"
)

// Generated from example definition: 2025-04-01-preview/Firmwares_Create_MaximumSet_Gen.json
func ExampleFirmwaresClient_Create_firmwaresCreateMaximumSetGenGeneratedByMaximumSetRule() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armiotfirmwaredefense.NewClientFactory("5C707B5F-6130-4F71-819E-953A28942E88", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewFirmwaresClient().Create(ctx, "rgiotfirmwaredefense", "exampleWorkspaceName", "00000000-0000-0000-0000-000000000000", armiotfirmwaredefense.Firmware{
		Properties: &armiotfirmwaredefense.FirmwareProperties{
			FileName: to.Ptr("dmnqhyxssutvnewntlb"),
			Vendor:   to.Ptr("hymojocxpxqhtblioaavylnzyg"),
			Model:    to.Ptr("wmyfbyjsggbvxcuin"),
			Version:  to.Ptr("nhtxzslgcbtptu"),
			FileSize: to.Ptr[int64](30),
			Status:   to.Ptr(armiotfirmwaredefense.StatusPending),
			StatusMessages: []*armiotfirmwaredefense.StatusMessage{
				{
					ErrorCode: to.Ptr[int64](20),
					Message:   to.Ptr("edtylkjvj"),
				},
			},
			Description: to.Ptr("sqt"),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armiotfirmwaredefense.FirmwaresClientCreateResponse{
	// 	Firmware: &armiotfirmwaredefense.Firmware{
	// 		Properties: &armiotfirmwaredefense.FirmwareProperties{
	// 			FileName: to.Ptr("dmnqhyxssutvnewntlb"),
	// 			Vendor: to.Ptr("hymojocxpxqhtblioaavylnzyg"),
	// 			Model: to.Ptr("wmyfbyjsggbvxcuin"),
	// 			Version: to.Ptr("nhtxzslgcbtptu"),
	// 			FileSize: to.Ptr[int64](30),
	// 			Status: to.Ptr(armiotfirmwaredefense.StatusPending),
	// 			StatusMessages: []*armiotfirmwaredefense.StatusMessage{
	// 				{
	// 					ErrorCode: to.Ptr[int64](20),
	// 					Message: to.Ptr("edtylkjvj"),
	// 				},
	// 			},
	// 			Description: to.Ptr("sqt"),
	// 			ProvisioningState: to.Ptr(armiotfirmwaredefense.ProvisioningStateSucceeded),
	// 		},
	// 		ID: to.Ptr("/subscriptions/07aed47b-60ad-4d6e-a07a-000000000000/resourceGroups/FirmwareAnalysisRG/providers/Microsoft.IoTFirmwareDefense/workspaces/default/firmwares/109a9886-50bf-85a8-9d75-000000000000/summaries/firmware"),
	// 		Name: to.Ptr("qobb"),
	// 		Type: to.Ptr("xf"),
	// 		SystemData: &armiotfirmwaredefense.SystemData{
	// 			CreatedBy: to.Ptr("nqisshvdzqcxzbujvacin"),
	// 			CreatedByType: to.Ptr(armiotfirmwaredefense.CreatedByTypeUser),
	// 			CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-06-13T15:22:45.940Z"); return t}()),
	// 			LastModifiedBy: to.Ptr("of"),
	// 			LastModifiedByType: to.Ptr(armiotfirmwaredefense.CreatedByTypeUser),
	// 			LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-06-13T15:22:45.940Z"); return t}()),
	// 		},
	// 	},
	// }
}

// Generated from example definition: 2025-04-01-preview/Firmwares_Create_MinimumSet_Gen.json
func ExampleFirmwaresClient_Create_firmwaresCreateMinimumSetGen() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armiotfirmwaredefense.NewClientFactory("685C0C6F-9867-4B1C-A534-AA3A05B54BCE", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewFirmwaresClient().Create(ctx, "rgworkspaces-firmwares", "A7", "umrkdttp", armiotfirmwaredefense.Firmware{
		Properties: &armiotfirmwaredefense.FirmwareProperties{},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armiotfirmwaredefense.FirmwaresClientCreateResponse{
	// 	Firmware: &armiotfirmwaredefense.Firmware{
	// 	},
	// }
}

// Generated from example definition: 2025-04-01-preview/Firmwares_Delete_MaximumSet_Gen.json
func ExampleFirmwaresClient_Delete_firmwaresDeleteMaximumSetGenGeneratedByMaximumSetRule() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armiotfirmwaredefense.NewClientFactory("5C707B5F-6130-4F71-819E-953A28942E88", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewFirmwaresClient().Delete(ctx, "rgiotfirmwaredefense", "exampleWorkspaceName", "00000000-0000-0000-0000-000000000000", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armiotfirmwaredefense.FirmwaresClientDeleteResponse{
	// }
}

// Generated from example definition: 2025-04-01-preview/Firmwares_Delete_MinimumSet_Gen.json
func ExampleFirmwaresClient_Delete_firmwaresDeleteMinimumSetGen() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armiotfirmwaredefense.NewClientFactory("685C0C6F-9867-4B1C-A534-AA3A05B54BCE", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewFirmwaresClient().Delete(ctx, "rgworkspaces-firmwares", "A7", "umrkdttp", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armiotfirmwaredefense.FirmwaresClientDeleteResponse{
	// }
}

// Generated from example definition: 2025-04-01-preview/Firmwares_Get_MaximumSet_Gen.json
func ExampleFirmwaresClient_Get_firmwaresGetMaximumSetGenGeneratedByMaximumSetRule() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armiotfirmwaredefense.NewClientFactory("5C707B5F-6130-4F71-819E-953A28942E88", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewFirmwaresClient().Get(ctx, "rgiotfirmwaredefense", "exampleWorkspaceName", "00000000-0000-0000-0000-000000000000", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armiotfirmwaredefense.FirmwaresClientGetResponse{
	// 	Firmware: &armiotfirmwaredefense.Firmware{
	// 		Properties: &armiotfirmwaredefense.FirmwareProperties{
	// 			FileName: to.Ptr("dmnqhyxssutvnewntlb"),
	// 			Vendor: to.Ptr("hymojocxpxqhtblioaavylnzyg"),
	// 			Model: to.Ptr("wmyfbyjsggbvxcuin"),
	// 			Version: to.Ptr("nhtxzslgcbtptu"),
	// 			FileSize: to.Ptr[int64](30),
	// 			Status: to.Ptr(armiotfirmwaredefense.StatusPending),
	// 			StatusMessages: []*armiotfirmwaredefense.StatusMessage{
	// 				{
	// 					ErrorCode: to.Ptr[int64](20),
	// 					Message: to.Ptr("edtylkjvj"),
	// 				},
	// 			},
	// 			Description: to.Ptr("sqt"),
	// 			ProvisioningState: to.Ptr(armiotfirmwaredefense.ProvisioningStateSucceeded),
	// 		},
	// 		ID: to.Ptr("/subscriptions/07aed47b-60ad-4d6e-a07a-000000000000/resourceGroups/FirmwareAnalysisRG/providers/Microsoft.IoTFirmwareDefense/workspaces/default/firmwares/109a9886-50bf-85a8-9d75-000000000000/summaries/firmware"),
	// 		Name: to.Ptr("qobb"),
	// 		Type: to.Ptr("xf"),
	// 		SystemData: &armiotfirmwaredefense.SystemData{
	// 			CreatedBy: to.Ptr("nqisshvdzqcxzbujvacin"),
	// 			CreatedByType: to.Ptr(armiotfirmwaredefense.CreatedByTypeUser),
	// 			CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-06-13T15:22:45.940Z"); return t}()),
	// 			LastModifiedBy: to.Ptr("of"),
	// 			LastModifiedByType: to.Ptr(armiotfirmwaredefense.CreatedByTypeUser),
	// 			LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-06-13T15:22:45.940Z"); return t}()),
	// 		},
	// 	},
	// }
}

// Generated from example definition: 2025-04-01-preview/Firmwares_Get_MinimumSet_Gen.json
func ExampleFirmwaresClient_Get_firmwaresGetMinimumSetGen() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armiotfirmwaredefense.NewClientFactory("685C0C6F-9867-4B1C-A534-AA3A05B54BCE", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewFirmwaresClient().Get(ctx, "rgworkspaces-firmwares", "A7", "umrkdttp", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armiotfirmwaredefense.FirmwaresClientGetResponse{
	// 	Firmware: &armiotfirmwaredefense.Firmware{
	// 	},
	// }
}

// Generated from example definition: 2025-04-01-preview/Firmwares_ListByWorkspace_MaximumSet_Gen.json
func ExampleFirmwaresClient_NewListByWorkspacePager_firmwaresListByWorkspaceMaximumSetGenGeneratedByMaximumSetRule() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armiotfirmwaredefense.NewClientFactory("5C707B5F-6130-4F71-819E-953A28942E88", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewFirmwaresClient().NewListByWorkspacePager("rgiotfirmwaredefense", "exampleWorkspaceName", nil)
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
		// page = armiotfirmwaredefense.FirmwaresClientListByWorkspaceResponse{
		// 	FirmwareListResult: armiotfirmwaredefense.FirmwareListResult{
		// 		Value: []*armiotfirmwaredefense.Firmware{
		// 			{
		// 				Properties: &armiotfirmwaredefense.FirmwareProperties{
		// 					FileName: to.Ptr("dmnqhyxssutvnewntlb"),
		// 					Vendor: to.Ptr("hymojocxpxqhtblioaavylnzyg"),
		// 					Model: to.Ptr("wmyfbyjsggbvxcuin"),
		// 					Version: to.Ptr("nhtxzslgcbtptu"),
		// 					FileSize: to.Ptr[int64](30),
		// 					Status: to.Ptr(armiotfirmwaredefense.StatusPending),
		// 					StatusMessages: []*armiotfirmwaredefense.StatusMessage{
		// 						{
		// 							ErrorCode: to.Ptr[int64](20),
		// 							Message: to.Ptr("edtylkjvj"),
		// 						},
		// 					},
		// 					ProvisioningState: to.Ptr(armiotfirmwaredefense.ProvisioningStateSucceeded),
		// 					Description: to.Ptr("zgzvtldpivwm"),
		// 				},
		// 				ID: to.Ptr("/subscriptions/07aed47b-60ad-4d6e-a07a-000000000000/resourceGroups/FirmwareAnalysisRG/providers/Microsoft.IoTFirmwareDefense/workspaces/default/firmwares/109a9886-50bf-85a8-9d75-000000000000/summaries/firmware"),
		// 				Name: to.Ptr("qobb"),
		// 				Type: to.Ptr("xf"),
		// 				SystemData: &armiotfirmwaredefense.SystemData{
		// 					CreatedBy: to.Ptr("nqisshvdzqcxzbujvacin"),
		// 					CreatedByType: to.Ptr(armiotfirmwaredefense.CreatedByTypeUser),
		// 					CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-06-13T15:22:45.940Z"); return t}()),
		// 					LastModifiedBy: to.Ptr("of"),
		// 					LastModifiedByType: to.Ptr(armiotfirmwaredefense.CreatedByTypeUser),
		// 					LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-06-13T15:22:45.940Z"); return t}()),
		// 				},
		// 			},
		// 		},
		// 		NextLink: to.Ptr("https://microsoft.com/a"),
		// 	},
		// }
	}
}

// Generated from example definition: 2025-04-01-preview/Firmwares_ListByWorkspace_MinimumSet_Gen.json
func ExampleFirmwaresClient_NewListByWorkspacePager_firmwaresListByWorkspaceMinimumSetGen() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armiotfirmwaredefense.NewClientFactory("685C0C6F-9867-4B1C-A534-AA3A05B54BCE", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewFirmwaresClient().NewListByWorkspacePager("rgworkspaces-firmwares", "A7", nil)
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
		// page = armiotfirmwaredefense.FirmwaresClientListByWorkspaceResponse{
		// 	FirmwareListResult: armiotfirmwaredefense.FirmwareListResult{
		// 		Value: []*armiotfirmwaredefense.Firmware{
		// 		},
		// 	},
		// }
	}
}

// Generated from example definition: 2025-04-01-preview/Firmwares_Update_MaximumSet_Gen.json
func ExampleFirmwaresClient_Update_firmwaresUpdateMaximumSetGenGeneratedByMaximumSetRule() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armiotfirmwaredefense.NewClientFactory("5C707B5F-6130-4F71-819E-953A28942E88", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewFirmwaresClient().Update(ctx, "rgiotfirmwaredefense", "exampleWorkspaceName", "00000000-0000-0000-0000-000000000000", armiotfirmwaredefense.FirmwareUpdateDefinition{
		Properties: &armiotfirmwaredefense.FirmwareProperties{
			FileName: to.Ptr("dmnqhyxssutvnewntlb"),
			Vendor:   to.Ptr("hymojocxpxqhtblioaavylnzyg"),
			Model:    to.Ptr("wmyfbyjsggbvxcuin"),
			Version:  to.Ptr("nhtxzslgcbtptu"),
			FileSize: to.Ptr[int64](30),
			Status:   to.Ptr(armiotfirmwaredefense.StatusPending),
			StatusMessages: []*armiotfirmwaredefense.StatusMessage{
				{
					ErrorCode: to.Ptr[int64](20),
					Message:   to.Ptr("edtylkjvj"),
				},
			},
			Description: to.Ptr("nknvqnkgumzbupxe"),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armiotfirmwaredefense.FirmwaresClientUpdateResponse{
	// 	Firmware: &armiotfirmwaredefense.Firmware{
	// 		Properties: &armiotfirmwaredefense.FirmwareProperties{
	// 			FileName: to.Ptr("dmnqhyxssutvnewntlb"),
	// 			Vendor: to.Ptr("hymojocxpxqhtblioaavylnzyg"),
	// 			Model: to.Ptr("wmyfbyjsggbvxcuin"),
	// 			Version: to.Ptr("nhtxzslgcbtptu"),
	// 			FileSize: to.Ptr[int64](30),
	// 			Status: to.Ptr(armiotfirmwaredefense.StatusPending),
	// 			StatusMessages: []*armiotfirmwaredefense.StatusMessage{
	// 				{
	// 					ErrorCode: to.Ptr[int64](20),
	// 					Message: to.Ptr("edtylkjvj"),
	// 				},
	// 			},
	// 			Description: to.Ptr("sqt"),
	// 			ProvisioningState: to.Ptr(armiotfirmwaredefense.ProvisioningStateSucceeded),
	// 		},
	// 		ID: to.Ptr("/subscriptions/07aed47b-60ad-4d6e-a07a-000000000000/resourceGroups/FirmwareAnalysisRG/providers/Microsoft.IoTFirmwareDefense/workspaces/default/firmwares/109a9886-50bf-85a8-9d75-000000000000/summaries/firmware"),
	// 		Name: to.Ptr("qobb"),
	// 		Type: to.Ptr("xf"),
	// 		SystemData: &armiotfirmwaredefense.SystemData{
	// 			CreatedBy: to.Ptr("nqisshvdzqcxzbujvacin"),
	// 			CreatedByType: to.Ptr(armiotfirmwaredefense.CreatedByTypeUser),
	// 			CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-06-13T15:22:45.940Z"); return t}()),
	// 			LastModifiedBy: to.Ptr("of"),
	// 			LastModifiedByType: to.Ptr(armiotfirmwaredefense.CreatedByTypeUser),
	// 			LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-06-13T15:22:45.940Z"); return t}()),
	// 		},
	// 	},
	// }
}

// Generated from example definition: 2025-04-01-preview/Firmwares_Update_MinimumSet_Gen.json
func ExampleFirmwaresClient_Update_firmwaresUpdateMinimumSetGen() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armiotfirmwaredefense.NewClientFactory("685C0C6F-9867-4B1C-A534-AA3A05B54BCE", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewFirmwaresClient().Update(ctx, "rgworkspaces-firmwares", "A7", "umrkdttp", armiotfirmwaredefense.FirmwareUpdateDefinition{}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armiotfirmwaredefense.FirmwaresClientUpdateResponse{
	// 	Firmware: &armiotfirmwaredefense.Firmware{
	// 	},
	// }
}
