package servicebus

import (
	"context"
	"sync"
	"time"
)

// MessageSession represents and allows for interaction with a Service Bus Session. Service Bus Sessions
type MessageSession struct {
	mu sync.RWMutex
	*receiver
	*sender
	SessionID string
}

func newMessageSession(s *sender, r *receiver) *MessageSession {
	panic("not implemented")
}

func (ms *MessageSession) LockedUntil() time.Time {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	panic("not implemented")
}

// RenewSessionLock requests that the Service Bus Server renews this client's lock on an existing Session.
func (ms *MessageSession) RenewSessionLock(ctx context.Context) error {
	panic("not implemented")
}

// SetState updates the current State associated with this Session.
func (ms *MessageSession) SetState(ctx context.Context, state []byte) error {
	ms.receiver.newSessionAndLink(ctx)
	panic("not implemented")
}

// State retrieves the current State associated with this Session.
func (ms *MessageSession) State(ctx context.Context) ([]byte, error) {
	panic("not implemented")
}
