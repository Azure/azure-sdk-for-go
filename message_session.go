package servicebus

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Azure/azure-amqp-common-go/rpc"
	"github.com/Azure/azure-amqp-common-go/uuid"
	"pack.ag/amqp"
)

// MessageSession represents and allows for interaction with a Service Bus Session. Service Bus Sessions
type MessageSession struct {
	mu sync.RWMutex
	*entity
	*receiver
	sessionID      string
	lockExpiration time.Time
}

func newMessageSession(r *receiver, e *entity, sessionID string) (retval *MessageSession, _ error) {
	retval = &MessageSession{
		receiver:       r,
		entity:         e,
		sessionID:      sessionID,
		lockExpiration: time.Now(),
	}

	return
}

// LockedUntil fetches the moment in time when the Session lock held by this receiver
func (ms *MessageSession) LockedUntil() time.Time {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	return ms.lockExpiration
}

// Renew requests that the Service Bus Server renews this client's lock on an existing Session.
func (ms *MessageSession) Renew(ctx context.Context) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	link, err := rpc.NewLink(ms.receiver.connection, ms.entity.ManagementPath())
	if err != nil {
		return err
	}

	reqID, err := uuid.NewV4()
	if err != nil {
		return err
	}

	msg := amqp.NewMessage([]byte{})
	msg.ApplicationProperties["operation"] = "com.microsoft:renew-session-lock"
	msg.ApplicationProperties["com.microsoft:tracking-id"] = reqID.String()

	if deadline, ok := ctx.Deadline(); ok {
		msg.ApplicationProperties["com.microsoft:server-timeout"] = durationTo8601Seconds(time.Until(deadline))
	}

	resp, err := link.RetryableRPC(ctx, 5, 5*time.Second, msg)
	if err != nil {
		return err
	}

	if rawUpdatedLockExpiration, ok := resp.Message.ApplicationProperties["expiration"]; ok {
		ms.lockExpiration = rawUpdatedLockExpiration.(time.Time)
	} else {
		return errors.New("lock expiration property not present")
	}

	return nil

}

// SetState updates the current State associated with this Session.
func (ms *MessageSession) SetState(ctx context.Context, state []byte) error {
	link, err := rpc.NewLink(ms.receiver.connection, ms.entity.ManagementPath())
	if err != nil {
		return err
	}

	reqID, err := uuid.NewV4()
	if err != nil {
		return err
	}

	msg := amqp.NewMessage([]byte{})
	msg.ApplicationProperties = map[string]interface{}{
		"operation":                 "com.microsoft:set-session-state",
		"com.microsoft:tracking-id": reqID.String(),
	}

	msg.Value = map[string]interface{}{
		"session-id":    ms.SessionID(),
		"session-state": state,
	}

	rsp, err := link.RetryableRPC(ctx, 5, 5*time.Second, msg)
	if err != nil {
		return err
	}

	if rsp.Code != 200 {
		return fmt.Errorf("amqp error (%d): %q", rsp.Code, rsp.Description)
	}
	return nil
}

// State retrieves the current State associated with this Session.
// https://docs.microsoft.com/en-us/azure/service-bus-messaging/service-bus-amqp-request-response#get-session-state
func (ms *MessageSession) State(ctx context.Context) ([]byte, error) {
	link, err := rpc.NewLink(ms.receiver.connection, ms.entity.ManagementPath())
	if err != nil {
		return []byte{}, err
	}

	reqID, err := uuid.NewV4()
	if err != nil {
		return []byte{}, err
	}

	msg := amqp.NewMessage([]byte{})

	msg.ApplicationProperties = map[string]interface{}{
		"operation":                     "com.microsoft:get-session-state",
		"com.microsoft.com:tracking-id": reqID.String(),
	}

	msg.Value = map[string]interface{}{
		"session-id": ms.SessionID(),
	}

	rsp, err := link.RetryableRPC(ctx, 5, 5*time.Second, msg)
	if err != nil {
		return []byte{}, err
	}

	if rsp.Code != 200 {
		return []byte{}, fmt.Errorf("amqp error (%d): %q", rsp.Code, rsp.Description)
	}

	if val, ok := rsp.Message.Value.(map[string]interface{}); ok {
		if rawState, ok := val["session-state"]; ok {
			if state, ok := rawState.([]byte); ok {
				return state, nil
			}
			return []byte{}, errors.New("server error: response value \"session-state\" is not a byte array")
		}
		return []byte{}, errors.New("server error: response did not contain value \"session-state\")")
	}
	return []byte{}, errors.New("server error: response value was not of expected type map[string]interface{}")
}

// SessionID gets the unique identifier of the session being interacted with by this MessageSession.
func (ms *MessageSession) SessionID() string {
	return ms.sessionID
}
