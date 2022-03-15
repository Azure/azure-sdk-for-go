//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcertificates

import (
	"context"
	"errors"
	"net/http"
	"reflect"
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

// ClientOptions are optional parameters for NewClient
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

// BeginCreateCertificateOptions contains optional parameters for Client.BeginCreateCertificate
type BeginCreateCertificateOptions struct {
	// The attributes of the certificate (optional).
	CertificateAttributes *CertificateProperties `json:"attributes,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]string `json:"tags,omitempty"`
}

func (b BeginCreateCertificateOptions) toGenerated() *generated.KeyVaultClientCreateCertificateOptions {
	return &generated.KeyVaultClientCreateCertificateOptions{}
}

// CreateCertificateResponse contains response fields for Client.BeginCreateCertificate
type CreateCertificateResponse struct {
	CertificateOperation
}

// CreateCertificatePoller is the poller returned by the Client.BeginCreateCertificate
type CreateCertificatePoller struct {
	certName       string
	certVersion    string
	vaultURL       string
	client         *generated.KeyVaultClient
	createResponse CreateCertificateResponse
	lastResponse   generated.KeyVaultClientGetCertificateResponse
	rawResponse    *http.Response
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
	return nil, err
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
		b.rawResponse = resp
		if b.Done() {
			break
		}
		time.Sleep(t)
	}
	return b.createResponse, nil
}

// CreateCertificatePollerResponse contains response fields for Client.BeginCreateCertificate
type CreateCertificatePollerResponse struct {
	// PollUntilDone will poll the service endpoint until a terminal state is reached or an error occurs
	PollUntilDone func(context.Context, time.Duration) (CreateCertificateResponse, error)

	// Poller contains an initialized WidgetPoller
	Poller CreateCertificatePoller
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
		PollUntilDone: p.pollUntilDone,
	}, nil
}

// GetCertificateOptions contains optional parameters for Client.GetCertificate
type GetCertificateOptions struct {
	Version string
}

// GetCertificateResponse contains response fields for Client.GetCertificate
type GetCertificateResponse struct {
	KeyVaultCertificateWithPolicy
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

// GetCertificateOperationOptions contains optional parameters for Client.GetCertificateOperation
type GetCertificateOperationOptions struct {
	// placeholder for future optional parameters.
}

func (g *GetCertificateOperationOptions) toGenerated() *generated.KeyVaultClientGetCertificateOperationOptions {
	return &generated.KeyVaultClientGetCertificateOperationOptions{}
}

// GetCertificateOperationResponse contains response field for Client.GetCertificateOperation
type GetCertificateOperationResponse struct {
	CertificateOperation
}

// GetCertificateOperation gets the creation operation associated with a specified certificate. This operation requires the certificates/get permission.
func (c *Client) GetCertificateOperation(ctx context.Context, certName string, options *GetCertificateOperationOptions) (GetCertificateOperationResponse, error) {
	resp, err := c.genClient.GetCertificateOperation(ctx, c.vaultURL, certName, options.toGenerated())
	if err != nil {
		return GetCertificateOperationResponse{}, err
	}

	return GetCertificateOperationResponse{
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

// BeginDeleteCertificateOptions contains optional parameters for Client.BeginDeleteCertificate
type BeginDeleteCertificateOptions struct {
	// placeholder for future optional parameters.
}

// convert public options to generated options struct
func (b *BeginDeleteCertificateOptions) toGenerated() *generated.KeyVaultClientDeleteCertificateOptions {
	return &generated.KeyVaultClientDeleteCertificateOptions{}
}

// DeleteCertificateResponse contains response fields for BeginDeleteCertificatePoller.FinalResponse
type DeleteCertificateResponse struct {
	DeletedCertificate
}

func deleteCertificateResponseFromGenerated(g *generated.KeyVaultClientDeleteCertificateResponse) DeleteCertificateResponse {
	if g == nil {
		return DeleteCertificateResponse{}
	}
	return DeleteCertificateResponse{
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

// DeleteCertificatePoller is the poller returned by the Client.BeginDeleteCertificate operation
type DeleteCertificatePoller struct {
	certificateName string // This is the certificate to Poll for in GetDeletedCertificate
	vaultURL        string
	client          *generated.KeyVaultClient
	deleteResponse  generated.KeyVaultClientDeleteCertificateResponse
	lastResponse    generated.KeyVaultClientGetDeletedCertificateResponse
	rawResponse     *http.Response
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
		s.rawResponse = resp
		if s.Done() {
			break
		}
		time.Sleep(t)
	}
	return deleteCertificateResponseFromGenerated(&s.deleteResponse), nil
}

// DeleteCertificatePollerResponse contains response fields for Client.BeginDeleteCertificate
type DeleteCertificatePollerResponse struct {
	// PollUntilDone will poll the service endpoint until a terminal state is reached or an error occurs
	PollUntilDone func(context.Context, time.Duration) (DeleteCertificateResponse, error)

	// Poller contains an initialized WidgetPoller
	Poller DeleteCertificatePoller
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
		PollUntilDone: s.pollUntilDone,
	}, nil
}

// PurgeDeletedCertificateOptions contains optional parameters for Client.PurgeDeletedCertificateOptions
type PurgeDeletedCertificateOptions struct {
	// placeholder for future optional parameters.
}

func (p *PurgeDeletedCertificateOptions) toGenerated() *generated.KeyVaultClientPurgeDeletedCertificateOptions {
	return &generated.KeyVaultClientPurgeDeletedCertificateOptions{}
}

// PurgeDeletedCertificateResponse contains response fields for Client.PurgeDeletedCertificate
type PurgeDeletedCertificateResponse struct {
	// placeholder for future reponse fields
}

// PurgeDeletedCertificate operation performs an irreversible deletion of the specified certificate, without possibility for recovery. The operation
// is not available if the recovery level does not specify 'Purgeable'. This operation requires the certificate/purge permission.
func (c *Client) PurgeDeletedCertificate(ctx context.Context, certName string, options *PurgeDeletedCertificateOptions) (PurgeDeletedCertificateResponse, error) {
	_, err := c.genClient.PurgeDeletedCertificate(ctx, c.vaultURL, certName, options.toGenerated())
	if err != nil {
		return PurgeDeletedCertificateResponse{}, err
	}

	return PurgeDeletedCertificateResponse{}, nil
}

// GetDeletedCertificateOptions contains optional parameters for Client.GetDeletedCertificate
type GetDeletedCertificateOptions struct {
	// placeholder for future optional parameters.
}

func (g *GetDeletedCertificateOptions) toGenerated() *generated.KeyVaultClientGetDeletedCertificateOptions {
	return &generated.KeyVaultClientGetDeletedCertificateOptions{}
}

// GetDeletedCertificateResponse contains response field for Client.GetDeletedCertificate
type GetDeletedCertificateResponse struct {
	DeletedCertificate
}

// GetDeletedCertificate retrieves the deleted certificate information plus its attributes, such as retention interval, scheduled permanent deletion
// and the current deletion recovery level. This operation requires the certificates/get permission.
func (c *Client) GetDeletedCertificate(ctx context.Context, certName string, options *GetDeletedCertificateOptions) (GetDeletedCertificateResponse, error) {
	resp, err := c.genClient.GetDeletedCertificate(ctx, c.vaultURL, certName, options.toGenerated())
	if err != nil {
		return GetDeletedCertificateResponse{}, err
	}

	return GetDeletedCertificateResponse{
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

// BackupCertificateOptions contains optional parameters for Client.BackupCertificateOptions
type BackupCertificateOptions struct {
	// placeholder for future optional parameters.
}

func (b *BackupCertificateOptions) toGenerated() *generated.KeyVaultClientBackupCertificateOptions {
	return &generated.KeyVaultClientBackupCertificateOptions{}
}

// BackupCertificateResponse contains response field for Client.BackupCertificate
type BackupCertificateResponse struct {
	// READ-ONLY; The backup blob containing the backed up certificate.
	Value []byte `json:"value,omitempty" azure:"ro"`
}

// BackupCertificate requests that a backup of the specified certificate be downloaded to the client. All versions of the certificate will be downloaded.
// This operation requires the certificates/backup permission.
func (c *Client) BackupCertificate(ctx context.Context, certName string, options *BackupCertificateOptions) (BackupCertificateResponse, error) {
	resp, err := c.genClient.BackupCertificate(ctx, c.vaultURL, certName, options.toGenerated())
	if err != nil {
		return BackupCertificateResponse{}, err
	}

	return BackupCertificateResponse{
		Value: resp.Value,
	}, nil
}

// ImportCertificateOptions contains optional parameters for Client.ImportCertificate
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

// ImportCertificateResponse contains response fields for Client.ImportCertificate
type ImportCertificateResponse struct {
	KeyVaultCertificateWithPolicy
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
	vaultURL  string
	genClient *generated.KeyVaultClient
	nextLink  *string
}

// More returns true if there are more pages to return
func (l *ListCertificatesPager) More() bool {
	if !reflect.ValueOf(l.nextLink).IsZero() {
		if l.nextLink == nil || len(*l.nextLink) == 0 {
			return false
		}
	}
	return true
}

// NextPage returns the current page of results
func (l *ListCertificatesPager) NextPage(ctx context.Context) (ListCertificatesPageResponse, error) {
	var resp *http.Response
	var err error
	if l.nextLink == nil {
		req, err := l.genClient.GetCertificatesCreateRequest(
			ctx,
			l.vaultURL,
			&generated.KeyVaultClientGetCertificatesOptions{},
		)
		if err != nil {
			return ListCertificatesPageResponse{}, err
		}
		resp, err = l.genClient.Pl.Do(req)
		if err != nil {
			return ListCertificatesPageResponse{}, err
		}
	} else {
		req, err := runtime.NewRequest(ctx, http.MethodGet, *l.nextLink)
		if err != nil {
			return ListCertificatesPageResponse{}, err
		}
		resp, err = l.genClient.Pl.Do(req)
		if err != nil {
			return ListCertificatesPageResponse{}, err
		}
	}
	if err != nil {
		return ListCertificatesPageResponse{}, err
	}
	result, err := l.genClient.GetCertificatesHandleResponse(resp)
	if err != nil {
		return ListCertificatesPageResponse{}, err
	}
	if result.NextLink == nil {
		// Set it to the zero value
		result.NextLink = to.StringPtr("")
	}
	l.nextLink = result.NextLink
	return listKeysPageFromGenerated(result), nil
}

// ListCertificatesOptions contains optional parameters for Client.ListCertificates
type ListCertificatesOptions struct {
	// placeholder for future optional parameters.
}

// ListCertificatesPageResponse contains response fields for ListCertificatesPager.NextPage
type ListCertificatesPageResponse struct {
	// READ-ONLY; A response message containing a list of certificates in the key vault along with a link to the next page of certificates.
	Certificates []*CertificateItem `json:"value,omitempty" azure:"ro"`
}

// convert internal Response to ListCertificatesPage
func listKeysPageFromGenerated(i generated.KeyVaultClientGetCertificatesResponse) ListCertificatesPageResponse {
	var vals []*CertificateItem

	for _, v := range i.Value {
		vals = append(vals, &CertificateItem{
			Properties:     certificateAttributesFromGenerated(v.Attributes),
			ID:             v.ID,
			Tags:           convertGeneratedMap(v.Tags),
			X509Thumbprint: v.X509Thumbprint,
		})
	}

	return ListCertificatesPageResponse{
		Certificates: vals,
	}
}

// ListCertificates retrieves a list of the certificates in the Key Vault as JSON Web Key structures that contain the
// public part of a stored certificate. The LIST operation is applicable to all certificate types, however only the
// base certificate identifier, attributes, and tags are provided in the response. Individual versions of a
// certificate are not listed in the response. This operation requires the certificates/list permission.
func (c *Client) ListCertificates(options *ListCertificatesOptions) ListCertificatesPager {
	return ListCertificatesPager{
		vaultURL:  c.vaultURL,
		genClient: c.genClient,
		nextLink:  nil,
	}
}

// ListCertificateVersionsPager is the pager returned by Client.ListCertificateVersions
type ListCertificateVersionsPager struct {
	vaultURL  string
	genClient *generated.KeyVaultClient
	nextLink  *string
	certName  string
}

// More returns true if there are more pages to return
func (l *ListCertificateVersionsPager) More() bool {
	if !reflect.ValueOf(l.nextLink).IsZero() {
		if l.nextLink == nil || len(*l.nextLink) == 0 {
			return false
		}
	}
	return true
}

// NextPage returns the current page of results
func (l *ListCertificateVersionsPager) NextPage(ctx context.Context) (ListCertificateVersionsPageResponse, error) {
	var resp *http.Response
	var err error
	if l.nextLink == nil {
		req, err := l.genClient.GetCertificateVersionsCreateRequest(
			ctx,
			l.vaultURL,
			l.certName,
			&generated.KeyVaultClientGetCertificateVersionsOptions{},
		)
		if err != nil {
			return ListCertificateVersionsPageResponse{}, err
		}
		resp, err = l.genClient.Pl.Do(req)
		if err != nil {
			return ListCertificateVersionsPageResponse{}, err
		}
	} else {
		req, err := runtime.NewRequest(ctx, http.MethodGet, *l.nextLink)
		if err != nil {
			return ListCertificateVersionsPageResponse{}, err
		}
		resp, err = l.genClient.Pl.Do(req)
		if err != nil {
			return ListCertificateVersionsPageResponse{}, err
		}
	}
	if err != nil {
		return ListCertificateVersionsPageResponse{}, err
	}
	result, err := l.genClient.GetCertificateVersionsHandleResponse(resp)
	if err != nil {
		return ListCertificateVersionsPageResponse{}, err
	}
	if result.NextLink == nil {
		// Set it to the zero value
		result.NextLink = to.StringPtr("")
	}
	l.nextLink = result.NextLink
	return listCertificateVersionsPageFromGenerated(result), nil
}

// ListCertificateVersionsOptions contains optional parameters for Client.ListCertificateVersions
type ListCertificateVersionsOptions struct {
	// placeholder for future optional parameters.
}

// ListCertificateVersionsPageResponse contains response fields for ListCertificateVersionsPager.NextPage
type ListCertificateVersionsPageResponse struct {
	// READ-ONLY; A response message containing a list of certificates in the key vault along with a link to the next page of certificates.
	Certificates []*CertificateItem `json:"value,omitempty" azure:"ro"`
}

// create ListCertificatesPage from generated pager
func listCertificateVersionsPageFromGenerated(i generated.KeyVaultClientGetCertificateVersionsResponse) ListCertificateVersionsPageResponse {
	var vals []*CertificateItem
	for _, v := range i.Value {
		vals = append(vals, &CertificateItem{
			Properties:     certificateAttributesFromGenerated(v.Attributes),
			ID:             v.ID,
			Tags:           convertGeneratedMap(v.Tags),
			X509Thumbprint: v.X509Thumbprint,
		})
	}

	return ListCertificateVersionsPageResponse{
		Certificates: vals,
	}
}

// ListCertificateVersions lists all versions of the specified certificate. The full certificate identifer and
// attributes are provided in the response. No values are returned for the certificates. This operation
// requires the certificates/list permission.
func (c *Client) ListCertificateVersions(certificateName string, options *ListCertificateVersionsOptions) ListCertificateVersionsPager {
	return ListCertificateVersionsPager{
		certName:  certificateName,
		vaultURL:  c.vaultURL,
		genClient: c.genClient,
		nextLink:  nil,
	}
}

// CreateIssuerOptions contains optional parameters for Client.CreateIssuer
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

// CreateIssuerResponse contains response fields for Client.CreateIssuer
type CreateIssuerResponse struct {
	CertificateIssuer
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

	cr := CreateIssuerResponse{}
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

// GetIssuerOptions contains optional parameters for Client.GetIssuer
type GetIssuerOptions struct {
	// placeholder for future optional parameters.
}

func (g *GetIssuerOptions) toGenerated() *generated.KeyVaultClientGetCertificateIssuerOptions {
	return &generated.KeyVaultClientGetCertificateIssuerOptions{}
}

// GetIssuerResponse contains response fields for ClientGetIssuer
type GetIssuerResponse struct {
	CertificateIssuer
}

// GetIssuer returns the specified certificate issuer resources in the specified key vault. This operation
// requires the certificates/manageissuers/getissuers permission.
func (c *Client) GetIssuer(ctx context.Context, issuerName string, options *GetIssuerOptions) (GetIssuerResponse, error) {
	resp, err := c.genClient.GetCertificateIssuer(ctx, c.vaultURL, issuerName, options.toGenerated())
	if err != nil {
		return GetIssuerResponse{}, err
	}

	g := GetIssuerResponse{}
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
	vaultURL  string
	genClient *generated.KeyVaultClient
	nextLink  *string
}

// More returns true if there are more pages to return
func (l *ListPropertiesOfIssuersPager) More() bool {
	if !reflect.ValueOf(l.nextLink).IsZero() {
		if l.nextLink == nil || len(*l.nextLink) == 0 {
			return false
		}
	}
	return true
}

// NextPage returns the current page of results
func (l *ListPropertiesOfIssuersPager) NextPage(ctx context.Context) (ListIssuersPropertiesOfIssuersPageResponse, error) {
	var resp *http.Response
	var err error
	if l.nextLink == nil {
		req, err := l.genClient.GetCertificateIssuersCreateRequest(
			ctx,
			l.vaultURL,
			&generated.KeyVaultClientGetCertificateIssuersOptions{},
		)
		if err != nil {
			return ListIssuersPropertiesOfIssuersPageResponse{}, err
		}
		resp, err = l.genClient.Pl.Do(req)
		if err != nil {
			return ListIssuersPropertiesOfIssuersPageResponse{}, err
		}
	} else {
		req, err := runtime.NewRequest(ctx, http.MethodGet, *l.nextLink)
		if err != nil {
			return ListIssuersPropertiesOfIssuersPageResponse{}, err
		}
		resp, err = l.genClient.Pl.Do(req)
		if err != nil {
			return ListIssuersPropertiesOfIssuersPageResponse{}, err
		}
	}
	if err != nil {
		return ListIssuersPropertiesOfIssuersPageResponse{}, err
	}
	result, err := l.genClient.GetCertificateIssuersHandleResponse(resp)
	if err != nil {
		return ListIssuersPropertiesOfIssuersPageResponse{}, err
	}
	if result.NextLink == nil {
		// Set it to the zero value
		result.NextLink = to.StringPtr("")
	}
	l.nextLink = result.NextLink
	return listIssuersPageFromGenerated(result), nil
}

// ListPropertiesOfIssuersOptions contains optional parameters for Client.ListIssuers
type ListPropertiesOfIssuersOptions struct {
	// placeholder for future optional parameters
}

// ListIssuersPropertiesOfIssuersPageResponse contains response fields for ListPropertiesOfIssuersPager.NextPage
type ListIssuersPropertiesOfIssuersPageResponse struct {
	// READ-ONLY; A response message containing a list of certificates in the key vault along with a link to the next page of certificates.
	Issuers []*CertificateIssuerItem `json:"value,omitempty" azure:"ro"`
}

// convert internal Response to ListPropertiesOfIssuersPage
func listIssuersPageFromGenerated(i generated.KeyVaultClientGetCertificateIssuersResponse) ListIssuersPropertiesOfIssuersPageResponse {
	var vals []*CertificateIssuerItem

	for _, v := range i.Value {
		vals = append(vals, certificateIssuerItemFromGenerated(v))
	}

	return ListIssuersPropertiesOfIssuersPageResponse{Issuers: vals}
}

// ListPropertiesOfIssuers returns a pager that can be used to get the set of certificate issuer resources in the specified key vault. This operation
// requires the certificates/manageissuers/getissuers permission.
func (c *Client) ListPropertiesOfIssuers(options *ListPropertiesOfIssuersOptions) ListPropertiesOfIssuersPager {
	return ListPropertiesOfIssuersPager{
		genClient: c.genClient,
		vaultURL:  c.vaultURL,
		nextLink:  nil,
	}
}

// DeleteIssuerOptions contains optional parameters for Client.DeleteIssuer
type DeleteIssuerOptions struct {
	// placeholder for future optional parameters.
}

func (d *DeleteIssuerOptions) toGenerated() *generated.KeyVaultClientDeleteCertificateIssuerOptions {
	return &generated.KeyVaultClientDeleteCertificateIssuerOptions{}
}

// DeleteIssuerResponse contains response fields for Client.DeleteIssuer
type DeleteIssuerResponse struct {
	CertificateIssuer
}

// DeleteIssuer permanently removes the specified certificate issuer from the vault. This operation requires the certificates/manageissuers/deleteissuers permission.
func (c *Client) DeleteIssuer(ctx context.Context, issuerName string, options *DeleteIssuerOptions) (DeleteIssuerResponse, error) {
	resp, err := c.genClient.DeleteCertificateIssuer(ctx, c.vaultURL, issuerName, options.toGenerated())
	if err != nil {
		return DeleteIssuerResponse{}, err
	}

	d := DeleteIssuerResponse{}
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

// UpdateIssuerOptions contains optional parameters for Client.UpdateIssuer
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

// UpdateIssuerResponse contains response fields for Client.UpdateIssuer
type UpdateIssuerResponse struct {
	CertificateIssuer
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

	u := UpdateIssuerResponse{}
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

// SetContactsOptions contains optional parameters for Client.CreateContacts
type SetContactsOptions struct {
	// placeholder for future optional parameters.
}

func (s *SetContactsOptions) toGenerated() *generated.KeyVaultClientSetCertificateContactsOptions {
	return &generated.KeyVaultClientSetCertificateContactsOptions{}
}

// SetContactsResponse contains response fields for Client.CreateContacts
type SetContactsResponse struct {
	Contacts
}

// SetContacts sets the certificate contacts for the specified key vault. This operation requires the certificates/managecontacts permission.
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
		Contacts: Contacts{
			ID:          resp.ID,
			ContactList: contactListFromGenerated(resp.ContactList),
		},
	}, nil
}

// GetContactsOptions contains optional parameters for Client.GetContacts
type GetContactsOptions struct {
	// placeholder for future optional parameters.
}

func (g *GetContactsOptions) toGenerated() *generated.KeyVaultClientGetCertificateContactsOptions {
	return &generated.KeyVaultClientGetCertificateContactsOptions{}
}

// GetContactsResponse contains response fields for Client.GetContacts
type GetContactsResponse struct {
	Contacts
}

// GetContacts returns the set of certificate contact resources in the specified key vault. This operation
// requires the certificates/managecontacts permission.
func (c *Client) GetContacts(ctx context.Context, options *GetContactsOptions) (GetContactsResponse, error) {
	resp, err := c.genClient.GetCertificateContacts(ctx, c.vaultURL, options.toGenerated())
	if err != nil {
		return GetContactsResponse{}, err
	}

	return GetContactsResponse{
		Contacts: Contacts{
			ID:          resp.ID,
			ContactList: contactListFromGenerated(resp.ContactList),
		},
	}, nil
}

// DeleteContactsOptions contains optional parameters for Client.DeleteContacts
type DeleteContactsOptions struct {
	// placeholder for future optional parameters.
}

func (d *DeleteContactsOptions) toGenerated() *generated.KeyVaultClientDeleteCertificateContactsOptions {
	return &generated.KeyVaultClientDeleteCertificateContactsOptions{}
}

// DeleteContactsResponse contains response field for Client.DeleteContacts
type DeleteContactsResponse struct {
	Contacts
}

// DeleteContacts deletes the certificate contacts for a specified key vault certificate. This operation requires the certificates/managecontacts permission.
func (c *Client) DeleteContacts(ctx context.Context, options *DeleteContactsOptions) (DeleteContactsResponse, error) {
	resp, err := c.genClient.DeleteCertificateContacts(ctx, c.vaultURL, options.toGenerated())
	if err != nil {
		return DeleteContactsResponse{}, err
	}

	return DeleteContactsResponse{
		Contacts: Contacts{
			ContactList: contactListFromGenerated(resp.ContactList),
			ID:          resp.ID,
		},
	}, nil
}

// UpdateCertificatePolicyOptions contains optional parameters for Client.UpdateCertificatePolicy
type UpdateCertificatePolicyOptions struct {
	// placeholder for future optional parameters.
}

func (u *UpdateCertificatePolicyOptions) toGenerated() *generated.KeyVaultClientUpdateCertificatePolicyOptions {
	return &generated.KeyVaultClientUpdateCertificatePolicyOptions{}
}

// UpdateCertificatePolicyResponse contains response fields for Client.UpdateCertificatePolicy
type UpdateCertificatePolicyResponse struct {
	CertificatePolicy
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
		CertificatePolicy: *certificatePolicyFromGenerated(&resp.CertificatePolicy),
	}, nil
}

// GetCertificatePolicyOptions contains optional parameters for Client.GetCertificatePolicy
type GetCertificatePolicyOptions struct {
	// placeholder for future optional parameters.
}

func (g *GetCertificatePolicyOptions) toGenerated() *generated.KeyVaultClientGetCertificatePolicyOptions {
	return &generated.KeyVaultClientGetCertificatePolicyOptions{}
}

// GetCertificatePolicyResponse contains response fields for Client.GetCertificatePolicy
type GetCertificatePolicyResponse struct {
	CertificatePolicy
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
		CertificatePolicy: *certificatePolicyFromGenerated(&resp.CertificatePolicy),
	}, nil
}

// UpdateCertificatePropertiesOptions contains optional parameters for Client.UpdateCertificateProperties
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

// UpdateCertificatePropertiesResponse contains response fields for Client.UpdateCertificateProperties
type UpdateCertificatePropertiesResponse struct {
	KeyVaultCertificate
}

// UpdateCertificateProperties applies the specified update on the given certificate; the only elements updated are the certificate's
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
		KeyVaultCertificate: certificateFromGenerated(&resp.CertificateBundle),
	}, nil
}

// MergeCertificateOptions contains optional parameters for Client.MergeCertificate
type MergeCertificateOptions struct {
	// The attributes of the certificate (optional).
	CertificateAttributes *CertificateProperties `json:"attributes,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]string `json:"tags,omitempty"`
}

func (m *MergeCertificateOptions) toGenerated() *generated.KeyVaultClientMergeCertificateOptions {
	return &generated.KeyVaultClientMergeCertificateOptions{}
}

// MergeCertificateResponse contains response fields for Client.MergeCertificate
type MergeCertificateResponse struct {
	KeyVaultCertificateWithPolicy
}

// MergeCertificate operation performs the merging of a certificate or certificate chain with a key pair currently available in the service. This operation requires the certificates/create permission.
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

// RestoreCertificateBackupOptions contains optional parameters for Client.RestoreCertificateBackup
type RestoreCertificateBackupOptions struct {
	// placeholder for future optional parameters.
}

func (r *RestoreCertificateBackupOptions) toGenerated() *generated.KeyVaultClientRestoreCertificateOptions {
	return &generated.KeyVaultClientRestoreCertificateOptions{}
}

// RestoreCertificateBackupResponse contains response fields for Client.RestoreCertificateBackup
type RestoreCertificateBackupResponse struct {
	KeyVaultCertificateWithPolicy
}

// RestoreCertificateBackup performs the reversal of the Delete operation. The operation is applicable in vaults
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

// BeginRecoverDeletedCertificateOptions contains optional parameters for Client.BeginRecoverDeletedCertificate
type BeginRecoverDeletedCertificateOptions struct {
	// placeholder for future optional parameters.
}

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
	rawResponse     *http.Response
}

// Done returns true when the polling operation is completed
func (b *RecoverDeletedCertificatePoller) Done() bool {
	return b.rawResponse.StatusCode == http.StatusOK
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
			b.rawResponse = resp
		}
		if b.Done() {
			break
		}
		b.rawResponse = resp
		time.Sleep(t)
	}
	return recoverDeletedCertificateResponseFromGenerated(b.recoverResponse), nil
}

// RecoverDeletedCertificateResponse contains response fields for Client.RecoverDeletedCertificate
type RecoverDeletedCertificateResponse struct {
	KeyVaultCertificate
}

// change recover deleted certificate reponse to the generated version.
func recoverDeletedCertificateResponseFromGenerated(i generated.KeyVaultClientRecoverDeletedCertificateResponse) RecoverDeletedCertificateResponse {
	return RecoverDeletedCertificateResponse{
		KeyVaultCertificate: certificateFromGenerated(&i.CertificateBundle),
	}
}

// RecoverDeletedCertificatePollerResponse contains response fields for Client.BeginRecoverDeletedCertificate
type RecoverDeletedCertificatePollerResponse struct {
	// PollUntilDone will poll the service endpoint until a terminal state is reached or an error occurs
	PollUntilDone func(context.Context, time.Duration) (RecoverDeletedCertificateResponse, error)

	// Poller contains an initialized RecoverDeletedCertificatePoller
	Poller RecoverDeletedCertificatePoller
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
		rawResponse:     getResp.RawResponse,
	}

	return RecoverDeletedCertificatePollerResponse{
		PollUntilDone: p.pollUntilDone,
		Poller:        p,
	}, nil
}

// ListDeletedCertificatesPager is the pager returned by Client.ListDeletedCertificates
type ListDeletedCertificatesPager struct {
	vaultURL  string
	genClient *generated.KeyVaultClient
	nextLink  *string
}

// More returns true if there are more pages to return
func (l *ListDeletedCertificatesPager) More() bool {
	if !reflect.ValueOf(l.nextLink).IsZero() {
		if l.nextLink == nil || len(*l.nextLink) == 0 {
			return false
		}
	}
	return true
}

// NextPage returns the current page of results
func (l *ListDeletedCertificatesPager) NextPage(ctx context.Context) (ListDeletedCertificatesPageResponse, error) {
	var resp *http.Response
	var err error
	if l.nextLink == nil {
		req, err := l.genClient.GetDeletedCertificatesCreateRequest(
			ctx,
			l.vaultURL,
			&generated.KeyVaultClientGetDeletedCertificatesOptions{},
		)
		if err != nil {
			return ListDeletedCertificatesPageResponse{}, err
		}
		resp, err = l.genClient.Pl.Do(req)
		if err != nil {
			return ListDeletedCertificatesPageResponse{}, err
		}
	} else {
		req, err := runtime.NewRequest(ctx, http.MethodGet, *l.nextLink)
		if err != nil {
			return ListDeletedCertificatesPageResponse{}, err
		}
		resp, err = l.genClient.Pl.Do(req)
		if err != nil {
			return ListDeletedCertificatesPageResponse{}, err
		}
	}
	if err != nil {
		return ListDeletedCertificatesPageResponse{}, err
	}
	result, err := l.genClient.GetDeletedCertificatesHandleResponse(resp)
	if err != nil {
		return ListDeletedCertificatesPageResponse{}, err
	}
	if result.NextLink == nil {
		// Set it to the zero value
		result.NextLink = to.StringPtr("")
	}
	l.nextLink = result.NextLink
	return listDeletedCertsPageFromGenerated(result), nil
}

// ListDeletedCertificatesPageResponse contains response field for ListDeletedCertificatesPager.NextPage
type ListDeletedCertificatesPageResponse struct {
	// READ-ONLY; A response message containing a list of deleted certificates in the vault along with a link to the next page of deleted certificates
	Certificates []*DeletedCertificateItem `json:"value,omitempty" azure:"ro"`
}

func listDeletedCertsPageFromGenerated(g generated.KeyVaultClientGetDeletedCertificatesResponse) ListDeletedCertificatesPageResponse {
	var certs []*DeletedCertificateItem

	if len(g.Value) > 0 {
		certs = make([]*DeletedCertificateItem, len(g.Value))

		for i, c := range g.Value {
			certs[i] = &DeletedCertificateItem{
				Properties:         certificateAttributesFromGenerated(c.Attributes),
				ID:                 c.ID,
				RecoveryID:         c.RecoveryID,
				Tags:               convertGeneratedMap(c.Tags),
				X509Thumbprint:     c.X509Thumbprint,
				DeletedDate:        c.DeletedDate,
				ScheduledPurgeDate: c.ScheduledPurgeDate,
			}
		}
	}

	return ListDeletedCertificatesPageResponse{
		Certificates: certs,
	}
}

// ListDeletedCertificatesOptions contains optional parameters for Client.ListDeletedCertificates
type ListDeletedCertificatesOptions struct {
	// placeholder for future optional parameters
}

// ListDeletedCertificates retrieves the certificates in the current vault which are in a deleted state and ready for recovery or purging.
// This operation includes deletion-specific information. This operation requires the certificates/get/list permission. This operation can
// only be enabled on soft-delete enabled vaults.
func (c *Client) ListDeletedCertificates(options *ListDeletedCertificatesOptions) ListDeletedCertificatesPager {
	return ListDeletedCertificatesPager{
		vaultURL:  c.vaultURL,
		genClient: c.genClient,
		nextLink:  nil,
	}
}

// CancelCertificateOperationOptions contains optional parameters for Client.CancelCertificateOperation
type CancelCertificateOperationOptions struct {
	// placeholder for future optional parameters.
}

func (c *CancelCertificateOperationOptions) toGenerated() *generated.KeyVaultClientUpdateCertificateOperationOptions {
	return &generated.KeyVaultClientUpdateCertificateOperationOptions{}
}

// CancelCertificateOperationResponse contains response fields for Client.CancelCertificateOperation
type CancelCertificateOperationResponse struct {
	CertificateOperation
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
		CertificateOperation: certificateOperationFromGenerated(resp.CertificateOperation),
	}, nil
}

// DeleteCertificateOperationOptions contains optional parameters for Client.DeleteCertificateOperation
type DeleteCertificateOperationOptions struct {
	// placeholder for future optional parameters.
}

func (d *DeleteCertificateOperationOptions) toGenerated() *generated.KeyVaultClientDeleteCertificateOperationOptions {
	return &generated.KeyVaultClientDeleteCertificateOperationOptions{}
}

// DeleteCertificateOperationResponse contains response fields for Client.DeleteCertificateOperation
type DeleteCertificateOperationResponse struct {
	CertificateOperation
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
		CertificateOperation: certificateOperationFromGenerated(resp.CertificateOperation),
	}, nil
}
