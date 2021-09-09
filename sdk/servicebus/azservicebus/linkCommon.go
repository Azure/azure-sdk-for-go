package azservicebus

import "context"

// linkState is a wrapper around a context that lets us manage the
// "user has closed this link, don't recreate it" logic in both
// Senders and Receivers.
type linkState struct {
	context.Context
	Close context.CancelFunc
	err   error
}

func newLinkState(ctx context.Context, errorOnCancel error) *linkState {
	ctx, cancel := context.WithCancel(ctx)

	return &linkState{
		Context: ctx,
		Close:   cancel,
		err:     errorOnCancel,
	}
}

// Err returns the "link is closed" error when the link has been closed
// or nil otherwise.
func (ls *linkState) Err() error {
	if ls.Context.Err() != nil {
		return ls.err
	}

	return nil
}

// Closed checks if the link has been closed by the user. Can be used as
// a quick check before running any operation like sending or receiving messages.
func (ls *linkState) Closed() bool {
	select {
	case <-ls.Context.Done():
		return true
	default:
		return false
	}
}
