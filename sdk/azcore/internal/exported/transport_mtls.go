// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package exported

import (
	"crypto"
	"crypto/tls"
	"errors"
	"net/http"
	"sync/atomic"
)

// MTLSRequest is an operation value recording the binding certificate, if any, of the token authorizing
// a request. Authentication policies set it; transportPolicy uses it to choose the request's client.
type MTLSRequest struct {
	Certificate *tls.Certificate
}

// mtlsClient is a certificate-scoped client and the most recently acquired copy of its certificate
type mtlsClient struct {
	client *http.Client
	cert   atomic.Pointer[tls.Certificate]
}

// clientFor returns the client for requests authorized by a token bound to cert. Each certificate gets
// its own connection pool so rotation never reuses a connection authenticated with an old certificate.
// The cache is bounded across rotations.
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
	// Credentials such as MSAL return a new *tls.Certificate for the same certificate on each
	// acquisition, so clients are keyed by the leaf's DER encoding rather than the pointer.
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
	tp.muClients.Lock()
	defer tp.muClients.Unlock()
	if c := tp.clients[string(leaf)]; c != nil {
		c.cert.Store(cert)
		return c.client, nil
	}
	client, ok := tp.trans.(*http.Client)
	if !ok || client == nil {
		return nil, errors.New("mTLS PoP requires Transport to be an *http.Client")
	}
	rt := client.Transport
	if rt == nil {
		rt = http.DefaultTransport
	}
	tr, ok := rt.(*http.Transport)
	if !ok || tr == nil {
		return nil, errors.New("mTLS PoP requires a standard *http.Transport")
	}
	//nolint:staticcheck // Both hooks bypass TLSClientConfig.
	if tr.DialTLS != nil || tr.DialTLSContext != nil {
		return nil, errors.New("mTLS PoP doesn't support custom TLS dialing")
	}
	tr = tr.Clone()
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
	entry := &mtlsClient{}
	entry.cert.Store(cert)
	config.GetClientCertificate = func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
		return entry.cert.Load(), nil
	}
	config.NextProtos = []string{"http/1.1"}
	tr.TLSClientConfig = config
	tr.ForceAttemptHTTP2 = false
	// Discard protocol callbacks that can retain the base transport's connection pools.
	tr.TLSNextProto = map[string]func(string, *tls.Conn) http.RoundTripper{}
	tr.Protocols = new(http.Protocols)
	tr.Protocols.SetHTTP1(true)
	derived := *client
	derived.Transport = tr
	derived.CheckRedirect = func(*http.Request, []*http.Request) error {
		return errors.New("redirects aren't permitted for mTLS PoP")
	}
	if tp.clients == nil {
		tp.clients = map[string]*mtlsClient{}
	}
	if len(tp.clients) >= 8 {
		for old, c := range tp.clients {
			c.client.CloseIdleConnections()
			delete(tp.clients, old)
			break
		}
	}
	entry.client = &derived
	tp.clients[string(leaf)] = entry
	return &derived, nil
}
