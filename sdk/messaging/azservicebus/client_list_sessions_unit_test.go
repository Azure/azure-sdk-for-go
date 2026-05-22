// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
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
	closed    int
}

func (l *scriptedRPCLink) Close(ctx context.Context) error {
	l.closed++
	return nil
}

func (l *scriptedRPCLink) RPC(ctx context.Context, msg *amqp.Message) (*amqpwrap.RPCResponse, error) {
	idx := len(l.calls)
	l.calls = append(l.calls, msg)
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

func TestClient_ListSessionsForQueue_PaginatesUntilShortPage(t *testing.T) {
	// Arrange three full pages (100 each) followed by a final partial page (37).
	// The implementation should issue exactly 4 RPC calls with skip values
	// 0, 100, 200, 300 and stop on the partial page rather than issuing a 5th
	// call.
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

	got, err := client.ListSessionsForQueue(context.Background(), "myqueue", nil)
	require.NoError(t, err)

	expected := append(append(append(append([]string{}, page1...), page2...), page3...), page4...)
	require.Equal(t, expected, got)
	require.Len(t, link.calls, 4)

	// All paginated calls should reuse a single RPC link addressed to the
	// queue's management endpoint.
	require.Equal(t, []string{"myqueue/$management"}, ns.ManagementPaths())

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

	got, err := client.ListSessionsForQueue(context.Background(), "myqueue", nil)
	require.NoError(t, err)
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

	_, err := client.ListSessionsForQueue(context.Background(), "myqueue", nil)
	require.NoError(t, err)
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

	_, err := client.ListSessionsForQueue(context.Background(), "myqueue",
		&ListSessionsOptions{SessionStateUpdatedAfter: &sessionStateUpdatedAfter})
	require.NoError(t, err)
	require.Len(t, link.calls, 1)

	body := link.calls[0].Value.(map[string]any)
	ts, ok := body["last-updated-time"].(time.Time)
	require.True(t, ok)
	require.True(t, ts.Equal(sessionStateUpdatedAfter), "expected %v, got %v", sessionStateUpdatedAfter, ts)
}

func TestClient_ListSessionsForQueue_EmptyNameReturnsError(t *testing.T) {
	// An empty queue name must be rejected client-side, before any RPC call,
	// so callers get a clear error instead of a malformed management address
	// like "/$management".
	link := &scriptedRPCLink{
		t: t,
		// No responses scripted: any RPC call would fail the test.
	}
	client, ns := newClientForListSessionsUnitTest(t, link)

	got, err := client.ListSessionsForQueue(context.Background(), "", nil)
	require.Error(t, err)
	require.Nil(t, got)
	require.Empty(t, link.calls)
	require.Empty(t, ns.ManagementPaths(), "NewRPCLink must not be called when name validation fails")
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

	got, err := client.ListSessionsForSubscription(context.Background(), "mytopic", "mysub", nil)
	require.NoError(t, err)
	require.Equal(t, append(append([]string{}, page1...), page2...), got)
	require.Len(t, link.calls, 2)
	require.Equal(t, int32(0), link.calls[0].Value.(map[string]any)["skip"])
	require.Equal(t, int32(100), link.calls[1].Value.(map[string]any)["skip"])
	require.Equal(t, []string{"mytopic/Subscriptions/mysub/$management"}, ns.ManagementPaths(),
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

			got, err := client.ListSessionsForSubscription(context.Background(), tc.topic, tc.sub, nil)
			require.Error(t, err)
			require.Nil(t, got)
			require.Empty(t, link.calls)
			require.Empty(t, ns.ManagementPaths(),
				"NewRPCLink must not be called when name validation fails")
		})
	}
}
