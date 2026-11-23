// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package exported

import (
	"crypto"
	"crypto/tls"
	"errors"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"weak"
)

// MTLSRequest is an operation value recording the binding certificate, if any, of the token authorizing
// a request. Authentication policies set it; transportPolicy uses it to choose the request's client.
type MTLSRequest struct {
	Certificate *tls.Certificate
}

// mtlsTransportsPerBase bounds the certificate-scoped transports derived from one base transport, so
// certificate rotation doesn't accumulate connection pools
const mtlsTransportsPerBase = 8

// mtlsIdleConnTimeout is the idle timeout of derived transports whose base has none. It's a variable so
// tests can shorten it.
var mtlsIdleConnTimeout = 90 * time.Second

// mtlsTransport is a certificate-scoped transport and the most recently acquired copy of its certificate
type mtlsTransport struct {
	tr   *http.Transport
	cert atomic.Pointer[tls.Certificate]
	// lastUsed orders entries for eviction; guarded by mtlsTransports.mu
	lastUsed uint64
}

// mtlsTransports caches certificate-scoped transports for the whole process. Pipelines sharing a base
// transport (by default, every pipeline) therefore share each certificate's connection pool, as they share
// the base transport's pool for other requests. Base transports are held weakly: when one is collected,
// its derived transports are closed and removed. Only requests authorized by bound tokens use this cache.
var mtlsTransports struct {
	mu   sync.Mutex
	tick uint64
	// bases maps a base transport to its derived transports, keyed by leaf certificate DER (see clientFor)
	bases map[weak.Pointer[http.Transport]]map[string]*mtlsTransport
}

// clientFor returns the client for requests authorized by a token bound to cert. Each certificate gets
// its own connection pool so rotation never reuses a connection authenticated with an old certificate.
func (tp *transportPolicy) clientFor(cert *tls.Certificate) (*http.Client, error) {
	// crypto/tls never reads private key material. To authenticate a client it calls Sign to prove
	// possession of the key, so PrivateKey must implement crypto.Signer. This is what allows
	// non-exportable keys: a KeyGuard key's Signer wraps an OS key handle and signs inside the VBS
	// trustlet. Checking here reports a missing or unusable key clearly rather than mid-handshake.
	var signer crypto.Signer
	if cert != nil {
		signer, _ = cert.PrivateKey.(crypto.Signer)
	}
	if signer == nil || len(cert.Certificate) == 0 {
		return nil, errors.New("mTLS PoP requires a certificate with a private key that implements crypto.Signer")
	}
	client, ok := tp.trans.(*http.Client)
	if !ok || client == nil {
		return nil, errors.New("mTLS PoP requires Transport to be an *http.Client")
	}
	rt := client.Transport
	if rt == nil {
		rt = http.DefaultTransport
	}
	base, ok := rt.(*http.Transport)
	if !ok || base == nil {
		return nil, errors.New("mTLS PoP requires a standard *http.Transport")
	}
	tr, err := mtlsTransportFor(base, cert)
	if err != nil {
		return nil, err
	}
	// The client is built for each request rather than cached, so that after an eviction no pipeline
	// keeps using a transport the cache no longer manages. The connection pool is in the transport.
	derived := *client
	derived.Transport = tr
	derived.CheckRedirect = func(*http.Request, []*http.Request) error {
		return errors.New("redirects aren't permitted for mTLS PoP")
	}
	return &derived, nil
}

// mtlsTransportFor returns the transport, derived from base, that presents cert
func mtlsTransportFor(base *http.Transport, cert *tls.Certificate) (*http.Transport, error) {
	// Credentials such as MSAL return a new *tls.Certificate for the same certificate on each
	// acquisition, so transports are keyed by the leaf's DER encoding rather than the pointer.
	// DER is canonical: the same certificate always has the same bytes, and a rotated certificate
	// (new serial number, validity, signature) never does. The bytes include the public key, so
	// equal keys imply the same private key, and a bound token's cnf claim (x5t#S256) is the
	// SHA-256 of these bytes. Intermediates don't affect binding and are excluded. The string
	// conversion stays inline because m[string(b)] doesn't allocate, whereas storing it in a
	// variable copies the certificate on every call.
	//
	// Copies of a certificate share its key but may carry different signers, for example a new OS
	// key handle after a previous one was closed. New handshakes therefore present the most recently
	// acquired copy, while established connections, already authenticated, remain in the pool.
	leaf := cert.Certificate[0]
	key := weak.Make(base)
	c := &mtlsTransports
	c.mu.Lock()
	defer c.mu.Unlock()
	c.tick++
	derived := c.bases[key]
	if t := derived[string(leaf)]; t != nil {
		t.cert.Store(cert)
		t.lastUsed = c.tick
		return t.tr, nil
	}
	entry, err := newMTLSTransport(base)
	if err != nil {
		return nil, err
	}
	if derived == nil {
		derived = map[string]*mtlsTransport{}
		if c.bases == nil {
			c.bases = map[weak.Pointer[http.Transport]]map[string]*mtlsTransport{}
		}
		c.bases[key] = derived
		// arg must not reference base, or base would never be collected
		runtime.AddCleanup(base, removeMTLSTransports, key)
	}
	if len(derived) >= mtlsTransportsPerBase {
		var lru string
		var oldest *mtlsTransport
		for k, t := range derived {
			if oldest == nil || t.lastUsed < oldest.lastUsed {
				lru, oldest = k, t
			}
		}
		oldest.tr.CloseIdleConnections()
		delete(derived, lru)
	}
	entry.cert.Store(cert)
	entry.lastUsed = c.tick
	derived[string(leaf)] = entry
	return entry.tr, nil
}

// newMTLSTransport derives a transport from base that presents the entry's certificate. The derived
// transport must not reference base; see mtlsTransportFor.
func newMTLSTransport(base *http.Transport) (*mtlsTransport, error) {
	//nolint:staticcheck // Both hooks bypass TLSClientConfig.
	if base.DialTLS != nil || base.DialTLSContext != nil {
		return nil, errors.New("mTLS PoP doesn't support custom TLS dialing")
	}
	tr := base.Clone()
	config := tr.TLSClientConfig.Clone()
	if config == nil {
		config = &tls.Config{}
	}
	if config.MinVersion > tls.VersionTLS12 || (config.MaxVersion != 0 && config.MaxVersion < tls.VersionTLS12) {
		return nil, errors.New("mTLS PoP requires TLS 1.2")
	}
	config.MinVersion = tls.VersionTLS12
	config.MaxVersion = tls.VersionTLS12
	config.Renegotiation = tls.RenegotiateOnceAsClient
	config.Certificates = nil
	// A shared session cache could resume a session authenticated with a different certificate.
	config.ClientSessionCache = nil
	entry := &mtlsTransport{tr: tr}
	config.GetClientCertificate = func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
		return entry.cert.Load(), nil
	}
	config.NextProtos = []string{"http/1.1"}
	tr.TLSClientConfig = config
	tr.ForceAttemptHTTP2 = false
	// Discard protocol callbacks that can retain the base transport and its connection pools.
	tr.TLSNextProto = map[string]func(string, *tls.Conn) http.RoundTripper{}
	tr.Protocols = new(http.Protocols)
	tr.Protocols.SetHTTP1(true)
	// Eviction closes a transport's idle connections, but a request that selected the transport just
	// before eviction can still open one and return it to the pool afterward, where nothing would close
	// it. Tracking in-flight requests can't prevent this, because net/http may return a connection to
	// the pool after the response body's Close returns. A finite idle timeout ensures such connections
	// close, after which the transport and its certificate can be collected.
	if tr.IdleConnTimeout <= 0 {
		tr.IdleConnTimeout = mtlsIdleConnTimeout
	}
	return entry, nil
}

// removeMTLSTransports releases the transports derived from a collected base transport
func removeMTLSTransports(key weak.Pointer[http.Transport]) {
	c := &mtlsTransports
	c.mu.Lock()
	derived := c.bases[key]
	delete(c.bases, key)
	c.mu.Unlock()
	for _, t := range derived {
		t.tr.CloseIdleConnections()
	}
}
