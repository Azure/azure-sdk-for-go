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
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
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
	if err := configureCustomTokenEndpoint(&caco.ClientOptions); err != nil {
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

func parseAndValidateCustomTokenEndpoint(endpoint string) (*url.URL, error) {
	tokenEndpoint, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse custom token endpoint URL %q: %w", endpoint, err)
	}
	if tokenEndpoint.Scheme != "https" {
		return nil, fmt.Errorf("custom token endpoint must use https scheme, got %q", tokenEndpoint.Scheme)
	}
	if tokenEndpoint.User != nil {
		return nil, fmt.Errorf("custom token endpoint URL %q must not contain user info", tokenEndpoint)
	}
	if tokenEndpoint.RawQuery != "" {
		return nil, fmt.Errorf("custom token endpoint URL %q must not contain a query", tokenEndpoint)
	}
	if tokenEndpoint.EscapedFragment() != "" {
		return nil, fmt.Errorf("custom token endpoint URL %q must not contain a fragment", tokenEndpoint)
	}
	return tokenEndpoint, nil
}

// configureCustomTokenEndpoint configures custom token endpoint mode if the required environment variables are present
func configureCustomTokenEndpoint(clientOptions *policy.ClientOptions) error {
	// check for custom token endpoint environment variables
	kubernetesTokenEndpointStr := os.Getenv(azureKubernetesTokenEndpoint)
	if kubernetesTokenEndpointStr == "" {
		// custom token endpoint is not set, skip configuration
		return nil
	}
	tokenEndpoint, err := parseAndValidateCustomTokenEndpoint(kubernetesTokenEndpointStr)
	if err != nil {
		return err
	}

	// capture values of kubernetesSNIName, kubernetesCAFile, and kubernetesCAData
	kubernetesSNIName := os.Getenv(azureKubernetesSNIName)
	kubernetesCAFile := os.Getenv(azureKubernetesCAFile)
	kubernetesCAData := os.Getenv(azureKubernetesCAData)

	// CAFile and CAData are mutually exclusive, at most one can be set.
	// If none of CAFile or CAData are set, the default system CA pool will be used.
	if kubernetesCAFile != "" && kubernetesCAData != "" {
		return fmt.Errorf("only one of AZURE_KUBERNETES_CA_FILE and AZURE_KUBERNETES_CA_DATA can be specified")
	}

	// creating a new azcore client to maintain the same behavior as confidential client's usage.
	// This fallback client is expected to be used for requesting the tenant discovery settings (non skippable).
	fallbackClient, err := azcore.NewClient(module, version, runtime.PipelineOptions{
		Tracing: runtime.TracingOptions{
			Namespace: traceNamespace,
		},
	}, clientOptions)
	if err != nil {
		return err
	}

	baseTransport := func() *http.Transport {
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

	clientOptions.Transport = &customTokenEndpointTransport{
		caFile:         kubernetesCAFile,
		caData:         kubernetesCAData,
		sniName:        kubernetesSNIName,
		tokenEndpoint:  tokenEndpoint,
		fallbackClient: fallbackClient,
		baseTransport:  baseTransport,
	}
	return nil
}

// customTokenEndpointTransport is a custom HTTP transport that redirects token requests
// to the custom token endpoint when in custom token endpoint mode.
//
// All non-token requests are forwarded to the fallback client.
type customTokenEndpointTransport struct {
	caFile         string
	caData         string
	sniName        string
	tokenEndpoint  *url.URL
	fallbackClient *azcore.Client
	baseTransport  *http.Transport
}

func (i *customTokenEndpointTransport) Do(req *http.Request) (*http.Response, error) {
	if !i.isTokenRequest(req) {
		// not a token request, fallback to the original transporter
		return doForClient(i.fallbackClient, req)
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

// loadCAPool loads the CA certificate pool to use.
// If neither CA file nor CA data is provided, the default system CA pool will be used.
func (i *customTokenEndpointTransport) loadCAPool() (*x509.CertPool, error) {
	if i.caFile == "" && i.caData == "" {
		// nil CertPool indicates that the default system CA pool should be used
		return nil, nil
	}

	var caDataBytes []byte
	var err error

	if i.caFile != "" {
		caDataBytes, err = os.ReadFile(i.caFile)
		if err != nil {
			return nil, fmt.Errorf("read CA file %q: %w", i.caFile, err)
		}
	} else if i.caData != "" {
		caDataBytes = []byte(i.caData)
	} else {
		return nil, fmt.Errorf("missing CA: neither CA file nor CA data provided")
	}

	// Error out if CA data is empty
	if len(caDataBytes) == 0 {
		if i.caFile != "" {
			return nil, fmt.Errorf("CA file %q is empty", i.caFile)
		} else {
			return nil, fmt.Errorf("CA data is empty")
		}
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caDataBytes) {
		if i.caFile != "" {
			return nil, fmt.Errorf("parse CA file %q: no valid certificates found", i.caFile)
		} else {
			return nil, fmt.Errorf("parse CA data: no valid certificates found")
		}
	}

	return caPool, nil
}

// getTokenTransporter rebuilds the HTTP transport every time it's called.
// This approach is acceptable because token requests are infrequent due to token caching at higher levels.
func (i *customTokenEndpointTransport) getTokenTransporter() (*http.Transport, error) {
	transport := i.baseTransport.Clone()

	caPool, err := i.loadCAPool()
	if err != nil {
		return nil, err
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

const (
	tokenRequestClientAssertionField     = "client_assertion"
	tokenRequestClientAssertionTypeField = "client_assertion_type"
)

func (i *customTokenEndpointTransport) isTokenRequest(req *http.Request) bool {
	if !strings.EqualFold(req.Method, http.MethodPost) {
		return false
	}

	if req.Body == nil || req.Body == http.NoBody {
		return false
	}

	b, err := io.ReadAll(req.Body)
	if err != nil {
		return false
	}
	req.Body.Close()                                   // close the original body to avoid resource leaks
	req.Body = streaming.NopCloser(bytes.NewReader(b)) // reset the body for future reads

	qs, err := url.ParseQuery(string(b))
	if err != nil {
		return false // unable to process
	}
	if qs.Has(tokenRequestClientAssertionField) && qs.Has(tokenRequestClientAssertionTypeField) {
		// this is a token request with client assertion set
		return true
	}

	return false
}
