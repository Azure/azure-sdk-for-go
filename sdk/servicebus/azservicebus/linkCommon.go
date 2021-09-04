package azservicebus

import "context"

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

func (ls *linkState) Err() error {
	if ls.Context.Err() != nil {
		return ls.err
	}

	return nil
}

func (ls *linkState) Closed() bool {
	select {
	case <-ls.Context.Done():
		return true
	default:
		return false
	}
}
