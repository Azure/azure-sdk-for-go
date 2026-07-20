// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/stretchr/testify/require"
)

// TestClient_ListSessionsForQueue_Live exercises the get-message-sessions operation
// end-to-end against a live namespace. It is live-only: active-messages mode relies on
// the service clamping the far-future sentinel timestamp to DateTime.MaxValue, which
// cannot be reproduced with recorded responses.
func TestClient_ListSessionsForQueue_Live(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("live-only: session listing exercises service-side sentinel clamping and cannot be recorded")
	}

	queue, cleanup := createQueue(t, nil, &admin.QueueProperties{
		RequiresSession: to.Ptr(true),
	})
	defer cleanup()

	client := newServiceBusClientForTest(t, nil)
	defer func() { require.NoError(t, client.Close(context.Background())) }()

	ctx := context.Background()

	// One message to each of 120 distinct sessions forces skip/top pagination across
	// two pages (100 + 20).
	const sessionCount = 120
	want := make([]string, sessionCount)

	sender, err := client.NewSender(queue, nil)
	require.NoError(t, err)

	for i := 0; i < sessionCount; i++ {
		want[i] = fmt.Sprintf("session-%03d", i)
		require.NoError(t, sender.SendMessage(ctx, &Message{
			Body:      []byte("hello"),
			SessionID: &want[i],
		}, nil))
	}
	require.NoError(t, sender.Close(ctx))

	// Active-messages mode (nil options => far-future sentinel): the service must return
	// exactly the sessions that currently have active messages, across all pages.
	var got []string
	pager := client.NewListSessionsForQueuePager(queue, nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		require.NoError(t, err)
		got = append(got, page.Sessions...)
	}

	require.ElementsMatch(t, want, got)
}

// TestClient_ListSessionsForQueue_SessionStateUpdatedAfter_Live exercises the
// SessionStateUpdatedAfter filter (the second listing mode) end-to-end against a live
// namespace. Unlike active-messages mode, this mode passes a real last-updated-time to the
// service, so it can only be validated live. Listing with a cutoff on either side of the
// state update must return opposite sets: a past cutoff lists the state sessions, a future
// cutoff excludes them. That difference proves the service applies the real last-updated-time
// rather than ignoring the filter.
func TestClient_ListSessionsForQueue_SessionStateUpdatedAfter_Live(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("live-only: session-state-updated-time filtering is a service-side behavior and cannot be recorded")
	}

	queue, cleanup := createQueue(t, nil, &admin.QueueProperties{
		RequiresSession: to.Ptr(true),
	})
	defer cleanup()

	client := newServiceBusClientForTest(t, nil)
	defer func() { require.NoError(t, client.Close(context.Background())) }()

	ctx := context.Background()

	// Three sessions that carry only session state (no active messages), plus one session
	// that has an active message but no session state. Mode two must list the state sessions
	// and must not list the message-only session.
	stateSessions := []string{"state-0", "state-1", "state-2"}
	const messageOnlySession = "msg-only-0"

	sender, err := client.NewSender(queue, nil)
	require.NoError(t, err)
	require.NoError(t, sender.SendMessage(ctx, &Message{
		Body:      []byte("hello"),
		SessionID: to.Ptr(messageOnlySession),
	}, nil))
	require.NoError(t, sender.Close(ctx))

	// A cutoff comfortably before the state updates, tolerant of client/service clock skew.
	before := time.Now().Add(-1 * time.Hour)

	for _, sessionID := range stateSessions {
		sessionReceiver, err := client.AcceptSessionForQueue(ctx, queue, sessionID, nil)
		require.NoError(t, err)
		require.NoError(t, sessionReceiver.SetSessionState(ctx, []byte("state"), nil))
		require.NoError(t, sessionReceiver.Close(ctx))
	}

	// A cutoff comfortably after the state updates.
	after := time.Now().Add(1 * time.Hour)

	// A past cutoff lists the three state sessions and excludes the message-only session
	// (which has no session state).
	gotBefore := collectAllSessions(t, ctx, client.NewListSessionsForQueuePager(queue,
		&ListSessionsOptions{SessionStateUpdatedAfter: &before}))
	for _, sessionID := range stateSessions {
		require.Containsf(t, gotBefore, sessionID, "a past cutoff must list state session %s", sessionID)
	}
	require.NotContains(t, gotBefore, messageOnlySession,
		"mode two must not list a session that has a message but no session state")

	// A future cutoff excludes the same sessions, proving the filter is applied rather than
	// ignored (if it were ignored, both cutoffs would return the same set).
	gotAfter := collectAllSessions(t, ctx, client.NewListSessionsForQueuePager(queue,
		&ListSessionsOptions{SessionStateUpdatedAfter: &after}))
	for _, sessionID := range stateSessions {
		require.NotContainsf(t, gotAfter, sessionID, "a future cutoff must not list state session %s", sessionID)
	}
}

// collectAllSessions drains a session pager and returns every session ID across all pages.
func collectAllSessions(t *testing.T, ctx context.Context, pager *runtime.Pager[ListSessionsResponse]) []string {
	t.Helper()
	var got []string
	for pager.More() {
		page, err := pager.NextPage(ctx)
		require.NoError(t, err)
		got = append(got, page.Sessions...)
	}
	return got
}
