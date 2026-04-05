// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"context"
	"crypto/tls"
	"net/url"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------------
// Pool Options Tests
// -----------------------------------------------------------------------------

func TestDefaultPoolOptions(t *testing.T) {
	opts := DefaultPoolOptions()

	require.Equal(t, 10, opts.MaxConnectionsPerEndpoint)
	require.Equal(t, 30, opts.MaxRequestsPerConnection)
	require.Equal(t, 60*time.Second, opts.ConnectionAcquisitionTimeout)
	require.Equal(t, time.Duration(0), opts.IdleConnectionTimeout)
	require.Equal(t, time.Hour, opts.IdleEndpointTimeout)
	require.True(t, opts.HealthCheckOnAcquire)
	require.True(t, opts.HealthCheckOnRelease)
	require.NotNil(t, opts.ConnectionOptions)
}

// -----------------------------------------------------------------------------
// Health Checker Tests
// -----------------------------------------------------------------------------

func TestHealthCheckTimestamps(t *testing.T) {
	ts := NewHealthCheckTimestamps()

	// Initial values should be set
	require.Greater(t, ts.lastReadNanos.Load(), int64(0))
	require.Greater(t, ts.lastWriteNanos.Load(), int64(0))
	require.Greater(t, ts.lastWriteAttemptNanos.Load(), int64(0))

	// Record operations
	time.Sleep(10 * time.Millisecond)
	ts.RecordRead()
	newRead := ts.lastReadNanos.Load()
	require.Greater(t, newRead, int64(0))

	time.Sleep(10 * time.Millisecond)
	ts.RecordWrite()
	newWrite := ts.lastWriteNanos.Load()
	require.Greater(t, newWrite, newRead)

	time.Sleep(10 * time.Millisecond)
	ts.RecordWriteAttempt()
	newAttempt := ts.lastWriteAttemptNanos.Load()
	require.Greater(t, newAttempt, newWrite)
}

func TestHealthChecker_RecentRead(t *testing.T) {
	opts := DefaultPoolOptions()
	hc := NewHealthChecker(opts)

	// Create mock connection and timestamps
	server := newMockServer(t)
	defer server.Close()

	go func() {
		conn := server.Accept(t)
		server.HandleContextNegotiation(t, conn)
		time.Sleep(time.Second)
	}()

	ctx := context.Background()
	connOpts := DefaultConnectionOptions()
	connOpts.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	conn, err := Dial(ctx, server.Address(), connOpts)
	require.NoError(t, err)
	defer conn.Close()

	// Recent read should be healthy
	ts := NewHealthCheckTimestamps()
	ts.RecordRead() // Just recorded a read
	result := hc.IsHealthy(conn, ts)
	require.True(t, result.Healthy)
	require.Equal(t, "recent read", result.Reason)
}

func TestHealthChecker_ClosedConnection(t *testing.T) {
	opts := DefaultPoolOptions()
	hc := NewHealthChecker(opts)

	// nil connection
	result := hc.IsHealthy(nil, NewHealthCheckTimestamps())
	require.False(t, result.Healthy)
	require.Equal(t, "connection is nil or closed", result.Reason)

	// Create and close a connection
	server := newMockServer(t)
	defer server.Close()

	go func() {
		conn := server.Accept(t)
		server.HandleContextNegotiation(t, conn)
	}()

	ctx := context.Background()
	connOpts := DefaultConnectionOptions()
	connOpts.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	conn, err := Dial(ctx, server.Address(), connOpts)
	require.NoError(t, err)
	conn.Close()

	ts := NewHealthCheckTimestamps()
	result = hc.IsHealthy(conn, ts)
	require.False(t, result.Healthy)
	require.Equal(t, "connection is nil or closed", result.Reason)
}

func TestHealthChecker_IdleTimeout(t *testing.T) {
	opts := DefaultPoolOptions()
	opts.IdleConnectionTimeout = 100 * time.Millisecond
	hc := NewHealthChecker(opts)

	// Create mock connection
	server := newMockServer(t)
	defer server.Close()

	go func() {
		conn := server.Accept(t)
		server.HandleContextNegotiation(t, conn)
		time.Sleep(time.Second)
	}()

	ctx := context.Background()
	connOpts := DefaultConnectionOptions()
	connOpts.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	conn, err := Dial(ctx, server.Address(), connOpts)
	require.NoError(t, err)
	defer conn.Close()

	// Create timestamps that simulate idle connection
	// Must be older than both:
	// 1. recentReadWindow (1s) - so we don't short-circuit to healthy
	// 2. IdleConnectionTimeout (100ms) - so we trigger idle timeout
	ts := NewHealthCheckTimestamps()
	ts.lastReadNanos.Store(time.Now().Add(-2 * time.Second).UnixNano())
	ts.lastWriteNanos.Store(time.Now().Add(-2 * time.Second).UnixNano())
	ts.lastWriteAttemptNanos.Store(time.Now().Add(-2 * time.Second).UnixNano())

	result := hc.IsHealthy(conn, ts)
	require.False(t, result.Healthy)
	require.Contains(t, result.Reason, "idle timeout")
}

// -----------------------------------------------------------------------------
// Connection Pool Tests
// -----------------------------------------------------------------------------

func TestConnectionPool_Create(t *testing.T) {
	addr, _ := url.Parse("rntbd://localhost:10255/")
	opts := DefaultPoolOptions()

	pool := NewConnectionPool(addr, opts)
	require.NotNil(t, pool)
	defer pool.Close()

	require.False(t, pool.IsClosed())
	stats := pool.Stats()
	require.Equal(t, 0, stats.AvailableConnections)
	require.Equal(t, 0, stats.AcquiredConnections)
	require.Equal(t, 0, stats.TotalConnections)
}

func TestConnectionPool_AcquireAndRelease(t *testing.T) {
	server := newMockServer(t)
	defer server.Close()

	// Start server handler for one connection
	go func() {
		conn := server.Accept(t)
		server.HandleContextNegotiation(t, conn)
		time.Sleep(2 * time.Second)
	}()

	opts := DefaultPoolOptions()
	opts.MaxConnectionsPerEndpoint = 5
	opts.ConnectionOptions = DefaultConnectionOptions()
	opts.ConnectionOptions.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	pool := NewConnectionPool(server.Address(), opts)
	defer pool.Close()

	ctx := context.Background()

	// Acquire first connection
	conn1, err := pool.Acquire(ctx)
	require.NoError(t, err)
	require.NotNil(t, conn1)

	stats := pool.Stats()
	require.Equal(t, 0, stats.AvailableConnections)
	require.Equal(t, 1, stats.AcquiredConnections)
	require.Equal(t, int64(1), stats.TotalCreated)

	// Release connection
	conn1.Release()

	stats = pool.Stats()
	require.Equal(t, 1, stats.AvailableConnections)
	require.Equal(t, 0, stats.AcquiredConnections)
}

func TestConnectionPool_ReuseConnection(t *testing.T) {
	server := newMockServer(t)
	defer server.Close()

	// Start server handler
	go func() {
		conn := server.Accept(t)
		server.HandleContextNegotiation(t, conn)
		time.Sleep(2 * time.Second)
	}()

	opts := DefaultPoolOptions()
	opts.MaxConnectionsPerEndpoint = 5
	opts.HealthCheckOnAcquire = false // Disable for this test
	opts.ConnectionOptions = DefaultConnectionOptions()
	opts.ConnectionOptions.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	pool := NewConnectionPool(server.Address(), opts)
	defer pool.Close()

	ctx := context.Background()

	// Acquire and release
	conn1, err := pool.Acquire(ctx)
	require.NoError(t, err)
	conn1ID := conn1.id
	conn1.Release()

	// Acquire again - should get same connection
	conn2, err := pool.Acquire(ctx)
	require.NoError(t, err)
	require.Equal(t, conn1ID, conn2.id)
	conn2.Release()

	// Only one connection should have been created
	stats := pool.Stats()
	require.Equal(t, int64(1), stats.TotalCreated)
}

func TestConnectionPool_MaxConnections(t *testing.T) {
	server := newMockServer(t)
	defer server.Close()

	maxConns := 2

	go func() {
		for i := 0; i < maxConns; i++ {
			conn := server.SafeAccept()
			if conn == nil {
				return
			}
			server.SafeHandleContextNegotiation(conn)
		}
	}()

	opts := DefaultPoolOptions()
	opts.MaxConnectionsPerEndpoint = maxConns
	opts.MaxRequestsPerConnection = 100
	opts.ConnectionAcquisitionTimeout = 500 * time.Millisecond
	opts.ConnectionOptions = DefaultConnectionOptions()
	opts.ConnectionOptions.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	pool := NewConnectionPool(server.Address(), opts)
	defer pool.Close()

	ctx := context.Background()

	var conns []*PooledConnection
	for i := 0; i < maxConns; i++ {
		conn, err := pool.Acquire(ctx)
		require.NoError(t, err)
		conns = append(conns, conn)
	}

	stats := pool.Stats()
	require.Equal(t, maxConns, stats.AcquiredConnections)
	require.Equal(t, int64(maxConns), stats.TotalCreated)

	shortCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	_, err := pool.Acquire(shortCtx)
	require.Error(t, err)
	require.ErrorIs(t, err, context.DeadlineExceeded)

	for _, conn := range conns {
		conn.Release()
	}
}

func TestConnectionPool_Close(t *testing.T) {
	server := newMockServer(t)
	defer server.Close()

	go func() {
		conn := server.SafeAccept()
		if conn != nil {
			server.SafeHandleContextNegotiation(conn)
		}
	}()

	opts := DefaultPoolOptions()
	opts.ConnectionOptions = DefaultConnectionOptions()
	opts.ConnectionOptions.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	pool := NewConnectionPool(server.Address(), opts)

	ctx := context.Background()

	conn, err := pool.Acquire(ctx)
	require.NoError(t, err)

	pool.Close()
	require.True(t, pool.IsClosed())

	conn.Release()

	_, err = pool.Acquire(ctx)
	require.ErrorIs(t, err, ErrPoolClosed)
}

func TestConnectionPool_ConcurrentAcquireRelease(t *testing.T) {
	server := newMockServer(t)
	defer server.Close()

	maxConns := 3
	numWorkers := 10
	opsPerWorker := 20

	go func() {
		for i := 0; i < maxConns+5; i++ {
			conn := server.SafeAccept()
			if conn == nil {
				return
			}
			server.SafeHandleContextNegotiation(conn)
		}
	}()

	opts := DefaultPoolOptions()
	opts.MaxConnectionsPerEndpoint = maxConns
	opts.MaxRequestsPerConnection = 100
	opts.HealthCheckOnAcquire = false
	opts.HealthCheckOnRelease = false
	opts.ConnectionOptions = DefaultConnectionOptions()
	opts.ConnectionOptions.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	pool := NewConnectionPool(server.Address(), opts)
	defer pool.Close()

	var wg sync.WaitGroup
	var successCount atomic.Int64

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx := context.Background()

			for j := 0; j < opsPerWorker; j++ {
				conn, err := pool.Acquire(ctx)
				if err != nil {
					continue
				}
				successCount.Add(1)
				time.Sleep(time.Millisecond)
				conn.Release()
			}
		}()
	}

	wg.Wait()

	require.Greater(t, successCount.Load(), int64(50))

	stats := pool.Stats()
	require.Equal(t, 0, stats.AcquiredConnections)
	require.LessOrEqual(t, stats.AvailableConnections, maxConns)
}

// -----------------------------------------------------------------------------
// Endpoint Provider Tests
// -----------------------------------------------------------------------------

func TestEndpointProvider_Create(t *testing.T) {
	opts := DefaultPoolOptions()
	provider := NewEndpointProvider(opts)
	require.NotNil(t, provider)
	defer provider.Close()

	require.False(t, provider.IsClosed())
	require.Equal(t, 0, provider.Count())
	require.Equal(t, int64(0), provider.Evictions())
}

func TestEndpointProvider_GetOrCreate(t *testing.T) {
	server := newMockServer(t)
	defer server.Close()

	go func() {
		conn := server.SafeAccept()
		if conn != nil {
			server.SafeHandleContextNegotiation(conn)
		}
	}()

	opts := DefaultPoolOptions()
	opts.ConnectionOptions = DefaultConnectionOptions()
	opts.ConnectionOptions.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	provider := NewEndpointProvider(opts)
	defer provider.Close()

	endpoint1, err := provider.GetOrCreate(server.Address())
	require.NoError(t, err)
	require.NotNil(t, endpoint1)
	require.Equal(t, 1, provider.Count())

	endpoint2, err := provider.GetOrCreate(server.Address())
	require.NoError(t, err)
	require.Equal(t, endpoint1, endpoint2)
	require.Equal(t, 1, provider.Count())

	_, err = provider.GetOrCreate(nil)
	require.Error(t, err)
}

func TestEndpointProvider_MultipleEndpoints(t *testing.T) {
	server1 := newMockServer(t)
	defer server1.Close()

	server2 := newMockServer(t)
	defer server2.Close()

	go func() {
		conn := server1.SafeAccept()
		if conn != nil {
			server1.SafeHandleContextNegotiation(conn)
		}
	}()

	go func() {
		conn := server2.SafeAccept()
		if conn != nil {
			server2.SafeHandleContextNegotiation(conn)
		}
	}()

	opts := DefaultPoolOptions()
	opts.ConnectionOptions = DefaultConnectionOptions()
	opts.ConnectionOptions.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	provider := NewEndpointProvider(opts)
	defer provider.Close()

	endpoint1, err := provider.GetOrCreate(server1.Address())
	require.NoError(t, err)

	endpoint2, err := provider.GetOrCreate(server2.Address())
	require.NoError(t, err)

	require.NotEqual(t, endpoint1, endpoint2)
	require.Equal(t, 2, provider.Count())

	endpoints := provider.List()
	require.Len(t, endpoints, 2)
}

func TestEndpointProvider_Evict(t *testing.T) {
	opts := DefaultPoolOptions()
	provider := NewEndpointProvider(opts)
	defer provider.Close()

	addr, _ := url.Parse("rntbd://localhost:10255/")
	endpoint, err := provider.GetOrCreate(addr)
	require.NoError(t, err)
	require.Equal(t, 1, provider.Count())

	// Evict endpoint
	provider.Evict(endpoint)
	require.Equal(t, 0, provider.Count())
	require.Equal(t, int64(1), provider.Evictions())

	// Get returns nil now
	got := provider.Get(addr)
	require.Nil(t, got)

	// Double evict should be safe
	provider.Evict(endpoint)
	require.Equal(t, int64(1), provider.Evictions())
}

func TestEndpointProvider_Close(t *testing.T) {
	server := newMockServer(t)
	defer server.Close()

	go func() {
		conn := server.SafeAccept()
		if conn != nil {
			server.SafeHandleContextNegotiation(conn)
		}
	}()

	opts := DefaultPoolOptions()
	opts.ConnectionOptions = DefaultConnectionOptions()
	opts.ConnectionOptions.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	provider := NewEndpointProvider(opts)

	endpoint, err := provider.GetOrCreate(server.Address())
	require.NoError(t, err)

	provider.Close()
	require.True(t, provider.IsClosed())
	require.True(t, endpoint.IsClosed())

	_, err = provider.GetOrCreate(server.Address())
	require.ErrorIs(t, err, ErrProviderClosed)
}

// -----------------------------------------------------------------------------
// Endpoint Tests
// -----------------------------------------------------------------------------

func TestEndpoint_Acquire(t *testing.T) {
	server := newMockServer(t)
	defer server.Close()

	go func() {
		conn := server.SafeAccept()
		if conn != nil {
			server.SafeHandleContextNegotiation(conn)
		}
	}()

	opts := DefaultPoolOptions()
	opts.ConnectionOptions = DefaultConnectionOptions()
	opts.ConnectionOptions.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	provider := NewEndpointProvider(opts)
	defer provider.Close()

	endpoint, err := provider.GetOrCreate(server.Address())
	require.NoError(t, err)

	ctx := context.Background()

	conn, err := endpoint.Acquire(ctx)
	require.NoError(t, err)
	require.NotNil(t, conn)
	defer conn.Release()

	stats := endpoint.Stats()
	require.Equal(t, int64(1), stats.RequestCount)
}

func TestEndpoint_Stats(t *testing.T) {
	opts := DefaultPoolOptions()
	provider := NewEndpointProvider(opts)
	defer provider.Close()

	addr, _ := url.Parse("rntbd://localhost:10255/")
	endpoint, err := provider.GetOrCreate(addr)
	require.NoError(t, err)

	stats := endpoint.Stats()
	require.Equal(t, int64(0), stats.RequestCount)
	require.Equal(t, int64(0), stats.SuccessCount)
	require.Equal(t, int64(0), stats.ErrorCount)
	require.False(t, stats.CreatedAt.IsZero())
	require.False(t, stats.LastRequest.IsZero())

	// Record events
	endpoint.RecordSuccess()
	endpoint.RecordSuccess()
	endpoint.RecordError()

	stats = endpoint.Stats()
	require.Equal(t, int64(2), stats.SuccessCount)
	require.Equal(t, int64(1), stats.ErrorCount)
}

func TestEndpoint_Close(t *testing.T) {
	opts := DefaultPoolOptions()
	provider := NewEndpointProvider(opts)
	defer provider.Close()

	addr, _ := url.Parse("rntbd://localhost:10255/")
	endpoint, err := provider.GetOrCreate(addr)
	require.NoError(t, err)
	require.Equal(t, 1, provider.Count())

	// Close endpoint
	endpoint.Close()
	require.True(t, endpoint.IsClosed())

	// Should be evicted from provider
	require.Equal(t, 0, provider.Count())
	require.Equal(t, int64(1), provider.Evictions())
}

// -----------------------------------------------------------------------------
// Pooled Connection Tests
// -----------------------------------------------------------------------------

func TestPooledConnection_Release(t *testing.T) {
	server := newMockServer(t)
	defer server.Close()

	go func() {
		conn := server.SafeAccept()
		if conn != nil {
			server.SafeHandleContextNegotiation(conn)
		}
	}()

	opts := DefaultPoolOptions()
	opts.HealthCheckOnAcquire = false
	opts.HealthCheckOnRelease = false
	opts.ConnectionOptions = DefaultConnectionOptions()
	opts.ConnectionOptions.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	pool := NewConnectionPool(server.Address(), opts)
	defer pool.Close()

	ctx := context.Background()

	conn, err := pool.Acquire(ctx)
	require.NoError(t, err)

	stats := pool.Stats()
	require.Equal(t, 1, stats.AcquiredConnections)
	require.Equal(t, 0, stats.AvailableConnections)

	conn.Release()

	stats = pool.Stats()
	require.Equal(t, 0, stats.AcquiredConnections)
	require.Equal(t, 1, stats.AvailableConnections)
}

// -----------------------------------------------------------------------------
// Error Tests
// -----------------------------------------------------------------------------

func TestPoolErrors(t *testing.T) {
	// Test that error types exist and are distinct
	require.NotEqual(t, ErrPoolClosed, ErrProviderClosed)
	require.NotEqual(t, ErrPoolClosed, ErrAcquisitionTimeout)
	require.NotEqual(t, ErrProviderClosed, ErrAcquisitionTimeout)

	// Error messages
	require.Contains(t, ErrPoolClosed.Error(), "pool")
	require.Contains(t, ErrProviderClosed.Error(), "provider")
	require.Contains(t, ErrAcquisitionTimeout.Error(), "timeout")
}
