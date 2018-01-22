package servicebus

import (
	"context"
	"pack.ag/amqp"
)

// Sender provides session and link handling for an sending entity path
type Sender struct {
	client     *amqp.Client
	session    *amqp.Session
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
	err := s.session.Close()
	if err != nil {
		return err
	}

	err = s.sender.Close()
	if err != nil {
		return err
	}
	return nil
}

// Send will send a message using the session and link
func (s *Sender) Send(ctx context.Context, msg *amqp.Message) error {
	// TODO: Add in recovery logic in case the link / session has gone down
	err := s.sender.Send(ctx, msg)
	if err != nil {
		return err
	}
	return nil
}

// newSessionAndLink will replace the existing session and link
func (s *Sender) newSessionAndLink() error {
	session, err := s.client.NewSession()
	if err != nil {
		return err
	}

	amqpSender, err := session.NewSender(amqp.LinkAddress(s.entityPath))
	if err != nil {
		return err
	}

	s.session = session
	s.sender = amqpSender
	return nil
}
