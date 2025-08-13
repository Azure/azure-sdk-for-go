//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
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

	// configure custom token endpoint if environment variables are present.
	// In custom token endpoint mode, a dedicated transport will be used for proxying token requests to a dedicated endpoint.
	if err := w.configureCustomTokenEndpoint(caco); err != nil {
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

// configureCustomTokenEndpoint configures custom token endpoint mode if the required environment variables are present
func (w *WorkloadIdentityCredential) configureCustomTokenEndpoint(caco *ClientAssertionCredentialOptions) error {
	// check for custom token endpoint environment variables
	kubernetesTokenEndpointStr := os.Getenv(azureKubernetesTokenEndpoint)
	
	if kubernetesTokenEndpointStr == "" {
		// custom token endpoint is not set, skip configuration
		return nil
	}

	// capture values of kubernetesSNIName and kubernetesCAFile
	kubernetesSNIName := os.Getenv(azureKubernetesSNIName)
	kubernetesCAFile := os.Getenv(azureKubernetesCAFile)

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
)

// customTokenEndpointTransport is a custom HTTP transport that redirects token requests
// to the Kubernetes token endpoint when in custom token endpoint mode
type customTokenEndpointTransport struct {
	caFile              string
	sniName             string
	tokenEndpoint       *url.URL
	fallbackTransporter policy.Transporter
	baseTransport       *http.Transport
}

func newIdentityBindingTransport(
	caFile, sniName, tokenEndpointStr string,
	fallbackTransporter policy.Transporter,
) (*customTokenEndpointTransport, error) {
	tokenEndpoint, err := url.Parse(tokenEndpointStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token endpoint URL %q: %w", tokenEndpointStr, err)
	}

	// Require the tokenEndpoint to use https scheme
	if tokenEndpoint.Scheme != "https" {
		return nil, fmt.Errorf("token endpoint must use https scheme, got %q", tokenEndpoint.Scheme)
	}

	// Use default Kubernetes CA file location if none specified
	if caFile == "" {
		caFile = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
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

	tr := &customTokenEndpointTransport{
		caFile:              caFile,
		sniName:             sniName,
		tokenEndpoint:       tokenEndpoint,
		fallbackTransporter: fallbackTransporter,
		baseTransport:       initialTransport,
	}

	// Validate CA file exists and can be parsed
	if _, err := os.Stat(caFile); err != nil {
		return nil, fmt.Errorf("read CA file %q: %w", caFile, err)
	}
	
	// Test parsing the CA file
	caData, err := os.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("read CA file %q: %w", caFile, err)
	}
	
	if len(caData) > 0 {
		caPool := x509.NewCertPool()
		if !caPool.AppendCertsFromPEM(caData) {
			return nil, fmt.Errorf("parse CA file %q: no valid certificates found", caFile)
		}
	}

	return tr, nil
}

func (i *customTokenEndpointTransport) Do(req *http.Request) (*http.Response, error) {
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
	newReq.URL.Path = i.tokenEndpoint.Path
	newReq.Host = i.tokenEndpoint.Host

	return tr.RoundTrip(newReq)
}

// getTokenTransporter rebuilds the HTTP transport every time it's called.
// This approach is acceptable because token requests are infrequent due to token caching at higher levels.
func (i *customTokenEndpointTransport) getTokenTransporter() (*http.Transport, error) {
	transport := i.baseTransport.Clone()
	
	// Load and configure custom CA
	caData, err := os.ReadFile(i.caFile)
	if err != nil {
		return nil, fmt.Errorf("read CA file %q: %w", i.caFile, err)
	}
	
	// Handle empty files during rotation
	if len(caData) == 0 {
		// CA file might be empty during rotation, use base transport without custom CA
		return transport, nil
	}
	
	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caData) {
		return nil, fmt.Errorf("parse CA file %q: no valid certificates found", i.caFile)
	}
	
	if transport.TLSClientConfig == nil {
		transport.TLSClientConfig = &tls.Config{}
	}
	transport.TLSClientConfig.RootCAs = caPool
	
	// Configure SNI if specified
	if i.sniName != "" {
		transport.TLSClientConfig.ServerName = i.sniName
	}
	
	return transport, nil
}
