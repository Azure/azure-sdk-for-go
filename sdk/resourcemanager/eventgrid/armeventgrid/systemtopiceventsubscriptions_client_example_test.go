//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armeventgrid_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/eventgrid/armeventgrid/v2"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/ee1eec42dcc710ff88db2d1bf574b2f9afe3d654/specification/eventgrid/resource-manager/Microsoft.EventGrid/stable/2025-02-15/examples/SystemTopicEventSubscriptions_Get.json
func ExampleSystemTopicEventSubscriptionsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armeventgrid.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSystemTopicEventSubscriptionsClient().Get(ctx, "examplerg", "exampleSystemTopic1", "examplesubscription1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.EventSubscription = armeventgrid.EventSubscription{
	// 	Name: to.Ptr("examplesubscription1"),
	// 	Type: to.Ptr("Microsoft.EventGrid/systemTopics/eventSubscriptions"),
	// 	ID: to.Ptr("/subscriptions/5b4b650e-28b9-4790-b3ab-ddbd88d727c4/resourceGroups/examplerg/providers/Microsoft.EventGrid/systemTopics/exampleSystemTopic1/eventSubscriptions/examplesubscription1"),
	// 	Properties: &armeventgrid.EventSubscriptionProperties{
	// 		Destination: &armeventgrid.StorageQueueEventSubscriptionDestination{
	// 			EndpointType: to.Ptr(armeventgrid.EndpointTypeStorageQueue),
	// 			Properties: &armeventgrid.StorageQueueEventSubscriptionDestinationProperties{
	// 				QueueName: to.Ptr("que"),
	// 				ResourceID: to.Ptr("/subscriptions/5b4b650e-28b9-4790-b3ab-ddbd88d727c4/resourceGroups/examplerg/providers/Microsoft.Storage/storageAccounts/testtrackedsource"),
	// 			},
	// 		},
	// 		EventDeliverySchema: to.Ptr(armeventgrid.EventDeliverySchemaEventGridSchema),
	// 		Filter: &armeventgrid.EventSubscriptionFilter{
	// 			IncludedEventTypes: []*string{
	// 				to.Ptr("Microsoft.Storage.BlobCreated"),
	// 				to.Ptr("Microsoft.Storage.BlobDeleted")},
	// 				IsSubjectCaseSensitive: to.Ptr(false),
	// 				SubjectBeginsWith: to.Ptr("ExamplePrefix"),
	// 				SubjectEndsWith: to.Ptr("ExampleSuffix"),
	// 			},
	// 			Labels: []*string{
	// 				to.Ptr("label1"),
	// 				to.Ptr("label2")},
	// 				ProvisioningState: to.Ptr(armeventgrid.EventSubscriptionProvisioningStateSucceeded),
	// 				RetryPolicy: &armeventgrid.RetryPolicy{
	// 					EventTimeToLiveInMinutes: to.Ptr[int32](1440),
	// 					MaxDeliveryAttempts: to.Ptr[int32](30),
	// 				},
	// 				Topic: to.Ptr("/subscriptions/5b4b650e-28b9-4790-b3ab-ddbd88d727c4/resourceGroups/examplerg/providers/Microsoft.EventGrid/systemTopics/exampleSystemTopic1"),
	// 			},
	// 		}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/ee1eec42dcc710ff88db2d1bf574b2f9afe3d654/specification/eventgrid/resource-manager/Microsoft.EventGrid/stable/2025-02-15/examples/SystemTopicEventSubscriptions_CreateOrUpdate.json
func ExampleSystemTopicEventSubscriptionsClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armeventgrid.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewSystemTopicEventSubscriptionsClient().BeginCreateOrUpdate(ctx, "examplerg", "exampleSystemTopic1", "exampleEventSubscriptionName1", armeventgrid.EventSubscription{
		Properties: &armeventgrid.EventSubscriptionProperties{
			Destination: &armeventgrid.WebHookEventSubscriptionDestination{
				EndpointType: to.Ptr(armeventgrid.EndpointTypeWebHook),
				Properties: &armeventgrid.WebHookEventSubscriptionDestinationProperties{
					EndpointURL: to.Ptr("https://requestb.in/15ksip71"),
				},
			},
			Filter: &armeventgrid.EventSubscriptionFilter{
				IsSubjectCaseSensitive: to.Ptr(false),
				SubjectBeginsWith:      to.Ptr("ExamplePrefix"),
				SubjectEndsWith:        to.Ptr("ExampleSuffix"),
			},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/ee1eec42dcc710ff88db2d1bf574b2f9afe3d654/specification/eventgrid/resource-manager/Microsoft.EventGrid/stable/2025-02-15/examples/SystemTopicEventSubscriptions_Delete.json
func ExampleSystemTopicEventSubscriptionsClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armeventgrid.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewSystemTopicEventSubscriptionsClient().BeginDelete(ctx, "examplerg", "exampleSystemTopic1", "examplesubscription1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/ee1eec42dcc710ff88db2d1bf574b2f9afe3d654/specification/eventgrid/resource-manager/Microsoft.EventGrid/stable/2025-02-15/examples/SystemTopicEventSubscriptions_Update.json
func ExampleSystemTopicEventSubscriptionsClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armeventgrid.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewSystemTopicEventSubscriptionsClient().BeginUpdate(ctx, "examplerg", "exampleSystemTopic1", "exampleEventSubscriptionName1", armeventgrid.EventSubscriptionUpdateParameters{
		Destination: &armeventgrid.WebHookEventSubscriptionDestination{
			EndpointType: to.Ptr(armeventgrid.EndpointTypeWebHook),
			Properties: &armeventgrid.WebHookEventSubscriptionDestinationProperties{
				EndpointURL: to.Ptr("https://requestb.in/15ksip71"),
			},
		},
		Filter: &armeventgrid.EventSubscriptionFilter{
			IsSubjectCaseSensitive: to.Ptr(true),
			SubjectBeginsWith:      to.Ptr("existingPrefix"),
			SubjectEndsWith:        to.Ptr("newSuffix"),
		},
		Labels: []*string{
			to.Ptr("label1"),
			to.Ptr("label2")},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/ee1eec42dcc710ff88db2d1bf574b2f9afe3d654/specification/eventgrid/resource-manager/Microsoft.EventGrid/stable/2025-02-15/examples/SystemTopicEventSubscriptions_GetFullUrl.json
func ExampleSystemTopicEventSubscriptionsClient_GetFullURL() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armeventgrid.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSystemTopicEventSubscriptionsClient().GetFullURL(ctx, "examplerg", "exampleSystemTopic1", "examplesubscription1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.EventSubscriptionFullURL = armeventgrid.EventSubscriptionFullURL{
	// 	EndpointURL: to.Ptr("https://requestb.in/15ksip71"),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/ee1eec42dcc710ff88db2d1bf574b2f9afe3d654/specification/eventgrid/resource-manager/Microsoft.EventGrid/stable/2025-02-15/examples/SystemTopicEventSubscriptions_ListBySystemTopic.json
func ExampleSystemTopicEventSubscriptionsClient_NewListBySystemTopicPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armeventgrid.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewSystemTopicEventSubscriptionsClient().NewListBySystemTopicPager("examplerg", "exampleSystemTopic1", &armeventgrid.SystemTopicEventSubscriptionsClientListBySystemTopicOptions{Filter: nil,
		Top: nil,
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
		// page.EventSubscriptionsListResult = armeventgrid.EventSubscriptionsListResult{
		// 	Value: []*armeventgrid.EventSubscription{
		// 		{
		// 			Name: to.Ptr("examplesubscription1"),
		// 			Type: to.Ptr("Microsoft.EventGrid/systemTopics/eventSubscriptions"),
		// 			ID: to.Ptr("/subscriptions/5b4b650e-28b9-4790-b3ab-ddbd88d727c4/resourceGroups/examplerg/providers/Microsoft.EventGrid/systemTopics/exampleSystemTopic1/eventSubscriptions/examplesubscription1"),
		// 			Properties: &armeventgrid.EventSubscriptionProperties{
		// 				Destination: &armeventgrid.StorageQueueEventSubscriptionDestination{
		// 					EndpointType: to.Ptr(armeventgrid.EndpointTypeStorageQueue),
		// 					Properties: &armeventgrid.StorageQueueEventSubscriptionDestinationProperties{
		// 						QueueName: to.Ptr("que"),
		// 						ResourceID: to.Ptr("/subscriptions/5b4b650e-28b9-4790-b3ab-ddbd88d727c4/resourceGroups/examplerg/providers/Microsoft.Storage/storageAccounts/testtrackedsource"),
		// 					},
		// 				},
		// 				EventDeliverySchema: to.Ptr(armeventgrid.EventDeliverySchemaEventGridSchema),
		// 				Filter: &armeventgrid.EventSubscriptionFilter{
		// 					SubjectBeginsWith: to.Ptr(""),
		// 					SubjectEndsWith: to.Ptr(""),
		// 				},
		// 				Labels: []*string{
		// 				},
		// 				ProvisioningState: to.Ptr(armeventgrid.EventSubscriptionProvisioningStateSucceeded),
		// 				RetryPolicy: &armeventgrid.RetryPolicy{
		// 					EventTimeToLiveInMinutes: to.Ptr[int32](1440),
		// 					MaxDeliveryAttempts: to.Ptr[int32](10),
		// 				},
		// 				Topic: to.Ptr("/subscriptions/5b4b650e-28b9-4790-b3ab-ddbd88d727c4/resourceGroups/examplerg/providers/Microsoft.EventGrid/systemTopics/exampleSystemTopic1"),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("examplesubscription2"),
		// 			Type: to.Ptr("Microsoft.EventGrid/systemTopics/eventSubscriptions"),
		// 			ID: to.Ptr("/subscriptions/5b4b650e-28b9-4790-b3ab-ddbd88d727c4/resourceGroups/examplerg/providers/Microsoft.EventGrid/systemTopics/exampleSystemTopic1/eventSubscriptions/examplesubscription2"),
		// 			Properties: &armeventgrid.EventSubscriptionProperties{
		// 				Destination: &armeventgrid.StorageQueueEventSubscriptionDestination{
		// 					EndpointType: to.Ptr(armeventgrid.EndpointTypeStorageQueue),
		// 					Properties: &armeventgrid.StorageQueueEventSubscriptionDestinationProperties{
		// 						QueueName: to.Ptr("que"),
		// 						ResourceID: to.Ptr("/subscriptions/5b4b650e-28b9-4790-b3ab-ddbd88d727c4/resourceGroups/examplerg/providers/Microsoft.Storage/storageAccounts/testtrackedsource"),
		// 					},
		// 				},
		// 				EventDeliverySchema: to.Ptr(armeventgrid.EventDeliverySchemaEventGridSchema),
		// 				Filter: &armeventgrid.EventSubscriptionFilter{
		// 					IncludedEventTypes: []*string{
		// 						to.Ptr("Microsoft.Storage.BlobCreated"),
		// 						to.Ptr("Microsoft.Storage.BlobDeleted")},
		// 						IsSubjectCaseSensitive: to.Ptr(false),
		// 						SubjectBeginsWith: to.Ptr("ExamplePrefix"),
		// 						SubjectEndsWith: to.Ptr("ExampleSuffix"),
		// 					},
		// 					Labels: []*string{
		// 						to.Ptr("label1"),
		// 						to.Ptr("label2")},
		// 						ProvisioningState: to.Ptr(armeventgrid.EventSubscriptionProvisioningStateSucceeded),
		// 						RetryPolicy: &armeventgrid.RetryPolicy{
		// 							EventTimeToLiveInMinutes: to.Ptr[int32](1440),
		// 							MaxDeliveryAttempts: to.Ptr[int32](30),
		// 						},
		// 						Topic: to.Ptr("/subscriptions/5b4b650e-28b9-4790-b3ab-ddbd88d727c4/resourceGroups/examplerg/providers/Microsoft.EventGrid/systemTopics/exampleSystemTopic1"),
		// 					},
		// 			}},
		// 		}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/ee1eec42dcc710ff88db2d1bf574b2f9afe3d654/specification/eventgrid/resource-manager/Microsoft.EventGrid/stable/2025-02-15/examples/SystemTopicEventSubscriptions_GetDeliveryAttributes.json
func ExampleSystemTopicEventSubscriptionsClient_GetDeliveryAttributes() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armeventgrid.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSystemTopicEventSubscriptionsClient().GetDeliveryAttributes(ctx, "examplerg", "exampleSystemTopic1", "examplesubscription1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.DeliveryAttributeListResult = armeventgrid.DeliveryAttributeListResult{
	// 	Value: []armeventgrid.DeliveryAttributeMappingClassification{
	// 		&armeventgrid.StaticDeliveryAttributeMapping{
	// 			Name: to.Ptr("header1"),
	// 			Type: to.Ptr(armeventgrid.DeliveryAttributeMappingTypeStatic),
	// 			Properties: &armeventgrid.StaticDeliveryAttributeMappingProperties{
	// 				IsSecret: to.Ptr(false),
	// 				Value: to.Ptr("NormalValue"),
	// 			},
	// 		},
	// 		&armeventgrid.DynamicDeliveryAttributeMapping{
	// 			Name: to.Ptr("header2"),
	// 			Type: to.Ptr(armeventgrid.DeliveryAttributeMappingTypeDynamic),
	// 			Properties: &armeventgrid.DynamicDeliveryAttributeMappingProperties{
	// 				SourceField: to.Ptr("data.foo"),
	// 			},
	// 		},
	// 		&armeventgrid.StaticDeliveryAttributeMapping{
	// 			Name: to.Ptr("header3"),
	// 			Type: to.Ptr(armeventgrid.DeliveryAttributeMappingTypeStatic),
	// 			Properties: &armeventgrid.StaticDeliveryAttributeMappingProperties{
	// 				IsSecret: to.Ptr(true),
	// 				Value: to.Ptr("mySecretValue"),
	// 			},
	// 	}},
	// }
}
