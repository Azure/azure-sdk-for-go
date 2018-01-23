package servicebus

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net"
	"pack.ag/amqp"
	"time"
)

// Receiver provides session and link handling for a receiving entity path
type Receiver struct {
	client     *amqp.Client
	session    *Session
	receiver   *amqp.Receiver
	entityPath string
	done       chan struct{}
}

// NewReceiver creates a new Service Bus message listener given an AMQP client and an entity path
func NewReceiver(client *amqp.Client, entityPath string) (*Receiver, error) {
	receiver := &Receiver{
		client:     client,
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
func (r *Receiver) Close() error {
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
func (r *Receiver) Recover() error {
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
func (r *Receiver) Listen(handler Handler) {
	messages := make(chan *amqp.Message)
	go r.listenForMessages(messages)
	go r.handleMessages(messages, handler)
}

func (r *Receiver) handleMessages(messages chan *amqp.Message, handler Handler) {
	for {
		select {
		case <-r.done:
			log.Debug("done handling messages")
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

func (r *Receiver) listenForMessages(msgChan chan *amqp.Message) {
	for {
		select {
		case <-r.done:
			log.Debug("done listenting for messages")
			close(msgChan)
			return
		default:
			log.Debug("attempting to receive messages")
			waitCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			msg, err := r.receiver.Receive(waitCtx)
			cancel()
			if err, ok := err.(net.Error); ok && err.Timeout() {
				log.Debug("attempting to receive messages timed out")
				continue
			} else if err != nil {
				log.Fatalln(err)
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
func (r *Receiver) newSessionAndLink() error {
	amqpSession, err := r.client.NewSession()
	if err != nil {
		return err
	}

	amqpReceiver, err := amqpSession.NewReceiver(
		amqp.LinkAddress(r.entityPath),
		amqp.LinkCredit(10),
		amqp.LinkBatching(true))
	if err != nil {
		return err
	}

	r.session = NewSession(amqpSession)
	r.receiver = amqpReceiver

	return nil
}
