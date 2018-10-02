package servicebus

import (
	"context"
	"sync"
	"time"
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
	panic("not implemented")
}

// SetState updates the current State associated with this Session.
func (ms *MessageSession) SetState(ctx context.Context, state []byte) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	panic("not implemented")
}

// State retrieves the current State associated with this Session.
func (ms *MessageSession) State(ctx context.Context) ([]byte, error) {
	ms.mu.RLock()
	defer ms.mu.Unlock()

	panic("not implemented")
}

func (ms *MessageSession) SessionID() string {
	return ms.sessionID
}
