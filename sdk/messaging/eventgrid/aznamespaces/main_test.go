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
	"github.com/stretchr/testify/require"
)

const recordingDirectory = "sdk/messaging/eventgrid/aznamespaces/testdata"

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

func createTopicAndUpdateEnv(t *testing.T) func() {
	if recording.GetRecordMode() == recording.PlaybackMode {
		return func() {}
	}

	azSubID := os.Getenv("AZNAMESPACES_SUBSCRIPTION_ID")
	require.NotEmpty(t, azSubID, "AZNAMESPACES_SUBSCRIPTION_ID is defined")

	resGroup := os.Getenv("AZNAMESPACES_RESOURCE_GROUP")
	require.NotEmpty(t, resGroup, "AZNAMESPACES_RESOURCE_GROUP is defined")

	nsURL, err := url.Parse(os.Getenv("EVENTGRID_ENDPOINT"))
	require.NoError(t, err)

	nsHost := strings.Split(nsURL.Host, ".")[0]

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(t, err)

	topicClient, err := armeventgrid.NewNamespaceTopicsClient(azSubID, cred, nil)
	require.NoError(t, err)

	subClient, err := armeventgrid.NewNamespaceTopicEventSubscriptionsClient(azSubID, cred, nil)
	require.NoError(t, err)

	topicName := fmt.Sprintf("topic-%d", time.Now().UnixNano())

	os.Setenv("EVENTGRID_TOPIC", topicName)

	subName := "testsubscription1"

	err = PollUntilDone(context.Background(), func() (*runtime.Poller[armeventgrid.NamespaceTopicsClientCreateOrUpdateResponse], error) {
		return topicClient.BeginCreateOrUpdate(context.Background(), resGroup, nsHost, topicName, armeventgrid.NamespaceTopic{}, nil)
	})
	require.NoError(t, err)

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
	require.NoError(t, err)

	fmt.Printf("Created topic %s\n", topicName)

	return func() {
		// NOTE: this cleanup will happen after the entire test is complete - don't use 't'
		if _, err = topicClient.BeginDelete(context.Background(), resGroup, nsHost, topicName, nil); err != nil {
			fmt.Printf("Failed to start the delete for our test topic %s: %s", topicName, err)
			return
		}

		fmt.Printf("Deleted topic %s\n", topicName)
	}
}
