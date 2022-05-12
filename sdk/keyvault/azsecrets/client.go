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

// VaultURL returns the vault URL string for the client
func (c *Client) VaultURL() string {
	return c.vaultUrl
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

// GetSecret gets a specified secret from a given key vault. The GET operation is applicable to any secret
// stored in Azure Key Vault. This operation requires the secrets/get permission
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

// SetSecret sets a secret in a specifed key vault. The set operation adds a secret to the Azure Key Vault, if the named secret
// already exists, Azure Key Vault creates a new version of that secret. This operation requires the secrets/set permission.
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

// DeleteSecretResponse contains the response for a Client.DeleteSecret operation.
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

// BeginDeleteSecretOptions contains the optional parameters for the Client.BeginDeleteSecret method.
type BeginDeleteSecretOptions struct {
	// ResumeToken is a string to rehydrate a poller for an operation that has already begun.
	ResumeToken *string
}

// convert public options to generated options struct
func (b *BeginDeleteSecretOptions) toGenerated() *internal.KeyVaultClientDeleteSecretOptions {
	return &internal.KeyVaultClientDeleteSecretOptions{}
}

// DeleteSecretPoller is the poller returned by the Client.StartDeleteSecret operation
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

// Poll fetches the latest state of the LRO. It returns an HTTP response or error.(
// If the LRO has completed successfully, the poller's state is updated and the HTTP response is returned.
// If the LRO has completed with failure or was cancelled, the poller's state is updated and the error is returned.)
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

// FinalResponse returns the final response after the operations has finished
func (s *DeleteSecretPoller) FinalResponse(ctx context.Context) (DeleteSecretResponse, error) {
	return deleteSecretResponseFromGenerated(s.deleteResponse), nil
}

// PollUntilDone continually calls the Poll operation until the operation is completed. In between each
// Poll is a wait determined by the t parameter.
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

// BeginDeleteSecret deletes a secret from the keyvault. Delete cannot be applied to an individual version of a secret. This operation
// requires the secrets/delete permission. This response contains a Poller struct that can be used to Poll for a response, or the
// response PollUntilDone function can be used to poll until completion.
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

// GetDeletedSecret gets the specified deleted secret. The operation returns the deleted secret along with its attributes.
// This operation requires the secrets/get permission.
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

// UpdateSecretPropertiesOptions contains the optional parameters for the Client.UpdateSecretProperties method.
type UpdateSecretPropertiesOptions struct {
	// placeholder for future optional parameters
}

// UpdateSecretPropertiesResponse contains the underlying response object for the UpdateSecretProperties method
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

// UpdateSecretProperties updates the attributes associated with a specified secret in a given key vault. The update
// operation changes specified attributes of an existing stored secret, attributes that are not specified in the
// request are left unchanged. The value of a secret itself cannot be changed. This operation requires the secrets/set permission.
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
}

// convert generated response to the publicly exposed version.
func backupSecretResponseFromGenerated(i internal.KeyVaultClientBackupSecretResponse) BackupSecretResponse {
	return BackupSecretResponse{
		Value: i.Value,
	}
}

// BackupSecrets backs up the specified secret. Requests that a backup of the specified secret be downloaded to the client.
// All versions of the secret will be downloaded. This operation requires the secrets/backup permission.
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
	// placeholder for future response fields
}

// Converts the generated response to the publicly exposed version.
func purgeDeletedSecretResponseFromGenerated(i internal.KeyVaultClientPurgeDeletedSecretResponse) PurgeDeletedSecretResponse {
	return PurgeDeletedSecretResponse{}
}

// PurgeDeletedSecret deletes the specified secret. The purge deleted secret operation removes the secret permanently, without the possibility of recovery.
// This operation can only be enabled on a soft-delete enabled vault. This operation requires the secrets/purge permission.
func (c *Client) PurgeDeletedSecret(ctx context.Context, name string, options *PurgeDeletedSecretOptions) (PurgeDeletedSecretResponse, error) {
	if options == nil {
		options = &PurgeDeletedSecretOptions{}
	}
	resp, err := c.kvClient.PurgeDeletedSecret(ctx, c.vaultUrl, name, options.toGenerated())
	return purgeDeletedSecretResponseFromGenerated(resp), err
}

// RecoverDeletedSecretPoller is the poller returned by Client.BeginRecoverDeletedSecret
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

// Poll fetches the latest state of the LRO. It returns an HTTP response or error.
// If the LRO has completed successfully, the poller's state is updated and the HTTP response is returned.
// If the LRO has completed with failure or was cancelled, the poller's state is updated and the error is returned.
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

// FinalResponse returns the final response after the operations has finished
func (b *RecoverDeletedSecretPoller) FinalResponse(ctx context.Context) (RecoverDeletedSecretResponse, error) {
	return recoverDeletedSecretResponseFromGenerated(b.recoverResponse), nil
}

// PollUntilDone continually calls the Poll operation until the operation is completed. In between each
// Poll is a wait determined by the t parameter.
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

// BeginRecoverDeletedSecretOptions contains the optional parameters for the Client.BeginRecoverDeletedSecret operation
type BeginRecoverDeletedSecretOptions struct {
	// ResumeToken is a string to rehydrate a poller for an operation that has already begun.
	ResumeToken *string
}

// Convert the publicly exposed options object to the generated version
func (b BeginRecoverDeletedSecretOptions) toGenerated() *internal.KeyVaultClientRecoverDeletedSecretOptions {
	return &internal.KeyVaultClientRecoverDeletedSecretOptions{}
}

// RecoverDeletedSecretResponse is the response object for the Client.RecoverDeletedSecret operation.
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

// BeginRecoverDeletedSecret recovers the deleted secret in the specified vault to the latest version.
// This operation can only be performed on a soft-delete enabled vault. This operation requires the secrets/recover permission.
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

// ListDeletedSecretsResponse holds the data for a single page.
type ListDeletedSecretsResponse struct {
	// READ-ONLY; The URL to get the next set of deleted secrets.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// READ-ONLY; A response message containing a list of the deleted secrets in the vault along with a link to the next page of deleted secrets
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

// ListDeletedSecretsOptions contains the optional parameters for the Client.ListDeletedSecrets operation.
type ListDeletedSecretsOptions struct {
	// placeholder for future optional parameters
}

// ListDeletedSecrets lists all versions of the specified secret. The full secret identifier and attributes are provided
// in the response. No values are returned for the secrets. This operation requires the secrets/list permission.
func (c *Client) ListDeletedSecrets(options *ListDeletedSecretsOptions) *runtime.Pager[ListDeletedSecretsResponse] {
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

// ListSecretVersionsOptions contains the options for the ListSecretVersions operations
type ListSecretVersionsOptions struct {
	// placeholder for future optional parameters
}

// ListPropertiesOfSecretVersionsResponse contains response field for ListSecretVersionsPager.NextPage
type ListPropertiesOfSecretVersionsResponse struct {
	// READ-ONLY; The URL to get the next set of secrets.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// READ-ONLY; A response message containing a list of secrets in the key vault along with a link to the next page of secrets.
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

// ListPropertiesOfSecretVersions lists all versions of the specified secret. The full secret identifer and
// attributes are provided in the response. No values are returned for the secrets. This operation
// requires the secrets/list permission.
func (c *Client) ListPropertiesOfSecretVersions(name string, options *ListSecretVersionsOptions) *runtime.Pager[ListPropertiesOfSecretVersionsResponse] {
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

// ListSecretsOptions contains the options for the ListSecretVersions operations
type ListSecretsOptions struct {
	// placeholder for future optional parameters.
}

// ListPropertiesOfSecretsResponse contains the current page of results for the Client.ListSecrets operation.
type ListPropertiesOfSecretsResponse struct {
	// READ-ONLY; The URL to get the next set of secrets.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// READ-ONLY; A response message containing a list of secrets in the key vault along with a link to the next page of secrets.
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

// ListPropertiesOfSecrets list all secrets in a specified key vault. The ListPropertiesOfSecrets operation is applicable to the entire vault,
// however, only the base secret identifier and its attributes are provided in the response. Individual
// secret versions are not listed in the response. This operation requires the secrets/list permission.
func (c *Client) ListPropertiesOfSecrets(options *ListSecretsOptions) *runtime.Pager[ListPropertiesOfSecretsResponse] {
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
