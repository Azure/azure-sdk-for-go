//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

const credNameWorkloadIdentity = "WorkloadIdentityCredential"

// WorkloadIdentityCredential supports Azure workload identity on Kubernetes.
// See [Azure Kubernetes Service documentation] for more information.
//
// [Azure Kubernetes Service documentation]: https://learn.microsoft.com/azure/aks/workload-identity-overview
type WorkloadIdentityCredential struct {
	assertion, file string
	cred            *ClientAssertionCredential
	expires         time.Time
	mtx             *sync.RWMutex
}

// WorkloadIdentityCredentialOptions contains optional parameters for WorkloadIdentityCredential.
type WorkloadIdentityCredentialOptions struct {
	azcore.ClientOptions

	// AdditionallyAllowedTenants specifies additional tenants for which the credential may acquire tokens.
	// Add the wildcard value "*" to allow the credential to acquire tokens for any tenant in which the
	// application is registered.
	AdditionallyAllowedTenants []string

	// Cache is a persistent cache the credential will use to store the tokens it acquires, making
	// them available to other processes and credential instances. The default, zero value means the
	// credential will store tokens in memory and not share them with any other credential instance.
	Cache Cache

	// ClientID of the service principal. Defaults to the value of the environment variable AZURE_CLIENT_ID.
	ClientID string

	// DisableInstanceDiscovery should be set true only by applications authenticating in disconnected clouds, or
	// private clouds such as Azure Stack. It determines whether the credential requests Microsoft Entra instance metadata
	// from https://login.microsoft.com before authenticating. Setting this to true will skip this request, making
	// the application responsible for ensuring the configured authority is valid and trustworthy.
	DisableInstanceDiscovery bool

	// TenantID of the service principal. Defaults to the value of the environment variable AZURE_TENANT_ID.
	TenantID string

	// TokenFilePath is the path of a file containing a Kubernetes service account token. Defaults to the value of the
	// environment variable AZURE_FEDERATED_TOKEN_FILE.
	TokenFilePath string
}

// NewWorkloadIdentityCredential constructs a WorkloadIdentityCredential. Service principal configuration is read
// from environment variables as set by the Azure workload identity webhook. Set options to override those values.
func NewWorkloadIdentityCredential(options *WorkloadIdentityCredentialOptions) (*WorkloadIdentityCredential, error) {
	if options == nil {
		options = &WorkloadIdentityCredentialOptions{}
	}
	ok := false
	clientID := options.ClientID
	if clientID == "" {
		if clientID, ok = os.LookupEnv(azureClientID); !ok {
			return nil, errors.New("no client ID specified. Check pod configuration or set ClientID in the options")
		}
	}
	file := options.TokenFilePath
	if file == "" {
		if file, ok = os.LookupEnv(azureFederatedTokenFile); !ok {
			return nil, errors.New("no token file specified. Check pod configuration or set TokenFilePath in the options")
		}
	}
	tenantID := options.TenantID
	if tenantID == "" {
		if tenantID, ok = os.LookupEnv(azureTenantID); !ok {
			return nil, errors.New("no tenant ID specified. Check pod configuration or set TenantID in the options")
		}
	}

	w := WorkloadIdentityCredential{file: file, mtx: &sync.RWMutex{}}
	caco := &ClientAssertionCredentialOptions{
		AdditionallyAllowedTenants: options.AdditionallyAllowedTenants,
		Cache:                      options.Cache,
		ClientOptions:              options.ClientOptions,
		DisableInstanceDiscovery:   options.DisableInstanceDiscovery,
	}

	// configure identity binding if environment variables are present.
	// In identity binding enabled mode, a dedicated transport will be used for proxying token requests to a dedicated endpoint.
	if err := w.configureIdentityBinding(caco); err != nil {
		return nil, err
	}

	cred, err := NewClientAssertionCredential(tenantID, clientID, w.getAssertion, caco)
	if err != nil {
		return nil, err
	}
	// we want "WorkloadIdentityCredential" in log messages, not "ClientAssertionCredential"
	cred.client.name = credNameWorkloadIdentity
	w.cred = cred
	return &w, nil
}

// GetToken requests an access token from Microsoft Entra ID. Azure SDK clients call this method automatically.
func (w *WorkloadIdentityCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, credNameWorkloadIdentity+"."+traceOpGetToken, w.cred.client.azClient.Tracer(), nil)
	defer func() { endSpan(err) }()
	tk, err := w.cred.GetToken(ctx, opts)
	return tk, err
}

// getAssertion returns the specified file's content, which is expected to be a Kubernetes service account token.
// Kubernetes is responsible for updating the file as service account tokens expire.
func (w *WorkloadIdentityCredential) getAssertion(context.Context) (string, error) {
	w.mtx.RLock()
	if w.expires.Before(time.Now()) {
		// ensure only one goroutine at a time updates the assertion
		w.mtx.RUnlock()
		w.mtx.Lock()
		defer w.mtx.Unlock()
		// double check because another goroutine may have acquired the write lock first and done the update
		if now := time.Now(); w.expires.Before(now) {
			content, err := os.ReadFile(w.file)
			if err != nil {
				return "", err
			}
			w.assertion = string(content)
			// Kubernetes rotates service account tokens when they reach 80% of their total TTL. The shortest TTL
			// is 1 hour. That implies the token we just read is valid for at least 12 minutes (20% of 1 hour),
			// but we add some margin for safety.
			w.expires = now.Add(10 * time.Minute)
		}
	} else {
		defer w.mtx.RUnlock()
	}
	return w.assertion, nil
}

// configureIdentityBinding configures identity binding mode if the required environment variables are present
func (w *WorkloadIdentityCredential) configureIdentityBinding(caco *ClientAssertionCredentialOptions) error {
	// check for identity binding mode environment variables
	kubernetesTokenEndpointStr := os.Getenv(azureKubernetesTokenEndpoint)
	kubernetesSNIName := os.Getenv(azureKubernetesSNIName)
	kubernetesCAFile := os.Getenv(azureKubernetesCAFile)

	if kubernetesTokenEndpointStr == "" && kubernetesSNIName == "" && kubernetesCAFile == "" {
		// identity binding is not set
		return nil
	}

	// All three variables must be present for identity binding mode
	if kubernetesTokenEndpointStr == "" || kubernetesSNIName == "" || kubernetesCAFile == "" {
		return errors.New("identity binding mode requires all three environment variables: AZURE_KUBERNETES_TOKEN_ENDPOINT, AZURE_KUBERNETES_SNI_NAME, and AZURE_KUBERNETES_CA_FILE")
	}

	transporter, err := newIdentityBindingTransport(
		kubernetesCAFile, kubernetesSNIName, kubernetesTokenEndpointStr,
		caco.Transport,
	)
	if err != nil {
		return err
	}
	caco.Transport = transporter
	return nil
}

const (
	tokenEndpointSuffix = "/oauth2/v2.0/token"
	caReloadInterval    = 10 * time.Minute
)

// identityBindingTransport is a custom HTTP transport that redirects token requests
// to the Kubernetes token endpoint when in identity binding mode
type identityBindingTransport struct {
	caFile              string
	sniName             string
	tokenEndpoint       *url.URL
	fallbackTransporter policy.Transporter

	mtx *sync.RWMutex

	nextRead  time.Time
	currentCA []byte
	transport *http.Transport
}

func newIdentityBindingTransport(
	caFile, sniName, tokenEndpointStr string,
	fallbackTransporter policy.Transporter,
) (*identityBindingTransport, error) {
	tokenEndpoint, err := url.Parse(tokenEndpointStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token endpoint URL %q: %w", tokenEndpointStr, err)
	}

	if fallbackTransporter == nil {
		// FIXME: can we callback to the defaultHTTPClient from azcore/runtime?
		fallbackTransporter = http.DefaultClient
	}

	initialTransport := func() *http.Transport {
		// try reusing the user provided transport if available
		if httpClient, ok := fallbackTransporter.(*http.Client); ok {
			if transport, ok := httpClient.Transport.(*http.Transport); ok {
				return transport.Clone()
			}
		}

		// if the user did not provide a policy.Transporter or it's not a *http.Client,
		// we fall back to the default one.
		// FIXME: can we callback to the defaultHTTPClient from azcore/runtime?
		if transport, ok := http.DefaultTransport.(*http.Transport); ok {
			return transport.Clone()
		}

		// this should not happen, but if the user mutates the net/http.DefaultTransport
		// to something else, we fall back to a sane default
		return &http.Transport{
			ForceAttemptHTTP2:   true,
			MaxIdleConns:        100,
			IdleConnTimeout:     90 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
		}
	}()

	tr := &identityBindingTransport{
		caFile:              caFile,
		sniName:             sniName,
		tokenEndpoint:       tokenEndpoint,
		fallbackTransporter: fallbackTransporter,
		mtx:                 &sync.RWMutex{},
		transport:           initialTransport,
	}

	// perform an initial load to surface any issues with the CA file and transport settings.
	// Lock is not held here as this is called in the constructor
	if err := tr.reloadCA(); err != nil {
		return nil, err
	}

	return tr, nil
}

func (i *identityBindingTransport) Do(req *http.Request) (*http.Response, error) {
	if !strings.HasSuffix(req.URL.Path, tokenEndpointSuffix) {
		// not a token request, fallback to the original transporter
		return i.fallbackTransporter.Do(req)
	}

	tr, err := i.getTokenTransporter()
	if err != nil {
		return nil, err
	}

	newReq := req.Clone(req.Context())
	newReq.URL.Scheme = i.tokenEndpoint.Scheme // this will always be https
	newReq.URL.Host = i.tokenEndpoint.Host
	newReq.URL.Path = ""
	newReq.Host = i.tokenEndpoint.Host

	return tr.RoundTrip(newReq)
}

func (i *identityBindingTransport) getTokenTransporter() (*http.Transport, error) {
	i.mtx.RLock()
	if i.nextRead.Before(time.Now()) {
		i.mtx.RUnlock()
		i.mtx.Lock()
		defer i.mtx.Unlock()
		// double check on the read time
		if now := time.Now(); i.nextRead.Before(now) {
			if err := i.reloadCA(); err != nil {
				// we return error if any attempt of reloading CA fails
				// This should surface in the token calls and we expect the caller to
				// have proper error handling / rate limit so we don't fall into deadloop here
				// due to scenario like broken CA file.
				return nil, err
			}
		}
	} else {
		defer i.mtx.RUnlock()
	}
	return i.transport, nil
}

func (i *identityBindingTransport) createTransportWithCAPool(
	fromTransport *http.Transport,
	caPool *x509.CertPool,
) *http.Transport {
	transport := fromTransport.Clone()
	if transport.TLSClientConfig == nil {
		transport.TLSClientConfig = &tls.Config{}
	}
	transport.TLSClientConfig.ServerName = i.sniName
	transport.TLSClientConfig.RootCAs = caPool
	return transport
}

// reloadCA attempts to read the latest CA from the CA file and updates the transport if the content has changed.
// If a new CA is discovered, the existing transport will be replaced with a new one that uses the new CA.
// It expects the caller to hold the write lock on i.mtx to ensure thread safety.
func (i *identityBindingTransport) reloadCA() error {
	newCA, err := os.ReadFile(i.caFile)
	if err != nil {
		return fmt.Errorf("read CA file %q: %w", i.caFile, err)
	}

	if len(newCA) == 0 {
		// the CA file might be in the middle of rotation without the content written.
		// We return nil and rely on next check.
		return nil
	}

	if bytes.Equal(i.currentCA, newCA) {
		// no change in CA content, no need to replace
		i.nextRead = time.Now().Add(caReloadInterval)
		return nil
	}

	newCAPool := x509.NewCertPool()
	if !newCAPool.AppendCertsFromPEM(newCA) {
		return fmt.Errorf("parse CA file %q: no valid certificates found", i.caFile)
	}

	newTransport := i.createTransportWithCAPool(i.transport, newCAPool)
	oldTransport := i.transport

	i.transport = newTransport
	i.currentCA = newCA
	i.nextRead = time.Now().Add(caReloadInterval)

	if oldTransport != nil {
		// drop any idle connections from previous transport so new requests can be
		// moved to the new transport
		oldTransport.CloseIdleConnections()
	}

	return nil
}
