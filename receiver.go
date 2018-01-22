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
	session    *amqp.Session
	receiver   *amqp.Receiver
	entityPath string
}

// NewReceiver creates a new Service Bus message listener given an AMQP client and an entity path
func NewReceiver(client *amqp.Client, entityPath string) (*Receiver, error) {
	receiver := &Receiver{
		client:     client,
		entityPath: entityPath,
	}
	err := receiver.newSessionAndLink()
	if err != nil {
		return nil, err
	}
	return receiver, nil
}

// Close will close the AMQP session and link of the receiver
func (r *Receiver) Close() error {
	err := r.session.Close()
	if err != nil {
		return err
	}

	err = r.receiver.Close()
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

// Receive start a listener for messages sent to the entity path
func (r *Receiver) Receive(ctx context.Context, handler Handler) {
	go func() {
	loop:
		for {
			select {
			case <-ctx.Done():
				r.Close()
				break loop
			default:
				err := r.handleMessageWithAcceptReject(handler)
				if err != nil {
					log.Warnln(err)
				}
			}
		}
	}()
}

// newSessionAndLink will replace the session and link on the receiver
func (r *Receiver) newSessionAndLink() error {
	session, err := r.client.NewSession()
	if err != nil {
		return err
	}

	amqpReceiver, err := session.NewReceiver(
		amqp.LinkAddress(r.entityPath),
		amqp.LinkCredit(10),
		amqp.LinkBatching(true))
	if err != nil {
		return err
	}

	r.session = session
	r.receiver = amqpReceiver

	return nil
}

// handleMessageWithAcceptReject will fetch a message from the entity path and accept if handler returns
// a nil error or it will reject the message if the handler returns a non-nil error
func (r *Receiver) handleMessageWithAcceptReject(handler Handler) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// Receive next message
	msg, err := r.receiver.Receive(ctx)
	cancel()
	if err, ok := err.(net.Error); ok && err.Timeout() {
		return nil
	} else if err != nil {
		log.Fatalln(err)
	}

	if msg != nil {
		id := interface{}("null")
		if msg.Properties != nil {
			id = msg.Properties.MessageID
		}
		log.Debugf("Message received: %s, id: %s", msg.Data, id)

		err = handler(ctx, msg)
		if err != nil {
			msg.Reject()
			log.Debugln("Message rejected")
		} else {
			// Accept message
			msg.Accept()
			log.Debugln("Message accepted")
		}
	}
	return nil
}
