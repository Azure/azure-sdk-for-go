package servicebus

import (
	"context"
	"pack.ag/amqp"
)

// Sender provides session and link handling for an sending entity path
type Sender struct {
	client     *amqp.Client
	session    *Session
	sender     *amqp.Sender
	entityPath string
	Name       string
}

// NewSender creates a new Service Bus message sender given an AMQP client and entity path
func NewSender(client *amqp.Client, entityPath string) (*Sender, error) {
	sender := &Sender{
		client:     client,
		entityPath: entityPath,
	}

	err := sender.newSessionAndLink()
	if err != nil {
		return nil, err
	}

	return sender, nil
}

// Recover will attempt to close the current session and link, then rebuild them
func (s *Sender) Recover() error {
	err := s.Close()
	if err != nil {
		return err
	}

	err = s.newSessionAndLink()
	if err != nil {
		return err
	}

	return nil
}

// Close will close the AMQP session and link of the sender
func (s *Sender) Close() error {
	err := s.sender.Close()
	if err != nil {
		return err
	}

	err = s.session.Close()
	if err != nil {
		return err
	}
	return nil
}

// Send will send a message to the entity path with options
func (s *Sender) Send(ctx context.Context, msg *amqp.Message, opts ...SendOption) error {
	// TODO: Add in recovery logic in case the link / session has gone down
	s.prepareMessage(msg)
	for _, opt := range opts {
		opt(msg)
	}

	err := s.sender.Send(ctx, msg)
	if err != nil {
		return err
	}
	return nil
}

func (s *Sender) prepareMessage(msg *amqp.Message) {
	if msg.Properties == nil {
		msg.Properties = &amqp.MessageProperties{}
	}

	if msg.Properties.GroupID == "" {
		msg.Properties.GroupID = s.session.SessionID
	}
}

// newSessionAndLink will replace the existing session and link
func (s *Sender) newSessionAndLink() error {
	amqpSession, err := s.client.NewSession()
	if err != nil {
		return err
	}

	amqpSender, err := amqpSession.NewSender(amqp.LinkAddress(s.entityPath))
	if err != nil {
		return err
	}

	s.session = NewSession(amqpSession)
	s.sender = amqpSender
	return nil
}

// SendOption provides a way to customize a message on sending
type SendOption func(message *amqp.Message) error

// SendWithMessageID provides an option of adding a message ID for the sent message
func SendWithMessageID(msgID interface{}) SendOption {
	return func(msg *amqp.Message) error {
		msg.Properties.MessageID = msgID
		return nil
	}
}

// SendWithoutSessionID will set the SessionID to nil. If sending to a partitioned Service Bus queue, this will cause
// the queue distributed the message in a round robin fashion to the next available partition with the effect of not
// enforcing FIFO ordering of messages, but enabling more efficient distribution of messages across partitions.
func SendWithoutSessionID() SendOption {
	return func(msg *amqp.Message) error {
		msg.Properties.GroupID = ""
		return nil
	}
}
