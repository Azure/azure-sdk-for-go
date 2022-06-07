//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsecrets

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets/internal/generated"
	shared "github.com/Azure/azure-sdk-for-go/sdk/keyvault/internal"
)

// Client interacts with Key Vault secrets.
type Client struct {
	kvClient *generated.KeyVaultClient
	vaultUrl string
}

// ClientOptions are the configurable options for a Client.
type ClientOptions struct {
	azcore.ClientOptions
}

// NewClient constructs a Client that accesses a Key Vault's secrets.
func NewClient(vaultURL string, credential azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}
	plOpts := runtime.PipelineOptions{
		PerRetry: []policy.Policy{shared.NewKeyVaultChallengePolicy(credential)},
	}
	pl := runtime.NewPipeline(moduleName, version, plOpts, &options.ClientOptions)
	return &Client{
		kvClient: generated.NewKeyVaultClient(pl),
		vaultUrl: vaultURL,
	}, nil
}

// VaultURL returns the URL for the client's Key Vault.
func (c *Client) VaultURL() string {
	return c.vaultUrl
}

// GetSecretOptions contains optional parameters for GetSecret.
type GetSecretOptions struct {
	// Version specifies which version of the secret the client should get.
	// If unspecified, the client will get the most recent version.
	Version string
}

// convert the exposed options struct to the internal one.
func (g *GetSecretOptions) toGenerated() *generated.KeyVaultClientGetSecretOptions {
	if g == nil {
		return &generated.KeyVaultClientGetSecretOptions{}
	}
	return &generated.KeyVaultClientGetSecretOptions{}
}

// GetSecretResponse is returned by GetSecret.
type GetSecretResponse struct {
	Secret
}

func getSecretResponseFromGenerated(i generated.KeyVaultClientGetSecretResponse) GetSecretResponse {
	props := secretPropertiesFromGenerated(i.Attributes, i.ID, i.ContentType, i.Kid, i.Managed, i.Tags)
	return GetSecretResponse{
		Secret: Secret{
			ID:         i.ID,
			Name:       props.Name,
			Properties: props,
			Value:      i.Value,
		},
	}
}

// GetSecret gets a secret value from the vault.
func (c *Client) GetSecret(ctx context.Context, name string, options *GetSecretOptions) (GetSecretResponse, error) {
	if options == nil {
		options = &GetSecretOptions{}
	}
	resp, err := c.kvClient.GetSecret(ctx, c.vaultUrl, name, options.Version, options.toGenerated())
	if err != nil {
		return GetSecretResponse{}, err
	}
	return getSecretResponseFromGenerated(resp), nil
}

// SetSecretOptions contains optional parameters for SetSecret.
type SetSecretOptions struct {
	// Type of the secret value such as a password.
	ContentType *string `json:"contentType,omitempty"`

	// The secret management attributes.
	Properties *Properties `json:"attributes,omitempty"`
}

// Convert the exposed struct to the generated code version
func (s *SetSecretOptions) toGenerated() *generated.KeyVaultClientSetSecretOptions {
	if s == nil {
		return nil
	}
	return &generated.KeyVaultClientSetSecretOptions{}
}

// SetSecretResponse is returned by SetSecret.
type SetSecretResponse struct {
	Secret
}

// convert generated response to publicly exposed response.
func setSecretResponseFromGenerated(i generated.KeyVaultClientSetSecretResponse) SetSecretResponse {
	props := secretPropertiesFromGenerated(i.Attributes, i.ID, i.ContentType, i.Kid, i.Managed, i.Tags)
	return SetSecretResponse{
		Secret: Secret{
			ID:         i.ID,
			Name:       props.Name,
			Properties: props,
			Value:      i.Value,
		},
	}
}

// SetSecret sets the value of a secret. If the secret already exists, this will create a new version of the secret.
func (c *Client) SetSecret(ctx context.Context, name string, value string, options *SetSecretOptions) (SetSecretResponse, error) {
	if options == nil {
		options = &SetSecretOptions{}
	}
	var secretAttribs generated.SecretAttributes
	var tags map[string]*string
	if options.Properties != nil {
		secretAttribs = *options.Properties.toGenerated()
		tags = options.Properties.Tags
	}
	resp, err := c.kvClient.SetSecret(ctx, c.vaultUrl, name, generated.SecretSetParameters{
		Value:            &value,
		ContentType:      options.ContentType,
		SecretAttributes: &secretAttribs,
		Tags:             tags,
	}, options.toGenerated())
	if err != nil {
		return SetSecretResponse{}, err
	}
	return setSecretResponseFromGenerated(resp), nil
}

// DeleteSecretResponse is returned by DeleteSecret.
type DeleteSecretResponse struct {
	DeletedSecret
}

func deleteSecretResponseFromGenerated(i generated.KeyVaultClientDeleteSecretResponse) DeleteSecretResponse {
	props := secretPropertiesFromGenerated(i.Attributes, i.ID, i.ContentType, i.Kid, i.Managed, i.Tags)
	return DeleteSecretResponse{
		DeletedSecret: DeletedSecret{
			ID:                 i.ID,
			Name:               props.Name,
			Properties:         props,
			RecoveryID:         i.RecoveryID,
			DeletedOn:          i.DeletedDate,
			ScheduledPurgeDate: i.ScheduledPurgeDate,
		},
	}
}

// BeginDeleteSecretOptions contains optional parameters for BeginDeleteSecret.
type BeginDeleteSecretOptions struct {
	// ResumeToken is a string to rehydrate a poller for an operation that has already begun.
	ResumeToken string
}

// convert public options to generated options struct
func (b *BeginDeleteSecretOptions) toGenerated() *generated.KeyVaultClientDeleteSecretOptions {
	return &generated.KeyVaultClientDeleteSecretOptions{}
}

// BeginDeleteSecret deletes all versions of a secret. It returns a Poller that enables waiting for Key Vault to finish
// deleting the secret.
func (c *Client) BeginDeleteSecret(ctx context.Context, name string, options *BeginDeleteSecretOptions) (*runtime.Poller[DeleteSecretResponse], error) {
	if options == nil {
		options = &BeginDeleteSecretOptions{}
	}

	handler := beginDeleteSecretOperation{
		poll: func(ctx context.Context) (*http.Response, error) {
			req, err := c.kvClient.GetDeletedSecretCreateRequest(ctx, c.vaultUrl, name, nil)
			if err != nil {
				return nil, err
			}
			return c.kvClient.Pipeline().Do(req)
		},
	}

	if options.ResumeToken != "" {
		return runtime.NewPollerFromResumeToken(
			options.ResumeToken, c.kvClient.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[DeleteSecretResponse]{Handler: &handler},
		)
	}

	var rawResp *http.Response
	ctx = runtime.WithCaptureResponse(ctx, &rawResp)
	if _, err := c.kvClient.DeleteSecret(ctx, c.vaultUrl, name, options.toGenerated()); err != nil {
		return nil, err
	}

	return runtime.NewPoller(rawResp, c.kvClient.Pipeline(), &runtime.NewPollerOptions[DeleteSecretResponse]{Handler: &handler})
}

// GetDeletedSecretOptions contains optional parameters for GetDeletedSecret.
type GetDeletedSecretOptions struct {
	// placeholder for future optional parameters
}

func (g *GetDeletedSecretOptions) toGenerated() *generated.KeyVaultClientGetDeletedSecretOptions {
	return &generated.KeyVaultClientGetDeletedSecretOptions{}
}

// GetDeletedSecretResponse is returned by GetDeletedSecret.
type GetDeletedSecretResponse struct {
	DeletedSecret
}

// Convert the generated response to the publicly exposed version
func getDeletedSecretResponseFromGenerated(i generated.KeyVaultClientGetDeletedSecretResponse) GetDeletedSecretResponse {
	props := secretPropertiesFromGenerated(i.Attributes, i.ID, i.ContentType, i.Kid, i.Managed, i.Tags)
	return GetDeletedSecretResponse{
		DeletedSecret: DeletedSecret{
			DeletedOn:          i.DeletedDate,
			ID:                 i.ID,
			Name:               props.Name,
			Properties:         props,
			RecoveryID:         i.RecoveryID,
			ScheduledPurgeDate: i.ScheduledPurgeDate,
		},
	}
}

// GetDeletedSecret gets a deleted secret.
func (c *Client) GetDeletedSecret(ctx context.Context, name string, options *GetDeletedSecretOptions) (GetDeletedSecretResponse, error) {
	if options == nil {
		options = &GetDeletedSecretOptions{}
	}
	resp, err := c.kvClient.GetDeletedSecret(ctx, c.vaultUrl, name, options.toGenerated())
	if err != nil {
		return GetDeletedSecretResponse{}, err
	}
	return getDeletedSecretResponseFromGenerated(resp), nil
}

// UpdateSecretPropertiesOptions contains optional parameters for UpdateSecretProperties.
type UpdateSecretPropertiesOptions struct {
	// placeholder for future optional parameters
}

// UpdateSecretPropertiesResponse is returned by UpdateSecretProperties.
type UpdateSecretPropertiesResponse struct {
	Secret
}

func updateSecretPropertiesResponseFromGenerated(i generated.KeyVaultClientUpdateSecretResponse) UpdateSecretPropertiesResponse {
	props := secretPropertiesFromGenerated(i.Attributes, i.ID, i.ContentType, i.Kid, i.Managed, i.Tags)
	return UpdateSecretPropertiesResponse{
		Secret: Secret{
			ID:         i.ID,
			Name:       props.Name,
			Properties: props,
			Value:      i.Value,
		},
	}
}

// UpdateSecretProperties updates a secret's properties, such as whether it's enabled. See the Properties type for a complete list.
// nil fields will keep their current values. This method can't change the secret's value; use SetSecret to do that.
func (c *Client) UpdateSecretProperties(ctx context.Context, properties Properties, options *UpdateSecretPropertiesOptions) (UpdateSecretPropertiesResponse, error) {
	name, version := "", ""
	if properties.Name != nil {
		name = *properties.Name
	}
	if properties.Version != nil {
		version = *properties.Version
	}
	resp, err := c.kvClient.UpdateSecret(
		ctx,
		c.vaultUrl,
		name,
		version,
		generated.SecretUpdateParameters{
			ContentType:      properties.ContentType,
			SecretAttributes: properties.toGenerated(),
			Tags:             properties.Tags,
		},
		nil,
	)
	if err != nil {
		return UpdateSecretPropertiesResponse{}, err
	}

	return updateSecretPropertiesResponseFromGenerated(resp), err
}

// BackupSecretOptions contains optional parameters for BackupSecret.
type BackupSecretOptions struct {
	// placeholder for future optional parameters
}

func (b *BackupSecretOptions) toGenerated() *generated.KeyVaultClientBackupSecretOptions {
	return &generated.KeyVaultClientBackupSecretOptions{}
}

// BackupSecretResponse is returned by BackupSecret.
type BackupSecretResponse struct {
	// READ-ONLY; The backup blob containing the backed up secret.
	Value []byte `json:"value,omitempty" azure:"ro"`
}

// convert generated response to the publicly exposed version.
func backupSecretResponseFromGenerated(i generated.KeyVaultClientBackupSecretResponse) BackupSecretResponse {
	return BackupSecretResponse{
		Value: i.Value,
	}
}

// BackupSecret requests an encrypted backup of all versions of a secret, readable only by Key Vault. Call RestoreSecret to restore a backup.
func (c *Client) BackupSecret(ctx context.Context, name string, options *BackupSecretOptions) (BackupSecretResponse, error) {
	if options == nil {
		options = &BackupSecretOptions{}
	}

	resp, err := c.kvClient.BackupSecret(ctx, c.vaultUrl, name, options.toGenerated())
	if err != nil {
		return BackupSecretResponse{}, err
	}

	return backupSecretResponseFromGenerated(resp), nil
}

// RestoreSecretBackupOptions contains optional parameters for RestoreSecret.
type RestoreSecretBackupOptions struct {
	// placeholder for future optional parameters
}

func (r RestoreSecretBackupOptions) toGenerated() *generated.KeyVaultClientRestoreSecretOptions {
	return &generated.KeyVaultClientRestoreSecretOptions{}
}

// RestoreSecretBackupResponse is returned by RestoreSecretBackup.
type RestoreSecretBackupResponse struct {
	Secret
}

// converts the generated response to the publicly exposed version.
func restoreSecretBackupResponseFromGenerated(i generated.KeyVaultClientRestoreSecretResponse) RestoreSecretBackupResponse {
	props := secretPropertiesFromGenerated(i.Attributes, i.ID, i.ContentType, i.Kid, i.Managed, i.Tags)
	return RestoreSecretBackupResponse{
		Secret: Secret{
			ID:         i.ID,
			Name:       props.Name,
			Properties: props,
			Value:      i.Value,
		},
	}
}

// RestoreSecretBackup restores a secret backup, as returned by BackupSecret, to the vault. This will restore all versions of
// the secret in the backup.
func (c *Client) RestoreSecretBackup(ctx context.Context, backup []byte, options *RestoreSecretBackupOptions) (RestoreSecretBackupResponse, error) {
	if options == nil {
		options = &RestoreSecretBackupOptions{}
	}

	resp, err := c.kvClient.RestoreSecret(ctx, c.vaultUrl, generated.SecretRestoreParameters{SecretBundleBackup: backup}, options.toGenerated())
	if err != nil {
		return RestoreSecretBackupResponse{}, err
	}

	return restoreSecretBackupResponseFromGenerated(resp), nil
}

// PurgeDeletedSecretOptions contains options for Client.PurgeDeletedSecret.
type PurgeDeletedSecretOptions struct {
	// placeholder for future optional parameters
}

func (p *PurgeDeletedSecretOptions) toGenerated() *generated.KeyVaultClientPurgeDeletedSecretOptions {
	return &generated.KeyVaultClientPurgeDeletedSecretOptions{}
}

// PurgeDeletedSecretResponse contains the response from method Client.PurgeDeletedSecret.
type PurgeDeletedSecretResponse struct {
	// placeholder for future response fields
}

// Converts the generated response to the publicly exposed version.
func purgeDeletedSecretResponseFromGenerated(i generated.KeyVaultClientPurgeDeletedSecretResponse) PurgeDeletedSecretResponse {
	return PurgeDeletedSecretResponse{}
}

// PurgeDeletedSecret permanently deletes a deleted secret.
func (c *Client) PurgeDeletedSecret(ctx context.Context, name string, options *PurgeDeletedSecretOptions) (PurgeDeletedSecretResponse, error) {
	if options == nil {
		options = &PurgeDeletedSecretOptions{}
	}
	resp, err := c.kvClient.PurgeDeletedSecret(ctx, c.vaultUrl, name, options.toGenerated())
	return purgeDeletedSecretResponseFromGenerated(resp), err
}

// BeginRecoverDeletedSecretOptions contains optional parameters for BeginRecoverDeletedSecret.
type BeginRecoverDeletedSecretOptions struct {
	// ResumeToken is a string to rehydrate a poller for an operation that has already begun.
	ResumeToken string
}

// Convert the publicly exposed options object to the generated version
func (b BeginRecoverDeletedSecretOptions) toGenerated() *generated.KeyVaultClientRecoverDeletedSecretOptions {
	return &generated.KeyVaultClientRecoverDeletedSecretOptions{}
}

// RecoverDeletedSecretResponse is returned by RecoverDeletedSecret.
type RecoverDeletedSecretResponse struct {
	SecretItem
}

// change recover deleted secret reponse to the generated version.
func recoverDeletedSecretResponseFromGenerated(i generated.KeyVaultClientRecoverDeletedSecretResponse) RecoverDeletedSecretResponse {
	props := secretPropertiesFromGenerated(i.Attributes, i.ID, i.ContentType, i.Kid, i.Managed, i.Tags)
	return RecoverDeletedSecretResponse{
		SecretItem: SecretItem{
			Properties: props,
			ID:         i.ID,
			Name:       props.Name,
		},
	}
}

// BeginRecoverDeletedSecret recovers a deleted secret to its latest version. Recovery may take several seconds. This method
// therefore returns a poller that enables waiting until recovery is complete.
func (c *Client) BeginRecoverDeletedSecret(ctx context.Context, name string, options *BeginRecoverDeletedSecretOptions) (*runtime.Poller[RecoverDeletedSecretResponse], error) {
	if options == nil {
		options = &BeginRecoverDeletedSecretOptions{}
	}

	handler := beginRecoverDeletedSecretOperation{
		poll: func(ctx context.Context) (*http.Response, error) {
			req, err := c.kvClient.GetSecretCreateRequest(ctx, c.vaultUrl, name, "", nil)
			if err != nil {
				return nil, err
			}
			return c.kvClient.Pipeline().Do(req)
		},
	}

	if options.ResumeToken != "" {
		return runtime.NewPollerFromResumeToken(
			options.ResumeToken, c.kvClient.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[RecoverDeletedSecretResponse]{Handler: &handler},
		)
	}

	var rawResp *http.Response
	ctx = runtime.WithCaptureResponse(ctx, &rawResp)
	if _, err := c.kvClient.RecoverDeletedSecret(ctx, c.vaultUrl, name, options.toGenerated()); err != nil {
		return nil, err
	}

	return runtime.NewPoller(
		rawResp, c.kvClient.Pipeline(), &runtime.NewPollerOptions[RecoverDeletedSecretResponse]{Handler: &handler},
	)
}

// ListDeletedSecretsResponse contains a page of deleted secrets.
type ListDeletedSecretsResponse struct {
	// NextLink is the URL to get the next page.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// DeletedSecrets is the page's content.
	DeletedSecrets []*DeletedSecretItem `json:"value,omitempty" azure:"ro"`
}

func listDeletedSecretsPageFromGenerated(g generated.KeyVaultClientGetDeletedSecretsResponse) ListDeletedSecretsResponse {
	var items []*DeletedSecretItem
	for _, v := range g.DeletedSecretListResult.Value {
		items = append(items, deletedSecretItemFromGenerated(v))
	}

	return ListDeletedSecretsResponse{
		NextLink:       g.NextLink,
		DeletedSecrets: items,
	}
}

// ListDeletedSecretsOptions contains optional parameters for ListDeletedSecrets.
type ListDeletedSecretsOptions struct {
	// placeholder for future optional parameters
}

// NewListDeletedSecretsPager creates a pager that lists all versions of a deleted secret, including their identifiers and other properties but no secret values.
func (c *Client) NewListDeletedSecretsPager(options *ListDeletedSecretsOptions) *runtime.Pager[ListDeletedSecretsResponse] {
	return runtime.NewPager(runtime.PagingHandler[ListDeletedSecretsResponse]{
		More: func(page ListDeletedSecretsResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ListDeletedSecretsResponse) (ListDeletedSecretsResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = c.kvClient.GetDeletedSecretsCreateRequest(ctx, c.vaultUrl, &generated.KeyVaultClientGetDeletedSecretsOptions{})
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ListDeletedSecretsResponse{}, err
			}
			resp, err := c.kvClient.Pipeline().Do(req)
			if err != nil {
				return ListDeletedSecretsResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ListDeletedSecretsResponse{}, runtime.NewResponseError(resp)
			}
			genResp, err := c.kvClient.GetDeletedSecretsHandleResponse(resp)
			if err != nil {
				return ListDeletedSecretsResponse{}, err
			}
			return listDeletedSecretsPageFromGenerated(genResp), nil
		},
	})
}

// ListPropertiesOfSecretVersionsOptions contains options for NewListPropertiesOfSecretVersionsPager.
type ListPropertiesOfSecretVersionsOptions struct {
	// placeholder for future optional parameters
}

// ListPropertiesOfSecretVersionsResponse contains a page of secret versions.
type ListPropertiesOfSecretVersionsResponse struct {
	// NextLink is the URL to get the next page.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// Secrets is the page's content.
	Secrets []*SecretItem `json:"value,omitempty" azure:"ro"`
}

// create ListSecretsPage from generated pager
func listSecretVersionsPageFromGenerated(i generated.KeyVaultClientGetSecretVersionsResponse) ListPropertiesOfSecretVersionsResponse {
	var secrets []*SecretItem
	for _, s := range i.Value {
		secrets = append(secrets, secretItemFromGenerated(s))
	}
	return ListPropertiesOfSecretVersionsResponse{
		NextLink: i.NextLink,
		Secrets:  secrets,
	}
}

// NewListPropertiesOfSecretVersionsPager creates a pager that lists the properties of all versions of a secret, not including their secret values.
func (c *Client) NewListPropertiesOfSecretVersionsPager(name string, options *ListPropertiesOfSecretVersionsOptions) *runtime.Pager[ListPropertiesOfSecretVersionsResponse] {
	return runtime.NewPager(runtime.PagingHandler[ListPropertiesOfSecretVersionsResponse]{
		More: func(page ListPropertiesOfSecretVersionsResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ListPropertiesOfSecretVersionsResponse) (ListPropertiesOfSecretVersionsResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = c.kvClient.GetSecretVersionsCreateRequest(ctx, c.vaultUrl, name, &generated.KeyVaultClientGetSecretVersionsOptions{})
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ListPropertiesOfSecretVersionsResponse{}, err
			}
			resp, err := c.kvClient.Pipeline().Do(req)
			if err != nil {
				return ListPropertiesOfSecretVersionsResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ListPropertiesOfSecretVersionsResponse{}, runtime.NewResponseError(resp)
			}
			genResp, err := c.kvClient.GetSecretVersionsHandleResponse(resp)
			if err != nil {
				return ListPropertiesOfSecretVersionsResponse{}, err
			}
			return listSecretVersionsPageFromGenerated(genResp), nil
		},
	})
}

// ListPropertiesOfSecretsOptions contains options for NewListPropertiesOfSecretsPager.
type ListPropertiesOfSecretsOptions struct {
	// placeholder for future optional parameters.
}

// ListPropertiesOfSecretsResponse contains a page of secret properties.
type ListPropertiesOfSecretsResponse struct {
	// NextLink is the URL to get the next page.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// Secrets is the page's content.
	Secrets []*SecretItem `json:"value,omitempty" azure:"ro"`
}

// create a ListSecretsPage from a generated code response
func listSecretsPageFromGenerated(i generated.KeyVaultClientGetSecretsResponse) ListPropertiesOfSecretsResponse {
	var secrets []*SecretItem
	for _, s := range i.Value {
		secrets = append(secrets, secretItemFromGenerated(s))
	}
	return ListPropertiesOfSecretsResponse{
		NextLink: i.NextLink,
		Secrets:  secrets,
	}
}

// NewListPropertiesOfSecretsPager constructs a pager that lists the properties of all secrets in the Key Vault, not including their secret values.
func (c *Client) NewListPropertiesOfSecretsPager(options *ListPropertiesOfSecretsOptions) *runtime.Pager[ListPropertiesOfSecretsResponse] {
	return runtime.NewPager(runtime.PagingHandler[ListPropertiesOfSecretsResponse]{
		More: func(page ListPropertiesOfSecretsResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ListPropertiesOfSecretsResponse) (ListPropertiesOfSecretsResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = c.kvClient.GetSecretsCreateRequest(ctx, c.vaultUrl, &generated.KeyVaultClientGetSecretsOptions{})
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ListPropertiesOfSecretsResponse{}, err
			}
			resp, err := c.kvClient.Pipeline().Do(req)
			if err != nil {
				return ListPropertiesOfSecretsResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ListPropertiesOfSecretsResponse{}, runtime.NewResponseError(resp)
			}
			genResp, err := c.kvClient.GetSecretsHandleResponse(resp)
			if err != nil {
				return ListPropertiesOfSecretsResponse{}, err
			}
			return listSecretsPageFromGenerated(genResp), nil
		},
	})
}
