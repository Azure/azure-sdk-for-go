package servicebus

import (
	"context"

	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"pack.ag/amqp"
)

// receiver provides session and link handling for a receiving entity path
type (
	receiver struct {
		sb         *serviceBus
		session    *session
		receiver   *amqp.Receiver
		entityPath string
		done       func()
		Name       uuid.UUID
	}
)

// newReceiver creates a new Service Bus message listener given an AMQP client and an entity path
func (sb *serviceBus) newReceiver(entityPath string) (*receiver, error) {
	receiver := &receiver{
		sb:         sb,
		entityPath: entityPath,
	}
	err := receiver.newSessionAndLink()
	return receiver, err
}

// Close will close the AMQP session and link of the receiver
func (r *receiver) Close() error {
	// This isn't safe to be called concurrently with Listen
	if r.done != nil {
		r.done()
	}
	err := r.receiver.Close()
	if err != nil {
		// ensure session is closed on receiver error
		_ = r.session.Close()
		return err
	}

	return r.session.Close()
}

// Recover will attempt to close the current session and link, then rebuild them
func (r *receiver) Recover() error {
	err := r.Close()
	if err != nil {
		return err
	}

	return r.newSessionAndLink()
}

// Listen start a listener for messages sent to the entity path
func (r *receiver) Listen(handler Handler) {
	ctx, done := context.WithCancel(context.Background())
	r.done = done
	messages := make(chan *amqp.Message)
	go r.listenForMessages(ctx, messages)
	go r.handleMessages(ctx, messages, handler)
}

func (r *receiver) handleMessages(ctx context.Context, messages chan *amqp.Message, handler Handler) {
	for {
		select {
		case <-ctx.Done():
			log.Debug("done handling messages")
			return
		case msg := <-messages:
			var id interface{} = "null"
			if msg.Properties != nil {
				id = msg.Properties.MessageID
			}

			log.Debugf("Message id: %s is being passed to handler", id)
			err := handler(ctx, msg)
			if err != nil {
				msg.Reject()
				log.Debugf("Message rejected: id: %s", id)
				continue
			}

			// Accept message
			msg.Accept()
			log.Debugf("Message accepted: id: %s", id)
		}
	}
}

func (r *receiver) listenForMessages(ctx context.Context, msgChan chan *amqp.Message) {
	for {
		//log.Debug("attempting to receive messages")
		msg, err := r.receiver.Receive(ctx)
		// TODO (vcabbage): This previously checked `net.Error.Timeout() == true`, which
		//                  should never happen. If it does it's a bug in pack.ag/amqp.
		if err != nil {
			if ctx.Err() != nil {
				return
			}

			// TODO (vcabbage): I'm not sure what the appropriate action is here, this was
			//                 previously a call to `log.Fatalln`, which calls os.Exit(1).
			log.Error(err)
			return
		}

		var id interface{} = "null"
		if msg.Properties != nil {
			id = msg.Properties.MessageID
		}
		log.Debugf("Message received: %s", id)

		select {
		case msgChan <- msg:
		case <-ctx.Done():
			return
		}
	}
}

// newSessionAndLink will replace the session and link on the receiver
func (r *receiver) newSessionAndLink() error {
	if r.sb.claimsBasedSecurityEnabled() {
		err := r.sb.negotiateClaim(r.entityPath)
		if err != nil {
			return err
		}
	}

	amqpSession, err := r.sb.newSession()
	if err != nil {
		return err
	}

	amqpReceiver, err := amqpSession.NewReceiver(
		amqp.LinkSourceAddress(r.entityPath),
		amqp.LinkCredit(10),
	)
	if err != nil {
		return err
	}

	r.session = newSession(amqpSession)
	r.receiver = amqpReceiver

	return nil
}
