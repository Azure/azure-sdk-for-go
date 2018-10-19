package servicebus

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-amqp-common-go/log"
	"github.com/Azure/azure-amqp-common-go/rpc"
	"go.opencensus.io/trace"
	"pack.ag/amqp"
)

//RenewLocks renews the locks on messages provided
func (e *entity) RenewLocks(ctx context.Context, messages []*Message) error {
	span, ctx := e.startSpanFromContext(ctx, "sb.entity.renewLocks")
	defer span.Finish()

	lockTokens := make([]amqp.UUID, 0, len(messages))
	for _, m := range messages {
		if m.LockToken == nil {
			log.For(ctx).Error(fmt.Errorf("failed: message has nil lock token, cannot renew lock"), trace.StringAttribute("messageId", m.ID))
			continue
		}

		amqpLockToken := amqp.UUID(*m.LockToken)
		lockTokens = append(lockTokens, amqpLockToken)
	}

	if len(lockTokens) < 1 {
		log.For(ctx).Info("no lock tokens present to renew")
		return nil
	}

	e.renewMessageLockMutex.Lock()
	defer e.renewMessageLockMutex.Unlock()

	renewRequestMsg := &amqp.Message{
		ApplicationProperties: map[string]interface{}{
			operationFieldName: serviceBuslockRenewalOperationName,
		},
		Value: map[string]interface{}{
			lockTokensFieldName: lockTokens,
		},
	}

	entityManagementAddress := e.ManagementPath()
	conn, err := e.namespace.newConnection()
	if err != nil {
		return err
	}
	err = e.namespace.negotiateClaim(ctx, conn, entityManagementAddress)
	if err != nil {
		log.For(ctx).Error(err)
		return err
	}

	rpcLink, err := rpc.NewLink(conn, entityManagementAddress)
	if err != nil {
		return err
	}

	response, err := rpcLink.RetryableRPC(ctx, 3, 1*time.Second, renewRequestMsg)
	if err != nil {
		return err
	}

	if response.Code != 200 {
		return fmt.Errorf("error renewing locks: %v", response.Description)
	}

	return nil
}
