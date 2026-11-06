// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"container/list"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"
)

type cosmosEndpoint struct {
	scheme  string
	address string
}

type cosmosEndpointPool struct {
	endpoint cosmosEndpoint
	conns    []*cosmosPoolConn
	dialing  bool
	users    int
	idle     int
	signal   cosmosPoolSignal
}

// The state hook can run inside Reserve or Close, so it must not take the pool
// lock. Wake one waiter; successful acquisitions pass the notification along.
type cosmosPoolSignal struct {
	changed chan struct{}
}

func (s *cosmosPoolSignal) channel() <-chan struct{} {
	return s.changed
}

func (s *cosmosPoolSignal) notify() {
	select {
	case s.changed <- struct{}{}:
	default:
	}
}

type cosmosConnectionPool struct {
	transport *http.Transport
	limit     int
	mu        sync.Mutex
	endpoints map[cosmosEndpoint]*cosmosEndpointPool
	idle      list.List
}

type cosmosPoolConn struct {
	pool        *cosmosConnectionPool
	endpoint    *cosmosEndpointPool
	conn        *http.ClientConn
	socket      *cosmosPoolSocket
	users       int
	used        bool
	ready       bool
	closing     bool
	released    bool
	idle        *list.Element
	idleSince   time.Time
	timer       *time.Timer
	releaseOnce sync.Once
}

type cosmosPoolDialKey struct{}

type cosmosPoolDial struct {
	conn    *cosmosPoolConn
	request *http.Request
}

func newCosmosConnectionPool(transport *http.Transport, limit int) *cosmosConnectionPool {
	p := &cosmosConnectionPool{
		transport: transport,
		limit:     limit,
		endpoints: make(map[cosmosEndpoint]*cosmosEndpointPool),
	}
	dial := transport.DialContext
	transport.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		attempt := ctx.Value(cosmosPoolDialKey{}).(cosmosPoolDial)
		conn, err := dial(ctx, network, address)
		if err != nil {
			return nil, err
		}
		socket := &cosmosPoolSocket{Conn: conn, release: func() {
			// net/http can close a socket while holding its connection lock.
			// Never acquire the pool lock from that callback.
			go attempt.conn.release()
		}}
		p.mu.Lock()
		attempt.conn.socket = socket
		p.mu.Unlock()
		return socket, nil
	}
	proxy := transport.Proxy
	if proxy != nil {
		transport.Proxy = func(req *http.Request) (*url.URL, error) {
			// NewClientConn supplies a synthetic request without the port.
			// Evaluate proxy exclusions using the original request instead.
			attempt := req.Context().Value(cosmosPoolDialKey{}).(cosmosPoolDial)
			return proxy(attempt.request)
		}
	}
	return p
}

func (p *cosmosConnectionPool) RoundTrip(req *http.Request) (_ *http.Response, err error) {
	originalRequest := req
	if req.Body != nil {
		clone := req.Clone(req.Context())
		clone.Body = &cosmosPoolRequestBody{ReadCloser: req.Body}
		req = clone
		defer func() {
			if err != nil {
				_ = req.Body.Close()
			}
		}()
	}
	if req.URL == nil || req.URL.Hostname() == "" {
		return nil, errors.New("azcosmos: HTTP request requires an endpoint")
	}
	scheme := req.URL.Scheme
	if scheme != "http" && scheme != "https" {
		return nil, errors.New("azcosmos: HTTP request requires http or https")
	}
	port := req.URL.Port()
	if port == "" {
		port = "80"
		if scheme == "https" {
			port = "443"
		}
	}
	key := cosmosEndpoint{scheme: scheme, address: net.JoinHostPort(strings.ToLower(req.URL.Hostname()), port)}
	trace := httptrace.ContextClientTrace(req.Context())
	if trace != nil && trace.GetConn != nil {
		trace.GetConn(key.address)
	}
	p.mu.Lock()
	e := p.endpoints[key]
	if e == nil {
		e = &cosmosEndpointPool{endpoint: key, signal: cosmosPoolSignal{changed: make(chan struct{}, 1)}}
		p.endpoints[key] = e
	}
	e.users++
	p.mu.Unlock()
	defer func() {
		p.mu.Lock()
		e.users--
		p.removeEndpointLocked(e)
		p.mu.Unlock()
		e.signal.notify()
	}()

	pc, info, err := p.acquire(req, e)
	if err != nil {
		return nil, err
	}
	if trace != nil && trace.GotConn != nil {
		trace.GotConn(info)
	}
	resp, err := pc.conn.RoundTrip(req)
	if err != nil {
		pc.done()
		return nil, err
	}
	p.mu.Lock()
	pc.ready = true
	p.mu.Unlock()
	e.signal.notify()
	resp.Request = originalRequest
	if resp.Body == nil || resp.Body == http.NoBody {
		pc.done()
	} else {
		resp.Body = &cosmosPoolResponseBody{ReadCloser: resp.Body, done: pc.done}
	}
	return resp, nil
}

func (p *cosmosConnectionPool) acquire(req *http.Request, e *cosmosEndpointPool) (*cosmosPoolConn, httptrace.GotConnInfo, error) {
	for {
		changed := e.signal.channel()
		if err := req.Context().Err(); err != nil {
			return nil, httptrace.GotConnInfo{}, err
		}
		p.mu.Lock()
		var retired []*cosmosPoolConn
		for _, pc := range e.conns {
			// Until the first response, H2's initial optimistic stream limit
			// may exceed the server's SETTINGS. Send only one initial request.
			if pc.conn == nil || pc.closing || (!pc.ready && pc.users != 0) {
				continue
			}
			if pc.conn.Err() != nil {
				p.retireLocked(pc)
				retired = append(retired, pc)
				continue
			}
			if len(retired) != 0 || pc.conn.Available() == 0 || pc.conn.Reserve() != nil {
				continue
			}
			info := httptrace.GotConnInfo{Conn: pc.socket, Reused: pc.used, WasIdle: pc.idle != nil}
			if info.WasIdle {
				info.IdleTime = time.Since(pc.idleSince)
			}
			p.removeIdleLocked(pc)
			pc.users++
			pc.used = true
			p.mu.Unlock()
			e.signal.notify()
			return pc, info, nil
		}
		if len(retired) != 0 {
			p.mu.Unlock()
			for _, pc := range retired {
				pc.close()
			}
			continue
		}
		if len(e.conns) < p.limit && !e.dialing {
			pc := &cosmosPoolConn{pool: p, endpoint: e}
			e.conns = append(e.conns, pc)
			e.dialing = true
			p.mu.Unlock()
			ctx := context.WithValue(req.Context(), cosmosPoolDialKey{}, cosmosPoolDial{conn: pc, request: req})
			conn, err := p.transport.NewClientConn(ctx, e.endpoint.scheme, e.endpoint.address)
			p.mu.Lock()
			pc.conn = conn
			pc.users = 1
			e.dialing = false
			p.mu.Unlock()
			if err == nil {
				err = req.Context().Err()
			}
			if err != nil {
				pc.close()
				return nil, httptrace.GotConnInfo{}, err
			}
			conn.SetStateHook(func(*http.ClientConn) { e.signal.notify() })
			pc.done()
			e.signal.notify()
			continue
		}
		p.mu.Unlock()
		select {
		case <-changed:
		case <-req.Context().Done():
			return nil, httptrace.GotConnInfo{}, req.Context().Err()
		}
	}
}

func (p *cosmosConnectionPool) removeEndpointLocked(e *cosmosEndpointPool) {
	if e.users == 0 && len(e.conns) == 0 {
		delete(p.endpoints, e.endpoint)
	}
}

func (p *cosmosConnectionPool) removeIdleLocked(pc *cosmosPoolConn) {
	if pc.idle != nil {
		p.idle.Remove(pc.idle)
		pc.idle = nil
		pc.endpoint.idle--
	}
	if pc.timer != nil {
		pc.timer.Stop()
		pc.timer = nil
	}
}

func (p *cosmosConnectionPool) retireLocked(pc *cosmosPoolConn) {
	pc.closing = true
	p.removeIdleLocked(pc)
}

func (pc *cosmosPoolConn) done() {
	p := pc.pool
	p.mu.Lock()
	pc.users--
	if pc.users != 0 || pc.closing || pc.released {
		p.mu.Unlock()
		return
	}
	pc.idleSince = time.Now()
	pc.idle = p.idle.PushBack(pc)
	pc.endpoint.idle++
	idle := pc.idle
	pc.timer = time.AfterFunc(p.transport.IdleConnTimeout, func() {
		p.mu.Lock()
		closeConn := pc.idle == idle
		if closeConn {
			p.retireLocked(pc)
		}
		p.mu.Unlock()
		if closeConn {
			pc.close()
		}
	})
	var retired []*cosmosPoolConn
	for elem := p.idle.Front(); elem != nil; {
		next := elem.Next()
		candidate := elem.Value.(*cosmosPoolConn)
		if p.idle.Len() > p.transport.MaxIdleConns ||
			(candidate.endpoint == pc.endpoint && pc.endpoint.idle > p.transport.MaxIdleConnsPerHost) {
			p.retireLocked(candidate)
			retired = append(retired, candidate)
		}
		elem = next
	}
	p.mu.Unlock()
	pc.endpoint.signal.notify()
	for _, candidate := range retired {
		candidate.close()
	}
}

func (pc *cosmosPoolConn) close() {
	if pc.conn != nil {
		_ = pc.conn.Close()
	}
	if pc.socket != nil {
		_ = pc.socket.Close()
	}
	pc.release()
}

func (pc *cosmosPoolConn) release() {
	pc.releaseOnce.Do(func() {
		p := pc.pool
		p.mu.Lock()
		pc.released = true
		p.removeIdleLocked(pc)
		for i, entry := range pc.endpoint.conns {
			if entry == pc {
				pc.endpoint.conns = slices.Delete(pc.endpoint.conns, i, i+1)
				break
			}
		}
		p.removeEndpointLocked(pc.endpoint)
		p.mu.Unlock()
		pc.endpoint.signal.notify()
	})
}

func (p *cosmosConnectionPool) CloseIdleConnections() {
	p.mu.Lock()
	var retired []*cosmosPoolConn
	for elem := p.idle.Front(); elem != nil; elem = p.idle.Front() {
		pc := elem.Value.(*cosmosPoolConn)
		p.retireLocked(pc)
		retired = append(retired, pc)
	}
	p.mu.Unlock()
	for _, pc := range retired {
		pc.close()
	}
}

type cosmosPoolSocket struct {
	net.Conn
	once    sync.Once
	release func()
	err     error
}

func (c *cosmosPoolSocket) Close() error {
	c.once.Do(func() {
		c.err = c.Conn.Close()
		c.release()
	})
	return c.err
}

type cosmosPoolRequestBody struct {
	io.ReadCloser
	once sync.Once
	err  error
}

func (b *cosmosPoolRequestBody) Close() error {
	b.once.Do(func() { b.err = b.ReadCloser.Close() })
	return b.err
}

type cosmosPoolResponseBody struct {
	io.ReadCloser
	once sync.Once
	done func()
}

func (b *cosmosPoolResponseBody) Read(buf []byte) (int, error) {
	n, err := b.ReadCloser.Read(buf)
	if err != nil {
		b.once.Do(b.done)
	}
	return n, err
}

func (b *cosmosPoolResponseBody) Close() error {
	err := b.ReadCloser.Close()
	b.once.Do(b.done)
	return err
}
