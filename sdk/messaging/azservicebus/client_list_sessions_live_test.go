// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"fmt"
	"testing"

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
