// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------------
// Test Helpers - Mock Server
// -----------------------------------------------------------------------------

// mockServer is a simple mock RNTBD server for testing.
type mockServer struct {
	listener    net.Listener
	tlsConfig   *tls.Config
	connections []net.Conn
	connMu      sync.Mutex
	closed      bool
}

// newMockServer creates a new mock server with TLS.
func newMockServer(t *testing.T) *mockServer {
	// Generate self-signed cert for testing
	cert, err := tls.X509KeyPair(testCert, testKey)
	require.NoError(t, err)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	listener, err := tls.Listen("tcp", "127.0.0.1:0", tlsConfig)
	require.NoError(t, err)

	return &mockServer{
		listener:  listener,
		tlsConfig: tlsConfig,
	}
}

// Address returns the server's address as a URL.
func (s *mockServer) Address() *url.URL {
	addr := s.listener.Addr().String()
	return &url.URL{
		Scheme: "rntbd",
		Host:   addr,
	}
}

// Accept accepts a connection and handles context negotiation.
func (s *mockServer) Accept(t *testing.T) net.Conn {
	conn, err := s.listener.Accept()
	require.NoError(t, err)

	s.connMu.Lock()
	s.connections = append(s.connections, conn)
	s.connMu.Unlock()

	return conn
}

// SafeAccept accepts a connection without using t.Fatal, safe for goroutines.
// Returns nil if the listener is closed or an error occurs.
func (s *mockServer) SafeAccept() net.Conn {
	conn, err := s.listener.Accept()
	if err != nil {
		return nil
	}

	s.connMu.Lock()
	s.connections = append(s.connections, conn)
	s.connMu.Unlock()

	return conn
}

// SafeHandleContextNegotiation handles context negotiation without using t.Fatal.
// Returns true on success, false on failure.
func (s *mockServer) SafeHandleContextNegotiation(conn net.Conn) bool {
	if conn == nil {
		return false
	}

	// Read context request
	req, err := DecodeContextRequest(conn)
	if err != nil {
		return false
	}

	// Send context response
	resp := ContextFrom(req, "MockServer", "1.0.0", StatusOK)
	data, err := resp.EncodeToBytes()
	if err != nil {
		return false
	}

	_, err = conn.Write(data)
	return err == nil
}

// HandleContextNegotiation reads a context request and sends a successful response.
func (s *mockServer) HandleContextNegotiation(t *testing.T, conn net.Conn) {
	// Read context request
	req, err := DecodeContextRequest(conn)
	require.NoError(t, err)

	// Send context response
	resp := ContextFrom(req, "MockServer", "1.0.0", StatusOK)
	data, err := resp.EncodeToBytes()
	require.NoError(t, err)

	_, err = conn.Write(data)
	require.NoError(t, err)
}

// HandleRequest reads a request and sends a response.
// It echoes the TransportRequestID from the request to the response, matching server behavior.
func (s *mockServer) HandleRequest(t *testing.T, conn net.Conn, handler func(*RequestMessage) *ResponseMessage) {
	// Read request
	req, err := DecodeRequestMessage(conn)
	require.NoError(t, err)

	// Generate response
	resp := handler(req)

	// Echo TransportRequestID from request to response (server behavior)
	if token := req.Headers.Get(uint16(RequestHeaderTransportRequestID)); token != nil && token.IsPresent() {
		if val, err := token.GetValue(); err == nil {
			_ = resp.Headers.SetValue(uint16(ResponseHeaderTransportRequestID), TokenULong, val)
		}
	}

	// Send response
	data, err := EncodeResponseToBytes(resp)
	require.NoError(t, err)

	_, err = conn.Write(data)
	require.NoError(t, err)
}

// Close closes the server and all connections.
func (s *mockServer) Close() {
	s.connMu.Lock()
	defer s.connMu.Unlock()

	if s.closed {
		return
	}
	s.closed = true

	for _, conn := range s.connections {
		_ = conn.Close()
	}
	_ = s.listener.Close()
}

// -----------------------------------------------------------------------------
// Test Certificates (self-signed for testing)
// -----------------------------------------------------------------------------

var testCert = []byte(`-----BEGIN CERTIFICATE-----
MIIBXjCCAQSgAwIBAgIBATAKBggqhkjOPQQDAjAPMQ0wCwYDVQQKEwRUZXN0MB4X
DTI2MDQwNDE4NTgzOVoXDTI3MDQwNDE4NTgzOVowDzENMAsGA1UEChMEVGVzdDBZ
MBMGByqGSM49AgEGCCqGSM49AwEHA0IABKBRMTB07d/30SvpIc/R5iFIAFoh07id
qujdbeRAbRiV2+gPEEV5O4ajAbg9FxUTGA6m7RzcHHuDZY6ddJzS7GWjUTBPMA4G
A1UdDwEB/wQEAwIFoDATBgNVHSUEDDAKBggrBgEFBQcDATAMBgNVHRMBAf8EAjAA
MBoGA1UdEQQTMBGCCWxvY2FsaG9zdIcEfwAAATAKBggqhkjOPQQDAgNIADBFAiBy
cYDdfKjN04OyUtH2pqDET1QFFwxU1+2vzZ4Okf5KGAIhAIsBEjeKqZb4kb+SgTUC
EBn63lLhi0B0FUA3bwTpO94X
-----END CERTIFICATE-----
`)

var testKey = []byte(`-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIPoceTiuL0Ck5+hZYxINB+tZg++i4Z32qYlq0rq5VGUXoAoGCCqGSM49
AwEHoUQDQgAEoFExMHTt3/fRK+khz9HmIUgAWiHTuJ2q6N1t5EBtGJXb6A8QRXk7
hqMBuD0XFRMYDqbtHNwce4Nljp10nNLsZQ==
-----END EC PRIVATE KEY-----
`)

// -----------------------------------------------------------------------------
// Connection Tests
// -----------------------------------------------------------------------------

func TestDefaultConnectionOptions(t *testing.T) {
	opts := DefaultConnectionOptions()

	require.Equal(t, 5*time.Second, opts.ConnectTimeout)
	require.Equal(t, 60*time.Second, opts.RequestTimeout)
	require.Equal(t, "azcosmos-go/"+ClientVersion, opts.UserAgent)
	require.Equal(t, 64*1024, opts.ReadBufferSize)
	require.Equal(t, 64*1024, opts.WriteBufferSize)
}

func TestDial_InvalidAddress(t *testing.T) {
	ctx := context.Background()

	// Nil address
	_, err := Dial(ctx, nil, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "address is required")

	// Wrong scheme
	addr := &url.URL{Scheme: "https", Host: "localhost:443"}
	_, err = Dial(ctx, addr, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported scheme")

	// Missing host
	addr = &url.URL{Scheme: "rntbd"}
	_, err = Dial(ctx, addr, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "host is required")
}

func TestDial_ConnectTimeout(t *testing.T) {
	// Try to connect to a non-routable address (should timeout)
	ctx := context.Background()
	addr := &url.URL{Scheme: "rntbd", Host: "10.255.255.1:12345"}
	opts := &ConnectionOptions{
		ConnectTimeout: 100 * time.Millisecond,
	}

	start := time.Now()
	_, err := Dial(ctx, addr, opts)
	elapsed := time.Since(start)

	require.Error(t, err)
	require.Less(t, elapsed, 2*time.Second, "should timeout quickly")
}

func TestDial_Success(t *testing.T) {
	server := newMockServer(t)
	defer server.Close()

	// Start server handler in goroutine
	go func() {
		conn := server.Accept(t)
		server.HandleContextNegotiation(t, conn)
	}()

	// Connect
	ctx := context.Background()
	opts := DefaultConnectionOptions()
	opts.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	conn, err := Dial(ctx, server.Address(), opts)
	require.NoError(t, err)
	require.NotNil(t, conn)
	defer conn.Close()

	// Verify connection state
	require.False(t, conn.IsClosed())
	require.NotNil(t, conn.Context())
	require.Equal(t, "MockServer", conn.Context().ServerAgent)
	require.Equal(t, "1.0.0", conn.Context().ServerVersion)
}

func TestConnection_Send_Success(t *testing.T) {
	server := newMockServer(t)
	defer server.Close()

	var serverConn net.Conn

	// Start server handler
	go func() {
		serverConn = server.Accept(t)
		server.HandleContextNegotiation(t, serverConn)

		// Handle one request
		server.HandleRequest(t, serverConn, func(req *RequestMessage) *ResponseMessage {
			resp := NewResponseMessage(StatusOK, req.Frame.ActivityID)
			resp.Payload = []byte(`{"result": "success"}`)
			resp.Headers.SetValue(uint16(ResponseHeaderPayloadPresent), TokenByte, byte(1))
			return resp
		})
	}()

	// Connect
	ctx := context.Background()
	opts := DefaultConnectionOptions()
	opts.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	conn, err := Dial(ctx, server.Address(), opts)
	require.NoError(t, err)
	defer conn.Close()

	// Send request
	activityID := uuid.New()
	req := NewRequestMessage(ResourceDocument, OperationRead, activityID)
	req.Headers.SetValue(uint16(RequestHeaderPayloadPresent), TokenByte, byte(0))

	resp, err := conn.Send(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, int32(StatusOK), resp.Frame.Status)
	require.Equal(t, activityID, resp.Frame.ActivityID)
	require.Equal(t, `{"result": "success"}`, string(resp.Payload))

	// Check stats
	stats := conn.Stats()
	require.Equal(t, int64(1), stats.RequestCount)
	require.Equal(t, int64(1), stats.ResponseCount)
	require.Equal(t, int64(0), stats.ErrorCount)
}

func TestConnection_Send_Multiplexing(t *testing.T) {
	server := newMockServer(t)
	defer server.Close()

	var serverConn net.Conn

	// Start server handler
	go func() {
		serverConn = server.Accept(t)
		server.HandleContextNegotiation(t, serverConn)

		// Handle multiple requests concurrently
		for i := 0; i < 5; i++ {
			server.HandleRequest(t, serverConn, func(req *RequestMessage) *ResponseMessage {
				// Add small delay to simulate processing
				time.Sleep(10 * time.Millisecond)
				return NewResponseMessage(StatusOK, req.Frame.ActivityID)
			})
		}
	}()

	// Connect
	ctx := context.Background()
	opts := DefaultConnectionOptions()
	opts.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	conn, err := Dial(ctx, server.Address(), opts)
	require.NoError(t, err)
	defer conn.Close()

	// Send multiple requests concurrently
	var wg sync.WaitGroup
	results := make(chan error, 5)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			activityID := uuid.New()
			req := NewRequestMessage(ResourceDocument, OperationRead, activityID)
			req.Headers.SetValue(uint16(RequestHeaderPayloadPresent), TokenByte, byte(0))

			resp, err := conn.Send(ctx, req)
			if err != nil {
				results <- err
				return
			}
			if resp.Frame.ActivityID != activityID {
				results <- io.ErrUnexpectedEOF // Wrong activity ID
				return
			}
			results <- nil
		}()
	}

	wg.Wait()
	close(results)

	// Check all requests succeeded
	for err := range results {
		require.NoError(t, err)
	}

	// Check stats
	stats := conn.Stats()
	require.Equal(t, int64(5), stats.RequestCount)
	require.Equal(t, int64(5), stats.ResponseCount)
}

func TestConnection_Send_OnClosedConnection(t *testing.T) {
	server := newMockServer(t)
	defer server.Close()

	// Start server handler
	go func() {
		conn := server.Accept(t)
		server.HandleContextNegotiation(t, conn)
	}()

	// Connect
	ctx := context.Background()
	opts := DefaultConnectionOptions()
	opts.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	conn, err := Dial(ctx, server.Address(), opts)
	require.NoError(t, err)

	// Close connection
	conn.Close()
	require.True(t, conn.IsClosed())

	// Try to send
	req := NewRequestMessage(ResourceDocument, OperationRead, uuid.New())
	_, err = conn.Send(ctx, req)
	require.Error(t, err)
	require.Equal(t, ErrConnectionClosed, err)
}

func TestConnection_Send_MissingActivityID(t *testing.T) {
	server := newMockServer(t)
	defer server.Close()

	// Start server handler
	go func() {
		conn := server.Accept(t)
		server.HandleContextNegotiation(t, conn)
	}()

	// Connect
	ctx := context.Background()
	opts := DefaultConnectionOptions()
	opts.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	conn, err := Dial(ctx, server.Address(), opts)
	require.NoError(t, err)
	defer conn.Close()

	// Send request without activity ID
	req := NewRequestMessage(ResourceDocument, OperationRead, uuid.Nil)
	_, err = conn.Send(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "ActivityID")
}

func TestConnection_Close_FailsPendingRequests(t *testing.T) {
	server := newMockServer(t)
	defer server.Close()

	var serverConn net.Conn

	// Start server handler that doesn't respond
	go func() {
		serverConn = server.Accept(t)
		server.HandleContextNegotiation(t, serverConn)
		// Don't respond to requests - just wait
		time.Sleep(5 * time.Second)
	}()

	// Connect
	ctx := context.Background()
	opts := DefaultConnectionOptions()
	opts.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	opts.RequestTimeout = 5 * time.Second // Long timeout

	conn, err := Dial(ctx, server.Address(), opts)
	require.NoError(t, err)

	// Start a request in background
	errCh := make(chan error, 1)
	go func() {
		req := NewRequestMessage(ResourceDocument, OperationRead, uuid.New())
		_, err := conn.Send(ctx, req)
		errCh <- err
	}()

	// Wait a bit for request to be sent
	time.Sleep(50 * time.Millisecond)

	// Close connection
	conn.Close()

	// Request should fail
	select {
	case err := <-errCh:
		require.Error(t, err)
	case <-time.After(1 * time.Second):
		t.Fatal("request should have failed when connection closed")
	}
}

func TestConnection_Stats(t *testing.T) {
	server := newMockServer(t)
	defer server.Close()

	// Start server handler
	go func() {
		conn := server.Accept(t)
		server.HandleContextNegotiation(t, conn)
	}()

	// Connect
	ctx := context.Background()
	opts := DefaultConnectionOptions()
	opts.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	conn, err := Dial(ctx, server.Address(), opts)
	require.NoError(t, err)
	defer conn.Close()

	// Check initial stats
	stats := conn.Stats()
	require.Equal(t, int64(0), stats.RequestCount)
	require.Equal(t, int64(0), stats.ResponseCount)
	require.Equal(t, int64(0), stats.ErrorCount)
	require.False(t, stats.LastUsed.IsZero())
}

func TestConnection_IsIdle(t *testing.T) {
	server := newMockServer(t)
	defer server.Close()

	// Start server handler
	go func() {
		conn := server.Accept(t)
		server.HandleContextNegotiation(t, conn)
	}()

	// Connect
	ctx := context.Background()
	opts := DefaultConnectionOptions()
	opts.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	conn, err := Dial(ctx, server.Address(), opts)
	require.NoError(t, err)
	defer conn.Close()

	// Just connected - should not be idle
	require.False(t, conn.IsIdle(100*time.Millisecond))

	// Wait and check again
	time.Sleep(150 * time.Millisecond)
	require.True(t, conn.IsIdle(100*time.Millisecond))
}

func TestConnection_PendingRequests(t *testing.T) {
	server := newMockServer(t)
	defer server.Close()

	// Start server handler
	go func() {
		conn := server.Accept(t)
		server.HandleContextNegotiation(t, conn)
		// Don't respond
		time.Sleep(5 * time.Second)
	}()

	// Connect
	ctx := context.Background()
	opts := DefaultConnectionOptions()
	opts.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	opts.RequestTimeout = 5 * time.Second

	conn, err := Dial(ctx, server.Address(), opts)
	require.NoError(t, err)
	defer conn.Close()

	// No pending requests initially
	require.Equal(t, 0, conn.PendingRequests())

	// Start a request in background
	go func() {
		req := NewRequestMessage(ResourceDocument, OperationRead, uuid.New())
		conn.Send(ctx, req)
	}()

	// Wait for request to be registered
	time.Sleep(50 * time.Millisecond)
	require.Equal(t, 1, conn.PendingRequests())
}
