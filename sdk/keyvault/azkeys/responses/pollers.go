//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package responses

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal/convert"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/models"
)

// DeleteKeyPoller is the interface for the Client.DeleteKey operation.
type DeleteKeyPoller struct {
	keyName        string // This is the key to Poll for in GetDeletedKey
	vaultUrl       string
	client         *generated.KeyVaultClient
	deleteResponse generated.KeyVaultClientDeleteKeyResponse
	lastResponse   generated.KeyVaultClientGetDeletedKeyResponse
	rawResponse    *http.Response
}

type NewDeleteKeyPollerParams struct {
	KeyName        string
	VaultUrl       string
	Client         *generated.KeyVaultClient
	DeleteResponse generated.KeyVaultClientDeleteKeyResponse
	LastResponse   generated.KeyVaultClientGetDeletedKeyResponse
}

func NewDeleteKeyPoller(params NewDeleteKeyPollerParams) *DeleteKeyPoller {
	return &DeleteKeyPoller{
		vaultUrl:       params.VaultUrl,
		keyName:        params.KeyName,
		client:         params.Client,
		deleteResponse: params.DeleteResponse,
		lastResponse:   params.LastResponse,
	}
}

// Done returns true if the LRO has reached a terminal state
func (s *DeleteKeyPoller) Done() bool {
	return s.lastResponse.RawResponse != nil
}

// Poll fetches the latest state of the LRO. It returns an HTTP response or error.(
// If the LRO has completed successfully, the poller's state is updated and the HTTP response is returned.
// If the LRO has completed with failure or was cancelled, the poller's state is updated and the error is returned.)
func (s *DeleteKeyPoller) Poll(ctx context.Context) (*http.Response, error) {
	resp, err := s.client.GetDeletedKey(ctx, s.vaultUrl, s.keyName, nil)
	if err == nil {
		// Service recognizes DeletedKey, operation is done
		s.lastResponse = resp
		return resp.RawResponse, nil
	}

	var httpResponseErr *azcore.ResponseError
	if errors.As(err, &httpResponseErr) {
		if httpResponseErr.StatusCode == http.StatusNotFound {
			// This is the expected result
			return s.deleteResponse.RawResponse, nil
		}
	}
	return s.deleteResponse.RawResponse, err
}

// FinalResponse returns the final response after the operations has finished
func (s *DeleteKeyPoller) FinalResponse(ctx context.Context) (DeleteKey, error) {
	return *deleteKeyResponseFromGenerated(&s.deleteResponse), nil
}

// pollUntilDone continually calls the Poll operation until the operation is completed. In between each
// Poll is a wait determined by the t parameter.
func (s *DeleteKeyPoller) pollUntilDone(ctx context.Context, t time.Duration) (DeleteKey, error) {
	for {
		resp, err := s.Poll(ctx)
		if err != nil {
			return DeleteKey{}, err
		}
		s.rawResponse = resp
		if s.Done() {
			break
		}
		time.Sleep(t)
	}
	return DeleteKey{}, nil
}

// RecoverDeletedKeyPoller is the interface for the Client.RecoverDeletedKey operation
type RecoverDeletedKeyPoller struct {
	keyName         string
	vaultUrl        string
	client          *generated.KeyVaultClient
	recoverResponse generated.KeyVaultClientRecoverDeletedKeyResponse
	lastResponse    generated.KeyVaultClientGetKeyResponse
	rawResponse     *http.Response
}

type NewRecoverDeletedKeyPollerParams struct {
	KeyName         string
	VaultUrl        string
	Client          *generated.KeyVaultClient
	RecoverResponse generated.KeyVaultClientRecoverDeletedKeyResponse
	LastResponse    generated.KeyVaultClientGetKeyResponse
}

func NewRecoverDeletedKeyPoller(params NewRecoverDeletedKeyPollerParams) *RecoverDeletedKeyPoller {
	return &RecoverDeletedKeyPoller{
		keyName:         params.KeyName,
		vaultUrl:        params.VaultUrl,
		client:          params.Client,
		recoverResponse: params.RecoverResponse,
		lastResponse:    params.LastResponse,
	}
}

// Done returns true when the polling operation is completed
func (b *RecoverDeletedKeyPoller) Done() bool {
	return b.rawResponse.StatusCode == http.StatusOK
}

// Poll fetches the latest state of the LRO. It returns an HTTP response or error.
// If the LRO has completed successfully, the poller's state is updated and the HTTP response is returned.
// If the LRO has completed with failure or was cancelled, the poller's state is updated and the error is returned.
func (b *RecoverDeletedKeyPoller) Poll(ctx context.Context) (*http.Response, error) {
	resp, err := b.client.GetKey(ctx, b.vaultUrl, b.keyName, "", nil)
	b.lastResponse = resp
	var httpErr *azcore.ResponseError
	if errors.As(err, &httpErr) {
		return httpErr.RawResponse, err
	}
	return resp.RawResponse, nil
}

// FinalResponse returns the final response after the operations has finished
func (b *RecoverDeletedKeyPoller) FinalResponse(ctx context.Context) (RecoverDeletedKey, error) {
	return recoverDeletedKeyResponseFromGenerated(b.recoverResponse), nil
}

// pollUntilDone is the method for the Response.PollUntilDone struct
func (b *RecoverDeletedKeyPoller) pollUntilDone(ctx context.Context, t time.Duration) (RecoverDeletedKey, error) {
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
	return recoverDeletedKeyResponseFromGenerated(b.recoverResponse), nil
}

// convert interal response to DeleteKeyResponse
func deleteKeyResponseFromGenerated(i *generated.KeyVaultClientDeleteKeyResponse) *DeleteKey {
	if i == nil {
		return nil
	}
	return &DeleteKey{
		RawResponse: i.RawResponse,
	}
}

// change recover deleted key reponse to the generated version.
func recoverDeletedKeyResponseFromGenerated(i generated.KeyVaultClientRecoverDeletedKeyResponse) RecoverDeletedKey {
	return RecoverDeletedKey{
		RawResponse: i.RawResponse,
		KeyBundle: models.KeyBundle{
			Attributes: convert.KeyAttributesFromGenerated(i.Attributes),
			Key:        convert.JSONWebKeyFromGenerated(i.Key),
			Tags:       convert.FromGeneratedMap(i.Tags),
			Managed:    i.Managed,
		},
	}
}
