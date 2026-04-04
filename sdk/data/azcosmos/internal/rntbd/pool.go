// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

// -----------------------------------------------------------------------------
// Pool Configuration
// -----------------------------------------------------------------------------

// PoolOptions configures connection pool behavior.
type PoolOptions struct {
	// MaxConnectionsPerEndpoint is the maximum number of connections per endpoint.
	// Default: 10.
	MaxConnectionsPerEndpoint int

	// MaxRequestsPerConnection is the maximum number of concurrent requests per connection.
	// Default: 30.
	MaxRequestsPerConnection int

	// ConnectionAcquisitionTimeout is how long to wait when acquiring a connection.
	// Default: 60 seconds.
	ConnectionAcquisitionTimeout time.Duration

	// IdleConnectionTimeout is how long a connection can be idle before it's closed.
	// 0 means use server-provided value.
	// Default: 0.
	IdleConnectionTimeout time.Duration

	// IdleEndpointTimeout is how long an endpoint can be idle before it's evicted.
	// Default: 1 hour.
	IdleEndpointTimeout time.Duration

	// HealthCheckOnAcquire controls whether to run health checks when acquiring a connection.
	// Default: true.
	HealthCheckOnAcquire bool

	// HealthCheckOnRelease controls whether to run health checks when releasing a connection.
	// Default: true.
	HealthCheckOnRelease bool

	// MaxConcurrentRequestsPerEndpoint is the maximum total concurrent requests per endpoint.
	// Used for fail-fast admission control before pool acquire.
	// Default: 0 (no limit; uses MaxConnectionsPerEndpoint * MaxRequestsPerConnection).
	MaxConcurrentRequestsPerEndpoint int

	// ConnectionOptions are passed to new connections.
	ConnectionOptions *ConnectionOptions
}

// DefaultPoolOptions returns the default pool options.
func DefaultPoolOptions() *PoolOptions {
	return &PoolOptions{
		MaxConnectionsPerEndpoint:    10,
		MaxRequestsPerConnection:     30,
		ConnectionAcquisitionTimeout: 60 * time.Second,
		IdleConnectionTimeout:        0, // Use server value
		IdleEndpointTimeout:          time.Hour,
		HealthCheckOnAcquire:         true,
		HealthCheckOnRelease:         true,
		ConnectionOptions:            DefaultConnectionOptions(),
	}
}

// -----------------------------------------------------------------------------
// Health Checker
// -----------------------------------------------------------------------------

// HealthCheckTimestamps tracks timing information for health checks.
type HealthCheckTimestamps struct {
	lastReadNanos         atomic.Int64
	lastWriteNanos        atomic.Int64
	lastWriteAttemptNanos atomic.Int64
}

// NewHealthCheckTimestamps creates a new HealthCheckTimestamps.
func NewHealthCheckTimestamps() *HealthCheckTimestamps {
	ts := &HealthCheckTimestamps{}
	now := time.Now().UnixNano()
	ts.lastReadNanos.Store(now)
	ts.lastWriteNanos.Store(now)
	ts.lastWriteAttemptNanos.Store(now)
	return ts
}

// RecordRead records a successful read.
func (ts *HealthCheckTimestamps) RecordRead() {
	ts.lastReadNanos.Store(time.Now().UnixNano())
}

// RecordWrite records a successful write.
func (ts *HealthCheckTimestamps) RecordWrite() {
	ts.lastWriteNanos.Store(time.Now().UnixNano())
}

// RecordWriteAttempt records a write attempt.
func (ts *HealthCheckTimestamps) RecordWriteAttempt() {
	ts.lastWriteAttemptNanos.Store(time.Now().UnixNano())
}

// Health check thresholds (matching Java SDK)
const (
	// recentReadWindow - if a read happened within this window, the connection is healthy
	recentReadWindow = 1 * time.Second

	// writeHangGracePeriod - grace period for write hang detection
	writeHangGracePeriod = 2 * time.Second

	// readHangGracePeriod - grace period for read hang detection
	// Java uses (45 + 10) = 55s to avoid false negatives under high CPU (e.g. Spark)
	readHangGracePeriod = 55 * time.Second

	// defaultWriteHangDetectionTime - default write delay limit
	defaultWriteHangDetectionTime = 10 * time.Second

	// defaultReadHangDetectionTime - default read delay limit
	defaultReadHangDetectionTime = 65 * time.Second
)

// HealthChecker checks whether a connection is healthy.
type HealthChecker struct {
	writeHangDetectionTime time.Duration
	readHangDetectionTime  time.Duration
	idleConnectionTimeout  time.Duration
}

// NewHealthChecker creates a new HealthChecker with the given options.
func NewHealthChecker(opts *PoolOptions) *HealthChecker {
	hc := &HealthChecker{
		writeHangDetectionTime: defaultWriteHangDetectionTime,
		readHangDetectionTime:  defaultReadHangDetectionTime,
		idleConnectionTimeout:  opts.IdleConnectionTimeout,
	}
	return hc
}

// HealthCheckResult represents the result of a health check.
type HealthCheckResult struct {
	Healthy bool
	Reason  string
}

// IsHealthy checks if a connection is healthy based on its timestamps.
func (hc *HealthChecker) IsHealthy(conn *Connection, ts *HealthCheckTimestamps) HealthCheckResult {
	if conn == nil || conn.IsClosed() {
		return HealthCheckResult{Healthy: false, Reason: "connection is nil or closed"}
	}

	now := time.Now().UnixNano()
	lastRead := ts.lastReadNanos.Load()
	lastWrite := ts.lastWriteNanos.Load()
	lastWriteAttempt := ts.lastWriteAttemptNanos.Load()

	// Shortcut: if a read happened recently, the connection is healthy
	if now-lastRead < int64(recentReadWindow) {
		return HealthCheckResult{Healthy: true, Reason: "recent read"}
	}

	// Check for hung write
	writeDelay := lastWriteAttempt - lastWrite
	writeHangDuration := now - lastWriteAttempt
	if writeDelay > int64(hc.writeHangDetectionTime) && writeHangDuration > int64(writeHangGracePeriod) {
		return HealthCheckResult{
			Healthy: false,
			Reason:  fmt.Sprintf("hung write: writeDelay=%v, writeHangDuration=%v", time.Duration(writeDelay), time.Duration(writeHangDuration)),
		}
	}

	// Check for hung read
	readDelay := lastWrite - lastRead
	readHangDuration := now - lastWrite
	if readDelay > int64(hc.readHangDetectionTime) && readHangDuration > int64(readHangGracePeriod) {
		return HealthCheckResult{
			Healthy: false,
			Reason:  fmt.Sprintf("hung read: readDelay=%v, readHangDuration=%v", time.Duration(readDelay), time.Duration(readHangDuration)),
		}
	}

	// Check idle timeout
	if hc.idleConnectionTimeout > 0 && now-lastRead > int64(hc.idleConnectionTimeout) {
		return HealthCheckResult{
			Healthy: false,
			Reason:  fmt.Sprintf("idle timeout: idle=%v, threshold=%v", time.Duration(now-lastRead), hc.idleConnectionTimeout),
		}
	}

	return HealthCheckResult{Healthy: true, Reason: "all checks passed"}
}

// -----------------------------------------------------------------------------
// Pooled Connection
// -----------------------------------------------------------------------------

// PooledConnection wraps a Connection with pool-specific metadata.
type PooledConnection struct {
	*Connection
	pool       *ConnectionPool
	endpoint   *Endpoint
	timestamps *HealthCheckTimestamps
	createdAt  time.Time
	id         uint64
}

// Release returns this connection to the pool and decrements the endpoint's concurrent request counter.
func (pc *PooledConnection) Release() {
	if pc.endpoint != nil {
		pc.endpoint.concurrentRequests.Add(-1)
	}
	if pc.pool != nil {
		pc.pool.Release(pc)
	}
}

// -----------------------------------------------------------------------------
// Connection Pool
// -----------------------------------------------------------------------------

// ConnectionPool manages a pool of connections to a single endpoint.
type ConnectionPool struct {
	// Configuration
	address *url.URL
	options *PoolOptions

	// Pool state (protected by mu)
	mu              sync.Mutex
	availableConns  *list.List                   // List of *PooledConnection (LIFO)
	acquiredConns   map[uint64]*PooledConnection // Map of ID -> *PooledConnection
	connecting      bool                         // True if a connection attempt is in progress
	pendingAcquires *list.List                   // List of pending acquisition requests
	nextConnID      uint64
	closed          bool

	// Health checking
	healthChecker *HealthChecker

	// Statistics
	totalCreated atomic.Int64
	totalClosed  atomic.Int64
	totalFailed  atomic.Int64

	// Lifecycle
	closeCh   chan struct{}
	closeOnce sync.Once
}

// pendingAcquire represents a pending connection acquisition request.
type pendingAcquire struct {
	ctx      context.Context
	resultCh chan acquireResult
	deadline time.Time
}

// acquireResult is the result of an acquisition attempt.
type acquireResult struct {
	conn *PooledConnection
	err  error
}

// NewConnectionPool creates a new connection pool for the given address.
func NewConnectionPool(address *url.URL, opts *PoolOptions) *ConnectionPool {
	if opts == nil {
		opts = DefaultPoolOptions()
	}

	pool := &ConnectionPool{
		address:         address,
		options:         opts,
		availableConns:  list.New(),
		acquiredConns:   make(map[uint64]*PooledConnection),
		pendingAcquires: list.New(),
		healthChecker:   NewHealthChecker(opts),
		closeCh:         make(chan struct{}),
	}

	// Start idle eviction goroutine
	go pool.idleEvictionLoop()

	return pool
}

// Acquire gets a connection from the pool, creating one if necessary.
// The returned connection must be released via Release() when done.
func (p *ConnectionPool) Acquire(ctx context.Context) (*PooledConnection, error) {
	// Check for timeout
	deadline, hasDeadline := ctx.Deadline()
	if !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, p.options.ConnectionAcquisitionTimeout)
		defer cancel()
		deadline, _ = ctx.Deadline()
	}

	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil, ErrPoolClosed
	}

	// Try to get an available connection
	conn := p.tryAcquireAvailable()
	if conn != nil {
		p.mu.Unlock()
		return conn, nil
	}

	// Check if we can create a new connection
	totalConns := p.availableConns.Len() + len(p.acquiredConns)
	if p.connecting {
		totalConns++ // Account for in-progress connection
	}

	if totalConns < p.options.MaxConnectionsPerEndpoint && !p.connecting {
		// Create a new connection
		p.connecting = true
		p.mu.Unlock()
		return p.createConnection(ctx)
	}

	// Pool is at capacity or a connection is being created
	// Queue this request
	pending := &pendingAcquire{
		ctx:      ctx,
		resultCh: make(chan acquireResult, 1),
		deadline: deadline,
	}
	p.pendingAcquires.PushBack(pending)
	p.mu.Unlock()

	// Wait for a connection
	select {
	case result := <-pending.resultCh:
		return result.conn, result.err
	case <-ctx.Done():
		// Remove from pending queue
		p.mu.Lock()
		for e := p.pendingAcquires.Front(); e != nil; e = e.Next() {
			if e.Value.(*pendingAcquire) == pending {
				p.pendingAcquires.Remove(e)
				break
			}
		}
		p.mu.Unlock()
		return nil, ctx.Err()
	case <-p.closeCh:
		return nil, ErrPoolClosed
	}
}

// tryAcquireAvailable tries to get an available, serviceable connection.
// Caller must hold p.mu.
func (p *ConnectionPool) tryAcquireAvailable() *PooledConnection {
	for p.availableConns.Len() > 0 {
		// LIFO: get from back
		elem := p.availableConns.Back()
		conn := elem.Value.(*PooledConnection)
		p.availableConns.Remove(elem)

		// Check if serviceable (not too many pending requests)
		if conn.IsClosed() {
			p.totalClosed.Add(1)
			continue
		}

		// Health check on acquire
		if p.options.HealthCheckOnAcquire {
			result := p.healthChecker.IsHealthy(conn.Connection, conn.timestamps)
			if !result.Healthy {
				conn.Close()
				p.totalClosed.Add(1)
				continue
			}
		}

		// Check pending request limit
		if conn.PendingRequests() >= p.options.MaxRequestsPerConnection {
			// Put back for now
			p.availableConns.PushBack(conn)
			return nil
		}

		// Mark as acquired
		p.acquiredConns[conn.id] = conn
		return conn
	}
	return nil
}

// createConnection creates a new connection.
func (p *ConnectionPool) createConnection(ctx context.Context) (*PooledConnection, error) {
	defer func() {
		p.mu.Lock()
		p.connecting = false
		p.processPendingAcquires()
		p.mu.Unlock()
	}()

	conn, err := Dial(ctx, p.address, p.options.ConnectionOptions)
	if err != nil {
		p.totalFailed.Add(1)
		return nil, fmt.Errorf("rntbd: failed to create connection: %w", err)
	}

	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		conn.Close()
		return nil, ErrPoolClosed
	}

	pc := &PooledConnection{
		Connection: conn,
		pool:       p,
		timestamps: NewHealthCheckTimestamps(),
		createdAt:  time.Now(),
		id:         p.nextConnID,
	}
	p.nextConnID++
	p.acquiredConns[pc.id] = pc
	p.totalCreated.Add(1)
	p.mu.Unlock()

	return pc, nil
}

// Release returns a connection to the pool.
func (p *ConnectionPool) Release(conn *PooledConnection) {
	if conn == nil {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		conn.Close()
		p.totalClosed.Add(1)
		return
	}

	// Remove from acquired
	delete(p.acquiredConns, conn.id)

	// Health check on release
	if p.options.HealthCheckOnRelease && !conn.IsClosed() {
		result := p.healthChecker.IsHealthy(conn.Connection, conn.timestamps)
		if !result.Healthy {
			conn.Close()
			p.totalClosed.Add(1)
			p.processPendingAcquires()
			return
		}
	}

	if conn.IsClosed() {
		p.totalClosed.Add(1)
		p.processPendingAcquires()
		return
	}

	// Return to available pool (LIFO)
	p.availableConns.PushBack(conn)

	// Try to satisfy pending acquisitions
	p.processPendingAcquires()
}

// processPendingAcquires tries to satisfy pending acquisition requests.
// Caller must hold p.mu.
func (p *ConnectionPool) processPendingAcquires() {
	for p.pendingAcquires.Len() > 0 && p.availableConns.Len() > 0 {
		// Get oldest pending request
		pendingElem := p.pendingAcquires.Front()
		pending := pendingElem.Value.(*pendingAcquire)

		// Check if context is still valid
		select {
		case <-pending.ctx.Done():
			p.pendingAcquires.Remove(pendingElem)
			continue
		default:
		}

		// Try to get an available connection
		conn := p.tryAcquireAvailable()
		if conn == nil {
			break
		}

		p.pendingAcquires.Remove(pendingElem)
		pending.resultCh <- acquireResult{conn: conn}
	}
}

// idleEvictionLoop periodically evicts idle connections.
func (p *ConnectionPool) idleEvictionLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.evictIdleConnections()
		case <-p.closeCh:
			return
		}
	}
}

// evictIdleConnections removes idle connections from the pool.
func (p *ConnectionPool) evictIdleConnections() {
	if p.options.IdleConnectionTimeout <= 0 {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return
	}

	// Check available connections for idle timeout
	var toRemove []*list.Element
	for e := p.availableConns.Front(); e != nil; e = e.Next() {
		conn := e.Value.(*PooledConnection)
		if conn.IsIdle(p.options.IdleConnectionTimeout) {
			toRemove = append(toRemove, e)
		}
	}

	for _, e := range toRemove {
		conn := e.Value.(*PooledConnection)
		p.availableConns.Remove(e)
		conn.Close()
		p.totalClosed.Add(1)
	}
}

// Close closes the pool and all connections.
func (p *ConnectionPool) Close() error {
	p.closeOnce.Do(func() {
		p.mu.Lock()
		p.closed = true
		close(p.closeCh)

		// Close all available connections
		for e := p.availableConns.Front(); e != nil; e = e.Next() {
			conn := e.Value.(*PooledConnection)
			conn.Connection.Close()
			p.totalClosed.Add(1)
		}
		p.availableConns.Init()

		// Close all acquired connections
		for _, conn := range p.acquiredConns {
			conn.Connection.Close()
			p.totalClosed.Add(1)
		}
		p.acquiredConns = make(map[uint64]*PooledConnection)

		// Fail all pending acquisitions
		for e := p.pendingAcquires.Front(); e != nil; e = e.Next() {
			pending := e.Value.(*pendingAcquire)
			pending.resultCh <- acquireResult{err: ErrPoolClosed}
		}
		p.pendingAcquires.Init()

		p.mu.Unlock()
	})
	return nil
}

// Stats returns pool statistics.
type PoolStats struct {
	AvailableConnections int
	AcquiredConnections  int
	TotalConnections     int
	PendingAcquisitions  int
	TotalCreated         int64
	TotalClosed          int64
	TotalFailed          int64
}

// Stats returns the current pool statistics.
func (p *ConnectionPool) Stats() PoolStats {
	p.mu.Lock()
	defer p.mu.Unlock()
	return PoolStats{
		AvailableConnections: p.availableConns.Len(),
		AcquiredConnections:  len(p.acquiredConns),
		TotalConnections:     p.availableConns.Len() + len(p.acquiredConns),
		PendingAcquisitions:  p.pendingAcquires.Len(),
		TotalCreated:         p.totalCreated.Load(),
		TotalClosed:          p.totalClosed.Load(),
		TotalFailed:          p.totalFailed.Load(),
	}
}

// IsClosed returns true if the pool is closed.
func (p *ConnectionPool) IsClosed() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.closed
}

// -----------------------------------------------------------------------------
// Endpoint Provider
// -----------------------------------------------------------------------------

// EndpointProvider manages endpoints (connection pools) for multiple addresses.
type EndpointProvider struct {
	mu        sync.RWMutex
	endpoints map[string]*Endpoint // key is address authority (host:port)
	options   *PoolOptions
	closed    bool
	closeCh   chan struct{}

	// Statistics
	evictions atomic.Int64
}

// Endpoint represents a connection pool to a specific Cosmos DB replica.
type Endpoint struct {
	provider *EndpointProvider
	address  *url.URL
	pool     *ConnectionPool

	// Tracking
	id               uint64
	createdAt        time.Time
	lastRequestNanos atomic.Int64

	// Admission control
	concurrentRequests    atomic.Int32
	maxConcurrentRequests int

	// Statistics
	requestCount atomic.Int64
	errorCount   atomic.Int64
	successCount atomic.Int64
}

// NewEndpointProvider creates a new EndpointProvider.
func NewEndpointProvider(opts *PoolOptions) *EndpointProvider {
	if opts == nil {
		opts = DefaultPoolOptions()
	}
	ep := &EndpointProvider{
		endpoints: make(map[string]*Endpoint),
		options:   opts,
		closeCh:   make(chan struct{}),
	}

	// Start idle endpoint eviction goroutine
	go ep.idleEndpointEvictionLoop()

	return ep
}

// GetOrCreate returns the endpoint for the given address, creating one if necessary.
func (ep *EndpointProvider) GetOrCreate(address *url.URL) (*Endpoint, error) {
	if address == nil {
		return nil, errors.New("rntbd: address is required")
	}

	key := address.Host // authority (host:port)

	// Fast path: check if endpoint exists
	ep.mu.RLock()
	if ep.closed {
		ep.mu.RUnlock()
		return nil, ErrProviderClosed
	}
	endpoint, ok := ep.endpoints[key]
	ep.mu.RUnlock()
	if ok {
		return endpoint, nil
	}

	// Slow path: create new endpoint
	ep.mu.Lock()
	defer ep.mu.Unlock()

	if ep.closed {
		return nil, ErrProviderClosed
	}

	// Double-check after acquiring write lock
	endpoint, ok = ep.endpoints[key]
	if ok {
		return endpoint, nil
	}

	// Create new endpoint
	pool := NewConnectionPool(address, ep.options)
	maxConcurrent := ep.options.MaxConcurrentRequestsPerEndpoint
	if maxConcurrent <= 0 {
		maxConcurrent = ep.options.MaxConnectionsPerEndpoint * ep.options.MaxRequestsPerConnection
	}
	endpoint = &Endpoint{
		provider:              ep,
		address:               address,
		pool:                  pool,
		id:                    uint64(len(ep.endpoints)),
		createdAt:             time.Now(),
		maxConcurrentRequests: maxConcurrent,
	}
	endpoint.lastRequestNanos.Store(time.Now().UnixNano())
	ep.endpoints[key] = endpoint

	return endpoint, nil
}

// Get returns the endpoint for the given address, or nil if not found.
func (ep *EndpointProvider) Get(address *url.URL) *Endpoint {
	if address == nil {
		return nil
	}
	key := address.Host

	ep.mu.RLock()
	defer ep.mu.RUnlock()
	return ep.endpoints[key]
}

// List returns all endpoints.
func (ep *EndpointProvider) List() []*Endpoint {
	ep.mu.RLock()
	defer ep.mu.RUnlock()

	result := make([]*Endpoint, 0, len(ep.endpoints))
	for _, endpoint := range ep.endpoints {
		result = append(result, endpoint)
	}
	return result
}

// Evict removes an endpoint from the provider.
func (ep *EndpointProvider) Evict(endpoint *Endpoint) {
	if endpoint == nil {
		return
	}

	ep.mu.Lock()
	defer ep.mu.Unlock()

	key := endpoint.address.Host
	if existing, ok := ep.endpoints[key]; ok && existing == endpoint {
		delete(ep.endpoints, key)
		ep.evictions.Add(1)
	}
}

// Count returns the number of endpoints.
func (ep *EndpointProvider) Count() int {
	ep.mu.RLock()
	defer ep.mu.RUnlock()
	return len(ep.endpoints)
}

// Evictions returns the total number of endpoint evictions.
func (ep *EndpointProvider) Evictions() int64 {
	return ep.evictions.Load()
}

// idleEndpointEvictionLoop periodically evicts idle endpoints.
func (ep *EndpointProvider) idleEndpointEvictionLoop() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ep.evictIdleEndpoints()
		case <-ep.closeCh:
			return
		}
	}
}

// evictIdleEndpoints removes endpoints that have been idle longer than IdleEndpointTimeout.
func (ep *EndpointProvider) evictIdleEndpoints() {
	if ep.options.IdleEndpointTimeout <= 0 {
		return
	}

	ep.mu.Lock()
	defer ep.mu.Unlock()

	if ep.closed {
		return
	}

	now := time.Now().UnixNano()
	threshold := int64(ep.options.IdleEndpointTimeout)

	for key, endpoint := range ep.endpoints {
		lastRequest := endpoint.lastRequestNanos.Load()
		if now-lastRequest > threshold {
			endpoint.pool.Close()
			delete(ep.endpoints, key)
			ep.evictions.Add(1)
		}
	}
}

// Close closes the provider and all endpoints.
func (ep *EndpointProvider) Close() error {
	ep.mu.Lock()
	defer ep.mu.Unlock()

	if ep.closed {
		return nil
	}
	ep.closed = true
	close(ep.closeCh)

	// Close all endpoints
	for _, endpoint := range ep.endpoints {
		endpoint.pool.Close()
	}
	ep.endpoints = make(map[string]*Endpoint)

	return nil
}

// IsClosed returns true if the provider is closed.
func (ep *EndpointProvider) IsClosed() bool {
	ep.mu.RLock()
	defer ep.mu.RUnlock()
	return ep.closed
}

// Options returns the pool options.
func (ep *EndpointProvider) Options() *PoolOptions {
	return ep.options
}

// -----------------------------------------------------------------------------
// Endpoint Methods
// -----------------------------------------------------------------------------

// Acquire gets a connection from this endpoint's pool.
// Performs fail-fast admission control before attempting pool acquire.
func (e *Endpoint) Acquire(ctx context.Context) (*PooledConnection, error) {
	e.lastRequestNanos.Store(time.Now().UnixNano())
	e.requestCount.Add(1)

	// Fail-fast admission control matching Java RntbdServiceEndpoint.request()
	current := e.concurrentRequests.Add(1)
	if int(current) > e.maxConcurrentRequests {
		e.concurrentRequests.Add(-1)
		e.errorCount.Add(1)
		return nil, ErrEndpointOverloaded
	}

	conn, err := e.pool.Acquire(ctx)
	if err != nil {
		e.concurrentRequests.Add(-1)
		e.errorCount.Add(1)
		return nil, err
	}
	conn.endpoint = e
	return conn, nil
}

// Address returns the endpoint address.
func (e *Endpoint) Address() *url.URL {
	return e.address
}

// ID returns the endpoint ID.
func (e *Endpoint) ID() uint64 {
	return e.id
}

// CreatedAt returns when the endpoint was created.
func (e *Endpoint) CreatedAt() time.Time {
	return e.createdAt
}

// LastRequestTime returns the time of the last request.
func (e *Endpoint) LastRequestTime() time.Time {
	return time.Unix(0, e.lastRequestNanos.Load())
}

// RecordSuccess records a successful request.
func (e *Endpoint) RecordSuccess() {
	e.successCount.Add(1)
}

// RecordError records a failed request.
func (e *Endpoint) RecordError() {
	e.errorCount.Add(1)
}

// Stats returns endpoint statistics.
type EndpointStats struct {
	RequestCount int64
	SuccessCount int64
	ErrorCount   int64
	PoolStats    PoolStats
	LastRequest  time.Time
	CreatedAt    time.Time
}

// Stats returns the endpoint statistics.
func (e *Endpoint) Stats() EndpointStats {
	return EndpointStats{
		RequestCount: e.requestCount.Load(),
		SuccessCount: e.successCount.Load(),
		ErrorCount:   e.errorCount.Load(),
		PoolStats:    e.pool.Stats(),
		LastRequest:  e.LastRequestTime(),
		CreatedAt:    e.createdAt,
	}
}

// Close closes this endpoint and removes it from the provider.
func (e *Endpoint) Close() error {
	if e.provider != nil {
		e.provider.Evict(e)
	}
	return e.pool.Close()
}

// IsClosed returns true if the endpoint's pool is closed.
func (e *Endpoint) IsClosed() bool {
	return e.pool.IsClosed()
}

// -----------------------------------------------------------------------------
// Errors
// -----------------------------------------------------------------------------

// ErrPoolClosed is returned when an operation is attempted on a closed pool.
var ErrPoolClosed = errors.New("rntbd: connection pool is closed")

// ErrProviderClosed is returned when an operation is attempted on a closed provider.
var ErrProviderClosed = errors.New("rntbd: endpoint provider is closed")

// ErrAcquisitionTimeout is returned when connection acquisition times out.
var ErrAcquisitionTimeout = errors.New("rntbd: connection acquisition timeout")

// ErrEndpointOverloaded is returned when the endpoint has too many concurrent requests.
var ErrEndpointOverloaded = errors.New("rntbd: endpoint overloaded, too many concurrent requests")
