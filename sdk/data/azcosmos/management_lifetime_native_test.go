// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"context"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestManagementCancellationRetainsNativeLifetime(t *testing.T) {
	started, release := make(chan struct{}), make(chan struct{})
	unblock := sync.OnceFunc(func() { close(release) })
	defer unblock()
	var once sync.Once
	server := newMetadataServer(t, func(w http.ResponseWriter, r *http.Request) bool {
		if r.URL.Path != "/dbs" {
			return false
		}
		once.Do(func() { close(started) })
		<-release
		w.WriteHeader(http.StatusCreated)
		_, _ = io.WriteString(w, `{"id":"db","_rid":"db-rid"}`)
		return true
	})
	client, err := NewClientWithKey(t.Context(), server.URL, mustKeyCredential(t), nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	done := make(chan error, 1)
	go func() { _, err := client.CreateDatabase(ctx, DatabaseProperties{ID: "db"}, nil); done <- err }()
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("management request did not reach native transport")
	}
	cancel()
	select {
	case err := <-done:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(5 * time.Second):
		t.Fatal("canceled management operation did not return")
	}
	closed := make(chan error, 1)
	go func() { closed <- client.Close() }()
	select {
	case <-closed:
		t.Fatal("Close freed native state before late management completion")
	case <-time.After(20 * time.Millisecond):
	}
	unblock()
	select {
	case err := <-closed:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Close did not drain late management completion")
	}
}
