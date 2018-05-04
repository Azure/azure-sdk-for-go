package servicebus

import (
	"context"
	"fmt"
	"time"

	"pack.ag/amqp"

	"github.com/Azure/azure-amqp-common-go"
	"github.com/Azure/azure-amqp-common-go/log"
	"github.com/opentracing/opentracing-go"
)

// receiver provides session and link handling for a receiving entity path
type (
	receiver struct {
		namespace         *Namespace
		connection        *amqp.Client
		session           *session
		receiver          *amqp.Receiver
		entityPath        string
		done              func()
		Name              string
		requiredSessionID *string
		lastError         error
	}

	// ReceiverOptions provides a structure for configuring receivers
	ReceiverOptions func(receiver *receiver) error

	// ListenerHandle provides the ability to close or listen to the close of a Receiver
	ListenerHandle struct {
		r   *receiver
		ctx context.Context
	}
)

// newReceiver creates a new Service Bus message listener given an AMQP client and an entity path
func (ns *Namespace) newReceiver(ctx context.Context, entityPath string, opts ...ReceiverOptions) (*receiver, error) {
	span, ctx := ns.startSpanFromContext(ctx, "servicebus.Hub.newReceiver")
	defer span.Finish()

	receiver := &receiver{
		namespace:  ns,
		entityPath: entityPath,
	}

	for _, opt := range opts {
		if err := opt(receiver); err != nil {
			return nil, err
		}
	}

	err := receiver.newSessionAndLink(ctx)
	return receiver, err
}

// Close will close the AMQP session and link of the receiver
func (r *receiver) Close(ctx context.Context) error {
	if r.done != nil {
		r.done()
	}

	return r.connection.Close()
}

// Recover will attempt to close the current session and link, then rebuild them
func (r *receiver) Recover(ctx context.Context) error {
	_ = r.Close(ctx) // we expect the receiver is in an error state
	return r.newSessionAndLink(ctx)
}

// Listen start a listener for messages sent to the entity path
func (r *receiver) Listen(handler Handler) *ListenerHandle {
	ctx, done := context.WithCancel(context.Background())
	r.done = done

	span, ctx := r.startConsumerSpanFromContext(ctx, "servicebus.receiver.Listen")
	defer span.Finish()

	messages := make(chan *amqp.Message)
	go r.listenForMessages(ctx, messages)
	go r.handleMessages(ctx, messages, handler)

	return &ListenerHandle{
		r:   r,
		ctx: ctx,
	}
}

func (r *receiver) handleMessages(ctx context.Context, messages chan *amqp.Message, handler Handler) {
	span, ctx := r.startConsumerSpanFromContext(ctx, "servicebus.receiver.handleMessages")
	defer span.Finish()
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-messages:
			r.handleMessage(ctx, msg, handler)
		}
	}
}

func (r *receiver) handleMessage(ctx context.Context, msg *amqp.Message, handler Handler) {
	event := eventFromMsg(msg)
	var span opentracing.Span
	wireContext, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, event)
	if err == nil {
		span, ctx = r.startConsumerSpanFromWire(ctx, "servicebus.receiver.handleMessage", wireContext)
	} else {
		span, ctx = r.startConsumerSpanFromContext(ctx, "servicebus.receiver.handleMessage")
	}
	defer span.Finish()

	id := messageID(msg)
	span.SetTag("amqp.message-id", id)

	err = handler(ctx, event)
	if err != nil {
		msg.Reject()
		log.For(ctx).Error(fmt.Errorf("message rejected: id: %v", id))
		return
	}
	msg.Accept()
}

func (r *receiver) listenForMessages(ctx context.Context, msgChan chan *amqp.Message) {
	span, ctx := r.startConsumerSpanFromContext(ctx, "servicebus.receiver.listenForMessages")
	defer span.Finish()

	for {
		msg, err := r.listenForMessage(ctx)
		if ctx.Err() != nil && ctx.Err() == context.DeadlineExceeded {
			return
		}

		if err != nil {
			_, retryErr := common.Retry(5, 10*time.Second, func() (interface{}, error) {
				sp, ctx := r.startConsumerSpanFromContext(ctx, "servicebus.receiver.listenForMessages.tryRecover")
				defer sp.Finish()

				err := r.Recover(ctx)
				if ctx.Err() != nil && ctx.Err() == context.DeadlineExceeded {
					return nil, ctx.Err()
				}

				if err != nil {
					log.For(ctx).Error(err)
					return nil, common.Retryable(err.Error())
				}
				return nil, nil
			})

			if retryErr != nil {
				r.lastError = retryErr
				r.Close(ctx)
				return
			}
			continue
		}
		select {
		case msgChan <- msg:
		case <-ctx.Done():
			return
		}
	}
}

func (r *receiver) listenForMessage(ctx context.Context) (*amqp.Message, error) {
	span, ctx := r.startConsumerSpanFromContext(ctx, "servicebus.receiver.listenForMessage")
	defer span.Finish()

	msg, err := r.receiver.Receive(ctx)
	if err != nil {
		log.For(ctx).Debug(err.Error())
		return nil, err
	}

	id := messageID(msg)
	span.SetTag("amqp.message-id", id)
	return msg, nil
}

// newSessionAndLink will replace the session and link on the receiver
func (r *receiver) newSessionAndLink(ctx context.Context) error {
	connection, err := r.namespace.newConnection()
	if err != nil {
		return err
	}
	r.connection = connection

	err = r.namespace.negotiateClaim(ctx, connection, r.entityPath)
	if err != nil {
		log.For(ctx).Error(err)
		return err
	}

	amqpSession, err := connection.NewSession()
	if err != nil {
		log.For(ctx).Error(err)
		return err
	}

	r.session, err = newSession(amqpSession)
	if err != nil {
		log.For(ctx).Error(err)
		return err
	}

	opts := []amqp.LinkOption{
		amqp.LinkSourceAddress(r.entityPath),
		amqp.LinkCredit(100),
	}

	// TODO: fix this with after SB team replies with bug fix for session filters
	//if r.requiredSessionID != nil {
	//	opts = append(opts, amqp.LinkSourceFilterString("com.microsoft:session-filter", *r.requiredSessionID))
	//	r.session.SessionID = *r.requiredSessionID
	//}

	amqpReceiver, err := amqpSession.NewReceiver(opts...)
	if err != nil {
		return err
	}

	r.receiver = amqpReceiver
	return nil
}

// ReceiverWithSession configures a receiver to use a session
func ReceiverWithSession(sessionID string) ReceiverOptions {
	return func(r *receiver) error {
		r.requiredSessionID = &sessionID
		return nil
	}
}

func messageID(msg *amqp.Message) interface{} {
	var id interface{} = "null"
	if msg.Properties != nil {
		id = msg.Properties.MessageID
	}
	return id
}

// Close will close the listener
func (lc *ListenerHandle) Close(ctx context.Context) error {
	return lc.r.Close(ctx)
}

// Done will close the channel when the listener has stopped
func (lc *ListenerHandle) Done() <-chan struct{} {
	return lc.ctx.Done()
}

// Err will return the last error encountered
func (lc *ListenerHandle) Err() error {
	if lc.r.lastError != nil {
		return lc.r.lastError
	}
	return lc.ctx.Err()
}
