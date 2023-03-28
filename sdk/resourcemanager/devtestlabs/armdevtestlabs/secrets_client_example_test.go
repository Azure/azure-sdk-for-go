//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armdevtestlabs_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/devtestlabs/armdevtestlabs/v2"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d55b8005f05b040b852c15e74a0f3e36494a15e1/specification/devtestlabs/resource-manager/Microsoft.DevTestLab/stable/2018-09-15/examples/Secrets_List.json
func ExampleSecretsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdevtestlabs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewSecretsClient().NewListPager("resourceGroupName", "{labName}", "{userName}", &armdevtestlabs.SecretsClientListOptions{Expand: nil,
		Filter:  nil,
		Top:     nil,
		Orderby: nil,
	})
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
		// page.SecretList = armdevtestlabs.SecretList{
		// 	Value: []*armdevtestlabs.Secret{
		// 		{
		// 			Name: to.Ptr("secret1"),
		// 			Type: to.Ptr("Microsoft.DevTestLab/labs/users/secrets"),
		// 			ID: to.Ptr("/subscriptions/{subscriptionId}/resourcegroups/resourceGroupName/providers/microsoft.devtestlab/labs/{labName}/users/{userName}/secrets/secret1"),
		// 			Properties: &armdevtestlabs.SecretProperties{
		// 				UniqueIdentifier: to.Ptr("00000000-0000-0000-0000-000000000000"),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("secret2"),
		// 			Type: to.Ptr("Microsoft.DevTestLab/labs/users/secrets"),
		// 			ID: to.Ptr("/subscriptions/{subscriptionId}/resourcegroups/resourceGroupName/providers/microsoft.devtestlab/labs/{labName}/users/{userName}/secrets/secret2"),
		// 			Properties: &armdevtestlabs.SecretProperties{
		// 				UniqueIdentifier: to.Ptr("00000000-0000-0000-0000-000000000000"),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d55b8005f05b040b852c15e74a0f3e36494a15e1/specification/devtestlabs/resource-manager/Microsoft.DevTestLab/stable/2018-09-15/examples/Secrets_Get.json
func ExampleSecretsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdevtestlabs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSecretsClient().Get(ctx, "resourceGroupName", "{labName}", "{userName}", "{secretName}", &armdevtestlabs.SecretsClientGetOptions{Expand: nil})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Secret = armdevtestlabs.Secret{
	// 	Name: to.Ptr("{secretName}"),
	// 	Type: to.Ptr("Microsoft.DevTestLab/labs/users/secrets"),
	// 	ID: to.Ptr("/subscriptions/{subscriptionId}/resourcegroups/resourceGroupName/providers/microsoft.devtestlab/labs/{labName}/users/{userName}/secrets/{secretName}"),
	// 	Properties: &armdevtestlabs.SecretProperties{
	// 		UniqueIdentifier: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d55b8005f05b040b852c15e74a0f3e36494a15e1/specification/devtestlabs/resource-manager/Microsoft.DevTestLab/stable/2018-09-15/examples/Secrets_CreateOrUpdate.json
func ExampleSecretsClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdevtestlabs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewSecretsClient().BeginCreateOrUpdate(ctx, "resourceGroupName", "{labName}", "{userName}", "{secretName}", armdevtestlabs.Secret{
		Properties: &armdevtestlabs.SecretProperties{
			Value: to.Ptr("{secret}"),
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
	// res.Secret = armdevtestlabs.Secret{
	// 	Name: to.Ptr("{secretName}"),
	// 	Type: to.Ptr("Microsoft.DevTestLab/labs/users/secrets"),
	// 	ID: to.Ptr("/subscriptions/{subscriptionId}/resourcegroups/resourceGroupName/providers/microsoft.devtestlab/labs/{labName}/users/{userName}/secrets/{secretName}"),
	// 	Properties: &armdevtestlabs.SecretProperties{
	// 		UniqueIdentifier: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d55b8005f05b040b852c15e74a0f3e36494a15e1/specification/devtestlabs/resource-manager/Microsoft.DevTestLab/stable/2018-09-15/examples/Secrets_Delete.json
func ExampleSecretsClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdevtestlabs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewSecretsClient().Delete(ctx, "resourceGroupName", "{labName}", "{userName}", "{secretName}", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d55b8005f05b040b852c15e74a0f3e36494a15e1/specification/devtestlabs/resource-manager/Microsoft.DevTestLab/stable/2018-09-15/examples/Secrets_Update.json
func ExampleSecretsClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdevtestlabs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSecretsClient().Update(ctx, "resourceGroupName", "{labName}", "{userName}", "{secretName}", armdevtestlabs.SecretFragment{
		Tags: map[string]*string{
			"tagName1": to.Ptr("tagValue1"),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Secret = armdevtestlabs.Secret{
	// 	Name: to.Ptr("{secretName}"),
	// 	Type: to.Ptr("Microsoft.DevTestLab/labs/users/secrets"),
	// 	ID: to.Ptr("/subscriptions/{subscriptionId}/resourcegroups/resourceGroupName/providers/microsoft.devtestlab/labs/{labName}/users/{userName}/secrets/{secretName}"),
	// 	Properties: &armdevtestlabs.SecretProperties{
	// 		UniqueIdentifier: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 	},
	// }
}
