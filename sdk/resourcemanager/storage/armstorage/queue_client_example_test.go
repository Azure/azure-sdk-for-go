//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armstorage_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/97ee23a6db6078abcbec7b75bf9af8c503e9bb8b/specification/storage/resource-manager/Microsoft.Storage/stable/2024-01-01/examples/QueueOperationPut.json
func ExampleQueueClient_Create_queueOperationPut() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armstorage.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewQueueClient().Create(ctx, "res3376", "sto328", "queue6185", armstorage.Queue{}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Queue = armstorage.Queue{
	// 	Name: to.Ptr("queue6185"),
	// 	Type: to.Ptr("Microsoft.Storage/storageAccounts/queueServices/queues"),
	// 	ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/res3376/providers/Microsoft.Storage/storageAccounts/sto328/queueServices/default/queues/queue6185"),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/97ee23a6db6078abcbec7b75bf9af8c503e9bb8b/specification/storage/resource-manager/Microsoft.Storage/stable/2024-01-01/examples/QueueOperationPutWithMetadata.json
func ExampleQueueClient_Create_queueOperationPutWithMetadata() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armstorage.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewQueueClient().Create(ctx, "res3376", "sto328", "queue6185", armstorage.Queue{
		QueueProperties: &armstorage.QueueProperties{
			Metadata: map[string]*string{
				"sample1": to.Ptr("meta1"),
				"sample2": to.Ptr("meta2"),
			},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Queue = armstorage.Queue{
	// 	Name: to.Ptr("queue6185"),
	// 	Type: to.Ptr("Microsoft.Storage/storageAccounts/queueServices/queues"),
	// 	ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/res3376/providers/Microsoft.Storage/storageAccounts/sto328/queueServices/default/queues/queue6185"),
	// 	QueueProperties: &armstorage.QueueProperties{
	// 		Metadata: map[string]*string{
	// 			"sample1": to.Ptr("meta1"),
	// 			"sample2": to.Ptr("meta2"),
	// 		},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/97ee23a6db6078abcbec7b75bf9af8c503e9bb8b/specification/storage/resource-manager/Microsoft.Storage/stable/2024-01-01/examples/QueueOperationPatch.json
func ExampleQueueClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armstorage.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewQueueClient().Update(ctx, "res3376", "sto328", "queue6185", armstorage.Queue{}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Queue = armstorage.Queue{
	// 	Name: to.Ptr("queue6185"),
	// 	Type: to.Ptr("Microsoft.Storage/storageAccounts/queueServices/queues"),
	// 	ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/res3376/providers/Microsoft.Storage/storageAccounts/sto328/queueServices/default/queues/queue6185"),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/97ee23a6db6078abcbec7b75bf9af8c503e9bb8b/specification/storage/resource-manager/Microsoft.Storage/stable/2024-01-01/examples/QueueOperationGet.json
func ExampleQueueClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armstorage.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewQueueClient().Get(ctx, "res3376", "sto328", "queue6185", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Queue = armstorage.Queue{
	// 	Name: to.Ptr("queue6185"),
	// 	Type: to.Ptr("Microsoft.Storage/storageAccounts/queueServices/queues"),
	// 	ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/res3376/providers/Microsoft.Storage/storageAccounts/sto328/queueServices/default/queues/queue6185"),
	// 	QueueProperties: &armstorage.QueueProperties{
	// 		Metadata: map[string]*string{
	// 			"sample1": to.Ptr("meta1"),
	// 			"sample2": to.Ptr("meta2"),
	// 		},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/97ee23a6db6078abcbec7b75bf9af8c503e9bb8b/specification/storage/resource-manager/Microsoft.Storage/stable/2024-01-01/examples/QueueOperationDelete.json
func ExampleQueueClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armstorage.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewQueueClient().Delete(ctx, "res3376", "sto328", "queue6185", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/97ee23a6db6078abcbec7b75bf9af8c503e9bb8b/specification/storage/resource-manager/Microsoft.Storage/stable/2024-01-01/examples/QueueOperationList.json
func ExampleQueueClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armstorage.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewQueueClient().NewListPager("res9290", "sto328", &armstorage.QueueClientListOptions{Maxpagesize: nil,
		Filter: nil,
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
		// page.ListQueueResource = armstorage.ListQueueResource{
		// 	Value: []*armstorage.ListQueue{
		// 		{
		// 			Name: to.Ptr("queue6185"),
		// 			Type: to.Ptr("Microsoft.Storage/storageAccounts/queueServices/queues"),
		// 			ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/res3376/providers/Microsoft.Storage/storageAccounts/sto328/queueServices/default/queues/queue6185"),
		// 			QueueProperties: &armstorage.ListQueueProperties{
		// 				Metadata: map[string]*string{
		// 					"sample1": to.Ptr("meta1"),
		// 					"sample2": to.Ptr("meta2"),
		// 				},
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("queue6186"),
		// 			Type: to.Ptr("Microsoft.Storage/storageAccounts/queueServices/queues"),
		// 			ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/res3376/providers/Microsoft.Storage/storageAccounts/sto328/queueServices/default/queues/queue6186"),
		// 			QueueProperties: &armstorage.ListQueueProperties{
		// 				Metadata: map[string]*string{
		// 					"sample1": to.Ptr("meta1"),
		// 					"sample2": to.Ptr("meta2"),
		// 				},
		// 			},
		// 	}},
		// }
	}
}
