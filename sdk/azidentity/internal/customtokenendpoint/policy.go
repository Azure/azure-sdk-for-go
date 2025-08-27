// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package customtokenendpoint

import (
	"bytes"
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
	errCustomEndpointEnvSetWithoutTokenEndpoint = errors.New(
		"AZURE_KUBERNETES_TOKEN_ENDPOINT is not set but other custom endpoint-related environment variables are present",
	)
	errCustomEndpointMultipleCASourcesSet = errors.New(
		"only one of AZURE_KUBERNETES_CA_FILE and AZURE_KUBERNETES_CA_DATA can be specified",
	)
)

func createTransport(sniName string, caPool *x509.CertPool) *http.Transport {
	var transport *http.Transport
	if tr, ok := http.DefaultTransport.(*http.Transport); ok {
		transport = tr.Clone()
	} else {
		// this should not happen, but if the user mutates the net/http.DefaultTransport
		// to something else, we fall back to a sane default
		transport = &http.Transport{
			ForceAttemptHTTP2:   true,
			MaxIdleConns:        100,
			IdleConnTimeout:     90 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
		}
	}

	if transport.TLSClientConfig == nil {
		transport.TLSClientConfig = &tls.Config{}
	}
	transport.TLSClientConfig.ServerName = sniName
	transport.TLSClientConfig.RootCAs = caPool

	return transport
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

	// preload the transport
	p := &customTokenEndpointPolicy{
		caFile:        kubernetesCAFile,
		caData:        []byte(kubernetesCAData),
		sniName:       kubernetesSNIName,
		tokenEndpoint: tokenEndpoint,
	}
	if _, err := p.getTokenTransporter(); err != nil {
		return err
	}

	clientOptions.PerRetryPolicies = append(clientOptions.PerRetryPolicies, p)
	return nil
}

// customTokenEndpointPolicy is a custom HTTP transport that redirects token requests
// to the custom token endpoint when in custom token endpoint mode.
//
// Only token request will be handled by this transport.
// Lock is not needed for internal caData as this policy is called under confidentialClient's lock.
type customTokenEndpointPolicy struct {
	caFile        string
	caData        []byte
	sniName       string
	tokenEndpoint *url.URL
	transport     *http.Transport
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

// getTokenTransporter provides the token transport to use for the request.
//
// There are a few scenarios need to be handled:
//  1. no CA overrides, use default transport
//  2. CA data override provided, use a transport with custom CA pool.
//     This transport is fixed after set.
//  3. CA file override is provided, use a transport with custom CA pool.
//     This transport needs to be recreated if the CA file content changes.
func (i *customTokenEndpointPolicy) getTokenTransporter() (*http.Transport, error) {
	if len(i.caData) == 0 && i.caFile == "" {
		// no custom CA overrides
		if i.transport == nil {
			i.transport = createTransport(i.sniName, nil)
		}
		return i.transport, nil
	}

	if i.caFile == "" {
		// host provided CA bytes in AZURE_KUBERNETES_CA_DATA and can't change
		// them now, so we need to create a client only if we haven't done so yet
		if i.transport != nil {
			return i.transport, nil
		}

		caPool := x509.NewCertPool()
		if !caPool.AppendCertsFromPEM([]byte(i.caData)) {
			return nil, fmt.Errorf("parse CA data: no valid certificates found")
		}

		i.transport = createTransport(i.sniName, caPool)
		return i.transport, nil
	}

	// host provided the CA bytes in a file whose contents it can change,
	// so we must read that file and maybe create a new client
	b, err := os.ReadFile(i.caFile)
	if err != nil {
		return nil, fmt.Errorf("read CA file %q: %s", i.caFile, err)
	}
	if len(b) == 0 {
		// this can happen during the middle of CA rotation on the host.
		// Erroring out here to force the client to retry the token call.
		return nil, fmt.Errorf("CA file %q is empty", i.caFile)
	}
	if !bytes.Equal(b, i.caData) {
		// CA has changed, rebuild the transport with new CA pool
		// invariant: i.transport is nil when i.caData is nil (initial call)
		caPool := x509.NewCertPool()
		if !caPool.AppendCertsFromPEM([]byte(b)) {
			return nil, fmt.Errorf("parse CA file %q: no valid certificates found", i.caFile)
		}
		i.transport = createTransport(i.sniName, caPool)
		i.caData = b
	}

	return i.transport, nil
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
