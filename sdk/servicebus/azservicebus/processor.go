package azservicebus

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/servicebus/azservicebus/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/servicebus/azservicebus/internal/utils"
)

var ErrReceiverClosed = errors.New("receiver has been closed and can no longer be used")

type processorConfig struct {
	fullEntityPath string
	retryPolicy    RetryPolicy
	receiveMode    ReceiveMode

	entity struct {
		Subqueue     SubQueue
		Queue        string
		Topic        string
		Subscription string
	}

	// determines if auto completion or abandonment of messages
	// happens based on the return value of user's processMessage handler.
	shouldAutoComplete bool
}

type Processor struct {
	mu *sync.Mutex

	receiversCtx    context.Context
	cancelReceivers func()

	processorCtx    context.Context
	cancelProcessor func()

	// configuration data that is read-only after the ServiceBusProcessor has been created
	config processorConfig

	activeReceiversWg *sync.WaitGroup
	createReceiver    func(ctx context.Context, maxConcurrency uint32) (internal.LegacyReceiver, error)
}

type ProcessorOption func(processor *Processor) error

// ProcessorWithSubQueue allows you to open the sub queue (ie: dead letter queues, transfer dead letter queues)
// for a queue or subscription.
func ProcessorWithSubQueue(subQueue SubQueue) ProcessorOption {
	return func(receiver *Processor) error {
		switch subQueue {
		case SubQueueDeadLetter:
		case SubQueueTransfer:
		case SubQueueNone:
			receiver.config.entity.Subqueue = subQueue
		default:
			return fmt.Errorf("unknown SubQueue %s", subQueue)
		}

		return nil
	}
}

// ProcessorWithReceiveMode controls the receive mode for the processor.
func ProcessorWithReceiveMode(receiveMode ReceiveMode) ProcessorOption {
	return func(processor *Processor) error {
		if receiveMode != ReceiveModePeekLock && receiveMode != ReceiveModeReceiveAndDelete {
			return fmt.Errorf("invalid receive mode specified %s", receiveMode)
		}

		processor.config.receiveMode = receiveMode
		return nil
	}
}

func ProcessorWithQueue(queue string) ProcessorOption {
	return func(processor *Processor) error {
		processor.config.entity.Queue = queue
		return nil
	}
}

func ProcessorWithSubscription(topic string, subscription string) ProcessorOption {
	return func(processor *Processor) error {
		processor.config.entity.Topic = topic
		processor.config.entity.Subscription = subscription
		return nil
	}
}

// ProcessorWithAutoComplete enables or disables auto-completion/abandon of messages
// When this option is enabled the result of the `processMessage` handler determines whether
// the message is abandoned (if an `error` is returned) or completed (if `nil` is returned).
// This option is enabled, by default.
func ProcessorWithAutoComplete(enableAutoCompleteMessages bool) ProcessorOption {
	return func(processor *Processor) error {
		processor.config.shouldAutoComplete = enableAutoCompleteMessages
		return nil
	}
}

func newProcessor(ns *internal.Namespace, options ...ProcessorOption) (*Processor, error) {
	processor := &Processor{
		config: processorConfig{
			receiveMode:        ReceiveModePeekLock,
			shouldAutoComplete: true,
		},

		mu:                &sync.Mutex{},
		activeReceiversWg: &sync.WaitGroup{},
	}

	for _, opt := range options {
		if err := opt(processor); err != nil {
			return nil, err
		}
	}

	entityPath, err := formatEntity(processor.config.entity)

	if err != nil {
		return nil, err
	}

	processor.config.fullEntityPath = entityPath

	receiveModeOption := internal.ReceiverWithReceiveMode(internal.PeekLockMode)

	if processor.config.receiveMode == "receiveAndDelete" {
		receiveModeOption = internal.ReceiverWithReceiveMode(internal.ReceiveAndDeleteMode)
	}

	processor.createReceiver = func(ctx context.Context, maxConcurrency uint32) (internal.LegacyReceiver, error) {
		// TODO: prefetch isn't _quite_ the right fit here (ie, we don't have a way to shut off the
		// spigot that I'm aware of).
		// But it's a pretty close approximation.
		return ns.NewReceiver(ctx, entityPath, receiveModeOption, internal.ReceiverWithPrefetchCount(maxConcurrency))
	}

	return processor, nil
}

type StartProcessorOptions struct {
	maxConcurrency uint32
}

type StartProcessorOption func(options *StartProcessorOptions) error

func StartWithConcurrency(maxConcurrency uint32) StartProcessorOption {
	return func(options *StartProcessorOptions) error {
		options.maxConcurrency = maxConcurrency
		return nil
	}
}

func newStartProcessorOptions(options ...StartProcessorOption) (*StartProcessorOptions, error) {
	subscribeOptions := &StartProcessorOptions{
		maxConcurrency: 1,
	}

	for _, opt := range options {
		if err := opt(subscribeOptions); err != nil {
			return nil, err
		}
	}

	return subscribeOptions, nil
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
// no user invention is required.
func (p *Processor) Start(handleMessage func(message *ReceivedMessage) error, handleError func(err error), options ...StartProcessorOption) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	startProcessorOptions, err := newStartProcessorOptions(options...)

	if err != nil {
		return err
	}

	p.processorCtx, p.cancelProcessor = context.WithCancel(context.Background())
	p.receiversCtx, p.cancelReceivers = context.WithCancel(context.Background())

	p.activeReceiversWg.Add(1)
	go func(ctx context.Context) {
		defer p.activeReceiversWg.Done()

		for {
			receiver, err := p.createReceiver(ctx, startProcessorOptions.maxConcurrency)

			if err != nil {
				// notify the user and then fall into doing a retry
				handleError(err)
			} else {
				if retry := subscribe(ctx, receiver, p.config.shouldAutoComplete, handleMessage, handleError); !retry {
					break
				}
			}

			if err := p.config.retryPolicy.Wait(ctx, 0); err != nil {
				handleError(err)
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

// Stop will wait for any pending callbacks to complete.
func (p *Processor) Close(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// processor is already stopped.
	if p.processorCtx == nil {
		return nil
	}

	p.cancelReceivers()

	err := utils.WaitForGroupOrContext(p.activeReceiversWg, ctx)

	// now unlock anyone _external_ to the processor that's waiting for us to exit or close.
	p.cancelProcessor()

	p.processorCtx = nil
	p.cancelProcessor = nil
	p.receiversCtx = nil
	p.cancelReceivers = nil

	return err
}

func subscribe(
	ctx context.Context,
	receiver internal.LegacyReceiver,
	shouldAutoComplete bool,
	handleMessage func(message *ReceivedMessage) error,
	notifyError func(err error)) bool {

	defer receiver.Close(ctx)

	activeCallbacksWg := &sync.WaitGroup{}
	notifyErrorAsync := wrapNotifyError(notifyError, activeCallbacksWg)

	const shouldRetry = true

	// TODO: `listen` doesn't give you additional messages unless you've settled previous ones.
	// In track 2 SDKs you should get a new message, regardless of settlement (ie: it's simply gated on the # of outstanding callbacks)
	listenHandle := receiver.Listen(ctx, internal.HandlerFunc(func(ctx context.Context, legacyMessage *internal.Message) error {
		activeCallbacksWg.Add(1)
		defer activeCallbacksWg.Done()

		// this shouldn't happen since we do a `select` above that prevents it.
		err := handleMessage(convertToReceivedMessage(legacyMessage))

		if err != nil {
			notifyErrorAsync(err)
		}

		var settleErr error

		if shouldAutoComplete {
			if err != nil {
				settleErr = receiver.AbandonMessage(ctx, legacyMessage)
			} else {
				settleErr = receiver.CompleteMessage(ctx, legacyMessage)
			}

			if settleErr != nil {
				notifyErrorAsync(settleErr)
			}
		}

		// errors from their handler are sent to their error handler but do not terminate the
		// subscription.
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

func formatEntity(entity struct {
	Subqueue     SubQueue
	Queue        string
	Topic        string
	Subscription string
}) (string, error) {
	entityPath := ""

	if entity.Queue != "" {
		entityPath = entity.Queue
	} else if entity.Topic != "" && entity.Subscription != "" {
		entityPath = fmt.Sprintf("%s/Subscriptions/%s", entity.Topic, entity.Subscription)
	} else {
		return "", errors.New("a queue or subscription was not specified")
	}

	if entity.Subqueue == SubQueueDeadLetter {
		entityPath += "/$DeadLetterQueue"
	} else if entity.Subqueue == SubQueueTransfer {
		entityPath += "/$Transfer/$DeadLetterQueue"
	}

	return entityPath, nil
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
