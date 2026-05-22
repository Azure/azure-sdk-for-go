// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azeventhubs_test

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/internal/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestProducerClient_SAS(t *testing.T) {
	getLogsFn := test.CaptureLogsForTest()

	testParams := test.GetConnectionParamsForTest(t)
	sasCS, err := sas.CreateConnectionStringWithSASUsingExpiry(testParams.CS(t).Primary, time.Now().UTC().Add(time.Hour))
	require.NoError(t, err)

	// sanity check - we did actually generate a connection string with an embedded SharedAccessSignature
	require.Contains(t, sasCS, "SharedAccessSignature=SharedAccessSignature")

	producerClient, err := azeventhubs.NewProducerClientFromConnectionString(sasCS, testParams.EventHubName, nil)
	require.NoError(t, err)

	defer test.RequireClose(t, producerClient)

	consumerClient, err := azeventhubs.NewConsumerClientFromConnectionString(sasCS, testParams.EventHubName, azeventhubs.DefaultConsumerGroup, nil)
	require.NoError(t, err)

	defer test.RequireClose(t, consumerClient)

	beforeProps, err := producerClient.GetPartitionProperties(context.Background(), "0", nil)
	require.NoError(t, err)

	batch, err := producerClient.NewEventDataBatch(context.Background(), &azeventhubs.EventDataBatchOptions{
		PartitionID: to.Ptr("0"),
	})
	require.NoError(t, err)

	err = batch.AddEventData(&azeventhubs.EventData{
		Body: []byte("TestProducerClient_SAS"),
	}, nil)
	require.NoError(t, err)

	err = producerClient.SendEventDataBatch(context.Background(), batch, nil)
	require.NoError(t, err)

	partClient, err := consumerClient.NewPartitionClient("0", &azeventhubs.PartitionClientOptions{
		StartPosition: getStartPosition(beforeProps),
	})
	require.NoError(t, err)

	defer test.RequireClose(t, partClient)

	events, err := partClient.ReceiveEvents(context.Background(), 1, nil)
	require.NoError(t, err)
	require.NotEmpty(t, events)

	logs := getLogsFn()
	require.Contains(t, logs, backgroundRenewalDisabledMsg)
}

const backgroundRenewalDisabledMsg = "[azeh.Auth] Token does not have an expiration date, no background renewal needed."

func TestClientsUnauthorizedCreds(t *testing.T) {
	testParams := test.GetConnectionParamsForTest(t)

	t.Run("ListenOnly with Producer", func(t *testing.T) {
		pc, err := azeventhubs.NewProducerClientFromConnectionString(testParams.CS(t).ListenOnly, testParams.EventHubName, nil)
		require.NoError(t, err)
		defer test.RequireClose(t, pc)

		batch, err := pc.NewEventDataBatch(context.Background(), nil)

		var ehErr *azeventhubs.Error
		require.ErrorAs(t, err, &ehErr)
		require.Equal(t, azeventhubs.ErrorCodeUnauthorizedAccess, ehErr.Code)
		require.Contains(t, err.Error(), "Description: Unauthorized access. 'Send' claim(s) are required to perform this operation")
		require.Nil(t, batch)
	})

	t.Run("SendOnly with Consumer", func(t *testing.T) {
		client, err := azeventhubs.NewConsumerClientFromConnectionString(testParams.CS(t).SendOnly, testParams.EventHubName, azeventhubs.DefaultConsumerGroup, nil)
		require.NoError(t, err)
		defer test.RequireClose(t, client)

		pc, err := client.NewPartitionClient("0", nil)
		require.NoError(t, err)
		defer test.RequireClose(t, pc)

		events, err := pc.ReceiveEvents(context.Background(), 1, nil)

		var ehErr *azeventhubs.Error
		require.ErrorAs(t, err, &ehErr)
		require.Equal(t, azeventhubs.ErrorCodeUnauthorizedAccess, ehErr.Code)
		require.Contains(t, err.Error(), "Description: Unauthorized access. 'Listen' claim(s) are required to perform this operation")
		require.Empty(t, events)
	})

	t.Run("Expired SAS", func(t *testing.T) {
		expiredCS, err := sas.CreateConnectionStringWithSASUsingExpiry(testParams.CS(t).Primary, time.Now().Add(-10*time.Minute))
		require.NoError(t, err)

		cc, err := azeventhubs.NewConsumerClientFromConnectionString(expiredCS, testParams.EventHubName, azeventhubs.DefaultConsumerGroup, nil)
		require.NoError(t, err)
		defer test.RequireClose(t, cc)

		pc, err := cc.NewPartitionClient("0", nil)
		require.NoError(t, err)
		defer test.RequireClose(t, pc)

		events, err := pc.ReceiveEvents(context.Background(), 1, nil)

		var ehErr *azeventhubs.Error
		require.ErrorAs(t, err, &ehErr)
		require.Equal(t, azeventhubs.ErrorCodeUnauthorizedAccess, ehErr.Code)
		require.Contains(t, err.Error(), "rpc: failed, status code 401 and description: ExpiredToken: The token is expired. Expiration time:")
		require.Empty(t, events)

		prodClient, err := azeventhubs.NewProducerClientFromConnectionString(expiredCS, testParams.EventHubName, nil)
		require.NoError(t, err)
		defer test.RequireClose(t, prodClient)

		batch, err := prodClient.NewEventDataBatch(context.Background(), nil)
		require.ErrorAs(t, err, &ehErr)
		require.Equal(t, azeventhubs.ErrorCodeUnauthorizedAccess, ehErr.Code)
		require.Contains(t, err.Error(), "rpc: failed, status code 401 and description: ExpiredToken: The token is expired. Expiration time:")
		require.Nil(t, batch)
	})

	t.Run("invalid identity creds", func(t *testing.T) {
		var cred fakeTokenCred

		prodClient, err := azeventhubs.NewProducerClient(testParams.EventHubNamespace, testParams.EventHubName, cred, nil)
		require.NoError(t, err)
		defer test.RequireClose(t, prodClient)

		batch, err := prodClient.NewEventDataBatch(context.Background(), nil)
		var authFailedErr *azidentity.AuthenticationFailedError
		require.ErrorAs(t, err, &authFailedErr)
		require.Nil(t, batch)

		cc, err := azeventhubs.NewConsumerClient(testParams.EventHubNamespace, testParams.EventHubName, azeventhubs.DefaultConsumerGroup, cred, nil)
		require.NoError(t, err)
		defer test.RequireClose(t, cc)

		partClient, err := cc.NewPartitionClient("0", nil)
		require.NoError(t, err)

		events, err := partClient.ReceiveEvents(context.Background(), 1, nil)
		require.ErrorAs(t, err, &authFailedErr)
		require.Nil(t, batch)
		require.Empty(t, events)
	})
}

type fakeTokenCred struct{}

var _ azcore.TokenCredential = fakeTokenCred{}

func (tc fakeTokenCred) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{}, &azidentity.AuthenticationFailedError{}
}

func TestProducerClient_GetHubAndPartitionProperties(t *testing.T) {
	getLogsFn := test.CaptureLogsForTest()
	testParams := test.GetConnectionParamsForTest(t)

	producer, err := azeventhubs.NewProducerClient(testParams.EventHubNamespace, testParams.EventHubName, testParams.Cred, nil)
	require.NoError(t, err)

	defer func() {
		err := producer.Close(context.Background())
		require.NoError(t, err)
	}()

	hubProps, err := producer.GetEventHubProperties(context.Background(), nil)
	require.NoError(t, err)
	require.NotEmpty(t, hubProps.PartitionIDs)

	wg := sync.WaitGroup{}

	for _, partitionID := range hubProps.PartitionIDs {
		wg.Add(1)

		go func(pid string) {
			defer wg.Done()

			t.Run(fmt.Sprintf("Partition%s", pid), func(t *testing.T) {
				sendAndReceiveToPartitionTest(t, testParams, pid)
			})
		}(partitionID)
	}

	wg.Wait()
	logs := getLogsFn()
	checkForTokenRefresh(t, logs, testParams.EventHubName)
}

// checkForTokenRefresh just makes sure that background token refresh has been started
// and that we haven't somehow fallen into the trap of marking all tokens are expired.
func checkForTokenRefresh(t *testing.T, logs []string, eventHubName string) {
	require.NotContains(t, logs, backgroundRenewalDisabledMsg)

	for _, log := range logs {
		if strings.HasPrefix(log, fmt.Sprintf("[azeh.Auth] (%s/$management) next refresh in ", eventHubName)) {
			return
		}
	}
	require.Failf(t, "No token negotiation log lines", "logs:%s", strings.Join(logs, "\n"))
}

func TestProducerClient_GetEventHubsProperties(t *testing.T) {
	testParams := test.GetConnectionParamsForTest(t)

	producer, err := azeventhubs.NewProducerClient(testParams.EventHubNamespace, testParams.EventHubName, testParams.Cred, nil)
	require.NoError(t, err)

	defer func() {
		err := producer.Close(context.Background())
		require.NoError(t, err)
	}()

	props, err := producer.GetEventHubProperties(context.Background(), nil)
	require.NoError(t, err)
	require.NotEmpty(t, props)
	require.NotEmpty(t, props.PartitionIDs)

	for _, pid := range props.PartitionIDs {
		props, err := producer.GetPartitionProperties(context.Background(), pid, nil)

		require.NoError(t, err)
		require.NotEmpty(t, props)

		require.Equal(t, pid, props.PartitionID)
	}
}

func TestProducerClient_SendToAny(t *testing.T) {
	// there are two ways to "send to any" partition
	// 1. Don't specify a partition ID or a partition key when creating the batch
	// 2. Specify a partition key. This is useful if you want to send events and have them
	//    be placed into the same partition but let the overall distribution of the partition keys
	//    happen through Event Hubs.

	t.Run("no partition key, no client instanceID", func(t *testing.T) {
		testSendAny(t, struct {
			testName     string
			instanceID   string
			partitionKey *string
		}{testName: "no partition key, no client instanceID"})
	})

	t.Run("no partition key, with client instanceID", func(t *testing.T) {
		testSendAny(t, struct {
			testName     string
			instanceID   string
			partitionKey *string
		}{
			testName:   "no partition key, with client instanceID",
			instanceID: "client ID",
		})
	})

	t.Run("actual partition key, no client instanceID", func(t *testing.T) {
		testSendAny(t, struct {
			testName     string
			instanceID   string
			partitionKey *string
		}{
			testName:     "actual partition key, no client instanceID",
			partitionKey: to.Ptr("my special partition key"),
		})
	})

	t.Run("actual partition key, with client instanceID", func(t *testing.T) {
		testSendAny(t, struct {
			testName     string
			instanceID   string
			partitionKey *string
		}{
			testName:     "actual partition key, with client instanceID",
			instanceID:   "client ID",
			partitionKey: to.Ptr("my special partition key"),
		})
	})
}

func testSendAny(t *testing.T, args struct {
	testName     string
	instanceID   string
	partitionKey *string
}) {
	testParams := test.GetConnectionParamsForTest(t)

	producer, err := azeventhubs.NewProducerClient(testParams.EventHubNamespace, testParams.EventHubName, testParams.Cred, nil)
	require.NoError(t, err)

	defer test.RequireClose(t, producer)

	batch, err := producer.NewEventDataBatch(context.Background(), &azeventhubs.EventDataBatchOptions{
		PartitionKey: args.partitionKey,
	})
	require.NoError(t, err)

	err = batch.AddEventData(&azeventhubs.EventData{
		Body:          []byte(args.testName),
		ContentType:   to.Ptr("content type"),
		CorrelationID: "correlation id",
		MessageID:     to.Ptr("message id"),
		Properties: map[string]any{
			"hello": "world",
		},
	}, nil)
	require.NoError(t, err)

	partitionsBeforeSend := getAllPartitionProperties(t, producer)

	err = producer.SendEventDataBatch(context.Background(), batch, nil)
	require.NoError(t, err)

	consumer, err := azeventhubs.NewConsumerClient(testParams.EventHubNamespace, testParams.EventHubName, azeventhubs.DefaultConsumerGroup, testParams.Cred, &azeventhubs.ConsumerClientOptions{
		InstanceID: args.instanceID,
	})
	require.NoError(t, err)

	defer test.RequireClose(t, consumer)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	receivedEvent := receiveEventFromAnyPartition(ctx, t, consumer, partitionsBeforeSend)

	require.Equal(t, azeventhubs.EventData{
		Body:          []byte(args.testName),
		ContentType:   to.Ptr("content type"),
		CorrelationID: "correlation id",
		MessageID:     to.Ptr("message id"),
		Properties: map[string]any{
			"hello": "world",
		}}, receivedEvent.EventData)

	require.GreaterOrEqual(t, receivedEvent.SequenceNumber, int64(0))
	require.NotNil(t, receivedEvent.Offset)
	require.NotZero(t, receivedEvent.EnqueuedTime)

	if args.partitionKey == nil {
		require.Nil(t, receivedEvent.PartitionKey)
	} else {
		require.NotNil(t, receivedEvent.PartitionKey)
		require.Equal(t, *args.partitionKey, *receivedEvent.PartitionKey)
	}
}

func TestProducerClient_AMQPAnnotatedMessages(t *testing.T) {
	testParams := test.GetConnectionParamsForTest(t)

	producer, err := azeventhubs.NewProducerClient(testParams.EventHubNamespace, testParams.EventHubName, testParams.Cred, nil)
	require.NoError(t, err)

	defer test.RequireClose(t, producer)

	beforeProps, err := producer.GetPartitionProperties(context.Background(), "0", nil)
	require.NoError(t, err)

	numEvents := int64(0)

	// send the events we need, encoding several AMQP body types and exercising all the fields.
	{
		batch, err := producer.NewEventDataBatch(context.Background(), &azeventhubs.EventDataBatchOptions{
			PartitionID: to.Ptr("0"),
		})
		require.NoError(t, err)

		// AMQP messages

		// sequence body
		err = batch.AddAMQPAnnotatedMessage(&azeventhubs.AMQPAnnotatedMessage{
			Body: azeventhubs.AMQPAnnotatedMessageBody{
				Sequence: [][]any{
					{"hello", "world"},
					{"howdy", "world"},
				},
			},
			ApplicationProperties: map[string]any{
				"appProperty1": "appProperty1Value",
			},
			// It doesn't appear that we can't round-trip these attributes:
			// Issue: https://github.com/Azure/azure-sdk-for-go/issues/19154
			DeliveryAnnotations: map[any]any{
				"deliveryAnnotation1": "deliveryAnnotation1Value",
			},
			Header: &azeventhubs.AMQPAnnotatedMessageHeader{
				DeliveryCount: 100,
			},
			Footer: map[any]any{
				"footerField1": "footerValue1",
			},
			MessageAnnotations: map[any]any{
				"messageAnnotation1": 101,
			},
			Properties: &azeventhubs.AMQPAnnotatedMessageProperties{
				GroupID: to.Ptr("custom-group-id"),
			},
		}, nil)
		require.NoError(t, err)

		// value body
		err = batch.AddAMQPAnnotatedMessage(&azeventhubs.AMQPAnnotatedMessage{
			Body: azeventhubs.AMQPAnnotatedMessageBody{
				Value: 999,
			},
		}, nil)
		require.NoError(t, err)

		// data body (multiple arrays, will be 'nil' in a normal ReceivedEventData)
		err = batch.AddAMQPAnnotatedMessage(&azeventhubs.AMQPAnnotatedMessage{
			Body: azeventhubs.AMQPAnnotatedMessageBody{
				Data: [][]byte{
					[]byte("hello"),
					[]byte("world"),
				},
			},
		}, nil)
		require.NoError(t, err)

		err = producer.SendEventDataBatch(context.Background(), batch, nil)
		require.NoError(t, err)

		numEvents = int64(batch.NumEvents())
	}

	afterProps, err := producer.GetPartitionProperties(context.Background(), "0", nil)
	require.NoError(t, err)

	require.Equal(t, numEvents, afterProps.LastEnqueuedSequenceNumber-beforeProps.LastEnqueuedSequenceNumber)

	consumerClient, err := azeventhubs.NewConsumerClient(testParams.EventHubNamespace, testParams.EventHubName, azeventhubs.DefaultConsumerGroup, testParams.Cred, nil)
	require.NoError(t, err)

	defer test.RequireClose(t, consumerClient)

	partitionClient, err := consumerClient.NewPartitionClient("0", &azeventhubs.PartitionClientOptions{
		StartPosition: getStartPosition(beforeProps),
	})
	require.NoError(t, err)

	defer test.RequireClose(t, partitionClient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	receivedEvents, err := partitionClient.ReceiveEvents(ctx, int(numEvents), nil)
	require.NoError(t, err)

	// all of these events have Body encodings that can't be represented in EventData,
	// so it'll be 'nil'. Users can get it from the inner RawAMQPMessage instead.
	for _, e := range receivedEvents {
		require.Nil(t, e.Body)
	}

	sequenceMessage, valueMessage, multiarrayDataMessage := receivedEvents[0].RawAMQPMessage, receivedEvents[1].RawAMQPMessage, receivedEvents[2].RawAMQPMessage

	require.Equal(t, [][]any{
		{"hello", "world"},
		{"howdy", "world"},
	}, sequenceMessage.Body.Sequence)

	require.Equal(t, map[string]any{
		"appProperty1": "appProperty1Value",
	}, sequenceMessage.ApplicationProperties)

	// It doesn't appear that we can round-trip this attribute:
	// https://github.com/Azure/azure-sdk-for-go/issues/19154
	// require.Equal(t, uint32(101), sequenceMessage.Header.DeliveryCount)
	// require.Equal(t, map[any]any{
	// 	"deliveryAnnotation1": "deliveryAnnotation1Value",
	// }, sequenceMessage.DeliveryAnnotations)

	require.Equal(t, map[any]any{
		"footerField1": "footerValue1",
	}, sequenceMessage.Footer)

	require.Equal(t, int64(101), sequenceMessage.MessageAnnotations["messageAnnotation1"])
	require.Equal(t, "custom-group-id", *sequenceMessage.Properties.GroupID)

	require.Equal(t, int64(999), valueMessage.Body.Value)

	// data body (multiple arrays, will be 'nil' in a normal ReceivedEventData)
	require.Equal(t, [][]byte{
		[]byte("hello"),
		[]byte("world"),
	}, multiarrayDataMessage.Body.Data)

	require.Equal(t, int(numEvents), len(receivedEvents))
}

func TestProducerClient_SendBatchExample(t *testing.T) {
	testParams := test.GetConnectionParamsForTest(t)

	producerClient, err := azeventhubs.NewProducerClient(testParams.EventHubNamespace, testParams.EventHubName, testParams.Cred, nil)
	require.NoError(t, err)

	beforeSend, err := producerClient.GetPartitionProperties(context.Background(), "0", nil)
	require.NoError(t, err)

	// this is a replicate of the code we use in the example "example_producer_events.go"
	// just testing to make sure it works the way we expect it to.
	newBatchOptions := &azeventhubs.EventDataBatchOptions{
		MaxBytes:    300,
		PartitionID: to.Ptr("0"),
	}

	const messageSize = 40

	events := []*azeventhubs.EventData{
		{
			Body: makeByteSlice(0, messageSize),
		},
		{
			Body: makeByteSlice(1, messageSize),
		},
		{
			Body: makeByteSlice(2, messageSize),
		},
		{
			Body: makeByteSlice(3, messageSize),
		},
		{
			Body: makeByteSlice(4, messageSize),
		},
	}

	batchesSentFromExcess := 0
	var numMessagesPerBatch []int

	// (example start)
	batch, err := producerClient.NewEventDataBatch(context.TODO(), newBatchOptions)

	if err != nil {
		panic(err)
	}

	for i := 0; i < len(events); i++ {
		err = batch.AddEventData(events[i], nil)

		if errors.Is(err, azeventhubs.ErrEventDataTooLarge) {
			if batch.NumEvents() == 0 {
				// This one event is too large for this batch, even on its own. No matter what we do it
				// will not be sendable at its current size.
				panic(err)
			}

			// This batch is full - we can send it and create a new one and continue
			// packaging and sending events.
			if err := producerClient.SendEventDataBatch(context.TODO(), batch, nil); err != nil {
				panic(err)
			}

			numMessagesPerBatch = append(numMessagesPerBatch, int(batch.NumEvents()))

			tmpBatch, err := producerClient.NewEventDataBatch(context.TODO(), newBatchOptions)

			if err != nil {
				panic(err)
			}

			batch = tmpBatch

			// rewind so we can retry adding this event to a batch
			i--
		} else if err != nil {
			panic(err)
		}
	}

	// if we have any events in the last batch, send it
	if batch.NumEvents() > 0 {
		if err := producerClient.SendEventDataBatch(context.TODO(), batch, nil); err != nil {
			panic(err)
		}

		batchesSentFromExcess++
	}
	// (example end)

	require.Equal(t, 1, batchesSentFromExcess)
	require.Equal(t, []int{2, 2}, numMessagesPerBatch)

	consumerClient, err := azeventhubs.NewConsumerClient(testParams.EventHubNamespace, testParams.EventHubName, azeventhubs.DefaultConsumerGroup, testParams.Cred, nil)
	require.NoError(t, err)

	defer func() { _ = consumerClient.Close(context.Background()) }()

	partitionClient, err := consumerClient.NewPartitionClient("0", &azeventhubs.PartitionClientOptions{
		StartPosition: getStartPosition(beforeSend),
	})
	require.NoError(t, err)

	defer func() { _ = partitionClient.Close(context.Background()) }()

	receivedEvents, err := partitionClient.ReceiveEvents(context.Background(), 5, nil)
	require.NoError(t, err)

	sort.Slice(events, func(i, j int) bool {
		return strings.Compare(string(receivedEvents[i].Body), string(receivedEvents[j].Body)) < 0
	})

	for i := 0; i < 5; i++ {
		require.Equal(t, string(makeByteSlice(i, messageSize)), string(receivedEvents[i].Body))
	}
}

func TestProducerClientUsingCustomEndpoint(t *testing.T) {
	testParams := test.GetConnectionParamsForTest(t)

	producerClient, err := azeventhubs.NewProducerClient(testParams.EventHubNamespace, testParams.EventHubName, testParams.Cred, &azeventhubs.ProducerClientOptions{
		CustomEndpoint: "127.0.0.1",
		RetryOptions: azeventhubs.RetryOptions{
			MaxRetries: -1,
		},
	})
	require.NoError(t, err)

	_, err = producerClient.NewEventDataBatch(context.Background(), nil)

	// NOTE, this is a little silly, but we just want to prove
	// that CustomEndpoint does get used as the actual TCP endpoint we connect to.
	require.Contains(t, err.Error(), "127.0.0.1:5671")
}

func makeByteSlice(index int, total int) []byte {
	// ie: %0<total>d, so it'll be zero padded up to the length we want
	text := fmt.Sprintf("%0"+fmt.Sprintf("%d", total)+"d", index)
	return []byte(text)
}

// receiveEventFromAnyPartition returns when it receives an event from any partition. Useful for tests where you're
// letting the service route the event and you're not sure where it'll end up.
func receiveEventFromAnyPartition(ctx context.Context, t *testing.T, consumer *azeventhubs.ConsumerClient, allPartitions []azeventhubs.PartitionProperties) *azeventhubs.ReceivedEventData {
	eventCh := make(chan *azeventhubs.ReceivedEventData, len(allPartitions))

	partitionClientContext, cancelPartitionReceiving := context.WithCancel(ctx)
	defer cancelPartitionReceiving()

	for _, partProps := range allPartitions {
		go func(ctx context.Context, partProps azeventhubs.PartitionProperties) {
			partClient, err := consumer.NewPartitionClient(partProps.PartitionID, &azeventhubs.PartitionClientOptions{
				StartPosition: getStartPosition(partProps),
			})
			require.NoError(t, err)

			defer test.RequireClose(t, partClient)

			events, err := partClient.ReceiveEvents(ctx, 1, nil)

			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
					t.Logf("Cancellation for partition %s", partProps.PartitionID)
					eventCh <- nil
					return
				}

				require.NoError(t, err)
			}

			if len(events) >= 1 {
				t.Logf("Event found on partition %s", partProps.PartitionID)
				eventCh <- events[0]
				cancelPartitionReceiving()
			}

			t.Logf("No error returned and no events received")
			eventCh <- nil
		}(partitionClientContext, partProps)
	}

	var receivedEvent *azeventhubs.ReceivedEventData

	for i := 0; i < len(allPartitions); i++ {
		select {
		case evt := <-eventCh:
			if evt != nil {
				t.Logf("Event received, waiting for other receives to exit")
				receivedEvent = evt
			}
		case <-ctx.Done():
			require.Fail(t, "No event received!")
			return nil
		}
	}

	require.NotNil(t, receivedEvent)
	return receivedEvent
}

func getAllPartitionProperties(t *testing.T, client interface {
	GetEventHubProperties(ctx context.Context, options *azeventhubs.GetEventHubPropertiesOptions) (azeventhubs.EventHubProperties, error)
	GetPartitionProperties(ctx context.Context, partitionID string, options *azeventhubs.GetPartitionPropertiesOptions) (azeventhubs.PartitionProperties, error)
}) []azeventhubs.PartitionProperties {
	hubProps, err := client.GetEventHubProperties(context.Background(), nil)
	require.NoError(t, err)

	var partitions []azeventhubs.PartitionProperties

	for _, partitionID := range hubProps.PartitionIDs {
		partProps, err := client.GetPartitionProperties(context.Background(), partitionID, nil)
		require.NoError(t, err)

		partitions = append(partitions, partProps)
	}

	sort.Slice(partitions, func(i, j int) bool {
		return partitions[i].PartitionID < partitions[j].PartitionID
	})

	return partitions
}

func sendAndReceiveToPartitionTest(t *testing.T, testParams test.ConnectionParamsForTest, partitionID string) {
	producer, err := azeventhubs.NewProducerClient(testParams.EventHubNamespace, testParams.EventHubName, testParams.Cred, nil)
	require.NoError(t, err)

	defer func() {
		err := producer.Close(context.Background())
		require.NoError(t, err)
	}()

	partProps, err := producer.GetPartitionProperties(context.Background(), partitionID, &azeventhubs.GetPartitionPropertiesOptions{})
	require.NoError(t, err)

	consumer, err := azeventhubs.NewConsumerClient(testParams.EventHubNamespace, testParams.EventHubName, azeventhubs.DefaultConsumerGroup, testParams.Cred, nil)
	require.NoError(t, err)

	defer func() {
		err := consumer.Close(context.Background())
		require.NoError(t, err)
	}()

	batch, err := producer.NewEventDataBatch(context.Background(), &azeventhubs.EventDataBatchOptions{
		PartitionID: &partitionID,
	})
	require.NoError(t, err)

	runID := time.Now().UnixNano()
	var expectedBodies []string

	for i := 0; i < 200; i++ {
		msg := fmt.Sprintf("%05d", i)

		err = batch.AddEventData(&azeventhubs.EventData{
			Body: []byte(msg),
			Properties: map[string]any{
				"PartitionID": partitionID,
				"RunID":       runID,
			},
		}, nil)
		require.NoError(t, err)

		expectedBodies = append(expectedBodies, msg)
	}

	err = producer.SendEventDataBatch(context.Background(), batch, nil)
	require.NoError(t, err)

	// give us 60 seconds to receive all 100 messages we sent in the batch
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var actualBodies []string

	subscription, err := consumer.NewPartitionClient(partitionID, &azeventhubs.PartitionClientOptions{
		StartPosition: getStartPosition(partProps),
	})
	require.NoError(t, err)

	for {
		events, err := subscription.ReceiveEvents(ctx, 100, nil)
		require.NoError(t, err)

		for _, event := range events {
			actualBodies = append(actualBodies, string(event.Body))

			require.Equal(t, partitionID, event.Properties["PartitionID"], "No messages from other partitions")
			require.Equal(t, runID, event.Properties["RunID"], "No messages from older runs")
		}

		if len(actualBodies) == len(expectedBodies) {
			break
		}
	}

	sort.Strings(actualBodies)
	require.Equal(t, expectedBodies, actualBodies)
}
