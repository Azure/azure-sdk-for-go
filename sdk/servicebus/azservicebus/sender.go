package azservicebus

import (
	"context"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/servicebus/azservicebus/internal"
)

// SenderOption specifies an option that can configure a Sender.
type SenderOption func(sender *Sender) error

// Sender is used to send messages as well as schedule them to be delivered at a later date.
type Sender struct {
	config struct {
		queueOrTopic string
	}

	createSender func(ctx context.Context) (legacySender, error)

	mu           *sync.Mutex
	legacySender legacySender

	linkState *linkState
}

// SendMessage sends a message to a queue or topic.
func (s *Sender) SendMessage(ctx context.Context, message *Message) error {
	if s.linkState.Closed() {
		return s.linkState.Err()
	}

	legacySender, err := s.createSender(ctx)

	if err != nil {
		return err
	}

	return legacySender.Send(ctx, &internal.Message{
		ID:   message.ID,
		Data: message.Body,
		// TODO: more fields
	})
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

// ie: `*internal.Namespace`
type legacySenderNamespace interface {
	NewSender(ctx context.Context, entityPath string, opts ...internal.SenderOption) (*internal.Sender, error)
}

// ie: `*internal.Sender`
type legacySender interface {
	Send(ctx context.Context, msg *internal.Message, opts ...internal.SendOption) error
	Close(ctx context.Context) error
}

func newSender(ns legacySenderNamespace, queueOrTopic string, options ...SenderOption) (*Sender, error) {
	sender := &Sender{
		config: struct {
			queueOrTopic string
		}{
			queueOrTopic: queueOrTopic,
		},
		mu:        &sync.Mutex{},
		linkState: newLinkState(context.Background(), errClosed{link: "sender"}),
	}

	for _, opt := range options {
		if err := opt(sender); err != nil {
			return nil, err
		}
	}

	sender.createSender = func(ctx context.Context) (legacySender, error) {
		sender.mu.Lock()
		defer sender.mu.Unlock()

		if sender.legacySender != nil {
			return sender.legacySender, nil
		}

		// TODO: allow passing in relevant options if needed
		legacySender, err := ns.NewSender(ctx, queueOrTopic)

		if err != nil {
			return nil, err
		}

		sender.legacySender = legacySender
		return legacySender, nil
	}

	return sender, nil
}
