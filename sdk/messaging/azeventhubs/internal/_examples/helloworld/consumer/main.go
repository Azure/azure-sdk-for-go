// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-amqp-common-go/v3/aad"
	mgmt "github.com/Azure/azure-sdk-for-go/services/eventhub/mgmt/2017-04-01/eventhub"
	"github.com/Azure/go-autorest/autorest/azure"
	azauth "github.com/Azure/go-autorest/autorest/azure/auth"

	eventhub "github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal"
)

const (
	Location          = "eastus"
	ResourceGroupName = "ehtest"
	HubName           = "producerConsumer"
)

func main() {
	hub, partitions := initHub()
	exit := make(chan struct{})

	handler := func(ctx context.Context, event *eventhub.Event) error {
		text := string(event.Data)
		if text == "exit\n" {
			fmt.Println("Oh snap!! Someone told me to exit!")
			exit <- *new(struct{})
		} else {
			fmt.Println(string(event.Data))
		}
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	for _, partitionID := range partitions {
		_, err := hub.Receive(ctx, partitionID, handler, eventhub.ReceiveWithLatestOffset())
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}
	cancel()

	fmt.Println("I am listening...")

	select {
	case <-exit:
		fmt.Println("closing after 2 seconds")
		select {
		case <-time.After(2 * time.Second):
			return
		}
	}
}

func initHub() (*eventhub.Hub, []string) {
	namespace := mustGetenv("EVENTHUB_NAMESPACE")
	hubMgmt, err := ensureEventHub(context.Background(), HubName)
	if err != nil {
		log.Fatal(err)
	}

	provider, err := aad.NewJWTProvider(aad.JWTProviderWithEnvironmentVars())
	if err != nil {
		log.Fatal(err)
	}
	hub, err := eventhub.NewHub(namespace, HubName, provider)
	if err != nil {
		panic(err)
	}
	return hub, *hubMgmt.PartitionIds
}

func mustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("Environment variable '" + key + "' required for integration tests.")
	}
	return v
}

func ensureEventHub(ctx context.Context, name string) (*mgmt.Model, error) {
	namespace := mustGetenv("EVENTHUB_NAMESPACE")
	client := getEventHubMgmtClient()
	hub, err := client.Get(ctx, ResourceGroupName, namespace, name)

	partitionCount := int64(4)
	if err != nil {
		newHub := &mgmt.Model{
			Name: &name,
			Properties: &mgmt.Properties{
				PartitionCount: &partitionCount,
			},
		}

		hub, err = client.CreateOrUpdate(ctx, ResourceGroupName, namespace, name, *newHub)
		if err != nil {
			return nil, err
		}
	}
	return &hub, nil
}

func getEventHubMgmtClient() *mgmt.EventHubsClient {
	subID := mustGetenv("AZURE_SUBSCRIPTION_ID")
	client := mgmt.NewEventHubsClientWithBaseURI(azure.PublicCloud.ResourceManagerEndpoint, subID)
	a, err := azauth.NewAuthorizerFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	client.Authorizer = a
	return &client
}
