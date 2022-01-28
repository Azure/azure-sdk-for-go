// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package shared

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/joho/godotenv"
)

type TestContext struct {
	*StressContext
	Client *azservicebus.Client
}

func MustGenerateMessages(sc *StressContext, sender *azservicebus.Sender, messageLimit int, numExtraBytes int, stats *Stats) {
	ctx, cancel := context.WithCancel(sc.Context)
	defer cancel()

	log.Printf("Sending %d messages", messageLimit)

	streamingBatch, err := NewStreamingMessageBatch(ctx, &senderWrapper{inner: sender}, stats)
	sc.PanicOnError("failed to create streaming batch", err)

	extraBytes := make([]byte, numExtraBytes)

	for i := 0; i < messageLimit; i++ {
		err := streamingBatch.Add(ctx, &azservicebus.Message{
			Body: extraBytes,
			ApplicationProperties: map[string]interface{}{
				"Number": i,
			},
		})
		sc.PanicOnError("failed add/sending a batch", err)
	}

	err = streamingBatch.Close(ctx)
	sc.PanicOnError("failed flushing final batch", err)
}

// MustCreateAutoDeletingQueue creates a queue that will auto-delete 10 minutes after activity has ceased.
func MustCreateAutoDeletingQueue(sc *StressContext, queueName string, qp *admin.QueueProperties) {
	adminClient, err := admin.NewClientFromConnectionString(sc.ConnectionString, nil)
	sc.PanicOnError("failed to create adminClient", err)

	autoDeleteOnIdle := 10 * time.Minute

	var newQP admin.QueueProperties

	if qp != nil {
		newQP = *qp
	}

	newQP.AutoDeleteOnIdle = &autoDeleteOnIdle

	// mostly useful for tracking backwards in case something goes wrong.
	newQP.UserMetadata = &sc.TestRunID

	_, err = adminClient.CreateQueue(context.Background(), queueName, &newQP, nil)
	sc.PanicOnError("failed to create queue", err)
}

func MustCreateSubscriptions(sc *StressContext, topicName string, subscriptionNames []string) func() {
	log.Printf("[BEGIN] Creating topic %s", topicName)
	defer log.Printf("[END] Creating topic %s", topicName)

	ac, err := admin.NewClientFromConnectionString(sc.ConnectionString, nil)
	sc.PanicOnError("Failed to create a topic manager", err)

	_, err = ac.CreateTopic(context.Background(), topicName, nil, nil)
	sc.PanicOnError("Failed to create topic", err)

	for _, name := range subscriptionNames {
		_, err := ac.CreateSubscription(context.Background(), topicName, name, nil, nil)
		sc.PanicOnError("Failed to create subscription manager", err)
	}

	return func() {
		_, err := ac.DeleteTopic(context.Background(), topicName, nil)
		sc.PanicOnError(fmt.Sprintf("Failed to delete topic %s", topicName), err)
	}
}

// ConstantlyUpdateQueue updates queue, changing the MaxDeliveryCount properly between 11 and 10, every `updateInterval`
// This will cause Service Bus to issue force-detaches to our links, allowing us to exercise our recovery logic.
func ConstantlyUpdateQueue(ctx context.Context, adminClient *admin.Client, queue string, updateInterval time.Duration) error {
	// updates the entity, which will in turn force a detach for clients.
	ticker := time.NewTicker(updateInterval)

	for range ticker.C {
		if err := ForceQueueDetach(ctx, adminClient, queue); err != nil {
			return err
		}
	}

	return nil
}

func ForceQueueDetach(ctx context.Context, adminClient *admin.Client, queue string) error {
	resp, err := adminClient.GetQueue(ctx, queue, nil)

	if err != nil {
		return err
	}

	if *resp.MaxDeliveryCount == 10 {
		*resp.MaxDeliveryCount = 11
	} else {
		*resp.MaxDeliveryCount = 10
	}

	_, err = adminClient.UpdateQueue(ctx, queue, resp.QueueProperties, nil)

	if err != nil {
		return err
	}

	return nil
}

// LoadEnvironment loads an .env file.
// If the env var `ENV_FILE` exists, we assume the value is a path to an .env file
// Otherwise we fall back to loading from the current directory.
func LoadEnvironment() error {
	var err error
	envFilePath := os.Getenv("ENV_FILE")

	if envFilePath == "" {
		// assume same directory
		err = godotenv.Load()
	} else {
		err = godotenv.Load(envFilePath)
	}

	if err != nil {
		return fmt.Errorf("failed to load .env file from path '%s': %s", envFilePath, err.Error())
	}

	return nil
}

// AddAuthFlags adds the flags needed for authenticating to Service Bus.
// Returns a function that can be called after the flags have been parsed, which will create the an *azservicebus.Client.
func AddAuthFlags(fs *flag.FlagSet) func() (*azservicebus.Client, *admin.Client, error) {
	connectionStringEnvVar := fs.String("cs", "SERVICEBUS_CONNECTION_STRING", "Environment variable containing a connection string for authentication.")
	fullyQualifiedNamespace := fs.String("ns", "", "A Service Bus namespace (ex: <server>.servicebus.windows.net). azidentity.DefaultAzureCredential will be used for authentication.")

	return func() (*azservicebus.Client, *admin.Client, error) {
		var serviceBusClient *azservicebus.Client

		if *fullyQualifiedNamespace != "" {
			// the DefaultAzureCredential will try multiple methods to authenticate, including using cached Azure CLI
			// credentials, pulling authentication variables from the environment and others!
			defaultAzureCredential, err := azidentity.NewDefaultAzureCredential(nil)

			if err != nil {
				return nil, nil, fmt.Errorf("failed to create a DefaultAzureCredential: %w", err)
			}

			serviceBusClient, err = azservicebus.NewClient(*fullyQualifiedNamespace, defaultAzureCredential, nil)

			if err != nil {
				return nil, nil, fmt.Errorf("failed to create an azservicebus.Client using the azidentity.DefaultAzureCredential: %w", err)
			}

			adminClient, err := admin.NewClient(*fullyQualifiedNamespace, defaultAzureCredential, nil)

			if err != nil {
				return nil, nil, fmt.Errorf("failed to create an admin.Client using the azidentity.DefaultAzureCredential: %w", err)
			}

			return serviceBusClient, adminClient, nil
		}

		// assume connection string based authentication, via the environment
		cs := os.Getenv(*connectionStringEnvVar)

		if cs == "" {
			return nil, nil, fmt.Errorf("no connection string in environment variable '%s'", *connectionStringEnvVar)
		}

		var err error
		serviceBusClient, err = azservicebus.NewClientFromConnectionString(cs, nil)

		if err != nil {
			return nil, nil, fmt.Errorf("failed to create an azservicebus.Client using a connection string: %w", err)
		}

		adminClient, err := admin.NewClientFromConnectionString(cs, nil)

		if err != nil {
			return nil, nil, fmt.Errorf("failed to create an admin.Client using the azidentity.DefaultAzureCredential: %w", err)
		}

		return serviceBusClient, adminClient, nil
	}
}

// NewCtrlCContext creates a context that cancels if the user hits ctrl+c.
func NewCtrlCContext() (context.Context, context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-c
		close(c)
		cancel()
	}()

	return ctx, cancel
}
