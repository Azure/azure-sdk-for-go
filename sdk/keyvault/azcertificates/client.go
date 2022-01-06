//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcertificates

import (
	"context"
	"net/http"

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

// BeginCreateCertificateResponse contains the response from method Client.BeginCreateCertificate.
type BeginCreateCertificateResponse struct {
	CertificateOperation

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func (c *Client) BeginCreateCertificate(ctx context.Context, certName string, policy CertificatePolicy, options *BeginCreateCertificateOptions) (BeginCreateCertificateResponse, error) {
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
		return BeginCreateCertificateResponse{}, err
	}

	return BeginCreateCertificateResponse{
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
