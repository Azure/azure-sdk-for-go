// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package armbicep_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armbicep"
	"log"
)

// Generated from example definition: 2023-11-01/DecompileBicep.json
func ExampleDecompileOperationGroupClient_Bicep() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armbicep.NewClientFactory("00000000-0000-0000-0000-000000000000", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewDecompileOperationGroupClient().Bicep(ctx, armbicep.DecompileOperationRequest{
		Template: to.Ptr("{\r\n \"$schema\": \"https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#\",\r\n \"contentVersion\": \"1.0.0.0\",\r\n \"metadata\": {\r\n \"_generator\": {\r\n \"name\": \"bicep\",\r\n \"version\": \"0.15.31.15270\",\r\n \"templateHash\": \"9249505596133208719\"\r\n }\r\n },\r\n \"parameters\": {\r\n \"storageAccountName\": {\r\n \"type\": \"string\"\r\n }\r\n },\r\n \"resources\": []\r\n}"),
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armbicep.DecompileOperationGroupClientBicepResponse{
	// 	DecompileOperationSuccessResponse: &armbicep.DecompileOperationSuccessResponse{
	// 		EntryPoint: to.Ptr("main.bicep"),
	// 		Files: []*armbicep.FileDefinition{
	// 			{
	// 				Path: to.Ptr("main.bicep"),
	// 				Contents: to.Ptr("param storageAccountName string"),
	// 			},
	// 		},
	// 	},
	// }
}
