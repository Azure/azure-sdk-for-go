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
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azcertificates/internal/generated"
	shared "github.com/Azure/azure-sdk-for-go/sdk/keyvault/internal"
)

// Client is the struct for interacting with a Key Vault Certificates instance.
// Don't use this type directly, use NewClient() instead.
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
func NewClient(vaultURL string, credential azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	genOptions := options.toConnectionOptions()

	genOptions.PerRetryPolicies = append(
		genOptions.PerRetryPolicies,
		shared.NewKeyVaultChallengePolicy(credential),
	)

	pl := runtime.NewPipeline(generated.ModuleName, generated.ModuleVersion, runtime.PipelineOptions{}, genOptions)

	return &Client{
		genClient: generated.NewKeyVaultClient(pl),
		vaultURL:  vaultURL,
	}, nil
}

// Optional parameters for the Client.BeginCreateCertificate function
type BeginCreateCertificateOptions struct {
	// The attributes of the certificate (optional).
	CertificateAttributes *CertificateProperties `json:"attributes,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]string `json:"tags,omitempty"`
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

// the poller returned by the Client.BeginCreateCertificate
type CreateCertificatePoller struct {
	certName       string
	certVersion    string
	vaultURL       string
	client         *generated.KeyVaultClient
	createResponse CreateCertificateResponse
	lastResponse   generated.KeyVaultClientGetCertificateResponse
	RawResponse    *http.Response
}

// Done returns true if the LRO has reached a terminal state
func (b *CreateCertificatePoller) Done() bool {
	return b.lastResponse.RawResponse.StatusCode == http.StatusOK
}

// Poll fetches the latest state of the operations. It returns an HTTP response or error.
// If the LRO has completed successfully, the poller's state is updated and the HTTP response is returned.
// If the LRO has completed with failure or was cancelled, the poller's state is updated and the error is returned.)
func (b *CreateCertificatePoller) Poll(ctx context.Context) (*http.Response, error) {
	resp, err := b.client.GetCertificate(ctx, b.vaultURL, b.certName, b.certVersion, nil)
	if err == nil {
		b.lastResponse = resp
		b.createResponse.ID = b.lastResponse.ID
		return resp.RawResponse, nil
	}

	var respErr *azcore.ResponseError
	if errors.As(err, &respErr) {
		if respErr.RawResponse.StatusCode == http.StatusNotFound {
			// The certificate has not been fully created yet
			return resp.RawResponse, nil
		}
	}

	// There was an error in this operation, return the original raw response and the error
	return b.createResponse.RawResponse, err
}

// FinalResponse returns the final response after the operations has finished
func (b *CreateCertificatePoller) FinalResponse(ctx context.Context) (CreateCertificateResponse, error) {
	return b.createResponse, nil
}

// pollUntilDone continuallys polls the service with a 't' delay until completion.
func (b *CreateCertificatePoller) pollUntilDone(ctx context.Context, t time.Duration) (CreateCertificateResponse, error) {
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

// BeginCreateCertificate creates a new certificate resource, if a certificate with this name already exists, a new version is created. This operation requires the certificates/create permission.
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
			Tags:                  convertToGeneratedMap(options.Tags),
			CertificateAttributes: options.CertificateAttributes.toGenerated(),
		},
		options.toGenerated(),
	)

	if err != nil {
		return CreateCertificatePollerResponse{}, err
	}

	p := CreateCertificatePoller{
		certName:    certName,
		certVersion: "",
		vaultURL:    c.vaultURL,
		client:      c.genClient,
		createResponse: CreateCertificateResponse{
			RawResponse: resp.RawResponse,
			CertificateOperation: CertificateOperation{
				CancellationRequested: resp.CancellationRequested,
				Csr:                   resp.Csr,
				Error:                 certificateErrorFromGenerated(resp.Error),
				IssuerParameters:      issuerParametersFromGenerated(resp.IssuerParameters),
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
	KeyVaultCertificateWithPolicy

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// GetCertificate gets information about a specific certificate. This operation requires the certificates/get permission.
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
		KeyVaultCertificateWithPolicy: KeyVaultCertificateWithPolicy{
			Properties:     certificateAttributesFromGenerated(resp.Attributes),
			Cer:            resp.Cer,
			ContentType:    resp.ContentType,
			Tags:           convertGeneratedMap(resp.Tags),
			ID:             resp.ID,
			KeyID:          resp.Kid,
			SecretID:       resp.Sid,
			X509Thumbprint: resp.X509Thumbprint,
			Policy:         certificatePolicyFromGenerated(resp.Policy),
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

// GetCertificateOperation gets the creation operation associated with a specified certificate. This operation requires the certificates/get permission.
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
			Error:                 certificateErrorFromGenerated(resp.Error),
			IssuerParameters:      issuerParametersFromGenerated(resp.IssuerParameters),
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

// DeleteCertificateResponse contains the response structure for the BeginDeleteCertificatePoller.FinalResponse function
type DeleteCertificateResponse struct {
	DeletedCertificate

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func deleteCertificateResponseFromGenerated(g *generated.KeyVaultClientDeleteCertificateResponse) DeleteCertificateResponse {
	if g == nil {
		return DeleteCertificateResponse{}
	}
	return DeleteCertificateResponse{
		RawResponse: g.RawResponse,
		DeletedCertificate: DeletedCertificate{
			RecoveryID:         g.RecoveryID,
			DeletedDate:        g.DeletedDate,
			ScheduledPurgeDate: g.ScheduledPurgeDate,
			Properties:         certificateAttributesFromGenerated(g.Attributes),
			Cer:                g.Cer,
			ContentType:        g.ContentType,
			Tags:               convertGeneratedMap(g.Tags),
			ID:                 g.ID,
			KeyID:              g.Kid,
			Policy:             certificatePolicyFromGenerated(g.Policy),
			SecretID:           g.Sid,
			X509Thumbprint:     g.X509Thumbprint,
		},
	}
}

// The poller returned by the Client.BeginDeleteCertificate operation
type DeleteCertificatePoller struct {
	certificateName string // This is the certificate to Poll for in GetDeletedCertificate
	vaultURL        string
	client          *generated.KeyVaultClient
	deleteResponse  generated.KeyVaultClientDeleteCertificateResponse
	lastResponse    generated.KeyVaultClientGetDeletedCertificateResponse
	RawResponse     *http.Response
}

// Done returns true if the LRO has reached a terminal state
func (s *DeleteCertificatePoller) Done() bool {
	return s.lastResponse.RawResponse != nil
}

// Poll fetches the latest state of the LRO. It returns an HTTP response or error.(
// If the LRO has completed successfully, the poller's state is updated and the HTTP response is returned.
// If the LRO has completed with failure or was cancelled, the poller's state is updated and the error is returned.)
func (s *DeleteCertificatePoller) Poll(ctx context.Context) (*http.Response, error) {
	resp, err := s.client.GetDeletedCertificate(ctx, s.vaultURL, s.certificateName, nil)
	if err == nil {
		// Service recognizes DeletedKey, operation is done
		s.lastResponse = resp
		return resp.RawResponse, nil
	}

	var httpResponseErr *azcore.ResponseError
	if errors.As(err, &httpResponseErr) {
		if httpResponseErr.RawResponse.StatusCode == http.StatusNotFound {
			// This is the expected result
			return s.deleteResponse.RawResponse, nil
		}
	}
	return s.deleteResponse.RawResponse, err
}

// FinalResponse returns the final response after the operations has finished
func (s *DeleteCertificatePoller) FinalResponse(ctx context.Context) (DeleteCertificateResponse, error) {
	return deleteCertificateResponseFromGenerated(&s.deleteResponse), nil
}

// pollUntilDone continually calls the Poll operation until the operation is completed. In between each
// Poll is a wait determined by the t parameter.
func (s *DeleteCertificatePoller) pollUntilDone(ctx context.Context, t time.Duration) (DeleteCertificateResponse, error) {
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
// requires the certificate/delete permission. This response contains a response with a Poller struct that can be used to Poll for a response, or the
// DeleteCertificatePollerResponse.PollUntilDone function can be used to poll until completion.
func (c *Client) BeginDeleteCertificate(ctx context.Context, certificateName string, options *BeginDeleteCertificateOptions) (DeleteCertificatePollerResponse, error) {
	if options == nil {
		options = &BeginDeleteCertificateOptions{}
	}
	resp, err := c.genClient.DeleteCertificate(ctx, c.vaultURL, certificateName, options.toGenerated())
	if err != nil {
		return DeleteCertificatePollerResponse{}, err
	}

	getResp, err := c.genClient.GetDeletedCertificate(ctx, c.vaultURL, certificateName, nil)
	var httpErr *azcore.ResponseError
	if errors.As(err, &httpErr) {
		if httpErr.RawResponse.StatusCode != http.StatusNotFound {
			return DeleteCertificatePollerResponse{}, err
		}
	}

	s := DeleteCertificatePoller{
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

// PurgeDeletedCertificate operation performs an irreversible deletion of the specified certificate, without possibility for recovery. The operation
// is not available if the recovery level does not specify 'Purgeable'. This operation requires the certificate/purge permission.
func (c *Client) PurgeDeletedCertificate(ctx context.Context, certName string, options *PurgeDeletedCertificateOptions) (PurgeDeletedCertificateResponse, error) {
	resp, err := c.genClient.PurgeDeletedCertificate(ctx, c.vaultURL, certName, options.toGenerated())
	if err != nil {
		return PurgeDeletedCertificateResponse{}, err
	}

	return PurgeDeletedCertificateResponse{
		RawResponse: resp.RawResponse,
	}, nil
}

// Optional parameters for the Client.GetDeletedCertificate function
type GetDeletedCertificateOptions struct{}

func (g *GetDeletedCertificateOptions) toGenerated() *generated.KeyVaultClientGetDeletedCertificateOptions {
	return &generated.KeyVaultClientGetDeletedCertificateOptions{}
}

// GetDeletedCertificateResponse is the response struct for the Client.GetDeletedCertificate function.
type GetDeletedCertificateResponse struct {
	DeletedCertificate

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// GetDeletedCertificate retrieves the deleted certificate information plus its attributes, such as retention interval, scheduled permanent deletion
// and the current deletion recovery level. This operation requires the certificates/get permission.
func (c *Client) GetDeletedCertificate(ctx context.Context, certName string, options *GetDeletedCertificateOptions) (GetDeletedCertificateResponse, error) {
	resp, err := c.genClient.GetDeletedCertificate(ctx, c.vaultURL, certName, options.toGenerated())
	if err != nil {
		return GetDeletedCertificateResponse{}, err
	}

	return GetDeletedCertificateResponse{
		RawResponse: resp.RawResponse,
		DeletedCertificate: DeletedCertificate{
			RecoveryID:         resp.RecoveryID,
			DeletedDate:        resp.DeletedDate,
			ScheduledPurgeDate: resp.ScheduledPurgeDate,
			Properties:         certificateAttributesFromGenerated(resp.Attributes),
			Cer:                resp.Cer,
			ContentType:        resp.ContentType,
			Tags:               convertGeneratedMap(resp.Tags),
			ID:                 resp.ID,
			KeyID:              resp.Kid,
			Policy:             certificatePolicyFromGenerated(resp.Policy),
			SecretID:           resp.Sid,
			X509Thumbprint:     resp.X509Thumbprint,
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

// BackupCertificate requests that a backup of the specified certificate be downloaded to the client. All versions of the certificate will be downloaded.
// This operation requires the certificates/backup permission.
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

// ImportCertificateOptions contains the optional parameters for the Client.ImportCertificate function.
type ImportCertificateOptions struct {
	// The attributes of the certificate (optional).
	CertificateAttributes *CertificateProperties `json:"attributes,omitempty"`

	// The management policy for the certificate.
	CertificatePolicy *CertificatePolicy `json:"policy,omitempty"`

	// If the private key in base64EncodedCertificate is encrypted, the password used for encryption.
	Password *string `json:"pwd,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]string `json:"tags,omitempty"`
}

func (i *ImportCertificateOptions) toGenerated() *generated.KeyVaultClientImportCertificateOptions {
	return &generated.KeyVaultClientImportCertificateOptions{}
}

// ImportCertificateResponse is the response struct for the Client.ImportCertificate function.
type ImportCertificateResponse struct {
	KeyVaultCertificateWithPolicy

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// ImportCertificate imports an existing valid certificate, containing a private key, into Azure Key Vault. This operation requires the
// certificates/import permission. The certificate to be imported can be in either PFX or PEM format. If the certificate is in PEM format
// the PEM file must contain the key as well as x509 certificates. Key Vault will only accept a key in PKCS#8 format.
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
			Tags:                     convertToGeneratedMap(options.Tags),
		},
		options.toGenerated(),
	)
	if err != nil {
		return ImportCertificateResponse{}, err
	}

	return ImportCertificateResponse{
		RawResponse: resp.RawResponse,
		KeyVaultCertificateWithPolicy: KeyVaultCertificateWithPolicy{
			Properties:     certificateAttributesFromGenerated(resp.Attributes),
			Cer:            resp.Cer,
			ContentType:    resp.ContentType,
			Tags:           convertGeneratedMap(resp.Tags),
			ID:             resp.ID,
			KeyID:          resp.Kid,
			SecretID:       resp.Sid,
			X509Thumbprint: resp.X509Thumbprint,
			Policy:         certificatePolicyFromGenerated(resp.Policy),
		},
	}, nil
}

// ListCertificatesPager implements the ListCertificatesPager interface
type ListCertificatesPager struct {
	genPager *generated.KeyVaultClientGetCertificatesPager
}

// PageResponse returns the results from the page most recently fetched from the service
func (l *ListCertificatesPager) PageResponse() ListCertificatesPage {
	return listKeysPageFromGenerated(l.genPager.PageResponse())
}

// Err returns an error value if the most recent call to NextPage was not successful, else nil
func (l *ListCertificatesPager) Err() error {
	return l.genPager.Err()
}

// NextPage fetches the next available page of results from the service. If the fetched page
// contains results, the return value is true, else false. Results fetched from the service
// can be evaluated by calling PageResponse on this Pager.
func (l *ListCertificatesPager) NextPage(ctx context.Context) bool {
	return l.genPager.NextPage(ctx)
}

// ListCertificatesOptions contains the optional parameters for the Client.ListCertificates method
type ListCertificatesOptions struct{}

// convert ListCertificatesOptions to generated options
func (l *ListCertificatesOptions) toGenerated() *generated.KeyVaultClientGetCertificatesOptions {
	if l == nil {
		return &generated.KeyVaultClientGetCertificatesOptions{}
	}

	return &generated.KeyVaultClientGetCertificatesOptions{}
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
			Properties:     certificateAttributesFromGenerated(v.Attributes),
			ID:             v.ID,
			Tags:           convertGeneratedMap(v.Tags),
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
	return ListCertificatesPager{
		genPager: c.genClient.GetCertificates(c.vaultURL, options.toGenerated()),
	}
}

// ListCertificateVersionsPager is the pager returned by Client.ListCertificateVersions
type ListCertificateVersionsPager struct {
	genPager *generated.KeyVaultClientGetCertificateVersionsPager
}

// PageResponse returns the results from the page most recently fetched from the service.
func (l *ListCertificateVersionsPager) PageResponse() ListCertificateVersionsPage {
	return listKeyVersionsPageFromGenerated(l.genPager.PageResponse())
}

// Err returns an error value if the most recent call to NextPage was not successful, else nil.
func (l *ListCertificateVersionsPager) Err() error {
	return l.genPager.Err()
}

// NextPage fetches the next available page of results from the service. If the fetched page
// contains results, the return value is true, else false. Results fetched from the service
// can be evaluated by calling PageResponse on this Pager.
func (l *ListCertificateVersionsPager) NextPage(ctx context.Context) bool {
	return l.genPager.NextPage(ctx)
}

// ListCertificateVersionsOptions contains the options for the ListCertificateVersions operations
type ListCertificateVersionsOptions struct{}

// convert the public ListCertificateVersionsOptions to the generated version
func (l *ListCertificateVersionsOptions) toGenerated() *generated.KeyVaultClientGetCertificateVersionsOptions {
	if l == nil {
		return &generated.KeyVaultClientGetCertificateVersionsOptions{}
	}

	return &generated.KeyVaultClientGetCertificateVersionsOptions{}
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
			Properties:     certificateAttributesFromGenerated(v.Attributes),
			ID:             v.ID,
			Tags:           convertGeneratedMap(v.Tags),
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
	return ListCertificateVersionsPager{
		genPager: c.genClient.GetCertificateVersions(
			c.vaultURL,
			certificateName,
			options.toGenerated(),
		),
	}
}

// CreateIssuerOptions contains the optional parameters for the Client.CreateIssuer function
type CreateIssuerOptions struct {
	// Determines whether the issuer is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// The credentials to be used for the issuer.
	Credentials *IssuerCredentials `json:"credentials,omitempty"`

	// Details of the organization administrator.
	AdministratorContacts []*AdministratorContact `json:"admin_details,omitempty"`

	// Id of the organization.
	OrganizationID *string `json:"id,omitempty"`
}

func (c *CreateIssuerOptions) toGenerated() *generated.KeyVaultClientSetCertificateIssuerOptions {
	return &generated.KeyVaultClientSetCertificateIssuerOptions{}
}

// CreateIssuerResponse is the response struct for the Client.CreateIssuer function
type CreateIssuerResponse struct {
	CertificateIssuer
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// CreateIssuer adds or updates the specified certificate issuer. This operation requires the certificates/setissuers permission.
func (c *Client) CreateIssuer(ctx context.Context, issuerName string, provider string, options *CreateIssuerOptions) (CreateIssuerResponse, error) {
	if options == nil {
		options = &CreateIssuerOptions{}
	}

	var orgDetails *generated.OrganizationDetails
	if options.AdministratorContacts != nil || options.OrganizationID != nil {
		orgDetails = &generated.OrganizationDetails{}
		if options.OrganizationID != nil {
			orgDetails.ID = options.OrganizationID
		}

		if options.AdministratorContacts != nil {
			a := make([]*generated.AdministratorDetails, len(options.AdministratorContacts))
			for idx, v := range options.AdministratorContacts {
				a[idx] = &generated.AdministratorDetails{
					EmailAddress: v.EmailAddress,
					FirstName:    v.FirstName,
					LastName:     v.LastName,
					Phone:        v.Phone,
				}
			}
			orgDetails.AdminDetails = a
		}
	}

	resp, err := c.genClient.SetCertificateIssuer(
		ctx,
		c.vaultURL,
		issuerName,
		generated.CertificateIssuerSetParameters{
			Provider:            &provider,
			Attributes:          &generated.IssuerAttributes{Enabled: options.Enabled},
			Credentials:         options.Credentials.toGenerated(),
			OrganizationDetails: orgDetails,
		},
		options.toGenerated(),
	)

	if err != nil {
		return CreateIssuerResponse{}, err
	}

	cr := CreateIssuerResponse{RawResponse: resp.RawResponse}
	cr.CertificateIssuer = CertificateIssuer{
		Credentials: issuerCredentialsFromGenerated(resp.Credentials),
		Provider:    resp.Provider,
		ID:          resp.ID,
	}

	if resp.Attributes != nil {
		cr.CertificateIssuer.Created = resp.Attributes.Created
		cr.CertificateIssuer.Enabled = resp.Attributes.Enabled
		cr.CertificateIssuer.Updated = resp.Attributes.Updated
	}
	if resp.OrganizationDetails != nil {
		cr.OrganizationID = resp.OrganizationDetails.ID
		var adminDetails []*AdministratorContact
		if resp.OrganizationDetails.AdminDetails != nil {
			adminDetails = make([]*AdministratorContact, len(resp.OrganizationDetails.AdminDetails))
			for idx, v := range resp.OrganizationDetails.AdminDetails {
				adminDetails[idx] = &AdministratorContact{
					EmailAddress: v.EmailAddress,
					FirstName:    v.FirstName,
					LastName:     v.LastName,
					Phone:        v.Phone,
				}
			}
		}
		cr.AdministratorContacts = adminDetails
	}

	return cr, nil
}

// GetIssuerOptions contains the optional parameters for the Client.GetIssuer function
type GetIssuerOptions struct{}

func (g *GetIssuerOptions) toGenerated() *generated.KeyVaultClientGetCertificateIssuerOptions {
	return &generated.KeyVaultClientGetCertificateIssuerOptions{}
}

// GetIssuerResponse contains the response from method Client.GetIssuer.
type GetIssuerResponse struct {
	CertificateIssuer
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// GetIssuer returns the specified certificate issuer resources in the specified key vault. This operation
// requires the certificates/manageissuers/getissuers permission.
func (c *Client) GetIssuer(ctx context.Context, issuerName string, options *GetIssuerOptions) (GetIssuerResponse, error) {
	resp, err := c.genClient.GetCertificateIssuer(ctx, c.vaultURL, issuerName, options.toGenerated())
	if err != nil {
		return GetIssuerResponse{}, err
	}

	g := GetIssuerResponse{RawResponse: resp.RawResponse}
	g.CertificateIssuer = CertificateIssuer{
		ID:          resp.ID,
		Provider:    resp.Provider,
		Credentials: issuerCredentialsFromGenerated(resp.Credentials),
	}

	if resp.Attributes != nil {
		g.CertificateIssuer.Created = resp.Attributes.Created
		g.CertificateIssuer.Enabled = resp.Attributes.Enabled
		g.CertificateIssuer.Updated = resp.Attributes.Updated
	}
	if resp.OrganizationDetails != nil {
		g.OrganizationID = resp.OrganizationDetails.ID
		var adminDetails []*AdministratorContact
		if resp.OrganizationDetails.AdminDetails != nil {
			adminDetails = make([]*AdministratorContact, len(resp.OrganizationDetails.AdminDetails))
			for idx, v := range resp.OrganizationDetails.AdminDetails {
				adminDetails[idx] = &AdministratorContact{
					EmailAddress: v.EmailAddress,
					FirstName:    v.FirstName,
					LastName:     v.LastName,
					Phone:        v.Phone,
				}
			}
		}
		g.AdministratorContacts = adminDetails
	}

	return g, nil
}

// ListPropertiesOfIssuersPager is the pager returned by Client.ListIssuers
type ListPropertiesOfIssuersPager struct {
	genPager *generated.KeyVaultClientGetCertificateIssuersPager
}

// PageResponse returns the results from the page most recently fetched from the service
func (l *ListPropertiesOfIssuersPager) PageResponse() ListIssuersPropertiesOfIssuersPage {
	return listIssuersPageFromGenerated(l.genPager.PageResponse())
}

// Err returns an error value if the most recent call to NextPage was not successful, else nil
func (l *ListPropertiesOfIssuersPager) Err() error {
	return l.genPager.Err()
}

// NextPage fetches the next available page of results from the service. If the fetched page
// contains results, the return value is true, else false. Results fetched from the service
// can be evaluated by calling PageResponse on this Pager.
func (l *ListPropertiesOfIssuersPager) NextPage(ctx context.Context) bool {
	return l.genPager.NextPage(ctx)
}

// ListPropertiesOfIssuersOptions contains the optional parameters for the Client.ListIssuers method
type ListPropertiesOfIssuersOptions struct{}

// convert ListIssuersOptions to generated options
func (l *ListPropertiesOfIssuersOptions) toGenerated() *generated.KeyVaultClientGetCertificateIssuersOptions {
	if l == nil {
		return &generated.KeyVaultClientGetCertificateIssuersOptions{}
	}

	return &generated.KeyVaultClientGetCertificateIssuersOptions{}
}

// ListIssuersPropertiesOfIssuersPage contains the current page of results for the Client.ListSecrets operation
type ListIssuersPropertiesOfIssuersPage struct {
	// READ-ONLY; A response message containing a list of certificates in the key vault along with a link to the next page of certificates.
	Issuers []*CertificateIssuerItem `json:"value,omitempty" azure:"ro"`

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// convert internal Response to ListPropertiesOfIssuersPage
func listIssuersPageFromGenerated(i generated.KeyVaultClientGetCertificateIssuersResponse) ListIssuersPropertiesOfIssuersPage {
	var vals []*CertificateIssuerItem

	for _, v := range i.Value {
		vals = append(vals, certificateIssuerItemFromGenerated(v))
	}

	return ListIssuersPropertiesOfIssuersPage{
		RawResponse: i.RawResponse,
		Issuers:     vals,
	}
}

// ListPropertiesOfIssuers returns a pager that can be used to get the set of certificate issuer resources in the specified key vault. This operation
// requires the certificates/manageissuers/getissuers permission.
func (c *Client) ListPropertiesOfIssuers(options *ListPropertiesOfIssuersOptions) ListPropertiesOfIssuersPager {
	return ListPropertiesOfIssuersPager{
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
	CertificateIssuer
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// DeleteIssuer permanently removes the specified certificate issuer from the vault. This operation requires the certificates/manageissuers/deleteissuers permission.
func (c *Client) DeleteIssuer(ctx context.Context, issuerName string, options *DeleteIssuerOptions) (DeleteIssuerResponse, error) {
	resp, err := c.genClient.DeleteCertificateIssuer(ctx, c.vaultURL, issuerName, options.toGenerated())
	if err != nil {
		return DeleteIssuerResponse{}, err
	}

	d := DeleteIssuerResponse{RawResponse: resp.RawResponse}
	d.CertificateIssuer = CertificateIssuer{
		ID:          resp.ID,
		Provider:    resp.Provider,
		Credentials: issuerCredentialsFromGenerated(resp.Credentials),
	}

	if resp.Attributes != nil {
		d.CertificateIssuer.Created = resp.Attributes.Created
		d.CertificateIssuer.Enabled = resp.Attributes.Enabled
		d.CertificateIssuer.Updated = resp.Attributes.Updated
	}
	if resp.OrganizationDetails != nil {
		d.OrganizationID = resp.OrganizationDetails.ID
		var adminDetails []*AdministratorContact
		if resp.OrganizationDetails.AdminDetails != nil {
			adminDetails = make([]*AdministratorContact, len(resp.OrganizationDetails.AdminDetails))
			for idx, v := range resp.OrganizationDetails.AdminDetails {
				adminDetails[idx] = &AdministratorContact{
					EmailAddress: v.EmailAddress,
					FirstName:    v.FirstName,
					LastName:     v.LastName,
					Phone:        v.Phone,
				}
			}
		}
		d.AdministratorContacts = adminDetails
	}

	return d, nil
}

// UpdateIssuerOptions contains the optional parameters for the Client.UpdateIssuer function
type UpdateIssuerOptions struct {
	// Determines whether the issuer is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// The credentials to be used for the issuer.
	Credentials *IssuerCredentials `json:"credentials,omitempty"`

	// Details of the organization administrator.
	AdministratorContacts []*AdministratorContact `json:"admin_details,omitempty"`

	// Id of the organization.
	OrganizationID *string `json:"id,omitempty"`

	// The issuer provider.
	Provider *string `json:"provider,omitempty"`
}

func (u *UpdateIssuerOptions) toUpdateParameters() generated.CertificateIssuerUpdateParameters {
	if u == nil {
		return generated.CertificateIssuerUpdateParameters{}
	}
	var attrib *generated.IssuerAttributes
	if u.Enabled != nil {
		attrib = &generated.IssuerAttributes{Enabled: u.Enabled}
	}

	var orgDetail *generated.OrganizationDetails
	if u.OrganizationID != nil || u.AdministratorContacts != nil {
		orgDetail = &generated.OrganizationDetails{}
		if u.OrganizationID != nil {
			orgDetail.ID = u.OrganizationID
		}

		if u.AdministratorContacts != nil {
			a := make([]*generated.AdministratorDetails, len(u.AdministratorContacts))
			for idx, v := range u.AdministratorContacts {
				a[idx] = &generated.AdministratorDetails{
					EmailAddress: v.EmailAddress,
					FirstName:    v.FirstName,
					LastName:     v.LastName,
					Phone:        v.Phone,
				}
			}

			orgDetail.AdminDetails = a
		}
	}

	return generated.CertificateIssuerUpdateParameters{
		Attributes:          attrib,
		Credentials:         u.Credentials.toGenerated(),
		OrganizationDetails: orgDetail,
		Provider:            u.Provider,
	}
}

// UpdateIssuerResponse contains the response from method Client.UpdateIssuer.
type UpdateIssuerResponse struct {
	CertificateIssuer

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// UpdateIssuer performs an update on the specified certificate issuer entity. This operation requires
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

	u := UpdateIssuerResponse{RawResponse: resp.RawResponse}
	u.CertificateIssuer = CertificateIssuer{
		ID:          resp.ID,
		Provider:    resp.Provider,
		Credentials: issuerCredentialsFromGenerated(resp.Credentials),
	}

	if resp.Attributes != nil {
		u.CertificateIssuer.Created = resp.Attributes.Created
		u.CertificateIssuer.Enabled = resp.Attributes.Enabled
		u.CertificateIssuer.Updated = resp.Attributes.Updated
	}
	if resp.OrganizationDetails != nil {
		u.OrganizationID = resp.OrganizationDetails.ID
		var adminDetails []*AdministratorContact
		if resp.OrganizationDetails.AdminDetails != nil {
			adminDetails = make([]*AdministratorContact, len(resp.OrganizationDetails.AdminDetails))
			for idx, v := range resp.OrganizationDetails.AdminDetails {
				adminDetails[idx] = &AdministratorContact{
					EmailAddress: v.EmailAddress,
					FirstName:    v.FirstName,
					LastName:     v.LastName,
					Phone:        v.Phone,
				}
			}
		}
		u.AdministratorContacts = adminDetails
	}

	return u, nil
}

// SetContactsOptions contains the optional parameters for the Client.CreateContacts function
type SetContactsOptions struct{}

func (s *SetContactsOptions) toGenerated() *generated.KeyVaultClientSetCertificateContactsOptions {
	return &generated.KeyVaultClientSetCertificateContactsOptions{}
}

// SetContactsResponse contains the response from method Client.CreateContacts.
type SetContactsResponse struct {
	Contacts
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// SetCertificateContacts sets the certificate contacts for the specified key vault. This operation requires the certificates/managecontacts permission.
func (c *Client) SetContacts(ctx context.Context, contacts Contacts, options *SetContactsOptions) (SetContactsResponse, error) {
	resp, err := c.genClient.SetCertificateContacts(
		ctx,
		c.vaultURL,
		contacts.toGenerated(),
		options.toGenerated(),
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

// GetCertificateContacts returns the set of certificate contact resources in the specified key vault. This operation
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

// DeleteContacts deletes the certificate contacts for a specified key vault certificate. This operation requires the certificates/managecontacts permission.
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

// UpdateCertificatePolicy sets specified members in the certificate policy, leave others as null. This operation requires the certificates/update permission.
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

// GetCertificatePolicy returns the specified certificate policy resources in the specified key vault. This operation requires the certificates/get permission.
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
	CertificateAttributes *CertificateProperties `json:"attributes,omitempty"`

	// The management policy for the certificate.
	CertificatePolicy *CertificatePolicy `json:"policy,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]string `json:"tags,omitempty"`
}

func (u *UpdateCertificatePropertiesOptions) toGenerated() *generated.KeyVaultClientUpdateCertificateOptions {
	return &generated.KeyVaultClientUpdateCertificateOptions{}
}

// UpdateCertificatePropertiesResponse contains the result from method Client.UpdateCertificateProperties.
type UpdateCertificatePropertiesResponse struct {
	KeyVaultCertificate
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// UpdateCertificate applies the specified update on the given certificate; the only elements updated are the certificate's
// attributes. This operation requires the certificates/update permission.
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
			Tags:                  convertToGeneratedMap(options.Tags),
		},
		options.toGenerated(),
	)
	if err != nil {
		return UpdateCertificatePropertiesResponse{}, err
	}
	return UpdateCertificatePropertiesResponse{
		RawResponse:         resp.RawResponse,
		KeyVaultCertificate: certificateFromGenerated(&resp.CertificateBundle),
	}, nil
}

// MergeCertificateOptions contains the optional parameters for the Client.MergeCertificate function.
type MergeCertificateOptions struct {
	// The attributes of the certificate (optional).
	CertificateAttributes *CertificateProperties `json:"attributes,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]string `json:"tags,omitempty"`
}

func (m *MergeCertificateOptions) toGenerated() *generated.KeyVaultClientMergeCertificateOptions {
	return &generated.KeyVaultClientMergeCertificateOptions{}
}

// MergeCertificateResponse contains the response from method Client.MergeCertificate.
type MergeCertificateResponse struct {
	KeyVaultCertificateWithPolicy

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
			Tags:                  convertToGeneratedMap(options.Tags),
		},
		options.toGenerated(),
	)
	if err != nil {
		return MergeCertificateResponse{}, err
	}

	return MergeCertificateResponse{
		RawResponse: resp.RawResponse,
		KeyVaultCertificateWithPolicy: KeyVaultCertificateWithPolicy{
			Properties:     certificateAttributesFromGenerated(resp.Attributes),
			Cer:            resp.Cer,
			ContentType:    resp.ContentType,
			Tags:           convertGeneratedMap(resp.Tags),
			ID:             resp.ID,
			KeyID:          resp.Kid,
			SecretID:       resp.Sid,
			X509Thumbprint: resp.X509Thumbprint,
			Policy:         certificatePolicyFromGenerated(resp.Policy),
		},
	}, nil
}

// RestoreCertificateBackupOptions contains the optional parameters for the Client.RestoreCertificateBackup method
type RestoreCertificateBackupOptions struct{}

func (r *RestoreCertificateBackupOptions) toGenerated() *generated.KeyVaultClientRestoreCertificateOptions {
	return &generated.KeyVaultClientRestoreCertificateOptions{}
}

// RestoreCertificateBackupResponse contains the response from method Client.RestoreCertificateBackup
type RestoreCertificateBackupResponse struct {
	KeyVaultCertificateWithPolicy

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
		RawResponse: resp.RawResponse,
		KeyVaultCertificateWithPolicy: KeyVaultCertificateWithPolicy{
			Properties:     certificateAttributesFromGenerated(resp.Attributes),
			Cer:            resp.Cer,
			ContentType:    resp.ContentType,
			Tags:           convertGeneratedMap(resp.Tags),
			ID:             resp.ID,
			KeyID:          resp.Kid,
			SecretID:       resp.Sid,
			X509Thumbprint: resp.X509Thumbprint,
			Policy:         certificatePolicyFromGenerated(resp.Policy),
		},
	}, nil
}

// BeginRecoverDeletedCertificateOptions contains the optional parameters for the Client.BeginRecoverDeletedCertificate function
type BeginRecoverDeletedCertificateOptions struct{}

func (b *BeginRecoverDeletedCertificateOptions) toGenerated() *generated.KeyVaultClientRecoverDeletedCertificateOptions {
	return &generated.KeyVaultClientRecoverDeletedCertificateOptions{}
}

// RecoverDeletedCertificatePoller is the poller for the Client.RecoverDeletedCertificate
type RecoverDeletedCertificatePoller struct {
	certName        string
	vaultUrl        string
	client          *generated.KeyVaultClient
	recoverResponse generated.KeyVaultClientRecoverDeletedCertificateResponse
	lastResponse    generated.KeyVaultClientGetCertificateResponse
	RawResponse     *http.Response
}

// Done returns true when the polling operation is completed
func (b *RecoverDeletedCertificatePoller) Done() bool {
	return b.RawResponse.StatusCode == http.StatusOK
}

// Poll fetches the latest state of the LRO. It returns an HTTP response or error.
// If the LRO has completed successfully, the poller's state is updated and the HTTP response is returned.
// If the LRO has completed with failure or was cancelled, the poller's state is updated and the error is returned.
func (b *RecoverDeletedCertificatePoller) Poll(ctx context.Context) (*http.Response, error) {
	resp, err := b.client.GetCertificate(ctx, b.vaultUrl, b.certName, "", nil)
	b.lastResponse = resp
	var httpErr *azcore.ResponseError
	if errors.As(err, &httpErr) {
		return httpErr.RawResponse, err
	}
	return resp.RawResponse, nil
}

// FinalResponse returns the final response after the operations has finished
func (b *RecoverDeletedCertificatePoller) FinalResponse(ctx context.Context) (RecoverDeletedCertificateResponse, error) {
	return recoverDeletedCertificateResponseFromGenerated(b.recoverResponse), nil
}

// pollUntilDone is the method for the Response.PollUntilDone struct
func (b *RecoverDeletedCertificatePoller) pollUntilDone(ctx context.Context, t time.Duration) (RecoverDeletedCertificateResponse, error) {
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
	KeyVaultCertificate

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// change recover deleted certificate reponse to the generated version.
func recoverDeletedCertificateResponseFromGenerated(i generated.KeyVaultClientRecoverDeletedCertificateResponse) RecoverDeletedCertificateResponse {
	return RecoverDeletedCertificateResponse{
		KeyVaultCertificate: certificateFromGenerated(&i.CertificateBundle),
		RawResponse:         i.RawResponse,
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
	var httpErr *azcore.ResponseError
	if errors.As(err, &httpErr) {
		if httpErr.RawResponse.StatusCode != http.StatusNotFound {
			return RecoverDeletedCertificatePollerResponse{}, err
		}
	}

	p := RecoverDeletedCertificatePoller{
		lastResponse:    getResp,
		certName:        certName,
		client:          c.genClient,
		vaultUrl:        c.vaultURL,
		recoverResponse: resp,
		RawResponse:     getResp.RawResponse,
	}

	return RecoverDeletedCertificatePollerResponse{
		PollUntilDone: p.pollUntilDone,
		Poller:        p,
		RawResponse:   getResp.RawResponse,
	}, nil
}

// ListDeletedCertificatesPager is the pager returned by Client.ListDeletedCertificates
type ListDeletedCertificatesPager struct {
	genPager *generated.KeyVaultClientGetDeletedCertificatesPager
}

// PageResponse returns the current page of results
func (l *ListDeletedCertificatesPager) PageResponse() ListDeletedCertificatesPage {
	resp := l.genPager.PageResponse()

	var vals []*DeletedCertificateItem

	for _, v := range resp.Value {
		vals = append(vals, &DeletedCertificateItem{
			RecoveryID:         v.RecoveryID,
			DeletedDate:        v.DeletedDate,
			ScheduledPurgeDate: v.ScheduledPurgeDate,
			Properties:         certificateAttributesFromGenerated(v.Attributes),
			ID:                 v.ID,
			Tags:               convertGeneratedMap(v.Tags),
			X509Thumbprint:     v.X509Thumbprint,
		})
	}

	return ListDeletedCertificatesPage{
		RawResponse:  resp.RawResponse,
		Certificates: vals,
	}
}

// Err returns an error if the last operation resulted in an error.
func (l *ListDeletedCertificatesPager) Err() error {
	return l.genPager.Err()
}

// NextPage fetches the next page of results.
func (l *ListDeletedCertificatesPager) NextPage(ctx context.Context) bool {
	return l.genPager.NextPage(ctx)
}

// ListDeletedCertificatesPage holds the data for a single page.
type ListDeletedCertificatesPage struct {
	// READ-ONLY; A response message containing a list of deleted certificates in the vault along with a link to the next page of deleted certificates
	Certificates []*DeletedCertificateItem `json:"value,omitempty" azure:"ro"`

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// ListDeletedCertificatesOptions contains the optional parameters for the Client.ListDeletedCertificates operation.
type ListDeletedCertificatesOptions struct {
}

// Convert publicly exposed options to the generated version.a
func (l *ListDeletedCertificatesOptions) toGenerated() *generated.KeyVaultClientGetDeletedCertificatesOptions {
	return &generated.KeyVaultClientGetDeletedCertificatesOptions{}
}

// ListDeletedCertificates retrieves the certificates in the current vault which are in a deleted state and ready for recovery or purging.
// This operation includes deletion-specific information. This operation requires the certificates/get/list permission. This operation can
// only be enabled on soft-delete enabled vaults.
func (c *Client) ListDeletedCertificates(options *ListDeletedCertificatesOptions) ListDeletedCertificatesPager {
	if options == nil {
		options = &ListDeletedCertificatesOptions{}
	}

	return ListDeletedCertificatesPager{
		genPager: c.genClient.GetDeletedCertificates(c.vaultURL, options.toGenerated()),
	}
}

// CancelCertificateOperationOptions contains the optional parameters for the Client.CancelCertificateOperation function
type CancelCertificateOperationOptions struct{}

func (c *CancelCertificateOperationOptions) toGenerated() *generated.KeyVaultClientUpdateCertificateOperationOptions {
	return &generated.KeyVaultClientUpdateCertificateOperationOptions{}
}

// CancelCertificateOperationResponse contains the response models for the Client.CancelCertificateOperation function
type CancelCertificateOperationResponse struct {
	CertificateOperation

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// CancelCertificateOperation cancels a certificate creation operation that is already in progress. This operation requires the certificates/update permission.
func (c *Client) CancelCertificateOperation(ctx context.Context, certName string, options *CancelCertificateOperationOptions) (CancelCertificateOperationResponse, error) {
	resp, err := c.genClient.UpdateCertificateOperation(
		ctx,
		c.vaultURL,
		certName,
		generated.CertificateOperationUpdateParameter{
			CancellationRequested: to.BoolPtr(true),
		},
		options.toGenerated(),
	)
	if err != nil {
		return CancelCertificateOperationResponse{}, err
	}

	return CancelCertificateOperationResponse{
		RawResponse:          resp.RawResponse,
		CertificateOperation: certificateOperationFromGenerated(resp.CertificateOperation),
	}, nil
}

// DeleteCertificateOperationsOptions contains the optional parameters for the Client.DeleteCertificateOperation function.
type DeleteCertificateOperationOptions struct{}

func (d *DeleteCertificateOperationOptions) toGenerated() *generated.KeyVaultClientDeleteCertificateOperationOptions {
	return &generated.KeyVaultClientDeleteCertificateOperationOptions{}
}

// DeleteCertificateOperationResponse contains the response for the Client.DeleteCertificateOperation function.
type DeleteCertificateOperationResponse struct {
	CertificateOperation
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// DeleteCertificateOperation deletes the creation operation for a specified certificate that is in the process of being created. The certificate is no
// longer created. This operation requires the certificates/update permission.
func (c *Client) DeleteCertificateOperation(ctx context.Context, certName string, options *DeleteCertificateOperationOptions) (DeleteCertificateOperationResponse, error) {
	resp, err := c.genClient.DeleteCertificateOperation(
		ctx,
		c.vaultURL,
		certName,
		options.toGenerated(),
	)

	if err != nil {
		return DeleteCertificateOperationResponse{}, err
	}

	return DeleteCertificateOperationResponse{
		RawResponse:          resp.RawResponse,
		CertificateOperation: certificateOperationFromGenerated(resp.CertificateOperation),
	}, nil
}
