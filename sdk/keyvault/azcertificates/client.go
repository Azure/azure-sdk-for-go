//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcertificates

import (
	"context"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azcertificates/internal/generated"
	shared "github.com/Azure/azure-sdk-for-go/sdk/keyvault/internal"
)

// Client is the struct for interacting with a Key Vault Certificates instance.
type Client struct {
	genClient *generated.KeyVaultClient
	vaultURL  string
}

// ClientOptions are the optional parameters for the NewClient function
type ClientOptions struct {
	azcore.ClientOptions
}

// converts ClientOptions to generated *generated.ConnectionOptions
func (c *ClientOptions) toConnectionOptions() *policy.ClientOptions {
	if c == nil {
		return &policy.ClientOptions{}
	}

	return &policy.ClientOptions{
		Logging:          c.Logging,
		Retry:            c.Retry,
		Telemetry:        c.Telemetry,
		Transport:        c.Transport,
		PerCallPolicies:  c.PerCallPolicies,
		PerRetryPolicies: c.PerRetryPolicies,
	}
}

// NewClient creates an instance of a Client for a Key Vault Certificate URL.
func NewClient(vaultURL string, credential azcore.TokenCredential, options *ClientOptions) (Client, error) {
	genOptions := options.toConnectionOptions()

	genOptions.PerRetryPolicies = append(
		genOptions.PerRetryPolicies,
		shared.NewKeyVaultChallengePolicy(credential),
	)

	conn := generated.NewConnection(genOptions)

	return Client{
		genClient: generated.NewKeyVaultClient(conn),
		vaultURL:  vaultURL,
	}, nil
}

// Optional parameters for the Client.BeginCreateCertificateOptions function
type BeginCreateCertificateOptions struct {
	// The attributes of the certificate (optional).
	CertificateAttributes *CertificateAttributes `json:"attributes,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]*string `json:"tags,omitempty"`
}

func (b BeginCreateCertificateOptions) toGenerated() *generated.KeyVaultClientCreateCertificateOptions {
	return &generated.KeyVaultClientCreateCertificateOptions{}
}

// CreateCertificateResponse contains the response from method Client.BeginCreateCertificate.
type CreateCertificateResponse struct {
	CertificateOperation

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// CreateCertificatePoller is the interface for the Client.BeginCreateCertificate operation.
type CreateCertificatePoller interface {
	// Done returns true if the LRO has reached a terminal state
	Done() bool

	// Poll fetches the latest state of the LRO. It returns an HTTP response or error.
	// If the LRO has completed successfully, the poller's state is updated and the HTTP response is returned.
	// If the LRO has completed with failure or was cancelled, the poller's state is updated and the error is returned.
	Poll(context.Context) (*http.Response, error)

	// FinalResponse returns the final response after the operations has finished
	FinalResponse(context.Context) (CreateCertificateResponse, error)
}

// the poller returned by the Client.BeginCreateCertificate
type beginCreateCertificatePoller struct {
	certName       string
	certVersion    string
	vaultURL       string
	client         *generated.KeyVaultClient
	createResponse CreateCertificateResponse
	lastResponse   generated.KeyVaultClientGetCertificateResponse
	RawResponse    *http.Response
}

func (b *beginCreateCertificatePoller) Done() bool {
	return b.lastResponse.RawResponse.StatusCode == http.StatusOK
}

func (b *beginCreateCertificatePoller) Poll(ctx context.Context) (*http.Response, error) {
	resp, err := b.client.GetCertificate(ctx, b.vaultURL, b.certName, b.certVersion, nil)
	if err == nil {
		b.lastResponse = resp
		return resp.RawResponse, nil
	}

	if resp.RawResponse.StatusCode == http.StatusNotFound {
		// The certificate has not been fully created yet
		return b.createResponse.RawResponse, nil
	}

	// There was an error in this operation, return the original raw response and the error
	return b.createResponse.RawResponse, err
}

func (b *beginCreateCertificatePoller) FinalResponse(ctx context.Context) (CreateCertificateResponse, error) {
	return b.createResponse, nil
}

func (b *beginCreateCertificatePoller) pollUntilDone(ctx context.Context, t time.Duration) (CreateCertificateResponse, error) {
	for {
		resp, err := b.Poll(ctx)
		if err != nil {
			return CreateCertificateResponse{}, err
		}
		b.RawResponse = resp
		if b.Done() {
			break
		}
		time.Sleep(t)
	}
	return b.createResponse, nil
}

// CreateCertificatePollerResponse contains the response from the Client.BeginCreateCertificate method
type CreateCertificatePollerResponse struct {
	// PollUntilDone will poll the service endpoint until a terminal state is reached or an error occurs
	PollUntilDone func(context.Context, time.Duration) (CreateCertificateResponse, error)

	// Poller contains an initialized WidgetPoller
	Poller CreateCertificatePoller

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func (c *Client) BeginCreateCertificate(ctx context.Context, certName string, policy CertificatePolicy, options *BeginCreateCertificateOptions) (CreateCertificatePollerResponse, error) {
	if options == nil {
		options = &BeginCreateCertificateOptions{}
	}

	resp, err := c.genClient.CreateCertificate(
		ctx,
		c.vaultURL,
		certName,
		generated.CertificateCreateParameters{
			CertificatePolicy:     policy.toGeneratedCertificateCreateParameters(),
			Tags:                  options.Tags,
			CertificateAttributes: options.CertificateAttributes.toGenerated(),
		},
		options.toGenerated(),
	)

	if err != nil {
		return CreateCertificatePollerResponse{}, err
	}

	p := &beginCreateCertificatePoller{
		certName:    certName,
		certVersion: "",
		vaultURL:    c.vaultURL,
		client:      c.genClient,
		createResponse: CreateCertificateResponse{
			RawResponse: resp.RawResponse,
			CertificateOperation: CertificateOperation{
				CancellationRequested: resp.CancellationRequested,
				Csr:                   resp.Csr,
				Error:                 resp.Error,
				IssuerParameters:      resp.IssuerParameters,
				RequestID:             resp.RequestID,
				Status:                resp.Status,
				StatusDetails:         resp.StatusDetails,
				Target:                resp.Target,
				ID:                    resp.ID,
			},
		},
		lastResponse: generated.KeyVaultClientGetCertificateResponse{},
	}

	return CreateCertificatePollerResponse{
		Poller:        p,
		RawResponse:   resp.RawResponse,
		PollUntilDone: p.pollUntilDone,
	}, nil
}

// KeyVaultClientGetCertificateOptions contains the optional parameters for the KeyVaultClient.GetCertificate method.
type GetCertificateOptions struct {
	Version string
}

// GetCertificateResponse contains the result from method Client.GetCertificate.
type GetCertificateResponse struct {
	CertificateBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// GetCertificate - Gets information about a specific certificate. This operation requires the certificates/get permission.
// If the operation fails it returns the *KeyVaultError error type.
func (c *Client) GetCertificate(ctx context.Context, certName string, options *GetCertificateOptions) (GetCertificateResponse, error) {
	if options == nil {
		options = &GetCertificateOptions{}
	}

	resp, err := c.genClient.GetCertificate(ctx, c.vaultURL, certName, options.Version, nil)
	if err != nil {
		return GetCertificateResponse{}, err
	}

	return GetCertificateResponse{
		RawResponse: resp.RawResponse,
		CertificateBundle: CertificateBundle{
			Attributes:     certificateAttributesFromGenerated(resp.Attributes),
			Cer:            resp.Cer,
			ContentType:    resp.ContentType,
			Tags:           resp.Tags,
			ID:             resp.ID,
			Kid:            resp.Kid,
			Policy:         certificatePolicyFromGenerated(resp.Policy),
			Sid:            resp.Sid,
			X509Thumbprint: resp.X509Thumbprint,
		},
	}, nil
}

type GetCertificateVersionOptions struct{}

type GetCertificateOperationResponse struct {
	CertificateOperation

	RawResponse *http.Response
}

func (c *Client) GetCertificateOperation(ctx context.Context, certName string, options *GetCertificateVersionOptions) (GetCertificateOperationResponse, error) {
	if options == nil {
		options = &GetCertificateVersionOptions{}
	}

	resp, err := c.genClient.GetCertificateOperation(ctx, c.vaultURL, certName, &generated.KeyVaultClientGetCertificateOperationOptions{})
	if err != nil {
		return GetCertificateOperationResponse{}, err
	}

	return GetCertificateOperationResponse{
		RawResponse: resp.RawResponse,
		CertificateOperation: CertificateOperation{
			CancellationRequested: resp.CancellationRequested,
			Csr:                   resp.Csr,
			Error:                 resp.Error,
			IssuerParameters:      resp.IssuerParameters,
			RequestID:             resp.RequestID,
			Status:                resp.Status,
			StatusDetails:         resp.StatusDetails,
			Target:                resp.Target,
			ID:                    resp.ID,
		},
	}, nil
}
