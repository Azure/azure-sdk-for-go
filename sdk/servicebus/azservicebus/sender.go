package azservicebus

import (
	"context"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/servicebus/azservicebus/internal"
	"github.com/Azure/go-amqp"
)

// SenderOption specifies an option that can configure a Sender.
type SenderOption func(sender *Sender) error

// Sender is used to send messages as well as schedule them to be delivered at a later date.
type Sender struct {
	config struct {
		queueOrTopic string
	}

	mu           *sync.Mutex
	legacySender internal.LegacySender

	linkState *linkState

	// for testing
	nsCreateLegacySender func(ctx context.Context, entityPath string) (internal.LegacySender, error)
}

type SendableMessage interface {
	toAMQPMessage() (*amqp.Message, error)
}

// CreateMessageBatch can be used to create a batch that contain multiple
// messages. Sending a batch of messages is more efficient than sending the
// messages one at a time.
func (s *Sender) CreateMessageBatch(ctx context.Context) (*MessageBatch, error) {
	if s.linkState.Closed() {
		return nil, s.linkState.Err()
	}

	legacySender, err := s.createAmqpSender(ctx)

	if err != nil {
		return nil, err
	}

	return newMessageBatch(int(legacySender.MaxMessageSize())), nil
}

// SendMessage sends a message to a queue or topic.
// Message can be a MessageBatch (created using `Sender.CreateMessageBatch`) or
// a Message.
func (s *Sender) SendMessage(ctx context.Context, message SendableMessage) error {
	if s.linkState.Closed() {
		return s.linkState.Err()
	}

	legacySender, err := s.createAmqpSender(ctx)

	if err != nil {
		return err
	}

	amqpMessage, err := message.toAMQPMessage()

	if err != nil {
		return err
	}

	return legacySender.SendAMQPMessage(ctx, amqpMessage)
}

// Close permanently closes the Sender.
func (s *Sender) Close(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	defer s.linkState.Close()

	var err error

	if s.legacySender != nil {
		err = s.legacySender.Close(ctx)
		s.legacySender = nil
	}

	return err
}

func (sender *Sender) createAmqpSender(ctx context.Context) (internal.LegacySender, error) {
	sender.mu.Lock()
	defer sender.mu.Unlock()

	if sender.legacySender != nil {
		return sender.legacySender, nil
	}

	// TODO: allow passing in relevant options if needed
	legacySender, err := sender.nsCreateLegacySender(ctx, sender.config.queueOrTopic)

	if err != nil {
		return nil, err
	}

	sender.legacySender = legacySender
	return legacySender, nil
}

// ie: `*internal.Namespace`
type legacySenderNamespace interface {
	NewLegacySender(ctx context.Context, entityPath string) (internal.LegacySender, error)
}

func newSender(ns legacySenderNamespace, queueOrTopic string) (*Sender, error) {
	sender := &Sender{
		config: struct {
			queueOrTopic string
		}{
			queueOrTopic: queueOrTopic,
		},
		mu:        &sync.Mutex{},
		linkState: newLinkState(context.Background(), errClosed{link: "sender"}),
	}

	sender.nsCreateLegacySender = ns.NewLegacySender
	return sender, nil
}
