// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"sync"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/internal/amqpwrap"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	otellog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/noop"
)

type lifecycleTestLogger struct {
	noop.Logger
	callback func(component, outcome string)
}

func (l lifecycleTestLogger) Enabled(context.Context, otellog.EnabledParameters) bool { return true }

func (l lifecycleTestLogger) Emit(_ context.Context, record otellog.Record) {
	var component, outcome string
	record.WalkAttributes(func(kv attribute.KeyValue) bool {
		switch kv.Key {
		case "component":
			component = kv.Value.AsString()
		case "outcome":
			outcome = kv.Value.AsString()
		}
		return true
	})
	l.callback(component, outcome)
}

func closeOnLifecycle(t *testing.T, mu *sync.RWMutex, trigger string, closeClient func() error, called *bool) *Diagnostics {
	t.Helper()
	return &Diagnostics{logger: lifecycleTestLogger{callback: func(_, outcome string) {
		// Check the lock before calling Close so a regression produces a useful assertion.
		require.True(t, mu.TryLock(), "application logger called under AMQP cache lock")
		mu.Unlock()
		if outcome == trigger {
			require.NoError(t, closeClient())
			*called = true
		}
	}}}
}

func TestLifecycleLoggerCanCloseConnection(t *testing.T) {
	for _, trigger := range []string{"started", "success", "invalidated"} {
		t.Run(trigger, func(t *testing.T) {
			ctx := t.Context()
			candidate := &fakeAMQPClient{}
			ns := &Namespace{newClientFn: func(context.Context, uint64) (amqpwrap.AMQPClient, error) { return candidate, nil }}
			var called bool
			ns.Diagnostics = closeOnLifecycle(t, &ns.clientMu, trigger, func() error { return ns.Close(ctx, true) }, &called)
			if trigger == "invalidated" {
				old := &fakeAMQPClient{}
				ns.client = old
				require.ErrorIs(t, ns.Recover(ctx, ns.connID), ErrClientClosed)
				require.Equal(t, 1, old.closeCalled)
			} else {
				_, _, err := ns.GetAMQPClientImpl(ctx)
				require.ErrorIs(t, err, ErrClientClosed)
				require.Equal(t, 1, candidate.closeCalled)
			}
			require.True(t, called)
			require.Nil(t, ns.client)
		})
	}
}

func TestLifecycleLoggerCanCloseLinks(t *testing.T) {
	for _, trigger := range []string{"started", "success", "invalidated"} {
		t.Run(trigger, func(t *testing.T) {
			ctx := t.Context()
			receiver := &FakeAMQPReceiver{NameForLink: "original"}
			links := NewLinks(&FakeNSForPartClient{}, "", func(string) string { return "0" }, func(context.Context, amqpwrap.AMQPSession, string, string) (*FakeAMQPReceiver, error) {
				return receiver, nil
			})
			var called bool
			links.WithDiagnostics(closeOnLifecycle(t, &links.linksMu, trigger, func() error { return links.Close(ctx) }, &called), "receiver")
			link, err := links.GetLink(ctx, "0")
			if trigger == "invalidated" {
				require.NoError(t, err)
				require.NoError(t, links.closePartitionLinkIfMatch(ctx, "0", "stale"))
				require.False(t, called)
				require.NoError(t, links.closePartitionLinkIfMatch(ctx, "0", link.Link().LinkName()))
			} else {
				require.ErrorIs(t, err, ErrClientClosed)
			}
			require.True(t, called)
			require.Equal(t, 1, receiver.CloseCalled)
			require.Nil(t, links.links)
		})
	}
}

func TestConcurrentColdConnection(t *testing.T) {
	entered, finish := make(chan struct{}), make(chan struct{})
	candidate := &fakeAMQPClient{}
	opens := 0
	ns := &Namespace{newClientFn: func(context.Context, uint64) (amqpwrap.AMQPClient, error) {
		opens++
		close(entered)
		<-finish
		return candidate, nil
	}}
	type result struct {
		conn amqpwrap.AMQPClient
		id   uint64
		err  error
	}
	results := make(chan result, 2)
	get := func() {
		conn, id, err := ns.GetAMQPClientImpl(t.Context())
		results <- result{conn, id, err}
	}
	go get()
	<-entered
	go get()
	close(finish)
	for range 2 {
		r := <-results
		require.NoError(t, r.err)
		require.Same(t, candidate, r.conn)
		require.EqualValues(t, 1, r.id)
	}
	require.Equal(t, 1, opens)
}
