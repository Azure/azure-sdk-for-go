package azservicebus

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqpCommon "github.com/Azure/azure-amqp-common-go/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/servicebus/azservicebus/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/servicebus/azservicebus/internal/utils"
)

type processorConfig struct {
	FullEntityPath string
	ReceiveMode    ReceiveMode

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
	// configuration data that is read-only after the Processor has been created
	config processorConfig

	mu *sync.Mutex

	receiversCtx    context.Context
	cancelReceivers func()

	processorCtx    context.Context
	cancelProcessor func()

	activeReceiversWg *sync.WaitGroup

	// replaceable for unit tests
	ns        legacyNamespace
	subscribe func(
		ctx context.Context,
		receiver internal.LegacyReceiver,
		shouldAutoComplete bool,
		handleMessage func(message *ReceivedMessage) error,
		notifyError func(err error)) bool
}

type ProcessorOption func(processor *Processor) error

// ProcessorWithSubQueue allows you to open the sub queue (ie: dead letter queues, transfer dead letter queues)
// for a queue or subscription.
func ProcessorWithSubQueue(subQueue SubQueue) ProcessorOption {
	return func(receiver *Processor) error {
		switch subQueue {
		case SubQueueDeadLetter:
		case SubQueueTransfer:
		case "":
			receiver.config.Entity.Subqueue = subQueue
		default:
			return fmt.Errorf("unknown SubQueue %s", subQueue)
		}

		return nil
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

func ProcessorWithQueue(queue string) ProcessorOption {
	return func(processor *Processor) error {
		processor.config.Entity.Queue = queue
		return nil
	}
}

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

type legacyNamespace interface {
	NewReceiver(ctx context.Context, entityPath string, opts ...internal.ReceiverOption) (internal.LegacyReceiver, error)
}

func newProcessor(ns legacyNamespace, options ...ProcessorOption) (*Processor, error) {
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
				Delay: time.Second * 5,
			},
		},

		mu:                &sync.Mutex{},
		activeReceiversWg: &sync.WaitGroup{},

		ns:        ns,
		subscribe: subscribe,
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

	processor.config.FullEntityPath = entityPath

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
	select {
	case <-p.Done():
		return ErrClosed{link: "processor"}
	default:
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.activeReceiversWg.Add(1)
	go func(ctx context.Context) {
		defer p.activeReceiversWg.Done()

		for {
			retry, _ := amqpCommon.Retry(p.config.RetryOptions.Times, p.config.RetryOptions.Delay, func() (interface{}, error) {
				receiver, err := p.ns.NewReceiver(ctx, p.config.FullEntityPath,
					internal.ReceiverWithReceiveMode(internal.ReceiveMode(p.config.ReceiveMode)))

				if err != nil {
					// notify the user and then fall into doing a retry
					handleError(err)
					return true, amqpCommon.Retryable("")
				}

				defer receiver.Close(ctx)

				if err := receiver.IssueCredit(uint32(p.config.MaxConcurrentCalls)); err != nil {
					// notify the user but there's no reason to restart because this failure must be
					// an internal error.
					handleError(err)
					return false, nil
				}

				// we retry infinitely, but do it in the pattern they specify via their retryOptions for each "round" of retries.
				retry := p.subscribe(ctx, receiver, p.config.ShouldAutoComplete, handleMessage, handleError)

				if retry {
					return true, amqpCommon.Retryable("")
				} else {
					return false, nil
				}
			})

			if !retry.(bool) {
				break
			}
		}
	}(p.receiversCtx)

	return nil
}

// Done returns a channel that will be close()'d when the Processor
// has been closed.
func (p *Processor) Done() <-chan struct{} {
	return p.processorCtx.Done()
}

// Close will wait for any pending callbacks to complete.
func (p *Processor) Close(ctx context.Context) error {
	select {
	case <-p.Done():
		return nil
	default:
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.cancelReceivers()

	err := utils.WaitForGroupOrContext(p.activeReceiversWg, ctx)

	// now unlock anyone _external_ to the processor that's waiting for us to exit or close.
	p.cancelProcessor()

	return err
}

func subscribe(
	ctx context.Context,
	receiver internal.LegacyReceiver,
	shouldAutoComplete bool,
	handleMessage func(message *ReceivedMessage) error,
	notifyError func(err error)) bool {

	activeCallbacksWg := &sync.WaitGroup{}
	notifyErrorAsync := wrapNotifyError(notifyError, activeCallbacksWg)

	const shouldRetry = true

	listenHandle := receiver.Listen(ctx, internal.HandlerFunc(func(ctx context.Context, legacyMessage *internal.Message) error {
		// this shouldn't happen since we do a `select` above that prevents it.
		// errors from their handler are sent to their error handler but do not terminate the
		// subscription.
		activeCallbacksWg.Add(1)
		defer activeCallbacksWg.Done()

		handleSingleMessage(handleMessage, notifyErrorAsync, shouldAutoComplete, receiver, legacyMessage)

		// user callback completes and they get a new credit
		if err := receiver.IssueCredit(1); err != nil {
			notifyErrorAsync(err)
		}

		return nil
	}))

	select {
	case <-ctx.Done():
		notifyErrorAsync(ctx.Err())
		activeCallbacksWg.Wait()
		return !shouldRetry
	case <-listenHandle.Done(): // TODO: eliminate this redundancy soon, and hopefully we only need to rely on ctx.Done() above.
		notifyErrorAsync(listenHandle.Err())
		activeCallbacksWg.Wait()

		// we should retry since the listen handle can be closed if we did a .Recover() on the receiver.
		return shouldRetry
	}
}

func handleSingleMessage(handleMessage func(message *ReceivedMessage) error, notifyErrorAsync func(err error), shouldAutoComplete bool, receiver internal.LegacyReceiver, legacyMessage *internal.Message) {
	err := handleMessage(convertToReceivedMessage(legacyMessage))

	if err != nil {
		notifyErrorAsync(err)
	}

	var settleErr error

	if shouldAutoComplete {
		// NOTE: we ignore the passed in context. Since we're settling behind the scenes
		// it's nice to wrap it up so users don't have to track it.
		if err != nil {
			settleErr = receiver.AbandonMessage(context.Background(), legacyMessage)
		} else {
			settleErr = receiver.CompleteMessage(context.Background(), legacyMessage)
		}

		if settleErr != nil {
			notifyErrorAsync(settleErr)
		}
	}
}

//
// settlement methods
// TODO: in other processor implementations this is implemented in the argument for the processMessage
// callback. You need some sort of association or else you have to track message <-> receiver mappings.
//

func (p *Processor) CompleteMessage(ctx context.Context, message *ReceivedMessage) error {
	return message.legacyMessage.Complete(ctx)
}

func (p *Processor) DeadLetterMessage(ctx context.Context, message *ReceivedMessage) error {
	// TODO: expand to let them set the reason and description.
	return message.legacyMessage.DeadLetter(ctx, nil)
}

func (p *Processor) AbandonMessage(ctx context.Context, message *ReceivedMessage) error {
	return message.legacyMessage.Abandon(ctx)
}

func (p *Processor) DeferMessage(ctx context.Context, message *ReceivedMessage) error {
	return message.legacyMessage.Defer(ctx)
}

func convertToReceivedMessage(legacyMessage *internal.Message) *ReceivedMessage {
	return &ReceivedMessage{
		// TODO: When we swap out the encoding from the legacy we should also make it so LockToken is simply a string, not expected to be a UUID.
		// Ie, it should be opaque to us.
		LockToken:      legacyMessage.LockToken.String(),
		SequenceNumber: *legacyMessage.SystemProperties.SequenceNumber,
		Message: Message{
			Body: legacyMessage.Data,
			ID:   legacyMessage.ID,
		},
		legacyMessage: legacyMessage,
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
