// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestProcessorReceiveWithDefaults(t *testing.T) {
	serviceBusClient, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	go func() {
		sender, err := serviceBusClient.NewSender(queueName)
		require.NoError(t, err)

		defer sender.Close(context.Background())

		// it's perfectly fine to have the processor started before the messages
		// have been sent.
		for i := 0; i < 5; i++ {
			err = sender.SendMessage(context.Background(), &Message{
				Body: []byte(fmt.Sprintf("hello world %d", i)),
			})

			time.Sleep(time.Second)
		}
		require.NoError(t, err)
	}()

	processor, err := newProcessorForQueue(serviceBusClient, queueName, nil)
	require.NoError(t, err)

	defer processor.Close(context.Background()) // multiple close is fine

	var messages []string
	mu := sync.Mutex{}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	err = processor.Start(ctx, func(m *ReceivedMessage) error {
		mu.Lock()
		defer mu.Unlock()

		messages = append(messages, string(m.Body))

		if len(messages) == 5 {
			cancel()
		}

		return nil
	}, func(err error) {
		if errors.Is(err, context.Canceled) {
			return
		}

		require.NoError(t, err)
	})

	require.ErrorIs(t, err, context.Canceled, "cancelling the context stops the processor")

	sort.Strings(messages)

	require.EqualValues(t, []string{
		"hello world 0",
		"hello world 1",
		"hello world 2",
		"hello world 3",
		"hello world 4",
	}, messages)

	require.NoError(t, processor.Close(context.Background()))
}

func TestProcessorReceiveWith100MessagesWithMaxConcurrency(t *testing.T) {
	serviceBusClient, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	const numMessages = 100
	var expectedBodies []string

	go func() {
		sender, err := serviceBusClient.NewSender(queueName)
		require.NoError(t, err)

		defer sender.Close(context.Background())

		batch, err := sender.NewMessageBatch(context.Background(), nil)
		require.NoError(t, err)

		// it's perfectly fine to have the processor started before the messages
		// have been sent.
		for i := 0; i < numMessages; i++ {
			expectedBodies = append(expectedBodies, fmt.Sprintf("hello world %03d", i))
			added, err := batch.Add(&Message{
				Body: []byte(expectedBodies[len(expectedBodies)-1]),
			})
			require.NoError(t, err)
			require.True(t, added)
		}

		require.NoError(t, sender.SendMessageBatch(context.Background(), batch))
	}()

	processor, err := newProcessorForQueue(
		serviceBusClient,
		queueName,
		&processorOptions{
			MaxConcurrentCalls: 20,
		})

	require.NoError(t, err)

	defer func() {
		require.NoError(t, processor.Close(context.Background())) // multiple close is fine
	}()

	var messages []string
	mu := sync.Mutex{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	err = processor.Start(ctx, func(m *ReceivedMessage) error {
		mu.Lock()
		defer mu.Unlock()

		messages = append(messages, string(m.Body))

		if len(messages) == 100 {
			go processor.Close(context.Background())
		}

		return nil
	}, func(err error) {
		if errors.Is(err, context.Canceled) {
			return
		}

		require.NoError(t, err)
	})

	require.NoError(t, err)

	sort.Strings(messages)
	require.EqualValues(t, expectedBodies, messages)

	require.NoError(t, processor.Close(ctx))
}

func TestProcessorUnitTests(t *testing.T) {
	p := &processor{}
	e := &entity{}

	require.NoError(t, applyProcessorOptions(p, e, nil))
	require.True(t, p.autoComplete)
	require.EqualValues(t, 1, p.maxConcurrentCalls)
	require.EqualValues(t, ReceiveModePeekLock, p.receiveMode)

	p = &processor{}
	e = &entity{
		Queue: "queue",
	}

	require.NoError(t, applyProcessorOptions(p, e, &processorOptions{
		ReceiveMode:         ReceiveModeReceiveAndDelete,
		SubQueue:            SubQueueDeadLetter,
		DisableAutoComplete: true,
		MaxConcurrentCalls:  101,
	}))

	require.False(t, p.autoComplete)
	require.EqualValues(t, 101, p.maxConcurrentCalls)
	require.EqualValues(t, ReceiveModeReceiveAndDelete, p.receiveMode)
	fullEntityPath, err := e.String()
	require.NoError(t, err)
	require.EqualValues(t, "queue/$DeadLetterQueue", fullEntityPath)
}

// func TestProcessorUnitTests(t *testing.T) {
// 	t.Run("Processor", func(t *testing.T) {
// 		t.Run("StartAndClose", func(t *testing.T) {

// 		})

// 		t.Run("CloseWaitsForActiveSubscribersToExit", func(t *testing.T) {
// 		})

// 		t.Run("CloseWithoutStart", func(t *testing.T) {

// 		})

// 		t.Run("DoubleClose", func(t *testing.T) {

// 		})
// 	})

// 	t.Run("subscribe", func(t *testing.T) {
// 		t.Run("cancelled by user does not retry", func(t *testing.T) {
// 			ctx, cancel := context.WithCancel(context.Background())
// 			cancel() // pre-cancel this context

// 			receiver := internal.NewFakeLegacyReceiver()
// 			var cancelledError error

// 			retry := subscribe(ctx, receiver, true, func(message *ReceivedMessage) error {
// 				return nil
// 			}, func(err error) {
// 				cancelledError = err
// 			})

// 			require.EqualError(t, cancelledError, context.Canceled.Error())
// 			require.False(t, retry, "User cancelling the context will not be retried")
// 			require.False(t, receiver.CloseCalled) // subscribe() is not responsible for the lifetime of the receiver
// 		})

// 		t.Run("error in the listener is retryable", func(t *testing.T) {
// 			receiver := internal.NewFakeLegacyReceiver()

// 			receiver.ListenImpl = func(ctx context.Context, handler internal.Handler) internal.ListenerHandle {
// 				ch := make(chan struct{})
// 				close(ch)
// 				return &internal.FakeListenerHandle{
// 					DoneChan: ch,
// 					ErrValue: errors.New("Some AMQP related error"),
// 				}
// 			}

// 			var errorFromListener error

// 			retry := subscribe(context.Background(), receiver, true, func(message *ReceivedMessage) error {
// 				return nil
// 			}, func(err error) {
// 				errorFromListener = err
// 			})

// 			require.EqualError(t, errorFromListener, "Some AMQP related error")
// 			require.True(t, retry, "AMQP errors will cause us to retry")
// 			require.False(t, receiver.CloseCalled) // subscribe() is not responsible for the lifetime of the receiver
// 		})

// 	})

// 	t.Run("handleSingleMessage", func(t *testing.T) {
// 		fakeMessage := &internal.Message{
// 			ID:        "fakeID",
// 			LockToken: &uuid.UUID{},
// 			SystemProperties: &internal.SystemProperties{
// 				SequenceNumber: to.Int64Ptr(1),
// 			},
// 		}

// 		setup := func() *internal.FakeInternalReceiver {
// 			receiver := internal.NewFakeLegacyReceiver()
// 			return receiver
// 		}

// 		t.Run("AutoCompleteCompleteMessage", func(t *testing.T) {
// 			receiver := setup()

// 			handleSingleMessage(func(message *ReceivedMessage) error {
// 				// successful return
// 				return nil
// 			}, func(err error) {
// 				require.NoError(t, err)
// 			}, true, receiver, fakeMessage)

// 			require.True(t, receiver.CompleteCalled)
// 			require.False(t, receiver.AbandonCalled)
// 		})

// 		t.Run("AutoCompleteAbandonMessage", func(t *testing.T) {
// 			receiver := setup()

// 			handleSingleMessage(func(message *ReceivedMessage) error {
// 				// error returned will abandon
// 				return errors.New("Purposefully reported error")
// 			}, func(err error) {
// 				require.EqualErrorf(t, err, "Purposefully reported error", "Error from the handler gets forwarded")
// 			}, true, receiver, fakeMessage)

// 			require.True(t, receiver.AbandonCalled)
// 			require.False(t, receiver.CompleteCalled)
// 		})

// 		t.Run("AutoCompleteAlreadySettledDoNotSettleTwice)", func(t *testing.T) {
// 			receiver := setup()

// 			handleSingleMessage(func(message *ReceivedMessage) error {
// 				// error returned will abandon
// 				return errors.New("Purposefully reported error")
// 			}, func(err error) {
// 				require.EqualErrorf(t, err, "Purposefully reported error", "Error from the handler gets forwarded")
// 			}, true, receiver, fakeMessage)

// 			// TODO: neither should be called - the message was already settled.
// 			require.True(t, receiver.AbandonCalled)
// 			require.False(t, receiver.CompleteCalled)
// 		})

// 		t.Run("autoComplete (off)", func(t *testing.T) {
// 			receiver := setup()

// 			handleSingleMessage(func(message *ReceivedMessage) error {
// 				// successful return
// 				return nil
// 			}, func(err error) {
// 				require.NoError(t, err)
// 			}, false, receiver, fakeMessage)

// 			require.False(t, receiver.CompleteCalled)
// 			require.False(t, receiver.AbandonCalled)
// 		})

// 		t.Run("SettlementErrorsAreForwarded(complete)", func(t *testing.T) {
// 			receiver := setup()

// 			receiver.CompleteMessageImpl = func(ctx context.Context, msg *internal.Message) error {
// 				return errors.New("Complete failed")
// 			}

// 			var settleError error

// 			handleSingleMessage(func(message *ReceivedMessage) error {
// 				return nil
// 			}, func(err error) {
// 				settleError = err
// 			}, true, receiver, fakeMessage)

// 			require.EqualError(t, settleError, "Complete failed")
// 		})

// 		t.Run("SettlementErrorsAreForwarded(abandon)", func(t *testing.T) {
// 			receiver := setup()

// 			receiver.AbandonMessageImpl = func(ctx context.Context, msg *internal.Message) error {
// 				return errors.New("Abandon failed")
// 			}

// 			var settleErrors []string

// 			handleSingleMessage(func(message *ReceivedMessage) error {
// 				return errors.New("Error that caused the abandon")
// 			}, func(err error) {
// 				settleErrors = append(settleErrors, err.Error())
// 			}, true, receiver, fakeMessage)

// 			require.EqualValues(t, settleErrors, []string{
// 				"Error that caused the abandon",
// 				"Abandon failed",
// 			})
// 		})
// 	})
// }
