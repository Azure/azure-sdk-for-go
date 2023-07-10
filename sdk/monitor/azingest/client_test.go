//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azingest_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azingest"
)

// for testing, create struct with all the data types
// remove computer field, not supported anymore
type ComputerInfo struct {
	InputTime         time.Time
	Computer          string
	AdditionalContext int
}

// test for greater than 1 mb
// generate a file and read from it

func TestUpload(t *testing.T) {
	azlog.SetListener(func(cls azlog.Event, msg string) {
		fmt.Println(msg)
	})

	endpoint := os.Getenv("MONITOR_INGESTION_DATA_COLLECTION_ENDPOINT")
	ruleID := os.Getenv("INGESTION_DATA_COLLECTION_RULE_IMMUTABLE_ID")
	streamName := os.Getenv("INGESTION_STREAM_NAME")
	clientID := os.Getenv("AZINGESTION_CLIENT_ID")
	clientSecret := os.Getenv("AZINGESTION_CLIENT_SECRET")
	tenantID := os.Getenv("AZINGESTION_TENANT_ID")

	credential, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	if err != nil {
		panic(err)
	}
	client, err := azingest.NewClient(endpoint, credential, &azingest.ClientOptions{azcore.ClientOptions{Logging: policy.LogOptions{IncludeBody: true}}})
	if err != nil {
		panic(err)
	}

	var data []ComputerInfo
	for i := 0; i < 10; i++ {
		data = append(data, ComputerInfo{
			InputTime:         time.Now().UTC(),
			Computer:          "Computer" + strconv.Itoa(i),
			AdditionalContext: i,
		})
	}
	data2, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	_, err = client.Upload(context.Background(), ruleID, streamName, data2, nil)
	if err != nil {
		panic(err)
	}

}
