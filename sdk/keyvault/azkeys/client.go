//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal"
)

// Client is the struct for interacting with a KeyVault Keys instance
type Client struct {
	kvClient *internal.KeyVaultClient
	vaultUrl string
}

// ClientOptions are the configurable options on a Client.
type ClientOptions struct {
	// Transport sets the transport for making HTTP requests.
	Transport policy.Transporter

	// Retry configures the built-in retry policy behavior.
	Retry policy.RetryOptions

	// Telemetry configures the built-in telemetry policy behavior.
	Telemetry policy.TelemetryOptions

	// Logging configures the built-in logging policy behavior.
	Logging policy.LogOptions

	// PerCallPolicies contains custom policies to inject into the pipeline.
	// Each policy is executed once per request.
	PerCallPolicies []policy.Policy

	// PerRetryPolicies contains custom policies to inject into the pipeline.
	// Each policy is executed once per request, and for each retry request.
	PerTryPolicies []policy.Policy
}

func (c *ClientOptions) toConnectionOptions() *internal.ConnectionOptions {
	if c == nil {
		return nil
	}

	return &internal.ConnectionOptions{
		HTTPClient:       c.Transport,
		Retry:            c.Retry,
		Telemetry:        c.Telemetry,
		Logging:          c.Logging,
		PerCallPolicies:  c.PerCallPolicies,
		PerRetryPolicies: c.PerTryPolicies,
	}
}

// NewClient returns a pointer to a Client object affinitized to a vaultUrl.
func NewClient(vaultUrl string, credential azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	conn := internal.NewConnection(credential, options.toConnectionOptions())

	return &Client{
		kvClient: &internal.KeyVaultClient{
			Con: conn,
		},
		vaultUrl: vaultUrl,
	}, nil
}

// CreateKeyOptions contains the optional parameters for the KeyVaultClient.CreateKey method.
type CreateKeyOptions struct {
	// Elliptic curve name. For valid values, see JsonWebKeyCurveName.
	Curve *internal.JSONWebKeyCurveName `json:"crv,omitempty"`

	// The attributes of a key managed by the key vault service.
	KeyAttributes *internal.KeyAttributes         `json:"attributes,omitempty"`
	KeyOps        []*internal.JSONWebKeyOperation `json:"key_ops,omitempty"`

	// The key size in bits. For example: 2048, 3072, or 4096 for RSA.
	KeySize *int32 `json:"key_size,omitempty"`

	// The public exponent for a RSA key.
	PublicExponent *int32 `json:"public_exponent,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]*string `json:"tags,omitempty"`
}

func (c *CreateKeyOptions) toGenerated() *internal.KeyVaultClientCreateKeyOptions {
	return &internal.KeyVaultClientCreateKeyOptions{}
}

func (c *CreateKeyOptions) toKeyCreateParameters(keyType KeyType) internal.KeyCreateParameters {
	return internal.KeyCreateParameters{
		Kty:            keyType.toGenerated(),
		Curve:          c.Curve,
		KeyAttributes:  c.KeyAttributes,
		KeyOps:         c.KeyOps,
		KeySize:        c.KeySize,
		PublicExponent: c.PublicExponent,
		Tags:           c.Tags,
	}
}

// KeyVaultClientCreateKeyResponse contains the response from method KeyVaultClient.CreateKey.
type CreateKeyResponse struct {
	KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func createKeyResponseFromGenerated(g internal.KeyVaultClientCreateKeyResponse) CreateKeyResponse {
	return CreateKeyResponse{
		RawResponse: g.RawResponse,
		KeyBundle: KeyBundle{
			Attributes: keyAttributesFromGenerated(g.Attributes),
			Key:        jsonWebKeyFromGenerated(g.Key),
			Tags:       g.Tags,
			Managed:    g.Managed,
		},
	}
}

func (c *Client) CreateKey(ctx context.Context, name string, keyType KeyType, options *CreateKeyOptions) (CreateKeyResponse, error) {
	if options == nil {
		options = &CreateKeyOptions{}
	}

	resp, err := c.kvClient.CreateKey(ctx, c.vaultUrl, name, options.toKeyCreateParameters(keyType), options.toGenerated())
	if err != nil {
		return CreateKeyResponse{}, err
	}

	return createKeyResponseFromGenerated(resp), nil
}

type CreateECKeyOptions struct {
	// Elliptic curve name. For valid values, see JsonWebKeyCurveName.
	CurveName *internal.JSONWebKeyCurveName `json:"crv,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]*string `json:"tags,omitempty"`

	// Whether to create an EC key with HSM protection
	HardwareProtected bool
}

func (c *CreateECKeyOptions) toKeyCreateParameters(keyType KeyType) internal.KeyCreateParameters {
	return internal.KeyCreateParameters{
		Kty:   keyType.toGenerated(),
		Curve: c.CurveName,
		Tags:  c.Tags,
	}
}

type CreateECKeyResponse struct {
	KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func createECKeyResponseFromGenerated(g internal.KeyVaultClientCreateKeyResponse) CreateECKeyResponse {
	return CreateECKeyResponse{
		RawResponse: g.RawResponse,
		KeyBundle: KeyBundle{
			Attributes: keyAttributesFromGenerated(g.Attributes),
			Key:        jsonWebKeyFromGenerated(g.Key),
			Tags:       g.Tags,
			Managed:    g.Managed,
		},
	}
}

func (c *Client) CreateECKey(ctx context.Context, name string, options *CreateECKeyOptions) (CreateECKeyResponse, error) {
	keyType := EC

	if options != nil && options.HardwareProtected {
		keyType = ECHSM
	} else if options == nil {
		options = &CreateECKeyOptions{}
	}

	resp, err := c.kvClient.CreateKey(ctx, c.vaultUrl, name, options.toKeyCreateParameters(keyType), &internal.KeyVaultClientCreateKeyOptions{})
	if err != nil {
		return CreateECKeyResponse{}, nil
	}

	return createECKeyResponseFromGenerated(resp), nil
}

type CreateOCTKeyOptions struct {
	// Hardware Protected OCT Key
	HardwareProtected bool

	// The key size in bits. For example: 2048, 3072, or 4096 for RSA.
	KeySize *int32 `json:"key_size,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]*string `json:"tags,omitempty"`
}

func (c *CreateOCTKeyOptions) toKeyCreateParameters(keyType KeyType) internal.KeyCreateParameters {
	return internal.KeyCreateParameters{
		Kty:     keyType.toGenerated(),
		KeySize: c.KeySize,
		Tags:    c.Tags,
	}
}

type CreateOCTKeyResponse struct {
	KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func createOCTKeyResponseFromGenerated(i internal.KeyVaultClientCreateKeyResponse) CreateOCTKeyResponse {
	return CreateOCTKeyResponse{
		RawResponse: i.RawResponse,
		KeyBundle: KeyBundle{
			Attributes: keyAttributesFromGenerated(i.Attributes),
			Key:        jsonWebKeyFromGenerated(i.Key),
			Tags:       i.Tags,
			Managed:    i.Managed,
		},
	}
}

func (c *Client) CreateOCTKey(ctx context.Context, name string, options *CreateOCTKeyOptions) (CreateOCTKeyResponse, error) {
	keyType := Oct

	if options != nil && options.HardwareProtected {
		keyType = OctHSM
	} else if options == nil {
		options = &CreateOCTKeyOptions{}
	}

	resp, err := c.kvClient.CreateKey(ctx, c.vaultUrl, name, options.toKeyCreateParameters(keyType), &internal.KeyVaultClientCreateKeyOptions{})

	return createOCTKeyResponseFromGenerated(resp), err
}

type CreateRSAKeyOptions struct {
	// Hardware Protected OCT Key
	HardwareProtected bool

	// The key size in bits. For example: 2048, 3072, or 4096 for RSA.
	KeySize *int32 `json:"key_size,omitempty"`

	// The public exponent for a RSA key.
	PublicExponent *int32 `json:"public_exponent,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]*string `json:"tags,omitempty"`
}

func (c CreateRSAKeyOptions) toKeyCreateParameters(k KeyType) internal.KeyCreateParameters {
	return internal.KeyCreateParameters{
		Kty:            k.toGenerated(),
		KeySize:        c.KeySize,
		PublicExponent: c.PublicExponent,
		Tags:           c.Tags,
	}
}

type CreateRSAKeyResponse struct {
	KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func createRSAKeyResponseFromGenerated(i internal.KeyVaultClientCreateKeyResponse) CreateRSAKeyResponse {
	return CreateRSAKeyResponse{
		RawResponse: i.RawResponse,
		KeyBundle: KeyBundle{
			Attributes: keyAttributesFromGenerated(i.Attributes),
			Key:        jsonWebKeyFromGenerated(i.Key),
			Tags:       i.Tags,
			Managed:    i.Managed,
		},
	}
}

func (c *Client) CreateRSAKey(ctx context.Context, name string, options *CreateRSAKeyOptions) (CreateRSAKeyResponse, error) {
	keyType := RSA

	if options != nil && options.HardwareProtected {
		keyType = RSAHSM
	} else if options == nil {
		options = &CreateRSAKeyOptions{}
	}

	resp, err := c.kvClient.CreateKey(ctx, c.vaultUrl, name, options.toKeyCreateParameters(keyType), &internal.KeyVaultClientCreateKeyOptions{})
	if err != nil {
		return CreateRSAKeyResponse{}, err
	}

	return createRSAKeyResponseFromGenerated(resp), nil
}

type ListKeysPager interface {
	// PageResponse returns the current ListKeysPage
	PageResponse() ListKeysPage

	// Err returns true if there is another page of data available, false if not
	Err() error

	// NextPage returns true if there is another page of data available, false if not
	NextPage(context.Context) bool
}

type listKeysPager struct {
	genPager internal.KeyVaultClientGetKeysPager
	err      error
}

func (l *listKeysPager) PageResponse() ListKeysPage {
	return listKeysPageFromGenerated(l.genPager.PageResponse())
}

func (l *listKeysPager) Err() error {
	return l.err
}

func (l *listKeysPager) NextPage(ctx context.Context) bool {
	return l.genPager.NextPage(ctx)
}

type ListKeysOptions struct {
	MaxResults *int32
}

func (l ListKeysOptions) toGenerated() *internal.KeyVaultClientGetKeysOptions {
	return &internal.KeyVaultClientGetKeysOptions{Maxresults: l.MaxResults}
}

type ListKeysPage struct {
	// READ-ONLY; The URL to get the next set of keys.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// READ-ONLY; A response message containing a list of keys in the key vault along with a link to the next page of keys.
	Keys []*KeyItem `json:"value,omitempty" azure:"ro"`
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func listKeysPageFromGenerated(i internal.KeyVaultClientGetKeysResponse) ListKeysPage {
	var keys []*KeyItem
	for _, k := range i.Value {
		keys = append(keys, keyItemFromGenerated(k))
	}
	return ListKeysPage{
		RawResponse: i.RawResponse,
		NextLink:    i.NextLink,
		Keys:        keys,
	}
}

func (c *Client) ListKeys(options *ListKeysOptions) ListKeysPager {
	if options == nil {
		options = &ListKeysOptions{}
	}
	p := c.kvClient.GetKeys(c.vaultUrl, options.toGenerated())

	return &listKeysPager{
		genPager: *p,
	}
}

type GetKeyOptions struct {
	Version string
}

type GetKeyResponse struct {
	KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func getKeyResponseFromGenerated(i internal.KeyVaultClientGetKeyResponse) GetKeyResponse {
	return GetKeyResponse{
		RawResponse: i.RawResponse,
		KeyBundle: KeyBundle{
			Attributes: keyAttributesFromGenerated(i.Attributes),
			Key:        jsonWebKeyFromGenerated(i.Key),
			Tags:       i.Tags,
			Managed:    i.Managed,
		},
	}
}

func (c *Client) GetKey(ctx context.Context, keyName string, options *GetKeyOptions) (GetKeyResponse, error) {
	if options == nil {
		options = &GetKeyOptions{}
	}

	resp, err := c.kvClient.GetKey(ctx, c.vaultUrl, keyName, options.Version, &internal.KeyVaultClientGetKeyOptions{})
	if err != nil {
		return GetKeyResponse{}, err
	}

	return getKeyResponseFromGenerated(resp), err
}

type GetDeletedKeyOptions struct{}

func (g GetDeletedKeyOptions) toGenerated() *internal.KeyVaultClientGetDeletedKeyOptions {
	return &internal.KeyVaultClientGetDeletedKeyOptions{}
}

type GetDeletedKeyResponse struct {
	DeletedKeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func getDeletedKeyResponseFromGenerated(i internal.KeyVaultClientGetDeletedKeyResponse) GetDeletedKeyResponse {
	return GetDeletedKeyResponse{
		RawResponse: i.RawResponse,
		DeletedKeyBundle: DeletedKeyBundle{
			KeyBundle: KeyBundle{
				Attributes: keyAttributesFromGenerated(i.Attributes),
				Key:        jsonWebKeyFromGenerated(i.Key),
				Tags:       i.Tags,
				Managed:    i.Managed,
			},
			RecoveryID:         i.RecoveryID,
			DeletedDate:        i.DeletedDate,
			ScheduledPurgeDate: i.ScheduledPurgeDate,
		},
	}
}

func (c *Client) GetDeletedKey(ctx context.Context, keyName string, options *GetDeletedKeyOptions) (GetDeletedKeyResponse, error) {
	if options == nil {
		options = &GetDeletedKeyOptions{}
	}

	resp, err := c.kvClient.GetDeletedKey(ctx, c.vaultUrl, keyName, options.toGenerated())
	if err != nil {
		return GetDeletedKeyResponse{}, err
	}

	return getDeletedKeyResponseFromGenerated(resp), nil
}

// PurgeDeletedKeyOptions is the struct for any future options for Client.PurgeDeletedKey.
type PurgeDeletedKeyOptions struct{}

func (p *PurgeDeletedKeyOptions) toGenerated() *internal.KeyVaultClientPurgeDeletedKeyOptions {
	return &internal.KeyVaultClientPurgeDeletedKeyOptions{}
}

// PurgeDeletedKeyResponse contains the response from method Client.PurgeDeletedKey.
type PurgeDeletedKeyResponse struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// Converts the generated response to the publicly exposed version.
func purgeDeletedKeyResponseFromGenerated(i internal.KeyVaultClientPurgeDeletedKeyResponse) PurgeDeletedKeyResponse {
	return PurgeDeletedKeyResponse{
		RawResponse: i.RawResponse,
	}
}

// PurgeDeletedKey deletes the specified key. The purge deleted key operation removes the key permanently, without the possibility of recovery.
// This operation can only be enabled on a soft-delete enabled vault. This operation requires the key/purge permission.
func (c *Client) PurgeDeletedKey(ctx context.Context, keyName string, options *PurgeDeletedKeyOptions) (PurgeDeletedKeyResponse, error) {
	if options == nil {
		options = &PurgeDeletedKeyOptions{}
	}
	resp, err := c.kvClient.PurgeDeletedKey(ctx, c.vaultUrl, keyName, options.toGenerated())
	return purgeDeletedKeyResponseFromGenerated(resp), err
}

// DeletedKeyResponse contains the response for a Client.BeginDeleteKey operation.
type DeleteKeyResponse struct {
	DeletedKeyBundle
	// RawResponse holds the underlying HTTP response
	RawResponse *http.Response
}

func deleteKeyResponseFromGenerated(i *internal.KeyVaultClientDeleteKeyResponse) *DeleteKeyResponse {
	if i == nil {
		return nil
	}
	return &DeleteKeyResponse{
		RawResponse: i.RawResponse,
	}
}

// BeginDeleteKeyOptions contains the optional parameters for the Client.BeginDeleteKey method.
type BeginDeleteKeyOptions struct{}

// convert public options to generated options struct
func (b *BeginDeleteKeyOptions) toGenerated() *internal.KeyVaultClientDeleteKeyOptions {
	return &internal.KeyVaultClientDeleteKeyOptions{}
}

// DeleteKeyPoller is the interface for the Client.DeleteKey operation.
type DeleteKeyPoller interface {
	// Done returns true if the LRO has reached a terminal state
	Done() bool

	// Poll fetches the latest state of the LRO. It returns an HTTP response or error.
	// If the LRO has completed successfully, the poller's state is updated and the HTTP response is returned.
	// If the LRO has completed with failure or was cancelled, the poller's state is updated and the error is returned.
	Poll(context.Context) (*http.Response, error)

	// FinalResponse returns the final response after the operations has finished
	FinalResponse(context.Context) (DeleteKeyResponse, error)
}

// The poller returned by the Client.StartDeleteKey operation
type startDeleteKeyPoller struct {
	keyName        string // This is the key to Poll for in GetDeletedKey
	vaultUrl       string
	client         *internal.KeyVaultClient
	deleteResponse internal.KeyVaultClientDeleteKeyResponse
	lastResponse   internal.KeyVaultClientGetDeletedKeyResponse
	RawResponse    *http.Response
}

// Done returns true if the LRO has reached a terminal state
func (s *startDeleteKeyPoller) Done() bool {
	return s.lastResponse.RawResponse != nil
}

// Poll fetches the latest state of the LRO. It returns an HTTP response or error.(
// If the LRO has completed successfully, the poller's state is updated and the HTTP response is returned.
// If the LRO has completed with failure or was cancelled, the poller's state is updated and the error is returned.)
func (s *startDeleteKeyPoller) Poll(ctx context.Context) (*http.Response, error) {
	resp, err := s.client.GetDeletedKey(ctx, s.vaultUrl, s.keyName, nil)
	if err == nil {
		// Service recognizes DeletedKey, operation is done
		s.lastResponse = resp
		return resp.RawResponse, nil
	} else if err != nil {
		return s.deleteResponse.RawResponse, nil
	}
	s.lastResponse = resp
	return resp.RawResponse, nil
}

// FinalResponse returns the final response after the operations has finished
func (s *startDeleteKeyPoller) FinalResponse(ctx context.Context) (DeleteKeyResponse, error) {
	return *deleteKeyResponseFromGenerated(&s.deleteResponse), nil
}

// pollUntilDone continually calls the Poll operation until the operation is completed. In between each
// Poll is a wait determined by the t parameter.
func (s *startDeleteKeyPoller) pollUntilDone(ctx context.Context, t time.Duration) (DeleteKeyResponse, error) {
	for {
		resp, err := s.Poll(ctx)
		if err != nil {
			return DeleteKeyResponse{}, err
		}
		s.RawResponse = resp
		if s.Done() {
			break
		}
		time.Sleep(t)
	}
	return DeleteKeyResponse{}, nil
}

type DeleteKeyPollerResponse struct {
	// PollUntilDone will poll the service endpoint until a terminal state is reached or an error occurs
	PollUntilDone func(context.Context, time.Duration) (DeleteKeyResponse, error)

	// Poller contains an initialized WidgetPoller
	Poller DeleteKeyPoller

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// BeginDeleteKey deletes a key from the keyvault. Delete cannot be applied to an individual version of a key. This operation
// requires the key/delete permission. This response contains a Poller struct that can be used to Poll for a response, or the
// response PollUntilDone function can be used to poll until completion.
func (c *Client) BeginDeleteKey(ctx context.Context, keyName string, options *BeginDeleteKeyOptions) (DeleteKeyPollerResponse, error) {
	if options == nil {
		options = &BeginDeleteKeyOptions{}
	}
	resp, err := c.kvClient.DeleteKey(ctx, c.vaultUrl, keyName, options.toGenerated())
	if err != nil {
		return DeleteKeyPollerResponse{}, err
	}

	getResp, err := c.kvClient.GetDeletedKey(ctx, c.vaultUrl, keyName, nil)
	var httpErr azcore.HTTPResponse
	if errors.As(err, &httpErr) {
		if httpErr.RawResponse().StatusCode != http.StatusNotFound {
			return DeleteKeyPollerResponse{}, err
		}
	}

	s := &startDeleteKeyPoller{
		vaultUrl:       c.vaultUrl,
		keyName:        keyName,
		client:         c.kvClient,
		deleteResponse: resp,
		lastResponse:   getResp,
	}

	return DeleteKeyPollerResponse{
		Poller:        s,
		RawResponse:   resp.RawResponse,
		PollUntilDone: s.pollUntilDone,
	}, nil
}

type BackupKeyOptions struct{}

func (b BackupKeyOptions) toGenerated() *internal.KeyVaultClientBackupKeyOptions {
	return &internal.KeyVaultClientBackupKeyOptions{}
}

type BackupKeyResponse struct {
	// READ-ONLY; The backup blob containing the backed up key.
	Value []byte `json:"value,omitempty" azure:"ro"`
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func backupKeyResponseFromGenerated(i internal.KeyVaultClientBackupKeyResponse) BackupKeyResponse {
	return BackupKeyResponse{
		RawResponse: i.RawResponse,
		Value:       i.Value,
	}
}

func (c *Client) BackupKey(ctx context.Context, keyName string, options *BackupKeyOptions) (BackupKeyResponse, error) {
	if options == nil {
		options = &BackupKeyOptions{}
	}

	resp, err := c.kvClient.BackupKey(ctx, c.vaultUrl, keyName, options.toGenerated())
	if err != nil {
		return BackupKeyResponse{}, err
	}

	return backupKeyResponseFromGenerated(resp), nil
}

// RecoverDeletedKeyPoller is the interface for the Client.RecoverDeletedKey operation
type RecoverDeletedKeyPoller interface {
	// Done returns true if the LRO has reached a terminal state
	Done() bool

	// Poll fetches the latest state of the LRO. It returns an HTTP response or error.
	// If the LRO has completed successfully, the poller's state is updated and the HTTP response is returned.
	// If the LRO has completed with failure or was cancelled, the poller's state is updated and the error is returned.
	Poll(context.Context) (*http.Response, error)

	// FinalResponse returns the final response after the operations has finished
	FinalResponse(context.Context) (RecoverDeletedKeyResponse, error)
}

type beginRecoverPoller struct {
	keyName         string
	vaultUrl        string
	client          *internal.KeyVaultClient
	recoverResponse internal.KeyVaultClientRecoverDeletedKeyResponse
	lastResponse    internal.KeyVaultClientGetKeyResponse
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
	resp, err := b.client.GetKey(ctx, b.vaultUrl, b.keyName, "", nil)
	b.lastResponse = resp
	var httpErr azcore.HTTPResponse
	if errors.As(err, &httpErr) {
		return httpErr.RawResponse(), err
	}
	return resp.RawResponse, nil
}

// FinalResponse returns the final response after the operations has finished
func (b *beginRecoverPoller) FinalResponse(ctx context.Context) (RecoverDeletedKeyResponse, error) {
	return recoverDeletedKeyResponseFromGenerated(b.recoverResponse), nil
}

func (b *beginRecoverPoller) pollUntilDone(ctx context.Context, t time.Duration) (RecoverDeletedKeyResponse, error) {
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
	return recoverDeletedKeyResponseFromGenerated(b.recoverResponse), nil
}

// BeginRecoverDeletedKeyOptions contains the optional parameters for the Client.BeginRecoverDeletedKey operation
type BeginRecoverDeletedKeyOptions struct{}

// Convert the publicly exposed options object to the generated version
func (b BeginRecoverDeletedKeyOptions) toGenerated() *internal.KeyVaultClientRecoverDeletedKeyOptions {
	return &internal.KeyVaultClientRecoverDeletedKeyOptions{}
}

// RecoverDeletedKeyResponse is the response object for the Client.RecoverDeletedKey operation.
type RecoverDeletedKeyResponse struct {
	KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// change recover deleted key reponse to the generated version.
func recoverDeletedKeyResponseFromGenerated(i internal.KeyVaultClientRecoverDeletedKeyResponse) RecoverDeletedKeyResponse {
	return RecoverDeletedKeyResponse{
		RawResponse: i.RawResponse,
		KeyBundle: KeyBundle{
			Attributes: keyAttributesFromGenerated(i.Attributes),
			Key:        jsonWebKeyFromGenerated(i.Key),
			Tags:       i.Tags,
			Managed:    i.Managed,
		},
	}
}

// RecoverDeletedKeyPollerResponse contains the response of the Client.BeginRecoverDeletedKey operations
type RecoverDeletedKeyPollerResponse struct {
	// PollUntilDone will poll the service endpoint until a terminal state is reached or an error occurs
	PollUntilDone func(context.Context, time.Duration) (RecoverDeletedKeyResponse, error)

	// Poller contains an initialized RecoverDeletedKeyPoller
	Poller RecoverDeletedKeyPoller

	// RawResponse cotains the underlying HTTP response
	RawResponse *http.Response
}

// BeginRecoverDeletedKey recovers the deleted key in the specified vault to the latest version.
// This operation can only be performed on a soft-delete enabled vault. This operation requires the keys/recover permission.
func (c *Client) BeginRecoverDeletedKey(ctx context.Context, keyName string, options *BeginRecoverDeletedKeyOptions) (RecoverDeletedKeyPollerResponse, error) {
	if options == nil {
		options = &BeginRecoverDeletedKeyOptions{}
	}
	resp, err := c.kvClient.RecoverDeletedKey(ctx, c.vaultUrl, keyName, options.toGenerated())
	if err != nil {
		return RecoverDeletedKeyPollerResponse{}, err
	}

	getResp, err := c.kvClient.GetKey(ctx, c.vaultUrl, keyName, "", nil)
	var httpErr azcore.HTTPResponse
	if errors.As(err, &httpErr) {
		if httpErr.RawResponse().StatusCode != http.StatusNotFound {
			return RecoverDeletedKeyPollerResponse{}, err
		}
	}

	b := &beginRecoverPoller{
		lastResponse:    getResp,
		keyName:         keyName,
		client:          c.kvClient,
		vaultUrl:        c.vaultUrl,
		recoverResponse: resp,
		RawResponse:     getResp.RawResponse,
	}

	return RecoverDeletedKeyPollerResponse{
		PollUntilDone: b.pollUntilDone,
		Poller:        b,
		RawResponse:   getResp.RawResponse,
	}, nil
}

type UpdateKeyPropertiesOptions struct {
	Version string
	// The attributes of a key managed by the key vault service.
	KeyAttributes *KeyAttributes `json:"attributes,omitempty"`

	// Json web key operations. For more information on possible key operations, see JsonWebKeyOperation.
	KeyOps []*internal.JSONWebKeyOperation `json:"key_ops,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]*string `json:"tags,omitempty"`
}

func (u UpdateKeyPropertiesOptions) toKeyUpdateParameters() internal.KeyUpdateParameters {
	var attribs *internal.KeyAttributes
	if u.KeyAttributes != nil {
		attribs = u.KeyAttributes.toGenerated()
	}
	return internal.KeyUpdateParameters{
		KeyOps:        u.KeyOps,
		KeyAttributes: attribs,
		Tags:          u.Tags,
	}
}

func (u UpdateKeyPropertiesOptions) toGeneratedOptions() *internal.KeyVaultClientUpdateKeyOptions {
	return &internal.KeyVaultClientUpdateKeyOptions{}
}

type UpdateKeyPropertiesResponse struct {
	KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func updateKeyPropertiesFromGenerated(i internal.KeyVaultClientUpdateKeyResponse) UpdateKeyPropertiesResponse {
	return UpdateKeyPropertiesResponse{
		RawResponse: i.RawResponse,
		KeyBundle: KeyBundle{
			Attributes: keyAttributesFromGenerated(i.Attributes),
			Key:        jsonWebKeyFromGenerated(i.Key),
			Tags:       i.Tags,
			Managed:    i.Managed,
		},
	}
}

func (c *Client) UpdateKeyProperties(ctx context.Context, keyName string, options *UpdateKeyPropertiesOptions) (UpdateKeyPropertiesResponse, error) {
	if options == nil {
		options = &UpdateKeyPropertiesOptions{}
	}
	resp, err := c.kvClient.UpdateKey(
		ctx,
		c.vaultUrl,
		keyName,
		options.Version,
		options.toKeyUpdateParameters(),
		options.toGeneratedOptions(),
	)
	if err != nil {
		return UpdateKeyPropertiesResponse{}, err
	}

	return updateKeyPropertiesFromGenerated(resp), nil
}

// ListDeletedKeys is the interface for the Client.ListDeletedKeys operation
type ListDeletedKeysPager interface {
	// PageResponse returns the current ListDeletedKeysPage
	PageResponse() ListDeletedKeysPage

	// Err returns true if there is another page of data available, false if not
	Err() error

	// NextPage returns true if there is another page of data available, false if not
	NextPage(context.Context) bool
}

// listDeletedKeysPager is the pager returned by Client.ListDeletedKeys
type listDeletedKeysPager struct {
	genPager *internal.KeyVaultClientGetDeletedKeysPager
}

// PageResponse returns the current page of results
func (l *listDeletedKeysPager) PageResponse() ListDeletedKeysPage {
	resp := l.genPager.PageResponse()

	var values []*DeletedKeyItem
	for _, d := range resp.Value {
		values = append(values, deletedKeyItemFromGenerated(d))
	}

	return ListDeletedKeysPage{
		RawResponse: resp.RawResponse,
		NextLink:    resp.NextLink,
		DeletedKeys: values,
	}
}

// Err returns an error if the last operation resulted in an error.
func (l *listDeletedKeysPager) Err() error {
	return l.genPager.Err()
}

// NextPage fetches the next page of results.
func (l *listDeletedKeysPager) NextPage(ctx context.Context) bool {
	return l.genPager.NextPage(ctx)
}

// ListDeletedKeysPage holds the data for a single page.
type ListDeletedKeysPage struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response

	// READ-ONLY; The URL to get the next set of deleted keys.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// READ-ONLY; A response message containing a list of the deleted keys in the vault along with a link to the next page of deleted keys
	DeletedKeys []*DeletedKeyItem `json:"value,omitempty" azure:"ro"`
}

// ListDeletedKeysOptions contains the optional parameters for the Client.ListDeletedKeys operation.
type ListDeletedKeysOptions struct {
	// Maximum number of results to return in a page. If not specified the service will return up to 25 results.
	MaxResults *int32
}

// Convert publicly exposed options to the generated version.a
func (l *ListDeletedKeysOptions) toGenerated() *internal.KeyVaultClientGetDeletedKeysOptions {
	return &internal.KeyVaultClientGetDeletedKeysOptions{
		Maxresults: l.MaxResults,
	}
}

// ListDeletedKeys lists all versions of the specified key. The full key identifier and attributes are provided
// in the response. No values are returned for the keys. This operation requires the keys/list permission.
func (c *Client) ListDeletedKeys(options *ListDeletedKeysOptions) ListDeletedKeysPager {
	if options == nil {
		options = &ListDeletedKeysOptions{}
	}

	return &listDeletedKeysPager{
		genPager: c.kvClient.GetDeletedKeys(c.vaultUrl, options.toGenerated()),
	}

}
