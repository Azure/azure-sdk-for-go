// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcontentunderstanding

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

const (
	// ServiceName identifies this service in azcore cloud configurations.
	ServiceName cloud.ServiceName = "contentunderstanding/azcontentunderstanding"
	// APIVersionGA selects the stable service API.
	APIVersionGA = "2025-11-01"
	// APIVersionPreview selects the preview API, including inline analysis.
	APIVersionPreview = "2026-06-01-preview"
)

func init() {
	cloud.AzurePublic.Services[ServiceName] = cloud.ServiceConfiguration{
		Audience: "https://cognitiveservices.azure.com",
	}
}

// ClientOptions configures authentication, transport, retries, logging and service version.
// APIVersion defaults to APIVersionPreview, matching the .NET preview client.
type ClientOptions struct {
	azcore.ClientOptions
}

// NewClient constructs a client with a caller-provided Microsoft Entra credential.
func NewClient(endpoint string, credential azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if credential == nil {
		return nil, errors.New("credential is required")
	}
	opts := clientOptions(options)
	configuredCloud := opts.Cloud
	if reflect.ValueOf(configuredCloud).IsZero() {
		configuredCloud = cloud.AzurePublic
	}
	service, ok := configuredCloud.Services[ServiceName]
	if !ok || service.Audience == "" {
		return nil, fmt.Errorf("cloud must configure an audience for %s", ServiceName)
	}
	auth := runtime.NewBearerTokenPolicy(credential, []string{strings.TrimRight(service.Audience, "/") + "/.default"}, nil)
	return newClient(endpoint, auth, opts)
}

// NewClientWithKeyCredential constructs a client with a resource key.
func NewClientWithKeyCredential(endpoint string, credential *azcore.KeyCredential, options *ClientOptions) (*Client, error) {
	if credential == nil {
		return nil, errors.New("key credential is required")
	}
	return newClient(endpoint, runtime.NewKeyCredentialPolicy(credential, "Ocp-Apim-Subscription-Key", nil), clientOptions(options))
}

func clientOptions(options *ClientOptions) ClientOptions {
	var result ClientOptions
	if options != nil {
		result = *options
	}
	if result.APIVersion == "" {
		result.APIVersion = APIVersionPreview
	}
	return result
}

func newClient(endpoint string, auth policy.Policy, options ClientOptions) (*Client, error) {
	u, err := url.Parse(strings.TrimSpace(endpoint))
	if err != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || (u.Path != "" && u.Path != "/") {
		return nil, errors.New("endpoint must be an absolute HTTPS resource URL without credentials, path, query or fragment")
	}
	internal, err := azcore.NewClient(moduleName, moduleVersion, runtime.PipelineOptions{
		APIVersion: runtime.APIVersionOptions{Name: "api-version"},
		PerCall:    []policy.Policy{runtime.NewRequestIDPolicy(), &clientPolicy{host: u.Host}},
		PerRetry:   []policy.Policy{auth},
	}, &options.ClientOptions)
	if err != nil {
		return nil, err
	}
	return &Client{internal: internal, endpoint: strings.TrimRight(u.String(), "/") + "/contentunderstanding"}, nil
}

type clientPolicy struct {
	host string
}

func (p *clientPolicy) Do(req *policy.Request) (*http.Response, error) {
	u := req.Raw().URL
	if u.Scheme != "https" || !strings.EqualFold(u.Host, p.host) || u.User != nil {
		return nil, errors.New("service continuation must remain on the HTTPS resource endpoint")
	}
	// The emitter currently loses the merge-patch media type on both PATCH operations.
	if req.Raw().Method == http.MethodPatch {
		req.Raw().Header.Set("Content-Type", "application/merge-patch+json")
	}
	binary := strings.HasSuffix(u.Path, ":analyzeBinary") || strings.HasSuffix(u.Path, ":analyzeBinaryInline")
	if binary && req.Raw().Header.Get("Content-Type") == "" {
		req.Raw().Header.Set("Content-Type", "application/octet-stream")
	}
	resp, err := req.Next()
	if err != nil {
		return resp, err
	}
	inline := strings.HasSuffix(u.Path, ":analyzeInline") || strings.HasSuffix(u.Path, ":analyzeBinaryInline")
	if !inline || resp.StatusCode != http.StatusOK {
		return resp, nil
	}
	var envelope struct {
		Status string          `json:"status"`
		Result json.RawMessage `json:"result"`
		Error  struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := runtime.UnmarshalAsJSON(resp, &envelope); err != nil {
		return resp, fmt.Errorf("decoding inline analysis: %w", err)
	}
	if envelope.Status != "Succeeded" {
		code := envelope.Error.Code
		if code == "" {
			code = "UnexpectedOperationStatus"
		}
		return resp, &azcore.ResponseError{StatusCode: resp.StatusCode, ErrorCode: code, RawResponse: resp}
	}
	if len(envelope.Result) == 0 || string(envelope.Result) == "null" {
		return resp, errors.New("successful inline analysis is missing its result")
	}
	return resp, nil
}
