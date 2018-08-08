package servicebus

import (
	"context"
	"fmt"
	"github.com/Azure/azure-amqp-common-go/log"
	"github.com/Azure/azure-amqp-common-go/uuid"
	otlogger "github.com/opentracing/opentracing-go/log"

	"pack.ag/amqp"
)

const (
	serviceBuslockRenewalOperationName = "com.microsoft:renew-lock"
)

func (e *entity) createManagementChannels(ctx context.Context, amqpPath string) (*sender, *receiver, error) {
	if e.namespace == nil {
		panic("expect namespace not nil")
	}

	subscriptionAddress := e.namespace.getEntityManagementPath(amqpPath)

	// receiver := uuid.NewV4().String()
	sender, err := e.namespace.newSender(ctx, subscriptionAddress)

	if err != nil {
		return nil, nil, err
	}

	reciever, err := e.namespace.newReceiver(ctx, subscriptionAddress)

	if err != nil {
		return nil, nil, err
	}

	return sender, reciever, nil
}

//RenewLocks renews the locks on messages provided
func (e *entity) RenewLocks(ctx context.Context, messages []*Message) error {
	lockTokens := make([]*uuid.UUID, 0, len(messages))
	for _, m := range messages {
		if m.LockToken == nil {
			log.For(ctx).Error(fmt.Errorf("failed: message has nil lock token, cannot renew lock"), otlogger.Object("messageId", m))
			continue
		}

		lockTokens = append(lockTokens, m.LockToken)
	}

	if len(lockTokens) < 1 {
		log.For(ctx).Info("no lock tokens present to renew")
		return nil
	}

	e.renewMessageLockMutex.Lock()
	defer e.renewMessageLockMutex.Unlock()

	sender, receiver, err := e.createManagementChannels(ctx, e.Name)
	if err != nil {
		return err
	}

	defer sender.Close(ctx)
	defer receiver.Close(ctx)

	messageID, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("error creating messageID: %+v", err)
	}

	replyToAddress, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("error creating replyToAddress: %+v", err)
	}

	renewRequestMsg := &amqp.Message{
		ApplicationProperties: map[string]interface{}{
			"operation": serviceBuslockRenewalOperationName,
		},
		Properties: &amqp.MessageProperties{
			MessageID: messageID,
			ReplyTo:   replyToAddress.String(),
		},
		Value: map[string]interface{}{
			"lock-tokens": lockTokens,
		},
	}

	err = sender.Send(ctx, &Message{
		message: renewRequestMsg,
	})
	if err != nil {
		return err
	}

	response, err := receiver.ReceiveOne(ctx)
	if err != nil {
		return err
	}

	responseCode, ok := response.UserProperties["statusCode"]
	if !ok {
		return fmt.Errorf("unexpected response from renewal request: %+v", response)
	}

	if responseCode != "200" {
		return fmt.Errorf("error renewing locks: %v", response.UserProperties["statusDescription"])
	}

	return nil
}
