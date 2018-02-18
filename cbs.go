package servicebus

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"pack.ag/amqp"
)

const (
	cbsAddress           = "$cbs"
	cbsReplyToPrefix     = "cbs-tmp-"
	cbsOperationKey      = "operation"
	cbsOperationPutToken = "put-token"
	cbsTokenTypeKey      = "type"
	cbsTokenTypeJwt      = "jwt"
	cbsAudienceKey       = "name"
	cbsExpirationKey     = "expiration"
	cbsStatusCodeKey     = "status-code"
	cbsDescriptionKey    = "status-description"
)

type (
	cbsLink struct {
		session       *amqp.Session
		receiver      *amqp.Receiver
		sender        *amqp.Sender
		clientAddress string
		negotiateMu   sync.Mutex
	}
)

func (sb *serviceBus) newCbsLink() (*cbsLink, error) {
	conn, err := sb.connection()
	if err != nil {
		return nil, err
	}
	authSession, err := conn.NewSession()
	if err != nil {
		return nil, err
	}

	authSender, err := authSession.NewSender(amqp.LinkTargetAddress(cbsAddress))
	if err != nil {
		return nil, err
	}

	cbsClientAddress := cbsReplyToPrefix + sb.name.String()
	authReceiver, err := authSession.NewReceiver(
		amqp.LinkSourceAddress(cbsAddress),
		amqp.LinkTargetAddress(cbsClientAddress),
	)
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

func (sb *serviceBus) ensureCbsLink() error {
	sb.cbsMu.Lock()
	defer sb.cbsMu.Unlock()

	if sb.cbsLink != nil {
		return nil
	}

	link, err := sb.newCbsLink()
	sb.cbsLink = link
	return err
}

func (sb *serviceBus) negotiateClaim(entityPath string) error {
	err := sb.ensureCbsLink()
	if err != nil {
		return err
	}
	sb.cbsLink.negotiateMu.Lock()
	defer sb.cbsLink.negotiateMu.Unlock()

	name := "amqp://" + sb.namespace + ".servicebus.windows.net/" + entityPath
	log.Debugf("sending to: %s, expiring on: %q, via: %s", name, sb.sbToken.ExpiresOn, sb.cbsLink.clientAddress)
	msg := &amqp.Message{
		Value: sb.sbToken.AccessToken,
		Properties: &amqp.MessageProperties{
			ReplyTo: sb.cbsLink.clientAddress,
		},
		ApplicationProperties: map[string]interface{}{
			cbsOperationKey:  cbsOperationPutToken,
			cbsTokenTypeKey:  cbsTokenTypeJwt,
			cbsAudienceKey:   name,
			cbsExpirationKey: sb.sbToken.ExpiresOn,
		},
	}

	_, err = retry(3, 1*time.Second, func() (interface{}, error) {
		log.Debugf("Attempting to negotiate cbs for %s in namespace %s", entityPath, sb.namespace)

		ctx := context.Background()

		err := sb.cbsLink.send(ctx, msg)
		if err != nil {
			return nil, err
		}

		res, err := sb.cbsLink.receive(ctx)
		if err != nil {
			return nil, err
		}

		statusCode, ok := res.ApplicationProperties[cbsStatusCodeKey].(int32)
		if !ok {
			return nil, &retryable{message: "cbs error: didn't understand the replied message status code"}
		}

		description, ok := res.ApplicationProperties[cbsDescriptionKey].(string)
		if !ok {
			return nil, &retryable{message: "cbs error: didn't understand the replied message description"}
		}

		switch {
		case statusCode >= 200 && statusCode < 300:
			log.Debugf("Successfully negotiated cbs for %s in namespace %s", entityPath, sb.namespace)
			return res, nil
		case statusCode >= 500:
			log.Debugf("Re-negotiating cbs for %s in namespace %s. Received status code: %d and error: %s", entityPath, sb.namespace, statusCode, description)
			return nil, &retryable{message: "cbs error: " + description}
		default:
			log.Debugf("Failed negotiating cbs for %s in namespace %s with error %d", entityPath, sb.namespace, statusCode)
			return nil, fmt.Errorf("cbs error: failed with code %d and message: %s", statusCode, description)
		}
	})

	return err
}

func (cl *cbsLink) send(ctx context.Context, msg *amqp.Message) error {
	return cl.sender.Send(ctx, msg)
}

func (cl *cbsLink) receive(ctx context.Context) (*amqp.Message, error) {
	return cl.receiver.Receive(ctx)
}
