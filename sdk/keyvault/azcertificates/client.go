//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcertificates

import (
	"context"
	"encoding/json"
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
	// Determines whether the object is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// Application specific metadata in the form of key-value pairs
	Tags map[string]string `json:"tags,omitempty"`

	// ResumeToken is a token for resuming long running operations from a previous poller
	ResumeToken *string
}

func (b BeginCreateCertificateOptions) toGenerated() *generated.KeyVaultClientCreateCertificateOptions {
	return &generated.KeyVaultClientCreateCertificateOptions{}
}

// CreateCertificateResponse contains response fields for Client.BeginCreateCertificate
type CreateCertificateResponse struct {
	Operation
}

// CreateCertificatePoller is the poller returned by the Client.BeginCreateCertificate
type CreateCertificatePoller struct {
	certName       string
	certVersion    string
	vaultURL       string
	client         *generated.KeyVaultClient
	createResponse CreateCertificateResponse
	lastResponse   generated.KeyVaultClientGetCertificateResponse
	getRawResponse *http.Response
	resumeToken    string
}

// ResumeToken returns a token for resuming polling at a later time
func (p *CreateCertificatePoller) ResumeToken() (string, error) {
	return p.resumeToken, nil
}

// Done returns true if the LRO has reached a terminal state
func (b *CreateCertificatePoller) Done() bool {
	if b.getRawResponse == nil {
		return false
	}
	return b.getRawResponse.StatusCode == http.StatusOK
}

// Poll fetches the latest state of the operations. It returns an HTTP response or error.
// If the LRO has completed successfully, the poller's state is updated and the HTTP response is returned.
// If the LRO has completed with failure or was cancelled, the poller's state is updated and the error is returned.)
func (b *CreateCertificatePoller) Poll(ctx context.Context) (*http.Response, error) {
	var getRawResp *http.Response
	ctx = runtime.WithCaptureResponse(ctx, &getRawResp)
	resp, err := b.client.GetCertificate(ctx, b.vaultURL, b.certName, b.certVersion, nil)
	if err == nil {
		b.getRawResponse = getRawResp
		b.lastResponse = resp
		b.createResponse.ID = b.lastResponse.ID
		return getRawResp, nil
	}

	if getRawResp != nil && getRawResp.StatusCode == http.StatusNotFound {
		// The certificate has not been fully created yet
		b.getRawResponse = getRawResp
		b.lastResponse = resp
		return b.getRawResponse, nil
	}

	// There was an error in this operation, return the original raw response and the error
	return getRawResp, err
}

// FinalResponse returns the final response after the operations has finished
func (b *CreateCertificatePoller) FinalResponse(ctx context.Context) (CreateCertificateResponse, error) {
	return b.createResponse, nil
}

// PollUntilDone continually polls the service with a 't' delay until completion.
func (b *CreateCertificatePoller) PollUntilDone(ctx context.Context, t time.Duration) (CreateCertificateResponse, error) {
	for {
		resp, err := b.Poll(ctx)
		if err != nil {
			return CreateCertificateResponse{}, err
		}
		b.getRawResponse = resp
		if b.Done() {
			break
		}
		time.Sleep(t)
	}
	return b.createResponse, nil
}

// BeginCreateCertificate creates a new certificate resource, if a certificate with this name already exists, a new version is created. This operation requires the certificates/create permission.
func (c *Client) BeginCreateCertificate(ctx context.Context, certificateName string, policy Policy, options *BeginCreateCertificateOptions) (*CreateCertificatePoller, error) {
	if options == nil {
		options = &BeginCreateCertificateOptions{}
	}

	var tags map[string]*string
	if options.Tags != nil {
		tags = convertToGeneratedMap(options.Tags)
	}

	var createResp generated.KeyVaultClientCreateCertificateResponse
	var err error
	var rt string

	if options.ResumeToken == nil {
		createResp, err = c.genClient.CreateCertificate(
			ctx,
			c.vaultURL,
			certificateName,
			generated.CertificateCreateParameters{
				CertificatePolicy:     policy.toGeneratedCertificateCreateParameters(),
				Tags:                  tags,
				CertificateAttributes: &generated.CertificateAttributes{Enabled: options.Enabled},
			},
			options.toGenerated(),
		)
		if err != nil {
			return nil, err
		}

		marshalled, err := json.Marshal(createResp)
		if err != nil {
			return nil, err
		}
		rt = string(marshalled)
	} else {
		rt = *options.ResumeToken
		err = json.Unmarshal([]byte(rt), &createResp)
		if err != nil {
			return nil, err
		}
	}

	return &CreateCertificatePoller{
		certName:    certificateName,
		certVersion: "",
		vaultURL:    c.vaultURL,
		client:      c.genClient,
		createResponse: CreateCertificateResponse{
			Operation: Operation{
				CancellationRequested: createResp.CancellationRequested,
				CSR:                   createResp.Csr,
				Error:                 certificateErrorFromGenerated(createResp.Error),
				IssuerParameters:      issuerParametersFromGenerated(createResp.IssuerParameters),
				RequestID:             createResp.RequestID,
				Status:                createResp.Status,
				StatusDetails:         createResp.StatusDetails,
				Target:                createResp.Target,
				ID:                    createResp.ID,
			},
		},
		lastResponse: generated.KeyVaultClientGetCertificateResponse{},
		resumeToken:  rt,
	}, nil
}

// GetCertificateOptions contains optional parameters for Client.GetCertificate
type GetCertificateOptions struct {
	Version string
}

// GetCertificateResponse contains response fields for Client.GetCertificate
type GetCertificateResponse struct {
	CertificateWithPolicy
}

// GetCertificate gets information about a specific certificate. This operation requires the certificates/get permission.
func (c *Client) GetCertificate(ctx context.Context, certificateName string, options *GetCertificateOptions) (GetCertificateResponse, error) {
	if options == nil {
		options = &GetCertificateOptions{}
	}

	resp, err := c.genClient.GetCertificate(ctx, c.vaultURL, certificateName, options.Version, nil)
	if err != nil {
		return GetCertificateResponse{}, err
	}

	return GetCertificateResponse{
		CertificateWithPolicy: CertificateWithPolicy{
			Properties:  propertiesFromGenerated(resp.Attributes, convertGeneratedMap(resp.Tags), resp.ID, resp.X509Thumbprint),
			CER:         resp.Cer,
			ContentType: resp.ContentType,
			ID:          resp.ID,
			KeyID:       resp.Kid,
			SecretID:    resp.Sid,
			Policy:      certificatePolicyFromGenerated(resp.Policy),
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
	Operation
}

// GetCertificateOperation gets the creation operation associated with a specified certificate. This operation requires the certificates/get permission.
func (c *Client) GetCertificateOperation(ctx context.Context, certificateName string, options *GetCertificateOperationOptions) (GetCertificateOperationResponse, error) {
	resp, err := c.genClient.GetCertificateOperation(ctx, c.vaultURL, certificateName, options.toGenerated())
	if err != nil {
		return GetCertificateOperationResponse{}, err
	}

	return GetCertificateOperationResponse{
		Operation: Operation{
			CancellationRequested: resp.CancellationRequested,
			CSR:                   resp.Csr,
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
	// ResumeToken is a string to begin polling from a previous operation
	ResumeToken *string
}

// convert public options to generated options struct
func (b *BeginDeleteCertificateOptions) toGenerated() *generated.KeyVaultClientDeleteCertificateOptions {
	return &generated.KeyVaultClientDeleteCertificateOptions{}
}

// DeleteCertificateResponse contains response fields for BeginDeleteCertificatePoller.FinalResponse
type DeleteCertificateResponse struct {
	DeletedCertificate
}

func deleteCertificateResponseFromGenerated(g generated.KeyVaultClientDeleteCertificateResponse) DeleteCertificateResponse {
	_, name, _ := shared.ParseID(g.ID)
	return DeleteCertificateResponse{
		DeletedCertificate: DeletedCertificate{
			RecoveryID:         g.RecoveryID,
			DeletedOn:          g.DeletedDate,
			ScheduledPurgeDate: g.ScheduledPurgeDate,
			Properties:         propertiesFromGenerated(g.Attributes, convertGeneratedMap(g.Tags), g.ID, g.X509Thumbprint),
			CER:                g.Cer,
			ContentType:        g.ContentType,
			ID:                 g.ID,
			Name:               name,
			KeyID:              g.Kid,
			Policy:             certificatePolicyFromGenerated(g.Policy),
			SecretID:           g.Sid,
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
	lastRawResponse *http.Response
	resumeToken     string
}

// ResumeToken returns a token for resuming polling at a later time
func (s *DeleteCertificatePoller) ResumeToken() (string, error) {
	return string(s.resumeToken), nil
}

// Done returns true if the LRO has reached a terminal state
func (s *DeleteCertificatePoller) Done() bool {
	if s.lastRawResponse == nil {
		return false
	}
	return s.lastRawResponse.StatusCode == http.StatusOK
}

// Poll fetches the latest state of the LRO. It returns an HTTP response or error.(
// If the LRO has completed successfully, the poller's state is updated and the HTTP response is returned.
// If the LRO has completed with failure or was cancelled, the poller's state is updated and the error is returned.)
func (s *DeleteCertificatePoller) Poll(ctx context.Context) (*http.Response, error) {
	var getRawResp *http.Response
	ctx = runtime.WithCaptureResponse(ctx, &getRawResp)
	resp, err := s.client.GetDeletedCertificate(ctx, s.vaultURL, s.certificateName, nil)
	if err == nil {
		// Service recognizes DeletedCertificate, operation is done
		s.lastRawResponse = getRawResp
		s.lastResponse = resp
		return s.lastRawResponse, nil
	}

	if getRawResp != nil && getRawResp.StatusCode == http.StatusNotFound {
		// This is the expected result
		s.lastRawResponse = getRawResp
		return s.lastRawResponse, nil
	}
	return s.lastRawResponse, err
}

// FinalResponse returns the final response after the operations has finished
func (s *DeleteCertificatePoller) FinalResponse(ctx context.Context) (DeleteCertificateResponse, error) {
	return deleteCertificateResponseFromGenerated(s.deleteResponse), nil
}

// PollUntilDone continually calls the Poll operation until the operation is completed. In between each
// Poll is a wait determined by the t parameter.
func (s *DeleteCertificatePoller) PollUntilDone(ctx context.Context, t time.Duration) (DeleteCertificateResponse, error) {
	for {
		resp, err := s.Poll(ctx)
		if err != nil {
			return DeleteCertificateResponse{}, err
		}
		s.lastRawResponse = resp
		if s.Done() {
			break
		}
		time.Sleep(t)
	}
	return deleteCertificateResponseFromGenerated(s.deleteResponse), nil
}

// BeginDeleteCertificate deletes a certificate from the keyvault. Delete cannot be applied to an individual version of a certificate. This operation
// requires the certificate/delete permission. This response contains a response with a Poller struct that can be used to Poll for a response, or the
// DeleteCertificatePollerResponse.PollUntilDone function can be used to poll until completion.
func (c *Client) BeginDeleteCertificate(ctx context.Context, certificateName string, options *BeginDeleteCertificateOptions) (*DeleteCertificatePoller, error) {
	if options == nil {
		options = &BeginDeleteCertificateOptions{}
	}
	var resumeToken string
	var delResp generated.KeyVaultClientDeleteCertificateResponse
	var err error
	if options.ResumeToken == nil {
		delResp, err = c.genClient.DeleteCertificate(ctx, c.vaultURL, certificateName, options.toGenerated())
		if err != nil {
			return nil, err
		}

		marshalled, err := json.Marshal(delResp)
		if err != nil {
			return nil, err
		}
		resumeToken = string(marshalled)
	} else {
		resumeToken = *options.ResumeToken
		err = json.Unmarshal([]byte(resumeToken), &delResp)
		if err != nil {
			return nil, err
		}
	}

	getResp, err := c.genClient.GetDeletedCertificate(ctx, c.vaultURL, certificateName, nil)
	var httpErr *azcore.ResponseError
	if errors.As(err, &httpErr) {
		if httpErr.RawResponse.StatusCode != http.StatusNotFound {
			return nil, err
		}
	}

	return &DeleteCertificatePoller{
		vaultURL:        c.vaultURL,
		certificateName: certificateName,
		client:          c.genClient,
		deleteResponse:  delResp,
		lastResponse:    getResp,
		resumeToken:     resumeToken,
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
func (c *Client) PurgeDeletedCertificate(ctx context.Context, certificateName string, options *PurgeDeletedCertificateOptions) (PurgeDeletedCertificateResponse, error) {
	_, err := c.genClient.PurgeDeletedCertificate(ctx, c.vaultURL, certificateName, options.toGenerated())
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
func (c *Client) GetDeletedCertificate(ctx context.Context, certificateName string, options *GetDeletedCertificateOptions) (GetDeletedCertificateResponse, error) {
	resp, err := c.genClient.GetDeletedCertificate(ctx, c.vaultURL, certificateName, options.toGenerated())
	if err != nil {
		return GetDeletedCertificateResponse{}, err
	}

	_, name, _ := shared.ParseID(resp.ID)
	return GetDeletedCertificateResponse{
		DeletedCertificate: DeletedCertificate{
			RecoveryID:         resp.RecoveryID,
			DeletedOn:          resp.DeletedDate,
			ScheduledPurgeDate: resp.ScheduledPurgeDate,
			Properties:         propertiesFromGenerated(resp.Attributes, convertGeneratedMap(resp.Tags), resp.ID, resp.X509Thumbprint),
			CER:                resp.Cer,
			ContentType:        resp.ContentType,
			ID:                 resp.ID,
			Name:               name,
			KeyID:              resp.Kid,
			Policy:             certificatePolicyFromGenerated(resp.Policy),
			SecretID:           resp.Sid,
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
func (c *Client) BackupCertificate(ctx context.Context, certificateName string, options *BackupCertificateOptions) (BackupCertificateResponse, error) {
	resp, err := c.genClient.BackupCertificate(ctx, c.vaultURL, certificateName, options.toGenerated())
	if err != nil {
		return BackupCertificateResponse{}, err
	}

	return BackupCertificateResponse{
		Value: resp.Value,
	}, nil
}

// ImportCertificateOptions contains optional parameters for Client.ImportCertificate
type ImportCertificateOptions struct {
	// The management policy for the certificate.
	CertificatePolicy *Policy `json:"policy,omitempty"`

	// Determines whether the object is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// If the private key in base64EncodedCertificate is encrypted, the password used for encryption.
	Password *string `json:"pwd,omitempty"`

	// Application specific metadata in the form of key-value pairs
	Tags map[string]string `json:"tags,omitempty"`
}

// ImportCertificateResponse contains response fields for Client.ImportCertificate
type ImportCertificateResponse struct {
	CertificateWithPolicy
}

// ImportCertificate imports an existing valid certificate, containing a private key, into Azure Key Vault. This operation requires the
// certificates/import permission. The certificate to be imported can be in either PFX or PEM format. If the certificate is in PEM format
// the PEM file must contain the key as well as x509 certificates. Key Vault will only accept a key in PKCS#8 format.
func (c *Client) ImportCertificate(ctx context.Context, certificateName string, certificate []byte, options *ImportCertificateOptions) (ImportCertificateResponse, error) {
	if options == nil {
		options = &ImportCertificateOptions{}
	}
	var tags map[string]*string
	if options.Tags != nil {
		tags = convertToGeneratedMap(options.Tags)
	}
	resp, err := c.genClient.ImportCertificate(
		ctx,
		c.vaultURL,
		certificateName,
		generated.CertificateImportParameters{
			Base64EncodedCertificate: to.Ptr(string(certificate)),
			CertificateAttributes: &generated.CertificateAttributes{
				Enabled: options.Enabled,
			},
			CertificatePolicy: options.CertificatePolicy.toGeneratedCertificateCreateParameters(),
			Password:          options.Password,
			Tags:              tags,
		},
		&generated.KeyVaultClientImportCertificateOptions{},
	)
	if err != nil {
		return ImportCertificateResponse{}, err
	}

	return ImportCertificateResponse{
		CertificateWithPolicy: CertificateWithPolicy{
			Properties:  propertiesFromGenerated(resp.Attributes, convertGeneratedMap(resp.Tags), resp.ID, resp.X509Thumbprint),
			CER:         resp.Cer,
			ContentType: resp.ContentType,
			ID:          resp.ID,
			KeyID:       resp.Kid,
			SecretID:    resp.Sid,
			Policy:      certificatePolicyFromGenerated(resp.Policy),
		},
	}, nil
}

// ListCertificatesOptions contains optional parameters for Client.ListCertificates
type ListCertificatesOptions struct {
	// placeholder for future optional parameters.
}

// ListPropertiesOfCertificatesResponse contains response fields for ListCertificatesPager.NextPage
type ListPropertiesOfCertificatesResponse struct {
	// READ-ONLY; A response message containing a list of certificates in the key vault along with a link to the next page of certificates.
	Certificates []*CertificateItem `json:"value,omitempty" azure:"ro"`

	// NextLink is a link to the next page of results
	NextLink *string
}

// convert internal Response to ListCertificatesPage
func listCertsPageFromGenerated(i generated.KeyVaultClientGetCertificatesResponse) ListPropertiesOfCertificatesResponse {
	var vals []*CertificateItem

	for _, v := range i.Value {
		vals = append(vals, &CertificateItem{
			Properties: propertiesFromGenerated(v.Attributes, convertGeneratedMap(v.Tags), v.ID, v.X509Thumbprint),
			ID:         v.ID,
		})
	}

	return ListPropertiesOfCertificatesResponse{
		Certificates: vals,
		NextLink:     i.NextLink,
	}
}

// NewListPropertiesOfCertificatesPager retrieves a list of the certificates in the Key Vault as JSON Web Key structures that contain the
// public part of a stored certificate. The LIST operation is applicable to all certificate types, however only the
// base certificate identifier, attributes, and tags are provided in the response. Individual versions of a
// certificate are not listed in the response. This operation requires the certificates/list permission.
func (c *Client) NewListPropertiesOfCertificatesPager(options *ListCertificatesOptions) *runtime.Pager[ListPropertiesOfCertificatesResponse] {
	pager := c.genClient.NewGetCertificatesPager(c.vaultURL, nil)
	return runtime.NewPager(runtime.PagingHandler[ListPropertiesOfCertificatesResponse]{
		More: func(page ListPropertiesOfCertificatesResponse) bool {
			return pager.More()
		},
		Fetcher: func(ctx context.Context, cur *ListPropertiesOfCertificatesResponse) (ListPropertiesOfCertificatesResponse, error) {
			page, err := pager.NextPage(ctx)
			if err != nil {
				return ListPropertiesOfCertificatesResponse{}, err
			}
			return listCertsPageFromGenerated(page), nil
		},
	})
}

// ListCertificateVersionsOptions contains optional parameters for Client.ListCertificateVersions
type ListCertificateVersionsOptions struct {
	// placeholder for future optional parameters.
}

// ListPropertiesOfCertificateVersionsResponse contains response fields for ListCertificateVersionsPager.NextPage
type ListPropertiesOfCertificateVersionsResponse struct {
	// READ-ONLY; A response message containing a list of certificates in the key vault along with a link to the next page of certificates.
	Certificates []*CertificateItem `json:"value,omitempty" azure:"ro"`

	// NextLink is a link to the next page of results to fetch
	NextLink *string
}

// create ListCertificatesPage from generated pager
func listCertificateVersionsPageFromGenerated(i generated.KeyVaultClientGetCertificateVersionsResponse) ListPropertiesOfCertificateVersionsResponse {
	var vals []*CertificateItem
	for _, v := range i.Value {
		vals = append(vals, &CertificateItem{
			Properties: propertiesFromGenerated(v.Attributes, convertGeneratedMap(v.Tags), v.ID, v.X509Thumbprint),
			ID:         v.ID,
		})
	}

	return ListPropertiesOfCertificateVersionsResponse{
		Certificates: vals,
		NextLink:     i.NextLink,
	}
}

// NewListPropertiesOfCertificateVersionsPager lists all versions of the specified certificate. The full certificate identifer and
// attributes are provided in the response. No values are returned for the certificates. This operation
// requires the certificates/list permission.
func (c *Client) NewListPropertiesOfCertificateVersionsPager(certificateName string, options *ListCertificateVersionsOptions) *runtime.Pager[ListPropertiesOfCertificateVersionsResponse] {
	pager := c.genClient.NewGetCertificateVersionsPager(c.vaultURL, certificateName, nil)
	return runtime.NewPager(runtime.PagingHandler[ListPropertiesOfCertificateVersionsResponse]{
		More: func(page ListPropertiesOfCertificateVersionsResponse) bool {
			return pager.More()
		},
		Fetcher: func(ctx context.Context, cur *ListPropertiesOfCertificateVersionsResponse) (ListPropertiesOfCertificateVersionsResponse, error) {
			page, err := pager.NextPage(ctx)
			if err != nil {
				return ListPropertiesOfCertificateVersionsResponse{}, err
			}
			return listCertificateVersionsPageFromGenerated(page), nil
		},
	})
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
	Issuer
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
					EmailAddress: v.Email,
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
	cr.Issuer = Issuer{
		Credentials: issuerCredentialsFromGenerated(resp.Credentials),
		Provider:    resp.Provider,
		ID:          resp.ID,
	}

	if resp.Attributes != nil {
		cr.Issuer.CreatedOn = resp.Attributes.Created
		cr.Issuer.Enabled = resp.Attributes.Enabled
		cr.Issuer.UpdatedOn = resp.Attributes.Updated
	}
	if resp.OrganizationDetails != nil {
		cr.Issuer.OrganizationID = resp.OrganizationDetails.ID
		var adminDetails []*AdministratorContact
		if resp.OrganizationDetails.AdminDetails != nil {
			adminDetails = make([]*AdministratorContact, len(resp.OrganizationDetails.AdminDetails))
			for idx, v := range resp.OrganizationDetails.AdminDetails {
				adminDetails[idx] = &AdministratorContact{
					Email:     v.EmailAddress,
					FirstName: v.FirstName,
					LastName:  v.LastName,
					Phone:     v.Phone,
				}
			}
		}
		cr.Issuer.AdministratorContacts = adminDetails
	}

	_, _, name := shared.ParseID(resp.ID)
	cr.Issuer.Name = name
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
	Issuer
}

// GetIssuer returns the specified certificate issuer resources in the specified key vault. This operation
// requires the certificates/manageissuers/getissuers permission.
func (c *Client) GetIssuer(ctx context.Context, issuerName string, options *GetIssuerOptions) (GetIssuerResponse, error) {
	resp, err := c.genClient.GetCertificateIssuer(ctx, c.vaultURL, issuerName, options.toGenerated())
	if err != nil {
		return GetIssuerResponse{}, err
	}

	g := GetIssuerResponse{}
	g.Issuer = Issuer{
		ID:          resp.ID,
		Provider:    resp.Provider,
		Credentials: issuerCredentialsFromGenerated(resp.Credentials),
	}

	if resp.Attributes != nil {
		g.Issuer.CreatedOn = resp.Attributes.Created
		g.Issuer.Enabled = resp.Attributes.Enabled
		g.Issuer.UpdatedOn = resp.Attributes.Updated
	}
	if resp.OrganizationDetails != nil {
		g.Issuer.OrganizationID = resp.OrganizationDetails.ID
		var adminDetails []*AdministratorContact
		if resp.OrganizationDetails.AdminDetails != nil {
			adminDetails = make([]*AdministratorContact, len(resp.OrganizationDetails.AdminDetails))
			for idx, v := range resp.OrganizationDetails.AdminDetails {
				adminDetails[idx] = &AdministratorContact{
					Email:     v.EmailAddress,
					FirstName: v.FirstName,
					LastName:  v.LastName,
					Phone:     v.Phone,
				}
			}
		}
		g.Issuer.AdministratorContacts = adminDetails
	}

	_, _, name := shared.ParseID(resp.ID)
	g.Issuer.Name = name
	return g, nil
}

// ListPropertiesOfIssuersOptions contains optional parameters for Client.ListIssuers
type ListPropertiesOfIssuersOptions struct {
	// placeholder for future optional parameters
}

// ListIssuersPropertiesOfIssuersResponse contains response fields for ListPropertiesOfIssuersPager.NextPage
type ListIssuersPropertiesOfIssuersResponse struct {
	// READ-ONLY; A response message containing a list of certificates in the key vault along with a link to the next page of certificates.
	Issuers []*IssuerItem `json:"value,omitempty" azure:"ro"`

	// NextLink is the next link of pages to fetch
	NextLink *string
}

// convert internal Response to ListPropertiesOfIssuersPage
func listIssuersPageFromGenerated(i generated.KeyVaultClientGetCertificateIssuersResponse) ListIssuersPropertiesOfIssuersResponse {
	var vals []*IssuerItem

	for _, v := range i.Value {
		vals = append(vals, certificateIssuerItemFromGenerated(v))
	}

	return ListIssuersPropertiesOfIssuersResponse{Issuers: vals, NextLink: i.NextLink}
}

// NewListPropertiesOfIssuersPager returns a pager that can be used to get the set of certificate issuer resources in the specified key vault. This operation
// requires the certificates/manageissuers/getissuers permission.
func (c *Client) NewListPropertiesOfIssuersPager(options *ListPropertiesOfIssuersOptions) *runtime.Pager[ListIssuersPropertiesOfIssuersResponse] {
	pager := c.genClient.NewGetCertificateIssuersPager(c.vaultURL, nil)
	return runtime.NewPager(runtime.PagingHandler[ListIssuersPropertiesOfIssuersResponse]{
		More: func(page ListIssuersPropertiesOfIssuersResponse) bool {
			return pager.More()
		},
		Fetcher: func(ctx context.Context, cur *ListIssuersPropertiesOfIssuersResponse) (ListIssuersPropertiesOfIssuersResponse, error) {
			page, err := pager.NextPage(ctx)
			if err != nil {
				return ListIssuersPropertiesOfIssuersResponse{}, err
			}
			return listIssuersPageFromGenerated(page), nil
		},
	})
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
	Issuer
}

// DeleteIssuer permanently removes the specified certificate issuer from the vault. This operation requires the certificates/manageissuers/deleteissuers permission.
func (c *Client) DeleteIssuer(ctx context.Context, issuerName string, options *DeleteIssuerOptions) (DeleteIssuerResponse, error) {
	resp, err := c.genClient.DeleteCertificateIssuer(ctx, c.vaultURL, issuerName, options.toGenerated())
	if err != nil {
		return DeleteIssuerResponse{}, err
	}

	d := DeleteIssuerResponse{}
	d.Issuer = Issuer{
		ID:          resp.ID,
		Provider:    resp.Provider,
		Credentials: issuerCredentialsFromGenerated(resp.Credentials),
	}

	if resp.Attributes != nil {
		d.Issuer.CreatedOn = resp.Attributes.Created
		d.Issuer.Enabled = resp.Attributes.Enabled
		d.Issuer.UpdatedOn = resp.Attributes.Updated
	}
	if resp.OrganizationDetails != nil {
		d.Issuer.OrganizationID = resp.OrganizationDetails.ID
		var adminDetails []*AdministratorContact
		if resp.OrganizationDetails.AdminDetails != nil {
			adminDetails = make([]*AdministratorContact, len(resp.OrganizationDetails.AdminDetails))
			for idx, v := range resp.OrganizationDetails.AdminDetails {
				adminDetails[idx] = &AdministratorContact{
					Email:     v.EmailAddress,
					FirstName: v.FirstName,
					LastName:  v.LastName,
					Phone:     v.Phone,
				}
			}
		}
		d.Issuer.AdministratorContacts = adminDetails
	}

	_, _, name := shared.ParseID(resp.ID)
	d.Issuer.Name = name
	return d, nil
}

// UpdateIssuerOptions contains optional parameters for Client.UpdateIssuer
type UpdateIssuerOptions struct {
	// placeholder for future optional parameters
}

func (i *Issuer) toUpdateParameters() generated.CertificateIssuerUpdateParameters {
	if i == nil {
		return generated.CertificateIssuerUpdateParameters{}
	}
	var attrib *generated.IssuerAttributes
	if i.Enabled != nil {
		attrib = &generated.IssuerAttributes{Enabled: i.Enabled}
	}

	var orgDetail *generated.OrganizationDetails
	if i.OrganizationID != nil || i.AdministratorContacts != nil {
		orgDetail = &generated.OrganizationDetails{}
		if i.OrganizationID != nil {
			orgDetail.ID = i.OrganizationID
		}

		if i.AdministratorContacts != nil {
			a := make([]*generated.AdministratorDetails, len(i.AdministratorContacts))
			for idx, v := range i.AdministratorContacts {
				a[idx] = &generated.AdministratorDetails{
					EmailAddress: v.Email,
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
		Credentials:         i.Credentials.toGenerated(),
		OrganizationDetails: orgDetail,
		Provider:            i.Provider,
	}
}

// UpdateIssuerResponse contains response fields for Client.UpdateIssuer
type UpdateIssuerResponse struct {
	Issuer
}

// UpdateIssuer performs an update on the specified certificate issuer entity. This operation requires
// the certificates/setissuers permission.
func (c *Client) UpdateIssuer(ctx context.Context, certificateIssuer Issuer, options *UpdateIssuerOptions) (UpdateIssuerResponse, error) {
	resp, err := c.genClient.UpdateCertificateIssuer(
		ctx,
		c.vaultURL,
		*certificateIssuer.Name,
		certificateIssuer.toUpdateParameters(),
		&generated.KeyVaultClientUpdateCertificateIssuerOptions{},
	)
	if err != nil {
		return UpdateIssuerResponse{}, err
	}

	u := UpdateIssuerResponse{}
	u.Issuer = Issuer{
		ID:          resp.ID,
		Provider:    resp.Provider,
		Credentials: issuerCredentialsFromGenerated(resp.Credentials),
	}

	if resp.Attributes != nil {
		u.Issuer.CreatedOn = resp.Attributes.Created
		u.Issuer.Enabled = resp.Attributes.Enabled
		u.Issuer.UpdatedOn = resp.Attributes.Updated
	}
	if resp.OrganizationDetails != nil {
		u.Issuer.OrganizationID = resp.OrganizationDetails.ID
		var adminDetails []*AdministratorContact
		if resp.OrganizationDetails.AdminDetails != nil {
			adminDetails = make([]*AdministratorContact, len(resp.OrganizationDetails.AdminDetails))
			for idx, v := range resp.OrganizationDetails.AdminDetails {
				adminDetails[idx] = &AdministratorContact{
					Email:     v.EmailAddress,
					FirstName: v.FirstName,
					LastName:  v.LastName,
					Phone:     v.Phone,
				}
			}
		}
		u.Issuer.AdministratorContacts = adminDetails
	}
	_, _, name := shared.ParseID(resp.ID)
	u.Issuer.Name = name
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
func (c *Client) SetContacts(ctx context.Context, contacts []*Contact, options *SetContactsOptions) (SetContactsResponse, error) {
	contactList := Contacts{ContactList: contacts}
	resp, err := c.genClient.SetCertificateContacts(
		ctx,
		c.vaultURL,
		contactList.toGenerated(),
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
	Policy
}

// UpdateCertificatePolicy sets specified members in the certificate policy, leave others as null. This operation requires the certificates/update permission.
func (c *Client) UpdateCertificatePolicy(ctx context.Context, certificateName string, policy Policy, options *UpdateCertificatePolicyOptions) (UpdateCertificatePolicyResponse, error) {
	resp, err := c.genClient.UpdateCertificatePolicy(
		ctx,
		c.vaultURL,
		certificateName,
		*policy.toGeneratedCertificateCreateParameters(),
		options.toGenerated(),
	)

	if err != nil {
		return UpdateCertificatePolicyResponse{}, err
	}

	return UpdateCertificatePolicyResponse{
		Policy: *certificatePolicyFromGenerated(&resp.CertificatePolicy),
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
	Policy
}

// GetCertificatePolicy returns the specified certificate policy resources in the specified key vault. This operation requires the certificates/get permission.
func (c *Client) GetCertificatePolicy(ctx context.Context, certificateName string, options *GetCertificatePolicyOptions) (GetCertificatePolicyResponse, error) {
	resp, err := c.genClient.GetCertificatePolicy(
		ctx,
		c.vaultURL,
		certificateName,
		options.toGenerated(),
	)
	if err != nil {
		return GetCertificatePolicyResponse{}, err
	}

	return GetCertificatePolicyResponse{
		Policy: *certificatePolicyFromGenerated(&resp.CertificatePolicy),
	}, nil
}

// UpdateCertificatePropertiesOptions contains optional parameters for Client.UpdateCertificateProperties
type UpdateCertificatePropertiesOptions struct {
	// The version of the certificate to update
	Version string

	// The attributes of the certificate (optional).
	Properties *Properties `json:"attributes,omitempty"`

	// The management policy for the certificate.
	CertificatePolicy *Policy `json:"policy,omitempty"`
}

func (u *UpdateCertificatePropertiesOptions) toGenerated() *generated.KeyVaultClientUpdateCertificateOptions {
	return &generated.KeyVaultClientUpdateCertificateOptions{}
}

// UpdateCertificatePropertiesResponse contains response fields for Client.UpdateCertificateProperties
type UpdateCertificatePropertiesResponse struct {
	Certificate
}

// UpdateCertificateProperties applies the specified update on the given certificate; the only elements updated are the certificate's
// attributes. This operation requires the certificates/update permission.
func (c *Client) UpdateCertificateProperties(ctx context.Context, certificateName string, properties Properties, options *UpdateCertificatePropertiesOptions) (UpdateCertificatePropertiesResponse, error) {
	if options == nil {
		options = &UpdateCertificatePropertiesOptions{}
	}
	var tags map[string]*string
	if properties.Tags != nil {
		tags = convertToGeneratedMap(properties.Tags)
	}
	resp, err := c.genClient.UpdateCertificate(
		ctx,
		c.vaultURL,
		certificateName,
		options.Version,
		generated.CertificateUpdateParameters{
			CertificateAttributes: properties.toGenerated(),
			Tags:                  tags,
		},
		options.toGenerated(),
	)
	if err != nil {
		return UpdateCertificatePropertiesResponse{}, err
	}
	return UpdateCertificatePropertiesResponse{
		Certificate: certificateFromGenerated(&resp.CertificateBundle),
	}, nil
}

// MergeCertificateOptions contains optional parameters for Client.MergeCertificate
type MergeCertificateOptions struct {
	// The attributes of the certificate (optional).
	Properties *Properties `json:"attributes,omitempty"`
}

func (m *MergeCertificateOptions) toGenerated() *generated.KeyVaultClientMergeCertificateOptions {
	return &generated.KeyVaultClientMergeCertificateOptions{}
}

// MergeCertificateResponse contains response fields for Client.MergeCertificate
type MergeCertificateResponse struct {
	CertificateWithPolicy
}

// MergeCertificate operation performs the merging of a certificate or certificate chain with a key pair currently available in the service. This operation requires the certificates/create permission.
func (c *Client) MergeCertificate(ctx context.Context, certificateName string, certificates [][]byte, options *MergeCertificateOptions) (MergeCertificateResponse, error) {
	if options == nil {
		options = &MergeCertificateOptions{}
	}
	var tags map[string]*string
	if options.Properties != nil && options.Properties.Tags != nil {
		tags = convertToGeneratedMap(options.Properties.Tags)
	}
	resp, err := c.genClient.MergeCertificate(
		ctx, c.vaultURL,
		certificateName,
		generated.CertificateMergeParameters{
			X509Certificates:      certificates,
			CertificateAttributes: options.Properties.toGenerated(),
			Tags:                  tags,
		},
		options.toGenerated(),
	)
	if err != nil {
		return MergeCertificateResponse{}, err
	}

	return MergeCertificateResponse{
		CertificateWithPolicy: CertificateWithPolicy{
			Properties:  propertiesFromGenerated(resp.Attributes, convertGeneratedMap(resp.Tags), resp.ID, resp.X509Thumbprint),
			CER:         resp.Cer,
			ContentType: resp.ContentType,
			ID:          resp.ID,
			KeyID:       resp.Kid,
			SecretID:    resp.Sid,
			Policy:      certificatePolicyFromGenerated(resp.Policy),
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
	CertificateWithPolicy
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
		CertificateWithPolicy: CertificateWithPolicy{
			Properties:  propertiesFromGenerated(resp.Attributes, convertGeneratedMap(resp.Tags), resp.ID, resp.X509Thumbprint),
			CER:         resp.Cer,
			ContentType: resp.ContentType,
			ID:          resp.ID,
			KeyID:       resp.Kid,
			SecretID:    resp.Sid,
			Policy:      certificatePolicyFromGenerated(resp.Policy),
		},
	}, nil
}

// BeginRecoverDeletedCertificateOptions contains optional parameters for Client.BeginRecoverDeletedCertificate
type BeginRecoverDeletedCertificateOptions struct {
	// ResumeToken is a token for resuming long running operations from a previous call.
	ResumeToken *string
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
	lastRawResponse *http.Response
	resumeToken     string
}

// ResumeToken returns a token for resuming polling at a later time
func (p *RecoverDeletedCertificatePoller) ResumeToken() (string, error) {
	return p.resumeToken, nil
}

// Done returns true when the polling operation is completed
func (b *RecoverDeletedCertificatePoller) Done() bool {
	if b.lastRawResponse == nil {
		return false
	}
	return b.lastRawResponse.StatusCode == http.StatusOK
}

// Poll fetches the latest state of the LRO. It returns an HTTP response or error.
// If the LRO has completed successfully, the poller's state is updated and the HTTP response is returned.
// If the LRO has completed with failure or was cancelled, the poller's state is updated and the error is returned.
func (b *RecoverDeletedCertificatePoller) Poll(ctx context.Context) (*http.Response, error) {
	var getRawResp *http.Response
	ctx = runtime.WithCaptureResponse(ctx, &getRawResp)
	resp, err := b.client.GetCertificate(ctx, b.vaultUrl, b.certName, "", nil)
	if err == nil {
		// Service has recovered certificate, operation is done
		b.lastRawResponse = getRawResp
		b.lastResponse = resp
		return b.lastRawResponse, nil
	}

	if getRawResp != nil && getRawResp.StatusCode == http.StatusNotFound {
		// This is our expected result
		b.lastRawResponse = getRawResp
		return b.lastRawResponse, nil
	}

	return getRawResp, err
}

// FinalResponse returns the final response after the operations has finished
func (b *RecoverDeletedCertificatePoller) FinalResponse(ctx context.Context) (RecoverDeletedCertificateResponse, error) {
	return recoverDeletedCertificateResponseFromGenerated(b.recoverResponse), nil
}

// PollUntilDone is the method for the Response.PollUntilDone struct
func (b *RecoverDeletedCertificatePoller) PollUntilDone(ctx context.Context, t time.Duration) (RecoverDeletedCertificateResponse, error) {
	for {
		resp, err := b.Poll(ctx)
		if err != nil {
			b.lastRawResponse = resp
		}
		if b.Done() {
			break
		}
		b.lastRawResponse = resp
		time.Sleep(t)
	}
	return recoverDeletedCertificateResponseFromGenerated(b.recoverResponse), nil
}

// RecoverDeletedCertificateResponse contains response fields for Client.RecoverDeletedCertificate
type RecoverDeletedCertificateResponse struct {
	Certificate
}

// change recover deleted certificate reponse to the generated version.
func recoverDeletedCertificateResponseFromGenerated(i generated.KeyVaultClientRecoverDeletedCertificateResponse) RecoverDeletedCertificateResponse {
	return RecoverDeletedCertificateResponse{
		Certificate: certificateFromGenerated(&i.CertificateBundle),
	}
}

// BeginRecoverDeletedCertificate recovers the deleted certificate in the specified vault to the latest version.
// This operation can only be performed on a soft-delete enabled vault. This operation requires the certificates/recover permission.
func (c *Client) BeginRecoverDeletedCertificate(ctx context.Context, certificateName string, options *BeginRecoverDeletedCertificateOptions) (*RecoverDeletedCertificatePoller, error) {
	if options == nil {
		options = &BeginRecoverDeletedCertificateOptions{}
	}
	var recoverResp generated.KeyVaultClientRecoverDeletedCertificateResponse
	var resumeToken string
	var err error
	if options.ResumeToken == nil {
		recoverResp, err = c.genClient.RecoverDeletedCertificate(ctx, c.vaultURL, certificateName, options.toGenerated())
		if err != nil {
			return nil, err
		}

		marshalled, err := json.Marshal(recoverResp)
		if err != nil {
			return nil, err
		}
		resumeToken = string(marshalled)
	} else {
		resumeToken = *options.ResumeToken
		err = json.Unmarshal([]byte(resumeToken), &recoverResp)
		if err != nil {
			return nil, err
		}
	}

	getResp, err := c.genClient.GetCertificate(ctx, c.vaultURL, certificateName, "", nil)
	var httpErr *azcore.ResponseError
	if errors.As(err, &httpErr) {
		if httpErr.RawResponse.StatusCode != http.StatusNotFound {
			return nil, err
		}
	}

	return &RecoverDeletedCertificatePoller{
		lastResponse:    getResp,
		certName:        certificateName,
		client:          c.genClient,
		vaultUrl:        c.vaultURL,
		recoverResponse: recoverResp,
		resumeToken:     resumeToken,
	}, nil
}

// ListDeletedCertificatesResponse contains response field for ListDeletedCertificatesPager.NextPage
type ListDeletedCertificatesResponse struct {
	// READ-ONLY; A response message containing a list of deleted certificates in the vault along with a link to the next page of deleted certificates
	Certificates []*DeletedCertificateItem `json:"value,omitempty" azure:"ro"`

	// NextLink gives the next page of items to fetch
	NextLink *string
}

func listDeletedCertsPageFromGenerated(g generated.KeyVaultClientGetDeletedCertificatesResponse) ListDeletedCertificatesResponse {
	var certs []*DeletedCertificateItem

	if len(g.Value) > 0 {
		certs = make([]*DeletedCertificateItem, len(g.Value))

		for i, c := range g.Value {
			_, name, _ := shared.ParseID(c.ID)
			certs[i] = &DeletedCertificateItem{
				Properties:         propertiesFromGenerated(c.Attributes, convertGeneratedMap(c.Tags), c.ID, c.X509Thumbprint),
				ID:                 c.ID,
				Name:               name,
				RecoveryID:         c.RecoveryID,
				DeletedOn:          c.DeletedDate,
				ScheduledPurgeDate: c.ScheduledPurgeDate,
			}
		}
	}

	return ListDeletedCertificatesResponse{
		Certificates: certs,
		NextLink:     g.NextLink,
	}
}

// ListDeletedCertificatesOptions contains optional parameters for Client.ListDeletedCertificates
type ListDeletedCertificatesOptions struct {
	// placeholder for future optional parameters
}

// NewListDeletedCertificatesPager retrieves the certificates in the current vault which are in a deleted state and ready for recovery or purging.
// This operation includes deletion-specific information. This operation requires the certificates/get/list permission. This operation can
// only be enabled on soft-delete enabled vaults.
func (c *Client) NewListDeletedCertificatesPager(options *ListDeletedCertificatesOptions) *runtime.Pager[ListDeletedCertificatesResponse] {
	pager := c.genClient.NewGetDeletedCertificatesPager(c.vaultURL, nil)
	return runtime.NewPager(runtime.PagingHandler[ListDeletedCertificatesResponse]{
		More: func(page ListDeletedCertificatesResponse) bool {
			return pager.More()
		},
		Fetcher: func(ctx context.Context, cur *ListDeletedCertificatesResponse) (ListDeletedCertificatesResponse, error) {
			page, err := pager.NextPage(ctx)
			if err != nil {
				return ListDeletedCertificatesResponse{}, err
			}
			return listDeletedCertsPageFromGenerated(page), nil
		},
	})
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
	Operation
}

// CancelCertificateOperation cancels a certificate creation operation that is already in progress. This operation requires the certificates/update permission.
func (c *Client) CancelCertificateOperation(ctx context.Context, certificateName string, options *CancelCertificateOperationOptions) (CancelCertificateOperationResponse, error) {
	resp, err := c.genClient.UpdateCertificateOperation(
		ctx,
		c.vaultURL,
		certificateName,
		generated.CertificateOperationUpdateParameter{
			CancellationRequested: to.Ptr(true),
		},
		options.toGenerated(),
	)
	if err != nil {
		return CancelCertificateOperationResponse{}, err
	}

	return CancelCertificateOperationResponse{
		Operation: certificateOperationFromGenerated(resp.CertificateOperation),
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
	Operation
}

// DeleteCertificateOperation deletes the creation operation for a specified certificate that is in the process of being created. The certificate is no
// longer created. This operation requires the certificates/update permission.
func (c *Client) DeleteCertificateOperation(ctx context.Context, certificateName string, options *DeleteCertificateOperationOptions) (DeleteCertificateOperationResponse, error) {
	resp, err := c.genClient.DeleteCertificateOperation(
		ctx,
		c.vaultURL,
		certificateName,
		options.toGenerated(),
	)

	if err != nil {
		return DeleteCertificateOperationResponse{}, err
	}

	return DeleteCertificateOperationResponse{
		Operation: certificateOperationFromGenerated(resp.CertificateOperation),
	}, nil
}
