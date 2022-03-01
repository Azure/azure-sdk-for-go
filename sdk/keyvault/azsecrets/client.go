//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsecrets

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets/internal"
	shared "github.com/Azure/azure-sdk-for-go/sdk/keyvault/internal"
)

// Client is the struct for interacting with a KeyVault Secrets instance
type Client struct {
	kvClient *internal.KeyVaultClient
	vaultUrl string
}

// ClientOptions are the configurable options on a Client.
type ClientOptions struct {
	azcore.ClientOptions
}

func (c *ClientOptions) toConnectionOptions() *policy.ClientOptions {
	if c == nil {
		return nil
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

// NewClient returns a pointer to a Client object affinitized to a vaultUrl.
func NewClient(vaultUrl string, credential azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	conOptions := options.toConnectionOptions()

	conOptions.PerRetryPolicies = append(
		conOptions.PerRetryPolicies,
		shared.NewKeyVaultChallengePolicy(credential),
	)

	return &Client{
		kvClient: internal.NewKeyVaultClient(conOptions),
		vaultUrl: vaultUrl,
	}, nil
}

// GetSecretOptions holds the optional parameters for the Client.GetSecret function
type GetSecretOptions struct {
	// Version specifies the version of a secret. If unspecified, the most recent version will be returned
	Version string
}

// convert the exposed options struct to the internal one.
func (g *GetSecretOptions) toGenerated() *internal.KeyVaultClientGetSecretOptions {
	if g == nil {
		return &internal.KeyVaultClientGetSecretOptions{}
	}
	return &internal.KeyVaultClientGetSecretOptions{}
}

// GetSecretResponse is the response object for the Client.GetSecret operation
type GetSecretResponse struct {
	Secret

	// RawResponse contains the underlying HTTP response
	RawResponse *http.Response
}

func getSecretResponseFromGenerated(i internal.KeyVaultClientGetSecretResponse) *GetSecretResponse {
	return &GetSecretResponse{
		RawResponse: i.RawResponse,
		Secret: Secret{
			Properties:  secretPropertiesFromGenerated(i.Attributes),
			ContentType: i.ContentType,
			ID:          i.ID,
			Tags:        convertPtrMap(i.Tags),
			Value:       i.Value,
			KeyID:       i.Kid,
			Managed:     i.Managed,
		},
	}
}

// GetSecret gets a specified secret from a given key vault. The GET operation is applicable to any secret
// stored in Azure Key Vault. This operation requires the secrets/get permission
func (c *Client) GetSecret(ctx context.Context, secretName string, options *GetSecretOptions) (GetSecretResponse, error) {
	if options == nil {
		options = &GetSecretOptions{}
	}
	resp, err := c.kvClient.GetSecret(ctx, c.vaultUrl, secretName, options.Version, options.toGenerated())
	return *getSecretResponseFromGenerated(resp), err
}

// SetSecretOptions contains the optional parameters for a Client.SetSecret operation
type SetSecretOptions struct {
	// Type of the secret value such as a password.
	ContentType *string `json:"contentType,omitempty"`

	// The secret management attributes.
	Properties *Properties `json:"attributes,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]string `json:"tags,omitempty"`
}

// Convert the exposed struct to the generated code version
func (s *SetSecretOptions) toGenerated() *internal.KeyVaultClientSetSecretOptions {
	if s == nil {
		return nil
	}
	return &internal.KeyVaultClientSetSecretOptions{}
}

// SetSecretResponse is the response struct for the Client.SetSecret operation.
type SetSecretResponse struct {
	Secret

	// RawResponse holds the underlying HTTP response
	RawResponse *http.Response
}

// convert generated response to publicly exposed response.
func setSecretResponseFromGenerated(i internal.KeyVaultClientSetSecretResponse) SetSecretResponse {
	return SetSecretResponse{
		RawResponse: i.RawResponse,
		Secret: Secret{
			Properties:  secretPropertiesFromGenerated(i.Attributes),
			ContentType: i.ContentType,
			ID:          i.ID,
			Tags:        convertPtrMap(i.Tags),
			Value:       i.Value,
			KeyID:       i.Kid,
			Managed:     i.Managed,
		},
	}
}

// SetSecret sets a secret in a specifed key vault. The set operation adds a secret to the Azure Key Vault, if the named secret
// already exists, Azure Key Vault creates a new version of that secret. This operation requires the secrets/set permission.
func (c *Client) SetSecret(ctx context.Context, secretName string, value string, options *SetSecretOptions) (SetSecretResponse, error) {
	if options == nil {
		options = &SetSecretOptions{}
	}
	var secretAttribs internal.SecretAttributes
	if options.Properties != nil {
		secretAttribs = *options.Properties.toGenerated()
	}
	resp, err := c.kvClient.SetSecret(ctx, c.vaultUrl, secretName, internal.SecretSetParameters{
		Value:            &value,
		ContentType:      options.ContentType,
		SecretAttributes: &secretAttribs,
		Tags:             convertToGeneratedMap(options.Tags),
	}, options.toGenerated())
	return setSecretResponseFromGenerated(resp), err
}

// DeletedSecretResponse contains the response for a Client.DeleteSecret operation.
type DeleteSecretResponse struct {
	DeletedSecret

	// RawResponse holds the underlying HTTP response
	RawResponse *http.Response
}

func deleteSecretResponseFromGenerated(i *internal.KeyVaultClientDeleteSecretResponse) *DeleteSecretResponse {
	if i == nil {
		return nil
	}
	return &DeleteSecretResponse{
		DeletedSecret: DeletedSecret{
			ContentType: i.ContentType,
			ID:          i.ID,
			Tags:        convertPtrMap(i.Tags),
			Value:       i.Value,
			KeyID:       i.Kid,
			Managed:     i.Managed,
			Properties: &Properties{
				Enabled:         i.Attributes.Enabled,
				ExpiresOn:       i.Attributes.Expires,
				NotBefore:       i.Attributes.NotBefore,
				CreatedOn:       i.Attributes.Created,
				UpdatedOn:       i.Attributes.Updated,
				RecoverableDays: i.Attributes.RecoverableDays,
				RecoveryLevel:   deletionRecoveryLevelFromGenerated(*i.Attributes.RecoveryLevel).ToPtr(),
			},
			RecoveryID:         i.RecoveryID,
			DeletedOn:          i.DeletedDate,
			ScheduledPurgeDate: i.ScheduledPurgeDate,
		},
		RawResponse: i.RawResponse,
	}
}

// BeginDeleteSecretOptions contains the optional parameters for the Client.BeginDeleteSecret method.
type BeginDeleteSecretOptions struct {
	// placeholder for future optional parameters
}

// convert public options to generated options struct
func (b *BeginDeleteSecretOptions) toGenerated() *internal.KeyVaultClientDeleteSecretOptions {
	return &internal.KeyVaultClientDeleteSecretOptions{}
}

// The poller returned by the Client.StartDeleteSecret operation
type DeleteSecretPoller struct {
	secretName     string // This is the secret to Poll for in GetDeletedSecret
	vaultUrl       string
	client         *internal.KeyVaultClient
	deleteResponse internal.KeyVaultClientDeleteSecretResponse
	lastResponse   internal.KeyVaultClientGetDeletedSecretResponse
	RawResponse    *http.Response
}

// Done returns true if the LRO has reached a terminal state
func (s *DeleteSecretPoller) Done() bool {
	return s.lastResponse.RawResponse != nil
}

// Poll fetches the latest state of the LRO. It returns an HTTP response or error.(
// If the LRO has completed successfully, the poller's state is updated and the HTTP response is returned.
// If the LRO has completed with failure or was cancelled, the poller's state is updated and the error is returned.)
func (s *DeleteSecretPoller) Poll(ctx context.Context) (*http.Response, error) {
	resp, err := s.client.GetDeletedSecret(ctx, s.vaultUrl, s.secretName, nil)
	if err == nil {
		// Service recognizes DeletedSecret, operation is done
		s.lastResponse = resp
		return resp.RawResponse, nil
	}
	var httpErr *azcore.ResponseError
	if errors.As(err, &httpErr) {
		if httpErr.StatusCode == http.StatusNotFound {
			// This is the expected result
			return s.deleteResponse.RawResponse, nil
		}
	}
	return s.deleteResponse.RawResponse, err
}

// FinalResponse returns the final response after the operations has finished
func (s *DeleteSecretPoller) FinalResponse(ctx context.Context) (DeleteSecretResponse, error) {
	return *deleteSecretResponseFromGenerated(&s.deleteResponse), nil
}

// pollUntilDone continually calls the Poll operation until the operation is completed. In between each
// Poll is a wait determined by the t parameter.
func (s *DeleteSecretPoller) pollUntilDone(ctx context.Context, t time.Duration) (DeleteSecretResponse, error) {
	for {
		resp, err := s.Poll(ctx)
		if err != nil {
			return DeleteSecretResponse{}, err
		}
		s.RawResponse = resp
		if s.Done() {
			break
		}
		time.Sleep(t)
	}
	return *deleteSecretResponseFromGenerated(&s.deleteResponse), nil
}

type DeleteSecretPollerResponse struct {
	// PollUntilDone will poll the service endpoint until a terminal state is reached or an error occurs
	PollUntilDone func(context.Context, time.Duration) (DeleteSecretResponse, error)

	// Poller contains an initialized WidgetPoller
	Poller DeleteSecretPoller

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// BeginDeleteSecret deletes a secret from the keyvault. Delete cannot be applied to an individual version of a secret. This operation
// requires the secrets/delete permission. This response contains a Poller struct that can be used to Poll for a response, or the
// response PollUntilDone function can be used to poll until completion.
func (c *Client) BeginDeleteSecret(ctx context.Context, secretName string, options *BeginDeleteSecretOptions) (DeleteSecretPollerResponse, error) {
	// TODO: this is kvSecretClient.DeleteSecret and a GetDeletedSecret under the hood for the polling version
	if options == nil {
		options = &BeginDeleteSecretOptions{}
	}
	resp, err := c.kvClient.DeleteSecret(ctx, c.vaultUrl, secretName, options.toGenerated())
	if err != nil {
		return DeleteSecretPollerResponse{}, err
	}

	getResp, err := c.kvClient.GetDeletedSecret(ctx, c.vaultUrl, secretName, nil)
	var httpErr *azcore.ResponseError
	if errors.As(err, &httpErr) {
		if httpErr.StatusCode != http.StatusNotFound {
			return DeleteSecretPollerResponse{}, err
		}
	}

	s := DeleteSecretPoller{
		vaultUrl:       c.vaultUrl,
		secretName:     secretName,
		client:         c.kvClient,
		deleteResponse: resp,
		lastResponse:   getResp,
	}

	return DeleteSecretPollerResponse{
		Poller:        s,
		RawResponse:   resp.RawResponse,
		PollUntilDone: s.pollUntilDone,
	}, nil
}

// GetDeletedSecretOptions contains the optional parameters for the Client.GetDeletedSecret method.
type GetDeletedSecretOptions struct {
	// placeholder for future optional parameters
}

func (g *GetDeletedSecretOptions) toGenerated() *internal.KeyVaultClientGetDeletedSecretOptions {
	return &internal.KeyVaultClientGetDeletedSecretOptions{}
}

// GetDeletedSecretResponse contains the response struct for the Client.GetDeletedSecret operation.
type GetDeletedSecretResponse struct {
	DeletedSecret

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// Convert the generated response to the publicly exposed version
func getDeletedSecretResponseFromGenerated(i internal.KeyVaultClientGetDeletedSecretResponse) GetDeletedSecretResponse {
	return GetDeletedSecretResponse{
		RawResponse: i.RawResponse,
		DeletedSecret: DeletedSecret{
			Properties:         secretPropertiesFromGenerated(i.Attributes),
			ContentType:        i.ContentType,
			ID:                 i.ID,
			RecoveryID:         i.RecoveryID,
			Tags:               convertPtrMap(i.Tags),
			Value:              i.Value,
			DeletedOn:          i.DeletedDate,
			KeyID:              i.Kid,
			Managed:            i.Managed,
			ScheduledPurgeDate: i.ScheduledPurgeDate,
		},
	}
}

// GetDeletedSecret gets the specified deleted secret. The operation returns the deleted secret along with its attributes.
// This operation requires the secrets/get permission.
func (c *Client) GetDeletedSecret(ctx context.Context, secretName string, options *GetDeletedSecretOptions) (GetDeletedSecretResponse, error) {
	if options == nil {
		options = &GetDeletedSecretOptions{}
	}
	resp, err := c.kvClient.GetDeletedSecret(ctx, c.vaultUrl, secretName, options.toGenerated())
	return getDeletedSecretResponseFromGenerated(resp), err
}

// UpdateSecretPropertiesOptions contains the optional parameters for the Client.UpdateSecretProperties method.
type UpdateSecretPropertiesOptions struct {
	// Version is the specific version of a Secret to update. If not specified it will update the most recent version.
	Version string

	// Type of the secret value such as a password.
	ContentType *string `json:"contentType,omitempty"`

	// The secret management attributes.
	Properties *Properties `json:"attributes,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]string `json:"tags,omitempty"`
}

func (u UpdateSecretPropertiesOptions) toGenerated() *internal.KeyVaultClientUpdateSecretOptions {
	return &internal.KeyVaultClientUpdateSecretOptions{}
}

func (u UpdateSecretPropertiesOptions) toGeneratedProperties() internal.SecretUpdateParameters {
	return internal.SecretUpdateParameters{
		ContentType:      u.ContentType,
		SecretAttributes: u.Properties.toGenerated(),
		Tags:             convertToGeneratedMap(u.Tags),
	}
}

// UpdateSecretPropertiesResponse contains the underlying response object for the UpdateSecretProperties method
type UpdateSecretPropertiesResponse struct {
	Secret
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func updateSecretPropertiesResponseFromGenerated(i internal.KeyVaultClientUpdateSecretResponse) UpdateSecretPropertiesResponse {
	return UpdateSecretPropertiesResponse{
		RawResponse: i.RawResponse,
		Secret: Secret{
			Properties:  secretPropertiesFromGenerated(i.Attributes),
			ContentType: i.ContentType,
			ID:          i.ID,
			Tags:        convertPtrMap(i.Tags),
			Value:       i.Value,
			KeyID:       i.Kid,
			Managed:     i.Managed,
		},
	}
}

// UpdateSecretProperties updates the attributes associated with a specified secret in a given key vault. The update
// operation changes specified attributes of an existing stored secret, attributes that are not specified in the
// request are left unchanged. The value of a secret itself cannot be changed. This operation requires the secrets/set permission.
func (c *Client) UpdateSecretProperties(ctx context.Context, secretName string, options *UpdateSecretPropertiesOptions) (UpdateSecretPropertiesResponse, error) {
	if options == nil {
		options = &UpdateSecretPropertiesOptions{}
	}

	resp, err := c.kvClient.UpdateSecret(
		ctx,
		c.vaultUrl,
		secretName,
		options.Version,
		options.toGeneratedProperties(),
		options.toGenerated(),
	)
	if err != nil {
		return UpdateSecretPropertiesResponse{}, err
	}

	return updateSecretPropertiesResponseFromGenerated(resp), err
}

// BackupSecretOptions contains the optional parameters for the Client.BackupSecret method.
type BackupSecretOptions struct {
	// placeholder for future optional parameters
}

func (b *BackupSecretOptions) toGenerated() *internal.KeyVaultClientBackupSecretOptions {
	return &internal.KeyVaultClientBackupSecretOptions{}
}

// BackupSecretResponse contains the response object for the Client.BackupSecret method.
type BackupSecretResponse struct {
	// READ-ONLY; The backup blob containing the backed up secret.
	Value []byte `json:"value,omitempty" azure:"ro"`

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// convert generated response to the publicly exposed version.
func backupSecretResponseFromGenerated(i internal.KeyVaultClientBackupSecretResponse) BackupSecretResponse {
	return BackupSecretResponse{
		RawResponse: i.RawResponse,
		Value:       i.Value,
	}
}

// BackupSecrets backs up the specified secret. Requests that a backup of the specified secret be downloaded to the client.
// All versions of the secret will be downloaded. This operation requires the secrets/backup permission.
func (c *Client) BackupSecret(ctx context.Context, secretName string, options *BackupSecretOptions) (BackupSecretResponse, error) {
	if options == nil {
		options = &BackupSecretOptions{}
	}

	resp, err := c.kvClient.BackupSecret(ctx, c.vaultUrl, secretName, options.toGenerated())
	if err != nil {
		return BackupSecretResponse{}, err
	}

	return backupSecretResponseFromGenerated(resp), nil
}

// RestoreSecretBackupOptions contains the optional parameters for the Client.RestoreSecret method.
type RestoreSecretBackupOptions struct {
	// placeholder for future optional parameters
}

func (r RestoreSecretBackupOptions) toGenerated() *internal.KeyVaultClientRestoreSecretOptions {
	return &internal.KeyVaultClientRestoreSecretOptions{}
}

// RestoreSecretBackupResponse contains the response object for the Client.RestoreSecretBackup operation.
type RestoreSecretBackupResponse struct {
	Secret

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// converts the generated response to the publicly exposed version.
func restoreSecretBackupResponseFromGenerated(i internal.KeyVaultClientRestoreSecretResponse) RestoreSecretBackupResponse {
	return RestoreSecretBackupResponse{
		RawResponse: i.RawResponse,
		Secret: Secret{
			ContentType: i.ContentType,
			ID:          i.ID,
			Tags:        convertPtrMap(i.Tags),
			Value:       i.Value,
			KeyID:       i.Kid,
			Managed:     i.Managed,
			Properties: &Properties{
				Enabled:         i.Attributes.Enabled,
				ExpiresOn:       i.Attributes.Expires,
				NotBefore:       i.Attributes.NotBefore,
				CreatedOn:       i.Attributes.Created,
				UpdatedOn:       i.Attributes.Updated,
				RecoverableDays: i.Attributes.RecoverableDays,
				RecoveryLevel:   deletionRecoveryLevelFromGenerated(*i.Attributes.RecoveryLevel).ToPtr(),
			},
		},
	}
}

// RestoreSecretBackup restores a backed up secret, and all its versions, to a vault. This operation requires the secrets/restore permission.
// The backup parameter is a blob of the secret to restore, this can be received from the Client.BackupSecret function.
func (c *Client) RestoreSecretBackup(ctx context.Context, backup []byte, options *RestoreSecretBackupOptions) (RestoreSecretBackupResponse, error) {
	if options == nil {
		options = &RestoreSecretBackupOptions{}
	}

	resp, err := c.kvClient.RestoreSecret(ctx, c.vaultUrl, internal.SecretRestoreParameters{SecretBundleBackup: backup}, options.toGenerated())
	if err != nil {
		return RestoreSecretBackupResponse{}, err
	}

	return restoreSecretBackupResponseFromGenerated(resp), nil
}

// PurgeDeletedSecretOptions is the struct for any future options for Client.PurgeDeletedSecret.
type PurgeDeletedSecretOptions struct {
	// placeholder for future optional parameters
}

func (p *PurgeDeletedSecretOptions) toGenerated() *internal.KeyVaultClientPurgeDeletedSecretOptions {
	return &internal.KeyVaultClientPurgeDeletedSecretOptions{}
}

// PurgeDeletedSecretResponse contains the response from method Client.PurgeDeletedSecret.
type PurgeDeletedSecretResponse struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// Converts the generated response to the publicly exposed version.
func purgeDeletedSecretResponseFromGenerated(i internal.KeyVaultClientPurgeDeletedSecretResponse) PurgeDeletedSecretResponse {
	return PurgeDeletedSecretResponse{
		RawResponse: i.RawResponse,
	}
}

// PurgeDeletedSecret deletes the specified secret. The purge deleted secret operation removes the secret permanently, without the possibility of recovery.
// This operation can only be enabled on a soft-delete enabled vault. This operation requires the secrets/purge permission.
func (c *Client) PurgeDeletedSecret(ctx context.Context, secretName string, options *PurgeDeletedSecretOptions) (PurgeDeletedSecretResponse, error) {
	if options == nil {
		options = &PurgeDeletedSecretOptions{}
	}
	resp, err := c.kvClient.PurgeDeletedSecret(ctx, c.vaultUrl, secretName, options.toGenerated())
	return purgeDeletedSecretResponseFromGenerated(resp), err
}

type RecoverDeletedSecretPoller struct {
	secretName      string
	vaultUrl        string
	client          *internal.KeyVaultClient
	recoverResponse internal.KeyVaultClientRecoverDeletedSecretResponse
	lastResponse    internal.KeyVaultClientGetSecretResponse
	RawResponse     *http.Response
}

// Done returns true when the polling operation is completed
func (b *RecoverDeletedSecretPoller) Done() bool {
	return b.RawResponse.StatusCode == http.StatusOK
}

// Poll fetches the latest state of the LRO. It returns an HTTP response or error.
// If the LRO has completed successfully, the poller's state is updated and the HTTP response is returned.
// If the LRO has completed with failure or was cancelled, the poller's state is updated and the error is returned.
func (b *RecoverDeletedSecretPoller) Poll(ctx context.Context) (*http.Response, error) {
	resp, err := b.client.GetSecret(ctx, b.vaultUrl, b.secretName, "", nil)
	b.lastResponse = resp
	var httpErr *azcore.ResponseError
	if errors.As(err, &httpErr) {
		return httpErr.RawResponse, err
	}
	return resp.RawResponse, nil
}

// FinalResponse returns the final response after the operations has finished
func (b *RecoverDeletedSecretPoller) FinalResponse(ctx context.Context) (RecoverDeletedSecretResponse, error) {
	return recoverDeletedSecretResponseFromGenerated(b.recoverResponse), nil
}

func (b *RecoverDeletedSecretPoller) pollUntilDone(ctx context.Context, t time.Duration) (RecoverDeletedSecretResponse, error) {
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
	return recoverDeletedSecretResponseFromGenerated(b.recoverResponse), nil
}

// BeginRecoverDeletedSecretOptions contains the optional parameters for the Client.BeginRecoverDeletedSecret operation
type BeginRecoverDeletedSecretOptions struct {
	// placeholder for future optional parameters
}

// Convert the publicly exposed options object to the generated version
func (b BeginRecoverDeletedSecretOptions) toGenerated() *internal.KeyVaultClientRecoverDeletedSecretOptions {
	return &internal.KeyVaultClientRecoverDeletedSecretOptions{}
}

// RecoverDeletedSecretResponse is the response object for the Client.RecoverDeletedSecret operation.
type RecoverDeletedSecretResponse struct {
	Secret
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// change recover deleted secret reponse to the generated version.
func recoverDeletedSecretResponseFromGenerated(i internal.KeyVaultClientRecoverDeletedSecretResponse) RecoverDeletedSecretResponse {
	var a *Properties
	if i.Attributes != nil {
		a = &Properties{
			Enabled:         i.Attributes.Enabled,
			ExpiresOn:       i.Attributes.Expires,
			NotBefore:       i.Attributes.NotBefore,
			CreatedOn:       i.Attributes.Created,
			UpdatedOn:       i.Attributes.Updated,
			RecoverableDays: i.Attributes.RecoverableDays,
			RecoveryLevel:   deletionRecoveryLevelFromGenerated(*i.Attributes.RecoveryLevel).ToPtr(),
		}
	}
	return RecoverDeletedSecretResponse{
		RawResponse: i.RawResponse,
		Secret: Secret{
			Properties:  a,
			ContentType: i.ContentType,
			ID:          i.ID,
			Tags:        convertPtrMap(i.Tags),
			Value:       i.Value,
			KeyID:       i.Kid,
			Managed:     i.Managed,
		},
	}
}

// RecoverDeletedSecretPollerResponse contains the response of the Client.BeginRecoverDeletedSecret operations
type RecoverDeletedSecretPollerResponse struct {
	// PollUntilDone will poll the service endpoint until a terminal state is reached or an error occurs
	PollUntilDone func(context.Context, time.Duration) (RecoverDeletedSecretResponse, error)

	// Poller contains an initialized RecoverDeletedSecretPoller
	Poller RecoverDeletedSecretPoller

	// RawResponse cotains the underlying HTTP response
	RawResponse *http.Response
}

// BeginRecoverDeletedSecret recovers the deleted secret in the specified vault to the latest version.
// This operation can only be performed on a soft-delete enabled vault. This operation requires the secrets/recover permission.
func (c *Client) BeginRecoverDeletedSecret(ctx context.Context, secretName string, options *BeginRecoverDeletedSecretOptions) (RecoverDeletedSecretPollerResponse, error) {
	if options == nil {
		options = &BeginRecoverDeletedSecretOptions{}
	}
	resp, err := c.kvClient.RecoverDeletedSecret(ctx, c.vaultUrl, secretName, options.toGenerated())
	if err != nil {
		return RecoverDeletedSecretPollerResponse{}, err
	}

	getResp, err := c.kvClient.GetSecret(ctx, c.vaultUrl, secretName, "", nil)
	var httpErr *azcore.ResponseError
	if errors.As(err, &httpErr) {
		if httpErr.StatusCode != http.StatusNotFound {
			return RecoverDeletedSecretPollerResponse{}, err
		}
	}

	b := RecoverDeletedSecretPoller{
		lastResponse:    getResp,
		secretName:      secretName,
		client:          c.kvClient,
		vaultUrl:        c.vaultUrl,
		recoverResponse: resp,
		RawResponse:     getResp.RawResponse,
	}

	return RecoverDeletedSecretPollerResponse{
		PollUntilDone: b.pollUntilDone,
		Poller:        b,
		RawResponse:   getResp.RawResponse,
	}, nil
}

// ListDeletedSecretsPager is the pager returned by Client.ListDeletedSecrets
type ListDeletedSecretsPager struct {
	genPager *internal.KeyVaultClientGetDeletedSecretsPager
}

// PageResponse returns the current page of results
func (l *ListDeletedSecretsPager) PageResponse() ListDeletedSecretsPage {
	resp := l.genPager.PageResponse()

	var values []DeletedSecretItem
	for _, d := range resp.Value {
		values = append(values, deletedSecretItemFromGenerated(d))
	}

	return ListDeletedSecretsPage{
		RawResponse:    resp.RawResponse,
		NextLink:       resp.NextLink,
		DeletedSecrets: values,
	}
}

// Err returns an error if the last operation resulted in an error.
func (l *ListDeletedSecretsPager) Err() error {
	return l.genPager.Err()
}

// NextPage fetches the next page of results.
func (l *ListDeletedSecretsPager) NextPage(ctx context.Context) bool {
	return l.genPager.NextPage(ctx)
}

// ListDeletedSecretsPage holds the data for a single page.
type ListDeletedSecretsPage struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response

	// READ-ONLY; The URL to get the next set of deleted secrets.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// READ-ONLY; A response message containing a list of the deleted secrets in the vault along with a link to the next page of deleted secrets
	DeletedSecrets []DeletedSecretItem `json:"value,omitempty" azure:"ro"`
}

// ListDeletedSecretsOptions contains the optional parameters for the Client.ListDeletedSecrets operation.
type ListDeletedSecretsOptions struct {
	// Maximum number of results to return in a page. If not specified the service will return up to 25 results.
	MaxResults *int32
}

// Convert publicly exposed options to the generated version.a
func (l *ListDeletedSecretsOptions) toGenerated() *internal.KeyVaultClientGetDeletedSecretsOptions {
	return &internal.KeyVaultClientGetDeletedSecretsOptions{
		Maxresults: l.MaxResults,
	}
}

// ListDeletedSecrets lists all versions of the specified secret. The full secret identifier and attributes are provided
// in the response. No values are returned for the secrets. This operation requires the secrets/list permission.
func (c *Client) ListDeletedSecrets(options *ListDeletedSecretsOptions) ListDeletedSecretsPager {
	if options == nil {
		options = &ListDeletedSecretsOptions{}
	}

	return ListDeletedSecretsPager{
		genPager: c.kvClient.GetDeletedSecrets(c.vaultUrl, options.toGenerated()),
	}

}

type ListSecretVersionsPager struct {
	genPager *internal.KeyVaultClientGetSecretVersionsPager
}

// PageResponse returns the results from the page most recently fetched from the service.
func (l *ListSecretVersionsPager) PageResponse() ListSecretVersionsPage {
	return listSecretVersionsPageFromGenerated(l.genPager.PageResponse())
}

// Err returns an error value if the most recent call to NextPage was not successful, else nil.
func (l *ListSecretVersionsPager) Err() error {
	return l.genPager.Err()
}

// NextPage fetches the next available page of results from the service. If the fetched page
// contains results, the return value is true, else false. Results fetched from the service
// can be evaluated by calling PageResponse on this Pager.
func (l *ListSecretVersionsPager) NextPage(ctx context.Context) bool {
	return l.genPager.NextPage(ctx)
}

// ListSecretVersionsOptions contains the options for the ListSecretVersions operations
type ListSecretVersionsOptions struct {
	// Maximum number of results to return in a page. If not specified the service will return up to 25 results.
	MaxResults *int32
}

// convert the public ListSecretVersionsOptions to the generated version
func (l *ListSecretVersionsOptions) toGenerated() *internal.KeyVaultClientGetSecretVersionsOptions {
	if l == nil {
		return &internal.KeyVaultClientGetSecretVersionsOptions{}
	}
	return &internal.KeyVaultClientGetSecretVersionsOptions{
		Maxresults: l.MaxResults,
	}
}

// The secret list result
type ListSecretVersionsPage struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response

	// READ-ONLY; The URL to get the next set of secrets.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// READ-ONLY; A response message containing a list of secrets in the key vault along with a link to the next page of secrets.
	Secrets []SecretItem `json:"value,omitempty" azure:"ro"`
}

// create ListSecretsPage from generated pager
func listSecretVersionsPageFromGenerated(i internal.KeyVaultClientGetSecretVersionsResponse) ListSecretVersionsPage {
	var secrets []SecretItem
	for _, s := range i.Value {
		secrets = append(secrets, secretItemFromGenerated(s))
	}
	return ListSecretVersionsPage{
		RawResponse: i.RawResponse,
		NextLink:    i.NextLink,
		Secrets:     secrets,
	}
}

// ListSecretVersions lists all versions of the specified secret. The full secret identifer and
// attributes are provided in the response. No values are returned for the secrets. This operation
// requires the secrets/list permission.
func (c *Client) ListSecretVersions(secretName string, options *ListSecretVersionsOptions) ListSecretVersionsPager {
	if options == nil {
		options = &ListSecretVersionsOptions{}
	}

	return ListSecretVersionsPager{
		genPager: c.kvClient.GetSecretVersions(
			c.vaultUrl,
			secretName,
			options.toGenerated(),
		),
	}
}

// ListSecretsPager implements the ListSecretsPager interface
type ListSecretsPager struct {
	genPager *internal.KeyVaultClientGetSecretsPager
}

// PageResponse returns the results from the page most recently fetched from the service
func (l *ListSecretsPager) PageResponse() ListSecretsPage {
	return listSecretsPageFromGenerated(l.genPager.PageResponse())
}

// Err returns an error value if the most recent call to NextPage was not successful, else nil.
func (l *ListSecretsPager) Err() error {
	return l.genPager.Err()
}

// NextPage fetches the next available page of results from the service. If the fetched page
// contains results, the return value is true, else false. Results fetched from the service
// can be evaluated by calling PageResponse on this Pager.
func (l *ListSecretsPager) NextPage(ctx context.Context) bool {
	return l.genPager.NextPage(ctx)
}

// ListSecretsOptions contains the options for the ListSecretVersions operations
type ListSecretsOptions struct {
	// Maximum number of results to return in a page. If not specified the service will return up to 25 results.
	MaxResults *int32
}

// converts the public struct to the generated code version.
func (l *ListSecretsOptions) toGenerated() *internal.KeyVaultClientGetSecretsOptions {
	if l == nil {
		return nil
	}
	return &internal.KeyVaultClientGetSecretsOptions{
		Maxresults: l.MaxResults,
	}
}

// ListSecretsPage contains the current page of results for the Client.ListSecrets operation.
type ListSecretsPage struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response

	// READ-ONLY; The URL to get the next set of secrets.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// READ-ONLY; A response message containing a list of secrets in the key vault along with a link to the next page of secrets.
	Secrets []SecretItem `json:"value,omitempty" azure:"ro"`
}

// create a ListSecretsPage from a generated code response
func listSecretsPageFromGenerated(i internal.KeyVaultClientGetSecretsResponse) ListSecretsPage {
	var secrets []SecretItem
	for _, s := range i.Value {
		secrets = append(secrets, secretItemFromGenerated(s))
	}
	return ListSecretsPage{
		RawResponse: i.RawResponse,
		NextLink:    i.NextLink,
		Secrets:     secrets,
	}
}

// ListSecrets list all secrets in a specified key vault. The ListSecrets operation is applicable to the entire vault,
// however, only the base secret identifier and its attributes are provided in the response. Individual
// secret versions are not listed in the response. This operation requires the secrets/list permission.
func (c *Client) ListSecrets(options *ListSecretsOptions) ListSecretsPager {
	if options == nil {
		options = &ListSecretsOptions{}
	}

	return ListSecretsPager{
		genPager: c.kvClient.GetSecrets(c.vaultUrl, options.toGenerated()),
	}
}
