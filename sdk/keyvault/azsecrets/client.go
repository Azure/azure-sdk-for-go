//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsecrets

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets/internal"
	shared "github.com/Azure/azure-sdk-for-go/sdk/keyvault/internal"
)

// Client interacts with Key Vault secrets.
type Client struct {
	kvClient *internal.KeyVaultClient
	vaultUrl string
}

// ClientOptions are the configurable options for a Client.
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

// NewClient constructs a Client that accesses a Key Vault's secrets.
func NewClient(vaultURL string, credential azcore.TokenCredential, options *ClientOptions) (*Client, error) {
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
func (g *GetSecretOptions) toGenerated() *internal.KeyVaultClientGetSecretOptions {
	if g == nil {
		return &internal.KeyVaultClientGetSecretOptions{}
	}
	return &internal.KeyVaultClientGetSecretOptions{}
}

// GetSecretResponse is returned by GetSecret.
type GetSecretResponse struct {
	Secret
}

func getSecretResponseFromGenerated(i internal.KeyVaultClientGetSecretResponse) GetSecretResponse {
	vaultURL, name, version := shared.ParseID(i.ID)
	return GetSecretResponse{
		Secret: Secret{
			Properties: &Properties{
				ContentType:     i.ContentType,
				CreatedOn:       i.Attributes.Created,
				Enabled:         i.Attributes.Enabled,
				ExpiresOn:       i.Attributes.Expires,
				IsManaged:       i.Managed,
				KeyID:           i.Kid,
				NotBefore:       i.Attributes.NotBefore,
				RecoverableDays: i.Attributes.RecoverableDays,
				RecoveryLevel:   (*string)(i.Attributes.RecoveryLevel),
				Tags:            convertPtrMap(i.Tags),
				UpdatedOn:       i.Attributes.Updated,
				VaultURL:        vaultURL,
				Version:         version,
				Name:            name,
			},
			ID:    i.ID,
			Name:  name,
			Value: i.Value,
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

// SetSecretResponse is returned by SetSecret.
type SetSecretResponse struct {
	Secret
}

// convert generated response to publicly exposed response.
func setSecretResponseFromGenerated(i internal.KeyVaultClientSetSecretResponse) SetSecretResponse {
	vaultURL, name, version := shared.ParseID(i.ID)
	return SetSecretResponse{
		Secret: Secret{
			Properties: &Properties{
				ContentType:     i.ContentType,
				CreatedOn:       i.Attributes.Created,
				Enabled:         i.Attributes.Enabled,
				ExpiresOn:       i.Attributes.Expires,
				IsManaged:       i.Managed,
				KeyID:           i.Kid,
				NotBefore:       i.Attributes.NotBefore,
				RecoverableDays: i.Attributes.RecoverableDays,
				RecoveryLevel:   (*string)(i.Attributes.RecoveryLevel),
				Tags:            convertPtrMap(i.Tags),
				UpdatedOn:       i.Attributes.Updated,
				VaultURL:        vaultURL,
				Version:         version,
				Name:            name,
			},
			ID:    i.ID,
			Name:  name,
			Value: i.Value,
		},
	}
}

// SetSecret sets the value of a secret. If the secret already exists, this will create a new version of the secret.
func (c *Client) SetSecret(ctx context.Context, name string, value string, options *SetSecretOptions) (SetSecretResponse, error) {
	if options == nil {
		options = &SetSecretOptions{}
	}
	var secretAttribs internal.SecretAttributes
	if options.Properties != nil {
		secretAttribs = *options.Properties.toGenerated()
	}
	resp, err := c.kvClient.SetSecret(ctx, c.vaultUrl, name, internal.SecretSetParameters{
		Value:            &value,
		ContentType:      options.ContentType,
		SecretAttributes: &secretAttribs,
		Tags:             convertToGeneratedMap(options.Tags),
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

func deleteSecretResponseFromGenerated(i internal.KeyVaultClientDeleteSecretResponse) DeleteSecretResponse {
	vaultURL, name, version := shared.ParseID(i.ID)
	return DeleteSecretResponse{
		DeletedSecret: DeletedSecret{
			ID:   i.ID,
			Name: name,
			Properties: &Properties{
				ContentType:     i.ContentType,
				CreatedOn:       i.Attributes.Created,
				Enabled:         i.Attributes.Enabled,
				ExpiresOn:       i.Attributes.Expires,
				IsManaged:       i.Managed,
				KeyID:           i.Kid,
				NotBefore:       i.Attributes.NotBefore,
				RecoverableDays: i.Attributes.RecoverableDays,
				RecoveryLevel:   (*string)(i.Attributes.RecoveryLevel),
				Tags:            convertPtrMap(i.Tags),
				UpdatedOn:       i.Attributes.Updated,
				VaultURL:        vaultURL,
				Version:         version,
				Name:            name,
			},
			RecoveryID:         i.RecoveryID,
			DeletedOn:          i.DeletedDate,
			ScheduledPurgeDate: i.ScheduledPurgeDate,
		},
	}
}

// BeginDeleteSecretOptions contains optional parameters for BeginDeleteSecret.
type BeginDeleteSecretOptions struct {
	// ResumeToken is a string to rehydrate a poller for an operation that has already begun.
	ResumeToken *string
}

// convert public options to generated options struct
func (b *BeginDeleteSecretOptions) toGenerated() *internal.KeyVaultClientDeleteSecretOptions {
	return &internal.KeyVaultClientDeleteSecretOptions{}
}

// DeleteSecretPoller is returned by BeginDeleteSecret.
type DeleteSecretPoller struct {
	secretName     string // This is the secret to Poll for in GetDeletedSecret
	vaultUrl       string
	client         *internal.KeyVaultClient
	deleteResponse internal.KeyVaultClientDeleteSecretResponse
	lastResponse   internal.KeyVaultClientGetDeletedSecretResponse
	rawResponse    *http.Response
	resumeToken    string
}

// Done returns true if the LRO has reached a terminal state
func (s *DeleteSecretPoller) Done() bool {
	if s.rawResponse == nil {
		return false
	}
	return s.rawResponse.StatusCode == http.StatusOK
}

// Poll fetches the latest state of the delete operation.
func (s *DeleteSecretPoller) Poll(ctx context.Context) (*http.Response, error) {
	var rawResp *http.Response
	ctx = runtime.WithCaptureResponse(ctx, &rawResp)
	resp, err := s.client.GetDeletedSecret(ctx, s.vaultUrl, s.secretName, nil)
	if err == nil {
		// Service recognizes DeletedSecret, operation is done
		s.lastResponse = resp
		s.rawResponse = rawResp
		return rawResp, nil
	}
	if rawResp != nil && rawResp.StatusCode == http.StatusNotFound {
		// This is the expected result
		s.rawResponse = rawResp
		return rawResp, nil
	}
	return rawResp, err
}

// FinalResponse returns the final response after the secret is deleted.
func (s *DeleteSecretPoller) FinalResponse(ctx context.Context) (DeleteSecretResponse, error) {
	return deleteSecretResponseFromGenerated(s.deleteResponse), nil
}

// PollUntilDone polls Key Vault until the secret is deleted. The t parameter determines the wait between polls.
func (s *DeleteSecretPoller) PollUntilDone(ctx context.Context, t time.Duration) (DeleteSecretResponse, error) {
	for {
		resp, err := s.Poll(ctx)
		if err != nil {
			return DeleteSecretResponse{}, err
		}
		s.rawResponse = resp
		if s.Done() {
			break
		}
		time.Sleep(t)
	}
	return deleteSecretResponseFromGenerated(s.deleteResponse), nil
}

// ResumeToken returns a token for resuming polling at a later time
func (s *DeleteSecretPoller) ResumeToken() (string, error) {
	return s.resumeToken, nil
}

// BeginDeleteSecret deletes all versions of a secret. It returns a Poller that enables waiting for Key Vault to finish
// deleting the secret.
func (c *Client) BeginDeleteSecret(ctx context.Context, name string, options *BeginDeleteSecretOptions) (*DeleteSecretPoller, error) {
	if options == nil {
		options = &BeginDeleteSecretOptions{}
	}
	var resumeToken string
	var delResp internal.KeyVaultClientDeleteSecretResponse
	var err error
	if options.ResumeToken == nil {
		delResp, err = c.kvClient.DeleteSecret(ctx, c.vaultUrl, name, options.toGenerated())
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

	getResp, err := c.kvClient.GetDeletedSecret(ctx, c.vaultUrl, name, nil)
	var httpErr *azcore.ResponseError
	if errors.As(err, &httpErr) {
		if httpErr.StatusCode != http.StatusNotFound {
			return nil, err
		}
	}

	return &DeleteSecretPoller{
		vaultUrl:       c.vaultUrl,
		secretName:     name,
		client:         c.kvClient,
		deleteResponse: delResp,
		lastResponse:   getResp,
		resumeToken:    resumeToken,
	}, nil
}

// GetDeletedSecretOptions contains optional parameters for GetDeletedSecret.
type GetDeletedSecretOptions struct {
	// placeholder for future optional parameters
}

func (g *GetDeletedSecretOptions) toGenerated() *internal.KeyVaultClientGetDeletedSecretOptions {
	return &internal.KeyVaultClientGetDeletedSecretOptions{}
}

// GetDeletedSecretResponse is returned by GetDeletedSecret.
type GetDeletedSecretResponse struct {
	DeletedSecret
}

// Convert the generated response to the publicly exposed version
func getDeletedSecretResponseFromGenerated(i internal.KeyVaultClientGetDeletedSecretResponse) GetDeletedSecretResponse {
	vaultURL, name, version := shared.ParseID(i.ID)
	return GetDeletedSecretResponse{
		DeletedSecret: DeletedSecret{
			Properties: &Properties{
				ContentType:     i.ContentType,
				CreatedOn:       i.Attributes.Created,
				Enabled:         i.Attributes.Enabled,
				ExpiresOn:       i.Attributes.Expires,
				IsManaged:       i.Managed,
				KeyID:           i.Kid,
				NotBefore:       i.Attributes.NotBefore,
				RecoverableDays: i.Attributes.RecoverableDays,
				RecoveryLevel:   (*string)(i.Attributes.RecoveryLevel),
				Tags:            convertPtrMap(i.Tags),
				UpdatedOn:       i.Attributes.Updated,
				VaultURL:        vaultURL,
				Version:         version,
				Name:            name,
			},
			ID:                 i.ID,
			Name:               name,
			RecoveryID:         i.RecoveryID,
			DeletedOn:          i.DeletedDate,
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

func updateSecretPropertiesResponseFromGenerated(i internal.KeyVaultClientUpdateSecretResponse) UpdateSecretPropertiesResponse {
	vaultURL, name, version := shared.ParseID(i.ID)
	return UpdateSecretPropertiesResponse{
		Secret: Secret{
			Properties: &Properties{
				ContentType:     i.ContentType,
				CreatedOn:       i.Attributes.Created,
				Enabled:         i.Attributes.Enabled,
				ExpiresOn:       i.Attributes.Expires,
				IsManaged:       i.Managed,
				KeyID:           i.Kid,
				NotBefore:       i.Attributes.NotBefore,
				RecoverableDays: i.Attributes.RecoverableDays,
				RecoveryLevel:   (*string)(i.Attributes.RecoveryLevel),
				Tags:            convertPtrMap(i.Tags),
				UpdatedOn:       i.Attributes.Updated,
				VaultURL:        vaultURL,
				Version:         version,
				Name:            name,
			},
			ID:    i.ID,
			Name:  name,
			Value: i.Value,
		},
	}
}

// UpdateSecretProperties updates a secret's properties, such as whether it's enabled. See the Properties type for a complete list.
// nil fields will keep their current values. This method can't change the secret's value; use SetSecret to do that.
func (c *Client) UpdateSecretProperties(ctx context.Context, secret Secret, options *UpdateSecretPropertiesOptions) (UpdateSecretPropertiesResponse, error) {
	name, version := "", ""
	if secret.Properties != nil && secret.Properties.Name != nil {
		name = *secret.Properties.Name
	}
	if secret.Properties != nil && secret.Properties.Version != nil {
		version = *secret.Properties.Version
	}

	resp, err := c.kvClient.UpdateSecret(
		ctx,
		c.vaultUrl,
		name,
		version,
		secret.toGeneratedProperties(),
		&internal.KeyVaultClientUpdateSecretOptions{},
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

func (b *BackupSecretOptions) toGenerated() *internal.KeyVaultClientBackupSecretOptions {
	return &internal.KeyVaultClientBackupSecretOptions{}
}

// BackupSecretResponse is returned by BackupSecret.
type BackupSecretResponse struct {
	// READ-ONLY; The backup blob containing the backed up secret.
	Value []byte `json:"value,omitempty" azure:"ro"`
}

// convert generated response to the publicly exposed version.
func backupSecretResponseFromGenerated(i internal.KeyVaultClientBackupSecretResponse) BackupSecretResponse {
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

func (r RestoreSecretBackupOptions) toGenerated() *internal.KeyVaultClientRestoreSecretOptions {
	return &internal.KeyVaultClientRestoreSecretOptions{}
}

// RestoreSecretBackupResponse is returned by RestoreSecretBackup.
type RestoreSecretBackupResponse struct {
	Secret
}

// converts the generated response to the publicly exposed version.
func restoreSecretBackupResponseFromGenerated(i internal.KeyVaultClientRestoreSecretResponse) RestoreSecretBackupResponse {
	vaultURL, name, version := shared.ParseID(i.ID)
	return RestoreSecretBackupResponse{
		Secret: Secret{
			ID:    i.ID,
			Name:  name,
			Value: i.Value,
			Properties: &Properties{
				ContentType:     i.ContentType,
				CreatedOn:       i.Attributes.Created,
				Enabled:         i.Attributes.Enabled,
				ExpiresOn:       i.Attributes.Expires,
				IsManaged:       i.Managed,
				KeyID:           i.Kid,
				NotBefore:       i.Attributes.NotBefore,
				RecoverableDays: i.Attributes.RecoverableDays,
				RecoveryLevel:   (*string)(i.Attributes.RecoveryLevel),
				Tags:            convertPtrMap(i.Tags),
				UpdatedOn:       i.Attributes.Updated,
				VaultURL:        vaultURL,
				Version:         version,
				Name:            name,
			},
		},
	}
}

// RestoreSecretBackup restores a secret backup, as returned by BackupSecret, to the vault. This will restore all versions of
// the secret in the backup.
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

// PurgeDeletedSecretOptions contains options for Client.PurgeDeletedSecret.
type PurgeDeletedSecretOptions struct {
	// placeholder for future optional parameters
}

func (p *PurgeDeletedSecretOptions) toGenerated() *internal.KeyVaultClientPurgeDeletedSecretOptions {
	return &internal.KeyVaultClientPurgeDeletedSecretOptions{}
}

// PurgeDeletedSecretResponse contains the response from method Client.PurgeDeletedSecret.
type PurgeDeletedSecretResponse struct {
	// placeholder for future response fields
}

// Converts the generated response to the publicly exposed version.
func purgeDeletedSecretResponseFromGenerated(i internal.KeyVaultClientPurgeDeletedSecretResponse) PurgeDeletedSecretResponse {
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

// RecoverDeletedSecretPoller is returned by BeginRecoverDeletedSecret.
type RecoverDeletedSecretPoller struct {
	secretName      string
	vaultUrl        string
	client          *internal.KeyVaultClient
	recoverResponse internal.KeyVaultClientRecoverDeletedSecretResponse
	lastResponse    internal.KeyVaultClientGetSecretResponse
	rawResponse     *http.Response
	resumeToken     string
}

// Done returns true when the polling operation is completed
func (b *RecoverDeletedSecretPoller) Done() bool {
	if b.rawResponse == nil {
		return false
	}
	return b.rawResponse.StatusCode == http.StatusOK
}

// Poll fetches the latest state of the recover operation.
func (b *RecoverDeletedSecretPoller) Poll(ctx context.Context) (*http.Response, error) {
	var rawResp *http.Response
	ctx = runtime.WithCaptureResponse(ctx, &rawResp)
	resp, err := b.client.GetSecret(ctx, b.vaultUrl, b.secretName, "", nil)
	if err == nil {
		// secret has been recovered, finish
		b.lastResponse = resp
		b.rawResponse = rawResp
		return b.rawResponse, nil
	}

	if rawResp != nil && rawResp.StatusCode == http.StatusNotFound {
		// this is the expected response
		b.lastResponse = resp
		b.rawResponse = rawResp
		return b.rawResponse, nil
	}

	return rawResp, err
}

// FinalResponse returns the final response after the recover operation is complete.
func (b *RecoverDeletedSecretPoller) FinalResponse(ctx context.Context) (RecoverDeletedSecretResponse, error) {
	return recoverDeletedSecretResponseFromGenerated(b.recoverResponse), nil
}

// PollUntilDone polls Key Vault until the recover operation is complete. The t parameter determines the wait between polls.
func (b *RecoverDeletedSecretPoller) PollUntilDone(ctx context.Context, t time.Duration) (RecoverDeletedSecretResponse, error) {
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
	return recoverDeletedSecretResponseFromGenerated(b.recoverResponse), nil
}

// ResumeToken returns a token for resuming polling at a later time
func (s *RecoverDeletedSecretPoller) ResumeToken() (string, error) {
	return s.resumeToken, nil
}

// BeginRecoverDeletedSecretOptions contains optional parameters for BeginRecoverDeletedSecret.
type BeginRecoverDeletedSecretOptions struct {
	// ResumeToken is a string to rehydrate a poller for an operation that has already begun.
	ResumeToken *string
}

// Convert the publicly exposed options object to the generated version
func (b BeginRecoverDeletedSecretOptions) toGenerated() *internal.KeyVaultClientRecoverDeletedSecretOptions {
	return &internal.KeyVaultClientRecoverDeletedSecretOptions{}
}

// RecoverDeletedSecretResponse is returned by RecoverDeletedSecret.
type RecoverDeletedSecretResponse struct {
	SecretItem
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
			RecoveryLevel:   (*string)(i.Attributes.RecoveryLevel),
		}
	}

	_, name, _ := shared.ParseID(i.ID)
	return RecoverDeletedSecretResponse{
		SecretItem: SecretItem{
			Properties:  a,
			ContentType: i.ContentType,
			ID:          i.ID,
			Name:        name,
			Tags:        convertPtrMap(i.Tags),
			IsManaged:   i.Managed,
		},
	}
}

// BeginRecoverDeletedSecret recovers a deleted secret to its latest version. Recovery may take several seconds. This method
// therefore returns a poller that enables waiting until recovery is complete.
func (c *Client) BeginRecoverDeletedSecret(ctx context.Context, name string, options *BeginRecoverDeletedSecretOptions) (*RecoverDeletedSecretPoller, error) {
	if options == nil {
		options = &BeginRecoverDeletedSecretOptions{}
	}
	var resumeToken string
	var recoverResp internal.KeyVaultClientRecoverDeletedSecretResponse
	var err error
	if options.ResumeToken == nil {
		recoverResp, err = c.kvClient.RecoverDeletedSecret(ctx, c.vaultUrl, name, options.toGenerated())
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

	getResp, err := c.kvClient.GetSecret(ctx, c.vaultUrl, name, "", nil)
	var httpErr *azcore.ResponseError
	if errors.As(err, &httpErr) {
		if httpErr.StatusCode != http.StatusNotFound {
			return nil, err
		}
	}

	return &RecoverDeletedSecretPoller{
		lastResponse:    getResp,
		secretName:      name,
		client:          c.kvClient,
		vaultUrl:        c.vaultUrl,
		recoverResponse: recoverResp,
		resumeToken:     resumeToken,
	}, nil
}

// ListDeletedSecretsResponse contains a page of deleted secrets.
type ListDeletedSecretsResponse struct {
	// NextLink is the URL to get the next page.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// DeletedSecrets is the page's content.
	DeletedSecrets []DeletedSecretItem `json:"value,omitempty" azure:"ro"`
}

func listDeletedSecretsPageFromGenerated(g internal.KeyVaultClientGetDeletedSecretsResponse) ListDeletedSecretsResponse {
	var items []DeletedSecretItem

	if len(g.DeletedSecretListResult.Value) > 0 {
		items = make([]DeletedSecretItem, len(g.DeletedSecretListResult.Value))
		for idx, v := range g.DeletedSecretListResult.Value {
			items[idx] = deletedSecretItemFromGenerated(v)
		}
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
				req, err = c.kvClient.GetDeletedSecretsCreateRequest(ctx, c.vaultUrl, &internal.KeyVaultClientGetDeletedSecretsOptions{})
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ListDeletedSecretsResponse{}, err
			}
			resp, err := c.kvClient.Pl.Do(req)
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
	Secrets []SecretItem `json:"value,omitempty" azure:"ro"`
}

// create ListSecretsPage from generated pager
func listSecretVersionsPageFromGenerated(i internal.KeyVaultClientGetSecretVersionsResponse) ListPropertiesOfSecretVersionsResponse {
	var secrets []SecretItem
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
				req, err = c.kvClient.GetSecretVersionsCreateRequest(ctx, c.vaultUrl, name, &internal.KeyVaultClientGetSecretVersionsOptions{})
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ListPropertiesOfSecretVersionsResponse{}, err
			}
			resp, err := c.kvClient.Pl.Do(req)
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
	Secrets []SecretItem `json:"value,omitempty" azure:"ro"`
}

// create a ListSecretsPage from a generated code response
func listSecretsPageFromGenerated(i internal.KeyVaultClientGetSecretsResponse) ListPropertiesOfSecretsResponse {
	var secrets []SecretItem
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
				req, err = c.kvClient.GetSecretsCreateRequest(ctx, c.vaultUrl, &internal.KeyVaultClientGetSecretsOptions{})
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ListPropertiesOfSecretsResponse{}, err
			}
			resp, err := c.kvClient.Pl.Do(req)
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
