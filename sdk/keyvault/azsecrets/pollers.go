//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsecrets

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets/internal/generated"
)

type beginDeleteSecretOperation struct {
	resp *http.Response
	poll func(context.Context) (*http.Response, error)
}

func (b *beginDeleteSecretOperation) Done() bool {
	return b.resp != nil && b.resp.StatusCode == http.StatusOK
}

func (b *beginDeleteSecretOperation) Poll(ctx context.Context) (*http.Response, error) {
	resp, err := b.poll(ctx)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNotFound) {
		return nil, runtime.NewResponseError(resp)
	}
	b.resp = resp
	return b.resp, nil
}

func (b *beginDeleteSecretOperation) Result(ctx context.Context, out *DeleteSecretResponse) error {
	var resp generated.KeyVaultClientDeleteSecretResponse
	data, err := runtime.Payload(b.resp)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}
	*out = deleteSecretResponseFromGenerated(resp)
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type beginRecoverDeletedSecretOperation struct {
	resp *http.Response
	poll func(context.Context) (*http.Response, error)
}

func (b *beginRecoverDeletedSecretOperation) Done() bool {
	return b.resp != nil && b.resp.StatusCode == http.StatusOK
}

func (b *beginRecoverDeletedSecretOperation) Poll(ctx context.Context) (*http.Response, error) {
	resp, err := b.poll(ctx)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNotFound) {
		return nil, runtime.NewResponseError(resp)
	}
	b.resp = resp
	return b.resp, nil
}

func (b *beginRecoverDeletedSecretOperation) Result(ctx context.Context, out *RecoverDeletedSecretResponse) error {
	var resp generated.KeyVaultClientRecoverDeletedSecretResponse
	data, err := runtime.Payload(b.resp)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}
	*out = recoverDeletedSecretResponseFromGenerated(resp)
	return nil
}
