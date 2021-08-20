package utils

import (
	"context"
	"sync"
)

// WaitForGroupOrContext will wait for wg or the context to complete.
func WaitForGroupOrContext(wg *sync.WaitGroup, ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}
