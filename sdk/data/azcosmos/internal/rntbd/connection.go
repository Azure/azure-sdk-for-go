// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

// -----------------------------------------------------------------------------
// Connection Options
// -----------------------------------------------------------------------------

// ConnectionOptions configures a Connection.
type ConnectionOptions struct {
	// TLSConfig is the TLS configuration. If nil, a default config is used.
	TLSConfig *tls.Config

	// ConnectTimeout is the timeout for establishing the TCP connection.
	// Default: 5 seconds.
	ConnectTimeout time.Duration

	// IdleTimeout is the duration after which an idle connection is closed.
	// Default: 0 (use server-provided value from context negotiation).
	IdleTimeout time.Duration

	// RequestTimeout is the default timeout for individual requests.
	// Default: 60 seconds.
	RequestTimeout time.Duration

	// UserAgent is the user agent string sent in context negotiation.
	// Default: "azcosmos-go/2.0.0".
	UserAgent string

	// ReadBufferSize is the size of the read buffer.
	// Default: 64KB.
	ReadBufferSize int

	// WriteBufferSize is the size of the write buffer.
	// Default: 64KB.
	WriteBufferSize int
}

// DefaultConnectionOptions returns the default connection options.
func DefaultConnectionOptions() *ConnectionOptions {
	return &ConnectionOptions{
		ConnectTimeout:  5 * time.Second,
		IdleTimeout:     0, // Use server value
		RequestTimeout:  60 * time.Second,
		UserAgent:       "azcosmos-go/" + ClientVersion,
		ReadBufferSize:  64 * 1024,
		WriteBufferSize: 64 * 1024,
	}
}

// -----------------------------------------------------------------------------
// Pending Request
// -----------------------------------------------------------------------------

// pendingRequest tracks an in-flight request awaiting a response.
type pendingRequest struct {
	transportRequestID uint32
	activityID         uuid.UUID
	responseCh         chan *ResponseMessage
	errCh              chan error
	deadline           time.Time
	sentAt             time.Time
}

// -----------------------------------------------------------------------------
// Connection
// -----------------------------------------------------------------------------

// Connection represents a single RNTBD connection to a Cosmos DB endpoint.
// It supports multiplexing multiple requests over the same TCP connection.
//
// The connection lifecycle:
//  1. Dial() establishes TCP + TLS connection
//  2. Context negotiation handshake
//  3. Ready for request/response multiplexing
//  4. Close() terminates the connection
type Connection struct {
	// Network connection
	conn    net.Conn
	reader  *bufio.Reader
	writer  *bufio.Writer
	writeMu sync.Mutex // Serialize writes

	// Connection state
	context   *Context // Negotiated context (nil until negotiation completes)
	address   *url.URL // Remote address
	options   *ConnectionOptions
	closed    atomic.Bool
	closeCh   chan struct{}
	closeOnce sync.Once
	closeErr  error

	// Request multiplexing
	pending            map[uint32]*pendingRequest
	pendingMu          sync.RWMutex
	nextTransportReqID atomic.Uint32

	// Statistics
	requestCount  atomic.Int64
	responseCount atomic.Int64
	errorCount    atomic.Int64
	lastUsed      atomic.Int64 // Unix timestamp in nanoseconds
}

// -----------------------------------------------------------------------------
// Connection Lifecycle
// -----------------------------------------------------------------------------

// Dial establishes a new RNTBD connection to the specified address.
// The address should be a URL like "rntbd://host:port/".
//
// The connection process:
//  1. Resolve address and establish TCP connection
//  2. TLS handshake
//  3. Context negotiation (RNTBD protocol handshake)
//
// Returns an error if any step fails.
func Dial(ctx context.Context, address *url.URL, opts *ConnectionOptions) (*Connection, error) {
	if opts == nil {
		opts = DefaultConnectionOptions()
	}

	// Validate address
	if address == nil {
		return nil, errors.New("rntbd: address is required")
	}
	if address.Scheme != "rntbd" {
		return nil, fmt.Errorf("rntbd: unsupported scheme %q, expected 'rntbd'", address.Scheme)
	}
	host := address.Host
	if host == "" {
		return nil, errors.New("rntbd: host is required")
	}

	// Apply connect timeout
	connectCtx := ctx
	if opts.ConnectTimeout > 0 {
		var cancel context.CancelFunc
		connectCtx, cancel = context.WithTimeout(ctx, opts.ConnectTimeout)
		defer cancel()
	}

	// Establish TCP connection
	dialer := &net.Dialer{}
	tcpConn, err := dialer.DialContext(connectCtx, "tcp", host)
	if err != nil {
		return nil, fmt.Errorf("rntbd: failed to connect to %s: %w", host, err)
	}

	// TLS handshake
	tlsConfig := opts.TLSConfig
	if tlsConfig == nil {
		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: address.Hostname(),
		}
	}
	tlsConn := tls.Client(tcpConn, tlsConfig)
	if err := tlsConn.HandshakeContext(connectCtx); err != nil {
		_ = tcpConn.Close()
		return nil, fmt.Errorf("rntbd: TLS handshake failed: %w", err)
	}

	// Create connection
	c := &Connection{
		conn:    tlsConn,
		reader:  bufio.NewReaderSize(tlsConn, opts.ReadBufferSize),
		writer:  bufio.NewWriterSize(tlsConn, opts.WriteBufferSize),
		address: address,
		options: opts,
		closeCh: make(chan struct{}),
		pending: make(map[uint32]*pendingRequest),
	}
	c.lastUsed.Store(time.Now().UnixNano())

	// Context negotiation
	if err := c.negotiateContext(connectCtx); err != nil {
		_ = c.conn.Close()
		return nil, fmt.Errorf("rntbd: context negotiation failed: %w", err)
	}

	// Start response reader goroutine
	go c.readLoop()

	return c, nil
}

// negotiateContext performs the RNTBD context negotiation handshake.
func (c *Connection) negotiateContext(ctx context.Context) error {
	// Create context request
	activityID := uuid.New()
	req := NewContextRequest(activityID, c.options.UserAgent)

	// Encode and send
	data, err := req.EncodeToBytes()
	if err != nil {
		return fmt.Errorf("failed to encode context request: %w", err)
	}

	c.writeMu.Lock()
	_, err = c.writer.Write(data)
	if err == nil {
		err = c.writer.Flush()
	}
	c.writeMu.Unlock()
	if err != nil {
		return fmt.Errorf("failed to send context request: %w", err)
	}

	// Read context response
	// Apply deadline if context has one
	if deadline, ok := ctx.Deadline(); ok {
		_ = c.conn.SetReadDeadline(deadline)
		defer func() { _ = c.conn.SetReadDeadline(time.Time{}) }()
	}

	context, err := DecodeContext(c.reader)
	if err != nil {
		return err
	}

	// Verify activity ID matches
	if context.ActivityID != activityID {
		return fmt.Errorf("activity ID mismatch: expected %s, got %s", activityID, context.ActivityID)
	}

	c.context = context
	return nil
}

// Close closes the connection.
func (c *Connection) Close() error {
	c.closeOnce.Do(func() {
		c.closed.Store(true)
		close(c.closeCh)

		// Fail all pending requests
		c.pendingMu.Lock()
		for _, req := range c.pending {
			select {
			case req.errCh <- ErrConnectionClosed:
			default:
			}
		}
		c.pending = make(map[uint32]*pendingRequest)
		c.pendingMu.Unlock()

		c.closeErr = c.conn.Close()
	})
	return c.closeErr
}

// IsClosed returns true if the connection is closed.
func (c *Connection) IsClosed() bool {
	return c.closed.Load()
}

// Context returns the negotiated context, or nil if not yet negotiated.
func (c *Connection) Context() *Context {
	return c.context
}

// Address returns the remote address.
func (c *Connection) Address() *url.URL {
	return c.address
}

// -----------------------------------------------------------------------------
// Request/Response
// -----------------------------------------------------------------------------

// Send sends a request and waits for the response.
// A monotonic transportRequestId is assigned to each request for response
// dispatching, matching the Java SDK's RntbdRequestManager pattern.
//
// If ctx has a deadline, it is used as the request timeout.
// Otherwise, the connection's default RequestTimeout is used.
func (c *Connection) Send(ctx context.Context, req *RequestMessage) (*ResponseMessage, error) {
	if c.IsClosed() {
		return nil, ErrConnectionClosed
	}

	activityID := req.Frame.ActivityID
	if activityID == uuid.Nil {
		return nil, errors.New("rntbd: request must have an ActivityID")
	}

	// Assign a monotonic transport request ID for response dispatching
	transportID := c.nextTransportReqID.Add(1)
	_ = req.Headers.SetValue(uint16(RequestHeaderTransportRequestID), TokenULong, transportID)

	// Determine deadline
	deadline, hasDeadline := ctx.Deadline()
	if !hasDeadline && c.options.RequestTimeout > 0 {
		deadline = time.Now().Add(c.options.RequestTimeout)
		hasDeadline = true
	}

	// Register pending request
	pending := &pendingRequest{
		transportRequestID: transportID,
		activityID:         activityID,
		responseCh:         make(chan *ResponseMessage, 1),
		errCh:              make(chan error, 1),
		deadline:           deadline,
		sentAt:             time.Now(),
	}

	c.pendingMu.Lock()
	c.pending[transportID] = pending
	c.pendingMu.Unlock()

	// Ensure cleanup
	defer func() {
		c.pendingMu.Lock()
		delete(c.pending, transportID)
		c.pendingMu.Unlock()
	}()

	// Encode request
	data, err := EncodeRequestToBytes(req)
	if err != nil {
		c.errorCount.Add(1)
		return nil, fmt.Errorf("rntbd: failed to encode request: %w", err)
	}

	// Send request
	c.writeMu.Lock()
	if hasDeadline {
		_ = c.conn.SetWriteDeadline(deadline)
	}
	_, err = c.writer.Write(data)
	if err == nil {
		err = c.writer.Flush()
	}
	if hasDeadline {
		_ = c.conn.SetWriteDeadline(time.Time{})
	}
	c.writeMu.Unlock()

	if err != nil {
		c.errorCount.Add(1)
		return nil, fmt.Errorf("rntbd: failed to send request: %w", err)
	}

	c.requestCount.Add(1)
	c.lastUsed.Store(time.Now().UnixNano())

	// Wait for response
	var timeoutCh <-chan time.Time
	if hasDeadline {
		timer := time.NewTimer(time.Until(deadline))
		defer timer.Stop()
		timeoutCh = timer.C
	}

	select {
	case resp := <-pending.responseCh:
		c.responseCount.Add(1)
		c.lastUsed.Store(time.Now().UnixNano())
		return resp, nil
	case err := <-pending.errCh:
		c.errorCount.Add(1)
		return nil, err
	case <-timeoutCh:
		c.errorCount.Add(1)
		return nil, fmt.Errorf("rntbd: request timeout for activity %s", activityID)
	case <-ctx.Done():
		c.errorCount.Add(1)
		return nil, ctx.Err()
	case <-c.closeCh:
		return nil, ErrConnectionClosed
	}
}

// -----------------------------------------------------------------------------
// Response Reader
// -----------------------------------------------------------------------------

// readLoop continuously reads responses from the connection and dispatches
// them to the appropriate pending request.
func (c *Connection) readLoop() {
	defer func() {
		if r := recover(); r != nil {
			readErr := fmt.Errorf("rntbd: readLoop panic: %v", r)
			c.failPendingRequests(readErr)
			_ = c.Close()
		}
	}()

	for {
		select {
		case <-c.closeCh:
			return
		default:
		}

		// Read response
		resp, err := DecodeResponseMessage(c.reader)
		if err != nil {
			if c.IsClosed() {
				return
			}
			// Connection error - close and fail all pending requests
			c.failPendingRequests(fmt.Errorf("rntbd: read error: %w", err))
			_ = c.Close()
			return
		}

		// Dispatch to pending request by transportRequestId
		var transportID uint32
		if token := resp.Headers.Get(uint16(ResponseHeaderTransportRequestID)); token != nil && token.IsPresent() {
			if val, err := token.GetValue(); err == nil {
				if v, ok := val.(int64); ok {
					transportID = uint32(v)
				}
			}
		}

		if transportID == 0 {
			// Response missing transportRequestId — cannot dispatch
			continue
		}

		c.pendingMu.RLock()
		pending, ok := c.pending[transportID]
		c.pendingMu.RUnlock()

		if ok {
			select {
			case pending.responseCh <- resp:
			default:
				// Channel full, response dropped (shouldn't happen with buffered channel)
			}
		}
		// If no pending request found, response is discarded (late response after timeout)
	}
}

// failPendingRequests sends the given error to all pending requests.
func (c *Connection) failPendingRequests(err error) {
	c.pendingMu.Lock()
	for _, req := range c.pending {
		select {
		case req.errCh <- err:
		default:
		}
	}
	c.pendingMu.Unlock()
}

// -----------------------------------------------------------------------------
// Statistics
// -----------------------------------------------------------------------------

// Stats returns connection statistics.
type ConnectionStats struct {
	RequestCount  int64
	ResponseCount int64
	ErrorCount    int64
	LastUsed      time.Time
	IdleTime      time.Duration
}

// Stats returns the connection statistics.
func (c *Connection) Stats() ConnectionStats {
	lastUsed := time.Unix(0, c.lastUsed.Load())
	return ConnectionStats{
		RequestCount:  c.requestCount.Load(),
		ResponseCount: c.responseCount.Load(),
		ErrorCount:    c.errorCount.Load(),
		LastUsed:      lastUsed,
		IdleTime:      time.Since(lastUsed),
	}
}

// IsIdle returns true if the connection has been idle longer than the given duration.
func (c *Connection) IsIdle(threshold time.Duration) bool {
	return c.Stats().IdleTime > threshold
}

// PendingRequests returns the number of pending requests.
func (c *Connection) PendingRequests() int {
	c.pendingMu.RLock()
	defer c.pendingMu.RUnlock()
	return len(c.pending)
}

// -----------------------------------------------------------------------------
// Errors
// -----------------------------------------------------------------------------

// ErrConnectionClosed is returned when an operation is attempted on a closed connection.
var ErrConnectionClosed = errors.New("rntbd: connection closed")
