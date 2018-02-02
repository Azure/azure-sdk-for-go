package servicebus

import (
	"context"
	"pack.ag/amqp"
	"fmt"
	"time"
	log "github.com/sirupsen/logrus"
)

const (
	CbsAddress       = "$cbs"
	CbsReplyToPrefix = "cbs-tmp-"
)

type (
	cbsLink struct {
		session       *amqp.Session
		receiver      *amqp.Receiver
		sender        *amqp.Sender
		clientAddress string
	}
)

func (sb *serviceBus) newCbsLink() (*cbsLink, error) {
	authSession, err := sb.client.NewSession()
	if err != nil {
		return nil, err
	}

	authSender, err := authSession.NewSender(amqp.LinkTargetAddress(CbsAddress))
	if err != nil {
		return nil, err
	}

	cbsClientAddress := CbsReplyToPrefix + sb.name.String()
	authReceiver, err := authSession.NewReceiver(
		amqp.LinkSourceAddress(CbsAddress),
		amqp.LinkTargetAddress(cbsClientAddress))
	if err != nil {
		return nil, err
	}

	return &cbsLink{
		sender:        authSender,
		receiver:      authReceiver,
		session:       authSession,
		clientAddress: cbsClientAddress,
	}, nil
}

func (sb *serviceBus) negotiateClaim(entityPath string) error {
	sb.cbsMu.Lock()
	defer sb.cbsMu.Unlock()

	name := "amqp://" + sb.namespace + ".servicebus.windows.net/" + entityPath
	log.Debugf("sending to: %s, expiring on: %q, via: %s", name, sb.sbToken.ExpiresOn, sb.cbsLink.clientAddress)
	msg := &amqp.Message{
		Value: sb.sbToken.AccessToken,
		Properties: &amqp.MessageProperties{
			ReplyTo:   sb.cbsLink.clientAddress,
		},
		ApplicationProperties: map[string]interface{}{
			"operation":  "put-token",
			"type":       "jwt",
			"name":       name,
			"expiration": sb.sbToken.ExpiresOn,
		},
	}

	_, err := retry(3, 1*time.Second, func() (interface{}, error) {
		log.Debugf("Attempting to negotiate cbs for %s in namespace %s", entityPath, sb.namespace)
		err := sb.cbsLink.send(context.Background(), msg)
		if err != nil {
			return nil, err
		}

		res, err := sb.cbsLink.receive(context.Background())
		if err != nil {
			return nil, err
		}

		if statusCode, ok := res.ApplicationProperties["status-code"].(int32); ok {
			description := res.ApplicationProperties["status-description"].(string)
			if statusCode >= 200 && statusCode < 300 {
				log.Debugf("Successfully negotiated cbs for %s in namespace %s", entityPath, sb.namespace)
				return res, nil
			} else if statusCode >= 500 {
				log.Debugf("Re-negotiating cbs for %s in namespace %s. Received status code: %d and error: %s", entityPath, sb.namespace, statusCode, description)
				return nil, &retryable{message: "cbs error: " + description}
			} else {
				log.Debugf("Failed negotiating cbs for %s in namespace %s with error %d", entityPath, sb.namespace, statusCode)
				return nil, fmt.Errorf("cbs error: failed with code %d and message: %s", statusCode, description)
			}
		}

		return nil, &retryable{message: "cbs error: didn't understand the replied message status code"}
	})

	return err
}

func (cl *cbsLink) forceClose() {
	if cl.sender != nil {
		cl.sender.Close()
	}

	if cl.receiver != nil {
		cl.receiver.Close()
	}

	if cl.session != nil {
		cl.session.Close()
	}
}

func (cl *cbsLink) send(ctx context.Context, msg *amqp.Message) error {
	return cl.sender.Send(ctx, msg)
}

func (cl *cbsLink) receive(ctx context.Context) (*amqp.Message, error) {
	return cl.receiver.Receive(ctx)
}
