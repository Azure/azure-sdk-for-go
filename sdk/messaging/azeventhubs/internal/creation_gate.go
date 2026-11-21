// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"sync"
)

// creationGate serializes AMQP creation separately from the cached resource.
// Close never waits on it, so a synchronous telemetry callback may close a client
// while a connection or link is being created. The creator checks for closure
// again before installing its result.
type creationGate struct {
	once  sync.Once
	token chan struct{}
}

func (g *creationGate) enter(ctx context.Context) error {
	g.once.Do(func() { g.token = make(chan struct{}, 1) })
	// Uncontended recovery must still invalidate an old connection if ctx is canceled.
	select {
	case g.token <- struct{}{}:
		return nil
	default:
	}
	select {
	case g.token <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (g *creationGate) leave() { <-g.token }
