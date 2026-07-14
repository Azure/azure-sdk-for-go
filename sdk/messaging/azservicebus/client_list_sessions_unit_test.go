// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"fmt"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/amqpwrap"
	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

// scriptedRPCLink scripts a sequence of responses for ListSessions pagination
// tests. Each call to RPC consumes the next entry; the request message is
// captured so tests can assert on operation name, paging parameters, and the
// last-updated-time mapping. RPC fails the test if exhausted to surface
// off-by-one pagination bugs.
type scriptedRPCLink struct {
	t         *testing.T
	calls     []*amqp.Message
	responses []*amqpwrap.RPCResponse
	// errs, when non-nil at a given call index, makes that RPC call return the
	// error instead of a response. Used to exercise retry/recovery paths.
	errs   []error
	closed int
}

func (l *scriptedRPCLink) Close(ctx context.Context) error {
	l.closed++
	return nil
}

func (l *scriptedRPCLink) RPC(ctx context.Context, msg *amqp.Message) (*amqpwrap.RPCResponse, error) {
	// A real RPC link honors context cancellation; mirror that so cancellation tests
	// exercise the same path as production.
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	idx := len(l.calls)
	l.calls = append(l.calls, msg)
	if idx < len(l.errs) && l.errs[idx] != nil {
		return nil, l.errs[idx]
	}
	if idx >= len(l.responses) {
		l.t.Fatalf("scriptedRPCLink: unexpected RPC call #%d (only %d responses scripted)", idx+1, len(l.responses))
	}
	return l.responses[idx], nil
}

// pathCapturingNS wraps internal.FakeNS to capture the management path
// passed to NewRPCLink. Tests use this to verify that ListSessionsForQueue
// and ListSessionsForSubscription construct the correct entity-relative
// management address.
type pathCapturingNS struct {
	*internal.FakeNS
	mu            sync.Mutex
	capturedPaths []string
	recoverCalls  []uint64
}

func (ns *pathCapturingNS) NewRPCLink(ctx context.Context, managementPath string) (amqpwrap.RPCLink, error) {
	ns.mu.Lock()
	ns.capturedPaths = append(ns.capturedPaths, managementPath)
	ns.mu.Unlock()
	return ns.FakeNS.NewRPCLink(ctx, managementPath)
}

func (ns *pathCapturingNS) ManagementPaths() []string {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	out := make([]string, len(ns.capturedPaths))
	copy(out, ns.capturedPaths)
	return out
}

// Recover records the connection revision it was called with so tests can assert
// that a connection-level failure triggered exactly one revision-guarded recovery.
func (ns *pathCapturingNS) Recover(ctx context.Context, clientRevision uint64) (bool, error) {
	ns.mu.Lock()
	ns.recoverCalls = append(ns.recoverCalls, clientRevision)
	ns.mu.Unlock()
	return ns.FakeNS.Recover(ctx, clientRevision)
}

func (ns *pathCapturingNS) RecoverCalls() []uint64 {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	out := make([]uint64, len(ns.recoverCalls))
	copy(out, ns.recoverCalls)
	return out
}

func newClientForListSessionsUnitTest(t *testing.T, rpcLink amqpwrap.RPCLink) (*Client, *pathCapturingNS) {
	t.Helper()
	fakeTokenCredential := struct{ azcore.TokenCredential }{}
	client, err := NewClient("fake.something", fakeTokenCredential, nil)
	require.NoError(t, err)
	ns := &pathCapturingNS{
		FakeNS: &internal.FakeNS{
			RPCLink: rpcLink,
		},
	}
	client.namespace = ns
	// Keep retries fast so the connection-recovery test does not sleep on the default backoff.
	client.retryOptions = RetryOptions{RetryDelay: time.Millisecond, MaxRetryDelay: time.Millisecond, MaxRetries: 3}
	return client, ns
}

func okPage(t *testing.T, ids ...string) *amqpwrap.RPCResponse {
	t.Helper()
	values := make([]any, len(ids))
	for i, id := range ids {
		values[i] = id
	}
	return &amqpwrap.RPCResponse{
		Code: 200,
		Message: &amqp.Message{
			Value: map[string]any{
				"sessions-ids": values,
			},
		},
	}
}

// makeIDs returns a slice of n synthetic session IDs of the form
// "prefix-000", "prefix-001", ... using a width that scales with n so the
// helper works correctly for any size, not just n <= 100.
func makeIDs(prefix string, n int) []string {
	width := len(fmt.Sprintf("%d", n-1))
	if width < 1 {
		width = 1
	}
	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = fmt.Sprintf("%s-%0*d", prefix, width, i)
	}
	return out
}

// drainPager iterates a ListSessions pager to completion, collecting all session
// IDs across pages and failing the test on any page error.
func drainPager(t *testing.T, pager *runtime.Pager[ListSessionsResponse]) []string {
	t.Helper()
	var got []string
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		got = append(got, page.Sessions...)
	}
	return got
}

func TestClient_ListSessionsForQueue_PaginatesUntilShortPage(t *testing.T) {
	// Three full pages (100 each) then a partial page (37): exactly 4 RPC calls
	// with skip 0, 100, 200, 300, stopping on the partial page rather than issuing
	// a 5th call.
	page1 := makeIDs("p1", 100)
	page2 := makeIDs("p2", 100)
	page3 := makeIDs("p3", 100)
	page4 := makeIDs("p4", 37)

	link := &scriptedRPCLink{
		t: t,
		responses: []*amqpwrap.RPCResponse{
			okPage(t, page1...),
			okPage(t, page2...),
			okPage(t, page3...),
			okPage(t, page4...),
		},
	}
	client, ns := newClientForListSessionsUnitTest(t, link)

	got := drainPager(t, client.NewListSessionsForQueuePager("myqueue", nil))

	expected := append(append(append(append([]string{}, page1...), page2...), page3...), page4...)
	require.Equal(t, expected, got)
	require.Len(t, link.calls, 4)

	// Each page opens its own RPC link, all addressed to the queue's management endpoint.
	require.Equal(t, []string{
		"myqueue/$management", "myqueue/$management",
		"myqueue/$management", "myqueue/$management",
	}, ns.ManagementPaths())

	// Skip values progress correctly across calls.
	for i, expectedSkip := range []int32{0, 100, 200, 300} {
		body, ok := link.calls[i].Value.(map[string]any)
		require.True(t, ok, "call %d: body shape", i)
		require.Equal(t, expectedSkip, body["skip"], "call %d: skip", i)
		require.Equal(t, int32(100), body["top"], "call %d: top", i)
		require.Equal(t, "com.microsoft:get-message-sessions",
			link.calls[i].ApplicationProperties["operation"])
	}
}

func TestClient_ListSessionsForQueue_StopsOnEmptyFirstPage(t *testing.T) {
	link := &scriptedRPCLink{
		t: t,
		responses: []*amqpwrap.RPCResponse{
			{
				Code: 200,
				Message: &amqp.Message{
					Value: map[string]any{
						"sessions-ids": []any{},
					},
				},
			},
		},
	}
	client, _ := newClientForListSessionsUnitTest(t, link)

	got := drainPager(t, client.NewListSessionsForQueuePager("myqueue", nil))
	require.Empty(t, got)
	require.Len(t, link.calls, 1)
}

func TestClient_ListSessionsForQueue_ActiveModeSendsSentinel(t *testing.T) {
	// Active-messages mode (SessionStateUpdatedAfter nil) must send the 10000-01-01 sentinel
	// (253402300800000 ms on the AMQP wire) so the service's .NET AMQP decoder
	// clamps it to DateTime.MaxValue, triggering active-messages mode.
	link := &scriptedRPCLink{
		t: t,
		responses: []*amqpwrap.RPCResponse{
			okPage(t, "active-1"),
		},
	}
	client, _ := newClientForListSessionsUnitTest(t, link)

	_ = drainPager(t, client.NewListSessionsForQueuePager("myqueue", nil))
	require.Len(t, link.calls, 1)

	body := link.calls[0].Value.(map[string]any)
	ts, ok := body["last-updated-time"].(time.Time)
	require.True(t, ok)
	require.Equal(t, 10000, ts.Year())
}

func TestClient_ListSessionsForQueue_SessionStateUpdatedAfterIsPropagated(t *testing.T) {
	sessionStateUpdatedAfter := time.Date(2026, 3, 15, 10, 30, 0, 0, time.UTC)

	link := &scriptedRPCLink{
		t: t,
		responses: []*amqpwrap.RPCResponse{
			okPage(t, "session-after"),
		},
	}
	client, _ := newClientForListSessionsUnitTest(t, link)

	_ = drainPager(t, client.NewListSessionsForQueuePager("myqueue",
		&ListSessionsOptions{SessionStateUpdatedAfter: &sessionStateUpdatedAfter}))
	require.Len(t, link.calls, 1)

	body := link.calls[0].Value.(map[string]any)
	ts, ok := body["last-updated-time"].(time.Time)
	require.True(t, ok)
	require.True(t, ts.Equal(sessionStateUpdatedAfter), "expected %v, got %v", sessionStateUpdatedAfter, ts)
}

func TestClient_ListSessionsForQueue_RecoversFromConnectionError(t *testing.T) {
	// A connection-level error (io.EOF => RecoveryKindConn) on the first RPC must be
	// recovered: the fetch retries, calls Recover with the connection revision captured
	// before the failure, then succeeds. Without the Recover call, NewRPCLink would keep
	// handing back the dead cached connection and the fetch would exhaust its retries.
	link := &scriptedRPCLink{
		t:    t,
		errs: []error{io.EOF, nil},
		responses: []*amqpwrap.RPCResponse{
			nil, // unused: call #0 returns io.EOF
			okPage(t, "recovered-1", "recovered-2"),
		},
	}
	client, ns := newClientForListSessionsUnitTest(t, link)

	got := drainPager(t, client.NewListSessionsForQueuePager("myqueue", nil))
	require.Equal(t, []string{"recovered-1", "recovered-2"}, got)
	require.Len(t, link.calls, 2, "expected one failed then one successful RPC")

	// Recover must have been called once, with the connection revision captured on the
	// failed attempt (FakeNS hands out revision 100 before any recovery).
	require.Equal(t, []uint64{100}, ns.RecoverCalls())
}

func TestClient_ListSessionsForQueue_EmptyNameReturnsError(t *testing.T) {
	// An empty queue name is reported from the first NextPage, before any RPC call,
	// so callers get a clear error instead of a malformed management address like
	// "/$management".
	link := &scriptedRPCLink{
		t: t,
		// No responses scripted: any RPC call would fail the test.
	}
	client, ns := newClientForListSessionsUnitTest(t, link)

	pager := client.NewListSessionsForQueuePager("", nil)
	page, err := pager.NextPage(context.Background())
	require.Error(t, err)
	require.Empty(t, page.Sessions)
	// Per the runtime.Pager contract the caller stops on the returned error; the
	// error must arrive before any management link is created.
	require.Empty(t, link.calls)
	require.Empty(t, ns.ManagementPaths(), "NewRPCLink must not be called when name validation fails")
}

func TestClient_ListSessionsForQueue_ContextCancellationStopsIteration(t *testing.T) {
	// A context cancelled between pages must surface as a terminal error from the next
	// NextPage (context cancellation is fatal, not retried), so the caller's loop exits.
	link := &scriptedRPCLink{
		t: t,
		responses: []*amqpwrap.RPCResponse{
			okPage(t, makeIDs("c", 100)...), // full first page => More() stays true
		},
	}
	client, _ := newClientForListSessionsUnitTest(t, link)

	ctx, cancel := context.WithCancel(context.Background())
	pager := client.NewListSessionsForQueuePager("myqueue", nil)

	first, err := pager.NextPage(ctx)
	require.NoError(t, err)
	require.Len(t, first.Sessions, 100)
	require.True(t, pager.More(), "a full first page means the pager reports more")

	cancel()

	_, err = pager.NextPage(ctx)
	require.ErrorIs(t, err, context.Canceled)
	require.Len(t, link.calls, 1, "the cancelled fetch must not complete a second RPC")
}

func TestClient_ListSessionsForQueue_RecoversFromConnectionErrorOnLaterPage(t *testing.T) {
	// A connection error on a LATER page (page 2, skip=100) must recover and resume at the
	// SAME skip offset, not restart from 0. This guards pagination state across recovery.
	page1 := makeIDs("p1", 100)
	page2 := makeIDs("p2", 20)

	link := &scriptedRPCLink{
		t:    t,
		errs: []error{nil, io.EOF, nil},
		responses: []*amqpwrap.RPCResponse{
			okPage(t, page1...),
			nil, // unused: call #1 returns io.EOF
			okPage(t, page2...),
		},
	}
	client, ns := newClientForListSessionsUnitTest(t, link)

	got := drainPager(t, client.NewListSessionsForQueuePager("myqueue", nil))
	require.Equal(t, append(append([]string{}, page1...), page2...), got)
	require.Len(t, link.calls, 3, "page1 + failed page2 + recovered page2")

	// Skip must be preserved across recovery: 0, then 100 (fail), then 100 (retry).
	require.Equal(t, int32(0), link.calls[0].Value.(map[string]any)["skip"])
	require.Equal(t, int32(100), link.calls[1].Value.(map[string]any)["skip"])
	require.Equal(t, int32(100), link.calls[2].Value.(map[string]any)["skip"])

	// Recover was called once, with the revision captured on page 2's failed attempt.
	require.Equal(t, []uint64{100}, ns.RecoverCalls())
}

func TestClient_ListSessionsForQueue_ExactPageBoundaryStops(t *testing.T) {
	// When the total is an exact multiple of the page size, the last full page is followed
	// by an empty page that terminates enumeration - the pager must stop there without looping.
	page1 := makeIDs("p1", 100)
	page2 := makeIDs("p2", 100)

	link := &scriptedRPCLink{
		t: t,
		responses: []*amqpwrap.RPCResponse{
			okPage(t, page1...),
			okPage(t, page2...),
			okPage(t), // empty final page => stop
		},
	}
	client, _ := newClientForListSessionsUnitTest(t, link)

	got := drainPager(t, client.NewListSessionsForQueuePager("myqueue", nil))
	require.Len(t, got, 200)
	require.Len(t, link.calls, 3, "two full pages plus one terminating empty page")
	require.Equal(t, int32(200), link.calls[2].Value.(map[string]any)["skip"])
}

func TestClient_ListSessionsForQueue_FatalErrorSurfacesWithoutRetry(t *testing.T) {
	// A non-retriable (fatal) error must surface immediately, with no retry attempt.
	link := &scriptedRPCLink{
		t:         t,
		errs:      []error{internal.NewErrNonRetriable("unauthorized")},
		responses: []*amqpwrap.RPCResponse{nil},
	}
	client, _ := newClientForListSessionsUnitTest(t, link)

	pager := client.NewListSessionsForQueuePager("myqueue", nil)
	_, err := pager.NextPage(context.Background())
	require.Error(t, err)
	require.Len(t, link.calls, 1, "a fatal error must not be retried")
}

func TestClient_ListSessionsForSubscription_PaginatesAndSendsCorrectPath(t *testing.T) {
	page1 := makeIDs("s1", 100)
	page2 := makeIDs("s2", 25)

	link := &scriptedRPCLink{
		t: t,
		responses: []*amqpwrap.RPCResponse{
			okPage(t, page1...),
			okPage(t, page2...),
		},
	}
	client, ns := newClientForListSessionsUnitTest(t, link)

	got := drainPager(t, client.NewListSessionsForSubscriptionPager("mytopic", "mysub", nil))
	require.Equal(t, append(append([]string{}, page1...), page2...), got)
	require.Len(t, link.calls, 2)
	require.Equal(t, int32(0), link.calls[0].Value.(map[string]any)["skip"])
	require.Equal(t, int32(100), link.calls[1].Value.(map[string]any)["skip"])
	require.Equal(t, []string{
		"mytopic/Subscriptions/mysub/$management",
		"mytopic/Subscriptions/mysub/$management",
	}, ns.ManagementPaths(),
		"subscription path must follow the topic/Subscriptions/<name> shape")
}

func TestClient_ListSessionsForSubscription_EmptyNamesReturnError(t *testing.T) {
	cases := []struct {
		name  string
		topic string
		sub   string
	}{
		{"empty topic", "", "mysub"},
		{"empty subscription", "mytopic", ""},
		{"both empty", "", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			link := &scriptedRPCLink{t: t}
			client, ns := newClientForListSessionsUnitTest(t, link)

			pager := client.NewListSessionsForSubscriptionPager(tc.topic, tc.sub, nil)
			page, err := pager.NextPage(context.Background())
			require.Error(t, err)
			require.Empty(t, page.Sessions)
			require.Empty(t, link.calls)
			require.Empty(t, ns.ManagementPaths(),
				"NewRPCLink must not be called when name validation fails")
		})
	}
}
