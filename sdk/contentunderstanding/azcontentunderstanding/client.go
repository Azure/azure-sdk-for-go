// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package azcontentunderstanding provides a prototype Azure AI Content Understanding client.
package azcontentunderstanding

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
)

const (
	moduleName = "github.com/Azure/azure-sdk-for-go/sdk/contentunderstanding/azcontentunderstanding"
	version    = "v0.1.0"
	// APIVersionGA is the stable service API version.
	APIVersionGA = "2025-11-01"
	// APIVersionPreview supports inline analysis.
	APIVersionPreview = "2026-06-01-preview"
)

// ClientOptions configures a Client.
type ClientOptions struct {
	azcore.ClientOptions
	// APIVersion overrides the embedded ClientOptions.APIVersion and defaults to APIVersionGA.
	APIVersion string
	// Audience overrides the default https://cognitiveservices.azure.com audience.
	Audience string
}

// Client calls the Content Understanding data-plane service.
// Methods are safe for concurrent use. Pollers and pagers are not.
type Client struct {
	internal   *azcore.Client
	endpoint   string
	apiVersion string
}

// NewClient creates a client using a caller-provided Microsoft Entra credential.
func NewClient(endpoint string, credential azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if credential == nil {
		return nil, errors.New("credential is required")
	}
	audience := "https://cognitiveservices.azure.com"
	if options != nil && options.Audience != "" {
		audience = strings.TrimRight(options.Audience, "/")
	}
	return newClient(endpoint, runtime.NewBearerTokenPolicy(credential, []string{audience + "/.default"}, nil), options)
}

// NewClientWithKeyCredential creates a client using a resource key.
func NewClientWithKeyCredential(endpoint string, credential *azcore.KeyCredential, options *ClientOptions) (*Client, error) {
	if credential == nil {
		return nil, errors.New("key credential is required")
	}
	return newClient(endpoint, runtime.NewKeyCredentialPolicy(credential, "Ocp-Apim-Subscription-Key", nil), options)
}

func newClient(endpoint string, auth policy.Policy, options *ClientOptions) (*Client, error) {
	u, err := url.Parse(strings.TrimSpace(endpoint))
	if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || (u.Path != "" && u.Path != "/") {
		return nil, errors.New("endpoint must be an absolute HTTPS resource URL without a path, query, credentials or fragment")
	}
	opts := ClientOptions{}
	if options != nil {
		opts = *options
	}
	if opts.APIVersion == "" {
		opts.APIVersion = opts.ClientOptions.APIVersion
	}
	if opts.APIVersion == "" {
		opts.APIVersion = APIVersionGA
	}
	opts.ClientOptions.APIVersion = opts.APIVersion
	client, err := azcore.NewClient(moduleName, version, runtime.PipelineOptions{
		APIVersion: runtime.APIVersionOptions{Name: "api-version"},
		PerCall:    []policy.Policy{runtime.NewRequestIDPolicy()},
		PerRetry:   []policy.Policy{&resourceCredentialPolicy{host: u.Host, auth: auth}},
	}, &opts.ClientOptions)
	if err != nil {
		return nil, err
	}
	return &Client{internal: client, endpoint: strings.TrimRight(u.String(), "/"), apiVersion: opts.APIVersion}, nil
}

type resourceCredentialPolicy struct {
	host string
	auth policy.Policy
}

func (p *resourceCredentialPolicy) Do(req *policy.Request) (*http.Response, error) {
	u := req.Raw().URL
	if u.Scheme != "https" || !strings.EqualFold(u.Host, p.host) || u.User != nil {
		return nil, errors.New("service continuation URL must remain on the HTTPS resource endpoint")
	}
	return p.auth.Do(req)
}

func analyzerPath(analyzerID string) (string, error) {
	if analyzerID == "" {
		return "", errors.New("analyzerID is required")
	}
	return "/analyzers/" + url.PathEscape(analyzerID), nil
}

func (c *Client) request(ctx context.Context, method, path string, query url.Values, body []byte, contentType string) (*policy.Request, error) {
	u, err := url.Parse(c.endpoint + "/contentunderstanding" + path)
	if err != nil {
		return nil, err
	}
	if query == nil {
		query = url.Values{}
	}
	query.Set("api-version", c.apiVersion)
	u.RawQuery = query.Encode()
	req, err := runtime.NewRequest(ctx, method, u.String())
	if err != nil {
		return nil, err
	}
	req.Raw().Header.Set("Accept", "application/json")
	if body != nil {
		if err := req.SetBody(streaming.NopCloser(bytes.NewReader(body)), contentType); err != nil {
			return nil, err
		}
	}
	return req, nil
}

func (c *Client) send(req *policy.Request, codes ...int) (*http.Response, error) {
	resp, err := c.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, codes...) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

func decodeResponse[T any](resp *http.Response) (T, error) {
	var result T
	payload, err := runtime.Payload(resp)
	if err != nil {
		return result, err
	}
	if !bytes.HasPrefix(bytes.TrimSpace(payload), []byte("{")) {
		return result, errors.New("Content Understanding response must be a JSON object")
	}
	if err := json.Unmarshal(payload, &result); err != nil {
		return result, fmt.Errorf("decoding Content Understanding response: %w", err)
	}
	return result, nil
}

func (c *Client) jsonRequest(ctx context.Context, method, path string, query url.Values, body any, contentType string) (*policy.Request, error) {
	var payload []byte
	var err error
	if body != nil {
		payload, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("encoding Content Understanding request: %w", err)
		}
	}
	return c.request(ctx, method, path, query, payload, contentType)
}
