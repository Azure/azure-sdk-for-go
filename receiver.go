package servicebus

import (
	"context"
	"net"
	"time"

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
		done       chan struct{}
		Name       uuid.UUID
	}
)

// newReceiver creates a new Service Bus message listener given an AMQP client and an entity path
func (sb *serviceBus) newReceiver(entityPath string) (*receiver, error) {
	receiver := &receiver{
		sb:         sb,
		entityPath: entityPath,
		done:       make(chan struct{}),
	}
	err := receiver.newSessionAndLink()
	if err != nil {
		return nil, err
	}
	return receiver, nil
}

// Close will close the AMQP session and link of the receiver
func (r *receiver) Close() error {
	close(r.done)

	err := r.receiver.Close()
	if err != nil {
		return err
	}

	err = r.session.Close()
	if err != nil {
		return err
	}

	return nil
}

// Recover will attempt to close the current session and link, then rebuild them
func (r *receiver) Recover() error {
	err := r.Close()
	if err != nil {
		return err
	}

	err = r.newSessionAndLink()
	if err != nil {
		return err
	}

	return nil
}

// Listen start a listener for messages sent to the entity path
func (r *receiver) Listen(handler Handler) {
	messages := make(chan *amqp.Message)
	go r.listenForMessages(messages)
	go r.handleMessages(messages, handler)
}

func (r *receiver) handleMessages(messages chan *amqp.Message, handler Handler) {
	for {
		select {
		case <-r.done:
			log.Debug("done handling messages")
			close(messages)
			return
		case msg := <-messages:
			ctx := context.Background()
			id := interface{}("null")
			if msg.Properties != nil {
				id = msg.Properties.MessageID
			}
			log.Debugf("Message id: %s is being passed to handler", id)
			err := handler(ctx, msg)

			if err != nil {
				msg.Reject()
				log.Debugf("Message rejected: id: %s", id)
			} else {
				// Accept message
				msg.Accept()
				log.Debugf("Message accepted: id: %s", id)
			}
		}
	}
}

func (r *receiver) listenForMessages(msgChan chan *amqp.Message) {
	for {
		select {
		case <-r.done:
			log.Debug("done listenting for messages")
			return
		default:
			//log.Debug("attempting to receive messages")
			waitCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			msg, err := r.receiver.Receive(waitCtx)
			cancel()

			// TODO: handle receive errors better. It's not sufficient to check only for timeout
			if err, ok := err.(net.Error); ok && err.Timeout() {
				log.Debug("attempting to receive messages timed out")
				continue
			} else if err != nil {
				log.Fatalln(err)
				time.Sleep(10 * time.Second)
			}
			if msg != nil {
				id := interface{}("null")
				if msg.Properties != nil {
					id = msg.Properties.MessageID
				}
				log.Debugf("Message received: %s", id)
				msgChan <- msg
			}
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
		amqp.LinkCredit(10))
	if err != nil {
		return err
	}

	r.session = newSession(amqpSession)
	r.receiver = amqpReceiver

	return nil
}
