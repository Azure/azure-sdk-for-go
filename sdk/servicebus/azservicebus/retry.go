package azservicebus

import "context"

type RetryPolicyOptions struct {
}

type RetryPolicy interface {
	Wait(ctx context.Context, c int) error
	Reset() error
}

type NullRetryPolicy struct{}

func (rp *NullRetryPolicy) Wait(ctx context.Context, c int) error {
	return nil
}

func (rp *NullRetryPolicy) Reset() error {
	return nil
}
