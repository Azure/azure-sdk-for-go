// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsystemevents_test

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/messaging"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/azsystemevents"
)

func Example_deserializingEventGridSchema() {
	// The method for extracting the payload will be different, depending on which data store you configured
	// as a data handler. For instance, if you used Service Bus, the payload would be in the azservicebus.Message.Body
	// field, as a []byte.

	// This particular payload is in the Event Grid Schema format, so we'll use the EventGridEvent to deserialize it.
	payload := []byte(`[
		{
			"id": "2d1781af-3a4c-4d7c-bd0c-e34b19da4e66",
			"topic": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
			"subject": "mySubject",
			"data": {
				"validationCode": "512d38b6-c7b8-40c8-89fe-f46f9e9622b6",
				"validationUrl": "https://rp-eastus2.eventgrid.azure.net:553/eventsubscriptions/estest/validate?id=B2E34264-7D71-453A-B5FB-B62D0FDC85EE&t=2018-04-26T20:30:54.4538837Z&apiVersion=2018-05-01-preview&token=1BNqCxBBSSE9OnNSfZM4%2b5H9zDegKMY6uJ%2fO2DFRkwQ%3d"
			},
			"eventType": "Microsoft.EventGrid.SubscriptionValidationEvent",
			"eventTime": "2018-01-25T22:12:19.4556811Z",
			"metadataVersion": "1",
			"dataVersion": "1"
		}
	]`)

	var eventGridEvents []azsystemevents.EventGridEvent

	err := json.Unmarshal(payload, &eventGridEvents)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	for _, envelope := range eventGridEvents {
		switch *envelope.EventType {
		case azsystemevents.TypeSubscriptionValidation:
			var eventData *azsystemevents.SubscriptionValidationEventData

			if err := json.Unmarshal(envelope.Data.([]byte), &eventData); err != nil {
				//  TODO: Update the following line with your application specific error handling logic
				log.Fatalf("ERROR: %s", err)
			}

			// print a field out from the event, showing what data is there.
			fmt.Printf("Validation code: %s\n", *eventData.ValidationCode)
			fmt.Printf("Validation URL: %s\n", *eventData.ValidationURL)
		default:
			//  TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: event type %s isn't handled", *envelope.EventType)
		}
	}

	// Output:
	// Validation code: 512d38b6-c7b8-40c8-89fe-f46f9e9622b6
	// Validation URL: https://rp-eastus2.eventgrid.azure.net:553/eventsubscriptions/estest/validate?id=B2E34264-7D71-453A-B5FB-B62D0FDC85EE&t=2018-04-26T20:30:54.4538837Z&apiVersion=2018-05-01-preview&token=1BNqCxBBSSE9OnNSfZM4%2b5H9zDegKMY6uJ%2fO2DFRkwQ%3d
}

func Example_deserializingCloudEventSchema() {
	// The method for extracting the payload will be different, depending on which data store you configured
	// as a data handler. For instance, if you used Service Bus, the payload would be in the azservicebus.Message.Body
	// field, as a []byte.

	// This particular payload is in the Cloud Event Schema format, so we'll use the messaging.CloudEvent, which comes from the `azcore` package,  to deserialize it.
	payload := []byte(`[
		{
			"specversion": "1.0",
			"id": "2d1781af-3a4c-4d7c-bd0c-e34b19da4e66",
			"source": "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
			"subject": "",
			"data": {
				"validationCode": "512d38b6-c7b8-40c8-89fe-f46f9e9622b6",
				"validationUrl": "https://rp-eastus2.eventgrid.azure.net:553/eventsubscriptions/estest/validate?id=B2E34264-7D71-453A-B5FB-B62D0FDC85EE&t=2018-04-26T20:30:54.4538837Z&apiVersion=2018-05-01-preview&token=1BNqCxBBSSE9OnNSfZM4%2b5H9zDegKMY6uJ%2fO2DFRkwQ%3d"
			},
			"type": "Microsoft.EventGrid.SubscriptionValidationEvent",
			"time": "2018-01-25T22:12:19.4556811Z",
			"specversion": "1.0"
		}
	]`)

	var cloudEvents []messaging.CloudEvent

	err := json.Unmarshal(payload, &cloudEvents)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	for _, envelope := range cloudEvents {
		switch envelope.Type {
		case string(azsystemevents.TypeSubscriptionValidation):
			var eventData *azsystemevents.SubscriptionValidationEventData

			if err := json.Unmarshal(envelope.Data.([]byte), &eventData); err != nil {
				//  TODO: Update the following line with your application specific error handling logic
				log.Fatalf("ERROR: %s", err)
			}

			// print a field out from the event, showing what data is there.
			fmt.Printf("Validation code: %s\n", *eventData.ValidationCode)
			fmt.Printf("Validation URL: %s\n", *eventData.ValidationURL)
		default:
			//  TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: event type %s isn't handled", envelope.Type)
		}
	}

	// Output:
	// Validation code: 512d38b6-c7b8-40c8-89fe-f46f9e9622b6
	// Validation URL: https://rp-eastus2.eventgrid.azure.net:553/eventsubscriptions/estest/validate?id=B2E34264-7D71-453A-B5FB-B62D0FDC85EE&t=2018-04-26T20:30:54.4538837Z&apiVersion=2018-05-01-preview&token=1BNqCxBBSSE9OnNSfZM4%2b5H9zDegKMY6uJ%2fO2DFRkwQ%3d
}
