// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package messaging_test

import (
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/messaging"
)

func Example_usingCloudEvent() {
	type sampleType struct {
		CustomField string `json:"custom_field"`
	}

	eventToSend, err := messaging.NewCloudEvent("source", "eventtype", &sampleType{
		CustomField: "hello, a custom field value",
	}, nil)

	if err != nil {
		panic(err)
	}

	receivedEvent, err := sendAndReceiveCloudEvent(eventToSend)

	if err != nil {
		panic(err)
	}

	var receivedData *sampleType

	if err := json.Unmarshal(receivedEvent.Data.([]byte), &receivedData); err != nil {
		panic(err)
	}

	fmt.Printf("Custom field = %s\n", receivedData.CustomField)

	// Output:
	// Custom field = hello, a custom field value
}

func sendAndReceiveCloudEvent(ce messaging.CloudEvent) (messaging.CloudEvent, error) {
	bytes, err := json.Marshal(ce)

	if err != nil {
		return messaging.CloudEvent{}, err
	}

	var received *messaging.CloudEvent

	if err := json.Unmarshal(bytes, &received); err != nil {
		return messaging.CloudEvent{}, err
	}

	return *received, nil
}
