package azcertificates

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azcertificates/internal/generated"
)

type BeginRecoverDeletedCertificateOptions struct{}

func (b *BeginRecoverDeletedCertificateOptions) toGenerated() *generated.KeyVaultClientRecoverDeletedCertificateOptions {
	return &generated.KeyVaultClientRecoverDeletedCertificateOptions{}
}

type BeginRecoverDeletedCertificateResponse struct{}

// RecoverDeletedCertificatePoller is the interface for the Client.RecoverDeletedCertificate operation
type RecoverDeletedCertificatePoller interface {
	// Done returns true if the LRO has reached a terminal state
	Done() bool

	// Poll fetches the latest state of the LRO. It returns an HTTP response or error.
	// If the LRO has completed successfully, the poller's state is updated and the HTTP response is returned.
	// If the LRO has completed with failure or was cancelled, the poller's state is updated and the error is returned.
	Poll(context.Context) (*http.Response, error)

	// FinalResponse returns the final response after the operations has finished
	FinalResponse(context.Context) (RecoverDeletedCertificateResponse, error)
}

// beginRecoverPoller implements the RecoverDeletedCertificatePoller interface
type beginRecoverPoller struct {
	certName        string
	vaultUrl        string
	client          *generated.KeyVaultClient
	recoverResponse generated.KeyVaultClientRecoverDeletedCertificateResponse
	lastResponse    generated.KeyVaultClientGetCertificateResponse
	RawResponse     *http.Response
}

// Done returns true when the polling operation is completed
func (b *beginRecoverPoller) Done() bool {
	return b.RawResponse.StatusCode == http.StatusOK
}

// Poll fetches the latest state of the LRO. It returns an HTTP response or error.
// If the LRO has completed successfully, the poller's state is updated and the HTTP response is returned.
// If the LRO has completed with failure or was cancelled, the poller's state is updated and the error is returned.
func (b *beginRecoverPoller) Poll(ctx context.Context) (*http.Response, error) {
	resp, err := b.client.GetCertificate(ctx, b.vaultUrl, b.certName, "", nil)
	b.lastResponse = resp
	var httpErr azcore.HTTPResponse
	if errors.As(err, &httpErr) {
		return httpErr.RawResponse(), err
	}
	return resp.RawResponse, nil
}

// FinalResponse returns the final response after the operations has finished
func (b *beginRecoverPoller) FinalResponse(ctx context.Context) (RecoverDeletedCertificateResponse, error) {
	return recoverDeletedCertificateResponseFromGenerated(b.recoverResponse), nil
}

// pollUntilDone is the method for the Response.PollUntilDone struct
func (b *beginRecoverPoller) pollUntilDone(ctx context.Context, t time.Duration) (RecoverDeletedCertificateResponse, error) {
	for {
		resp, err := b.Poll(ctx)
		if err != nil {
			b.RawResponse = resp
		}
		if b.Done() {
			break
		}
		b.RawResponse = resp
		time.Sleep(t)
	}
	return recoverDeletedCertificateResponseFromGenerated(b.recoverResponse), nil
}

// RecoverDeletedCertificateResponse is the response object for the Client.RecoverDeletedCertificate operation.
type RecoverDeletedCertificateResponse struct {
	CertificateBundle

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// change recover deleted certificate reponse to the generated version.
func recoverDeletedCertificateResponseFromGenerated(i generated.KeyVaultClientRecoverDeletedCertificateResponse) RecoverDeletedCertificateResponse {
	return RecoverDeletedCertificateResponse{
		CertificateBundle: certificateBundleFromGenerated(&i.CertificateBundle),
		RawResponse:       i.RawResponse,
	}
}

// RecoverDeletedCertificatePollerResponse contains the response of the Client.BeginRecoverDeletedCertificate operations
type RecoverDeletedCertificatePollerResponse struct {
	// PollUntilDone will poll the service endpoint until a terminal state is reached or an error occurs
	PollUntilDone func(context.Context, time.Duration) (RecoverDeletedCertificateResponse, error)

	// Poller contains an initialized RecoverDeletedCertificatePoller
	Poller RecoverDeletedCertificatePoller

	// RawResponse cotains the underlying HTTP response
	RawResponse *http.Response
}

// BeginRecoverDeletedCertificate recovers the deleted certificate in the specified vault to the latest version.
// This operation can only be performed on a soft-delete enabled vault. This operation requires the certificates/recover permission.
func (c *Client) BeginRecoverDeletedCertificate(ctx context.Context, certName string, options *BeginRecoverDeletedCertificateOptions) (RecoverDeletedCertificatePollerResponse, error) {
	if options == nil {
		options = &BeginRecoverDeletedCertificateOptions{}
	}
	resp, err := c.genClient.RecoverDeletedCertificate(ctx, c.vaultURL, certName, options.toGenerated())
	if err != nil {
		return RecoverDeletedCertificatePollerResponse{}, err
	}

	getResp, err := c.genClient.GetCertificate(ctx, c.vaultURL, certName, "", nil)
	var httpErr azcore.HTTPResponse
	if errors.As(err, &httpErr) {
		if httpErr.RawResponse().StatusCode != http.StatusNotFound {
			return RecoverDeletedCertificatePollerResponse{}, err
		}
	}

	b := &beginRecoverPoller{
		lastResponse:    getResp,
		certName:        certName,
		client:          c.genClient,
		vaultUrl:        c.vaultURL,
		recoverResponse: resp,
		RawResponse:     getResp.RawResponse,
	}

	return RecoverDeletedCertificatePollerResponse{
		PollUntilDone: b.pollUntilDone,
		Poller:        b,
		RawResponse:   getResp.RawResponse,
	}, nil
}
