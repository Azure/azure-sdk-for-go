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
}

func newMessageSession(s *sender, r *receiver) *MessageSession {
	panic("not implemented")
}

func (ms *MessageSession) LockedUntil() time.Time {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	panic("not implemented")
}

func (ms *MessageSession) RenewSessionLock() error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

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
