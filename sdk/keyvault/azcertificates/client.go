//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcertificates

import (
	"context"
	"errors"
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

// Optional parameters for the Client.BeginCreateCertificate function
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
		b.createResponse.ID = b.lastResponse.ID
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

// GetCertificateOptions contains the optional parameters for the Client.GetCertificate method.
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

// GetCertificateOperationOptions contains the optional parameters for the Client.GetCertificateOperation method.
type GetCertificateOperationOptions struct{}

func (g *GetCertificateOperationOptions) toGenerated() *generated.KeyVaultClientGetCertificateOperationOptions {
	return &generated.KeyVaultClientGetCertificateOperationOptions{}
}

// GetCertificateOperationResponse contains the result from method Client.GetCertificateOperation.
type GetCertificateOperationResponse struct {
	CertificateOperation

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// GetCertificateOperation - Gets the creation operation associated with a specified certificate. This operation requires the certificates/get permission.
// If the operation fails it returns the *KeyVaultError error type.
func (c *Client) GetCertificateOperation(ctx context.Context, certName string, options *GetCertificateOperationOptions) (GetCertificateOperationResponse, error) {
	resp, err := c.genClient.GetCertificateOperation(ctx, c.vaultURL, certName, options.toGenerated())
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

// BeginDeleteCertificateOptions contains the optional parameters for the Client.BeginDeleteCertificate method.
type BeginDeleteCertificateOptions struct{}

// convert public options to generated options struct
func (b *BeginDeleteCertificateOptions) toGenerated() *generated.KeyVaultClientDeleteCertificateOptions {
	return &generated.KeyVaultClientDeleteCertificateOptions{}
}

type DeleteCertificateResponse struct {
	DeletedCertificateBundle

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func deleteCertificateResponseFromGenerated(g *generated.KeyVaultClientDeleteCertificateResponse) DeleteCertificateResponse {
	if g == nil {
		return DeleteCertificateResponse{}
	}
	return DeleteCertificateResponse{
		RawResponse: g.RawResponse,
		DeletedCertificateBundle: DeletedCertificateBundle{
			RecoveryID:         g.RecoveryID,
			DeletedDate:        g.DeletedDate,
			ScheduledPurgeDate: g.ScheduledPurgeDate,
			CertificateBundle: CertificateBundle{
				Attributes:     certificateAttributesFromGenerated(g.Attributes),
				Cer:            g.Cer,
				ContentType:    g.ContentType,
				Tags:           g.Tags,
				ID:             g.ID,
				Kid:            g.Kid,
				Policy:         certificatePolicyFromGenerated(g.Policy),
				Sid:            g.Sid,
				X509Thumbprint: g.X509Thumbprint,
			},
		},
	}
}

// DeleteCertificatePoller is the interface for the Client.DeleteCertificate operation.
type DeleteCertificatePoller interface {
	// Done returns true if the LRO has reached a terminal state
	Done() bool

	// Poll fetches the latest state of the LRO. It returns an HTTP response or error.
	// If the LRO has completed successfully, the poller's state is updated and the HTTP response is returned.
	// If the LRO has completed with failure or was cancelled, the poller's state is updated and the error is returned.
	Poll(context.Context) (*http.Response, error)

	// FinalResponse returns the final response after the operations has finished
	FinalResponse(context.Context) (DeleteCertificateResponse, error)
}

// The poller returned by the Client.BeginDeleteCertificate operation
type beginDeleteCertificatePoller struct {
	certificateName string // This is the certificate to Poll for in GetDeletedKey
	vaultURL        string
	client          *generated.KeyVaultClient
	deleteResponse  generated.KeyVaultClientDeleteCertificateResponse
	lastResponse    generated.KeyVaultClientGetDeletedCertificateResponse
	RawResponse     *http.Response
}

// Done returns true if the LRO has reached a terminal state
func (s *beginDeleteCertificatePoller) Done() bool {
	return s.lastResponse.RawResponse != nil
}

// Poll fetches the latest state of the LRO. It returns an HTTP response or error.(
// If the LRO has completed successfully, the poller's state is updated and the HTTP response is returned.
// If the LRO has completed with failure or was cancelled, the poller's state is updated and the error is returned.)
func (s *beginDeleteCertificatePoller) Poll(ctx context.Context) (*http.Response, error) {
	resp, err := s.client.GetDeletedCertificate(ctx, s.vaultURL, s.certificateName, nil)
	if err == nil {
		// Service recognizes DeletedKey, operation is done
		s.lastResponse = resp
		return resp.RawResponse, nil
	}

	var httpResponseErr azcore.HTTPResponse
	if errors.As(err, &httpResponseErr) {
		if httpResponseErr.RawResponse().StatusCode == http.StatusNotFound {
			// This is the expected result
			return s.deleteResponse.RawResponse, nil
		}
	}
	return s.deleteResponse.RawResponse, err
}

// FinalResponse returns the final response after the operations has finished
func (s *beginDeleteCertificatePoller) FinalResponse(ctx context.Context) (DeleteCertificateResponse, error) {
	return deleteCertificateResponseFromGenerated(&s.deleteResponse), nil
}

// pollUntilDone continually calls the Poll operation until the operation is completed. In between each
// Poll is a wait determined by the t parameter.
func (s *beginDeleteCertificatePoller) pollUntilDone(ctx context.Context, t time.Duration) (DeleteCertificateResponse, error) {
	for {
		resp, err := s.Poll(ctx)
		if err != nil {
			return DeleteCertificateResponse{}, err
		}
		s.RawResponse = resp
		if s.Done() {
			break
		}
		time.Sleep(t)
	}
	return deleteCertificateResponseFromGenerated(&s.deleteResponse), nil
}

// DeleteCertificatePollerResponse contains the response from the Client.BeginDeleteCertificate method
type DeleteCertificatePollerResponse struct {
	// PollUntilDone will poll the service endpoint until a terminal state is reached or an error occurs
	PollUntilDone func(context.Context, time.Duration) (DeleteCertificateResponse, error)

	// Poller contains an initialized WidgetPoller
	Poller DeleteCertificatePoller

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// BeginDeleteCertificate deletes a certificate from the keyvault. Delete cannot be applied to an individual version of a certificate. This operation
// requires the certificate/delete permission. This response contains a Poller struct that can be used to Poll for a response, or the
// response PollUntilDone function can be used to poll until completion.
func (c *Client) BeginDeleteCertificate(ctx context.Context, certificateName string, options *BeginDeleteCertificateOptions) (DeleteCertificatePollerResponse, error) {
	if options == nil {
		options = &BeginDeleteCertificateOptions{}
	}
	resp, err := c.genClient.DeleteCertificate(ctx, c.vaultURL, certificateName, options.toGenerated())
	if err != nil {
		return DeleteCertificatePollerResponse{}, err
	}

	getResp, err := c.genClient.GetDeletedCertificate(ctx, c.vaultURL, certificateName, nil)
	var httpErr azcore.HTTPResponse
	if errors.As(err, &httpErr) {
		if httpErr.RawResponse().StatusCode != http.StatusNotFound {
			return DeleteCertificatePollerResponse{}, err
		}
	}

	s := &beginDeleteCertificatePoller{
		vaultURL:        c.vaultURL,
		certificateName: certificateName,
		client:          c.genClient,
		deleteResponse:  resp,
		lastResponse:    getResp,
	}

	return DeleteCertificatePollerResponse{
		Poller:        s,
		RawResponse:   resp.RawResponse,
		PollUntilDone: s.pollUntilDone,
	}, nil
}

// Optional parameters for the Client.PurgeDeletedCertificateOptions function
type PurgeDeletedCertificateOptions struct{}

func (p *PurgeDeletedCertificateOptions) toGenerated() *generated.KeyVaultClientPurgeDeletedCertificateOptions {
	return &generated.KeyVaultClientPurgeDeletedCertificateOptions{}
}

// PurgeDeletedCertificateResponse contains the response from method Client.PurgeDeletedCertificate.
type PurgeDeletedCertificateResponse struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func (c *Client) PurgeDeletedCertificate(ctx context.Context, certName string, options *PurgeDeletedCertificateOptions) (PurgeDeletedCertificateResponse, error) {
	resp, err := c.genClient.PurgeDeletedCertificate(ctx, c.vaultURL, certName, options.toGenerated())
	if err != nil {
		return PurgeDeletedCertificateResponse{}, err
	}

	return PurgeDeletedCertificateResponse{
		RawResponse: resp.RawResponse,
	}, nil
}

// Optional parameters for the Client.GetDeletedCertificateOptions function
type GetDeletedCertificateOptions struct{}

func (g *GetDeletedCertificateOptions) toGenerated() *generated.KeyVaultClientGetDeletedCertificateOptions {
	return &generated.KeyVaultClientGetDeletedCertificateOptions{}
}

type GetDeletedCertificateResponse struct {
	DeletedCertificateBundle

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func (c *Client) GetDeletedCertificate(ctx context.Context, certName string, options *GetDeletedCertificateOptions) (GetDeletedCertificateResponse, error) {
	resp, err := c.genClient.GetDeletedCertificate(ctx, c.vaultURL, certName, options.toGenerated())
	if err != nil {
		return GetDeletedCertificateResponse{}, err
	}

	return GetDeletedCertificateResponse{
		RawResponse: resp.RawResponse,
		DeletedCertificateBundle: DeletedCertificateBundle{
			RecoveryID:         resp.RecoveryID,
			DeletedDate:        resp.DeletedDate,
			ScheduledPurgeDate: resp.ScheduledPurgeDate,
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
		},
	}, nil
}

// Optional parameters for the Client.BackupCertificateOptions function
type BackupCertificateOptions struct{}

func (b *BackupCertificateOptions) toGenerated() *generated.KeyVaultClientBackupCertificateOptions {
	return &generated.KeyVaultClientBackupCertificateOptions{}
}

// BackupCertificateResponse contains the response from method Client.BackupCertificate.
type BackupCertificateResponse struct {
	// READ-ONLY; The backup blob containing the backed up certificate.
	Value []byte `json:"value,omitempty" azure:"ro"`

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func (c *Client) BackupCertificate(ctx context.Context, certName string, options *BackupCertificateOptions) (BackupCertificateResponse, error) {
	resp, err := c.genClient.BackupCertificate(ctx, c.vaultURL, certName, options.toGenerated())
	if err != nil {
		return BackupCertificateResponse{}, err
	}

	return BackupCertificateResponse{
		RawResponse: resp.RawResponse,
		Value:       resp.Value,
	}, nil
}

type ImportCertificateOptions struct {
	// The attributes of the certificate (optional).
	CertificateAttributes *CertificateAttributes `json:"attributes,omitempty"`

	// The management policy for the certificate.
	CertificatePolicy *CertificatePolicy `json:"policy,omitempty"`

	// If the private key in base64EncodedCertificate is encrypted, the password used for encryption.
	Password *string `json:"pwd,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]*string `json:"tags,omitempty"`
}

func (i *ImportCertificateOptions) toGenerated() *generated.KeyVaultClientImportCertificateOptions {
	return &generated.KeyVaultClientImportCertificateOptions{}
}

type ImportCertificateResponse struct {
	CertificateBundle

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func (c *Client) ImportCertificate(ctx context.Context, certName string, base64EncodedCertificate string, options *ImportCertificateOptions) (ImportCertificateResponse, error) {
	if options == nil {
		options = &ImportCertificateOptions{}
	}
	resp, err := c.genClient.ImportCertificate(
		ctx,
		c.vaultURL,
		certName,
		generated.CertificateImportParameters{
			Base64EncodedCertificate: &base64EncodedCertificate,
			CertificateAttributes:    options.CertificateAttributes.toGenerated(),
			CertificatePolicy:        options.CertificatePolicy.toGeneratedCertificateCreateParameters(),
			Password:                 options.Password,
			Tags:                     options.Tags,
		},
		options.toGenerated(),
	)
	if err != nil {
		return ImportCertificateResponse{}, err
	}

	return ImportCertificateResponse{
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

// ListCertificatesPager is a Pager for the Client.ListCertificates operation
type ListCertificatesPager interface {
	// PageResponse returns the current ListCertificatesPage
	PageResponse() ListCertificatesPage

	// Err returns true if there is another page of data available, false if not
	Err() error

	// NextPage returns true if there is another page of data available, false if not
	NextPage(context.Context) bool
}

// listKeysPager implements the ListCertificatesPager interface
type listKeysPager struct {
	genPager *generated.KeyVaultClientGetCertificatesPager
}

// PageResponse returns the results from the page most recently fetched from the service
func (l *listKeysPager) PageResponse() ListCertificatesPage {
	return listKeysPageFromGenerated(l.genPager.PageResponse())
}

// Err returns an error value if the most recent call to NextPage was not successful, else nil
func (l *listKeysPager) Err() error {
	return l.genPager.Err()
}

// NextPage fetches the next available page of results from the service. If the fetched page
// contains results, the return value is true, else false. Results fetched from the service
// can be evaluated by calling PageResponse on this Pager.
func (l *listKeysPager) NextPage(ctx context.Context) bool {
	return l.genPager.NextPage(ctx)
}

// ListCertificatesOptions contains the optional parameters for the Client.ListCertificates method
type ListCertificatesOptions struct {
	MaxResults *int32
}

// convert ListCertificatesOptions to generated options
func (l *ListCertificatesOptions) toGenerated() *generated.KeyVaultClientGetCertificatesOptions {
	if l == nil {
		return &generated.KeyVaultClientGetCertificatesOptions{}
	}

	return &generated.KeyVaultClientGetCertificatesOptions{Maxresults: l.MaxResults}
}

// ListCertificatesPage contains the current page of results for the Client.ListSecrets operation
type ListCertificatesPage struct {
	// READ-ONLY; A response message containing a list of certificates in the key vault along with a link to the next page of certificates.
	Certificates []*CertificateItem `json:"value,omitempty" azure:"ro"`

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// convert internal Response to ListCertificatesPage
func listKeysPageFromGenerated(i generated.KeyVaultClientGetCertificatesResponse) ListCertificatesPage {
	var vals []*CertificateItem

	for _, v := range i.Value {
		vals = append(vals, &CertificateItem{
			Attributes:     certificateAttributesFromGenerated(v.Attributes),
			ID:             v.ID,
			Tags:           v.Tags,
			X509Thumbprint: v.X509Thumbprint,
		})
	}

	return ListCertificatesPage{
		RawResponse:  i.RawResponse,
		Certificates: vals,
	}
}

// ListCertificates retrieves a list of the certificates in the Key Vault as JSON Web Key structures that contain the
// public part of a stored certificate. The LIST operation is applicable to all certificate types, however only the
// base certificate identifier, attributes, and tags are provided in the response. Individual versions of a
// certificate are not listed in the response. This operation requires the certificates/list permission.
func (c *Client) ListCertificates(options *ListCertificatesOptions) ListCertificatesPager {
	return &listKeysPager{
		genPager: c.genClient.GetCertificates(c.vaultURL, options.toGenerated()),
	}
}

// ListCertificateVersionsPager is a Pager for Client.ListCertificateVersions results
type ListCertificateVersionsPager interface {
	// PageResponse returns the current ListCertificateVersionsPage
	PageResponse() ListCertificateVersionsPage

	// Err returns true if there is another page of data available, false if not
	Err() error

	// NextPage returns true if there is another page of data available, false if not
	NextPage(context.Context) bool
}

// listKeyVersionsPager implements the ListCertificateVersionsPager interface
type listKeyVersionsPager struct {
	genPager *generated.KeyVaultClientGetCertificateVersionsPager
}

// PageResponse returns the results from the page most recently fetched from the service.
func (l *listKeyVersionsPager) PageResponse() ListCertificateVersionsPage {
	return listKeyVersionsPageFromGenerated(l.genPager.PageResponse())
}

// Err returns an error value if the most recent call to NextPage was not successful, else nil.
func (l *listKeyVersionsPager) Err() error {
	return l.genPager.Err()
}

// NextPage fetches the next available page of results from the service. If the fetched page
// contains results, the return value is true, else false. Results fetched from the service
// can be evaluated by calling PageResponse on this Pager.
func (l *listKeyVersionsPager) NextPage(ctx context.Context) bool {
	return l.genPager.NextPage(ctx)
}

// ListCertificateVersionsOptions contains the options for the ListCertificateVersions operations
type ListCertificateVersionsOptions struct {
	// Maximum number of results to return in a page. If not specified the service will return up to 25 results.
	MaxResults *int32
}

// convert the public ListCertificateVersionsOptions to the generated version
func (l *ListCertificateVersionsOptions) toGenerated() *generated.KeyVaultClientGetCertificateVersionsOptions {
	if l == nil {
		return &generated.KeyVaultClientGetCertificateVersionsOptions{}
	}

	return &generated.KeyVaultClientGetCertificateVersionsOptions{
		Maxresults: l.MaxResults,
	}
}

// ListCertificateVersionsPage contains the current page from a ListCertificateVersionsPager.PageResponse method
type ListCertificateVersionsPage struct {
	// READ-ONLY; A response message containing a list of certificates in the key vault along with a link to the next page of certificates.
	Certificates []*CertificateItem `json:"value,omitempty" azure:"ro"`

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// create ListKeysPage from generated pager
func listKeyVersionsPageFromGenerated(i generated.KeyVaultClientGetCertificateVersionsResponse) ListCertificateVersionsPage {
	var vals []*CertificateItem
	for _, v := range i.Value {
		vals = append(vals, &CertificateItem{
			Attributes:     certificateAttributesFromGenerated(v.Attributes),
			ID:             v.ID,
			Tags:           v.Tags,
			X509Thumbprint: v.X509Thumbprint,
		})
	}

	return ListCertificateVersionsPage{
		RawResponse:  i.RawResponse,
		Certificates: vals,
	}
}

// ListCertificateVersions lists all versions of the specified certificate. The full certificate identifer and
// attributes are provided in the response. No values are returned for the certificates. This operation
// requires the certificates/list permission.
func (c *Client) ListCertificateVersions(certificateName string, options *ListCertificateVersionsOptions) ListCertificateVersionsPager {
	return &listKeyVersionsPager{
		genPager: c.genClient.GetCertificateVersions(
			c.vaultURL,
			certificateName,
			options.toGenerated(),
		),
	}
}

// CreateIssuerOptions contains the optional parameters for the Client.CreateIssuer function
type CreateIssuerOptions struct {
	// Attributes of the issuer object.
	Attributes *IssuerAttributes `json:"attributes,omitempty"`

	// The credentials to be used for the issuer.
	Credentials *IssuerCredentials `json:"credentials,omitempty"`

	// Details of the organization as provided to the issuer.
	OrganizationDetails *OrganizationDetails `json:"org_details,omitempty"`
}

func (c *CreateIssuerOptions) toGenerated() *generated.KeyVaultClientSetCertificateIssuerOptions {
	return &generated.KeyVaultClientSetCertificateIssuerOptions{}
}

// CreateIssuerResponse contains the response body for the Client.CreateIssuer function
type CreateIssuerResponse struct {
	IssuerBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// SetCertificateContacts - Sets the certificate contacts for the specified key vault. This operation requires the certificates/managecontacts permission.
func (c *Client) CreateIssuer(ctx context.Context, issuerName string, provider string, options *CreateIssuerOptions) (CreateIssuerResponse, error) {
	if options == nil {
		options = &CreateIssuerOptions{}
	}

	resp, err := c.genClient.SetCertificateIssuer(
		ctx,
		c.vaultURL,
		issuerName,
		generated.CertificateIssuerSetParameters{
			Provider:            &provider,
			Attributes:          options.Attributes.toGenerated(),
			Credentials:         options.Credentials.toGenerated(),
			OrganizationDetails: options.OrganizationDetails.toGenerated(),
		},
		options.toGenerated(),
	)

	if err != nil {
		return CreateIssuerResponse{}, err
	}

	return CreateIssuerResponse{
		RawResponse: resp.RawResponse,
		IssuerBundle: IssuerBundle{
			Attributes:          issuerAttributesFromGenerated(resp.Attributes),
			Credentials:         issuerCredentialsFromGenerated(resp.Credentials),
			OrganizationDetails: organizationDetailsFromGenerated(resp.OrganizationDetails),
			Provider:            resp.Provider,
			ID:                  resp.ID,
		},
	}, nil
}

// GetIssuerOptions contains the optional parameters for the Client.GetIssuer function
type GetIssuerOptions struct{}

func (g *GetIssuerOptions) toGenerated() *generated.KeyVaultClientGetCertificateIssuerOptions {
	return &generated.KeyVaultClientGetCertificateIssuerOptions{}
}

// GetIssuerResponse contains the response from method Client.GetIssuer.
type GetIssuerResponse struct {
	IssuerBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// GetIssuer - The GetIssuer operation returns the specified certificate issuer resources in the specified key vault. This operation
// requires the certificates/manageissuers/getissuers permission.
func (c *Client) GetIssuer(ctx context.Context, issuerName string, options *GetIssuerOptions) (GetIssuerResponse, error) {
	resp, err := c.genClient.GetCertificateIssuer(ctx, c.vaultURL, issuerName, options.toGenerated())
	if err != nil {
		return GetIssuerResponse{}, err
	}

	return GetIssuerResponse{
		RawResponse: resp.RawResponse,
		IssuerBundle: IssuerBundle{
			ID:                  resp.ID,
			Provider:            resp.Provider,
			Attributes:          issuerAttributesFromGenerated(resp.Attributes),
			Credentials:         issuerCredentialsFromGenerated(resp.Credentials),
			OrganizationDetails: organizationDetailsFromGenerated(resp.OrganizationDetails),
		},
	}, nil
}

// ListIssuersPager is a Pager for the Client.ListIssuers operation
type ListIssuersPager interface {
	// PageResponse returns the current ListIssuersPage
	PageResponse() ListIssuersPage

	// Err returns true if there is another page of data available, false if not
	Err() error

	// NextPage returns true if there is another page of data available, false if not
	NextPage(context.Context) bool
}

// listIssuersPager implements the ListIssuersPager interface
type listIssuersPager struct {
	genPager *generated.KeyVaultClientGetCertificateIssuersPager
}

// PageResponse returns the results from the page most recently fetched from the service
func (l *listIssuersPager) PageResponse() ListIssuersPage {
	return listIssuersPageFromGenerated(l.genPager.PageResponse())
}

// Err returns an error value if the most recent call to NextPage was not successful, else nil
func (l *listIssuersPager) Err() error {
	return l.genPager.Err()
}

// NextPage fetches the next available page of results from the service. If the fetched page
// contains results, the return value is true, else false. Results fetched from the service
// can be evaluated by calling PageResponse on this Pager.
func (l *listIssuersPager) NextPage(ctx context.Context) bool {
	return l.genPager.NextPage(ctx)
}

// ListIssuersOptions contains the optional parameters for the Client.ListIssuers method
type ListIssuersOptions struct {
	MaxResults *int32
}

// convert ListIssuersOptions to generated options
func (l *ListIssuersOptions) toGenerated() *generated.KeyVaultClientGetCertificateIssuersOptions {
	if l == nil {
		return &generated.KeyVaultClientGetCertificateIssuersOptions{}
	}

	return &generated.KeyVaultClientGetCertificateIssuersOptions{Maxresults: l.MaxResults}
}

// ListIssuersPage contains the current page of results for the Client.ListSecrets operation
type ListIssuersPage struct {
	// READ-ONLY; A response message containing a list of certificates in the key vault along with a link to the next page of certificates.
	Issuers []*CertificateIssuerItem `json:"value,omitempty" azure:"ro"`

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// convert internal Response to ListIssuersPage
func listIssuersPageFromGenerated(i generated.KeyVaultClientGetCertificateIssuersResponse) ListIssuersPage {
	var vals []*CertificateIssuerItem

	for _, v := range i.Value {
		vals = append(vals, certificateIssuerItemFromGenerated(v))
	}

	return ListIssuersPage{
		RawResponse: i.RawResponse,
		Issuers:     vals,
	}
}

// ListIssuers - The ListIssuers operation returns the set of certificate issuer resources in the specified key vault. This operation
// requires the certificates/manageissuers/getissuers permission.
func (c *Client) ListIssuers(options *ListIssuersOptions) ListIssuersPager {
	return &listIssuersPager{
		genPager: c.genClient.GetCertificateIssuers(c.vaultURL, options.toGenerated()),
	}
}

// DeleteIssuerOptions contains the optional parameters for the Client.DeleteIssuer function
type DeleteIssuerOptions struct{}

func (d *DeleteIssuerOptions) toGenerated() *generated.KeyVaultClientDeleteCertificateIssuerOptions {
	return &generated.KeyVaultClientDeleteCertificateIssuerOptions{}
}

// DeleteIssuerResponse contains the response from method Client.DeleteIssuer.
type DeleteIssuerResponse struct {
	IssuerBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// DeleteIssuer - The DeleteIssuer operation permanently removes the specified certificate issuer from the vault. This operation requires
// the certificates/manageissuers/deleteissuers permission.
func (c *Client) DeleteIssuer(ctx context.Context, issuerName string, options *DeleteIssuerOptions) (DeleteIssuerResponse, error) {
	resp, err := c.genClient.DeleteCertificateIssuer(ctx, c.vaultURL, issuerName, options.toGenerated())
	if err != nil {
		return DeleteIssuerResponse{}, err
	}

	return DeleteIssuerResponse{
		RawResponse: resp.RawResponse,
		IssuerBundle: IssuerBundle{
			Attributes:          issuerAttributesFromGenerated(resp.Attributes),
			Credentials:         issuerCredentialsFromGenerated(resp.Credentials),
			OrganizationDetails: organizationDetailsFromGenerated(resp.OrganizationDetails),
			Provider:            resp.Provider,
			ID:                  resp.ID,
		},
	}, nil
}

// UpdateIssuerOptions contains the optional parameters for the Client.UpdateIssuer function
type UpdateIssuerOptions struct {
	// Attributes of the issuer object.
	Attributes *IssuerAttributes `json:"attributes,omitempty"`

	// The credentials to be used for the issuer.
	Credentials *IssuerCredentials `json:"credentials,omitempty"`

	// Details of the organization as provided to the issuer.
	OrganizationDetails *OrganizationDetails `json:"org_details,omitempty"`

	// The issuer provider.
	Provider *string `json:"provider,omitempty"`
}

func (u *UpdateIssuerOptions) toUpdateParameters() generated.CertificateIssuerUpdateParameters {
	if u == nil {
		return generated.CertificateIssuerUpdateParameters{}
	}

	return generated.CertificateIssuerUpdateParameters{
		Attributes:          u.Attributes.toGenerated(),
		Credentials:         u.Credentials.toGenerated(),
		OrganizationDetails: u.OrganizationDetails.toGenerated(),
		Provider:            u.Provider,
	}
}

// UpdateIssuerResponse contains the response from method Client.UpdateIssuer.
type UpdateIssuerResponse struct {
	IssuerBundle

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// UpdateIssuer - The UpdateIssuer operation performs an update on the specified certificate issuer entity. This operation requires
// the certificates/setissuers permission.
func (c *Client) UpdateIssuer(ctx context.Context, issuerName string, options *UpdateIssuerOptions) (UpdateIssuerResponse, error) {
	resp, err := c.genClient.UpdateCertificateIssuer(
		ctx,
		c.vaultURL,
		issuerName,
		options.toUpdateParameters(),
		&generated.KeyVaultClientUpdateCertificateIssuerOptions{},
	)
	if err != nil {
		return UpdateIssuerResponse{}, err
	}

	return UpdateIssuerResponse{
		RawResponse: resp.RawResponse,
		IssuerBundle: IssuerBundle{
			Attributes:          issuerAttributesFromGenerated(resp.Attributes),
			Credentials:         issuerCredentialsFromGenerated(resp.Credentials),
			OrganizationDetails: organizationDetailsFromGenerated(resp.OrganizationDetails),
			Provider:            resp.Provider,
			ID:                  resp.ID,
		},
	}, nil
}

// SetContactsOptions contains the optional parameters for the Client.CreateContacts function
type SetContactsOptions struct{}

// SetContactsResponse contains the response from method Client.CreateContacts.
type SetContactsResponse struct {
	Contacts
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// SetCertificateContacts - Sets the certificate contacts for the specified key vault. This operation requires the certificates/managecontacts permission.
func (c *Client) SetContacts(ctx context.Context, contacts Contacts, options *SetContactsOptions) (SetContactsResponse, error) {
	resp, err := c.genClient.SetCertificateContacts(
		ctx,
		c.vaultURL,
		contacts.toGenerated(),
		&generated.KeyVaultClientSetCertificateContactsOptions{},
	)

	if err != nil {
		return SetContactsResponse{}, err
	}

	return SetContactsResponse{
		RawResponse: resp.RawResponse,
		Contacts: Contacts{
			ID:          resp.ID,
			ContactList: contactListFromGenerated(resp.ContactList),
		},
	}, nil
}

// GetContactsOptions contains the optional parameters for the Client.GetContacts function
type GetContactsOptions struct{}

func (g *GetContactsOptions) toGenerated() *generated.KeyVaultClientGetCertificateContactsOptions {
	return &generated.KeyVaultClientGetCertificateContactsOptions{}
}

// GetContactsResponse contains the response from method Client.GetContacts.
type GetContactsResponse struct {
	Contacts

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// GetCertificateContacts - The GetCertificateContacts operation returns the set of certificate contact resources in the specified key vault. This operation
// requires the certificates/managecontacts permission.
func (c *Client) GetContacts(ctx context.Context, options *GetContactsOptions) (GetContactsResponse, error) {
	resp, err := c.genClient.GetCertificateContacts(ctx, c.vaultURL, options.toGenerated())
	if err != nil {
		return GetContactsResponse{}, err
	}

	return GetContactsResponse{
		RawResponse: resp.RawResponse,
		Contacts: Contacts{
			ID:          resp.ID,
			ContactList: contactListFromGenerated(resp.ContactList),
		},
	}, nil
}

// DeleteContactsOptions contains the optional parameters for the Client.DeleteContacts function
type DeleteContactsOptions struct{}

func (d *DeleteContactsOptions) toGenerated() *generated.KeyVaultClientDeleteCertificateContactsOptions {
	return &generated.KeyVaultClientDeleteCertificateContactsOptions{}
}

// DeleteContactsResponse contains the response from method Client.DeleteContacts.
type DeleteContactsResponse struct {
	Contacts

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// DeleteContacts - Deletes the certificate contacts for a specified key vault certificate. This operation requires the certificates/managecontacts permission.
func (c *Client) DeleteContacts(ctx context.Context, options *DeleteContactsOptions) (DeleteContactsResponse, error) {
	resp, err := c.genClient.DeleteCertificateContacts(ctx, c.vaultURL, options.toGenerated())
	if err != nil {
		return DeleteContactsResponse{}, err
	}

	return DeleteContactsResponse{
		RawResponse: resp.RawResponse,
		Contacts: Contacts{
			ContactList: contactListFromGenerated(resp.ContactList),
			ID:          resp.ID,
		},
	}, nil
}

// UpdateCertificatePolicyOptions contains the optional parameters for the Client.UpdateCertificatePolicy method.
type UpdateCertificatePolicyOptions struct{}

func (u *UpdateCertificatePolicyOptions) toGenerated() *generated.KeyVaultClientUpdateCertificatePolicyOptions {
	return &generated.KeyVaultClientUpdateCertificatePolicyOptions{}
}

// UpdateCertificatePolicyResponse contains the response from method Client.UpdateCertificatePolicy.
type UpdateCertificatePolicyResponse struct {
	CertificatePolicy

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// UpdateCertificatePolicy - Set specified members in the certificate policy. Leave others as null. This operation requires the certificates/update permission.
func (c *Client) UpdateCertificatePolicy(ctx context.Context, certName string, policy CertificatePolicy, options *UpdateCertificatePolicyOptions) (UpdateCertificatePolicyResponse, error) {
	resp, err := c.genClient.UpdateCertificatePolicy(
		ctx,
		c.vaultURL,
		certName,
		*policy.toGeneratedCertificateCreateParameters(),
		options.toGenerated(),
	)

	if err != nil {
		return UpdateCertificatePolicyResponse{}, err
	}

	return UpdateCertificatePolicyResponse{
		RawResponse:       resp.RawResponse,
		CertificatePolicy: *certificatePolicyFromGenerated(&resp.CertificatePolicy),
	}, nil
}

// GetCertificatePolicyOptions contains the optional parameters for the method Client.GetCertificatePolicy.
type GetCertificatePolicyOptions struct{}

func (g *GetCertificatePolicyOptions) toGenerated() *generated.KeyVaultClientGetCertificatePolicyOptions {
	return &generated.KeyVaultClientGetCertificatePolicyOptions{}
}

// GetCertificatePolicyResponse contains the response from method Client.GetCertificatePolicy.
type GetCertificatePolicyResponse struct {
	CertificatePolicy

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// GetCertificatePolicy - The GetCertificatePolicy operation returns the specified certificate policy resources in the specified key vault. This operation
// requires the certificates/get permission.
func (c *Client) GetCertificatePolicy(ctx context.Context, certName string, options *GetCertificatePolicyOptions) (GetCertificatePolicyResponse, error) {
	resp, err := c.genClient.GetCertificatePolicy(
		ctx,
		c.vaultURL,
		certName,
		options.toGenerated(),
	)
	if err != nil {
		return GetCertificatePolicyResponse{}, err
	}

	return GetCertificatePolicyResponse{
		RawResponse:       resp.RawResponse,
		CertificatePolicy: *certificatePolicyFromGenerated(&resp.CertificatePolicy),
	}, nil
}

// UpdateCertificatePropertiesOptions contains the optional parameters for the Client.UpdateCertificateProperties function
type UpdateCertificatePropertiesOptions struct {
	// The version of the certificate to update
	Version string

	// The attributes of the certificate (optional).
	CertificateAttributes *CertificateAttributes `json:"attributes,omitempty"`

	// The management policy for the certificate.
	CertificatePolicy *CertificatePolicy `json:"policy,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]*string `json:"tags,omitempty"`
}

func (u *UpdateCertificatePropertiesOptions) toGenerated() *generated.KeyVaultClientUpdateCertificateOptions {
	return &generated.KeyVaultClientUpdateCertificateOptions{}
}

// UpdateCertificatePropertiesResponse contains the result from method Client.UpdateCertificateProperties.
type UpdateCertificatePropertiesResponse struct {
	CertificateBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func (c *Client) UpdateCertificateProperties(ctx context.Context, certName string, options *UpdateCertificatePropertiesOptions) (UpdateCertificatePropertiesResponse, error) {
	if options == nil {
		options = &UpdateCertificatePropertiesOptions{}
	}
	resp, err := c.genClient.UpdateCertificate(
		ctx,
		c.vaultURL,
		certName,
		options.Version,
		generated.CertificateUpdateParameters{
			CertificateAttributes: options.CertificateAttributes.toGenerated(),
			CertificatePolicy:     options.CertificatePolicy.toGeneratedCertificateCreateParameters(),
			Tags:                  options.Tags,
		},
		options.toGenerated(),
	)
	if err != nil {
		return UpdateCertificatePropertiesResponse{}, err
	}
	return UpdateCertificatePropertiesResponse{
		RawResponse:       resp.RawResponse,
		CertificateBundle: certificateBundleFromGenerated(&resp.CertificateBundle),
	}, nil
}

// MergeCertificateOptions contains the optional parameters for the Client.MergeCertificate function.
type MergeCertificateOptions struct {
	// The attributes of the certificate (optional).
	CertificateAttributes *CertificateAttributes `json:"attributes,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]*string `json:"tags,omitempty"`
}

func (m *MergeCertificateOptions) toGenerated() *generated.KeyVaultClientMergeCertificateOptions {
	return &generated.KeyVaultClientMergeCertificateOptions{}
}

// MergeCertificateResponse contains the response from method Client.MergeCertificate.
type MergeCertificateResponse struct {
	CertificateBundle

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// The MergeCertificate operation performs the merging of a certificate or certificate chain with a key pair currently available in the service. This operation requires the certificates/create permission.
func (c *Client) MergeCertificate(ctx context.Context, certName string, certificates [][]byte, options *MergeCertificateOptions) (MergeCertificateResponse, error) {
	if options == nil {
		options = &MergeCertificateOptions{}
	}
	resp, err := c.genClient.MergeCertificate(
		ctx, c.vaultURL,
		certName,
		generated.CertificateMergeParameters{
			X509Certificates:      certificates,
			CertificateAttributes: options.CertificateAttributes.toGenerated(),
			Tags:                  options.Tags,
		},
		options.toGenerated(),
	)
	if err != nil {
		return MergeCertificateResponse{}, err
	}

	return MergeCertificateResponse{
		RawResponse:       resp.RawResponse,
		CertificateBundle: certificateBundleFromGenerated(&resp.CertificateBundle),
	}, nil
}

// RestoreCertificateBackupOptions contains the optional parameters for the Client.RestoreCertificateBackup method
type RestoreCertificateBackupOptions struct{}

func (r *RestoreCertificateBackupOptions) toGenerated() *generated.KeyVaultClientRestoreCertificateOptions {
	return &generated.KeyVaultClientRestoreCertificateOptions{}
}

// RestoreCertificateBackupResponse contains the response from method Client.RestoreCertificateBackup
type RestoreCertificateBackupResponse struct {
	CertificateBundle

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// The RecoverDeletedCertificate operation performs the reversal of the Delete operation. The operation is applicable in vaults
// enabled for soft-delete, and must be issued during the retention interval (available in the deleted certificate's attributes).
// This operation requires the certificates/recover permission.
func (c *Client) RestoreCertificateBackup(ctx context.Context, certificateBackup []byte, options *RestoreCertificateBackupOptions) (RestoreCertificateBackupResponse, error) {
	resp, err := c.genClient.RestoreCertificate(
		ctx,
		c.vaultURL,
		generated.CertificateRestoreParameters{CertificateBundleBackup: certificateBackup},
		options.toGenerated(),
	)
	if err != nil {
		return RestoreCertificateBackupResponse{}, err
	}

	return RestoreCertificateBackupResponse{
		RawResponse:       resp.RawResponse,
		CertificateBundle: certificateBundleFromGenerated(&resp.CertificateBundle),
	}, nil
}


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
