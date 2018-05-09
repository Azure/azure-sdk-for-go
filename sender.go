package servicebus

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-amqp-common-go"
	"github.com/Azure/azure-amqp-common-go/log"
	"github.com/Azure/azure-amqp-common-go/uuid"
	"github.com/opentracing/opentracing-go"
	"pack.ag/amqp"
)

// sender provides session and link handling for an sending entity path
type (
	sender struct {
		namespace  *Namespace
		connection *amqp.Client
		session    *session
		sender     *amqp.Sender
		entityPath string
		Name       string
	}

	// SendOption provides a way to customize a message on sending
	SendOption func(event *Event) error

	eventer interface {
		Set(key, value string)
		toMsg() *amqp.Message
	}
)

// newSender creates a new Service Bus message sender given an AMQP client and entity path
func (ns *Namespace) newSender(ctx context.Context, entityPath string) (*sender, error) {
	span, ctx := ns.startSpanFromContext(ctx, "sb.sender.newSender")
	defer span.Finish()

	s := &sender{
		namespace:  ns,
		entityPath: entityPath,
	}
	log.For(ctx).Debug(fmt.Sprintf("creating a new sender for entity path %s", s.entityPath))
	err := s.newSessionAndLink(ctx)
	return s, err
}

// Recover will attempt to close the current session and link, then rebuild them
func (s *sender) Recover(ctx context.Context) error {
	span, ctx := s.startProducerSpanFromContext(ctx, "sb.sender.Recover")
	defer span.Finish()
	_ = s.Close(ctx) // we expect the sender is in an error state
	return s.newSessionAndLink(ctx)
}

// Close will close the AMQP connection, session and link of the sender
func (s *sender) Close(ctx context.Context) error {
	span, _ := s.startProducerSpanFromContext(ctx, "sb.sender.Close")
	defer span.Finish()

	return s.connection.Close()
}

// Send will send a message to the entity path with options
//
// This will retry sending the message if the server responds with a busy error.
func (s *sender) Send(ctx context.Context, event *Event, opts ...SendOption) error {
	span, ctx := s.startProducerSpanFromContext(ctx, "sb.sender.Send")
	defer span.Finish()

	if event.GroupID == nil {
		event.GroupID = &s.session.SessionID
		next := s.session.getNext()
		event.GroupSequence = &next
	}

	if event.ID == "" {
		id, err := uuid.NewV4()
		if err != nil {
			return err
		}
		event.ID = id.String()
	}

	for _, opt := range opts {
		err := opt(event)
		if err != nil {
			return err
		}
	}

	return s.trySend(ctx, event)
}

func (s *sender) trySend(ctx context.Context, evt eventer) error {
	sp, ctx := s.startProducerSpanFromContext(ctx, "sb.sender.trySend")
	defer sp.Finish()

	times := 3
	delay := 10 * time.Second
	durationOfSend := 3 * time.Second
	if deadline, ok := ctx.Deadline(); ok {
		times = int(time.Until(deadline) / (delay + durationOfSend))
		times = max(times, 1) // give at least one chance at sending
	}
	_, err := common.Retry(times, delay, func() (interface{}, error) {
		sp, ctx := s.startProducerSpanFromContext(ctx, "sb.sender.trySend.transmit")
		defer sp.Finish()

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			innerCtx, cancel := context.WithTimeout(ctx, durationOfSend)
			defer cancel()

			err := opentracing.GlobalTracer().Inject(sp.Context(), opentracing.TextMap, evt)
			if err != nil {
				log.For(ctx).Error(err)
				return nil, err
			}

			msg := evt.toMsg()
			sp.SetTag("sb.message-id", msg.Properties.MessageID)
			err = s.sender.Send(innerCtx, msg)
			if err != nil {
				recoverErr := s.Recover(ctx)
				if recoverErr != nil {
					log.For(ctx).Error(recoverErr)
				}
			}

			if amqpErr, ok := err.(*amqp.Error); ok {
				if amqpErr.Condition == "com.microsoft:server-busy" {
					return nil, common.Retryable(amqpErr.Condition)
				}
			}

			return nil, err
		}
	})
	return err
}

func (s *sender) String() string {
	return s.Name
}

func (s *sender) getAddress() string {
	return s.entityPath
}

func (s *sender) getFullIdentifier() string {
	return s.namespace.getEntityAudience(s.getAddress())
}

// newSessionAndLink will replace the existing session and link
func (s *sender) newSessionAndLink(ctx context.Context) error {
	span, ctx := s.startProducerSpanFromContext(ctx, "sb.sender.newSessionAndLink")
	defer span.Finish()

	connection, err := s.namespace.newConnection()
	if err != nil {
		log.For(ctx).Error(err)
		return err
	}
	s.connection = connection

	err = s.namespace.negotiateClaim(ctx, connection, s.getAddress())
	if err != nil {
		log.For(ctx).Error(err)
		return err
	}

	amqpSession, err := connection.NewSession()
	if err != nil {
		log.For(ctx).Error(err)
		return err
	}

	amqpSender, err := amqpSession.NewSender(
		amqp.LinkTargetAddress(s.getAddress()),
		amqp.LinkReceiverSettle(amqp.ModeSecond))
	if err != nil {
		log.For(ctx).Error(err)
		return err
	}

	s.session, err = newSession(amqpSession)
	if err != nil {
		log.For(ctx).Error(err)
		return err
	}

	s.sender = amqpSender
	return nil
}

// SendWithMessageID configures the message with a message ID
func SendWithMessageID(messageID string) SendOption {
	return func(event *Event) error {
		event.ID = messageID
		return nil
	}
}

// SendWithSession configures the message to send with a specific session and sequence. By default, a sender has a
// default session (uuid.NewV4()) and sequence generator.
func SendWithSession(sessionID string, sequence uint32) SendOption {
	return func(event *Event) error {
		event.GroupID = &sessionID
		event.GroupSequence = &sequence
		return nil
	}
}

// SendWithoutSessionID will set the SessionID to nil. If sending to a partitioned Service Bus queue, this will cause
// the queue distributed the message in a round robin fashion to the next available partition with the effect of not
// enforcing FIFO ordering of messages, but enabling more efficient distribution of messages across partitions.
func SendWithoutSessionID() SendOption {
	return func(event *Event) error {
		event.GroupID = nil
		event.GroupSequence = nil
		return nil
	}
}
