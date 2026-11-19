// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"crypto/tls"
	"errors"
	"net/http"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/transport"
)

type mtlsRequest struct {
	certificate *tls.Certificate
}

// Each certificate gets its own connection pool so rotation never reuses a connection
// authenticated with an old certificate. The cache is bounded across rotations.
type mtlsPolicy struct {
	base    policy.Transporter
	mu      sync.Mutex
	clients map[*tls.Certificate]*http.Client
}

func (p *mtlsPolicy) Do(req *policy.Request) (*http.Response, error) {
	var mtls mtlsRequest
	req.OperationValue(&mtls)
	if mtls.certificate == nil {
		for header, values := range req.Raw().Header {
			if strings.EqualFold(header, "x-ms-tokenboundauth") {
				return nil, errorinfo.NonRetriableError(errors.New("mTLS PoP requires a token with a binding certificate"))
			}
			if strings.EqualFold(header, "Authorization") {
				for _, value := range values {
					if fields := strings.Fields(value); len(fields) > 0 && strings.EqualFold(fields[0], "mtls_pop") {
						return nil, errorinfo.NonRetriableError(errors.New("mTLS PoP requires a token with a binding certificate"))
					}
				}
			}
		}
		return req.Next()
	}
	if req.Raw().URL.Scheme != "https" {
		return nil, errorinfo.NonRetriableError(errors.New("mTLS PoP requires HTTPS"))
	}
	client, err := p.client(mtls.certificate)
	if err != nil {
		return nil, errorinfo.NonRetriableError(err)
	}
	return client.Do(req.Raw())
}

func (p *mtlsPolicy) client(cert *tls.Certificate) (*http.Client, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if client := p.clients[cert]; client != nil {
		return client, nil
	}
	if cert == nil || len(cert.Certificate) == 0 || cert.PrivateKey == nil {
		return nil, errors.New("mTLS PoP requires a certificate with a private key")
	}
	base := p.base
	if base == nil {
		base = transport.DefaultHTTPClient
	}
	client, ok := base.(*http.Client)
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
	config.GetClientCertificate = func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
		return cert, nil
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
	if p.clients == nil {
		p.clients = map[*tls.Certificate]*http.Client{}
	}
	if len(p.clients) >= 8 {
		for old, c := range p.clients {
			c.CloseIdleConnections()
			delete(p.clients, old)
			break
		}
	}
	p.clients[cert] = &derived
	return &derived, nil
}
