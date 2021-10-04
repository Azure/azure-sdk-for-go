// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/utils"
	"github.com/Azure/go-amqp"
	"github.com/devigned/tab"
)

type processorConfig struct {
	ReceiveMode ReceiveMode

	Entity entity

	// determines if auto completion or abandonment of messages
	// happens based on the return value of user's processMessage handler.
	ShouldAutoComplete bool
	MaxConcurrentCalls int

	baseRetrier internal.Retrier

	cleanupOnClose func()
}

// Processor is a push-based receiver for Service Bus.
type Processor struct {
	settler   settler
	amqpLinks internal.AMQPLinks

	mu *sync.Mutex

	userMessageHandler func(message *ReceivedMessage) error
	userErrorHandler   func(err error)

	receiversCtx    context.Context
	cancelReceivers func()

	wg sync.WaitGroup

	// configuration data that is read-only after the Processor has been created
	config processorConfig
}

// ProcessorOption represents an option on the Processor.
// Some examples:
// - `ProcessorWithReceiveMode` to configure the receive mode,
// - `ProcessorWithQueue` to target a queue.
type ProcessorOption func(processor *Processor) error

// ProcessorWithSubQueue allows you to open the sub queue (ie: dead letter queues, transfer dead letter queues)
// for a queue or subscription.
func ProcessorWithSubQueue(subQueue SubQueue) ProcessorOption {
	return func(p *Processor) error {
		return p.config.Entity.SetSubQueue(subQueue)
	}
}

// ProcessorWithReceiveMode controls the receive mode for the processor.
// The receive mode controls when a message is deleted from Service Bus.
//
// `azservicebus.PeekLock` is the default. The message is locked, preventing multiple
// receivers from processing the message at once. You control the lock state of the message
// using one of the message settlement functions like processor.CompleteMessage(), which removes
// it from Service Bus, or processor.AbandonMessage(), which makes it available again.
//
// `azservicebus.ReceiveAndDelete` causes Service Bus to remove the message as soon
// as it's received.
//
// More information about receive modes:
// https://docs.microsoft.com/azure/service-bus-messaging/message-transfers-locks-settlement#settling-receive-operations
func ProcessorWithReceiveMode(receiveMode ReceiveMode) ProcessorOption {
	return func(processor *Processor) error {
		if receiveMode != PeekLock && receiveMode != ReceiveAndDelete {
			return fmt.Errorf("invalid receive mode specified %d", receiveMode)
		}

		processor.config.ReceiveMode = receiveMode
		return nil
	}
}

// ProcessorWithAutoComplete enables or disables auto-completion/abandon of messages
// When this option is enabled the result of the `processMessage` handler determines whether
// the message is abandoned (if an `error` is returned) or completed (if `nil` is returned).
// This option is enabled, by default.
func ProcessorWithAutoComplete(enableAutoCompleteMessages bool) ProcessorOption {
	return func(processor *Processor) error {
		processor.config.ShouldAutoComplete = enableAutoCompleteMessages
		return nil
	}
}

// ProcessorWithMaxConcurrentCalls controls the maximum number of message processing
// goroutines that are active at any time.
// Default is 1.
func ProcessorWithMaxConcurrentCalls(maxConcurrentCalls int) ProcessorOption {
	return func(processor *Processor) error {
		processor.config.MaxConcurrentCalls = maxConcurrentCalls
		return nil
	}
}

func newProcessor(ns internal.NamespaceWithNewAMQPLinks, cleanupOnClose func(), options ...ProcessorOption) (*Processor, error) {
	processor := &Processor{
		config: processorConfig{
			ReceiveMode:        PeekLock,
			ShouldAutoComplete: true,
			MaxConcurrentCalls: 1,
			// TODO: make this configurable
			baseRetrier: internal.NewBackoffRetrier(internal.BackoffRetrierParams{
				Factor:     2,
				Min:        1,
				Max:        time.Minute,
				MaxRetries: 5,
			}),
			cleanupOnClose: cleanupOnClose,
		},

		mu: &sync.Mutex{},
	}

	for _, opt := range options {
		if err := opt(processor); err != nil {
			return nil, err
		}
	}

	entityPath, err := processor.config.Entity.String()

	if err != nil {
		return nil, err
	}

	processor.amqpLinks = ns.NewAMQPLinks(entityPath, func(ctx context.Context, session internal.AMQPSession) (internal.AMQPSenderCloser, internal.AMQPReceiverCloser, error) {
		linkOptions := createLinkOptions(processor.config.ReceiveMode, entityPath)
		_, receiver, err := createReceiverLink(ctx, session, linkOptions)

		if err != nil {
			return nil, nil, err
		}

		if err := receiver.IssueCredit(uint32(processor.config.MaxConcurrentCalls)); err != nil {
			_ = receiver.Close(ctx)
			return nil, nil, err
		}

		return nil, receiver, nil
	})

	processor.settler = newMessageSettler(processor.amqpLinks, processor.config.baseRetrier)
	processor.receiversCtx, processor.cancelReceivers = context.WithCancel(context.Background())

	return processor, nil
}

// Start will start receiving messages from the queue or subscription.
//
//   if err := processor.Start(context.TODO(), messageHandler, errorHandler); err != nil {
//     log.Fatalf("Processor failed to start: %s", err.Error())
//   }
//
// Any errors that occur (such as network disconnects, failures in handleMessage) will be
// sent to your handleError function. The processor will retry and restart as needed -
// no user intervention is required.
func (p *Processor) Start(ctx context.Context, handleMessage func(message *ReceivedMessage) error, handleError func(err error)) error {
	ctx, span := tab.StartSpan(ctx, internal.SpanProcessorLoop)
	defer span.End()

	err := func() error {
		p.mu.Lock()
		defer p.mu.Unlock()

		if p.userMessageHandler != nil {
			return errors.New("processor already started")
		}

		p.userMessageHandler = handleMessage
		p.userErrorHandler = handleError
		p.receiversCtx, p.cancelReceivers = context.WithCancel(ctx)

		return nil
	}()

	if err != nil {
		return err
	}

	for {
		if err := p.subscribe(p.receiversCtx); err != nil {

			if internal.IsCancelError(err) {
				break
			}

			p.userErrorHandler(err)

			if err := p.amqpLinks.RecoverIfNeeded(ctx, err); err != nil {
				p.userErrorHandler(err)
			}
		}
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

// Close will wait for any pending callbacks to complete.
// NOTE: Close() cannot be called synchronously in a message
// or error handler. You must run it asynchronously using
// `go processor.Close(ctx)` or similar.
func (p *Processor) Close(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.amqpLinks.ClosedPermanently() {
		return nil
	}

	ctx, span := tab.StartSpan(ctx, internal.SpanProcessorClose)
	defer span.End()

	defer func() {
		if err := p.amqpLinks.Close(ctx, true); err != nil {
			span.Logger().Debug(fmt.Sprintf("Error closing amqpLinks on processor.Close(): %s", err.Error()))
		}
	}()

	p.config.cleanupOnClose()

	_, receiver, _, _, err := p.amqpLinks.Get(ctx)

	if err != nil {
		span.Logger().Error(err)
		return err
	}

	if err := receiver.DrainCredit(ctx); err != nil {
		span.Logger().Error(err)
		// fall through for now and just let whatever is going on finish
		// otherwise they might not be able to actually close.
	}

	p.cancelReceivers()
	return utils.WaitForGroupOrContext(ctx, &p.wg)
}

// CompleteMessage completes a message, deleting it from the queue or subscription.
func (p *Processor) CompleteMessage(ctx context.Context, message *ReceivedMessage) error {
	return p.settler.CompleteMessage(ctx, message)
}

// AbandonMessage will cause a message to be returned to the queue or subscription.
// This will increment its delivery count, and potentially cause it to be dead lettered
// depending on your queue or subscription's configuration.
func (p *Processor) AbandonMessage(ctx context.Context, message *ReceivedMessage) error {
	return p.settler.AbandonMessage(ctx, message)
}

// DeferMessage will cause a message to be deferred. Deferred messages
// can be received using `Receiver.ReceiveDeferredMessages`.
func (p *Processor) DeferMessage(ctx context.Context, message *ReceivedMessage) error {
	return p.settler.DeferMessage(ctx, message)
}

// DeadLetterMessage settles a message by moving it to the dead letter queue for a
// queue or subscription. To receive these messages create a receiver with `Client.NewProcessor()`
// using the `ProcessorWithSubQueue()` option.
func (p *Processor) DeadLetterMessage(ctx context.Context, message *ReceivedMessage, options ...DeadLetterOption) error {
	return p.settler.DeadLetterMessage(ctx, message, options...)
}

// subscribe continually receives messages from Service Bus, stopping
// if a fatal link/connection error occurs.
func (p *Processor) subscribe(ctx context.Context) error {
	_, receiver, _, _, err := p.amqpLinks.Get(ctx)

	if err != nil {
		return err
	}

	p.wg.Add(1)
	defer p.wg.Done()

	for {
		amqpMessage, err := receiver.Receive(ctx)

		if err != nil {
			return err
		}

		if amqpMessage == nil {
			// amqpMessage shouldn't be nil here, but somehow it is.
			// need to track this down in the AMQP library.
			continue
		}

		p.wg.Add(1)

		go func() {
			defer p.wg.Done()

			// purposefully avoiding using `ctx`. We always let processing complete
			// for message threads to avoid potential message loss.
			_ = p.processMessage(context.Background(), receiver, amqpMessage)
		}()
	}
}

func (p *Processor) processMessage(ctx context.Context, receiver internal.AMQPReceiver, amqpMessage *amqp.Message) error {
	ctx, span := tab.StartSpan(ctx, internal.SpanProcessorMessage)
	defer span.End()

	receivedMessage := newReceivedMessage(ctx, amqpMessage)
	messageHandlerErr := p.userMessageHandler(receivedMessage)

	if messageHandlerErr != nil {
		p.userErrorHandler(messageHandlerErr)
	}

	if p.config.ShouldAutoComplete {
		var settleErr error

		if messageHandlerErr != nil {
			settleErr = p.settler.AbandonMessage(ctx, receivedMessage)
		} else {
			settleErr = p.settler.CompleteMessage(ctx, receivedMessage)
		}

		if settleErr != nil {
			p.userErrorHandler(fmt.Errorf("failed to settle message with ID '%s': %w", receivedMessage.ID, settleErr))
			return settleErr
		}
	}

	select {
	case <-p.receiversCtx.Done():
		return nil
	default:
	}

	if err := receiver.IssueCredit(1); err != nil {
		if !internal.IsDrainingError(err) {
			p.userErrorHandler(err)
			return fmt.Errorf("failed issuing additional credit, processor will be restarted: %w", err)
		}
	}

	return nil
}

// processorWithQueue configures a processor to connect to a queue.
func processorWithQueue(queue string) ProcessorOption {
	return func(processor *Processor) error {
		processor.config.Entity.Queue = queue
		return nil
	}
}

// processorWithSubscription configures a processor to connect to a subscription
// associated with a topic.
func processorWithSubscription(topic string, subscription string) ProcessorOption {
	return func(processor *Processor) error {
		processor.config.Entity.Topic = topic
		processor.config.Entity.Subscription = subscription
		return nil
	}
}
