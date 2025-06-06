// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package armportalservicescopilot_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/portalservicescopilot/armportalservicescopilot"
	"log"
)

// Generated from example definition: 2024-04-01-preview/CopilotSettings_CreateOrUpdate.json
func ExampleCopilotSettingsClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armportalservicescopilot.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewCopilotSettingsClient().CreateOrUpdate(ctx, armportalservicescopilot.CopilotSettingsResource{
		Properties: &armportalservicescopilot.CopilotSettingsProperties{
			AccessControlEnabled: to.Ptr(true),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armportalservicescopilot.CopilotSettingsClientCreateOrUpdateResponse{
	// 	CopilotSettingsResource: &armportalservicescopilot.CopilotSettingsResource{
	// 		Properties: &armportalservicescopilot.CopilotSettingsProperties{
	// 			AccessControlEnabled: to.Ptr(true),
	// 		},
	// 		ID: to.Ptr("/providers/Microsoft.Portal/copilotSettings/default"),
	// 		Name: to.Ptr("default"),
	// 		Type: to.Ptr("Microsoft.PortalServices/copilotSettings"),
	// 	},
	// }
}

// Generated from example definition: 2024-04-01-preview/CopilotSettings_Delete.json
func ExampleCopilotSettingsClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armportalservicescopilot.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewCopilotSettingsClient().Delete(ctx, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armportalservicescopilot.CopilotSettingsClientDeleteResponse{
	// }
}

// Generated from example definition: 2024-04-01-preview/CopilotSettings_Get.json
func ExampleCopilotSettingsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armportalservicescopilot.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewCopilotSettingsClient().Get(ctx, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armportalservicescopilot.CopilotSettingsClientGetResponse{
	// 	CopilotSettingsResource: &armportalservicescopilot.CopilotSettingsResource{
	// 		Properties: &armportalservicescopilot.CopilotSettingsProperties{
	// 			AccessControlEnabled: to.Ptr(true),
	// 		},
	// 		ID: to.Ptr("/providers/Microsoft.Portal/copilotSettings/default"),
	// 		Name: to.Ptr("default"),
	// 		Type: to.Ptr("Microsoft.PortalServices/copilotSettings"),
	// 	},
	// }
}

// Generated from example definition: 2024-04-01-preview/CopilotSettings_Update.json
func ExampleCopilotSettingsClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armportalservicescopilot.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewCopilotSettingsClient().Update(ctx, armportalservicescopilot.CopilotSettingsResourceUpdate{
		Properties: &armportalservicescopilot.CopilotSettingsResourceUpdateProperties{
			AccessControlEnabled: to.Ptr(true),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armportalservicescopilot.CopilotSettingsClientUpdateResponse{
	// 	CopilotSettingsResource: &armportalservicescopilot.CopilotSettingsResource{
	// 		Properties: &armportalservicescopilot.CopilotSettingsProperties{
	// 			AccessControlEnabled: to.Ptr(true),
	// 		},
	// 		ID: to.Ptr("/providers/Microsoft.Portal/copilotSettings/default"),
	// 		Name: to.Ptr("default"),
	// 		Type: to.Ptr("Microsoft.PortalServices/copilotSettings"),
	// 	},
	// }
}
