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
	// identity binding mode fields
	identityBinding         bool
	kubernetesCAFile        string
	kubernetesSNIName       string
	kubernetesTokenEndpoint *url.URL
	// CA certificate caching fields
	caCert     []byte
	caCertPool *x509.CertPool
	caExpires  time.Time
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

	// Check for identity binding mode environment variables
	kubernetesTokenEndpointStr := os.Getenv(azureKubernetesTokenEndpoint)
	kubernetesSNIName := os.Getenv(azureKubernetesSNIName)
	kubernetesCAFile := os.Getenv(azureKubernetesCAFile)

	// If any of the identity binding environment variables are present, enable identity binding mode
	if kubernetesTokenEndpointStr != "" || kubernetesSNIName != "" || kubernetesCAFile != "" {
		// All three variables must be present for identity binding mode
		if kubernetesTokenEndpointStr == "" || kubernetesSNIName == "" || kubernetesCAFile == "" {
			return nil, errors.New("identity binding mode requires all three environment variables: AZURE_KUBERNETES_TOKEN_ENDPOINT, AZURE_KUBERNETES_SNI_NAME, and AZURE_KUBERNETES_CA_FILE")
		}
		
		// Parse the Kubernetes token endpoint URL
		kubernetesTokenEndpoint, err := url.Parse(kubernetesTokenEndpointStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Kubernetes token endpoint URL: %w", err)
		}
		
		w.identityBinding = true
		w.kubernetesTokenEndpoint = kubernetesTokenEndpoint
		w.kubernetesSNIName = kubernetesSNIName
		w.kubernetesCAFile = kubernetesCAFile
	}

	caco := ClientAssertionCredentialOptions{
		AdditionallyAllowedTenants: options.AdditionallyAllowedTenants,
		Cache:                      options.Cache,
		ClientOptions:              options.ClientOptions,
		DisableInstanceDiscovery:   options.DisableInstanceDiscovery,
	}

	// If identity binding mode is enabled, configure a custom HTTP client
	if w.identityBinding {
		// Override the client options to use our custom transport
		caco.ClientOptions.Transport = &identityBindingTransport{
			credential:        &w,
			kubernetesSNIName: w.kubernetesSNIName,
		}
	}

	cred, err := NewClientAssertionCredential(tenantID, clientID, w.getAssertion, &caco)
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

// loadKubernetesCA loads and caches the Kubernetes CA certificate
func (w *WorkloadIdentityCredential) loadKubernetesCA() (*x509.CertPool, error) {
	w.mtx.RLock()
	if w.caExpires.After(time.Now()) && w.caCertPool != nil {
		defer w.mtx.RUnlock()
		return w.caCertPool, nil
	}
	
	// ensure only one goroutine at a time updates the CA cert
	w.mtx.RUnlock()
	w.mtx.Lock()
	defer w.mtx.Unlock()
	
	// double check because another goroutine may have acquired the write lock first and done the update
	if now := time.Now(); w.caExpires.After(now) && w.caCertPool != nil {
		return w.caCertPool, nil
	}
	
	// Load the CA certificate for the Kubernetes endpoint
	caCert, err := os.ReadFile(w.kubernetesCAFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read Kubernetes CA file: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, errors.New("failed to parse Kubernetes CA certificate")
	}
	
	// Cache the CA certificate for 10 minutes (same as token assertion)
	w.caCert = caCert
	w.caCertPool = caCertPool
	w.caExpires = time.Now().Add(10 * time.Minute)
	
	return caCertPool, nil
}

// identityBindingTransport is a custom HTTP transport that redirects token requests
// to the Kubernetes token endpoint when in identity binding mode
type identityBindingTransport struct {
	credential        *WorkloadIdentityCredential
	kubernetesSNIName string
}

func (t *identityBindingTransport) Do(req *http.Request) (*http.Response, error) {
	// Check if this is a token request to the Azure authority host
	if strings.HasSuffix(req.URL.Path, "/oauth2/v2.0/token") && (req.URL.Host == "login.microsoftonline.com" || 
		req.URL.Host == "login.microsoftonline.us" || 
		req.URL.Host == "login.partner.microsoftonline.cn" ||
		req.URL.Host == "login.microsoftonline.de") {
		// This is a token request, redirect to Kubernetes endpoint
		
		// Load the CA certificate (this will use cached version if still valid)
		caCertPool, err := t.credential.loadKubernetesCA()
		if err != nil {
			return nil, err
		}
		
		// Create custom transport with the CA and SNI configuration
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:    caCertPool,
				ServerName: t.kubernetesSNIName,
			},
		}
		
		// Clone the request to avoid modifying the original
		newReq := req.Clone(req.Context())
		
		// Update the URL to point to the Kubernetes endpoint
		newReq.URL.Scheme = t.credential.kubernetesTokenEndpoint.Scheme
		newReq.URL.Host = t.credential.kubernetesTokenEndpoint.Host
		newReq.Host = t.credential.kubernetesTokenEndpoint.Host
		
		// Preserve the original path (contains tenant ID and token endpoint path)
		// The path should be something like "/tenant-id/oauth2/v2.0/token"
		// Keep the original path to maintain the token request structure
		
		return transport.RoundTrip(newReq)
	}
	
	// For non-token requests, use the default transport
	return http.DefaultTransport.RoundTrip(req)
}
