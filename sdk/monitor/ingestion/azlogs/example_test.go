// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azlogs_test

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/ingestion/azlogs"
)

var client azlogs.Client

type Computer struct {
	Time              time.Time
	Computer          string
	AdditionalContext string
}

func ExampleNewClient() {
	endpoint = os.Getenv("DATA_COLLECTION_ENDPOINT")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		//TODO: handle error
	}

	client, err := azlogs.NewClient(endpoint, cred, nil)
	if err != nil {
		//TODO: handle error
	}
	_ = client
}

func ExampleClient_Upload() {
	// set necessary data collection rule variables
	ruleID := os.Getenv("DATA_COLLECTION_RULE_IMMUTABLE_ID")
	streamName := os.Getenv("DATA_COLLECTION_RULE_STREAM_NAME")

	// generating logs
	// logs should match the schema defined by the provided stream
	var data []Computer
	for i := 0; i < 10; i++ {
		data = append(data, Computer{
			Time:              time.Now().UTC(),
			Computer:          "Computer" + strconv.Itoa(i),
			AdditionalContext: "context",
		})
	}
	// Marshal data into []byte
	logs, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	// upload logs
	_, err = client.Upload(context.TODO(), ruleID, streamName, logs, nil)
	if err != nil {
		//TODO: handle error
	}
}
