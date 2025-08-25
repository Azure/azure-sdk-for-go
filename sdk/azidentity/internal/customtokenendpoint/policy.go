// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package customtokenendpoint

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

const (
	AzureKubernetesCAData        = "AZURE_KUBERNETES_CA_DATA"
	AzureKubernetesCAFile        = "AZURE_KUBERNETES_CA_FILE"
	AzureKubernetesSNIName       = "AZURE_KUBERNETES_SNI_NAME"
	AzureKubernetesTokenEndpoint = "AZURE_KUBERNETES_TOKEN_ENDPOINT"
)

func parseAndValidateCustomTokenEndpoint(endpoint string) (*url.URL, error) {
	tokenEndpoint, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse custom token endpoint URL %q: %s", endpoint, err)
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

var (
	errCustomEndpointEnvSetWithoutTokenEndpoint = fmt.Errorf(
		"AZURE_KUBERNETES_TOKEN_ENDPOINT is not set but other custom endpoint-related environment variables are present",
	)
	errCustomEndpointMultipleCASourcesSet = fmt.Errorf(
		"only one of AZURE_KUBERNETES_CA_FILE and AZURE_KUBERNETES_CA_DATA can be specified",
	)
)

func createBaseTransport() *http.Transport {
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
}

// Configure configures custom token endpoint mode if the required environment variables are present.
func Configure(clientOptions *policy.ClientOptions) error {
	kubernetesTokenEndpointStr := os.Getenv(AzureKubernetesTokenEndpoint)
	kubernetesSNIName := os.Getenv(AzureKubernetesSNIName)
	kubernetesCAFile := os.Getenv(AzureKubernetesCAFile)
	kubernetesCAData := os.Getenv(AzureKubernetesCAData)

	if kubernetesTokenEndpointStr == "" {
		// custom token endpoint is not set, while other Kubernetes-related environment variables are present,
		// this is likely a configuration issue so erroring out to avoid misconfiguration
		if kubernetesSNIName != "" || kubernetesCAFile != "" || kubernetesCAData != "" {
			return errCustomEndpointEnvSetWithoutTokenEndpoint
		}

		return nil
	}
	tokenEndpoint, err := parseAndValidateCustomTokenEndpoint(kubernetesTokenEndpointStr)
	if err != nil {
		return err
	}

	// CAFile and CAData are mutually exclusive, at most one can be set.
	// If none of CAFile or CAData are set, the default system CA pool will be used.
	if kubernetesCAFile != "" && kubernetesCAData != "" {
		return errCustomEndpointMultipleCASourcesSet
	}

	clientOptions.PerRetryPolicies = append(
		clientOptions.PerRetryPolicies,
		&customTokenEndpointPolicy{
			caFile:        kubernetesCAFile,
			caData:        kubernetesCAData,
			sniName:       kubernetesSNIName,
			tokenEndpoint: tokenEndpoint,
			baseTransport: createBaseTransport(),
		},
	)
	return nil
}

// customTokenEndpointPolicy is a custom HTTP transport that redirects token requests
// to the custom token endpoint when in custom token endpoint mode.
//
// Only token request will be handled by this transport.
type customTokenEndpointPolicy struct {
	caFile        string
	caData        string
	sniName       string
	tokenEndpoint *url.URL
	baseTransport *http.Transport
}

func (i *customTokenEndpointPolicy) Do(req *policy.Request) (*http.Response, error) {
	isTokenRequest, err := i.isTokenRequest(req)
	switch {
	case err != nil:
		// unable to determine if this is a token request, fail the request
		// We don't pass the request to the fallback client as the request body might be in half broken state.
		return nil, err
	case !isTokenRequest:
		// not a token request, fallback to the original transporter
		return req.Next()
	}

	tr, err := i.getTokenTransporter()
	if err != nil {
		return nil, err
	}

	rawReq := req.Raw()
	rawReq.URL.Scheme = i.tokenEndpoint.Scheme // this will always be https
	rawReq.URL.Host = i.tokenEndpoint.Host
	rawReq.URL.Path = i.tokenEndpoint.Path
	rawReq.Host = i.tokenEndpoint.Host

	resp, err := tr.RoundTrip(rawReq)
	if err == nil && resp == nil {
		// this policy is effectively a transport, so it must handle
		// this rare case. Returning an error makes the retry policy
		// try the request again
		err = errors.New("received nil response")
	}
	return resp, err
}

// loadCAPool loads the CA certificate pool to use.
// If neither CA file nor CA data is provided, the default system CA pool will be used.
func (i *customTokenEndpointPolicy) loadCAPool() (*x509.CertPool, error) {
	if i.caFile == "" && i.caData == "" {
		// nil CertPool indicates that the default system CA pool should be used
		return nil, nil
	}

	var caDataBytes []byte
	var err error

	if i.caFile != "" {
		caDataBytes, err = os.ReadFile(i.caFile)
		if err != nil {
			return nil, fmt.Errorf("read CA file %q: %s", i.caFile, err)
		}
		if len(caDataBytes) == 0 {
			return nil, fmt.Errorf("CA file %q is empty", i.caFile)
		}
	} else {
		caDataBytes = []byte(i.caData)
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
func (i *customTokenEndpointPolicy) getTokenTransporter() (*http.Transport, error) {
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
	tokenRequestContentTypePrefix        = "application/x-www-form-urlencoded"
)

func (i *customTokenEndpointPolicy) isTokenRequest(req *policy.Request) (bool, error) {
	if !strings.EqualFold(req.Raw().Method, http.MethodPost) {
		return false, nil
	}
	contentType := req.Raw().Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, tokenRequestContentTypePrefix) {
		// expecting token request to use application/x-www-form-urlencoded
		return false, nil
	}
	if req.Raw().ContentLength <= 0 {
		// expecting non-empty request body
		return false, nil
	}

	// peak the request body to confirm it's asking for token
	reqBody := req.Body()
	b, err := io.ReadAll(reqBody)
	if err != nil {
		// unable to read the form body at all, fail the whole token request in caller
		return false, err
	}
	_ = reqBody.Close()
	if err := req.RewindBody(); err != nil {
		return false, err
	}

	qs, err := url.ParseQuery(string(b))
	if err != nil {
		return false, nil // unable to process the form body, treat as non token request
	}

	// check if it's a token request with client assertion set
	return qs.Has(tokenRequestClientAssertionField) && qs.Has(tokenRequestClientAssertionTypeField), nil
}
