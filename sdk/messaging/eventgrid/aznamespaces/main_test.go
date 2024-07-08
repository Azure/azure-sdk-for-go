//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aznamespaces_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/aznamespaces"
	"github.com/joho/godotenv"
)

const recordingDirectory = "sdk/messaging/eventgrid/aznamespaces/testdata"

func TestMain(m *testing.M) {
	code := run(m)
	os.Exit(code)
}

func run(m *testing.M) int {
	if recording.GetRecordMode() != recording.PlaybackMode {
		if err := godotenv.Load(); err != nil {
			log.Printf("Failed to load .env file, no integration tests will run: %s", err)
		}

		purgeEvents()
	}

	if recording.GetRecordMode() == recording.PlaybackMode || recording.GetRecordMode() == recording.RecordingMode {
		proxy, err := recording.StartTestProxy(recordingDirectory, nil)
		if err != nil {
			panic(err)
		}

		defer func() {
			err := recording.StopTestProxy(proxy)
			if err != nil {
				panic(err)
			}
		}()
	}

	return m.Run()
}

func purgeEvents() {
	testVars, err := loadEnv()

	if err != nil {
		panic(err)
	}

	cred, err := credential.New(nil)

	if err != nil {
		panic(err)
	}

	receiver, err := aznamespaces.NewReceiverClient(testVars.Endpoint, testVars.Topic, testVars.Subscription, cred, &aznamespaces.ReceiverClientOptions{
		ClientOptions: policy.ClientOptions{
			Logging: policy.LogOptions{
				IncludeBody:        true,
				AllowedQueryParams: []string{"maxWaitTime", "maxEvents"},
			},
		},
	})

	if err != nil {
		panic(err)
	}

	log.Printf("(setup) Purging any events in %s/%s...", testVars.Topic, testVars.Subscription)

	for {
		recvResp, err := receiver.ReceiveEvents(context.Background(), &aznamespaces.ReceiveEventsOptions{
			MaxEvents:   to.Ptr[int32](100),
			MaxWaitTime: to.Ptr[int32](10),
		})

		if err != nil {
			panic(err)
		}

		if len(recvResp.Details) == 0 {
			break
		}

		var lockTokens []string

		log.Printf("(setup) Got %d events, deleting...", len(recvResp.Details))

		for _, event := range recvResp.Details {
			lockTokens = append(lockTokens, *event.BrokerProperties.LockToken)
		}

		if _, err := receiver.AcknowledgeEvents(context.Background(), lockTokens, nil); err != nil {
			panic(err)
		}
	}

	log.Printf("(setup) Done purging events from %s/%s...", testVars.Topic, testVars.Subscription)
}
