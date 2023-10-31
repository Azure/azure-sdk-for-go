//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azwebpubsub_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azwebpubsub"
	"github.com/stretchr/testify/require"
)

func TestClient_SendToAll(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode || testing.Short() {
		t.Skip()
	}

	client := newClientWrapper(t)
	const hub = "chat"

	_, err := client.SendToAll(context.Background(),
		hub, azwebpubsub.ContentTypeTextPlain, newStream("Hello world!"),
		&azwebpubsub.ClientSendToAllOptions{})
	require.NoError(t, err)
}

func TestClient_Exists(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode || testing.Short() {
		t.Skip()
	}

	client := newClientWrapper(t)
	const hub = "chat"

	_, err := client.SendToAll(context.Background(),
		hub, azwebpubsub.ContentTypeTextPlain, newStream("Hello world!"),
		&azwebpubsub.ClientSendToAllOptions{})
	require.NoError(t, err)
}
func newStream(message string) io.ReadSeekCloser {
	return streaming.NopCloser(bytes.NewReader([]byte(message)))
}
