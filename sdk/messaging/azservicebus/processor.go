// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqpCommon "github.com/Azure/azure-amqp-common-go/v3"
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

	RetryOptions struct {
		Times int
		Delay time.Duration
	}
}

// Processor is a push-based receiver for Service Bus.
type Processor struct {
	settler *messageSettler

	mu        *sync.Mutex
	amqpLinks internal.AMQPLinks

	receiversCtx    context.Context
	cancelReceivers func()

	processorCtx    context.Context
	cancelProcessor func()

	wg *sync.WaitGroup

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
func ProcessorWithReceiveMode(receiveMode ReceiveMode) ProcessorOption {
	return func(processor *Processor) error {
		if receiveMode != PeekLock && receiveMode != ReceiveAndDelete {
			return fmt.Errorf("invalid receive mode specified %d", receiveMode)
		}

		processor.config.ReceiveMode = receiveMode
		return nil
	}
}

// ProcessorWithQueue configures a processor to connect to a queue.
func ProcessorWithQueue(queue string) ProcessorOption {
	return func(processor *Processor) error {
		processor.config.Entity.Queue = queue
		return nil
	}
}

// ProcessorWithSubscription configures a processor to connect to a subscription
// associated with a topic.
func ProcessorWithSubscription(topic string, subscription string) ProcessorOption {
	return func(processor *Processor) error {
		processor.config.Entity.Topic = topic
		processor.config.Entity.Subscription = subscription
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

func newProcessor(ns internal.NamespaceWithNewAMQPLinks, options ...ProcessorOption) (*Processor, error) {
	processor := &Processor{
		config: processorConfig{
			ReceiveMode:        PeekLock,
			ShouldAutoComplete: true,
			MaxConcurrentCalls: 1,
			RetryOptions: struct {
				Times int
				Delay time.Duration
			}{
				// TODO: allow these to be configured.
				Times: 10,
				Delay: 5 * time.Second,
			},
		},

		mu: &sync.Mutex{},
		wg: &sync.WaitGroup{},
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
		return createReceiverLink(ctx, session, linkOptions)
	})
	processor.settler = &messageSettler{links: processor.amqpLinks}

	processor.processorCtx, processor.cancelProcessor = context.WithCancel(context.Background())
	processor.receiversCtx, processor.cancelReceivers = context.WithCancel(context.Background())

	return processor, nil
}

// Start will start receiving messages from the queue or subscription.
//
//   if err := processor.Start(messageHandler, errorHandler); err != nil {
//     log.Fatalf("Processor failed to start: %s", err.Error())
//   }
//
//   <- processor.Done()
//
// Any errors that occur (such as network disconnects, failures in handleMessage) will be
// sent to your handleError function. The processor will retry and restart as needed -
// no user intervention is required.
func (p *Processor) Start(handleMessage func(message *ReceivedMessage) error, handleError func(err error)) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// TODO: stop multiple invocations!
	/*
		p.handleMessage = ()
		p.handleError = func() {
			tab.For(ctx).Debug()
		}
	*/

	p.wg.Add(1)

	go func(ctx context.Context) {
		defer func() {
			p.wg.Done()
			tab.For(ctx).Debug("Exiting forever loop in Processor")
		}()

		for {
			retry, _ := amqpCommon.Retry(p.config.RetryOptions.Times, p.config.RetryOptions.Delay, func() (interface{}, error) {
				_, receiver, _, _, err := p.amqpLinks.Get(ctx)

				if err != nil {
					if !isRetryableLinksError(err) {
						tab.For(ctx).Debug(fmt.Sprintf("Non retriable error occurred ('%s'), stopping the processor loop", err.Error()))
						handleError(fmt.Errorf("Non-retryable error, stopping processor loop: %w", err))
						return false, nil
					}

					// notify the user and then fall into doing a retry
					handleError(fmt.Errorf("failed getting links for subscribe, will retry: %w", err))
					return true, amqpCommon.Retryable("")
				}

				if err := receiver.IssueCredit(uint32(p.config.MaxConcurrentCalls)); err != nil {
					// notify the user but there's no reason to restart because this failure must be
					// an internal error.
					handleError(fmt.Errorf("internal failure when issuing credit, will recreate links: %w", err))
					return true, nil
				}

				// we retry infinitely, but do it in the pattern they specify via their retryOptions for each "round" of retries.
				err = p.subscribe(ctx, receiver, handleMessage, handleError)

				if err != nil {
					if isRetryableSubscribeError(err) {
						if err := p.amqpLinks.Close(ctx, false); err != nil {
							tab.For(ctx).Debug(fmt.Sprintf("failed when closing links, will reopen: %s", err.Error()))
						}

						return true, amqpCommon.Retryable("")
					}

					return false, nil
				}

				return true, amqpCommon.Retryable("")
			})

			if !retry.(bool) {
				break
			}
		}
	}(p.receiversCtx)

	return nil
}

// Close will wait for any pending callbacks to complete.
func (p *Processor) Close(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.cancelReceivers()
	return utils.WaitForGroupOrContext(ctx, p.wg)
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

func (p *Processor) subscribe(
	ctx context.Context,
	receiver internal.AMQPReceiver,
	handleMessage func(message *ReceivedMessage) error,
	notifyError func(err error)) error {

	notifyErrorAsync := wrapNotifyError(notifyError, p.wg)

	p.wg.Add(1)
	defer p.wg.Done()

	for {
		amqpMessage, err := receiver.Receive(ctx)

		if err != nil {
			if !isRetryableSubscribeError(err) {
				return err
			}

			// TODO: need to get a backoff policy in here.
			continue
		}

		p.wg.Add(1)

		go func() {
			defer p.wg.Done()

			receivedMessage := newReceivedMessage(ctx, amqpMessage)
			err := handleMessage(receivedMessage)

			if err != nil {
				notifyErrorAsync(err)
			}

			var settleErr error

			if p.config.ShouldAutoComplete {
				// NOTE: we ignore the passed in context. Since we're settling behind the scenes
				// it's nice to wrap it up so users don't have to track it.
				if err != nil {
					settleErr = p.settler.AbandonMessage(context.Background(), receivedMessage)
				} else {
					settleErr = p.settler.CompleteMessage(context.Background(), receivedMessage)
				}

				if settleErr != nil {
					notifyErrorAsync(fmt.Errorf("failed settling message with ID '%s': %w", receivedMessage.ID, settleErr))
				}
			}

			if err := receiver.IssueCredit(1); err != nil {
				notifyErrorAsync(fmt.Errorf("failed issuing additional credit, processor will be restarted: %w", err))

				// close the links here and cause the processor for this receiver to restart
				if err := p.amqpLinks.Close(ctx, false); err != nil {
					tab.For(ctx).Debug(fmt.Sprintf("Failed when closing links, but was restarting anyways: %s", err.Error()))
				}
			}

			notifyErrorAsync(err)
		}()
	}
}

func wrapNotifyError(fn func(err error), wg *sync.WaitGroup) func(err error) {
	return func(err error) {
		if err == nil {
			return
		}

		wg.Add(1)
		go func() {
			fn(err)
			wg.Done()
		}()
	}
}

func isRetryableSubscribeError(err error) bool {
	isFatalLinkError := err == amqp.ErrConnClosed ||
		err == amqp.ErrLinkClosed ||
		err == amqp.ErrLinkDetached ||
		err == amqp.ErrSessionClosed

	if isFatalLinkError || isCancelled(err) {
		return false
	}

	return true
}

// isRetryableLinksError checks if an error is retryable if it
// was returned from links.Get().
// NOTE: this function panics if you pass it a nil error.
func isRetryableLinksError(err error) bool {
	if isCancelled(err) {
		return false
	}

	switch err.(type) {
	case interface{ NonRetriable() }:
		// we're done, they've closed the links permanently (ie, we're not meant to
		// recover from this)
		return false
	default:
		return true
	}
}
