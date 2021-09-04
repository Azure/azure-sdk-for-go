package azservicebus

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-amqp-common-go/v3/uuid"
	"github.com/Azure/azure-sdk-for-go/sdk/servicebus/azservicebus/internal"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/stretchr/testify/require"
)

func createQueue(t *testing.T, connectionString string, queueName string) func() {
	ns, err := internal.NewNamespace(internal.NamespaceWithConnectionString(connectionString))
	require.NoError(t, err)

	qm := ns.NewQueueManager()

	_, err = qm.Put(context.TODO(), queueName)
	require.NoError(t, err)

	return func() {
		if err := qm.Delete(context.TODO(), queueName); err != nil {
			require.NoError(t, err)
		}
	}
}

func TestProcessor(t *testing.T) {
	cs := os.Getenv("SERVICEBUS_CONNECTION_STRING")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	serviceBusClient, err := NewServiceBusClient(ServiceBusWithConnectionString(cs))
	require.NoError(t, err)

	nanoSeconds := time.Now().UnixNano()
	queueName := fmt.Sprintf("queue-%X", nanoSeconds)
	cleanupQueue := createQueue(t, cs, queueName)
	defer cleanupQueue()

	t.Run("ReceiveMessagesUsingProcessor", func(t *testing.T) {
		sender, err := serviceBusClient.NewSender(queueName)
		require.NoError(t, err)

		err = sender.SendMessage(ctx, &Message{
			Body: []byte("hello world"),
		})

		require.NoError(t, err)

		processor, err := serviceBusClient.NewProcessor(ProcessorWithQueue(queueName))
		require.NoError(t, err)

		defer processor.Close(ctx)

		messagesCh := make(chan *ReceivedMessage, 1)

		err = processor.Start(func(message *ReceivedMessage) error {
			select {
			case messagesCh <- message:
				break
			default:
				return fmt.Errorf("More messages than expected")
			}
			return nil
		}, func(err error) {
			if err == context.Canceled {
				return
			}

			require.NoError(t, err)
		})

		require.NoError(t, err)

		// wait for a period of time, but let's be reasonable
		select {
		case message := <-messagesCh:
			require.EqualValues(t, "hello world", string(message.Body))
		case <-processor.Done():
			t.Fatal("Processor was closed before messages arrived")
			break
		case <-ctx.Done():
			t.Fatal("Test finished before any messages arrived")
			break
		}
	})
}

func TestProcessorUnitTests(t *testing.T) {
	t.Parallel()

	t.Run("Processor", func(t *testing.T) {
		t.Run("StartAndClose", func(t *testing.T) {
			fakeNs := internal.NewFakeNamespace()
			processor, err := newProcessor(fakeNs, ProcessorWithQueue("hello"))
			require.NoError(t, err)

			subscribeCh := make(chan int)

			processor.subscribe = func(ctx context.Context, receiver internal.LegacyReceiver, shouldAutoComplete bool, handleMessage func(message *ReceivedMessage) error, notifyError func(err error)) bool {
				close(subscribeCh)
				return false
			}

			err = processor.Start(func(message *ReceivedMessage) error {
				return nil
			}, func(err error) {
			})

			require.NoError(t, err)

			<-subscribeCh

			require.NoError(t, processor.Close(context.Background()))
		})

		t.Run("CloseWaitsForActiveSubscribersToExit", func(t *testing.T) {
			fakeNs := internal.NewFakeNamespace()
			processor, err := newProcessor(fakeNs, ProcessorWithQueue("hello"))
			require.NoError(t, err)

			done := make(chan int)

			subscribeWg := &sync.WaitGroup{}
			subscribeWg.Add(1)

			processor.subscribe = func(ctx context.Context, receiver internal.LegacyReceiver, shouldAutoComplete bool, handleMessage func(message *ReceivedMessage) error, notifyError func(err error)) bool {
				subscribeWg.Done()
				// block subscribe from completing
				<-done
				return false
			}

			err = processor.Start(func(message *ReceivedMessage) error {
				return nil
			}, func(err error) {
			})

			require.NoError(t, err)

			// subscribe has been called (and is now blocked waiting on `done`!)
			subscribeWg.Wait()

			select {
			case <-time.After(time.Second * 2):
				close(done)
				require.NoError(t, processor.Close(context.Background()))
			case <-processor.Done():
				require.Fail(t, "processor is not done while receivers are active ")
			}

			<-processor.Done()
		})

		t.Run("CloseWithoutStart", func(t *testing.T) {
			fakeNs := internal.NewFakeNamespace()
			processor, err := newProcessor(fakeNs, ProcessorWithQueue("hello"))
			require.NoError(t, err)
			require.NoError(t, processor.Close(context.Background()))

			err = processor.Start(func(message *ReceivedMessage) error {
				t.Fail()
				return nil
			}, func(err error) {
				t.Fail()
			})

			require.EqualError(t, err, ErrProcessorClosed.Error())
		})

		t.Run("DoubleClose", func(t *testing.T) {
			fakeNs := internal.NewFakeNamespace()
			processor, err := newProcessor(fakeNs, ProcessorWithQueue("hello"))
			require.NoError(t, err)
			require.NoError(t, processor.Close(context.Background()))
			require.NoError(t, processor.Close(context.Background()))
		})
	})

	t.Run("subscribe", func(t *testing.T) {
		t.Run("cancelled by user does not retry", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel() // pre-cancel this context

			receiver := internal.NewFakeLegacyReceiver()
			var cancelledError error

			retry := subscribe(ctx, receiver, true, func(message *ReceivedMessage) error {
				return nil
			}, func(err error) {
				cancelledError = err
			})

			require.EqualError(t, cancelledError, context.Canceled.Error())
			require.False(t, retry, "User cancelling the context will not be retried")
			require.False(t, receiver.CloseCalled) // subscribe() is not responsible for the lifetime of the receiver
		})

		t.Run("error in the listener is retryable", func(t *testing.T) {
			receiver := internal.NewFakeLegacyReceiver()

			receiver.ListenImpl = func(ctx context.Context, handler internal.Handler) internal.ListenerHandle {
				ch := make(chan struct{})
				close(ch)
				return &internal.FakeListenerHandle{
					DoneChan: ch,
					ErrValue: errors.New("Some AMQP related error"),
				}
			}

			var errorFromListener error

			retry := subscribe(context.Background(), receiver, true, func(message *ReceivedMessage) error {
				return nil
			}, func(err error) {
				errorFromListener = err
			})

			require.EqualError(t, errorFromListener, "Some AMQP related error")
			require.True(t, retry, "AMQP errors will cause us to retry")
			require.False(t, receiver.CloseCalled) // subscribe() is not responsible for the lifetime of the receiver
		})

	})

	t.Run("handleSingleMessage", func(t *testing.T) {
		fakeMessage := &internal.Message{
			ID:        "fakeID",
			LockToken: &uuid.UUID{},
			SystemProperties: &internal.SystemProperties{
				SequenceNumber: to.Int64Ptr(1),
			},
		}

		setup := func() *internal.FakeInternalReceiver {
			receiver := internal.NewFakeLegacyReceiver()
			return receiver
		}

		t.Run("AutoCompleteCompleteMessage", func(t *testing.T) {
			receiver := setup()

			handleSingleMessage(func(message *ReceivedMessage) error {
				// successful return
				return nil
			}, func(err error) {
				require.NoError(t, err)
			}, true, receiver, fakeMessage)

			require.True(t, receiver.CompleteCalled)
			require.False(t, receiver.AbandonCalled)
		})

		t.Run("AutoCompleteAbandonMessage", func(t *testing.T) {
			receiver := setup()

			handleSingleMessage(func(message *ReceivedMessage) error {
				// error returned will abandon
				return errors.New("Purposefully reported error")
			}, func(err error) {
				require.EqualErrorf(t, err, "Purposefully reported error", "Error from the handler gets forwarded")
			}, true, receiver, fakeMessage)

			require.True(t, receiver.AbandonCalled)
			require.False(t, receiver.CompleteCalled)
		})

		t.Run("AutoCompleteAlreadySettledDoNotSettleTwice)", func(t *testing.T) {
			receiver := setup()

			handleSingleMessage(func(message *ReceivedMessage) error {
				// error returned will abandon
				return errors.New("Purposefully reported error")
			}, func(err error) {
				require.EqualErrorf(t, err, "Purposefully reported error", "Error from the handler gets forwarded")
			}, true, receiver, fakeMessage)

			// TODO: neither should be called - the message was already settled.
			require.True(t, receiver.AbandonCalled)
			require.False(t, receiver.CompleteCalled)
		})

		t.Run("autoComplete (off)", func(t *testing.T) {
			receiver := setup()

			handleSingleMessage(func(message *ReceivedMessage) error {
				// successful return
				return nil
			}, func(err error) {
				require.NoError(t, err)
			}, false, receiver, fakeMessage)

			require.False(t, receiver.CompleteCalled)
			require.False(t, receiver.AbandonCalled)
		})

		t.Run("SettlementErrorsAreForwarded(complete)", func(t *testing.T) {
			receiver := setup()

			receiver.CompleteMessageImpl = func(ctx context.Context, msg *internal.Message) error {
				return errors.New("Complete failed")
			}

			var settleError error

			handleSingleMessage(func(message *ReceivedMessage) error {
				return nil
			}, func(err error) {
				settleError = err
			}, true, receiver, fakeMessage)

			require.EqualError(t, settleError, "Complete failed")
		})

		t.Run("SettlementErrorsAreForwarded(abandon)", func(t *testing.T) {
			receiver := setup()

			receiver.AbandonMessageImpl = func(ctx context.Context, msg *internal.Message) error {
				return errors.New("Abandon failed")
			}

			var settleErrors []string

			handleSingleMessage(func(message *ReceivedMessage) error {
				return errors.New("Error that caused the abandon")
			}, func(err error) {
				settleErrors = append(settleErrors, err.Error())
			}, true, receiver, fakeMessage)

			require.EqualValues(t, settleErrors, []string{
				"Error that caused the abandon",
				"Abandon failed",
			})
		})
	})
}
