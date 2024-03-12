//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aznamespaces_test

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/eventgrid/armeventgrid/v2"
	"github.com/joho/godotenv"
)

const recordingDirectory = "sdk/messaging/azeventgrid/testdata"

func TestMain(m *testing.M) {
	code := run(m)
	os.Exit(code)
}

func run(m *testing.M) int {
	defer topicCleanup()
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
	if err := godotenv.Load(); err != nil {
		log.Printf("Failed to load .env file, no integration tests will run: %s", err)
	}

	return m.Run()
}

func PollUntilDone[T any](ctx context.Context, fn func() (*runtime.Poller[T], error)) error {
	poller, err := fn()

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, nil)
	return err
}

func createTopicAndUpdateEnv() func() {
	azSubID := os.Getenv("AZEVENTGRID_SUBSCRIPTION_ID")
	resGroup := os.Getenv("AZEVENTGRID_RESOURCE_GROUP")

	if azSubID == "" || resGroup == "" || recording.GetRecordMode() != recording.LiveMode {
		// ie, these are unit tests.
		return func() {}
	}

	nsURL, err := url.Parse(os.Getenv("EVENTGRID_ENDPOINT"))

	if err != nil {
		panic(err)
	}

	nsHost := strings.Split(nsURL.Host, ".")[0]

	cred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		panic(err)
	}

	topicClient, err := armeventgrid.NewNamespaceTopicsClient(azSubID, cred, nil)

	if err != nil {
		panic(err)
	}

	subClient, err := armeventgrid.NewNamespaceTopicEventSubscriptionsClient(azSubID, cred, nil)

	if err != nil {
		panic(err)
	}

	topicName := fmt.Sprintf("topic-%d", time.Now().UnixNano())

	os.Setenv("EVENTGRID_TOPIC", topicName)

	subName := "testsubscription1"

	err = PollUntilDone(context.Background(), func() (*runtime.Poller[armeventgrid.NamespaceTopicsClientCreateOrUpdateResponse], error) {
		return topicClient.BeginCreateOrUpdate(context.Background(), resGroup, nsHost, topicName, armeventgrid.NamespaceTopic{}, nil)
	})

	if err != nil {
		panic(err)
	}

	err = PollUntilDone(context.Background(), func() (*runtime.Poller[armeventgrid.NamespaceTopicEventSubscriptionsClientCreateOrUpdateResponse], error) {
		return subClient.BeginCreateOrUpdate(context.Background(),
			resGroup,
			nsHost,
			topicName,
			subName,
			armeventgrid.Subscription{
				Properties: &armeventgrid.SubscriptionProperties{
					DeliveryConfiguration: &armeventgrid.DeliveryConfiguration{
						DeliveryMode: to.Ptr(armeventgrid.DeliveryModeQueue),
					},
				},
			},
			nil)
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("Created topic %s\n", topicName)

	return func() {
		if _, err = topicClient.BeginDelete(context.Background(), resGroup, nsHost, topicName, nil); err != nil {
			fmt.Printf("Failed to start the delete for our test topic %s: %s", topicName, err)
			return
		}

		fmt.Printf("Deleted topic %s\n", topicName)
	}
}
